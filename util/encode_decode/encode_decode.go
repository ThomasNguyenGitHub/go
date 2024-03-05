package encode_decode

import (
	b64 "encoding/base64"
	"errors"
	"fmt"
	"strings"
)

const Separator = "__"

func EncodeData(data, url string) (string, error) {
	if data == "" || url == "" {
		return "", errors.New(fmt.Sprintf("[EncodeData] empty data or url, data: %s, url: %s", data, url))
	}

	var (
		step1 = b64.StdEncoding.EncodeToString([]byte(data))
		step2 = reverse(step1)
		step3 = b64.StdEncoding.EncodeToString([]byte(step2))
		step4 = b64.StdEncoding.EncodeToString([]byte(url))
		step5 = strings.Join([]string{step3, step4}, Separator)
	)

	return b64.StdEncoding.EncodeToString([]byte(step5)), nil
}

func DecodeData(data, url string) (string, error) {
	if data == "" {
		return data, nil
	}
	if url == "" {
		return "", errors.New(fmt.Sprintf("[DecodeData] empty data or url, data: %s, url: %s", data, url))
	}

	step1, err := b64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", errors.New(fmt.Sprintf("[DecodeData][step1] decode data failed, data: %v, err: %+v", data, err))
	}

	step2 := strings.Split(string(step1), Separator)
	if len(step2) < 2 {
		return "", errors.New(fmt.Sprintf("[DecodeData][step2] empty data or url, data: %s", string(step1)))
	}

	step3, err := b64.StdEncoding.DecodeString(step2[1])
	if err != nil {
		return "", errors.New(fmt.Sprintf("[DecodeData][step3] decode url failed, url: %v, err: %+v", step2[1], err))
	}

	if !strings.EqualFold(string(step3), url) {
		return "", errors.New(fmt.Sprintf("[DecodeData][step3] url not match, url: %s, inputUrl: %+v", string(step3), url))
	}

	step4, err := b64.StdEncoding.DecodeString(step2[0])
	if err != nil {
		return "", errors.New(fmt.Sprintf("[DecodeData][step4] decode data failed, data: %v, err: %+v", step2[0], err))
	}

	step5 := reverse(string(step4))
	step6, err := b64.StdEncoding.DecodeString(step5)
	if err != nil {
		return "", errors.New(fmt.Sprintf("[DecodeData][step6] decode data failed, data: %v, err: %+v", step5, err))
	}

	return string(step6), nil
}

// function, which takes a string as
// argument and return the reverse of string.
func reverse(s string) string {
	rns := []rune(s) // convert to rune
	for i, j := 0, len(rns)-1; i < j; i, j = i+1, j-1 {

		// swap the letters of the string,
		// like first with last and so on.
		rns[i], rns[j] = rns[j], rns[i]
	}

	// return the reversed string.
	return string(rns)
}
