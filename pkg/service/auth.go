package service

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	apiserver "mongo_db"
	"mongo_db/pkg/repository"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

const (
	tokenTTL = time.Hour * 1
)

type TokenClaims struct {
	jwt.StandardClaims
	UserId string
}

type AuthService struct {
	repository repository.Authorization
}

func NewAuthSerivce(repo repository.Authorization) *AuthService {
	return &AuthService{repository: repo}
}

func (a *AuthService) CreateUser(user apiserver.User) error {
	user.RefreshToken, _ = a.HashRefreshToken(user.RefreshToken)
	return a.repository.CreateUser(user)
}

func (a *AuthService) GetUserById(id string) (*apiserver.User, error) {
	return a.repository.GetUserById(id)
}

func (a *AuthService) GenerateRefreshToken() (string, error) {
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}
	refreshToken := base64.StdEncoding.EncodeToString(tokenBytes)

	return refreshToken, nil
}
func (a *AuthService) GenerateAccessToken(id string, refreshToken string) (string, error) {
	claims := &TokenClaims{
		UserId: id + refreshToken,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	accessToken, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (a *AuthService) HashRefreshToken(refreshToken string) (string, error) {
	hashToken, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashToken), nil
}

func (a *AuthService) UpdateRefreshToken(id string, refreshToken string) error {
	hashedToken, err := a.HashRefreshToken(refreshToken)
	if err != nil {
		return err
	}
	return a.repository.UpdateRefreshToken(id, hashedToken)
}

func (a *AuthService) ParseAccessToken(accessToken string) (string, error) {
	token, err := jwt.Parse(accessToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", fmt.Errorf("invalid method")
		}
		return []byte(os.Getenv("JWT_KEY")), nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userId := claims["UserId"].(string)
		return userId, nil
	}
	return "", fmt.Errorf("invalid token")
}
