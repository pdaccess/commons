package domain

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type PdaccessClaims struct {
	UserId string `json:"sub"`
	Urk    string `json:"x-urk"`
	Tpk    string `json:"x-tpk,omitempty"`
	Role   string `json:"x-role,omitempty"`
	Realm  string `json:"x-realm,omitempty"`
	AuthId string `json:"x-auth-id,omitempty"`

	Aud string `json:"aud,omitempty"`
	Nbf int64  `json:"nbf,omitempty"`
	Exp int64  `json:"exp"`
	Iat int64  `json:"iat"`
}

func (a *PdaccessClaims) GetAudience() (jwt.ClaimStrings, error) {
	return jwt.ClaimStrings([]string{a.Aud}), nil
}

func (a *PdaccessClaims) GetExpirationTime() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(time.Unix(a.Exp, 0)), nil
}

func (a *PdaccessClaims) GetIssuedAt() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(time.Unix(a.Iat, 0)), nil
}

func (a *PdaccessClaims) GetIssuer() (string, error) {
	return "pvault", nil
}

func (a *PdaccessClaims) GetNotBefore() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(time.Unix(a.Nbf, 0)), nil
}

func (a *PdaccessClaims) GetSubject() (string, error) {
	return a.UserId, nil
}

func (a *PdaccessClaims) Parse(content string) error {

	_, err := jwt.ParseWithClaims(content, a, nil)

	if err != nil && !errors.Is(err, jwt.ErrTokenUnverifiable) {
		return fmt.Errorf("token jwt decode :%w", err)
	}

	return nil
}
