package bump

import (
	"path/filepath"
	"sync"
)

// New returns a Version that has v.safety() pre-performed on it
//
// Example:
// 		version := bump.New()
func New() *Version {
	return &Version{
		parsed: make(map[string]interface{}),
		mu:     &sync.RWMutex{},
	}
}

// Parse returns a new Version (or error) once New() has been called and v.parse() is validated
//
// Example:
// 		version, err := bump.Parse("v1.2.3-alpha.4")
func Parse(version string) (*Version, error) {
	v := New()
	v.raw = []byte(version)
	v.path = filepath.Join(".", "VERSION")
	err := v.parse("VERSION", v.raw)
	if err != nil {
		return nil, err
	}
	return v, nil
}
