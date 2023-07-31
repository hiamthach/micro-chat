package util

func Contains(element string, slice []string) bool {
	for _, value := range slice {
		if element == value {
			return true
		}
	}
	return false
}
