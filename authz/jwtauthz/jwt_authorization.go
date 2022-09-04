package jwtauthz

// JwtAuthorization contains token, secret, and the underlying payload data
type JwtAuthorization struct {
	jwtTokenStr string
	secret      []byte
	JwtClaims   JwtClaims
}
