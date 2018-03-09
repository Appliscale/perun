package helpers

func SliceContains(list []string, a string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func IsPlainMap(mp map[string]interface{}) bool {
	// First we check is it more complex. If so - it is worth investigating and we should stop checking.
	for _, m := range mp {
		if _, ok := m.(map[string]interface{}); ok {
			return false
		} else if _, ok := m.([]interface{}); ok {
			return false
		}
	}
	// Ok, it isn't. So is there any <nil>?
	if mapContainsNil(mp) { // Yes, it is - so it is a map worth investigating. This is not the map we're looking for.
		return false
	}

	return true // There is no <nil> and no complexity - it is a plain, non-nil map.
}

func IsPlainSlice(slc []interface{}) bool {
	// The same flow as in `isPlainMap` function.
	for _, s := range slc {
		if _, ok := s.(map[string]interface{}); ok {
			return false
		} else if _, ok := s.([]interface{}); ok {
			return false
		}
	}

	if sliceContainsNil(slc) {
		return false
	}

	return true
}

func Discard(slice []interface{}, n interface{}) []interface{} {
	result := []interface{}{}
	for _, s := range slice {
		if s != n {
			result = append(result, s)
		}
	}
	return result
}

// We check if the element is non-string, non-float64, non-boolean. Then it is another node or <nil>. There is no other option.
func IsNonStringFloatBool(v interface{}) bool {
	var isString, isFloat, isBool bool
	if _, ok := v.(string); ok {
		isString = true
	} else if _, ok := v.(float64); ok {
		isFloat = true
	} else if _, ok := v.(bool); ok {
		isBool = true
	}
	if !isString && !isFloat && !isBool {
		return true
	}
	return false
}

func mapContainsNil(mp map[string]interface{}) bool {
	for _, m := range mp {
		if m == nil {
			return true
		}
	}
	return false
}

func sliceContainsNil(slice []interface{}) bool {
	for _, s := range slice {
		if s == nil {
			return true
		}
	}
	return false
}
