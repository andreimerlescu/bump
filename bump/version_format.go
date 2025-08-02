package bump

import "fmt"

// Format returns a formatted version string, allowing control over the 'v' prefix.
func (v *Version) Format(withPrefix bool) string {
	v.safety()
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.format(withPrefix)
}

// format is the internal, lock-free implementation for creating a version string.
func (v *Version) format(withPrefix bool) string {
	v.safety()
	baseFormat := "%d.%d.%d"
	if withPrefix && !v.noPrefix {
		baseFormat = "v%d.%d.%d"
	}

	// For shorthand forms, format back in the same style unless a full version is required.
	if v.useForm != "" {
		switch v.useForm {
		case FormA:
			return fmt.Sprintf(FormA, v.Major, v.Minor, v.Patch)
		case FormB:
			return fmt.Sprintf(FormB, v.Major, v.Minor, v.Patch, v.Alpha)
		case FormC:
			return fmt.Sprintf(FormC, v.Major, v.Minor, v.Patch, v.Beta)
		case FormD:
			return fmt.Sprintf(FormD, v.Major, v.Minor, v.Patch, v.RC)
		case FormE:
			return fmt.Sprintf(FormE, v.Major, v.Minor, v.Patch, v.Alpha, v.Beta)
		case FormF:
			return fmt.Sprintf(FormF, v.Major, v.Minor, v.Patch, v.Preview)
		case FormG:
			return fmt.Sprintf(FormG, v.Major, v.Minor, v.Patch)
		case FormH:
			return fmt.Sprintf(FormH, v.Major, v.Minor)
		case FormI:
			return fmt.Sprintf(FormI, v.Major)
		case FormJ:
			if withPrefix {
				return fmt.Sprintf(FormJ, v.Major, v.Minor)
			}
			return fmt.Sprintf(FormH, v.Major, v.Minor) // FormH not FormJ
		default:
		}
	}

	base := fmt.Sprintf(baseFormat, v.Major, v.Minor, v.Patch)
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
	return fmt.Sprintf("%s%s", base, preRelease)
}
