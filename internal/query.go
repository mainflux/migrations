package util

const (
	defLimit     = uint64(100)
	limit1000    = uint64(1000)
	limit10000   = uint64(10000)
	limit100000  = uint64(100000)
	limit1000000 = uint64(1000000)
)

// UpdateLimit is a helper function that updates the limit when querying
// based on the total number of items.
func UpdateLimit(total uint64) uint64 {
	switch {
	case total <= limit1000:
		return defLimit
	case total <= limit10000:
		return limit1000
	case total <= limit100000:
		return limit10000
	case total <= limit1000000:
		return limit100000
	default:
		return limit1000000
	}
}
