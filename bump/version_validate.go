package bump

import (
	"errors"
	"fmt"
)

// Validate ranges over formsInOrder to run the internal validateForm<T>() func on the Version struct
func (v *Version) Validate() error {
	v.safety()
	v.mu.RLock()
	defer v.mu.RUnlock()
	var major, minor, patch, preview, alpha, beta, rc int
	for _, t := range formsInOrder {
		rawStr := string(v.raw)
		switch t {
		case FormA:
			return v.validateFormA(rawStr, &major, &minor, &patch)
		case FormB:
			return v.validateFormB(rawStr, &major, &minor, &patch, &alpha)
		case FormC:
			return v.validateFormC(rawStr, &major, &minor, &patch, &beta)
		case FormD:
			return v.validateFormD(rawStr, &major, &minor, &patch, &rc)
		case FormE:
			return v.validateFormE(rawStr, &major, &minor, &patch, &alpha, &beta)
		case FormF:
			return v.validateFormF(rawStr, &major, &minor, &patch, &preview)
		case FormG:
			return v.validateFormG(rawStr, &major, &minor, &patch)
		case FormH:
			return v.validateFormH(rawStr, &major, &minor)
		case FormI:
			return v.validateFormI(rawStr, &major)
		case FormJ:
			return v.validateFormJ(rawStr, &major, &minor)

		}
	}
	return fmt.Errorf("unrecognized version format: \"%s\"", string(v.raw))
}

// validateFormA assigns FormA to the Version and uses fmt.Sscanf on the argument version components for verification
// and upstream access to verified results
func (v *Version) validateFormA(rawStr string, major, minor, patch *int) error {
	defer func() {
		v.useForm = FormA
	}()
	n, err := fmt.Sscanf(rawStr, FormA, major, minor, patch)
	if err == nil && n == Forms[FormA] {
		return nil
	}
	if *major == 0 && *minor == 0 && *patch == 0 {
		return fmt.Errorf("failed to parse type %s from %s", FormA, rawStr)
	}
	if err != nil {
		return fmt.Errorf("failed to scan string %s using %s threw err: %w", rawStr, FormA, err)
	}
	return fmt.Errorf("scan should have returned %d items, got %d", Forms[FormA], n)
}

// validateFormB assigns FormB to the Version and uses fmt.Sscanf on the argument version components for verification
// and upstream access to verified results
func (v *Version) validateFormB(rawStr string, major, minor, patch, alpha *int) error {
	defer func() {
		v.useForm = FormB
	}()
	n, err := fmt.Sscanf(rawStr, FormB, major, minor, patch, alpha)
	if err == nil && n == Forms[FormB] {
		return nil
	}
	if *major == 0 && *minor == 0 && *patch == 0 && *alpha == 0 {
		return fmt.Errorf("failed to parse type %s from %s", FormB, rawStr)
	}
	if err != nil {
		return fmt.Errorf("failed to scan string %s using %s threw err: %w", rawStr, FormB, err)
	}
	return fmt.Errorf("scan should have returned %d items, got %d", Forms[FormB], n)
}

// validateFormC assigns FormC to the Version and uses fmt.Sscanf on the argument version components for verification
// and upstream access to verified results
func (v *Version) validateFormC(rawStr string, major, minor, patch, beta *int) error {
	defer func() {
		v.useForm = FormC
	}()
	n, err := fmt.Sscanf(rawStr, FormC, major, minor, patch, beta)
	if err == nil && n == Forms[FormC] {
		return nil
	}
	if *major == 0 && *minor == 0 && *patch == 0 && *beta == 0 {
		return fmt.Errorf("failed to parse type %s from %s", FormC, rawStr)
	}
	if err != nil {
		return fmt.Errorf("failed to scan string %s using %s threw err: %w", rawStr, FormC, err)
	}
	return fmt.Errorf("scan should have returned %d items, got %d", Forms[FormC], n)
}

// validateFormD assigns FormD to the Version and uses fmt.Sscanf on the argument version components for verification
// and upstream access to verified results
func (v *Version) validateFormD(rawStr string, major, minor, patch, rc *int) error {
	defer func() {
		v.useForm = FormD
	}()
	n, err := fmt.Sscanf(rawStr, FormD, major, minor, patch, rc)
	if err == nil && n == Forms[FormD] {
		return nil
	}
	if *major == 0 && *minor == 0 && *patch == 0 && *rc == 0 {
		return fmt.Errorf("failed to parse type %s from %s", FormD, rawStr)
	}
	if err != nil {
		return fmt.Errorf("failed to scan string %s using %s threw err: %w", rawStr, FormD, err)
	}
	return fmt.Errorf("scan should have returned %d items, got %d", Forms[FormD], n)
}

// validateFormE assigns FormE to the Version and uses fmt.Sscanf on the argument version components for verification
// and upstream access to verified results
func (v *Version) validateFormE(rawStr string, major, minor, patch, beta, alpha *int) error {
	defer func() {
		v.useForm = FormE
	}()
	n, err := fmt.Sscanf(rawStr, FormE, major, minor, patch, beta, alpha)
	if err == nil && n == Forms[FormE] {
		return nil
	}
	if *major == 0 && *minor == 0 && *patch == 0 && *beta == 0 && *alpha == 0 {
		return fmt.Errorf("failed to parse type %s from %s", FormE, rawStr)
	}
	if err != nil {
		return fmt.Errorf("failed to scan string %s using %s threw err: %w", rawStr, FormE, err)
	}
	return fmt.Errorf("scan should have returned %d items, got %d", Forms[FormE], n)
}

// validateFormF assigns FormF to the Version and uses fmt.Sscanf on the argument version components for verification
// and upstream access to verified results
func (v *Version) validateFormF(rawStr string, major, minor, patch, preview *int) error {
	defer func() {
		v.useForm = FormF
	}()
	if major == nil {
		return errors.New("major cannot be nil")
	}
	if minor == nil {
		return errors.New("minor cannot be nil")
	}
	if patch == nil {
		return errors.New("patch cannot be nil")
	}
	if preview == nil {
		return errors.New("preview cannot be nil")
	}
	n, err := fmt.Sscanf(rawStr, FormF, &major, minor, patch, &preview)
	if err == nil && n == Forms[FormF] {
		return nil
	}
	if *major == 0 && *minor == 0 && *patch == 0 && *preview == 0 {
		return fmt.Errorf("failed to parse type %s from %s", FormF, rawStr)
	}
	if err != nil {
		return fmt.Errorf("failed to scan string %s using %s threw err: %w", rawStr, FormF, err)
	}
	return fmt.Errorf("scan should have returned %d items, got %d", Forms[FormF], n)
}

// validateFormG assigns FormG to the Version and uses fmt.Sscanf on the argument version components for verification
// and upstream access to verified results
func (v *Version) validateFormG(rawStr string, major, minor, patch *int) error {
	defer func() {
		v.useForm = FormG
	}()
	if major == nil {
		return errors.New("major cannot be nil")
	}
	if minor == nil {
		return errors.New("minor cannot be nil")
	}
	if patch == nil {
		return errors.New("patch cannot be nil")
	}
	n, err := fmt.Sscanf(rawStr, FormG, &major, minor, patch)
	if err == nil && n == Forms[FormG] {
		return nil
	}
	if *major == 0 && *minor == 0 && *patch == 0 {
		return fmt.Errorf("failed to parse type %s from %s", FormG, rawStr)
	}
	if err != nil {
		return fmt.Errorf("failed to scan string %s using %s threw err: %w", rawStr, FormG, err)
	}
	return fmt.Errorf("scan should have returned %d items, got %d", Forms[FormG], n)
}

// validateFormH assigns FormH to the Version and uses fmt.Sscanf on the argument version components for verification
// and upstream access to verified results
func (v *Version) validateFormH(rawStr string, major, minor *int) error {
	defer func() {
		v.useForm = FormH
	}()
	if major == nil {
		return errors.New("major cannot be nil")
	}
	if minor == nil {
		return errors.New("minor cannot be nil")
	}
	n, err := fmt.Sscanf(rawStr, FormH, major, minor)
	if err == nil && n == Forms[FormH] {
		return nil
	}
	if *major == 0 && *minor == 0 {
		return fmt.Errorf("failed to parse type %s from %s", FormH, rawStr)
	}
	if err != nil {
		return fmt.Errorf("failed to scan string %s using %s threw err: %w", rawStr, FormH, err)
	}
	return fmt.Errorf("scan should have returned %d items, got %d", Forms[FormH], n)
}

// validateFormI assigns FormI to the Version and uses fmt.Sscanf on the argument version components for verification
// and upstream access to verified results
func (v *Version) validateFormI(rawStr string, major *int) error {
	defer func() {
		v.useForm = FormI
	}()
	if major == nil {
		return errors.New("major cannot be nil")
	}
	n, err := fmt.Sscanf(rawStr, FormI, major)
	if err == nil && n == Forms[FormI] {
		return nil
	}
	if *major == 0 {
		return fmt.Errorf("failed to parse type %s from %s", FormI, rawStr)
	}
	if err != nil {
		return fmt.Errorf("failed to scan string %s using %s threw err: %w", rawStr, FormI, err)
	}
	return fmt.Errorf("scan should have returned %d items, got %d", Forms[FormI], n)
}

// validateFormJ assigns FormJ to the Version and uses fmt.Sscanf on the argument version components for verification
// and upstream access to verified results
func (v *Version) validateFormJ(rawStr string, major, minor *int) error {
	defer func() {
		v.useForm = FormJ
	}()
	if major == nil {
		return errors.New("major cannot be nil")
	}
	if minor == nil {
		return errors.New("minor cannot be nil")
	}
	n, err := fmt.Sscanf(rawStr, FormJ, major, minor)
	if err == nil && n == Forms[FormJ] {
		return nil
	}
	if *major == 0 && *minor == 0 {
		return fmt.Errorf("failed to parse type %s from %s", FormJ, rawStr)
	}
	if err != nil {
		return fmt.Errorf("failed to scan string %s using %s threw err: %w", rawStr, FormJ, err)
	}
	return fmt.Errorf("scan should have returned %d items, got %d", Forms[FormJ], n)
}
