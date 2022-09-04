// package jwtauthz provides a convenient tool for extracting payload data from a JWT token
package jwtauthz

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/dgrijalva/jwt-go"
)

// NewJwtAuthorization inits a new JwtAuthorization obj by supplying secret and bearerTokenStr
func NewJwtAuthorization(secret string, bearerTokenStr string) JwtAuthorization {
	jwtTokenStr := parseBearerToken(bearerTokenStr)
	return JwtAuthorization{jwtTokenStr, []byte(secret), JwtClaims{}}
}

// ExtractClaimsFromToken extracts the payload data from a JWT token and puts it into JwtClaims object
func (j *JwtAuthorization) ExtractClaimsFromToken() error {
	// Decode the token and throws error if there is any problem while decoding
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(j.jwtTokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		// check token signing method etc
		return j.secret, nil
	})
	if err != nil || !token.Valid {
		return err
	}

	// Store the sent data inside JwtClaims struct
	err = j.parseJwtClaims(claims)
	if err != nil {
		return err
	}

	return nil
}

// IsAuthorized checks whether a user is authorized for a certain role in authorizedRole
func (j *JwtAuthorization) IsAuthorized(authorizedRole string) bool {
	// Checks whether the user is authorized to access a resource
	for _, role := range j.JwtClaims.Roles {
		roleString, ok := role.(string)
		if !ok {
			return false
		}
		if roleString == authorizedRole {
			return true
		}
	}
	return false
}

func parseBearerToken(bearerTokenStr string) string {
	splitToken := strings.Split(bearerTokenStr, "Bearer ")
	if len(splitToken) < 2 {
		return ""
	}
	return splitToken[1]
}

func (j *JwtAuthorization) parseJwtClaims(claims jwt.MapClaims) error {
	// Parse ID
	id, ok := claims["id"]
	if !ok {
		return fmt.Errorf("key 'id' is not contained within JWT's claims")
	}
	idFloat64, ok := id.(float64)
	if !ok {
		return fmt.Errorf("key 'id' is in wrong format")
	}
	j.JwtClaims.ID = int(idFloat64)

	// Parse email
	email, ok := claims["email"]
	if !ok {
		return fmt.Errorf("key 'email' is not contained within JWT's claims")
	}
	j.JwtClaims.Email, ok = email.(string)
	if !ok {
		return fmt.Errorf("key 'email' is in wrong format")
	}

	// Parse name
	name, ok := claims["name"]
	if !ok {
		return fmt.Errorf("key 'name' is not contained within JWT's claims")
	}
	j.JwtClaims.Name, ok = name.(string)
	if !ok {
		return fmt.Errorf("key 'name' is in wrong format")
	}

	// Parse roles
	roles, ok := claims["roles"]
	if !ok {
		return fmt.Errorf("key 'roles' is not contained within JWT's claims")
	}
	j.JwtClaims.Roles, ok = roles.([]interface{})
	if !ok {
		return fmt.Errorf("key 'roles' is in wrong format")
	}

	// parse authz section
	authz, ok := claims["authz"]
	if !ok {
		return nil // for compability reason with legacy stuff we just ignore if authz key dont exists for now
	}

	if err := json.Unmarshal([]byte(authz.(string)), &j.JwtClaims.Authz); err != nil {
		return fmt.Errorf("key 'authz' is in wrong format: %v", err)
	}

	return nil
}

// AuthenticateJWTMiddleware is a convenience middleware function for JWT authentication within your beego application
//
// One Example: put this within the router section which you want to validate using JWT
//
//	 ns := beego.NewNamespace("/public",
//		 beego.NSNamespace("/v1",
//			 beego.NSNamespace("/admin",
//				 AuthenticateJWTMiddleware(conf.JWTSecret), // put this here
//				 beego.NSNamespace("/notification",
//					 beego.NSNamespace("/csv",
//						 beego.NSInclude(&v1.AdminCSVController{Publisher: publisher, AppConf: conf, RedisMaster: redisMaster, RedisSlave: redisSlave}),
//					 ),
//					 beego.NSInclude(&v1.Template{RedisMaster: redisMaster, RedisSlave: redisSlave, Publisher: publisher}),
//				 ),
//			 ),
//			 beego.NSInclude(&v1.SanityCheck{Environtment: conf.Environment, ClientRPC: rpcClient}),
//		 ),
//	 )
//
func AuthenticateJWTMiddleware(jwtSecret string) beego.LinkNamespace {
	return beego.NSBefore(func(ctx *context.Context) {
		bearerTokenStr := ctx.Request.Header.Get("Authorization")
		jwtAuthorization := NewJwtAuthorization(jwtSecret, bearerTokenStr)
		err := jwtAuthorization.ExtractClaimsFromToken()
		if err != nil {
			log.Printf("[router] failed to get jwt token. %+v\n", err)
			errMessage, _ := json.Marshal(map[string]interface{}{"message": "JWT is invalid"})
			ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
			ctx.ResponseWriter.Write(errMessage)
			return
		}
		if !jwtAuthorization.IsAuthorized("superuser") {
			errMessage, _ := json.Marshal(map[string]interface{}{"message": "you are not authorized to access this resource"})
			ctx.ResponseWriter.WriteHeader(http.StatusForbidden)
			ctx.ResponseWriter.Write(errMessage)
		}
		return
	})
}
