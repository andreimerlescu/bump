package bump

import (
	"bytes"
	"path/filepath"
)

// Fix attempts to correct a malformed raw value of the Version struct
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
