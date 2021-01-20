package alias

import (
	"io"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

var Instance Def

func Set(alias, fullpath string) {
	Instance.Set(alias, fullpath)
}

func Del(alias string) {
	Instance.Del(alias)
}

func Lookup(alias string) (fullpath string) {
	return Instance.Lookup(alias)
}

func Reverse(fullpath string) []string {
	return Instance.Reverse(fullpath)
}

func LoadInstance(filename string) (retErr error) {
	file, err := os.Open(filename)
	switch {
	case err == nil:
		defer func() {
			if err := file.Close(); err != nil && retErr == nil {
				retErr = err
				return
			}
		}()
		return DecodeInstance(file)
	case os.IsNotExist(err):
		return nil
	default:
		return err
	}
}

func DecodeInstance(r io.Reader) (retErr error) {
	return yaml.NewDecoder(r).Decode(&Instance)
}

func SaveInstance(filename string) (retErr error) {
	if err := os.MkdirAll(filepath.Dir(filename), 0700); err != nil {
		return err
	}
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	switch {
	case err == nil:
		defer func() {
			if err := file.Close(); err != nil && retErr == nil {
				retErr = err
				return
			}
		}()
		return EncodeInstance(file)
	default:
		return err
	}
}

func EncodeInstance(w io.Writer) (retErr error) {
	return yaml.NewEncoder(w).Encode(&Instance)
}
