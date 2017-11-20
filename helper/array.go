package helper

// InArray check string in array.
func InArray(needle string, haystack []string) ([]string, bool) {
	if len(haystack) == 0 {
		return haystack, false
	}

	for i, val := range haystack {
		if val == needle {
			haystack = append(haystack[:i], haystack[i+1:]...)
			return haystack, true
		}
	}

	return haystack, false
}

// DiffArray show difference in two array.
func DiffArray(s, t []string) []string {
	v := []string{}
	if len(s) == 0 && len(t) == 0 {
		return []string{}
	}

	if len(s) == 0 {
		return t
	}

	if len(t) == 0 {
		return s
	}

	if len(s) > len(t) {
		s, t = t, s
	}

	for _, val := range s {
		if newT, ok := InArray(val, t); ok {
			t = newT
			continue
		}

		v = append(v, val)
	}

	if len(t) > 0 {
		v = append(v, t...)
	}

	return v
}
