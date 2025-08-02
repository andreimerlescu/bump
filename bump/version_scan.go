package bump

import (
	"fmt"
	"strings"
)

// scan attempts to take a raw []byte and use formsInOrder to fmt.Sscanf that raw string value. If the tempV scan is
// successful, and Forms (of the formsInOrder as (t)) matches the number of assignments of the version components
func (v *Version) scan(raw []byte) error {
	v.Major, v.Minor, v.Patch, v.Alpha, v.Beta, v.RC, v.Preview = 0, 0, 0, 0, 0, 0, 0

	rawStr := string(raw)
	for _, t := range formsInOrder {
		var n int
		var err error
		tempV := &Version{}

		switch t {
		case FormA:
			n, err = fmt.Sscanf(rawStr, t, &tempV.Major, &tempV.Minor, &tempV.Patch)
		case FormB:
			n, err = fmt.Sscanf(rawStr, t, &tempV.Major, &tempV.Minor, &tempV.Patch, &tempV.Alpha)
		case FormC:
			n, err = fmt.Sscanf(rawStr, t, &tempV.Major, &tempV.Minor, &tempV.Patch, &tempV.Beta)
		case FormD:
			n, err = fmt.Sscanf(rawStr, t, &tempV.Major, &tempV.Minor, &tempV.Patch, &tempV.RC)
		case FormE:
			n, err = fmt.Sscanf(rawStr, t, &tempV.Major, &tempV.Minor, &tempV.Patch, &tempV.Beta, &tempV.Alpha)
		case FormF:
			n, err = fmt.Sscanf(rawStr, t, &tempV.Major, &tempV.Minor, &tempV.Patch, &tempV.Preview)
		case FormG:
			n, err = fmt.Sscanf(rawStr, t, &tempV.Major, &tempV.Minor, &tempV.Patch)
		case FormH:
			n, err = fmt.Sscanf(rawStr, t, &tempV.Major, &tempV.Minor)
		case FormI:
			n, err = fmt.Sscanf(rawStr, t, &tempV.Major)
		case FormJ:
			n, err = fmt.Sscanf(rawStr, t, &tempV.Major, &tempV.Minor)
		}

		if err == nil && n == Forms[t] {
			v.Major, v.Minor, v.Patch = tempV.Major, tempV.Minor, tempV.Patch
			v.Alpha, v.Beta, v.RC, v.Preview = tempV.Alpha, tempV.Beta, tempV.RC, tempV.Preview
			v.useForm = t
			v.noPrefix = strings.HasPrefix(t, "%d")
			return nil
		}
	}
	return fmt.Errorf("unrecognized version format: \"%s\"", rawStr)
}
