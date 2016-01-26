package conf

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"

	"github.com/kyoh86/gogh/env"
)

// Configures for GOGH
type Configures struct {
	AccessToken string `json:"accessToken"`
}

var (
	ConfigurationItemAccessToken = "access-token"
)

var ConfigureItems = []string{
	ConfigurationItemAccessToken,
}

var (
	// ErrNotUpdated stops `Set` func saving configurations
	ErrNotUpdated = errors.New("not updated")
)

// Set configures and save them
func Set(updater func(Configures) (Configures, error)) error {
	// open file
	file, err := os.OpenFile(env.ConfigurationsFile, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	// unmarshal
	var c Configures
	buf, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	if len(buf) > 0 {
		if err := json.Unmarshal(buf, &c); err != nil {
			return err
		}
	}

	// update
	c, err = updater(c)
	if err != nil {
		if err == ErrNotUpdated {
			return nil
		}
		return err
	}

	// save file
	if err := file.Truncate(0); err != nil {
		return err
	}
	buf, err = json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	_, err = file.Write(buf)
	if err != nil {
		return err
	}

	return nil
}

// Get configures from file
func Get() (*Configures, error) {
	file, err := os.OpenFile(env.ConfigurationsFile, os.O_CREATE|os.O_RDONLY, 0600)
	if err != nil {
		if os.IsNotExist(err) {
			return &Configures{}, nil
		}
		return nil, err
	}
	defer file.Close()

	dec := json.NewDecoder(file)
	var c Configures
	if err := dec.Decode(&c); err != nil {
		if err == io.EOF {
			return &c, nil
		}
		return nil, err
	}
	return &c, nil
}
