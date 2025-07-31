package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/andreimerlescu/bump/bump"
)

type result struct {
	Version string `json:"version"`
}

const (
	BinaryVersion = "v1.0.5-beta.1"

	envAlwaysWrite  = "BUMP_ALWAYS_WRITE"
	envDefaultInput = "BUMP_DEFAULT_INPUT"
	envNoAlphaBeta  = "BUMP_NO_ALPHA_BETA"
	envNoAlpha      = "BUMP_NO_ALPHA"
	envNoBeta       = "BUMP_NO_BETA"
	envNoRC         = "BUMP_NO_RC"
	envNoPreview    = "BUMP_NO_PREVIEW"
	envNeverFix     = "BUMP_NEVER_FIX"
)

var (
	initialInputFile = filepath.Join(".", "VERSION")

	setEnvVal    string
	defaultInput string
	inputFile    string

	showAbout   bool
	showEnv     bool
	gomod       bool
	showVersion bool
	major       bool
	neverFix    bool
	minor       bool
	noAlpha     bool
	noBeta      bool
	noRC        bool
	noPreview   bool
	patch       bool
	alpha       bool
	beta        bool
	rc          bool
	shouldFix   bool
	preview     bool
	checkFile   bool
	noAlphaBeta bool
	useJson     bool
	writeInput  bool
)

func appEnv(indent string) string {
	var out strings.Builder
	for e, v := range map[string]string{
		envAlwaysWrite:  strconv.FormatBool(envIs(envAlwaysWrite)),
		envDefaultInput: envVal(envDefaultInput, defaultInput),
		envNeverFix:     strconv.FormatBool(envIs(envNeverFix)),
		envNoAlpha:      strconv.FormatBool(envIs(envNoAlpha)),
		envNoBeta:       strconv.FormatBool(envIs(envNoBeta)),
		envNoAlphaBeta:  strconv.FormatBool(envIs(envNoAlphaBeta)),
		envNoRC:         strconv.FormatBool(envIs(envNoRC)),
		envNoPreview:    strconv.FormatBool(envIs(envNoPreview)),
	} {
		out.WriteString(fmt.Sprintf("%s%s=%s\n", indent, e, envVal(e, v)))
	}
	return out.String()
}

func envIs(name string) bool {
	v, ok := os.LookupEnv(name)
	if !ok {
		return false
	}
	vb, err := strconv.ParseBool(v)
	if err != nil {
		return false
	}
	return vb
}

func envVal(name, fallback string) string {
	v, ok := os.LookupEnv(name)
	if !ok {
		return fallback
	}
	return v
}

func about() {
	var out strings.Builder
	out.WriteString("Bump Version: " + BinaryVersion + "\n")
	out.WriteString("Usage:\n")
	out.WriteString("  bump -fix [-write] [-in=FILE]\n")
	out.WriteString("  bump -fix -gomod [-write] [-in=go.fixGoMod]\n")
	out.WriteString("  bump -check\n")
	out.WriteString("  bump -[major|minor|patch|alpha|beta|rc|preview]\n")
	out.WriteString("  bump -[major|minor|patch|alpha|beta|rc|preview] -write\n")
	out.WriteString("  bump -json -[major|minor|patch|alpha|beta|rc|preview]\n")
	out.WriteString("  bump -json -[major|minor|patch|alpha|beta|rc|preview] -write\n")
	out.WriteString("Supported File Types:\n")
	for _, t := range bump.SupportedFiles {
		out.WriteString(fmt.Sprintf("  %s\n", t))
	}
	out.WriteString("Defaults: \n")
	out.WriteString(fmt.Sprintf("  -in=%s [default: %s]\n", inputFile, defaultInput))
	out.WriteString("Environment Variables:\n")
	out.WriteString(appEnv("  "))
	out.WriteString("ORDER | Format\n")
	out.WriteString("------|------------------------------\n")
	for p := 0; p < len(bump.Priority); p++ {
		t := bump.Priority[p]
		if len(t) == 0 {
			continue
		}
		ps := fmt.Sprintf("% 5d", p)
		out.WriteString(fmt.Sprintf("%s | %s\n", ps, t))
	}
	fmt.Print(out.String())
}

func main() {
	// config sets up the initial flags for the binary
	config()

	if gomod {
		fixGoMod()
		return
	}

	version := bump.NewVersion()
	// version.ParseFile() opens the inputFile and loads the contents using fmt.Sscanf on the []byte from the contents
	err := version.ParseFile(inputFile)
	if err != nil {
		if !shouldFix || neverFix {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		fixErr := version.Fix()
		if fixErr != nil {
			_, _ = fmt.Fprintln(os.Stderr, fixErr.Error())
			os.Exit(1)
		}

		fmt.Println(version.String())

		if writeInput {
			writeErr := version.Save(inputFile)
			if writeErr != nil {
				_, _ = fmt.Fprintln(os.Stderr, "Error writing file:", writeErr)
				os.Exit(1)
			}
		}
		os.Exit(0) // Fix is a terminal operation.
	}

	// version.Raw() provides the string value of the []byte from the os.ReadFile(inputFile)
	originalVersion := version.Raw()

	if checkFile {
		if useJson {
			r := &result{}
			r.Version = originalVersion
			output, err := json.MarshalIndent(r, "", "   ")
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(output))
		} else {
			fmt.Printf("%s\n", originalVersion)
		}
		os.Exit(0)
	}

	// validate() ensures that you aren't combining -major -minor together in a -write operation
	bumpFlags, err := validate()
	check(err)

	// run() performs the ++ on the int values within the bump.Version struct itself
	run(version)

	// this renders a new output value using the pattern based on the type of the release based on its int values
	newVersion := version.String()

	// copy the rendered output to the struct for output so it can be programmatically accessed with -json
	version.Version = strings.Clone(newVersion)

	// determine if the version was changed or not
	wasBumped := !strings.EqualFold(originalVersion, newVersion)

	// finish() renders the final output to the user
	finish(version, wasBumped, bumpFlags, originalVersion, newVersion)
}

func config() {
	defaultInput = envVal(envDefaultInput, initialInputFile)
	flag.StringVar(&inputFile, "in", defaultInput, "input file")
	flag.BoolVar(&showAbout, "about", false, "show about")
	flag.BoolVar(&showEnv, "appEnv", false, "show appEnv")
	flag.StringVar(&setEnvVal, "set", "", "set appEnv to new value")
	flag.BoolVar(&major, "major", false, "major version bump")
	flag.BoolVar(&minor, "minor", false, "minor version bump")
	flag.BoolVar(&patch, "patch", false, "patch version bump")
	flag.BoolVar(&alpha, "alpha", false, "alpha version bump")
	flag.BoolVar(&beta, "beta", false, "beta version bump")
	flag.BoolVar(&rc, "rc", false, "rc version bump")
	flag.BoolVar(&preview, "preview", false, "preview version bump")
	flag.BoolVar(&useJson, "json", false, "use json version bump")
	flag.BoolVar(&showVersion, "v", false, "show version")
	flag.BoolVar(&writeInput, "write", envIs(envAlwaysWrite), "writeInput version file")
	flag.BoolVar(&checkFile, "check", false, "check version file")
	flag.BoolVar(&shouldFix, "fix", false, "fix version file")
	flag.BoolVar(&gomod, "gomod", false, "handle input as a go.fixGoMod file")
	flag.Parse()
	neverFix = envIs(envNeverFix)
	if neverFix {
		shouldFix = false
	}
	noAlphaBeta = envIs(envNoAlphaBeta)
	noAlpha = envIs(envNoAlpha)
	noBeta = envIs(envNoBeta)
	noRC = envIs(envNoRC)
	noPreview = envIs(envNoPreview)
	if showVersion {
		fmt.Println(BinaryVersion)
		os.Exit(0)
	}
	if showEnv {
		fmt.Println(appEnv(""))
		os.Exit(0)
	}
	if showAbout {
		about()
		os.Exit(0)
	}
}

// validate ensures that flags like -major and -minor aren't being combined. only -alpha and -beta can be combined
func validate() (int, error) {
	bumpFlags := 0
	if major {
		bumpFlags++
	}
	if minor {
		bumpFlags++
	}
	if patch {
		bumpFlags++
	}
	if alpha && (!noAlpha || !noAlphaBeta) {
		bumpFlags++
	}
	if beta && (!noBeta || !noAlphaBeta) {
		bumpFlags++
	}
	if rc && !noRC {
		bumpFlags++
	}
	if preview && !noPreview {
		bumpFlags++
	}

	if bumpFlags > 1 {
		return 0, fmt.Errorf("only one bump operation can be used at a time")
	}
	return bumpFlags, nil
}

// run performs the version.BumpMajor(), version.BumpMinor(), etc. func calls on the version itself based on the flags
func run(version *bump.Version) {
	if major {
		version.BumpMajor()
	}

	if minor {
		version.BumpMinor()
	}

	if patch {
		version.BumpPatch()
	}

	if rc && !noRC {
		version.BumpRC()
	}

	if beta && (!noBeta || !noAlphaBeta) {
		version.BumpBeta()
	}

	if alpha && (!noAlpha || !noAlphaBeta) {
		version.BumpAlpha()
	}

	if preview && !noPreview {
		version.BumpPreview()
	}
}

// finish renders the final output to the user respecting their choice of -json and -write
func finish(version *bump.Version, wasBumped bool, bumpFlags int, originalVersion, newVersion string) {
	if wasBumped {
		if writeInput {
			check(version.Save(inputFile))
			if useJson {
				output, err := json.MarshalIndent(version, "", "  ")
				check(err)
				fmt.Println(string(output))
			} else {
				fmt.Printf("Bumped %s → %s (saved to %s)\n", originalVersion, newVersion, inputFile)
			}
		} else {
			if useJson {
				output, err := json.MarshalIndent(version, "", "  ")
				check(err)
				fmt.Println(string(output))
			} else {
				fmt.Printf("Bumped %s → %s\n", originalVersion, newVersion)
			}
		}
	} else if writeInput {
		check(version.Save(inputFile))
		if useJson {
			output, err := json.MarshalIndent(version, "", "  ")
			check(err)
			fmt.Println(string(output))
		} else {
			fmt.Printf("Re-saved version %s to %s\n", newVersion, inputFile)
		}
	} else if bumpFlags == 0 && !checkFile && !writeInput {
		if useJson {
			output, err := json.MarshalIndent(version, "", "  ")
			check(err)
			fmt.Println(string(output))
		} else {
			fmt.Println("No bump operation specified. Use -major, -minor, -patch, etc. to bump the version.")
			fmt.Printf("Current version is: %s\n", originalVersion)
		}
	}
}

func fixPackageJSON(path string) error {
	bumpFlags, _ := validate()
	isBumpAttempted := bumpFlags > 0 || alpha || beta || rc || preview
	if isBumpAttempted {
		_, _ = fmt.Fprintln(os.Stderr, "error: bump commands (-major, -minor, etc.) are ineligible with the -gomod flag")
		os.Exit(1)
	}

	ver := bump.NewVersion()
	err := ver.ParseFile(path)
	if err != nil {
		return fmt.Errorf("bump.Version.ParseFile(%s) threw err: %w", path, err)
	}

	return nil

}

func fixGoMod() {
	bumpFlags, _ := validate()
	isBumpAttempted := bumpFlags > 0 || alpha || beta || rc || preview
	if isBumpAttempted {
		_, _ = fmt.Fprintln(os.Stderr, "error: bump commands (-major, -minor, etc.) are ineligible with the -gomod flag")
		os.Exit(1)
	}

	content, err := os.ReadFile(inputFile)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "error reading file:", err)
		os.Exit(1)
	}

	lines := strings.Split(string(content), "\n")
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
		_, _ = fmt.Fprintln(os.Stderr, "error: could not find 'go' directive in", inputFile)
		os.Exit(1)
	}

	if len(strings.Split(goVersion, ".")) != 2 {
		if !checkFile {
			fmt.Println(goVersion)
		}
		os.Exit(0) // Format is valid (e.g., 1.24.5), so exit 0.
	}

	// Format is `x.y`, which is considered invalid for check or requires a fix.
	if checkFile {
		os.Exit(1)
	}

	if !shouldFix {
		_, _ = fmt.Fprintln(os.Stderr, "error: go.fixGoMod version is in short format (e.g., 1.24), run with -fix to correctContents")
		os.Exit(1)
	}

	// --- Fix logic for go.fixGoMod ---
	fixedVersion := ""
	goVersionPath := filepath.Join(os.Getenv("HOME"), "go", "version")
	if _, err := os.Stat(goVersionPath); err == nil {
		if b, err := os.ReadFile(goVersionPath); err == nil {
			versionFromFile := strings.TrimSpace(string(b))
			if strings.HasPrefix(versionFromFile, goVersion+".") {
				fixedVersion = versionFromFile
			}
		}
	}

	if fixedVersion == "" { // Fallback if file doesn't provide a fix
		if goVersion == "1.24" {
			fixedVersion = "1.24.5"
		} else {
			_, _ = fmt.Fprintln(os.Stderr, "error: no fix available for go version", goVersion)
			os.Exit(1)
		}
	}

	fmt.Println(fixedVersion)
	if writeInput {
		lines[lineIndex] = "go " + fixedVersion
		newContent := strings.Join(lines, "\n")
		if err := os.WriteFile(inputFile, []byte(newContent), 0644); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "error writing to file:", err)
			os.Exit(1)
		}
	}
}

var check = func(what interface{}) {
	switch q := what.(type) {
	case error:
		if q != nil {
			log.Fatal(q)
		}
	default:
		if q != nil {
			log.Printf("%v\n", q)
		}
	}
}
