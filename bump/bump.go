package bump

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

type Version struct {
	mu      sync.RWMutex
	parsed  map[string]interface{} // contains unmarshal'd json|yaml|toml|ini key=>value pairs
	path    string
	raw     []byte
	Major   int    `json:"major"`
	Minor   int    `json:"minor"`
	Patch   int    `json:"patch"`
	Alpha   int    `json:"alpha"`
	Beta    int    `json:"beta"`
	RC      int    `json:"rc"`
	Preview int    `json:"preview"`
	Version string `json:"version"`

	noPrefix bool
	useForm  string
}

func NewVersion() *Version {
	return &Version{
		parsed: make(map[string]interface{}),
		mu:     sync.RWMutex{},
	}
}

// Version formats
const (
	FormA string = "v%d.%d.%d"
	FormB string = "v%d.%d.%d-alpha.%d"
	FormC string = "v%d.%d.%d-beta.%d"
	FormD string = "v%d.%d.%d-rc.%d"
	FormE string = "v%d.%d.%d-beta.%d-alpha.%d"
	FormF string = "v%d.%d.%d-preview.%d"
	FormG string = "%d.%d.%d" // SemVer Standard
	FormH string = "%d.%d"    // Shorthand SemVer (0.1 -> 100.100)
	FormI string = "v%d"      // v# (v1 -> v100)
	FormJ string = "v%d.%d"   // v#.# (v1.1 -> v100.100)

	PatternGoMod      string = `go {{.Version}}`
	PatternDockerfile string = `LABEL version="{{.Version}}"`

	FileVersion     string = "VERSION"
	FilePackageJson string = "package.json"
	FileMavenPom    string = "pom.xml"
	FileHelmChart   string = "Chart.yaml"
	FileDockerfile  string = "Dockerfile"
	FileGoMod       string = "go.mod"
)

var SupportedFiles = []string{
	FileVersion,
	FilePackageJson,
	FileMavenPom,
	FileHelmChart,
	FileDockerfile,
	FileGoMod,
}

var (
	reTwoPartPre = regexp.MustCompile(`^(\d+)\.(\d+)(-[a-zA-Z0-9-.]+)$`)
	reThreePart  = regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)$`)
	reTwoPart    = regexp.MustCompile(`^(\d+)\.(\d+)$`)
	reFuzzy      = regexp.MustCompile(`^(\d+)\.(\d+).*`)
)

var Priority = map[int]string{
	0: FormE,
	1: FormB,
	2: FormC,
	3: FormD,
	4: FormF,
	5: FormA,
	6: FormG,
	7: FormH,
	8: FormI,
	9: FormJ,
}

// Forms is a map of format strings to the expected number of scanned items.
var Forms = map[string]int{
	FormA: 3,
	FormB: 4,
	FormC: 4,
	FormD: 4,
	FormE: 5,
	FormF: 4,
	FormG: 3,
	FormH: 2,
	FormI: 1,
	FormJ: 2,
}

// formsInOrder defines a deterministic order for scanning, from most specific to least.
var formsInOrder = []string{FormE, FormB, FormC, FormD, FormF, FormA, FormG, FormH, FormI, FormJ}

// String formats the version struct into a standardized string.
func (v *Version) String() string {
	base := fmt.Sprintf(FormA, v.Major, v.Minor, v.Patch)
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

func (v *Version) Validate() error {
	var major, minor, patch, preview, alpha, beta, rc int
	for _, t := range formsInOrder {
		rawStr := string(v.raw)
		switch t {
		case FormA:
			return v.validateFormA(rawStr, &major, &minor, &patch)
		case FormB:
			return v.validateFormB(rawStr, &major, &minor, &patch, &alpha)
		case FormC:
			return v.validateFormC(rawStr, &major, &minor, &patch, &beta)
		case FormD:
			return v.validateFormD(rawStr, &major, &minor, &patch, &rc)
		case FormE:
			return v.validateFormE(rawStr, &major, &minor, &patch, &alpha, &beta)
		case FormF:
			return v.validateFormF(rawStr, &major, &minor, &patch, &preview)
		case FormG:
			return v.validateFormG(rawStr, &major, &minor, &patch)
		case FormH:
			return v.validateFormH(rawStr, &major, &minor)
		case FormI:
			return v.validateFormI(rawStr, &major)
		case FormJ:
			return v.validateFormJ(rawStr, &major, &minor)

		}
	}
	return fmt.Errorf("unrecognized version format: \"%s\"", string(v.raw))
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
	v.path = path
	return nil
}

// ParseFile loads and then scans a version file.
func (v *Version) ParseFile(path string) error {
	if err := v.LoadFile(path); err != nil {
		return err
	}
	v.path = path
	return v.Parse()
}

func (v *Version) Parse() error {
	base := filepath.Base(v.path)
	switch base {
	case FileVersion:
		return v.parseVersion()
	case FilePackageJson:
		return v.parsePackageJson()
	case FileHelmChart:
		return v.parseHelmChart()
	case FileDockerfile:
		return v.parseDockerfile()
	case FileGoMod:
		return v.parseGoMod()
	case FileMavenPom:
		return v.parseMavenPom()
	default:
		return v.Scan()
	}
}

// Scan attempts to parse the raw version string into the Version struct fields.
func (v *Version) Scan() error {
	v.Major, v.Minor, v.Patch, v.Alpha, v.Beta, v.RC, v.Preview = 0, 0, 0, 0, 0, 0, 0
	for _, t := range formsInOrder {
		var n int
		var err error
		rawStr := string(v.raw)
		switch t {
		case FormA:
			n, err = fmt.Sscanf(rawStr, t, &v.Major, &v.Minor, &v.Patch)
		case FormB:
			n, err = fmt.Sscanf(rawStr, t, &v.Major, &v.Minor, &v.Patch, &v.Alpha)
		case FormC:
			n, err = fmt.Sscanf(rawStr, t, &v.Major, &v.Minor, &v.Patch, &v.Beta)
		case FormD:
			n, err = fmt.Sscanf(rawStr, t, &v.Major, &v.Minor, &v.Patch, &v.RC)
		case FormE:
			n, err = fmt.Sscanf(rawStr, t, &v.Major, &v.Minor, &v.Patch, &v.Beta, &v.Alpha)
		case FormF:
			n, err = fmt.Sscanf(rawStr, t, &v.Major, &v.Minor, &v.Patch, &v.Preview)
		}
		if err == nil && n == Forms[t] {
			return nil
		}
		switch t {
		case FormG:
			n, err = fmt.Sscanf(rawStr, t, &v.Major, &v.Minor, &v.Patch)
		case FormH:
			n, err = fmt.Sscanf(rawStr, t, &v.Major, &v.Minor)
		case FormI:
			n, err = fmt.Sscanf(rawStr, t, &v.Major)
		case FormJ:
			n, err = fmt.Sscanf(rawStr, t, &v.Major, &v.Minor)
		}

		if err != nil || n == Forms[t] {
			return fmt.Errorf("err in Scan() t = %s ; n = %d ; err = %s", t, n, err)
		}
	}
	return fmt.Errorf("unrecognized version format: \"%s\"", string(v.raw))
}

func (v *Version) Fix() error {
	content := strings.Clone(string(v.raw))
	if len(content) == 0 {
		v.Major = 0
		v.Minor = 0
		v.Patch = 0
		v.RC = 0
		v.Alpha = 0
		v.Beta = 1
		v.Preview = 0
		return nil
	}

	// Pattern: 1.24-beta.1 -> v1.24.0-beta.1
	if matches := reTwoPartPre.FindStringSubmatch(content); len(matches) > 0 {
		majorI, err := strconv.Atoi(matches[1])
		if err != nil {
			return fmt.Errorf("could not parse major version: %w", err)
		}
		minorI, err := strconv.Atoi(matches[2])
		if err != nil {
			return fmt.Errorf("could not parse minor version: %w", err)
		}
		patchI, err := strconv.Atoi(matches[3])
		if err != nil {
			return fmt.Errorf("could not parse patch version: %w", err)
		}
		v.Major = majorI
		v.Minor = minorI
		v.Patch = patchI
		v.RC = 0
		v.Alpha = 0
		v.Beta = 0
		v.Preview = 0
		return nil
	}

	// Pattern: 1.24.5 -> v1.24.5
	if reThreePart.MatchString(content) {
		v.raw = []byte("v" + content)
		return v.Scan()
	}

	// Pattern: 1.24 -> v1.24.0
	if reTwoPart.MatchString(content) {
		v.raw = []byte("v" + content + ".0")
		return v.Scan()
	}

	// Pattern for "1.24.x-beta-q" -> "v1.24.0-beta.1"
	if matches := reFuzzy.FindStringSubmatch(content); len(matches) > 0 {
		v.raw = []byte(fmt.Sprintf("v%s.%s.0-beta.1", matches[1], matches[2]))
		return v.Scan()
	}
	return nil
}

// Format returns a string of the value parsed
func (v *Version) Format() string {
	base := filepath.Base(v.path)
	switch base {
	case FileVersion:
		return v.String()
	case FilePackageJson:
		v.useForm = FormG
		return v.String()
	case FileDockerfile:
		return v.String()
	case FileGoMod:
		return v.String()
	case FileMavenPom:
		return v.String()
	case FileHelmChart:
		return v.String()
	default:
		return fmt.Sprintf("invalid path \"%s\"", v.path)
	}
}

// Save writes the current version string to the specified file path.
func (v *Version) Save(path string) error {
	v.path = path
	base := filepath.Base(v.path)
	switch base {
	case FileVersion:
		return v.saveVersion()
	case FilePackageJson:
		return v.savePackageJson()
	case FileDockerfile:
		return v.saveDockerfile()
	case FileGoMod:
		return v.saveGoMod()
	case FileMavenPom:
		return v.saveMavenPom()
	case FileHelmChart:
		return v.saveHelmChart()
	default:
		return fmt.Errorf("invalid path \"%s\"", v.path)
	}
}

func (v *Version) saveVersion() error {
	return os.WriteFile(v.path, []byte(v.String()), 0644)
}

func (v *Version) savePackageJson() error {
	if v.parsed == nil {
		return fmt.Errorf("parsed map is nil for raw: %s", v.raw)
	}
	v.parsed["version"] = v.String()
	output, err := json.MarshalIndent(v.parsed, "", "  ")
	if err != nil {
		return fmt.Errorf("could not marshal version info: %s", err)
	}
	return os.WriteFile(v.path, output, 0644)
}

func (v *Version) saveHelmChart() error {
	return nil
}

func (v *Version) saveDockerfile() error {
	return nil
}

func (v *Version) saveGoMod() error {
	return nil
}

func (v *Version) saveMavenPom() error {
	return nil
}

func (v *Version) validateFormA(rawStr string, major, minor, patch *int) error {
	defer func() {
		v.useForm = FormA
	}()
	n, err := fmt.Sscanf(rawStr, FormA, major, minor, patch)
	if err == nil && n == Forms[FormA] {
		return nil
	}
	if *major == 0 && *minor == 0 && *patch == 0 {
		return fmt.Errorf("failed to parse type %s from %s", FormA, rawStr)
	}
	if err != nil {
		return fmt.Errorf("failed to scan string %s using %s threw err: %w", rawStr, FormA, err)
	}
	return fmt.Errorf("scan should have returned %d items, got %d", Forms[FormA], n)
}

func (v *Version) validateFormB(rawStr string, major, minor, patch, alpha *int) error {
	defer func() {
		v.useForm = FormB
	}()
	n, err := fmt.Sscanf(rawStr, FormB, major, minor, patch, alpha)
	if err == nil && n == Forms[FormB] {
		return nil
	}
	if *major == 0 && *minor == 0 && *patch == 0 && *alpha == 0 {
		return fmt.Errorf("failed to parse type %s from %s", FormB, rawStr)
	}
	if err != nil {
		return fmt.Errorf("failed to scan string %s using %s threw err: %w", rawStr, FormB, err)
	}
	return fmt.Errorf("scan should have returned %d items, got %d", Forms[FormB], n)
}

func (v *Version) validateFormC(rawStr string, major, minor, patch, beta *int) error {
	defer func() {
		v.useForm = FormC
	}()
	n, err := fmt.Sscanf(rawStr, FormC, major, minor, patch, beta)
	if err == nil && n == Forms[FormC] {
		return nil
	}
	if *major == 0 && *minor == 0 && *patch == 0 && *beta == 0 {
		return fmt.Errorf("failed to parse type %s from %s", FormC, rawStr)
	}
	if err != nil {
		return fmt.Errorf("failed to scan string %s using %s threw err: %w", rawStr, FormC, err)
	}
	return fmt.Errorf("scan should have returned %d items, got %d", Forms[FormC], n)
}

func (v *Version) validateFormD(rawStr string, major, minor, patch, rc *int) error {
	defer func() {
		v.useForm = FormD
	}()
	n, err := fmt.Sscanf(rawStr, FormD, major, minor, patch, rc)
	if err == nil && n == Forms[FormD] {
		return nil
	}
	if *major == 0 && *minor == 0 && *patch == 0 && *rc == 0 {
		return fmt.Errorf("failed to parse type %s from %s", FormD, rawStr)
	}
	if err != nil {
		return fmt.Errorf("failed to scan string %s using %s threw err: %w", rawStr, FormD, err)
	}
	return fmt.Errorf("scan should have returned %d items, got %d", Forms[FormD], n)
}

func (v *Version) validateFormE(rawStr string, major, minor, patch, beta, alpha *int) error {
	defer func() {
		v.useForm = FormE
	}()
	n, err := fmt.Sscanf(rawStr, FormE, major, minor, patch, beta, alpha)
	if err == nil && n == Forms[FormE] {
		return nil
	}
	if *major == 0 && *minor == 0 && *patch == 0 && *beta == 0 && *alpha == 0 {
		return fmt.Errorf("failed to parse type %s from %s", FormE, rawStr)
	}
	if err != nil {
		return fmt.Errorf("failed to scan string %s using %s threw err: %w", rawStr, FormE, err)
	}
	return fmt.Errorf("scan should have returned %d items, got %d", Forms[FormE], n)
}

func (v *Version) validateFormF(rawStr string, major, minor, patch, preview *int) error {
	defer func() {
		v.useForm = FormF
	}()
	if major == nil {
		return errors.New("major cannot be nil")
	}
	if minor == nil {
		return errors.New("minor cannot be nil")
	}
	if patch == nil {
		return errors.New("patch cannot be nil")
	}
	if preview == nil {
		return errors.New("preview cannot be nil")
	}
	n, err := fmt.Sscanf(rawStr, FormF, &major, minor, patch, &preview)
	if err == nil && n == Forms[FormF] {
		return nil
	}
	if *major == 0 && *minor == 0 && *patch == 0 && *preview == 0 {
		return fmt.Errorf("failed to parse type %s from %s", FormF, rawStr)
	}
	if err != nil {
		return fmt.Errorf("failed to scan string %s using %s threw err: %w", rawStr, FormF, err)
	}
	return fmt.Errorf("scan should have returned %d items, got %d", Forms[FormF], n)
}

func (v *Version) validateFormG(rawStr string, major, minor, patch *int) error {
	defer func() {
		v.useForm = FormG
	}()
	if major == nil {
		return errors.New("major cannot be nil")
	}
	if minor == nil {
		return errors.New("minor cannot be nil")
	}
	if patch == nil {
		return errors.New("patch cannot be nil")
	}
	n, err := fmt.Sscanf(rawStr, FormG, &major, minor, patch)
	if err == nil && n == Forms[FormG] {
		return nil
	}
	if *major == 0 && *minor == 0 && *patch == 0 {
		return fmt.Errorf("failed to parse type %s from %s", FormG, rawStr)
	}
	if err != nil {
		return fmt.Errorf("failed to scan string %s using %s threw err: %w", rawStr, FormG, err)
	}
	return fmt.Errorf("scan should have returned %d items, got %d", Forms[FormG], n)
}

func (v *Version) validateFormH(rawStr string, major, minor *int) error {
	defer func() {
		v.useForm = FormH
	}()
	if major == nil {
		return errors.New("major cannot be nil")
	}
	if minor == nil {
		return errors.New("minor cannot be nil")
	}
	n, err := fmt.Sscanf(rawStr, FormH, major, minor)
	if err == nil && n == Forms[FormH] {
		return nil
	}
	if *major == 0 && *minor == 0 {
		return fmt.Errorf("failed to parse type %s from %s", FormH, rawStr)
	}
	if err != nil {
		return fmt.Errorf("failed to scan string %s using %s threw err: %w", rawStr, FormH, err)
	}
	return fmt.Errorf("scan should have returned %d items, got %d", Forms[FormH], n)
}

func (v *Version) validateFormI(rawStr string, major *int) error {
	defer func() {
		v.useForm = FormI
	}()
	if major == nil {
		return errors.New("major cannot be nil")
	}
	n, err := fmt.Sscanf(rawStr, FormI, major)
	if err == nil && n == Forms[FormI] {
		return nil
	}
	if *major == 0 {
		return fmt.Errorf("failed to parse type %s from %s", FormI, rawStr)
	}
	if err != nil {
		return fmt.Errorf("failed to scan string %s using %s threw err: %w", rawStr, FormI, err)
	}
	return fmt.Errorf("scan should have returned %d items, got %d", Forms[FormI], n)
}

func (v *Version) validateFormJ(rawStr string, major, minor *int) error {
	defer func() {
		v.useForm = FormJ
	}()
	if major == nil {
		return errors.New("major cannot be nil")
	}
	if minor == nil {
		return errors.New("minor cannot be nil")
	}
	n, err := fmt.Sscanf(rawStr, FormJ, major, minor)
	if err == nil && n == Forms[FormJ] {
		return nil
	}
	if *major == 0 && *minor == 0 {
		return fmt.Errorf("failed to parse type %s from %s", FormJ, rawStr)
	}
	if err != nil {
		return fmt.Errorf("failed to scan string %s using %s threw err: %w", rawStr, FormJ, err)
	}
	return fmt.Errorf("scan should have returned %d items, got %d", Forms[FormJ], n)
}

func (v *Version) parseVersion() error {
	return v.Scan()
}

func (v *Version) parsePackageJson() error {
	m := make(map[string]interface{})
	err := json.Unmarshal(v.raw, &m)
	if err != nil {
		return err
	}
	ver, ok := m["version"]
	if !ok {
		return errors.New("version not found")
	}
	vs, ok := ver.(string)
	if !ok {
		return errors.New("version is not a string")
	}
	v.raw = []byte(vs)
	return v.Scan()
}

func (v *Version) parseHelmChart() error {
	m := make(map[string]interface{})
	err := yaml.Unmarshal(v.raw, &m)
	if err != nil {
		return err
	}
	ver, ok := m["version"]
	if !ok {
		return errors.New("version not found")
	}
	vs, ok := ver.(string)
	if !ok {
		return errors.New("version is not a string")
	}
	v.raw = []byte(vs)
	return v.Scan()
}

func (v *Version) parseDockerfile() error {
	return nil
}

// parseGoMod reads the v.raw to find the line "go #.#[.#]" as FormG with FormH as fallback
func (v *Version) parseGoMod() error {
	lines := strings.Split(string(v.raw), "\n")
	lineIndex, goVersion := -1, ""

	for i, line := range lines {
		if strings.HasPrefix(line, "go ") {
			if parts := strings.Fields(line); len(parts) == 2 {
				lineIndex, goVersion = i, parts[1]
				break
			}
		}
	}
	if lineIndex == -1 {
		return fmt.Errorf("could not find 'go' directive in %s", string(v.raw))
	}

	mv := NewVersion()

	errG := v.validateFormG(goVersion, &mv.Major, &mv.Minor, &mv.Patch)
	if errG != nil {
		return errG
	} else {
		errH := v.validateFormH(goVersion, &mv.Major, &mv.Minor)
		if errH != nil {
			return errH
		}
	}

	v.Major = mv.Major
	v.Minor = mv.Minor
	v.Patch = mv.Patch
	v.Alpha = mv.Alpha
	v.Beta = mv.Beta
	v.Preview = mv.Preview
	v.RC = mv.RC

	return nil
}

func (v *Version) parseMavenPom() error {
	return nil
}
