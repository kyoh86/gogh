package main

import (
	"encoding"
	"fmt"
	"os"
	"path/filepath"
)

type ExpandablePath struct {
	raw      string
	expanded string
}

func (p *ExpandablePath) UnmarshalText(raw []byte) error {
	ex, err := ParsePath(string(raw))
	if err != nil {
		return err
	}
	*p = ex
	return nil
}

func (p ExpandablePath) MarshalText() ([]byte, error) {
	return []byte(p.raw), nil
}

func expandPath(p string) (string, error) {
	p = os.ExpandEnv(p)
	runes := []rune(p)
	switch len(runes) {
	case 0:
		return p, nil
	case 1:
		if runes[0] == '~' {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return "", fmt.Errorf("search user home dir: %w", err)
			}
			return homeDir, nil
		}
	default:
		if runes[0] == '~' && (runes[1] == filepath.Separator || runes[1] == '/') {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return "", fmt.Errorf("search user home dir: %w", err)
			}
			return filepath.Join(homeDir, string(runes[2:])), nil
		}
	}
	return p, nil
}

func ParsePath(raw string) (ExpandablePath, error) {
	expanded, err := expandPath(raw)
	if err != nil {
		return ExpandablePath{}, fmt.Errorf("expand path: %w", err)
	}
	return ExpandablePath{
		raw:      raw,
		expanded: expanded,
	}, nil
}

var (
	_ encoding.TextUnmarshaler = (*ExpandablePath)(nil)
	_ encoding.TextMarshaler   = (*ExpandablePath)(nil)
)
