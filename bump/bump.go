package bump

import (
	"path/filepath"
	"sync"
)

// New returns a Version that has v.safety() pre-performed on it. A demo of this package can be viewed here:
//
// 		https://go.dev/play/p/JlPROpUD2B3
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

// Create will return a new *Version with a created file and parsed version string completely valid
//
// Example:
// 		package main
// 		import (
//			"github.com/andreimerlescu/bump/bump"
//			"github.com/andreimerlescu/checkfs"
// 			"github.com/andreimerlescu/checkfs/file"
// 		)
// 		func main() {
//			path := filepath.Join(os.TempDir(), "VERSION")
// 			version, err := bump.Create("v1.2.3-beta.4", path)
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 			err = checkfs.File(path, file.Options{ Exists: true })
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 		}
//
func Create(version, path string) (*Version, error) {
	v := New()
	v.raw = []byte(version)
	v.path = path
	err := v.parse("VERSION", v.raw)
	if err != nil {
		return nil, err
	}
	err = v.Save(v.path)
	if err != nil {
		return nil, err
	}
	return v, nil
}
