package domain

import "github.com/golang-jwt/jwt/v5"

type JwtModel struct {
	Iss string           `json:"iss"`
	Aud jwt.ClaimStrings `json:"aud"`
	Sub string           `json:"sub"`
	Azp string           `json:"azp"`
	Jti string           `json:"jti"`
	Exp *jwt.NumericDate `json:"exp"`
	Nbf *jwt.NumericDate `json:"nbf"`
	Iat *jwt.NumericDate `json:"iat"`
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
	NameTag    string `json:"nametag"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Email      string `json:"email"`
	Picture    string `json:"picture"`
}
