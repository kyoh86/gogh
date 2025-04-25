package config

import (
	"fmt"
	"os"
	"sync"
)

var (
	globalTokens = TokenManager{}
	tokensOnce   sync.Once
)

func TokensPath() (string, error) {
	path, err := appFilePath("GOGH_TOKENS_PATH", os.UserCacheDir, "tokens.yaml")
	if err != nil {
		return "", fmt.Errorf("search config path: %w", err)
	}
	return path, nil
}

func LoadTokens() (_ *TokenManager, retErr error) {
	tokensOnce.Do(func() {
		path, err := TokensPath()
		if err != nil {
			retErr = err
			return
		}

		if err := loadYAML(path, &globalTokens); err != nil {
			retErr = err
			return
		}
	})
	return &globalTokens, retErr
}

func SaveTokens() error {
	path, err := TokensPath()
	if err != nil {
		return err
	}
	return saveYAML(path, globalTokens)
}
