package bump

// compareInt is used to determine if a > b then 1 or a < b then -1 or 0
func compareInt(a, b int) int {
	if a > b {
		return 1
	}
	if a < b {
		return -1
	}
	return 0
}
