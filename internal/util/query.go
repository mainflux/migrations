package util

import "strconv"

// UpdateLimit is a helper function that updates the limit when querying
// based on the total number of items
func UpdateLimit(total uint64) uint64 {
	// Convert the number to a string and get the length
	numStr := strconv.Itoa(int(total))
	numLen := len(numStr)

	switch numLen {
	case 1:
		return uint64(1)
	case 2:
		return uint64(10)
	case 3:
		return uint64(100)
	case 4:
		return uint64(1000)
	case 5:
		return uint64(10000)
	case 6:
		return uint64(100000)
	case 7:
		return uint64(1000000)
	case 8:
		return uint64(10000000)
	case 9:
		return uint64(100000000)
	case 10:
		return uint64(1000000000)
	default:
		return uint64(100)
	}

}
