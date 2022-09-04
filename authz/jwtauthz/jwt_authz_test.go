package jwtauthz

import (
	"encoding/json"
	"testing"

	"github.com/astaxie/beego"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
)

func createJwtToken(claims jwt.MapClaims) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	tokenString, err := token.SignedString([]byte(beego.AppConfig.String("secret")))
	return tokenString, err
}

func TestJwtAuthorization(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	Convey("ExtractClaimsFromToken()", t, func() {
		Convey("When bearer token's format is invalid", func() {
			testBearerTokenStr := "bearer abcde" // should be "Bearer abcde"
			jwtAuthorization := NewJwtAuthorization(beego.AppConfig.String("secret"), testBearerTokenStr)
			So(jwtAuthorization.jwtTokenStr, ShouldEqual, "")
		})

		Convey("When JWT cannot be decoded", func() {
			testBearerTokenStr := "Bearer abcde"
			jwtAuthorization := NewJwtAuthorization(beego.AppConfig.String("secret"), testBearerTokenStr)
			err := jwtAuthorization.ExtractClaimsFromToken()
			So(err, ShouldNotBeNil)
		})

		Convey("When the JWT's claim has the wrong key", func() {
			// JWT's data claims should be in the format of (id int, name string, email string, roles []string)
			claims := make(jwt.MapClaims)
			claims["idd"] = 2
			claims["namez"] = "buddy"
			claims["emaily"] = "budiryan@hotmail.com"
			claims["rolez"] = []string{"superuser"}

			testTokenStr, _ := createJwtToken(claims)
			jwtAuthorization := NewJwtAuthorization(beego.AppConfig.String("secret"), "Bearer "+testTokenStr)
			err := jwtAuthorization.ExtractClaimsFromToken()
			So(err, ShouldNotBeNil)
		})

		Convey("When the JWT's claim has the wrong data format", func() {
			// JWT's data claims should be in the format of (id int, name string, email string, roles []string)
			claims := make(jwt.MapClaims)
			claims["id"] = "2"
			claims["name"] = "333"
			claims["email"] = "budiryan@hotmail.com"
			claims["roles"] = []string{"444"}

			testTokenStr, _ := createJwtToken(claims)
			jwtAuthorization := NewJwtAuthorization(beego.AppConfig.String("secret"), "Bearer "+testTokenStr)
			err := jwtAuthorization.ExtractClaimsFromToken()
			So(err, ShouldNotBeNil)
		})

		Convey("When extracting claims from the token is successful is successful", func() {
			// JWT's data claims should be in the format of (id int, name string, email string, roles []string)
			claims := make(jwt.MapClaims)
			claims["id"] = 2
			claims["name"] = "budi ryan"
			claims["email"] = "budiryan@hotmail.com"
			claims["roles"] = []string{"superuser", "admin"}

			testTokenStr, _ := createJwtToken(claims)
			jwtAuthorization := NewJwtAuthorization(beego.AppConfig.String("secret"), "Bearer "+testTokenStr)
			err := jwtAuthorization.ExtractClaimsFromToken()
			So(err, ShouldBeNil)
		})
	})

	Convey("IsAuthorized()", t, func() {
		Convey("When user is not authorized", func() {
			// JWT's data claims should be in the format of (id int, name string, email string, roles []string)
			claims := make(jwt.MapClaims)
			claims["id"] = 2
			claims["name"] = "budi ryan"
			claims["email"] = "budiryan@hotmail.com"
			claims["roles"] = []string{"admin"}

			testTokenStr, _ := createJwtToken(claims)
			jwtAuthorization := NewJwtAuthorization(beego.AppConfig.String("secret"), "Bearer "+testTokenStr)
			_ = jwtAuthorization.ExtractClaimsFromToken()

			authorized := jwtAuthorization.IsAuthorized("superuser")
			So(authorized, ShouldBeFalse)
		})

		Convey("When user is authorized to access the controller's resource", func() {
			// JWT's data claims should be in the format of (id int, name string, email string, roles []string)
			claims := make(jwt.MapClaims)
			claims["id"] = 2
			claims["name"] = "budi ryan"
			claims["email"] = "budiryan@hotmail.com"
			claims["roles"] = []string{"superuser"} // "user's privilege is superuser"

			testTokenStr, _ := createJwtToken(claims)
			jwtAuthorization := NewJwtAuthorization(beego.AppConfig.String("secret"), "Bearer "+testTokenStr)
			_ = jwtAuthorization.ExtractClaimsFromToken()

			authorized := jwtAuthorization.IsAuthorized("superuser") // only "superuser" is granted resource
			So(authorized, ShouldBeTrue)
		})

		Convey("it should contains authz information", func() {
			// JWT's data claims should be in the format of (id int, name string, email string, roles []string)
			authzObj := &Authz{
				Roles: []Role{
					Role{
						ID:   1,
						Name: "tes",
					},
				},
				Permissions: []string{"tes"},
			}
			authzSer, _ := json.Marshal(&authzObj)
			claims := make(jwt.MapClaims)
			claims["id"] = 2
			claims["name"] = "budi ryan"
			claims["email"] = "budiryan@hotmail.com"
			claims["roles"] = []string{"superuser"} // "user's privilege is superuser"
			claims["authz"] = string(authzSer)

			testTokenStr, _ := createJwtToken(claims)
			jwtAuthorization := NewJwtAuthorization(beego.AppConfig.String("secret"), "Bearer "+testTokenStr)
			_ = jwtAuthorization.ExtractClaimsFromToken()

			So(jwtAuthorization.JwtClaims.Authz, ShouldNotBeNil)
			So(jwtAuthorization.JwtClaims.Authz.Permissions[0], ShouldEqual, "tes")
		})
	})
}
