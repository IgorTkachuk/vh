package jwt

import (
	"encoding/json"
	"fmt"
	"github.com/cristalhq/jwt/v3"
	"github.com/google/uuid"
	"os"
	"time"
	"vh/internal/models"
	"vh/package/cache"
)

var _ Helper = &helper{}

type RT struct {
	RefreshToken string `json:"refresh_token,omitempty"`
}

type UserClaims struct {
	jwt.RegisteredClaims
	login string `json:"login"`
}

type Helper interface {
	GenerateAccessToken(u models.UserDto) ([]byte, error)
	UpdateRefreshToken(rt RT) ([]byte, error)
}

type helper struct {
	RTCache cache.Repository
}

func NewHelper(cache cache.Repository) *helper {
	return &helper{RTCache: cache}
}

func (h helper) GenerateAccessToken(u models.UserDto) ([]byte, error) {
	fmt.Printf("Create token for user %s\n", u.Login)
	secret := os.Getenv("JWT_SECRET")
	signer, err := jwt.NewSignerHS(jwt.HS256, []byte(secret))
	if err != nil {
		return nil, err
	}

	builder := jwt.NewBuilder(signer)
	claim := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        u.Login,
			Audience:  []string{"users"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)),
		},
	}

	token, err := builder.Build(claim)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Create refresh token for user %s\n", u.Login)
	refreshTokenUUID := uuid.New()
	userBytes, _ := json.Marshal(u)

	err = h.RTCache.Set([]byte(refreshTokenUUID.String()), userBytes, 0)
	if err != nil {
		return nil, err
	}

	tokensBytes, err := json.Marshal(map[string]string{
		"token":         token.String(),
		"refresh_token": refreshTokenUUID.String(),
	})

	if err != nil {
		return nil, err
	}

	return tokensBytes, nil
}

func (h helper) UpdateRefreshToken(rt RT) ([]byte, error) {
	defer h.RTCache.Del([]byte(rt.RefreshToken))

	userBytes, err := h.RTCache.Get([]byte(rt.RefreshToken))
	if err != nil {
		return nil, err
	}

	var user models.UserDto
	err = json.Unmarshal(userBytes, &user)
	if err != nil {
		return nil, err
	}

	return h.GenerateAccessToken(user)
}
