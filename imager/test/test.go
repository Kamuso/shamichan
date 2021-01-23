// Package test contains utility functions used throughout the project in tests
package test

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/onsi/gomega"
)

// LogUnexpected fails the test and prints the values in an
// `expected: X got: Y` format
func LogUnexpected(t *testing.T, got, expected interface{}) {
	t.Helper()

	// Allow comparison of structs with private fields
	var options []cmp.Option
	for _, v := range [...]interface{}{expected, got} {
		if v != nil && reflect.TypeOf(v).Kind() == reflect.Struct {
			options = append(options, cmp.AllowUnexported(v))
		}
	}
	t.Fatal("\n" + cmp.Diff(expected, got, options...))
}

// AssertEquals asserts two values are deeply equal or fails the test, if
// not
func AssertEquals(t *testing.T, res, std interface{}) {
	t.Helper()
	if !reflect.DeepEqual(res, std) {
		LogUnexpected(t, res, std)
	}
}

// UnexpectedError fails the test with an unexpected error message
func UnexpectedError(t *testing.T, err error) {
	t.Helper()
	t.Fatalf("unexpected error: %s", err)
}

// AssertFileEquals reads a file from disk and asserts it equals the standard
// buffer
func AssertFileEquals(t *testing.T, path string, std []byte) {
	t.Helper()

	buf, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	AssertBufferEquals(t, buf, std)
}

// AssertBufferEquals asserts two buffers are equal
func AssertBufferEquals(t *testing.T, buf, std []byte) {
	t.Helper()

	if !bytes.Equal(buf, std) {
		t.Fatalf("buffers not equal: `%s` : `%s`", string(std), string(buf))
	}
}

// GenString produces a random base64 string of passed length
func GenString(len int) string {
	buf := GenBuf(base64.RawURLEncoding.DecodedLen(len + 1))
	return base64.RawURLEncoding.EncodeToString(buf)[:len]
}

// Produces a random buffer of passed length
func GenBuf(len int) []byte {
	buf := make([]byte, len)
	rand.Read(buf)
	return buf
}

// ReadSample reads a sample file of passed file name
func ReadSample(t *testing.T, name string) []byte {
	t.Helper()

	path := filepath.Join("testdata", name)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		t.Error(err)
	}
	return data
}

// Opens a sample file for reading
func OpenSample(t *testing.T, name string) *os.File {
	t.Helper()

	f, err := os.Open(filepath.Join("testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	return f
}

// Skip this test, if run in a CI environment
func SkipInCI(t *testing.T) {
	t.Helper()
	if os.Getenv("CI") == "true" {
		t.Skip()
	}
}

// Decode JSON into dst from buffer
func DecodeJSON(t *testing.T, buf []byte, dst interface{}) {
	t.Helper()

	err := json.Unmarshal(buf, dst)
	if err != nil {
		t.Fatalf("%s:\n%s", err, string(buf))
	}
}

func AssertJSON(t *testing.T, r io.Reader, std interface{}) {
	t.Helper()

	res, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("json: %s", string(res))

	// Strip trailing newline - encoder artefact
	if l := len(res); l != 0 && res[l-1] == '\n' {
		res = res[:l-1]
	}

	stdJSON, err := json.Marshal(std)
	if err != nil {
		t.Fatal(err)
	}

	gomega.NewGomegaWithT(t).Expect(res).To(gomega.MatchJSON(stdJSON))
}
