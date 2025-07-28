package bump

import (
	"bytes"
	"fmt"
	"os"
)

type Version struct {
	raw     []byte
	Major   int    `json:"major"`
	Minor   int    `json:"minor"`
	Patch   int    `json:"patch"`
	Alpha   int    `json:"alpha"`
	Beta    int    `json:"beta"`
	RC      int    `json:"rc"`
	Preview int    `json:"preview"`
	Version string `json:"version"`
}

// Version formats
const (
	TypeA string = "v%d.%d.%d"
	TypeB string = "v%d.%d.%d-alpha.%d"
	TypeC string = "v%d.%d.%d-beta.%d"
	TypeD string = "v%d.%d.%d-rc.%d"
	TypeE string = "v%d.%d.%d-beta.%d-alpha.%d"
	TypeF string = "v%d.%d.%d-preview.%d"
)

var Priority = map[int]string{
	0: TypeE,
	1: TypeB,
	2: TypeC,
	3: TypeD,
	4: TypeF,
	5: TypeA,
}

// Types is a map of format strings to the expected number of scanned items.
var Types = map[string]int{TypeA: 3, TypeB: 4, TypeC: 4, TypeD: 4, TypeE: 5, TypeF: 4}

// typesInOrder defines a deterministic order for scanning, from most specific to least.
var typesInOrder = []string{TypeE, TypeB, TypeC, TypeD, TypeF, TypeA}

// String formats the version struct into a standardized string.
func (v *Version) String() string {
	base := fmt.Sprintf(TypeA, v.Major, v.Minor, v.Patch)
	var preRelease string

	if v.Preview > 0 {
		preRelease = fmt.Sprintf("-preview.%d", v.Preview)
	} else if v.RC > 0 {
		preRelease = fmt.Sprintf("-rc.%d", v.RC)
	} else if v.Beta > 0 && v.Alpha > 0 {
		preRelease = fmt.Sprintf("-beta.%d-alpha.%d", v.Beta, v.Alpha)
	} else if v.Beta > 0 {
		preRelease = fmt.Sprintf("-beta.%d", v.Beta)
	} else if v.Alpha > 0 {
		preRelease = fmt.Sprintf("-alpha.%d", v.Alpha)
	}

	return fmt.Sprintf("%s%s", base, preRelease)
}

// Raw returns the original byte slice read from the file.
func (v *Version) Raw() string {
	return string(v.raw)
}

// Compare checks if the current version is less than (-1), equal (0), or greater than (1) another version.
func (v *Version) Compare(other *Version) int {
	if v.Major != other.Major {
		if v.Major > other.Major {
			return 1
		}
		return -1
	}
	if v.Minor != other.Minor {
		if v.Minor > other.Minor {
			return 1
		}
		return -1
	}
	if v.Patch != other.Patch {
		if v.Patch > other.Patch {
			return 1
		}
		return -1
	}

	// Pre-release tag comparison (a version without pre-release is higher)
	vIsPre := v.RC > 0 || v.Beta > 0 || v.Alpha > 0 || v.Preview > 0
	otherIsPre := other.RC > 0 || other.Beta > 0 || other.Alpha > 0 || other.Preview > 0

	if !vIsPre && otherIsPre {
		return 1
	}
	if vIsPre && !otherIsPre {
		return -1
	}

	// Compare pre-release identifiers
	if v.RC != other.RC {
		if v.RC > other.RC {
			return 1
		}
		return -1
	}
	// Add other pre-release comparisons here if needed (e.g., beta, alpha)

	return 0
}

// BumpMajor increments the major version and resets all lower-order fields.
func (v *Version) BumpMajor() {
	v.Major++
	v.Minor = 0
	v.Patch = 0
	v.RC = 0
	v.Alpha = 0
	v.Beta = 0
	v.Preview = 0
}

// BumpMinor increments the minor version and resets relevant fields.
func (v *Version) BumpMinor() {
	v.Minor++
	v.Patch = 0
	v.RC = 0
	v.Alpha = 0
	v.Beta = 0
	v.Preview = 0
}

// BumpPatch increments the patch version and resets pre-release fields.
func (v *Version) BumpPatch() {
	v.Patch++
	v.RC = 0
	v.Alpha = 0
	v.Beta = 0
	v.Preview = 0
}

// BumpRC increments the release candidate number.
func (v *Version) BumpRC() {
	v.RC++
	v.Alpha = 0
	v.Beta = 0
	v.Preview = 0
}

// BumpAlpha increments the alpha number.
func (v *Version) BumpAlpha() {
	v.Alpha++
}

// BumpBeta increments the beta number.
func (v *Version) BumpBeta() {
	v.Beta++
}

// BumpPreview increments the preview number.
func (v *Version) BumpPreview() {
	v.Preview++
	v.Patch = 0
	v.Alpha = 0
	v.Beta = 0
	v.RC = 0
}

// LoadFile reads the content of a file into the Version's raw field.
func (v *Version) LoadFile(path string) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	v.raw = bytes.TrimSpace(raw)
	return nil
}

// ParseFile loads and then scans a version file.
func (v *Version) ParseFile(path string) error {
	if err := v.LoadFile(path); err != nil {
		return err
	}
	return v.Scan()
}

// Scan attempts to parse the raw version string into the Version struct fields.
func (v *Version) Scan() error {
	v.Major, v.Minor, v.Patch, v.Alpha, v.Beta, v.RC, v.Preview = 0, 0, 0, 0, 0, 0, 0
	for _, t := range typesInOrder {
		var n int
		var err error
		rawStr := string(v.raw)
		switch t {
		case TypeA:
			n, err = fmt.Sscanf(rawStr, t, &v.Major, &v.Minor, &v.Patch)
		case TypeB:
			n, err = fmt.Sscanf(rawStr, t, &v.Major, &v.Minor, &v.Patch, &v.Alpha)
		case TypeC:
			n, err = fmt.Sscanf(rawStr, t, &v.Major, &v.Minor, &v.Patch, &v.Beta)
		case TypeD:
			n, err = fmt.Sscanf(rawStr, t, &v.Major, &v.Minor, &v.Patch, &v.RC)
		case TypeE:
			n, err = fmt.Sscanf(rawStr, t, &v.Major, &v.Minor, &v.Patch, &v.Beta, &v.Alpha)
		case TypeF:
			n, err = fmt.Sscanf(rawStr, t, &v.Major, &v.Minor, &v.Patch, &v.Preview)
		}
		if err == nil && n == Types[t] {
			return nil
		}
	}
	return fmt.Errorf("unrecognized version format: \"%s\"", string(v.raw))
}

// Save writes the current version string to the specified file path.
func (v *Version) Save(path string) error {
	return os.WriteFile(path, []byte(v.String()), 0644)
}
