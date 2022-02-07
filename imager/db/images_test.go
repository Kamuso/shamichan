package db

import (
	"bytes"
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/bakape/shamichan/imager/assets"
	"github.com/bakape/shamichan/imager/common"
	"github.com/bakape/shamichan/imager/test"
	"github.com/jackc/pgx/v4"
)

// Create image diretories and return a function that deletes them
func setupImageDirs(t *testing.T) {
	t.Helper()

	if err := assets.CreateDirs(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := assets.DeleteDirs(); err != nil {
			t.Fatal(err)
		}
	})
}

func prepareSampleImage(t *testing.T) (
	img common.ImageCommon,
	files [2]*os.File,
) {
	t.Helper()

	clearTables(t, "images")
	setupImageDirs(t)

	img = common.ImageCommon{
		Width:       300,
		Height:      300,
		ThumbHeight: 150,
		ThumbWidth:  150,
		Size:        1 << 20,
	}
	copy(img.SHA1[:], test.GenBuf(20))
	copy(img.MD5[:], test.GenBuf(16))

	assertNoImage(t, img.SHA1)

	for i, name := range [...]string{"sample", "thumb"} {
		files[i] = test.OpenSample(t, name+".jpg")
	}
	t.Cleanup(func() {
		for _, f := range files {
			f.Close()
		}
	})
	err := InTransaction(context.Background(), func(tx pgx.Tx) error {
		return AllocateImage(context.Background(), tx, img, files[0], files[1])
	})
	if err != nil {
		t.Fatal(err)
	}

	return
}

func assertNoImage(t *testing.T, id common.SHA1Hash) {
	t.Helper()

	err := InTransaction(context.Background(), func(tx pgx.Tx) (err error) {
		_, err = GetImage(context.Background(), tx, id)
		return
	})
	test.AssertEquals(t, err, pgx.ErrNoRows)
}

func TestAllocateImage(t *testing.T) {
	std, files := prepareSampleImage(t)

	// Assert files
	t.Run("files", func(t *testing.T) {
		for i, path := range assets.GetFilePaths(
			std.SHA1,
			common.JPEG,
			common.JPEG,
		) {
			buf, err := ioutil.ReadFile(path)
			if err != nil {
				t.Error(err)
			}

			_, err = files[i].Seek(0, 0)
			if err != nil {
				t.Fatal(err)
			}
			res, err := ioutil.ReadAll(files[i])
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(buf, res) {
				t.Error("invalid file")
			}
		}
	})

	// Assert database record
	t.Run("db row", func(t *testing.T) {
		var img common.ImageCommon
		err := InTransaction(context.Background(), func(tx pgx.Tx) (err error) {
			img, err = GetImage(context.Background(), tx, std.SHA1)
			return
		})
		if err != nil {
			t.Fatal(err)
		}
		test.AssertEquals(t, img, std)
	})
}

// func TestInsertImage(t *testing.T) {
// 	clearTables(t, "threads")
// 	thread, authKey := insertSampleThread(t)
// 	img, _, close := prepareSampleImage(t)
// 	defer close()

// 	const name = "fuko_da.jpeg"

// 	assertCanInsert := func(std bool) {
// 		t.Helper()

// 		can, err := CanInsertImage(context.Background(), authKey)
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		test.AssertEquals(t, can, std)
// 	}

// 	assertCanInsert(true)

// 	err := InTransaction(context.Background(), func(tx pgx.Tx) (err error) {
// 		resThread, err := InsertImage(
// 			context.Background(),
// 			tx,
// 			thread,
// 			authKey,
// 			img.SHA1,
// 			name,
// 			false,
// 		)
// 		if err != nil {
// 			return
// 		}

// 		test.AssertEquals(t, resThread, thread)

// 		return
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	assertCanInsert(false)

// 	assertPost := func(spoilered bool) {
// 		t.Helper()

// 		res, err := GetPost(context.Background(), thread)
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		var co struct {
// 			Created_on int
// 		}
// 		err = json.Unmarshal(res, &co)
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		test.AssertJSON(t, bytes.NewReader(res), map[string]interface{}{
// 			"id":   thread,
// 			"body": map[string]interface{}{},
// 			"flag": nil,
// 			"name": nil,
// 			"open": true,
// 			"page": 0,
// 			"sage": false,
// 			"trip": nil,
// 			"image": map[string]interface{}{
// 				"md5":          hex.EncodeToString(img.MD5[:]),
// 				"name":         "fuko_da.jpeg",
// 				"sha1":         hex.EncodeToString(img.SHA1[:]),
// 				"size":         1048576,
// 				"audio":        false,
// 				"title":        nil,
// 				"video":        false,
// 				"width":        300,
// 				"artist":       nil,
// 				"height":       300,
// 				"duration":     0,
// 				"file_type":    "JPEG",
// 				"spoilered":    spoilered,
// 				"thumb_type":   "JPEG",
// 				"thumb_width":  150,
// 				"thumb_height": 150,
// 			},
// 			"thread":     thread,
// 			"created_on": co.Created_on,
// 		})
// 	}

// 	assertPost(false)

// 	err = SpoilerImage(context.Background(), thread)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	assertPost(true)
// }
