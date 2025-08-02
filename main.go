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
	BinaryVersion = "v1.0.5"

	envAlwaysWrite    = "BUMP_ALWAYS_WRITE"
	envDefaultInput   = "BUMP_DEFAULT_INPUT"
	envNoAlphaBeta    = "BUMP_NO_ALPHA_BETA"
	envNoAlpha        = "BUMP_NO_ALPHA"
	envNoBeta         = "BUMP_NO_BETA"
	envNoRC           = "BUMP_NO_RC"
	envNoPreview      = "BUMP_NO_PREVIEW"
	envNeverFix       = "BUMP_NEVER_FIX"
	envInitOnNotFound = "BUMP_INIT_ON_NOT_FOUND"
	envAlwaysFix      = "BUMP_ALWAYS_FIX"
)

var (
	initialInputFile = filepath.Join(".", "VERSION")

	inputFile string

	showAbout   bool
	showEnv     bool
	showVersion bool
	major       bool
	minor       bool
	patch       bool
	alpha       bool
	beta        bool
	rc          bool
	shouldParse string
	shouldInit  bool
	shouldFix   bool
	preview     bool
	checkFile   bool
	useJson     bool
	writeInput  bool
)

func appEnv(indent string) string {
	var out strings.Builder
	defaultInput := envVal(envDefaultInput, initialInputFile)
	for e, v := range map[string]string{
		envAlwaysWrite:    strconv.FormatBool(envIs(envAlwaysWrite)),
		envDefaultInput:   envVal(envDefaultInput, defaultInput),
		envNeverFix:       strconv.FormatBool(envIs(envNeverFix)),
		envNoAlpha:        strconv.FormatBool(envIs(envNoAlpha)),
		envNoBeta:         strconv.FormatBool(envIs(envNoBeta)),
		envNoAlphaBeta:    strconv.FormatBool(envIs(envNoAlphaBeta)),
		envNoRC:           strconv.FormatBool(envIs(envNoRC)),
		envNoPreview:      strconv.FormatBool(envIs(envNoPreview)),
		envInitOnNotFound: strconv.FormatBool(envIs(envInitOnNotFound)),
		envAlwaysFix:      strconv.FormatBool(envIs(envAlwaysFix)),
	} {
		out.WriteString(fmt.Sprintf("%s%s=%s\n", indent, e, v))
	}
	return out.String()
}

func envIs(name string) bool {
	v, ok := os.LookupEnv(name)
	if !ok {
		return false
	}
	vb, err := strconv.ParseBool(v)
	return err == nil && vb
}

func envVal(name, fallback string) string {
	if v, ok := os.LookupEnv(name); ok {
		return v
	}
	return fallback
}

func about() {
	var out strings.Builder
	out.WriteString("Bump Version: " + BinaryVersion + "\n")
	out.WriteString("Usage:\n")
	out.WriteString("  bump -check [-in=FILE]\n")
	out.WriteString("  bump -fix [-write] [-in=FILE]\n")
	out.WriteString("  bump -[major|minor|patch|alpha|beta|rc|preview] [-write] [-in=FILE] [-json]\n")
	out.WriteString("Supported File Types:\n")
	for _, t := range bump.SupportedFiles {
		out.WriteString(fmt.Sprintf("  %s\n", t))
	}
	out.WriteString("Defaults: \n")
	out.WriteString(fmt.Sprintf("  -in=%s [default: %s]\n", inputFile, envVal(envDefaultInput, initialInputFile)))
	out.WriteString("Environment Variables:\n")
	out.WriteString(appEnv("  "))
	fmt.Print(out.String())
}

func main() {
	config()

	if envIs(envAlwaysFix) && envIs(envNeverFix) {
		_, _ = fmt.Fprintf(os.Stderr, "env %s and %s cannot be used together", envAlwaysFix, envNeverFix)
		os.Exit(1)
	}

retry:
	version := bump.New()
	err := version.LoadFile(inputFile)
	if err != nil {
		if strings.HasSuffix(inputFile, "VERSION") && os.IsNotExist(err) && (envIs(envInitOnNotFound) || shouldInit) {
			var err2 error
			if len(shouldParse) > 0 {
				err2 = os.WriteFile(inputFile, []byte(shouldParse), 0644)
			} else {
				err2 = os.WriteFile(inputFile, []byte("v0.0.0-beta.1"), 0644)
			}
			if err2 != nil {
				log.Fatal(err2)
			}
			goto retry
		}
		_, _ = fmt.Fprintln(os.Stderr, "Error reading file:", err)
		os.Exit(1)
	}

	if len(shouldParse) > 0 {
		version.SetRaw([]byte(shouldParse))
	}

	// If -fix is requested, attempt to fix the version before parsing.
	if shouldFix {
		if err := version.Fix(); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "Error fixing version:", err)
			os.Exit(1)
		}
	}

	// Now, parse the (potentially fixed) version string.
	if err = version.Parse(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Error parsing version:", err)
		// If fix was not requested but might have helped, suggest it.
		if !shouldFix {
			v2 := bump.New()
			_ = v2.LoadFile(inputFile)
			if v2.Fix() == nil {
				_, _ = fmt.Fprintln(os.Stderr, "Hint: the version string may be fixable with the -fix flag.")
			}
		}
		os.Exit(1)
	}

	originalVersionStr := version.Format(!version.NoPrefix())

	if checkFile {
		if useJson {
			printJson(&result{Version: originalVersionStr})
		} else {
			fmt.Println(originalVersionStr)
		}
		os.Exit(0)
	}

	bumpFlags, err := validate()
	check(err)

	if bumpFlags > 0 {
		run(version)
	}

	newVersionStr := version.Format(!version.NoPrefix())
	version.Version = newVersionStr // For JSON output
	wasBumped := !strings.EqualFold(originalVersionStr, newVersionStr) || !strings.EqualFold(newVersionStr, shouldParse) || shouldInit

	finish(version, wasBumped, bumpFlags, originalVersionStr, newVersionStr)
}

func config() {
	// input actions
	defaultInput := envVal(envDefaultInput, initialInputFile)
	flag.StringVar(&inputFile, "in", defaultInput, fmt.Sprintf("input file (default: %s or BUMP_DEFAULT_INPUT)", initialInputFile))
	flag.StringVar(&shouldParse, "parse", "", "use value as input of new VERSION file")

	// information actions
	flag.BoolVar(&showVersion, "v", false, "show binary version")
	flag.BoolVar(&showAbout, "about", false, "show about")
	flag.BoolVar(&showEnv, "env", false, "show environment variables")

	// bump actions
	flag.BoolVar(&major, "major", false, "major version bump")
	flag.BoolVar(&minor, "minor", false, "minor version bump")
	flag.BoolVar(&patch, "patch", false, "patch version bump")
	flag.BoolVar(&alpha, "alpha", false, "alpha version bump")
	flag.BoolVar(&beta, "beta", false, "beta version bump")
	flag.BoolVar(&rc, "rc", false, "rc version bump")
	flag.BoolVar(&preview, "preview", false, "preview version bump")

	// flow control actions
	flag.BoolVar(&useJson, "json", false, "use json output")
	flag.BoolVar(&writeInput, "write", envIs(envAlwaysWrite), "write version back to file")
	flag.BoolVar(&checkFile, "check", false, "check version file and print it")
	flag.BoolVar(&shouldFix, "fix", envIs(envAlwaysFix), "fix malformed version string if possible")
	flag.BoolVar(&shouldInit, "init", envIs(envInitOnNotFound), "initialize version file")
	flag.Parse()

	if showVersion {
		fmt.Println(BinaryVersion)
		os.Exit(0)
	}
	if showEnv {
		fmt.Print(appEnv(""))
		os.Exit(0)
	}
	if showAbout {
		about()
		os.Exit(0)
	}

	// Load env-controlled settings
	if envIs(envNeverFix) {
		shouldFix = false
	}

}

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
	// Pre-release bumps can be combined with major/minor/patch, but not with each other.
	preReleaseFlags := 0
	if alpha && !envIs(envNoAlpha) && !envIs(envNoAlphaBeta) {
		preReleaseFlags++
	}
	if beta && !envIs(envNoBeta) && !envIs(envNoAlphaBeta) {
		preReleaseFlags++
	}
	if rc && !envIs(envNoRC) {
		preReleaseFlags++
	}
	if preview && !envIs(envNoPreview) {
		preReleaseFlags++
	}

	if bumpFlags > 1 {
		return 0, fmt.Errorf("only one of -major, -minor, or -patch can be used at a time")
	}
	if preReleaseFlags > 1 {
		// Exception: allow alpha and beta to be combined
		if !(preReleaseFlags == 2 && (alpha && beta)) {
			return 0, fmt.Errorf("only one pre-release bump can be used at a time (e.g., -alpha, -beta)")
		}
	}
	return bumpFlags + preReleaseFlags, nil
}

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
	if rc && !envIs(envNoRC) {
		version.BumpRC()
	}
	if beta && !envIs(envNoBeta) && !envIs(envNoAlphaBeta) {
		version.BumpBeta()
	}
	if alpha && !envIs(envNoAlpha) && !envIs(envNoAlphaBeta) {
		version.BumpAlpha()
	}
	if preview && !envIs(envNoPreview) {
		version.BumpPreview()
	}
}

func finish(version *bump.Version, wasBumped bool, bumpFlags int, originalVersion, newVersion string) {
	if useJson {
		printJson(version)
		if writeInput {
			check(version.Save(inputFile))
		}
		return
	}

	if !wasBumped && len(shouldParse) > 0 {
		wasBumped = true
		bumpFlags++
	}

	if shouldInit && strings.EqualFold(originalVersion, newVersion) {
		if writeInput {
			check(version.Save(inputFile))
			fmt.Printf("Initialized %s (saved to %s)\n", originalVersion, inputFile)
		} else {
			fmt.Printf("Initialized %s\n", originalVersion)
		}
	} else if wasBumped {
		if writeInput {
			check(version.Save(inputFile))
			if strings.EqualFold(originalVersion, newVersion) {
				fmt.Printf("Parsed %s (saved to %s)\n", newVersion, inputFile)
			} else {
				fmt.Printf("Bumped %s → %s (saved to %s)\n", originalVersion, newVersion, inputFile)
			}
		} else {
			fmt.Printf("Bumped %s → %s\n", originalVersion, newVersion)
		}
	} else if writeInput && shouldFix {
		check(version.Save(inputFile))
		fmt.Printf("Fixed and saved version %s to %s\n", newVersion, inputFile)
	} else if bumpFlags == 0 && !checkFile {
		if len(shouldParse) > 0 && !strings.EqualFold(originalVersion, newVersion) {
			check(version.Save(inputFile))
		} else {
			fmt.Println("No bump operation specified. Use -major, -minor, -patch, etc., to bump the version.")
		}
		fmt.Printf("Current version is: %s\n", originalVersion)
	} else {
		// No bump occurred (e.g., -alpha on a version that is already an alpha)
		fmt.Printf("Version is %s (no change)\n", originalVersion)
	}
}

func printJson(data interface{}) {
	output, err := json.MarshalIndent(data, "", "  ")
	check(err)
	fmt.Println(string(output))
}

var check = func(this any) {
	switch t := this.(type) {
	case error:
		if this != nil {
			log.Fatal(t)
		}
	case string:
		if len(t) > 0 {
			fmt.Println(t)
		}
	default:
		if t != nil {
			fmt.Println(t)
		}
	}
}
