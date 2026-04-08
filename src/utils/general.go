package utils

func Difference(all, remove []uint64) []uint64 {
	m := make(map[uint64]bool)

	for _, r := range remove {
		m[r] = true
	}

	var result []uint64
	for _, a := range all {
		if !m[a] {
			result = append(result, a)
		}
	}

	return result
}
