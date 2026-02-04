package domain

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type PdaccessClaims struct {
	UserId string `json:"sub"`
	Role   string `json:"uRole"`
	Realm  string `json:"realm"`
	AuthId string `json:"authId"`
	Aud    string `json:"aud"`
	Nbf    int64  `json:"nbf"`
	Exp    int64  `json:"exp"`
	Iat    int64  `json:"iat"`
}

// GetAudience implements jwt.Claims.
func (a *PdaccessClaims) GetAudience() (jwt.ClaimStrings, error) {
	return jwt.ClaimStrings([]string{a.Aud}), nil
}

// GetExpirationTime implements jwt.Claims.
func (a *PdaccessClaims) GetExpirationTime() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(time.UnixMicro(a.Exp)), nil
}

// GetIssuedAt implements jwt.Claims.
func (a *PdaccessClaims) GetIssuedAt() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(time.UnixMicro(a.Iat)), nil
}

// GetIssuer implements jwt.Claims.
func (a *PdaccessClaims) GetIssuer() (string, error) {
	return "pdaccess", nil
}

// GetNotBefore implements jwt.Claims.
func (a *PdaccessClaims) GetNotBefore() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(time.UnixMicro(a.Nbf)), nil
}

// GetSubject implements jwt.Claims.
func (a *PdaccessClaims) GetSubject() (string, error) {
	return a.UserId, nil
}

func (a *PdaccessClaims) Valid() error {
	if time.Unix(a.Exp, 0).Before(time.Now()) {
		return errors.New("invalid token")
	}
	return nil
}

func (a *PdaccessClaims) Parse(content string) error {

	_, err := jwt.ParseWithClaims(content, a, nil)

	if err != nil && !errors.Is(err, jwt.ErrTokenUnverifiable) {
		return fmt.Errorf("token jwt decode :%w", err)
	}

	return nil
}
