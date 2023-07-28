package util

const (
	defDivider = 10
	maxLimit   = uint64(1000000)
)

// UpdateLimit is a helper function that updates the limit when querying
// based on the total number of items.
func UpdateLimit(total uint64) uint64 {
	switch {
	case total >= maxLimit:
		return maxLimit
	default:
		return total / defDivider
	}
}
