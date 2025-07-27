package bump

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestBumpMajor ensures the major version is incremented correctly and resets other fields.
func TestBumpMajor(t *testing.T) {
	v := &Version{Major: 1, Minor: 2, Patch: 3, RC: 4, Beta: 5, Alpha: 6}
	v.BumpMajor()

	assert.Equal(t, 2, v.Major, "Major should be incremented")
	assert.Equal(t, 0, v.Minor, "Minor should be reset")
	assert.Equal(t, 0, v.Patch, "Patch should be reset")
	assert.Equal(t, 0, v.RC, "RC should be reset")
	assert.Equal(t, 0, v.Beta, "Beta should be reset")
	assert.Equal(t, 0, v.Alpha, "Alpha should be reset")
}

// TestBumpMinor ensures the minor version is incremented correctly and resets relevant fields.
func TestBumpMinor(t *testing.T) {
	v := &Version{Major: 1, Minor: 2, Patch: 3, RC: 4}
	v.BumpMinor()

	assert.Equal(t, 1, v.Major, "Major should not change")
	assert.Equal(t, 3, v.Minor, "Minor should be incremented")
	assert.Equal(t, 0, v.Patch, "Patch should be reset")
	assert.Equal(t, 0, v.RC, "RC should be reset")
}

// TestBumpPatch ensures the patch version is incremented correctly and resets pre-release fields.
func TestBumpPatch(t *testing.T) {
	v := &Version{Major: 1, Minor: 2, Patch: 3, RC: 4}
	v.BumpPatch()

	assert.Equal(t, 1, v.Major, "Major should not change")
	assert.Equal(t, 2, v.Minor, "Minor should not change")
	assert.Equal(t, 4, v.Patch, "Patch should be incremented")
	assert.Equal(t, 0, v.RC, "RC should be reset")
}

// TestStringerFormatting checks that the String() method formats versions correctly.
func TestStringerFormatting(t *testing.T) {
	testCases := []struct {
		name     string
		version  Version
		expected string
	}{
		{"Standard Version", Version{Major: 1, Minor: 2, Patch: 3}, "v1.2.3"},
		{"With Alpha", Version{Major: 1, Minor: 0, Patch: 0, Alpha: 5}, "v1.0.0-alpha.5"},
		{"With Beta", Version{Major: 2, Minor: 1, Patch: 0, Beta: 2}, "v2.1.0-beta.2"},
		{"With RC", Version{Major: 3, Minor: 0, Patch: 0, RC: 1}, "v3.0.0-rc.1"},
		{"With Beta and Alpha", Version{Major: 0, Minor: 1, Patch: 0, Beta: 4, Alpha: 1}, "v0.1.0-beta.4-alpha.1"},
		{"With Preview", Version{Major: 1, Minor: 2, Patch: 3, Preview: 7}, "v1.2.3-preview.7"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.version.String())
		})
	}
}

// TestSaveLoad ensures that saving and loading a version file works as expected.
func TestSaveLoad(t *testing.T) {
	tempDir := t.TempDir()
	versionFile := filepath.Join(tempDir, "VERSION")

	// 1. Create and save a version
	v1 := &Version{Major: 1, Minor: 8, Patch: 0, RC: 2}
	err := v1.Save(versionFile)
	assert.NoError(t, err, "Saving version should not produce an error")

	// 2. Load it into a new struct
	v2 := &Version{}
	err = v2.ParseFile(versionFile)
	assert.NoError(t, err, "Parsing version file should not produce an error")

	// 3. Verify the data is correct
	assert.Equal(t, "v1.8.0-rc.2", v2.String(), "Loaded version string should match saved version")
	assert.Equal(t, v1.Major, v2.Major)
	assert.Equal(t, v1.Minor, v2.Minor)
	assert.Equal(t, v1.Patch, v2.Patch)
	assert.Equal(t, v1.RC, v2.RC)
}

// FuzzParse runs a fuzz test against the Scan method to ensure it doesn't panic on unexpected inputs.
func FuzzParse(f *testing.F) {
	// Add seed values for both valid and invalid formats
	f.Add("v1.2.3")
	f.Add("v1.2.3-alpha.1")
	f.Add("v1.2.3-beta.1")
	f.Add("v1.2.3-rc.1")
	f.Add("v1.2.3-beta.1-alpha.1")
	f.Add("v1.2.3-preview.1")
	f.Add("vx.y.z")
	f.Add("1.2.3")
	f.Add("v1.2.3-garbage")
	f.Add("")
	f.Add("v-1.-2.-3")

	f.Fuzz(func(t *testing.T, input string) {
		v := &Version{}
		v.raw = []byte(input)

		// The goal is to ensure Scan() never panics, regardless of input.
		// The fuzz runner will report a failure if a panic occurs.
		err := v.Scan()

		if err == nil {
			// If parsing succeeded, the String() output should be parsable again.
			v2 := &Version{}
			v2.raw = []byte(v.String())
			assert.NoError(t, v2.Scan(), "A successfully parsed version should be re-parsable")
		}
	})
}

// BenchmarkScan measures the performance of the Scan method.
func BenchmarkScan(b *testing.B) {
	v := &Version{}
	// Use a complex version string to exercise the full scanning logic
	v.raw = []byte("v1.2.3-beta.4-alpha.5")
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = v.Scan()
	}
}
