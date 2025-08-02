package bump

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

// Save passes through based on the filepath.Base(path) provided
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

// saveVersion passes into os.WriteFile on the path and the v.format(!v.noPrefix) with 0644 permissions
func (v *Version) saveVersion() error {
	return os.WriteFile(v.path, []byte(v.format(!v.noPrefix)), 0644)
}

// savePackageJson assigns to the parsed value the v.format(false) to the "version" key before json.MarshalIndent the
// output before sending that into os.WriteFile on the path provided
func (v *Version) savePackageJson() error {
	if v.parsed == nil {
		return errors.New("cannot save package.json: file was not parsed")
	}
	v.useForm = FormG
	v.parsed["version"] = v.format(false)
	output, err := json.MarshalIndent(v.parsed, "", "  ")
	if err != nil {
		return fmt.Errorf("could not marshal version info: %w", err)
	}
	return os.WriteFile(v.path, output, 0644)
}

// saveHelmChart assigns tot he parsed value of the v.format(false) to the "version", uses yaml.NewEncoder on those bytes
// and then uses os.WriteFile to save the bytes to the path provided, the full parsed output
func (v *Version) saveHelmChart() error {
	if v.parsed == nil {
		return errors.New("cannot save Chart.yaml: file was not parsed")
	}
	v.useForm = ""
	newVersion := v.format(false)
	v.parsed["version"] = newVersion
	//  if _, ok := v.parsed["appVersion"]; ok {
	//		v.parsed["appVersion"] = newVersion
	//	}
	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(2)
	if err := encoder.Encode(v.parsed); err != nil {
		return fmt.Errorf("could not marshal helm chart: %w", err)
	}
	return os.WriteFile(v.path, buf.Bytes(), 0644)
}

// saveDockerfile replaces in raw using regex reDockerfileVersion to replace the LABEL provided in the Dockerfile file
// before running os.WriteFile on the provided path
func (v *Version) saveDockerfile() error {
	v.useForm = ""
	newVersion := v.format(true)
	newContent := reDockerfileVersion.ReplaceAll(v.raw, []byte("${1}"+newVersion+"${3}"))
	return os.WriteFile(v.path, newContent, 0644)
}

// saveGoMod replaces in raw using regex reGoModVersion to replace the "go #.#[.#]" with the new version before sending
// that to os.WriteFile on the provided path
func (v *Version) saveGoMod() error {
	v.useForm = FormG
	newVersion := v.format(false)
	newContent := reGoModVersion.ReplaceAll(v.raw, []byte("${1}"+newVersion))
	return os.WriteFile(v.path, newContent, 0644)
}

// saveMavenPom replaces in raw using regex reMavenVersion to replace the <project>..<version> tag in pom.xml before
// sending it to os.WriteFile on the provided path
func (v *Version) saveMavenPom() error {
	v.useForm = ""
	newVersion := v.format(false)
	if !reMavenVersion.Match(v.raw) {
		return errors.New("could not find <project>...<version> tag in pom.xml to update")
	}
	newContent := reMavenVersion.ReplaceAll(v.raw, []byte("${1}"+newVersion+"${3}"))
	return os.WriteFile(v.path, newContent, 0644)
}
