package utils

func Contains(ids []uint, candidate uint) bool {
	for _, id := range ids {
		if id == candidate {
			return true
		}
	}
	return false
}

func IsEqualPrice(a, b *float64) bool {
	switch {
	case a == nil && b == nil:
		return true
	case a == nil || b == nil:
		return false
	default:
		return *a == *b
	}
}
