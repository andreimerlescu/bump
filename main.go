package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/andreimerlescu/bump/bump"
)

const (
	envInitOnNotFound = "BUMP_INIT_ON_NOT_FOUND" // ENV create -in if not found
	envDefaultInput   = "BUMP_DEFAULT_INPUT"     // ENV defines default -in
	envNoAlphaBeta    = "BUMP_NO_ALPHA_BETA"     // ENV prevents -alpha -beta combo usage
	envAlwaysWrite    = "BUMP_ALWAYS_WRITE"      // ENV always sets -write
	envNoPreview      = "BUMP_NO_PREVIEW"        // ENV prevents -preview usage
	envAlwaysFix      = "BUMP_ALWAYS_FIX"        // ENV always sets -fix
	envNeverFix       = "BUMP_NEVER_FIX"         // ENV never allow -fix to be applied
	envNoAlpha        = "BUMP_NO_ALPHA"          // ENV prevents -alpha usage
	envNoBeta         = "BUMP_NO_BETA"           // ENV prevents -beta usage
	envNoRC           = "BUMP_NO_RC"             // ENV prevents -rc usage

	VFN = "VERSION"
)

// result stores the Version output as a string
type result struct {
	Version string `json:"version"`
}

// binaryVersionBytes contains the embedded VERSION file's contents
//
//go:embed VERSION
var binaryVersionBytes embed.FS

// binaryCurrentVersion is defined by BinaryVersion() and contains the contents of
// the VERSION file
var binaryCurrentVersion string

// BinaryVersion returns the embedded VERSION file of the igo repository as a string
// and cache that value into binaryCurrentVersion once os.ReadFile is complete
func BinaryVersion() string {
	if len(binaryCurrentVersion) == 0 {
		versionBytes, err := binaryVersionBytes.ReadFile(VFN)
		if err != nil {
			return "v0.0.0"
		}
		binaryCurrentVersion = strings.TrimSpace(string(versionBytes))
	}
	return binaryCurrentVersion
}

var (
	initialInputFile = filepath.Join(".", VFN)

	shouldParse string // flag.StringVar -parse
	inputFile   string // flag.StringVar -in

	showVersion bool // flag.BoolVar -v
	shouldInit  bool // flag.BoolVar -init
	writeInput  bool // flag.BoolVar -write
	showAbout   bool // flag.BoolVar -about
	shouldFix   bool // flag.BoolVar -fix
	checkFile   bool // flag.BoolVar -check
	showEnv     bool // flag.BoolVar -env
	useJson     bool // flag.BoolVar -json
	preview     bool // flag.BoolVar -preview
	major       bool // flag.BoolVar -major
	minor       bool // flag.BoolVar -minor
	patch       bool // flag.BoolVar -patch
	alpha       bool // flag.BoolVar -alpha
	beta        bool // flag.BoolVar -beta
	rc          bool // flag.BoolVar -rc
)

// appEnv renders a KEY=VAL\nKEY=VAL\n string of bump ENV variable customization options
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

// envIs combines os.LookupEnv to strconv.ParseBool
func envIs(name string) bool {
	v, ok := os.LookupEnv(name)
	if !ok {
		return false
	}
	vb, err := strconv.ParseBool(v)
	return err == nil && vb
}

// envVal combines os.LookupEnv to fallback
func envVal(name, fallback string) string {
	if v, ok := os.LookupEnv(name); ok {
		return v
	}
	return fallback
}

// about prints a helpful bump usage description to STDOUT
func about() {
	var out strings.Builder
	out.WriteString("Bump Version: " + BinaryVersion() + "\n")
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

var versionCalls = atomic.Int64{}

func NewVersion() *bump.Version {
	if versionCalls.Load() > 3 {
		return nil
	}
	v := bump.New()
	err := v.LoadFile(inputFile)
	if err != nil {
		if strings.HasSuffix(inputFile, VFN) && os.IsNotExist(err) && (envIs(envInitOnNotFound) || shouldInit) {
			var err2 error
			if len(shouldParse) > 0 {
				err2 = os.WriteFile(inputFile, []byte(shouldParse), 0644)
			} else {
				err2 = os.WriteFile(inputFile, []byte("v0.0.0-beta.1"), 0644)
			}
			if err2 != nil {
				log.Fatal(err2)
			}
			versionCalls.Add(1)
			return NewVersion()
		}
		_, _ = fmt.Fprintln(os.Stderr, "Error reading file:", err)
		os.Exit(1)
	}

	if len(shouldParse) > 0 {
		v.SetRaw([]byte(shouldParse))
	}

	// If -fix is requested, attempt to fix the version before parsing.
	if shouldFix {
		if err := v.Fix(); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "Error fixing version:", err)
			os.Exit(1)
		}
	}

	// Now, parse the (potentially fixed) version string.
	if err = v.Parse(); err != nil {
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
	if envIs(envAlwaysFix) && envIs(envNeverFix) {
		_, _ = fmt.Fprintf(os.Stderr, "env %s and %s cannot be used together", envAlwaysFix, envNeverFix)
		os.Exit(1)
	}
	return v
}

func main() {
	config()
	versionCalls.Store(0)
	version := NewVersion()
	if version == nil {
		fmt.Println("version is nil")
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

	newVersionStr := version.Format(version.NoPrefix() == false)
	version.Version = newVersionStr // For JSON output
	noChange := !strings.EqualFold(originalVersionStr, newVersionStr)
	wasParsed := !strings.EqualFold(newVersionStr, shouldParse)
	wasBumped := noChange || wasParsed || shouldInit

	finish(version, wasBumped, bumpFlags, originalVersionStr, newVersionStr)
}

// config gets the flag environment set up, parses if we are showing version, about, or env.
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
		fmt.Println(BinaryVersion())
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

// validate attempts to count the number of bump commands being executed
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

// run executes the bump commands using the bump package
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

// finish prints the summary output of the bump request
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
			if strings.EqualFold(originalVersion, newVersion) && len(shouldParse) > 0 {
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

// printJson uses check() on the error and prints to STDOUT the Indented JSON output
func printJson(data interface{}) {
	output, err := json.MarshalIndent(data, "", "  ")
	check(err)
	fmt.Println(string(output))
}

// check a variable assigned a func type can be redefined but falls through the logic when an err
// needs to be verified if nil or not. If this is an error, then we'll write to STDERR, otherwise,
// we write to STDOUT.
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
