package middleware

import (
	"fmt"
	"net/http"
	"strings"
)

func (c *Credential) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		token := strings.Replace(r.Header.Get("Authorization"), "Bearer ", "", -1)
		err := c.VerifyToken(token)

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "Access denied: invalid token")
			return
		}

		next(w, r)
	}
}
