package bump

import "strings"

// BumpMajor is responsible for increasing the Major field in the Version struct
func (v *Version) BumpMajor() {
	v.safety()
	v.mu.Lock()
	defer v.mu.Unlock()
	v.Major++
	v.Minor, v.Patch, v.RC, v.Alpha, v.Beta, v.Preview = 0, 0, 0, 0, 0, 0
}

// BumpMinor is responsible for increasing the Minor field in the Version struct
func (v *Version) BumpMinor() {
	v.safety()
	v.mu.Lock()
	defer v.mu.Unlock()
	v.Minor++
	v.Patch, v.RC, v.Alpha, v.Beta, v.Preview = 0, 0, 0, 0, 0
}

// BumpPatch is responsible for increasing the Patch field in the Version struct
func (v *Version) BumpPatch() {
	v.safety()
	v.mu.Lock()
	defer v.mu.Unlock()
	v.Patch++
	v.RC, v.Alpha, v.Beta, v.Preview = 0, 0, 0, 0
	if !strings.EqualFold(v.useForm, FormG) {
		v.useForm = FormA
	}
}

// BumpRC is responsible for increasing the RC field in the Version struct
func (v *Version) BumpRC() {
	v.safety()
	v.mu.Lock()
	defer v.mu.Unlock()
	v.RC++
	v.Alpha, v.Beta, v.Preview = 0, 0, 0
	v.useForm = FormD
}

// BumpAlpha is responsible for increasing the Alpha field in the Version struct
func (v *Version) BumpAlpha() {
	v.safety()
	v.mu.Lock()
	defer v.mu.Unlock()
	v.Alpha++
	if v.useForm == FormD {
		v.useForm = FormE
	} else {
		v.useForm = FormB
	}
}

// BumpBeta is responsible for increasing the Beta field in the Version struct
func (v *Version) BumpBeta() {
	v.safety()
	v.mu.Lock()
	defer v.mu.Unlock()
	v.Beta++
	v.useForm = FormB
}

// BumpPreview is responsible for increasing the Preview field in the Version struct
func (v *Version) BumpPreview() {
	v.safety()
	v.mu.Lock()
	defer v.mu.Unlock()
	v.Preview++
	v.Patch, v.Alpha, v.Beta, v.RC = 0, 0, 0, 0
	v.useForm = FormF
}
