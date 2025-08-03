package bump

import (
	"os"
	"sync"
)

// Version is a struct that is used to describe a VERSION file
type Version struct {
	mu         *sync.RWMutex          // used for protecting the parsed map in concurrent settings
	parsed     map[string]interface{} // contains unmarshal'd json|yaml|toml|ini key=>value pairs
	path       string                 // path to the source of the Version
	raw        []byte                 // raw file contents of Version file path
	noPrefix   bool                   // control whether "v" is prepended to the SemVer
	useForm    string                 // control which format to use for rendering the version
	isIgo      bool                   // determine whether or not igo is used
	igoVersion string                 // stored igo version

	Major   int    `json:"major"`
	Minor   int    `json:"minor"`
	Patch   int    `json:"patch"`
	Alpha   int    `json:"alpha"`
	Beta    int    `json:"beta"`
	RC      int    `json:"rc"`
	Preview int    `json:"preview"`
	Version string `json:"version"`
}

// LoadFile stores the []byte contents of the path into the raw property of the Version struct
//
// Example:
// 		v := bump.New()
// 		err := v.LoadFile(filepath.Join(".","VERSION"))
//      if err != nil {
// 			panic(err)
// 		}
func (v *Version) LoadFile(path string) error {
	v.safety()
	v.mu.Lock()
	defer v.mu.Unlock()
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	v.raw = raw
	v.path = path
	return nil
}

// safety is responsible for assuring that the mutex and map are not nil
func (v *Version) safety() {
	if v.mu == nil {
		v.mu = &sync.RWMutex{}
	}
	if v.parsed == nil {
		v.parsed = make(map[string]interface{})
	}
}

// String formats the version struct into a standardized string with a 'v' prefix.
func (v *Version) String() string {
	v.safety()
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.format(true)
}

// Raw returns the []byte stored in the original `-in`
func (v *Version) Raw() string {
	v.safety()
	return string(v.raw)
}

// NoPrefix indicates whether or not the Version should prepend a "v" prefix on the output Format
func (v *Version) NoPrefix() bool {
	v.safety()
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.noPrefix
}

// Compare is used to compare different Version structs for comparison
func (v *Version) Compare(o *Version) int {
	v.safety()
	if v.Major != o.Major {
		if v.Major > o.Major {
			return 1
		}
		return -1
	}
	if v.Minor != o.Minor {
		if v.Minor > o.Minor {
			return 1
		}
		return -1
	}
	if v.Patch != o.Patch {
		if v.Patch > o.Patch {
			return 1
		}
		return -1
	}

	var (
		vIsPre     = v.Preview > 0 || v.RC > 0 || v.Beta > 0 || v.Alpha > 0
		otherIsPre = o.Preview > 0 || o.RC > 0 || o.Beta > 0 || o.Alpha > 0

		check1 = vIsPre == false && otherIsPre == true
		check2 = vIsPre == true && otherIsPre == false
		check3 = vIsPre == false && otherIsPre == false
	)

	if check1 {
		return 1
	}
	if check2 {
		return -1
	}
	if check3 {
		return 0
	}

	// Both are pre-releases, compare them
	if v.Preview != o.Preview {
		return compareInt(v.Preview, o.Preview)
	}
	if v.RC != o.RC {
		return compareInt(v.RC, o.RC)
	}
	if v.Beta != o.Beta {
		return compareInt(v.Beta, o.Beta)
	}
	if v.Alpha != o.Alpha {
		return compareInt(v.Alpha, o.Alpha)
	}
	return 0
}

// SetRaw allows you to overwrite the contents of the `-in` file passed into the package
func (v *Version) SetRaw(raw []byte) {
	v.safety()
	v.mu.Lock()
	defer v.mu.Unlock()
	v.raw = raw
}
