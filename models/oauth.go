package models

import "github.com/golang-jwt/jwt/v5"

type jwtModel struct {
	Iss string
	Aud jwt.ClaimStrings
	Sub string
	Azp string
	Jti string
	Exp jwt.NumericDate
	Nbf jwt.NumericDate
	Iat jwt.NumericDate
}

type OAuthGoogleJWTModel struct {
	jwtModel
	Email          string
	Email_verified bool
	Name           string
	Picture        string
	Given_name     string
	Family_name    string
}

func (j jwtModel) GetExpirationTime() (*jwt.NumericDate, error) {
	return &j.Exp, nil
}

func (j jwtModel) GetIssuedAt() (*jwt.NumericDate, error) {
	return &j.Iat, nil
}

func (j jwtModel) GetNotBefore() (*jwt.NumericDate, error) {
	return &j.Nbf, nil
}

func (j jwtModel) GetIssuer() (string, error) {
	return j.Iss, nil
}

func (j jwtModel) GetSubject() (string, error) {
	return j.Sub, nil
}

func (j jwtModel) GetAudience() (jwt.ClaimStrings, error) {
	return j.Aud, nil
}
