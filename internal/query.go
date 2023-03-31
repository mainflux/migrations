package util

const defLimit = uint64(100)

// UpdateLimit is a helper function that updates the limit when querying
// based on the total number of items
func UpdateLimit(total uint64) uint64 {
	if total <= 1000 {
		return defLimit
	}
	if total <= 10000 {
		return uint64(1000)
	}
	if total <= 100000 {
		return uint64(10000)
	}
	if total <= 1000000 {
		return uint64(100000)
	}
	return uint64(1000000)
}
