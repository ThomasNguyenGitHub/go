package util

import (
	"bytes"
	"io"
	"strings"
	"unicode"

	"github.com/gabriel-vasile/mimetype"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

/*func ReadFile(file io.Reader) ([]byte, string, error) {
	var w = new(bytes.Buffer)
	if _, err := io.Copy(w, file); err != nil {
		return nil, "", err
	}
	fileBytes := w.Bytes()
	return fileBytes, http.DetectContentType(fileBytes), nil
}*/

func ReadFile(file io.Reader) ([]byte, string, error) {
	var w = new(bytes.Buffer)
	if _, err := io.Copy(w, file); err != nil {
		return nil, "", err
	}
	var (
		contentType string
		fileBytes   = w.Bytes()
		mt          = mimetype.Detect(fileBytes)
	)
	if mt != nil {
		contentType = mt.String()
	}
	return fileBytes, contentType, nil
}

func RemoveAccents(str string) string {
	var (
		t            = transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
		result, _, _ = transform.String(t, str)
		r            = strings.NewReplacer("đ", "d", "Đ", "D")
	)
	return r.Replace(result)
}

func NormalizeFileName(fileName string) string {
	return strings.TrimSpace(strings.ReplaceAll(RemoveAccents(fileName), " ", ""))
}
