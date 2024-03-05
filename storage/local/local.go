package local

import (
	"bytes"
	"errors"
	"io"
	logg "log"
	"os"
	"strconv"
	"strings"
	"syscall"
)

var ErrEnvVarEmpty = errors.New("getenv: environment variable empty")

func Getenv(key string) string {
	err := load()
	if err != nil {
		logg.Print("Error loading", err)
	}
	return os.Getenv(key)
}
func GetenvStr(key string) (string, error) {
	err := load()
	if err != nil {
		return "", ErrEnvVarEmpty
	}
	v := os.Getenv(key)
	if v == "" {
		return v, ErrEnvVarEmpty
	}
	return v, nil
}
func GetenvInt(key string) (int, error) {
	s, err := GetenvStr(key)
	if err != nil {
		return 0, err
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func GetenvBool(key string) (bool, error) {
	s, err := GetenvStr(key)
	if err != nil {
		return false, err
	}
	v, err := strconv.ParseBool(s)
	if err != nil {
		return false, err
	}
	return v, nil
}

func MustMapEnv(target *string, envKey string) {
	v := Getenv(envKey)
	if v == "" {
		logg.Print("environment variable %q not set", envKey)
	}
	*target = v
}
func Setenv(key, value string) error {
	err := syscall.Setenv(key, value)
	if err != nil {
		logg.Print("environment variable %q not set", key)
	}
	return nil
}

func loadFile(filename string, overload bool) error {
	envMap, err := readFile(filename)
	if err != nil {
		return err
	}

	currentEnv := map[string]bool{}
	rawEnv := os.Environ()
	for _, rawEnvLine := range rawEnv {
		key := strings.Split(rawEnvLine, "=")[0]
		currentEnv[key] = true
	}

	for key, value := range envMap {
		if !currentEnv[key] || overload {
			_ = os.Setenv(key, value)
		}
	}

	return nil
}
func load(filenames ...string) (err error) {
	filenames = filenamesOrDefault(filenames)

	for _, filename := range filenames {
		err = loadFile(filename, false)
		if err != nil {
			return // return early on a spazout
		}
	}
	return
}
func filenamesOrDefault(filenames []string) []string {
	if len(filenames) == 0 {
		return []string{".env"}
	}
	return filenames
}
func readFile(filename string) (envMap map[string]string, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	return parse(file)
}

// Parse reads an env file from io.Reader, returning a map of keys and values.
func parse(r io.Reader) (map[string]string, error) {
	var buf bytes.Buffer
	_, err := io.Copy(&buf, r)
	if err != nil {
		return nil, err
	}

	return unmarshalBytes(buf.Bytes())
}

// Unmarshal reads an env file from a string, returning a map of keys and values.
func unmarshal(str string) (envMap map[string]string, err error) {
	return unmarshalBytes([]byte(str))
}

// UnmarshalBytes parses env file from byte slice of chars, returning a map of keys and values.
func unmarshalBytes(src []byte) (map[string]string, error) {
	out := make(map[string]string)
	err := parseBytes(src, out)

	return out, err
}
