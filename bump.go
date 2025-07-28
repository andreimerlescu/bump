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

const (
	BinaryVersion = "v1.0.2"

	envAlwaysWrite  = "BUMP_ALWAYS_WRITE"
	envDefaultInput = "BUMP_DEFAULT_INPUT"
	envNoAlphaBeta  = "BUMP_NO_ALPHA_BETA"
	envNoAlpha      = "BUMP_NO_ALPHA"
	envNoBeta       = "BUMP_NO_BETA"
	envNoRC         = "BUMP_NO_RC"
	envNoPreview    = "BUMP_NO_PREVIEW"
)

var (
	initialInputFile = filepath.Join(".", "VERSION")

	showAbout    bool
	major        bool
	minor        bool
	patch        bool
	alpha        bool
	beta         bool
	rc           bool
	preview      bool
	useJson      bool
	showVersion  bool
	writeInput   bool
	checkFile    bool
	inputFile    string
	alwaysWrite  bool
	defaultInput string
	noAlphaBeta  bool
	noAlpha      bool
	noBeta       bool
	noRC         bool
	noPreview    bool
)

func about() {
	var out strings.Builder
	out.WriteString("Bump Version: " + BinaryVersion + "\n")
	out.WriteString("Usage:\n")
	out.WriteString("  bump -check\n")
	out.WriteString("  bump -[major|minor|patch|alpha|beta|rc|preview]\n")
	out.WriteString("  bump -[major|minor|patch|alpha|beta|rc|preview] -write\n")
	out.WriteString("  bump -json -[major|minor|patch|alpha|beta|rc|preview]\n")
	out.WriteString("  bump -json -[major|minor|patch|alpha|beta|rc|preview] -write\n")
	out.WriteString("Defaults: \n")
	out.WriteString(fmt.Sprintf("  -in=%s [default: %s]\n", inputFile, defaultInput))
	out.WriteString("Environment Variables:\n")
	out.WriteString(fmt.Sprintf("  ENV[%s]=%s\n", envAlwaysWrite, strconv.FormatBool(alwaysWrite)))
	out.WriteString(fmt.Sprintf("  ENV[%s]=%s\n", envDefaultInput, os.Getenv(envDefaultInput)))
	out.WriteString(fmt.Sprintf("  ENV[%s]=%s\n", envNoBeta, strconv.FormatBool(noBeta)))
	out.WriteString(fmt.Sprintf("  ENV[%s]=%s\n", envNoAlpha, strconv.FormatBool(noAlpha)))
	out.WriteString(fmt.Sprintf("  ENV[%s]=%s\n", envNoAlphaBeta, strconv.FormatBool(noAlphaBeta)))
	out.WriteString(fmt.Sprintf("  ENV[%s]=%s\n", envNoRC, strconv.FormatBool(noRC)))
	out.WriteString(fmt.Sprintf("  ENV[%s]=%s\n", envNoPreview, strconv.FormatBool(noPreview)))
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

	version := bump.Version{}
	// version.ParseFile() opens the inputFile and loads the contents using fmt.Sscanf on the []byte from the contents
	check(version.ParseFile(inputFile))

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
	run(&version)

	// this renders a new output value using the pattern based on the type of the release based on its int values
	newVersion := version.String()

	// copy the rendered output to the struct for output so it can be programmatically accessed with -json
	version.Version = strings.Clone(newVersion)

	// determine if the version was changed or not
	wasBumped := !strings.EqualFold(originalVersion, newVersion)

	// finish() renders the final output to the user
	finish(&version, wasBumped, bumpFlags, originalVersion, newVersion)
}

func config() {
	defaultInput = env(envDefaultInput, initialInputFile)
	flag.StringVar(&inputFile, "in", defaultInput, "input file")
	flag.BoolVar(&showAbout, "about", false, "show about")
	flag.BoolVar(&major, "major", false, "major version bump")
	flag.BoolVar(&minor, "minor", false, "minor version bump")
	flag.BoolVar(&patch, "patch", false, "patch version bump")
	flag.BoolVar(&alpha, "alpha", false, "alpha version bump")
	flag.BoolVar(&beta, "beta", false, "beta version bump")
	flag.BoolVar(&rc, "rc", false, "rc version bump")
	flag.BoolVar(&preview, "preview", false, "preview version bump")
	flag.BoolVar(&useJson, "json", false, "use json version bump")
	flag.BoolVar(&showVersion, "v", false, "show version")
	flag.BoolVar(&writeInput, "write", false, "writeInput version file")
	flag.BoolVar(&checkFile, "check", false, "check version file")
	flag.Parse()
	alwaysWrite = envIs(envAlwaysWrite)
	noAlphaBeta = envIs(envNoAlphaBeta)
	noAlpha = envIs(envNoAlpha)
	noBeta = envIs(envNoBeta)
	noRC = envIs(envNoRC)
	noPreview = envIs(envNoPreview)
	if showVersion {
		fmt.Println(BinaryVersion)
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
		if writeInput || alwaysWrite {
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
	} else if writeInput || alwaysWrite {
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

func env(name, fallback string) string {
	v, ok := os.LookupEnv(name)
	if !ok {
		return fallback
	}
	return v
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

type result struct {
	Version string `json:"version"`
}
