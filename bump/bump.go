package bump

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

type Version struct {
	mu       *sync.RWMutex
	parsed   map[string]interface{} // contains unmarshal'd json|yaml|toml|ini key=>value pairs
	path     string
	raw      []byte
	noPrefix bool
	useForm  string

	Major   int    `json:"major"`
	Minor   int    `json:"minor"`
	Patch   int    `json:"patch"`
	Alpha   int    `json:"alpha"`
	Beta    int    `json:"beta"`
	RC      int    `json:"rc"`
	Preview int    `json:"preview"`
	Version string `json:"version"`
}

func New() *Version {
	return &Version{
		parsed: make(map[string]interface{}),
		mu:     &sync.RWMutex{},
	}
}

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

func (v *Version) safety() {
	if v.mu == nil {
		v.mu = &sync.RWMutex{}
	}
	if v.parsed == nil {
		v.parsed = make(map[string]interface{})
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
	FormG string = "%d.%d.%d" // SemVer Standard (no v-prefix)
	FormH string = "%d.%d"    // Shorthand SemVer (e.g., 1.24)
	FormI string = "v%d"      // v# (v1 -> v100)
	FormJ string = "v%d.%d"   // v#.# (v1.1 -> v100.100)

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
	reTwoPart   = regexp.MustCompile(`^(\d+)\.(\d+)$`)
	reThreePart = regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)$`)
	reFuzzy     = regexp.MustCompile(`^v?(\d+)\.(\d+).*`)

	// Regex for file-specific parsing/saving
	reDockerfileVersion = regexp.MustCompile(`(LABEL\s+(?:org\.label-schema\.version|version)=")([^"]+)(")`)
	reGoModVersion      = regexp.MustCompile(`(go\s+)([0-9.]+)`)
	reMavenVersion      = regexp.MustCompile(`(?s)(<project.*?>.*?<version>)(.*?)(</version>)`)
)

// Forms is a map of format strings to the expected number of scanned items.
var Forms = map[string]int{
	FormA: 3, FormB: 4, FormC: 4, FormD: 4, FormE: 5, FormF: 4,
	FormG: 3, FormH: 2, FormI: 1, FormJ: 2,
}

// formsInOrder defines a deterministic order for scanning, from most specific to least.
var formsInOrder = []string{FormE, FormB, FormC, FormD, FormF, FormA, FormG, FormH, FormJ, FormI}

// String formats the version struct into a standardized string with a 'v' prefix.
func (v *Version) String() string {
	v.safety()
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.format(true)
}

// Format returns a formatted version string, allowing control over the 'v' prefix.
func (v *Version) Format(withPrefix bool) string {
	v.safety()
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.format(withPrefix)
}

// format is the internal, lock-free implementation for creating a version string.
func (v *Version) format(withPrefix bool) string {
	v.safety()
	baseFormat := "%d.%d.%d"
	if withPrefix && !v.noPrefix {
		baseFormat = "v%d.%d.%d"
	}

	// For shorthand forms, format back in the same style unless a full version is required.
	if v.useForm != "" {
		switch v.useForm {
		case FormA:
			return fmt.Sprintf(FormA, v.Major, v.Minor, v.Patch)
		case FormB:
			return fmt.Sprintf(FormB, v.Major, v.Minor, v.Patch, v.Alpha)
		case FormC:
			return fmt.Sprintf(FormC, v.Major, v.Minor, v.Patch, v.Beta)
		case FormD:
			return fmt.Sprintf(FormD, v.Major, v.Minor, v.Patch, v.RC)
		case FormE:
			return fmt.Sprintf(FormE, v.Major, v.Minor, v.Patch, v.Alpha, v.Beta)
		case FormF:
			return fmt.Sprintf(FormF, v.Major, v.Minor, v.Patch, v.Preview)
		case FormG:
			return fmt.Sprintf(FormG, v.Major, v.Minor, v.Patch)
		case FormH:
			return fmt.Sprintf(FormH, v.Major, v.Minor)
		case FormI:
			return fmt.Sprintf(FormI, v.Major)
		case FormJ:
			if withPrefix {
				return fmt.Sprintf(FormJ, v.Major, v.Minor)
			}
			return fmt.Sprintf(FormH, v.Major, v.Minor) // FormH not FormJ
		default:
		}
	}

	base := fmt.Sprintf(baseFormat, v.Major, v.Minor, v.Patch)
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

func (v *Version) Raw() string {
	v.safety()
	return string(v.raw)
}

func (v *Version) NoPrefix() bool {
	v.safety()
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.noPrefix
}

func (v *Version) Compare(other *Version) int {
	v.safety()
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

	vIsPre := v.Preview > 0 || v.RC > 0 || v.Beta > 0 || v.Alpha > 0
	otherIsPre := other.Preview > 0 || other.RC > 0 || other.Beta > 0 || other.Alpha > 0

	if !vIsPre && otherIsPre {
		return 1
	}
	if vIsPre && !otherIsPre {
		return -1
	}
	if !vIsPre && !otherIsPre {
		return 0
	}

	// Both are pre-releases, compare them
	if v.Preview != other.Preview {
		return compareInt(v.Preview, other.Preview)
	}
	if v.RC != other.RC {
		return compareInt(v.RC, other.RC)
	}
	if v.Beta != other.Beta {
		return compareInt(v.Beta, other.Beta)
	}
	if v.Alpha != other.Alpha {
		return compareInt(v.Alpha, other.Alpha)
	}
	return 0
}

func (v *Version) SetRaw(raw []byte) {
	v.safety()
	v.mu.Lock()
	defer v.mu.Unlock()
	v.raw = raw
}

func (v *Version) BumpMajor() {
	v.safety()
	v.mu.Lock()
	defer v.mu.Unlock()
	v.Major++
	v.Minor, v.Patch, v.RC, v.Alpha, v.Beta, v.Preview = 0, 0, 0, 0, 0, 0
}

func (v *Version) BumpMinor() {
	v.safety()
	v.mu.Lock()
	defer v.mu.Unlock()
	v.Minor++
	v.Patch, v.RC, v.Alpha, v.Beta, v.Preview = 0, 0, 0, 0, 0
}

func (v *Version) BumpPatch() {
	v.safety()
	v.mu.Lock()
	defer v.mu.Unlock()
	v.Patch++
	v.RC, v.Alpha, v.Beta, v.Preview = 0, 0, 0, 0
	v.useForm = FormA
}

func (v *Version) BumpRC() {
	v.safety()
	v.mu.Lock()
	defer v.mu.Unlock()
	v.RC++
	v.Alpha, v.Beta, v.Preview = 0, 0, 0
	v.useForm = FormD
}

func (v *Version) BumpAlpha() {
	v.safety()
	v.mu.Lock()
	defer v.mu.Unlock()
	v.Alpha++
	if v.useForm == FormD {
		v.useForm = FormE
	} else {
		v.useForm = FormB
	}
}

func (v *Version) BumpBeta() {
	v.safety()
	v.mu.Lock()
	defer v.mu.Unlock()
	v.Beta++
	v.useForm = FormB
}

func (v *Version) BumpPreview() {
	v.safety()
	v.mu.Lock()
	defer v.mu.Unlock()
	v.Preview++
	v.Patch, v.Alpha, v.Beta, v.RC = 0, 0, 0, 0
	v.useForm = FormF
}

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

func (v *Version) ParseFile(path string) error {
	v.safety()
	if err := v.LoadFile(path); err != nil {
		return err
	}
	return v.Parse()
}

func (v *Version) Parse() error {
	v.safety()
	v.mu.Lock()
	defer v.mu.Unlock()
	base := filepath.Base(v.path)
	parseableContent := bytes.TrimSpace(v.raw)
	return v.parse(base, parseableContent)
}

func (v *Version) parse(base string, content []byte) error {
	var err error
	switch base {
	case FileVersion:
		err = v.parseVersion(content)
	case FilePackageJson:
		err = v.parsePackageJson(content)
	case FileHelmChart:
		err = v.parseHelmChart(content)
	case FileDockerfile:
		err = v.parseDockerfile(content)
	case FileGoMod:
		err = v.parseGoMod(content)
	case FileMavenPom:
		err = v.parseMavenPom(content)
	default:
		err = v.scan(content)
	}
	return err
}

func (v *Version) parseVersion(content []byte) error {
	return v.scan(content)
}

func (v *Version) parsePackageJson(content []byte) error {
	m := make(map[string]interface{})
	if err := json.Unmarshal(content, &m); err != nil {
		return err
	}
	v.parsed = m
	ver, ok := m["version"]
	if !ok {
		return errors.New("version key not found in package.json")
	}
	vs, ok := ver.(string)
	if !ok {
		return errors.New("version is not a string in package.json")
	}
	return v.scan([]byte(vs))
}

func (v *Version) parseHelmChart(content []byte) error {
	m := make(map[string]interface{})
	if err := yaml.Unmarshal(content, &m); err != nil {
		return err
	}
	v.parsed = m
	ver, ok := m["version"]
	if !ok {
		return errors.New("version key not found in Chart.yaml")
	}
	vs, ok := ver.(string)
	if !ok {
		return errors.New("version is not a string in Chart.yaml")
	}
	return v.scan([]byte(vs))
}

func (v *Version) parseDockerfile(content []byte) error {
	matches := reDockerfileVersion.FindSubmatch(content)
	if len(matches) < 3 {
		return errors.New("could not find version LABEL in Dockerfile")
	}
	return v.scan(matches[2])
}

func (v *Version) parseGoMod(content []byte) error {
	matches := reGoModVersion.FindSubmatch(content)
	if len(matches) < 3 {
		return errors.New("could not find 'go' directive in go.mod")
	}
	return v.scan(matches[2])
}

func (v *Version) parseMavenPom(content []byte) error {
	matches := reMavenVersion.FindSubmatch(content)
	if len(matches) < 3 {
		return errors.New("could not find <version> tag inside <project> in pom.xml")
	}
	return v.scan(matches[2])
}

func (v *Version) scan(raw []byte) error {
	v.Major, v.Minor, v.Patch, v.Alpha, v.Beta, v.RC, v.Preview = 0, 0, 0, 0, 0, 0, 0

	rawStr := string(raw)
	for _, t := range formsInOrder {
		var n int
		var err error
		tempV := &Version{}

		switch t {
		case FormA:
			n, err = fmt.Sscanf(rawStr, t, &tempV.Major, &tempV.Minor, &tempV.Patch)
		case FormB:
			n, err = fmt.Sscanf(rawStr, t, &tempV.Major, &tempV.Minor, &tempV.Patch, &tempV.Alpha)
		case FormC:
			n, err = fmt.Sscanf(rawStr, t, &tempV.Major, &tempV.Minor, &tempV.Patch, &tempV.Beta)
		case FormD:
			n, err = fmt.Sscanf(rawStr, t, &tempV.Major, &tempV.Minor, &tempV.Patch, &tempV.RC)
		case FormE:
			n, err = fmt.Sscanf(rawStr, t, &tempV.Major, &tempV.Minor, &tempV.Patch, &tempV.Beta, &tempV.Alpha)
		case FormF:
			n, err = fmt.Sscanf(rawStr, t, &tempV.Major, &tempV.Minor, &tempV.Patch, &tempV.Preview)
		case FormG:
			n, err = fmt.Sscanf(rawStr, t, &tempV.Major, &tempV.Minor, &tempV.Patch)
		case FormH:
			n, err = fmt.Sscanf(rawStr, t, &tempV.Major, &tempV.Minor)
		case FormI:
			n, err = fmt.Sscanf(rawStr, t, &tempV.Major)
		case FormJ:
			n, err = fmt.Sscanf(rawStr, t, &tempV.Major, &tempV.Minor)
		}

		if err == nil && n == Forms[t] {
			v.Major, v.Minor, v.Patch = tempV.Major, tempV.Minor, tempV.Patch
			v.Alpha, v.Beta, v.RC, v.Preview = tempV.Alpha, tempV.Beta, tempV.RC, tempV.Preview
			v.useForm = t
			v.noPrefix = strings.HasPrefix(t, "%d")
			return nil
		}
	}
	return fmt.Errorf("unrecognized version format: \"%s\"", rawStr)
}

func (v *Version) Fix() error {
	v.safety()
	v.mu.Lock()
	defer v.mu.Unlock()

	if (v.Major > 0 || v.Minor > 0 || v.Patch > 0) && len(v.raw) == 0 {
		v.useForm = FormA
		v.raw = []byte(v.format(v.noPrefix))
		return v.scan(v.raw)
	}

	content := string(bytes.TrimSpace(v.raw))
	if content == "" {
		v.raw = []byte("v0.0.1")
		v.useForm = FormA
		return v.scan(v.raw)
	}

	var fixedContent string
	if reThreePart.MatchString(content) {
		fixedContent = "v" + content
	} else if reTwoPart.MatchString(content) {
		fixedContent = "v" + content + ".0"
	}

	if fixedContent != "" {
		v.raw = []byte(fixedContent)
		v.useForm = ""
		return v.scan(v.raw)
	}

	base := filepath.Base(v.path)
	parseableContent := bytes.TrimSpace(v.raw)
	return v.parse(base, parseableContent)
}

func (v *Version) Save(path string) error {
	v.safety()
	v.mu.Lock()
	defer v.mu.Unlock()
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
		return v.saveVersion()
	}
}

func (v *Version) saveVersion() error {
	return os.WriteFile(v.path, []byte(v.format(!v.noPrefix)), 0644)
}

func (v *Version) savePackageJson() error {
	if v.parsed == nil {
		return errors.New("cannot save package.json: file was not parsed")
	}
	v.useForm = "" // Ensure full version is written
	v.parsed["version"] = v.format(false)
	output, err := json.MarshalIndent(v.parsed, "", "  ")
	if err != nil {
		return fmt.Errorf("could not marshal version info: %w", err)
	}
	return os.WriteFile(v.path, output, 0644)
}

func (v *Version) saveHelmChart() error {
	if v.parsed == nil {
		return errors.New("cannot save Chart.yaml: file was not parsed")
	}
	v.useForm = ""
	newVersion := v.format(false)
	v.parsed["version"] = newVersion
	if _, ok := v.parsed["appVersion"]; ok {
		v.parsed["appVersion"] = newVersion
	}
	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(2)
	if err := encoder.Encode(v.parsed); err != nil {
		return fmt.Errorf("could not marshal helm chart: %w", err)
	}
	return os.WriteFile(v.path, buf.Bytes(), 0644)
}

func (v *Version) saveDockerfile() error {
	v.useForm = ""
	newVersion := v.format(true)
	newContent := reDockerfileVersion.ReplaceAll(v.raw, []byte("${1}"+newVersion+"${3}"))
	return os.WriteFile(v.path, newContent, 0644)
}

func (v *Version) saveGoMod() error {
	v.useForm = FormG
	newVersion := v.format(false)
	newContent := reGoModVersion.ReplaceAll(v.raw, []byte("${1}"+newVersion))
	return os.WriteFile(v.path, newContent, 0644)
}

func (v *Version) saveMavenPom() error {
	v.useForm = ""
	newVersion := v.format(false)
	if !reMavenVersion.Match(v.raw) {
		return errors.New("could not find <project>...<version> tag in pom.xml to update")
	}
	newContent := reMavenVersion.ReplaceAll(v.raw, []byte("${1}"+newVersion+"${3}"))
	return os.WriteFile(v.path, newContent, 0644)
}

func (v *Version) Validate() error {
	v.safety()
	v.mu.RLock()
	defer v.mu.RUnlock()
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

func compareInt(a, b int) int {
	if a > b {
		return 1
	}
	if a < b {
		return -1
	}
	return 0
}
