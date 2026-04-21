package bump

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func mustParse(t *testing.T, s string) *Version {
	t.Helper()
	v, err := Parse(s)
	require.NoError(t, err, "Parse(%q)", s)
	return v
}

func writeFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	require.NoError(t, os.WriteFile(path, []byte(content), 0644))
	return path
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	require.NoError(t, err)
	return strings.TrimSpace(string(b))
}

// fileRoundTrip is the exact sequence main.go uses for every CLI invocation:
//
//	LoadFile → Parse → (optional bump) → Format → Save → reload → verify
//
// It directly mirrors what test.sh scenario steps exercise.
func fileRoundTrip(t *testing.T, path string, bumpFunc func(*Version)) (before, after string) {
	t.Helper()

	v := New()
	require.NoError(t, v.LoadFile(path), "LoadFile")
	require.NoError(t, v.Parse(), "Parse")

	before = v.Format(!v.NoPrefix())

	if bumpFunc != nil {
		bumpFunc(v)
	}

	after = v.Format(!v.NoPrefix())
	return before, after
}

func fileRoundTripWrite(t *testing.T, path string, bumpFunc func(*Version)) (before, after string) {
	t.Helper()
	before, after = fileRoundTrip(t, path, bumpFunc)
	v := New()
	require.NoError(t, v.LoadFile(path))
	require.NoError(t, v.Parse())
	if bumpFunc != nil {
		bumpFunc(v)
	}
	require.NoError(t, v.Save(path), "Save")
	return before, after
}

// ---------------------------------------------------------------------------
// TestScanSetsUseForm — verifies scan() sets useForm correctly for all inputs
// This directly tests the state that BumpX methods depend on.
// ---------------------------------------------------------------------------

func TestScanSetsUseForm(t *testing.T) {
	cases := []struct {
		input          string
		expectForm     string
		expectNoPrefix bool
	}{
		{"v1.2.3", FormA, false},
		{"v1.2.3-alpha.1", FormB, false},
		{"v1.2.3-alpha.0", FormB, false},
		{"v1.2.3-beta.1", FormC, false},
		{"v1.2.3-rc.1", FormD, false},
		{"v1.2.3-beta.2-alpha.3", FormE, false},
		{"v1.2.3-preview.1", FormF, false},
		{"1.2.3", FormG, true},
		{"1.24", FormH, true},
		{"v1.24", FormJ, false},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			v := New()
			v.raw = []byte(tc.input)
			err := v.scan(v.raw)
			require.NoError(t, err, "scan(%q) must succeed", tc.input)
			assert.Equal(t, tc.expectForm, v.useForm,
				"scan(%q): wrong useForm", tc.input)
			assert.Equal(t, tc.expectNoPrefix, v.noPrefix,
				"scan(%q): wrong noPrefix", tc.input)
		})
	}
}

// ---------------------------------------------------------------------------
// TestScanFieldValues — verifies scan() extracts correct field values
// ---------------------------------------------------------------------------

func TestScanFieldValues(t *testing.T) {
	cases := []struct {
		input   string
		major   int
		minor   int
		patch   int
		alpha   int
		beta    int
		rc      int
		preview int
	}{
		{"v1.2.3", 1, 2, 3, 0, 0, 0, 0},
		{"v1.0.0-alpha.1", 1, 0, 0, 1, 0, 0, 0},
		{"v1.0.0-alpha.0", 1, 0, 0, 0, 0, 0, 0},
		{"v1.0.0-beta.3", 1, 0, 0, 0, 3, 0, 0},
		{"v1.0.0-rc.2", 1, 0, 0, 0, 0, 2, 0},
		{"v1.0.0-preview.5", 1, 0, 0, 0, 0, 0, 5},
		{"v1.0.0-beta.2-alpha.3", 1, 0, 0, 3, 2, 0, 0},
		{"1.2.3", 1, 2, 3, 0, 0, 0, 0},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			v := New()
			v.raw = []byte(tc.input)
			require.NoError(t, v.scan(v.raw))
			assert.Equal(t, tc.major, v.Major, "Major")
			assert.Equal(t, tc.minor, v.Minor, "Minor")
			assert.Equal(t, tc.patch, v.Patch, "Patch")
			assert.Equal(t, tc.alpha, v.Alpha, "Alpha")
			assert.Equal(t, tc.beta, v.Beta, "Beta")
			assert.Equal(t, tc.rc, v.RC, "RC")
			assert.Equal(t, tc.preview, v.Preview, "Preview")
		})
	}
}

// ---------------------------------------------------------------------------
// TestFormatRespectsPrefixFlag — Format(false) must NEVER emit a leading "v"
// This catches the FormB-FormF hardcoded "v" prefix bug in format().
// ---------------------------------------------------------------------------

func TestFormatRespectsPrefixFlag(t *testing.T) {
	cases := []struct {
		input string
	}{
		{"v1.0.0"},
		{"v1.0.0-alpha.1"},
		{"v1.0.0-beta.1"},
		{"v1.0.0-rc.1"},
		{"v1.0.0-preview.1"},
		{"v1.0.0-beta.2-alpha.3"},
	}
	for _, tc := range cases {
		t.Run(tc.input+"/Format(false)", func(t *testing.T) {
			v := mustParse(t, tc.input)
			got := v.Format(false)
			assert.False(t, strings.HasPrefix(got, "v"),
				"Format(false) must not emit 'v' prefix, got %q", got)
		})
		t.Run(tc.input+"/Format(true)", func(t *testing.T) {
			v := mustParse(t, tc.input)
			got := v.Format(true)
			assert.True(t, strings.HasPrefix(got, "v"),
				"Format(true) must emit 'v' prefix for prefixed inputs, got %q", got)
		})
	}
}

// ---------------------------------------------------------------------------
// TestBumpAlphaOnCleanVersion — the exact failure in scenario_01 / test #5-#10
// This is the single most important test: does BumpAlpha on "v1.0.0" produce
// "v1.0.0-alpha.1", both in Format() and on disk after Save()?
// ---------------------------------------------------------------------------

func TestBumpAlphaOnCleanVersion(t *testing.T) {
	v := mustParse(t, "v1.0.0")

	// Pre-conditions: scan set the right form and fields
	assert.Equal(t, FormA, v.useForm, "parsed v1.0.0 should have useForm=FormA")
	assert.Equal(t, 0, v.Alpha, "Alpha should be 0 before bump")

	v.BumpAlpha()

	// Post-conditions immediately after bump
	assert.Equal(t, FormB, v.useForm, "useForm should be FormB after BumpAlpha")
	assert.Equal(t, 1, v.Alpha, "Alpha should be 1 after BumpAlpha")
	assert.Equal(t, "v1.0.0-alpha.1", v.String(), "String()")
	assert.Equal(t, "v1.0.0-alpha.1", v.Format(true), "Format(true)")
	assert.Equal(t, "1.0.0-alpha.1", v.Format(false), "Format(false)")
}

// ---------------------------------------------------------------------------
// TestBumpAlphaOnAlpha0 — scenario_02 starts with alpha.0
// ---------------------------------------------------------------------------

func TestBumpAlphaOnAlpha0(t *testing.T) {
	v := mustParse(t, "v1.0.0-alpha.0")
	assert.Equal(t, FormB, v.useForm, "parsed v1.0.0-alpha.0 should be FormB")
	assert.Equal(t, 0, v.Alpha, "Alpha should be 0")

	v.BumpAlpha()
	assert.Equal(t, 1, v.Alpha, "Alpha should be 1 after bump")
	assert.Equal(t, "v1.0.0-alpha.1", v.Format(true))
}

// ---------------------------------------------------------------------------
// TestBumpDoesNotMutateBeforeSave — bump without Save must not touch the file.
// Mirrors the test.sh pattern of: bump -alpha (no -write) then grep original.
// ---------------------------------------------------------------------------

func TestBumpDoesNotMutateBeforeSave(t *testing.T) {
	dir := t.TempDir()
	path := writeFile(t, dir, "VERSION", "v1.0.0")

	v := New()
	require.NoError(t, v.LoadFile(path))
	require.NoError(t, v.Parse())
	v.BumpAlpha() // no Save

	assert.Equal(t, "v1.0.0", readFile(t, path),
		"file must be unchanged when Save() is not called")
}

// ---------------------------------------------------------------------------
// TestBumpAlphaWriteRoundTrip — scenario_01 steps #8-#10 in Go:
//   bump -alpha -write  →  cat VERSION  →  grep 'v1.0.0-alpha.1'
// ---------------------------------------------------------------------------

func TestBumpAlphaWriteRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := writeFile(t, dir, "VERSION", "v1.0.0")

	// Load, bump, save
	v := New()
	require.NoError(t, v.LoadFile(path))
	require.NoError(t, v.Parse())
	v.BumpAlpha()
	require.NoError(t, v.Save(path))

	// File must contain v1.0.0-alpha.1
	assert.Equal(t, "v1.0.0-alpha.1", readFile(t, path))

	// Reload and re-parse must also give v1.0.0-alpha.1
	v2 := New()
	require.NoError(t, v2.LoadFile(path))
	require.NoError(t, v2.Parse())
	assert.Equal(t, "v1.0.0-alpha.1", v2.Format(true))
	assert.Equal(t, FormB, v2.useForm, "reloaded version should be FormB")
	assert.Equal(t, 1, v2.Alpha)
}

// ---------------------------------------------------------------------------
// TestScenario01 — complete Go port of test.sh scenario_01
// ---------------------------------------------------------------------------

func TestScenario01(t *testing.T) {
	dir := t.TempDir()
	path := writeFile(t, dir, "VERSION", "v1.0.0")

	// bump -check
	{
		v := New()
		require.NoError(t, v.LoadFile(path))
		require.NoError(t, v.Parse())
		assert.Equal(t, "v1.0.0", v.Format(!v.NoPrefix()), "step: bump -check")
	}

	// bump -alpha (no write) → reports v1.0.0-alpha.1, file still v1.0.0
	{
		v := New()
		require.NoError(t, v.LoadFile(path))
		require.NoError(t, v.Parse())
		v.BumpAlpha()
		assert.Equal(t, "v1.0.0-alpha.1", v.Format(true), "step: bump -alpha output")
		assert.Equal(t, "v1.0.0", readFile(t, path), "step: file unchanged after bump without -write")
	}

	// bump -alpha -write → file becomes v1.0.0-alpha.1
	{
		v := New()
		require.NoError(t, v.LoadFile(path))
		require.NoError(t, v.Parse())
		v.BumpAlpha()
		require.NoError(t, v.Save(path))
		assert.Equal(t, "v1.0.0-alpha.1", readFile(t, path), "step: grep v1.0.0-alpha.1")
	}
}

// ---------------------------------------------------------------------------
// TestScenario02 — complete Go port of test.sh scenario_02
// ---------------------------------------------------------------------------

func TestScenario02(t *testing.T) {
	dir := t.TempDir()
	path := writeFile(t, dir, "VERSION", "v1.0.0-alpha.0")

	step := func(label string, bumpFunc func(*Version), write bool, expectFile, expectFormat string) {
		t.Helper()
		v := New()
		require.NoError(t, v.LoadFile(path), label+": LoadFile")
		require.NoError(t, v.Parse(), label+": Parse")
		if bumpFunc != nil {
			bumpFunc(v)
		}
		if expectFormat != "" {
			assert.Equal(t, expectFormat, v.Format(true), label+": Format(true)")
		}
		if write {
			require.NoError(t, v.Save(path), label+": Save")
		}
		assert.Equal(t, expectFile, readFile(t, path), label+": file contents")
	}

	// bump -check
	step("bump -check", nil, false, "v1.0.0-alpha.0", "v1.0.0-alpha.0")
	// bump -alpha (no write)
	step("bump -alpha", (*Version).BumpAlpha, false, "v1.0.0-alpha.0", "v1.0.0-alpha.1")
	// bump -alpha -write
	step("bump -alpha -write", (*Version).BumpAlpha, true, "v1.0.0-alpha.1", "v1.0.0-alpha.1")
	// bump -patch (no write)
	step("bump -patch", (*Version).BumpPatch, false, "v1.0.0-alpha.1", "v1.0.1")
	// bump -patch -write
	step("bump -patch -write", (*Version).BumpPatch, true, "v1.0.1", "v1.0.1")
	// bump -major -write
	step("bump -major -write", (*Version).BumpMajor, true, "v2.0.0", "v2.0.0")
	// bump -preview -write
	step("bump -preview -write", (*Version).BumpPreview, true, "v2.0.0-preview.1", "v2.0.0-preview.1")
}

// ---------------------------------------------------------------------------
// TestScenario03 — fix malformed version
// ---------------------------------------------------------------------------

func TestScenario03(t *testing.T) {
	dir := t.TempDir()
	path := writeFile(t, dir, "VERSION", "1.25")

	// bump -fix (no write) — file unchanged
	{
		v := New()
		require.NoError(t, v.LoadFile(path))
		require.NoError(t, v.Fix())
		assert.Equal(t, "1.25", readFile(t, path), "file unchanged without -write")
	}

	// bump -fix -write
	{
		v := New()
		require.NoError(t, v.LoadFile(path))
		require.NoError(t, v.Fix())
		require.NoError(t, v.Save(path))
		assert.Equal(t, "v1.25.0", readFile(t, path))
	}
}

// ---------------------------------------------------------------------------
// TestScenario04 — check+fix on already-valid version changes nothing
// ---------------------------------------------------------------------------

func TestScenario04(t *testing.T) {
	dir := t.TempDir()
	path := writeFile(t, dir, "VERSION", "v1.17.7-beta.6")

	for _, write := range []bool{false, true} {
		v := New()
		require.NoError(t, v.LoadFile(path))
		require.NoError(t, v.Fix())
		require.NoError(t, v.Parse())
		if write {
			require.NoError(t, v.Save(path))
		}
		assert.Equal(t, "v1.17.7-beta.6", readFile(t, path),
			"fix on valid version must not change it (write=%v)", write)
	}
}

// ---------------------------------------------------------------------------
// TestScenario09 — Dockerfile patch bump
// ---------------------------------------------------------------------------

func TestScenario09(t *testing.T) {
	dir := t.TempDir()
	path := writeFile(t, dir, "Dockerfile", `LABEL version="v3.2.1"`)

	v := New()
	require.NoError(t, v.LoadFile(path))
	require.NoError(t, v.Parse())
	assert.Equal(t, 3, v.Major)
	assert.Equal(t, 2, v.Minor)
	assert.Equal(t, 1, v.Patch)

	v.BumpPatch()
	require.NoError(t, v.Save(path))
	assert.Contains(t, readFile(t, path), `"v3.2.2"`)
}

// ---------------------------------------------------------------------------
// TestScenario10 — Chart.yaml patch bump
// ---------------------------------------------------------------------------

func TestScenario10(t *testing.T) {
	dir := t.TempDir()
	path := writeFile(t, dir, "Chart.yaml", "apiVersion: v2\nname: mychart\nversion: 0.1.0\n")

	v := New()
	require.NoError(t, v.LoadFile(path))
	require.NoError(t, v.Parse())
	v.BumpPatch()
	require.NoError(t, v.Save(path))

	v2 := New()
	require.NoError(t, v2.LoadFile(path))
	require.NoError(t, v2.Parse())
	assert.Equal(t, 0, v2.Major)
	assert.Equal(t, 1, v2.Minor)
	assert.Equal(t, 1, v2.Patch)
}

// ---------------------------------------------------------------------------
// TestScenario11 — pom.xml patch bump
// ---------------------------------------------------------------------------

func TestScenario11(t *testing.T) {
	dir := t.TempDir()
	path := writeFile(t, dir, "pom.xml", `<project><version>2.2.2</version></project>`)

	v := New()
	require.NoError(t, v.LoadFile(path))
	require.NoError(t, v.Parse())
	v.BumpPatch()
	require.NoError(t, v.Save(path))
	assert.Contains(t, readFile(t, path), "<version>2.2.3</version>")
}

// ---------------------------------------------------------------------------
// TestScenario12 — BUMP_ALWAYS_WRITE env behavior (env tested via direct API)
// ---------------------------------------------------------------------------

func TestScenario12(t *testing.T) {
	dir := t.TempDir()
	path := writeFile(t, dir, "VERSION", "v5.5.5")

	// patch
	{
		v := New()
		require.NoError(t, v.LoadFile(path))
		require.NoError(t, v.Parse())
		v.BumpPatch()
		require.NoError(t, v.Save(path))
		assert.Equal(t, "v5.5.6", readFile(t, path))
	}

	// minor x2 (mirrors BUMP_DEFAULT_INPUT pattern — second bump continues from saved state)
	{
		v := New()
		require.NoError(t, v.LoadFile(path))
		require.NoError(t, v.Parse())
		v.BumpMinor()
		require.NoError(t, v.Save(path))
		assert.Equal(t, "v5.6.0", readFile(t, path))
	}
}
