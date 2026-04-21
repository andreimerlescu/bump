package bump

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestAllBumps ensures all bump variants increment correctly and reset lower-order fields.
func TestAllBumps(t *testing.T) {
	t.Run("BumpMajor", func(t *testing.T) {
		v, err := Parse("v1.2.3-beta.1")
		assert.NoError(t, err)
		v.BumpMajor()
		assert.Equal(t, 2, v.Major, "Major should be incremented")
		assert.Equal(t, 0, v.Minor, "Minor should be reset")
		assert.Equal(t, 0, v.Patch, "Patch should be reset")
		assert.Equal(t, 0, v.RC, "RC should be reset")
		assert.Equal(t, 0, v.Beta, "Beta should be reset")
		assert.Equal(t, 0, v.Alpha, "Alpha should be reset")
		assert.Equal(t, 0, v.Preview, "Preview should be reset")
	})

	t.Run("BumpMinor", func(t *testing.T) {
		v, err := Parse("v1.2.3-beta.1")
		assert.NoError(t, err)
		v.BumpMinor()
		assert.Equal(t, 1, v.Major, "Major should not change")
		assert.Equal(t, 3, v.Minor, "Minor should be incremented")
		assert.Equal(t, 0, v.Patch, "Patch should be reset")
		assert.Equal(t, 0, v.RC, "RC should be reset")
		assert.Equal(t, 0, v.Preview, "Preview should be reset")
	})

	t.Run("BumpPatch", func(t *testing.T) {
		v, err := Parse("v1.2.3-beta.1")
		assert.NoError(t, err)
		v.BumpPatch()
		assert.Equal(t, 1, v.Major, "Major should not change")
		assert.Equal(t, 2, v.Minor, "Minor should not change")
		assert.Equal(t, 4, v.Patch, "Patch should be incremented")
		assert.Equal(t, 0, v.RC, "RC should be reset")
	})

	t.Run("BumpRC", func(t *testing.T) {
		v, err := Parse("v1.2.3-rc.4")
		assert.NoError(t, err)
		v.BumpRC()
		assert.Equal(t, 5, v.RC, "RC should be incremented")
	})

	t.Run("BumpBeta", func(t *testing.T) {
		v, err := Parse("v1.2.3-beta.4")
		assert.NoError(t, err)
		v.BumpBeta()
		assert.Equal(t, 5, v.Beta, "Beta should be incremented")
	})

	t.Run("BumpAlpha", func(t *testing.T) {
		v, err := Parse("v1.2.3-alpha.4")
		assert.NoError(t, err)
		v.BumpAlpha()
		assert.Equal(t, 5, v.Alpha, "Alpha should be incremented")
	})

	t.Run("BumpPreview", func(t *testing.T) {
		v, err := Parse("v1.2.3-preview.4")
		assert.NoError(t, err)
		v.BumpPreview()
		assert.Equal(t, 5, v.Preview, "Preview should be incremented")
	})
}

// TestFormatting checks that String() and Format() methods work correctly.
func TestFormatting(t *testing.T) {
	testCases := []struct {
		name             string
		version          Version
		expectedString   string
		expectedFormat   string // Format(false)
		expectedNoPrefix bool
	}{
		{"Standard Version", Version{Major: 1, Minor: 2, Patch: 3}, "v1.2.3", "1.2.3", false},
		{"No Prefix Version", Version{Major: 1, Minor: 2, Patch: 3, noPrefix: true}, "1.2.3", "1.2.3", true},
		{"With Alpha", Version{Major: 1, Minor: 0, Patch: 0, Alpha: 5}, "v1.0.0-alpha.5", "1.0.0-alpha.5", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedString, tc.version.String())
			assert.Equal(t, tc.expectedFormat, tc.version.Format(tc.expectedNoPrefix == true))
			assert.Equal(t, tc.expectedNoPrefix, tc.version.NoPrefix())
		})
	}
}

// TestCompare verifies the version comparison logic.
func TestCompare(t *testing.T) {
	testCases := []struct {
		name     string
		v1       *Version
		v2       *Version
		expected int
	}{
		{"Equal", &Version{Major: 1, Minor: 2, Patch: 3}, &Version{Major: 1, Minor: 2, Patch: 3}, 0},
		{"Greater Major", &Version{Major: 2, Minor: 0, Patch: 0}, &Version{Major: 1, Minor: 9, Patch: 9}, 1},
		{"Less Minor", &Version{Major: 1, Minor: 2, Patch: 3}, &Version{Major: 1, Minor: 3, Patch: 0}, -1},
		{"Release > Pre-release", &Version{Major: 1, Minor: 0, Patch: 0}, &Version{Major: 1, Minor: 0, Patch: 0, Alpha: 1}, 1},
		{"Pre-release < Release", &Version{Major: 1, Minor: 0, Patch: 0, RC: 1}, &Version{Major: 1, Minor: 0, Patch: 0}, -1},
		{"Pre-release Equal", &Version{Major: 1, Minor: 0, Patch: 0, RC: 1}, &Version{Major: 1, Minor: 0, Patch: 0, RC: 1}, 0},
		{"Pre-release Compare", &Version{Major: 1, Minor: 0, Patch: 0, RC: 2}, &Version{Major: 1, Minor: 0, Patch: 0, RC: 1}, 1},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.v1.Compare(tc.v2))
		})
	}
}

// TestFix validates the logic for fixing malformed version strings.
func TestFix(t *testing.T) {
	testCases := []struct {
		name      string
		input     string
		expectErr bool
		expected  Version
	}{
		{"Already Valid", "v1.2.3", false, Version{Major: 1, Minor: 2, Patch: 3}},
		{"Fixable No Prefix", "1.2.3", false, Version{Major: 1, Minor: 2, Patch: 3}},
		{"Fixable Two-Part", "1.2", false, Version{Major: 1, Minor: 2, Patch: 0}},
		{"Fixable Fuzzy", "v1.5.x-rc-thing", false, Version{Major: 1, Minor: 5}},
		{"Empty Input", "", false, Version{Major: 0, Minor: 0, Patch: 1}},
		{"Unfixable Gibberish", "not-a-version", true, Version{}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			v := New()
			v.raw = []byte(tc.input)
			err := v.Fix()

			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected.Major, v.Major)
				assert.Equal(t, tc.expected.Minor, v.Minor)
				assert.Equal(t, tc.expected.Patch, v.Patch)
				assert.Equal(t, tc.expected.Beta, v.Beta)
			}
		})
	}
}

// TestParseAndSaveFileTypes provides comprehensive testing for all supported file formats.
func TestParseAndSaveFileTypes(t *testing.T) {
	testCases := []struct {
		name           string
		filename       string
		initialContent string
		initialVersion Version
		bumpFunc       func(*Version)
		finalVersion   Version
	}{
		{
			name:           "VERSION file patch bump",
			filename:       "VERSION",
			initialContent: "v1.2.3",
			initialVersion: Version{Major: 1, Minor: 2, Patch: 3},
			bumpFunc:       (*Version).BumpPatch,
			finalVersion:   Version{Major: 1, Minor: 2, Patch: 4},
		},
		{
			name:           "package.json minor bump",
			filename:       "package.json",
			initialContent: `{"name": "test-app", "version": "1.2.3"}`,
			initialVersion: Version{Major: 1, Minor: 2, Patch: 3, noPrefix: true},
			bumpFunc:       (*Version).BumpMinor,
			finalVersion:   Version{Major: 1, Minor: 3, Patch: 0},
		},
		{
			name:           "Chart.yaml major bump",
			filename:       "Chart.yaml",
			initialContent: "apiVersion: v2\nname: my-chart\nversion: 1.2.3\nappVersion: \"1.2.3\"",
			initialVersion: Version{Major: 1, Minor: 2, Patch: 3, noPrefix: true},
			bumpFunc:       (*Version).BumpMajor,
			finalVersion:   Version{Major: 2, Minor: 0, Patch: 0},
		},
		{
			name:           "Dockerfile alpha bump",
			filename:       "Dockerfile",
			initialContent: `FROM alpine\nLABEL version="v1.2.3"`,
			initialVersion: Version{Major: 1, Minor: 2, Patch: 3},
			bumpFunc:       (*Version).BumpAlpha,
			finalVersion:   Version{Major: 1, Minor: 2, Patch: 3, Alpha: 1},
		},
		{
			name:           "pom.xml patch bump",
			filename:       "pom.xml",
			initialContent: `<project><version>1.2.3</version></project>`,
			initialVersion: Version{Major: 1, Minor: 2, Patch: 3},
			bumpFunc:       (*Version).BumpPatch,
			finalVersion:   Version{Major: 1, Minor: 2, Patch: 4},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tempDir := t.TempDir()
			filePath := filepath.Join(tempDir, tc.filename)
			err := os.WriteFile(filePath, []byte(tc.initialContent), 0644)
			assert.NoError(t, err)

			// 1. Parse initial file and verify
			v1 := New()
			err = v1.ParseFile(filePath)
			assert.NoError(t, err, "Parsing initial file should succeed")
			assert.Equal(t, tc.initialVersion.Major, v1.Major)
			assert.Equal(t, tc.initialVersion.Minor, v1.Minor)
			assert.Equal(t, tc.initialVersion.Patch, v1.Patch)

			// 2. Bump version and save
			tc.bumpFunc(v1)
			err = v1.Save(filePath)
			assert.NoError(t, err, "Saving modified file should succeed")

			// 3. Parse the modified file and verify the new version
			v2 := New()
			err = v2.ParseFile(filePath)
			assert.NoError(t, err, "Parsing modified file should succeed")
			assert.Equal(t, tc.finalVersion.Major, v2.Major, "Major version should match after bump")
			assert.Equal(t, tc.finalVersion.Minor, v2.Minor, "Minor version should match after bump")
			assert.Equal(t, tc.finalVersion.Patch, v2.Patch, "Patch version should match after bump")
			assert.Equal(t, tc.finalVersion.Alpha, v2.Alpha, "Alpha version should match after bump")
		})
	}
}

func FuzzParse(f *testing.F) {
	testcases := []string{
		"v1.2.3", "1.2.3", "1.21", "v1.2.3-alpha.1", "v1.2.3-beta.1", "v1.2.3-rc.1",
		"v1.2.3-beta.1-alpha.1", "v1.2.3-preview.1", "vx.y.z", "v1.2.3-garbage", "", "v-1.-2.-3",
	}
	for _, tc := range testcases {
		f.Add(tc)
	}
	f.Fuzz(func(t *testing.T, input string) {
		v := New()
		v.raw = []byte(input)
		err := v.Parse() // Fuzz the public Parse method

		if err == nil {
			formattedString := v.String()
			v2 := New()
			v2.raw = []byte(formattedString)
			err2 := v2.Parse()
			assert.NoError(t, err2, fmt.Sprintf("A successfully parsed version string '%s' (from input '%s') should be re-parsable", formattedString, input))
		}
	})
}

func BenchmarkScan(b *testing.B) {
	v := New()
	rawVersion := []byte("v1.2.3-beta.4-alpha.5")
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = v.scan(rawVersion)
	}
}

// TestBumpAlphaFromCleanVersion tests bumping alpha on a version with no pre-release tag.
// This is scenario_01 from the CLI tests — the case the unit tests currently miss entirely.
func TestBumpAlphaFromCleanVersion(t *testing.T) {
	v, err := Parse("v1.0.0")
	assert.NoError(t, err)
	assert.Equal(t, FormA, v.useForm, "clean version should parse as FormA")

	v.BumpAlpha()

	assert.Equal(t, 1, v.Alpha, "Alpha should be 1")
	assert.Equal(t, FormB, v.useForm, "useForm should be FormB after BumpAlpha on clean version")
	assert.Equal(t, "v1.0.0-alpha.1", v.String(), "formatted string should be v1.0.0-alpha.1")
	assert.Equal(t, "v1.0.0-alpha.1", v.Format(true), "Format(true) should be v1.0.0-alpha.1")
}

// TestBumpBetaFromCleanVersion tests bumping beta on a version with no pre-release tag.
func TestBumpBetaFromCleanVersion(t *testing.T) {
	v, err := Parse("v1.0.0")
	assert.NoError(t, err)

	v.BumpBeta()

	assert.Equal(t, 1, v.Beta, "Beta should be 1")
	assert.Equal(t, FormC, v.useForm, "useForm should be FormC after BumpBeta")
	assert.Equal(t, "v1.0.0-beta.1", v.String(), "formatted string should be v1.0.0-beta.1")
}

// TestBumpFormattedOutput tests that all bump operations produce correct formatted strings,
// not just correct field values. This is the critical gap in the existing tests.
func TestBumpFormattedOutput(t *testing.T) {
	testCases := []struct {
		name       string
		input      string
		bumpFunc   func(*Version)
		expected   string
		expectForm string
	}{
		{
			name:       "alpha from clean version",
			input:      "v1.0.0",
			bumpFunc:   (*Version).BumpAlpha,
			expected:   "v1.0.0-alpha.1",
			expectForm: FormB,
		},
		{
			name:       "alpha from existing alpha",
			input:      "v1.0.0-alpha.3",
			bumpFunc:   (*Version).BumpAlpha,
			expected:   "v1.0.0-alpha.4",
			expectForm: FormB,
		},
		{
			name:       "alpha from alpha.0",
			input:      "v1.0.0-alpha.0",
			bumpFunc:   (*Version).BumpAlpha,
			expected:   "v1.0.0-alpha.1",
			expectForm: FormB,
		},
		{
			name:       "beta from clean version",
			input:      "v1.0.0",
			bumpFunc:   (*Version).BumpBeta,
			expected:   "v1.0.0-beta.1",
			expectForm: FormC,
		},
		{
			name:       "beta from existing beta",
			input:      "v1.0.0-beta.4",
			bumpFunc:   (*Version).BumpBeta,
			expected:   "v1.0.0-beta.5",
			expectForm: FormC,
		},
		{
			name:       "patch from clean version",
			input:      "v1.0.0",
			bumpFunc:   (*Version).BumpPatch,
			expected:   "v1.0.1",
			expectForm: FormA,
		},
		{
			name:       "patch from alpha clears pre-release",
			input:      "v1.0.0-alpha.3",
			bumpFunc:   (*Version).BumpPatch,
			expected:   "v1.0.1",
			expectForm: FormA,
		},
		{
			name:       "major from beta clears everything",
			input:      "v1.0.0-beta.6",
			bumpFunc:   (*Version).BumpMajor,
			expected:   "v2.0.0",
			expectForm: FormA,
		},
		{
			name:       "preview from clean version",
			input:      "v2.0.0",
			bumpFunc:   (*Version).BumpPreview,
			expected:   "v2.0.0-preview.1",
			expectForm: FormF,
		},
		{
			name:       "rc from clean version",
			input:      "v1.0.0",
			bumpFunc:   (*Version).BumpRC,
			expected:   "v1.0.0-rc.1",
			expectForm: FormD,
		},
		{
			name:       "minor from clean version",
			input:      "v1.0.0",
			bumpFunc:   (*Version).BumpMinor,
			expected:   "v1.1.0",
			expectForm: FormA,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			v, err := Parse(tc.input)
			assert.NoError(t, err)

			tc.bumpFunc(v)

			assert.Equal(t, tc.expected, v.String(), "String() output mismatch")
			assert.Equal(t, tc.expected, v.Format(true), "Format(true) output mismatch")
			assert.Equal(t, tc.expectForm, v.useForm, "useForm mismatch")
		})
	}
}

// TestBumpWriteAndReparse tests the full round-trip: bump, save, reload, verify.
// This mirrors what the CLI tests do and catches save/parse interaction bugs.
func TestBumpWriteAndReparse(t *testing.T) {
	testCases := []struct {
		name     string
		initial  string
		bumpFunc func(*Version)
		expected string
	}{
		{"alpha from clean", "v1.0.0", (*Version).BumpAlpha, "v1.0.0-alpha.1"},
		{"alpha from alpha.0", "v1.0.0-alpha.0", (*Version).BumpAlpha, "v1.0.0-alpha.1"},
		{"beta from clean", "v1.0.0", (*Version).BumpBeta, "v1.0.0-beta.1"},
		{"patch from alpha", "v1.0.0-alpha.1", (*Version).BumpPatch, "v1.0.1"},
		{"major from clean", "v1.0.0", (*Version).BumpMajor, "v2.0.0"},
		{"preview from major", "v2.0.0", (*Version).BumpPreview, "v2.0.0-preview.1"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tempDir := t.TempDir()
			path := filepath.Join(tempDir, "VERSION")

			err := os.WriteFile(path, []byte(tc.initial), 0644)
			assert.NoError(t, err)

			// Load, bump, save
			v1 := New()
			err = v1.ParseFile(path)
			assert.NoError(t, err)

			tc.bumpFunc(v1)

			// Verify formatted output BEFORE saving
			assert.Equal(t, tc.expected, v1.Format(true),
				"Format(true) should show new version before save")

			err = v1.Save(path)
			assert.NoError(t, err)

			// Reload and verify
			v2 := New()
			err = v2.ParseFile(path)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, v2.Format(true),
				"reparsed file should contain new version")

			// Also verify file contents directly
			contents, err := os.ReadFile(path)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, strings.TrimSpace(string(contents)),
				"file contents should be new version")
		})
	}
}
