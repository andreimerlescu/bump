package bump

import "regexp"

// Version formats
const (
	FormA string = "v%d.%d.%d"                  // Prefix SemVer
	FormB string = "v%d.%d.%d-alpha.%d"         // Alpha Prefix SemVer
	FormC string = "v%d.%d.%d-beta.%d"          // Beta Prefix SemVer
	FormD string = "v%d.%d.%d-rc.%d"            // RC Prefix SemVer
	FormE string = "v%d.%d.%d-beta.%d-alpha.%d" // Alpha Beta Prefix SemVer
	FormF string = "v%d.%d.%d-preview.%d"       // Preview Prefix SemVer
	FormG string = "%d.%d.%d"                   // SemVer Standard (no v-prefix)
	FormH string = "%d.%d"                      // Shorthand SemVer (e.g., 1.24)
	FormI string = "v%d"                        // v# (v1 -> v100)
	FormJ string = "v%d.%d"                     // v#.# (v1.1 -> v100.100)

	FileVersion     string = "VERSION"      // Full Contents Replaced
	FilePackageJson string = "package.json" // Key "version" Replaced
	FileMavenPom    string = "pom.xml"      // Key "version" Replaced
	FileHelmChart   string = "Chart.yaml"   // Key "version" Replaced
	FileDockerfile  string = "Dockerfile"   // Label "version" Replaced
	FileGoMod       string = "go.mod"       // Line 3, aka "go #.#[.#]" Replaced
)

// SupportedFiles can be passed into `-in` when running bump
var SupportedFiles = []string{
	FileVersion,
	FilePackageJson,
	FileMavenPom,
	FileHelmChart,
	FileDockerfile,
	FileGoMod,
}

var (
	reTwoPart   = regexp.MustCompile(`^(\d+)\.(\d+)$`)        // Two Part Version Only
	reThreePart = regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)$`) // Three Part Version Only
	reFuzzy     = regexp.MustCompile(`^v?(\d+)\.(\d+).*`)     // Prefixed Two Part Version Only

	// Regex for file-specific parsing/saving
	reDockerfileVersion = regexp.MustCompile(`(LABEL\s+(?:org\.label-schema\.version|version)=")([^"]+)(")`) // Dockerfile Label
	reGoModVersion      = regexp.MustCompile(`(go\s+)([0-9.]+)`)                                             // Go Mod Version
	reMavenVersion      = regexp.MustCompile(`(?s)(<project.*?>.*?<version>)(.*?)(</version>)`)              // Maven Version
)

// Forms is a map of format strings to the expected number of scanned items.
var Forms = map[string]int{
	FormE: 5, // 1:major 2:minor 3:patch 4:beta 5:alpha
	FormB: 4, // 1:major 2:minor 3:patch 4:alpha
	FormC: 4, // 1:major 2:minor 3:patch 4:beta
	FormD: 4, // 1:major 2:minor 3:patch 4:rc
	FormF: 4, // 1:major 2:minor 3:patch 4:preview
	FormA: 3, // 1:major 2:minor 3:patch
	FormG: 3, // 1:major 2:minor 3:patch
	FormH: 2, // 1:major 2:minor
	FormJ: 2, // 1:major 2:minor
	FormI: 1, // 1:major
}

// formsInOrder defines a deterministic order for scanning, from most specific to least.
var formsInOrder = []string{FormE, FormB, FormC, FormD, FormF, FormA, FormG, FormH, FormJ, FormI}
