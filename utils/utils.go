package utils

func StripNonAlphaNumeric(s string) string {
	var result []rune

	for _, r := range s {
		if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' {
			result = append(result, r)
		}
	}

	return string(result)
}
