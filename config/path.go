package config

import (
	"encoding"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

type Path struct {
	raw      string
	expanded string
}

func (p *Path) Set(v string) error {
	p.raw = v
	p.expanded = v // NOTE: path is always expanded in the flag
	return nil
}

func (p Path) String() string {
	return p.raw
}

func(p Path) Expand() string {
	return p.expanded
}

func (p Path) Type() string {
	return "string"
}

func (p *Path) UnmarshalText(raw []byte) error {
	ex, err := parsePath(string(raw))
	if err != nil {
		return err
	}
	*p = ex
	return nil
}

func (p Path) MarshalText() ([]byte, error) {
	return []byte(p.raw), nil
}

func expand(p string) (string, error) {
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

func parsePath(raw string) (Path, error) {
	expanded, err := expand(raw)
	if err != nil {
		return Path{}, fmt.Errorf("expand path: %w", err)
	}
	return Path{
		raw:      raw,
		expanded: expanded,
	}, nil
}

var (
	_ encoding.TextUnmarshaler = (*Path)(nil)
	_ encoding.TextMarshaler   = (*Path)(nil)
	_ flag.Value               = (*Path)(nil)
)
