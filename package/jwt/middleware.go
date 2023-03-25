package jwt

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cristalhq/jwt/v3"
	"net/http"
	"os"
	"strings"
	"time"
)

func Middleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(authHeader) != 2 {
			unauthorized(w, fmt.Errorf("mailformed token"))
			return
		}

		jwtToken := authHeader[1]
		secret := os.Getenv("JWT_SECRET")
		verifier, err := jwt.NewVerifierHS(jwt.HS256, []byte(secret))
		if err != nil {
			unauthorized(w, err)
			return
		}

		token, err := jwt.ParseAndVerifyString(jwtToken, verifier)
		if err != nil {
			unauthorized(w, err)
			return
		}

		var uc UserClaims
		err = json.Unmarshal(token.RawClaims(), &uc)
		if err != nil {
			unauthorized(w, err)
			return
		}

		if valid := uc.IsValidAt(time.Now()); !valid {
			unauthorized(w, fmt.Errorf("access token is expired"))
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", uc.ID)
		h(w, r.WithContext(ctx))
	}
}

func unauthorized(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusForbidden)
	w.Write([]byte(err.Error()))
}
