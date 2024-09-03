package middleware

import "strings"

func (c *Credential) ExtractServiceMethod(fullMethod string) (string, string) {
	parts := strings.Split(fullMethod, "/")

	if len(parts) != 3 {
		return "", ""
	}

	return parts[1], parts[2]
}
