package alias

import (
	"os"

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
		return yaml.NewDecoder(file).Decode(&Instance)
	case os.IsNotExist(err):
		return nil
	default:
		return err
	}
}
