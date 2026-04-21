package bump

import (
	"fmt"
	"strings"
)

// Format returns a formatted version string, allowing control over the 'v' prefix.
func (v *Version) Format(withPrefix bool) string {
	v.safety()
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.format(withPrefix)
}

// format is the internal, lock-free implementation for creating a version string.
func (v *Version) format(withPrefix bool) string {
	if strings.HasSuffix(v.path, "go.mod") {
		v.useForm = FormG
	}

	// base is either "v%d.%d.%d" or "%d.%d.%d" depending on withPrefix and noPrefix.
	// All switch cases below build on this rather than using the Form* constants
	// directly, because Form* constants hardcode a 'v' prefix and would ignore
	// the withPrefix argument.
	base := "%d.%d.%d"
	if withPrefix && !v.noPrefix {
		base = "v%d.%d.%d"
	}

	if v.useForm != "" {
		switch v.useForm {
		case FormA:
			return fmt.Sprintf(base, v.Major, v.Minor, v.Patch)
		case FormB:
			return fmt.Sprintf(base+"-alpha.%d", v.Major, v.Minor, v.Patch, v.Alpha)
		case FormC:
			return fmt.Sprintf(base+"-beta.%d", v.Major, v.Minor, v.Patch, v.Beta)
		case FormD:
			return fmt.Sprintf(base+"-rc.%d", v.Major, v.Minor, v.Patch, v.RC)
		case FormE:
			return fmt.Sprintf(base+"-beta.%d-alpha.%d", v.Major, v.Minor, v.Patch, v.Beta, v.Alpha)
		case FormF:
			return fmt.Sprintf(base+"-preview.%d", v.Major, v.Minor, v.Patch, v.Preview)
		case FormG:
			// FormG is always no-prefix regardless of withPrefix (go.mod style: "1.2.3")
			return fmt.Sprintf(FormG, v.Major, v.Minor, v.Patch)
		case FormH:
			// FormH is always no-prefix (shorthand: "1.24")
			return fmt.Sprintf(FormH, v.Major, v.Minor)
		case FormI:
			// FormI always has a 'v' prefix ("v1") — it's part of the format itself
			return fmt.Sprintf(FormI, v.Major)
		case FormJ:
			if withPrefix {
				return fmt.Sprintf(FormJ, v.Major, v.Minor)
			}
			return fmt.Sprintf(FormH, v.Major, v.Minor)
		default:
		}
	}

	// Fallback: useForm is empty, build from fields directly.
	result := fmt.Sprintf(base, v.Major, v.Minor, v.Patch)
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
	return result + preRelease
}
