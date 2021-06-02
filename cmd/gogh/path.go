package main

import (
	"encoding"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

type expandedPath struct {
	raw      string
	expanded string
}

func (p *expandedPath) Set(v string) error {
	p.raw = v
	p.expanded = v // NOTE: path is always expanded in the flag
	return nil
}

func (p expandedPath) String() string {
	return p.raw
}

func (p expandedPath) Type() string {
	return "string"
}

func (p *expandedPath) UnmarshalText(raw []byte) error {
	ex, err := parsePath(string(raw))
	if err != nil {
		return err
	}
	*p = ex
	return nil
}

func (p expandedPath) MarshalText() ([]byte, error) {
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

func parsePath(raw string) (expandedPath, error) {
	expanded, err := expandPath(raw)
	if err != nil {
		return expandedPath{}, fmt.Errorf("expand path: %w", err)
	}
	return expandedPath{
		raw:      raw,
		expanded: expanded,
	}, nil
}

var (
	_ encoding.TextUnmarshaler = (*expandedPath)(nil)
	_ encoding.TextMarshaler   = (*expandedPath)(nil)
	_ flag.Value               = (*expandedPath)(nil)
)
