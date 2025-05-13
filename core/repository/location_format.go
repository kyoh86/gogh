package repository

import (
	"encoding/json"
	"strings"
)

//TODO: move to app/location_format/location_format_usecase.go

// LocationFormat defines the interface for formatting local repository references
type LocationFormat interface {
	Format(ref Location) (string, error)
}

// LocationFormatFunc is a function type that implements the LocalRepoFormat interface
type LocationFormatFunc func(Location) (string, error)

// Format calls the function itself to format the local repository reference
func (f LocationFormatFunc) Format(ref Location) (string, error) {
	return f(ref)
}

// LocalRepoFormatRelPath formats the local repository reference to its full path
var LocationFormatFullPath = LocationFormatFunc(func(ref Location) (string, error) {
	return ref.FullPath(), nil
})

// LocalRepoFormatRelFilePath formats the local repository reference to its path
var LocationFormatPath = LocationFormatFunc(func(ref Location) (string, error) {
	return ref.Path(), nil
})

// LocationFormatJSON formats the local repository reference to a JSON string
var LocationFormatJSON = LocationFormatFunc(func(ref Location) (string, error) {
	buf, _ := json.Marshal(map[string]any{
		"fullPath": ref.FullPath(),
		"path":     ref.Path(),
		"host":     ref.Host(),
		"owner":    ref.Owner(),
		"name":     ref.Name(),
	})
	return string(buf), nil
})

// LocationFormatFields formats the local repository reference to a string with specified fields
func LocationFormatFields(s string) LocationFormat {
	return LocationFormatFunc(func(ref Location) (string, error) {
		return strings.Join([]string{
			ref.FullPath(),
			ref.Path(),
			ref.Host(),
			ref.Owner(),
			ref.Name(),
		}, s), nil
	})
}
