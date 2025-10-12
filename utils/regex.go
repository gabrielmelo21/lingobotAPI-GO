package utils

import "regexp"

func ValidateName(name string) bool {
	re := regexp.MustCompile(`^[A-Za-zÀ-ÖØ-öø-ÿ\s]+$`)
	return re.MatchString(name)
}
