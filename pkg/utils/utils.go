package utils

func StringMapToSlice(old map[string]bool) []string {
	result := make([]string, 0)
	for s := range old {
		result = append(result, s)
	}
	return result
}
