// This file is automatically generated by qtc from "forms.html".
// See https://github.com/valyala/quicktemplate for details.

//line forms.html:1
package templates

//line forms.html:1
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line forms.html:1
import "github.com/bakape/meguca/config"

//line forms.html:2
import "github.com/bakape/meguca/lang"

// OwnedBoard renders a form for selecting one of several boards owned by the
// user

//line forms.html:6
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line forms.html:6
func StreamOwnedBoard(qw422016 *qt422016.Writer, boards config.BoardTitles, lang map[string]string) {
	//line forms.html:7
	if len(boards) != 0 {
		//line forms.html:7
		qw422016.N().S(`<select name="boards" required>`)
		//line forms.html:9
		for _, b := range boards {
			//line forms.html:9
			qw422016.N().S(`<option value="`)
			//line forms.html:10
			qw422016.N().S(b.ID)
			//line forms.html:10
			qw422016.N().S(`">`)
			//line forms.html:11
			streamformatTitle(qw422016, b.ID, b.Title)
			//line forms.html:11
			qw422016.N().S(`</option>`)
			//line forms.html:13
		}
		//line forms.html:13
		qw422016.N().S(`</select><br><input type="submit" value="`)
		//line forms.html:16
		qw422016.N().S(lang["submit"])
		//line forms.html:16
		qw422016.N().S(`">`)
		//line forms.html:17
	} else {
		//line forms.html:18
		qw422016.N().S(lang["ownNoBoards"])
		//line forms.html:18
		qw422016.N().S(`<br><br>`)
		//line forms.html:21
	}
	//line forms.html:21
	qw422016.N().S(`<input type="button" name="cancel" value="`)
	//line forms.html:22
	qw422016.N().S(lang["cancel"])
	//line forms.html:22
	qw422016.N().S(`"><div class="form-response admin"></div>`)
//line forms.html:24
}

//line forms.html:24
func WriteOwnedBoard(qq422016 qtio422016.Writer, boards config.BoardTitles, lang map[string]string) {
	//line forms.html:24
	qw422016 := qt422016.AcquireWriter(qq422016)
	//line forms.html:24
	StreamOwnedBoard(qw422016, boards, lang)
	//line forms.html:24
	qt422016.ReleaseWriter(qw422016)
//line forms.html:24
}

//line forms.html:24
func OwnedBoard(boards config.BoardTitles, lang map[string]string) string {
	//line forms.html:24
	qb422016 := qt422016.AcquireByteBuffer()
	//line forms.html:24
	WriteOwnedBoard(qb422016, boards, lang)
	//line forms.html:24
	qs422016 := string(qb422016.B)
	//line forms.html:24
	qt422016.ReleaseByteBuffer(qb422016)
	//line forms.html:24
	return qs422016
//line forms.html:24
}

//line forms.html:26
func streamformatTitle(qw422016 *qt422016.Writer, id, title string) {
	//line forms.html:26
	qw422016.N().S(`/`)
	//line forms.html:27
	qw422016.N().S(id)
	//line forms.html:27
	qw422016.N().S(`/ -`)
	//line forms.html:27
	qw422016.E().S(title)
//line forms.html:28
}

//line forms.html:28
func writeformatTitle(qq422016 qtio422016.Writer, id, title string) {
	//line forms.html:28
	qw422016 := qt422016.AcquireWriter(qq422016)
	//line forms.html:28
	streamformatTitle(qw422016, id, title)
	//line forms.html:28
	qt422016.ReleaseWriter(qw422016)
//line forms.html:28
}

//line forms.html:28
func formatTitle(id, title string) string {
	//line forms.html:28
	qb422016 := qt422016.AcquireByteBuffer()
	//line forms.html:28
	writeformatTitle(qb422016, id, title)
	//line forms.html:28
	qs422016 := string(qb422016.B)
	//line forms.html:28
	qt422016.ReleaseByteBuffer(qb422016)
	//line forms.html:28
	return qs422016
//line forms.html:28
}

// BoardNavigation renders a board selection and search form

//line forms.html:31
func StreamBoardNavigation(qw422016 *qt422016.Writer, lang map[string]string) {
	//line forms.html:31
	qw422016.N().S(`<input type="text" class="full-width" name="search" placeholder="`)
	//line forms.html:32
	qw422016.N().S(lang["search"])
	//line forms.html:32
	qw422016.N().S(`"><br><form><input type="submit" value="`)
	//line forms.html:35
	qw422016.N().S(lang["apply"])
	//line forms.html:35
	qw422016.N().S(`"><input type="button" name="cancel" value="`)
	//line forms.html:36
	qw422016.N().S(lang["cancel"])
	//line forms.html:36
	qw422016.N().S(`"><br>`)
	//line forms.html:38
	for _, b := range config.GetBoardTitles() {
		//line forms.html:38
		qw422016.N().S(`<label><input type="checkbox" name="`)
		//line forms.html:40
		qw422016.N().S(b.ID)
		//line forms.html:40
		qw422016.N().S(`"><a href="/`)
		//line forms.html:41
		qw422016.N().S(b.ID)
		//line forms.html:41
		qw422016.N().S(`/" class="history">`)
		//line forms.html:42
		streamformatTitle(qw422016, b.ID, b.Title)
		//line forms.html:42
		qw422016.N().S(`</a><br></label>`)
		//line forms.html:46
	}
	//line forms.html:46
	qw422016.N().S(`</form>`)
//line forms.html:48
}

//line forms.html:48
func WriteBoardNavigation(qq422016 qtio422016.Writer, lang map[string]string) {
	//line forms.html:48
	qw422016 := qt422016.AcquireWriter(qq422016)
	//line forms.html:48
	StreamBoardNavigation(qw422016, lang)
	//line forms.html:48
	qt422016.ReleaseWriter(qw422016)
//line forms.html:48
}

//line forms.html:48
func BoardNavigation(lang map[string]string) string {
	//line forms.html:48
	qb422016 := qt422016.AcquireByteBuffer()
	//line forms.html:48
	WriteBoardNavigation(qb422016, lang)
	//line forms.html:48
	qs422016 := string(qb422016.B)
	//line forms.html:48
	qt422016.ReleaseByteBuffer(qb422016)
	//line forms.html:48
	return qs422016
//line forms.html:48
}

// CreateBoard renders a the form for creating new boards

//line forms.html:51
func StreamCreateBoard(qw422016 *qt422016.Writer, ln lang.Pack) {
	//line forms.html:52
	qw422016.N().S(renderTable(specs["createBoard"], ln))
	//line forms.html:52
	qw422016.N().S(`<input type="submit" value="`)
	//line forms.html:53
	qw422016.N().S(ln.UI["submit"])
	//line forms.html:53
	qw422016.N().S(`"><input type="button" name="cancel" value="`)
	//line forms.html:54
	qw422016.N().S(ln.UI["cancel"])
	//line forms.html:54
	qw422016.N().S(`">`)
	//line forms.html:55
	streamcaptcha(qw422016, "create-board", ln.UI)
	//line forms.html:55
	qw422016.N().S(`<div class="form-response admin"></div>`)
//line forms.html:57
}

//line forms.html:57
func WriteCreateBoard(qq422016 qtio422016.Writer, ln lang.Pack) {
	//line forms.html:57
	qw422016 := qt422016.AcquireWriter(qq422016)
	//line forms.html:57
	StreamCreateBoard(qw422016, ln)
	//line forms.html:57
	qt422016.ReleaseWriter(qw422016)
//line forms.html:57
}

//line forms.html:57
func CreateBoard(ln lang.Pack) string {
	//line forms.html:57
	qb422016 := qt422016.AcquireByteBuffer()
	//line forms.html:57
	WriteCreateBoard(qb422016, ln)
	//line forms.html:57
	qs422016 := string(qb422016.B)
	//line forms.html:57
	qt422016.ReleaseByteBuffer(qb422016)
	//line forms.html:57
	return qs422016
//line forms.html:57
}

//line forms.html:59
func streamcaptcha(qw422016 *qt422016.Writer, id string, lang map[string]string) {
	//line forms.html:60
	conf := config.Get()

	//line forms.html:61
	if !conf.Captcha {
		//line forms.html:62
		return
		//line forms.html:63
	}
	//line forms.html:63
	qw422016.N().S(`<div class="captcha-container"><div id="adcopy-outer-`)
	//line forms.html:65
	qw422016.N().S(id)
	//line forms.html:65
	qw422016.N().S(`"><div id="adcopy-puzzle-image-`)
	//line forms.html:66
	qw422016.N().S(id)
	//line forms.html:66
	qw422016.N().S(`" class="captcha-image" title="`)
	//line forms.html:66
	qw422016.N().S(lang["reloadCaptcha"])
	//line forms.html:66
	qw422016.N().S(`"></div><div id="adcopy-puzzle-audio-`)
	//line forms.html:67
	qw422016.N().S(id)
	//line forms.html:67
	qw422016.N().S(`" class="hidden"></div><div id="adcopy-pixel-image-`)
	//line forms.html:68
	qw422016.N().S(id)
	//line forms.html:68
	qw422016.N().S(`" class="hidden"></div><div><span id="adcopy-instr-`)
	//line forms.html:70
	qw422016.N().S(id)
	//line forms.html:70
	qw422016.N().S(`" class="hidden"></span></div><input id="adcopy_response-`)
	//line forms.html:72
	qw422016.N().S(id)
	//line forms.html:72
	qw422016.N().S(`" name="adcopy_response" class="full-width" type="text" placeholder="`)
	//line forms.html:72
	qw422016.N().S(lang["focusForCaptcha"])
	//line forms.html:72
	qw422016.N().S(`" required><input type="hidden" name="adcopy_challenge" id="adcopy_challenge-`)
	//line forms.html:73
	qw422016.N().S(id)
	//line forms.html:73
	qw422016.N().S(`" hidden><a id="adcopy-link-refresh-`)
	//line forms.html:74
	qw422016.N().S(id)
	//line forms.html:74
	qw422016.N().S(`" class="hidden"></a><a id="adcopy-link-audio-`)
	//line forms.html:75
	qw422016.N().S(id)
	//line forms.html:75
	qw422016.N().S(`" class="hidden"></a><a id="adcopy-link-image-`)
	//line forms.html:76
	qw422016.N().S(id)
	//line forms.html:76
	qw422016.N().S(`" class="hidden"></a><a id="adcopy-link-info-`)
	//line forms.html:77
	qw422016.N().S(id)
	//line forms.html:77
	qw422016.N().S(`" class="hidden"></a><noscript><iframe src="http://api.solvemedia.com/papi/challenge.noscript?k=`)
	//line forms.html:79
	qw422016.N().U(conf.CaptchaPublicKey)
	//line forms.html:79
	qw422016.N().S(`"height="260" width="350" frameborder="0"></iframe><br><input name="adcopy_challenge" type="text" required></noscript></div></div>`)
//line forms.html:86
}

//line forms.html:86
func writecaptcha(qq422016 qtio422016.Writer, id string, lang map[string]string) {
	//line forms.html:86
	qw422016 := qt422016.AcquireWriter(qq422016)
	//line forms.html:86
	streamcaptcha(qw422016, id, lang)
	//line forms.html:86
	qt422016.ReleaseWriter(qw422016)
//line forms.html:86
}

//line forms.html:86
func captcha(id string, lang map[string]string) string {
	//line forms.html:86
	qb422016 := qt422016.AcquireByteBuffer()
	//line forms.html:86
	writecaptcha(qb422016, id, lang)
	//line forms.html:86
	qs422016 := string(qb422016.B)
	//line forms.html:86
	qt422016.ReleaseByteBuffer(qb422016)
	//line forms.html:86
	return qs422016
//line forms.html:86
}
