package domain

import "github.com/golang-jwt/jwt/v5"

type JwtModel struct {
	Iss string
	Aud jwt.ClaimStrings
	Sub string
	Azp string
	Jti string
	Exp *jwt.NumericDate
	Nbf *jwt.NumericDate
	Iat *jwt.NumericDate
}

type GoogleJWTModel struct {
	JwtModel
	Email          string
	Email_verified bool
	Name           string
	Picture        string
	Given_name     string
	Family_name    string
}

func (j JwtModel) GetExpirationTime() (*jwt.NumericDate, error) {
	return j.Exp, nil
}

func (j JwtModel) GetIssuedAt() (*jwt.NumericDate, error) {
	return j.Iat, nil
}

func (j JwtModel) GetNotBefore() (*jwt.NumericDate, error) {
	return j.Nbf, nil
}

func (j JwtModel) GetIssuer() (string, error) {
	return j.Iss, nil
}

func (j JwtModel) GetSubject() (string, error) {
	return j.Sub, nil
}

func (j JwtModel) GetAudience() (jwt.ClaimStrings, error) {
	return j.Aud, nil
}

type SessionJwtModel struct {
	JwtModel
	UserID     int64  `json:"userid"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Email      string `json:"email"`
	Picture    string `json:"picture"`
}
