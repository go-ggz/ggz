package helper

// InArray check string in array.
func InArray(needle string, haystack []string) bool {
	for _, val := range haystack {
		if val == needle {
			return true
		}
	}

	return false
}

// DiffArray show difference in two array.
func DiffArray(a, b []string) []string {
	var s, t, v []string
	if len(a) == 0 && len(b) == 0 {
		return []string{}
	}

	if len(a) == 0 {
		return b
	}

	if len(b) == 0 {
		return a
	}

	if len(a) > len(b) {
		s = a
		t = b
	} else {
		s = b
		t = a
	}

	for _, val := range s {
		if len(v)+len(t) == len(s) {
			continue
		}

		if ok := InArray(val, t); ok {
			continue
		}

		v = append(v, val)
	}

	return v
}
