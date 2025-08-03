package bump

import (
	"bytes"
	"encoding/json"
	"errors"
	"gopkg.in/yaml.v3"
	"path/filepath"
)

// ParseFile uses LoadFile on the path to return Parse()
func (v *Version) ParseFile(path string) error {
	v.safety()
	if err := v.LoadFile(path); err != nil {
		return err
	}
	return v.Parse()
}

// Parse trims the byte spaces of the raw field and captures the base of the path before passing both into the internal parse func
func (v *Version) Parse() error {
	v.safety()
	v.mu.Lock()
	defer v.mu.Unlock()
	base := filepath.Base(v.path)
	parseableContent := bytes.TrimSpace(v.raw)
	return v.parse(base, parseableContent)
}

// parse switches on the provided base to look for as File<Kind> ie FileVersion, FileGoMod, etc. and run the subsequent
// v.parse<Kind>() func with the provided []byte content, otherwise, we'll just v.scan() the content
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

// parseVersion passes through to v.scan()
func (v *Version) parseVersion(content []byte) error {
	return v.scan(content)
}

// parsePackageJson (json) unmarshal's the provided []bytes into a map[string]interface{}, extracts the "version" key and
// returns v.scan() on the []byte contents of the value of "version" key in the map
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

// parseHelmChart (yaml) unmarshal's the provided []bytes into a map[string]interface{}, extracts the "version" key and
// returns v.scan() on the []byte contents of the value of the "version" key in the map
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

// parseDockerfile (text) uses regex reDockerfileVersion to FindSubmatch on the content and return v.scan(matches[2]) of the result
func (v *Version) parseDockerfile(content []byte) error {
	matches := reDockerfileVersion.FindSubmatch(content)
	if len(matches) < 3 {
		return errors.New("could not find version LABEL in Dockerfile")
	}
	return v.scan(matches[2])
}

// parseGoMod (text) uses regex reGoModVersion to FindSubmatch on the content and return v.scan(matches[2]) of the result
func (v *Version) parseGoMod(content []byte) error {
	matches := reGoModVersion.FindSubmatch(content)
	if len(matches) < 3 {
		return errors.New("could not find 'go' directive in go.mod")
	}
	// matches[0] "go 1.24"
	// matches[1] = "go "
	// matches[2] = "1.24"
	v.useForm = FormG
	return v.scan(matches[2])
}

// parseMavenPom (text) uses regex reMavenVersion to FindSubmatch on the content and return v.scan(matches[2]) of the result
func (v *Version) parseMavenPom(content []byte) error {
	matches := reMavenVersion.FindSubmatch(content)
	if len(matches) < 3 {
		return errors.New("could not find <version> tag inside <project> in pom.xml")
	}
	return v.scan(matches[2])
}

func (v *Version) parseIgo() error {
	igoVersion, err := currentIgoVersion()
	if err != nil {
		return err
	}
	return v.parseVersion([]byte(igoVersion))
}
