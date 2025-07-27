package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/andreimerlescu/bump/bump"
)

const BinaryVersion = "v1.0.0"

var (
	inputFile   = flag.String("in", filepath.Join(".", "VERSION"), "input file")
	major       = flag.Bool("major", false, "major version bump")
	minor       = flag.Bool("minor", false, "minor version bump")
	patch       = flag.Bool("patch", false, "patch version bump")
	alpha       = flag.Bool("alpha", false, "alpha version bump")
	beta        = flag.Bool("beta", false, "beta version bump")
	rc          = flag.Bool("rc", false, "rc version bump")
	preview     = flag.Bool("preview", false, "preview version bump")
	useJson     = flag.Bool("json", false, "use json version bump")
	showVersion = flag.Bool("v", false, "show version")
	writeInput  = flag.Bool("write", false, "writeInput version file")
	checkFile   = flag.Bool("check", false, "check version file")
)

func main() {
	flag.Parse()
	if *showVersion {
		fmt.Println(BinaryVersion)
		os.Exit(0)
	}
	version := bump.Version{}
	check(version.ParseFile(*inputFile))
	originalVersion := version.Raw()
	if *checkFile {
		if *useJson {
			fmt.Printf(`{"version": "%s"}`, originalVersion)
		} else {
			fmt.Printf("%s\n", originalVersion)
		}
		os.Exit(0)
	}

	bumpFlags, err := validate()
	check(err)

	if *major {
		version.BumpMajor()
	}

	if *minor {
		version.BumpMinor()
	}

	if *patch {
		version.BumpPatch()
	}

	if *rc {
		version.BumpRC()
	}

	if *beta {
		version.BumpBeta()
	}

	if *alpha {
		version.BumpAlpha()
	}

	if *preview {
		version.BumpPreview()
	}

	newVersion := version.String()
	wasBumped := !strings.EqualFold(originalVersion, newVersion)

	if wasBumped {
		if *writeInput {
			check(version.Save(*inputFile))
			if *useJson {
				output, err := json.MarshalIndent(version, "", "  ")
				check(err)
				fmt.Println(string(output))
			} else {
				fmt.Printf("Bumped %s → %s (saved to %s)\n", originalVersion, newVersion, *inputFile)
			}
		} else {
			if *useJson {
				output, err := json.MarshalIndent(version, "", "  ")
				check(err)
				fmt.Println(string(output))
			} else {
				fmt.Printf("Bumped %s → %s\n", originalVersion, newVersion)
			}
		}
	} else if *writeInput {
		check(version.Save(*inputFile))
		if *useJson {
			output, err := json.MarshalIndent(version, "", "  ")
			check(err)
			fmt.Println(string(output))
		} else {
			fmt.Printf("Re-saved version %s to %s\n", newVersion, *inputFile)
		}
	} else if bumpFlags == 0 && !*checkFile && !*writeInput {
		if *useJson {
			output, err := json.MarshalIndent(version, "", "  ")
			check(err)
			fmt.Println(string(output))
		} else {
			fmt.Println("No bump operation specified. Use -major, -minor, -patch, etc. to bump the version.")
			fmt.Printf("Current version is: %s\n", originalVersion)
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

func validate() (int, error) {
	bumpFlags := 0
	if *major {
		bumpFlags++
	}
	if *minor {
		bumpFlags++
	}
	if *patch {
		bumpFlags++
	}
	if *alpha {
		bumpFlags++
	}
	if *beta {
		bumpFlags++
	}
	if *rc {
		bumpFlags++
	}
	if *preview {
		bumpFlags++
	}

	if bumpFlags > 1 {
		return 0, fmt.Errorf("only one bump operation can be used at a time")
	}
	return bumpFlags, nil
}
