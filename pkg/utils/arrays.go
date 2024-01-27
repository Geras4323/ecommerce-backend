package utils

func CheckIfInArray[S ~[]E, E comparable](s S, e E) bool {
	for i := range s {
		if e == s[i] {
			return true
		}
	}
	return false
}