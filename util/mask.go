package util

import (
	"encoding/json"
	"fmt"
	"github.com/ThomasNguyenGitHub/go/storage/local"
	"log"
	"strconv"
	"strings"
)

const (
	maskingChar = "•"
	maskingLen  = 4
)

var (
	maskingFields = map[string]bool{
		"pin":                   true,
		"password":              true,
		"passcode":              true,
		"new_passcode":          true,
		"new_password":          true,
		"otp":                   true,
		"secret":                true,
		"credential":            true,
		"access_token":          true,
		"token_access":          true,
		"authorization_app_key": true,
	}
	loadCustomMaskingField = false
)

// Mask masks something
func Mask(s string, noFirstChars, noLastChars int) string {
	return MaskByChar(s, maskingChar, noFirstChars, noLastChars, false)
}

func MaskNoSpace(s string, noFirstChars, noLastChars int) string {
	return MaskByChar(s, maskingChar, noFirstChars, noLastChars, true)
}

// MaskByChar masks something by specific char
func MaskByChar(s, maskChar string, noFirstChars, noLastChars int, noSpace bool) string {
	var (
		l           = len(s)
		noMaskChars = maskingLen
	)
	if l == 0 || l < noMaskChars || noFirstChars > noMaskChars || noLastChars > noMaskChars {
		return s
	}
	mc := maskChar
	if v := local.Getenv("MASKING_CHAR"); v != "" {
		mc = v
	}
	if mc == "" {
		mc = maskingChar
	}
	if v, _ := strconv.Atoi("MASKING_CHAR_LEN"); v > 0 {
		noMaskChars = v
	}
	var format = "%s %s %s"
	if noFirstChars <= 0 {
		noFirstChars = 0
		format = "%s%s %s"
	}
	if noSpace {
		format = strings.ReplaceAll(format, " ", "")
	}
	return fmt.Sprintf(format, s[:noFirstChars], strings.Repeat(mc, noMaskChars), s[l-noLastChars:])
}

// MaskEmail masks an email
func MaskEmail(s string, noFirstChars, noLastChars int) string {
	l := len(s)
	if l == 0 {
		return s
	}
	var (
		sep   = "@"
		parts = strings.Split(s, sep)
	)
	if len(parts) != 2 {
		return s
	}
	v := MaskNoSpace(parts[0], noFirstChars, noLastChars)
	return v + sep + parts[1]
}

func MaskFieldsFromBytes(b []byte) map[string]interface{} {
	var m map[string]interface{}
	if !json.Valid(b) {
		log.Printf("Invalid JSON: %s", string(b))
		return m
	}
	if !loadCustomMaskingField {
		for _, v := range strings.Split(local.Getenv("MASKING_FIELDS"), ",") {
			maskingFields[v] = true
		}
		loadCustomMaskingField = true
	}
	if err := json.Unmarshal(b, &m); err != nil {
		return m
	}
	maskFieldsRecursively(m)
	return m
}

func MaskFieldsFromInterface(v interface{}) map[string]interface{} {
	if v == nil {
		return make(map[string]interface{})
	}
	b, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	return MaskFieldsFromBytes(b)
}

func maskFieldsRecursively(m map[string]interface{}) {
	for k, v := range m {
		if v == nil {
			continue
		}
		nestedMap, ok := v.(map[string]interface{})
		if ok {
			maskFieldsRecursively(nestedMap)
		}
		if maskingFields[strings.ToLower(k)] {
			m[k] = MaskNoSpace(fmt.Sprint(nestedMap[k]), 0, 0)
		}
	}
}
