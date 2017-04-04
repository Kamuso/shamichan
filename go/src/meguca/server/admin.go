// Various administration endpoints for logged in users

package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"meguca/auth"
	"meguca/common"
	"meguca/config"
	"meguca/db"
	"meguca/server/websockets"
	"net/http"
	"regexp"
	"time"
)

const (
	// Body size limit for POST request JSON. Should never exceed 32 KB.
	// Consider anything bigger an attack.
	jsonLimit = 1 << 15

	maxAnswers      = 100  // Maximum number of eightball answers
	maxEightballLen = 2000 // Total chars in eightball
)

var (
	errTooManyAnswers   = errors.New("too many eightball answers")
	errEightballTooLong = common.ErrTooLong("eightball")
	errTitleTooLong     = common.ErrTooLong("board title")
	errNoticeTooLong    = common.ErrTooLong("notice")
	errRulesTooLong     = common.ErrTooLong("rules")
	errBanReasonTooLong = common.ErrTooLong("ban reason")
	errInvalidBoardName = errors.New("invalid board name")
	errBoardNameTaken   = errors.New("board name taken")
	errAccessDenied     = errors.New("access denied")
	errNoReason         = errors.New("no reason provided")
	errNoDuration       = errors.New("no ban duration provided")

	boardNameValidation = regexp.MustCompile(`^[a-z0-9]{1,3}$`)
)

// Request to set the board-specific configuration for a board
type boardConfigSettingRequest struct {
	auth.SessionCreds
	config.BoardConfigs
}

// Request for the current non-public board configuration
type boardConfigRequest struct {
	auth.SessionCreds
	ID string `json:"id"`
}

type configSettingRequest struct {
	auth.SessionCreds
	config.Configs
}

type boardCreationRequest struct {
	Name, Title, Captcha string
	auth.SessionCreds
}

type boardDeletionRequest struct {
	ID, Captcha string
	auth.SessionCreds
}

// Request to perform a moderation action on a specific set of posts
type postActionRequest struct {
	IDs []uint64
	auth.SessionCreds
}

// Decode JSON sent in a request with a read limit of 8 KB. Returns if the
// decoding succeeded.
func decodeJSON(w http.ResponseWriter, r *http.Request, dest interface{}) bool {
	decoder := json.NewDecoder(io.LimitReader(r.Body, jsonLimit))
	if err := decoder.Decode(dest); err != nil {
		http.Error(w, fmt.Sprintf("400 %s", err), 400)
		logError(r, err)
		return false
	}
	return true
}

// Set board-specific configurations to the user's owned board
func configureBoard(w http.ResponseWriter, r *http.Request) {
	var msg boardConfigSettingRequest
	isValid := decodeJSON(w, r, &msg) &&
		isLoggedIn(w, r, msg.SessionCreds) &&
		isBoardOwner(w, r, msg.ID, msg.UserID) &&
		validateBoardConfigs(w, msg.BoardConfigs)
	if !isValid {
		return
	}

	msg.BoardConfigs.ID = msg.ID
	if err := db.UpdateBoard(msg.BoardConfigs); err != nil {
		text500(w, r, err)
		return
	}
}

// Assert the user is one of the board's owners
func isBoardOwner(
	w http.ResponseWriter,
	r *http.Request,
	board, userID string,
) bool {
	pos, err := db.GetPositions(userID, board)
	switch {
	case err != nil:
		text500(w, r, err)
		return false
	case !pos["owners"]:
		http.Error(w, "403 Not board owner", 403)
		return false
	default:
		return true
	}
}

// Validate length limit compliance of various fields
func validateBoardConfigs(
	w http.ResponseWriter,
	conf config.BoardConfigs,
) bool {
	totalLen := 0
	for _, answer := range conf.Eightball {
		totalLen += len(answer)
	}

	var err error
	switch {
	case len(conf.Eightball) > maxAnswers:
		err = errTooManyAnswers
	case totalLen > maxEightballLen:
		err = errEightballTooLong
	case len(conf.Notice) > common.MaxLenNotice:
		err = errNoticeTooLong
	case len(conf.Rules) > common.MaxLenRules:
		err = errRulesTooLong
	case len(conf.Title) > common.MaxLenBoardTitle:
		err = errTitleTooLong
	}
	if err != nil {
		http.Error(w, fmt.Sprintf("400 %s", err), 400)
		return false
	}

	return true
}

// Serve the current board configurations to the client, including publically
// unexposed ones. Intended to be used before setting the the configs with
// configureBoard().
func servePrivateBoardConfigs(w http.ResponseWriter, r *http.Request) {
	conf, isValid := boardConfData(w, r)
	if !isValid {
		return
	}
	serveJSON(w, r, "", conf)
}

// Serve the current server configurations. Available only to the "admin"
// account
func servePrivateServerConfigs(w http.ResponseWriter, r *http.Request) {
	var msg auth.SessionCreds
	if !decodeJSON(w, r, &msg) || !isAdmin(w, r, msg) {
		return
	}
	serveJSON(w, r, "", config.Get())
}

func isAdmin(
	w http.ResponseWriter,
	r *http.Request,
	msg auth.SessionCreds,
) bool {
	if !(isLoggedIn(w, r, msg)) {
		return false
	}
	if msg.UserID != "admin" {
		text403(w, errAccessDenied)
		return false
	}
	return true
}

// Determine, if the client has access rights to the configurations, and return
// them, if so
func boardConfData(w http.ResponseWriter, r *http.Request) (
	config.BoardConfigs, bool,
) {
	var (
		msg  boardConfigRequest
		conf config.BoardConfigs
	)
	isValid := decodeJSON(w, r, &msg) &&
		isLoggedIn(w, r, msg.SessionCreds) &&
		isBoardOwner(w, r, msg.ID, msg.UserID)
	if !isValid {
		return conf, false
	}

	conf = config.GetBoardConfigs(msg.ID).BoardConfigs
	if conf.ID == "" {
		text404(w)
		return conf, false
	}

	return conf, true
}

// Handle requests to create a board
func createBoard(w http.ResponseWriter, r *http.Request) {
	var msg boardCreationRequest
	if !decodeJSON(w, r, &msg) || !isLoggedIn(w, r, msg.SessionCreds) {
		return
	}

	// Validate request data
	var err error
	switch {
	case msg.UserID != "admin" && config.Get().DisableUserBoards:
		err = errAccessDenied
	case !boardNameValidation.MatchString(msg.Name):
		err = errInvalidBoardName
	case len(msg.Title) > 100:
		err = errTitleTooLong
	case !auth.AuthenticateCaptcha(msg.Captcha):
		err = errInvalidCaptcha
	}
	if err != nil {
		text400(w, err)
		return
	}

	tx, err := db.StartTransaction()
	if err != nil {
		text500(w, r, err)
		return
	}
	defer db.RollbackOnError(tx, &err)

	err = db.WriteBoard(tx, db.BoardConfigs{
		Created: time.Now(),
		BoardConfigs: config.BoardConfigs{
			BoardPublic: config.BoardPublic{
				Title: msg.Title,
			},
			ID:        msg.Name,
			Eightball: config.EightballDefaults,
		},
	})
	switch {
	case err == nil:
	case db.IsConflictError(err):
		text400(w, errBoardNameTaken)
		return
	default:
		text500(w, r, err)
		return
	}

	err = db.WriteStaff(tx, msg.Name, map[string][]string{
		"owners": []string{msg.UserID},
	})
	if err != nil {
		text500(w, r, err)
		return
	}
	if err := tx.Commit(); err != nil {
		text500(w, r, err)
	}
}

// Set the server configuration to match the one sent from the admin account
// user
func configureServer(w http.ResponseWriter, r *http.Request) {
	var msg configSettingRequest
	if !decodeJSON(w, r, &msg) || !isAdmin(w, r, msg.SessionCreds) {
		return
	}
	if err := db.WriteConfigs(msg.Configs); err != nil {
		text500(w, r, err)
	}
}

// Delete a board owned by the client
func deleteBoard(w http.ResponseWriter, r *http.Request) {
	var msg boardDeletionRequest
	isValid := decodeJSON(w, r, &msg) &&
		isLoggedIn(w, r, msg.SessionCreds) &&
		isBoardOwner(w, r, msg.ID, msg.UserID)
	if !isValid {
		return
	}
	if !auth.AuthenticateCaptcha(msg.Captcha) {
		text403(w, errInvalidCaptcha)
		return
	}
	if err := db.DeleteBoard(msg.ID); err != nil {
		text500(w, r, err)
		return
	}
}

// Delete one or multiple posts on a moderated board
func deletePost(w http.ResponseWriter, r *http.Request) {
	var msg postActionRequest
	if !decodeJSON(w, r, &msg) || !isLoggedIn(w, r, msg.SessionCreds) {
		return
	}
	for _, id := range msg.IDs {
		board, err := db.GetPostBoard(id)
		switch err {
		case nil:
		case sql.ErrNoRows:
			text400(w, err)
			return
		default:
			text500(w, r, err)
			return
		}

		// TODO: More than just board owners should be able to delete posts
		if !isBoardOwner(w, r, board, msg.UserID) {
			return
		}

		err = db.DeletePost(board, id)
		switch err.(type) {
		case nil:
		case common.ErrInvalidPostID:
			text400(w, err)
			return
		default:
			text500(w, r, err)
			return
		}
	}
}

// Ban a specific IP from a specific board
func ban(w http.ResponseWriter, r *http.Request) {
	var msg struct {
		Global   bool
		Duration uint64
		Reason   string
		postActionRequest
	}

	// Decode and validate
	isValid := decodeJSON(w, r, &msg) &&
		isLoggedIn(w, r, msg.SessionCreds)
	if isValid && msg.Global {
		isValid = isAdmin(w, r, msg.SessionCreds)
	}
	switch {
	case !isValid:
		return
	case len(msg.Reason) > common.MaxBanReasonLength:
		text400(w, errBanReasonTooLong)
		return
	case msg.Reason == "":
		text400(w, errNoReason)
		return
	case msg.Duration == 0:
		text400(w, errNoDuration)
		return
	}

	// Group posts by board
	byBoard := make(map[string][]uint64, 2)
	if msg.Global {
		byBoard["all"] = msg.IDs
	} else {
		for _, id := range msg.IDs {
			board, err := db.GetPostBoard(id)
			switch err {
			case nil:
			case sql.ErrNoRows:
				text400(w, err)
				return
			default:
				text500(w, r, err)
				return
			}

			byBoard[board] = append(byBoard[board], id)
		}

		// Assert rights to moderate for all affected boards
		for b := range byBoard {
			if !isBoardOwner(w, r, b, msg.UserID) {
				return
			}
		}
	}

	// Apply bans
	expires := time.Now().Add(time.Duration(msg.Duration) * time.Minute)
	for b, ids := range byBoard {
		ips, err := db.Ban(b, msg.Reason, msg.UserID, expires, ids...)
		if err != nil {
			text500(w, r, err)
			return
		}

		// Redirect all banned connected clients to the /all/ board
		for ip := range ips {
			for _, cl := range common.Clients.GetByIP(ip) {
				cl.Redirect("all")
			}
		}
	}
}

// Send a textual message to all connected clients
func sendNotification(w http.ResponseWriter, r *http.Request) {
	var msg struct {
		Text string
		auth.SessionCreds
	}
	if !decodeJSON(w, r, &msg) || !isAdmin(w, r, msg.SessionCreds) {
		return
	}

	data, err := common.EncodeMessage(common.MessageNotification, msg.Text)
	if err != nil {
		text500(w, r, err)
		return
	}
	for _, cl := range websockets.Clients.All() {
		cl.Send(data)
	}
}

// Assign moderation staff to a board
func assignStaff(w http.ResponseWriter, r *http.Request) {
	var msg struct {
		Board, Captcha               string
		Owners, Moderators, Janitors []string
		auth.SessionCreds
	}

	isValid := decodeJSON(w, r, &msg) &&
		isLoggedIn(w, r, msg.SessionCreds) &&
		isBoardOwner(w, r, msg.Board, msg.UserID)
	switch {
	case !isValid:
		return
	// Ensure there always is at least one board owner
	case len(msg.Owners) == 0:
		text400(w, errors.New("no board owners set"))
		return
	case !auth.AuthenticateCaptcha(msg.Captcha):
		text400(w, errInvalidCaptcha)
		return
	default:
		// Maximum of 100 staff per position
		for _, s := range [...][]string{msg.Owners, msg.Moderators, msg.Janitors} {
			if len(s) > 100 {
				text400(w, errors.New("too many staff per position"))
				return
			}
		}
	}

	// Write to database
	tx, err := db.StartTransaction()
	if err != nil {
		text500(w, r, err)
		return
	}
	defer db.RollbackOnError(tx, &err)

	err = db.WriteStaff(tx, msg.Board, map[string][]string{
		"owners":     msg.Owners,
		"moderators": msg.Moderators,
		"janitors":   msg.Janitors,
	})
	if err != nil {
		text500(w, r, err)
		return
	}

	err = tx.Commit()
	if err != nil {
		text500(w, r, err)
	}
}
