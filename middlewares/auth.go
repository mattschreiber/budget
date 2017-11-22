package middlewares


//  Generate RSA signing files via shell (adjust as needed):
//
//  $ openssl genrsa -out app.rsa 1024
//  $ openssl rsa -in app.rsa -pubout > app.rsa.pub
//
// Code borrowed and modified from the following sources:
// https://www.youtube.com/watch?v=dgJFeqeXVKw
// https://goo.gl/ofVjK4
// https://github.com/dgrijalva/jwt-go
//

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
  "context"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
)

// const (
// 	// For simplicity these files are in the same folder as the app binary.
// 	// You shouldn't do this in production.
// 	privKeyPath = "app.rsa"
// 	pubKeyPath  = "app.rsa.pub"
// )

// var (
// 	verifyKey *rsa.PublicKey
// 	signKey   *rsa.PrivateKey
// )

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type TokenContextKey struct {
  Name string
}

type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// type User struct {
// 	ID       int    `json:"id"`
// 	Name     string `json:"name"`
// 	Username string `json:"username"`
// 	Password string `json:"password"`
// }


type Token struct {
	Token string `json:"token"`
}

// JwtCustomClaims are custom claims extending default ones.
type JwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.StandardClaims
}

type Exception struct {
    Message string `json:"message"`
}


func LoginHandler(w http.ResponseWriter, r *http.Request) {

	var user UserCredentials

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, "Error in request")
		return
	}
	if strings.ToLower(user.Username) != "someone" {
		if user.Password != "p@ssword" {
			w.WriteHeader(http.StatusForbidden)
			fmt.Println("Error logging in")
			fmt.Fprint(w, "Invalid credentials")
			return
		}
    return
	}


  // Set custom claims
	claims := &JwtCustomClaims{
		user.Username,
		true,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

  token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
  t, err := token.SignedString([]byte("secret"))
	if err != nil {
		fatal(err)
	}
  response := map[string]string{"token": t}
  json, err := json.Marshal(response)
  if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)

}

func ValidateToken(next http.HandlerFunc) http.HandlerFunc {
  signingKey := []byte("secret")
  return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
    token, err := request.ParseFromRequestWithClaims(req, request.AuthorizationHeaderExtractor, &JwtCustomClaims{},
      func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
        return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
    }
       return signingKey, nil
     })
     if err != nil {
        http.Error(w, err.Error(), http.StatusUnauthorized)
        return
     }

     if token.Valid {
       t := TokenContextKey{}
       if claims, ok := token.Claims.(*JwtCustomClaims); ok {
         ctx := context.WithValue(req.Context(), t.Name, claims.Name)
         fmt.Printf("%v %v\n", claims.Name, claims.StandardClaims.ExpiresAt)
         next(w, req.WithContext(ctx))
       } else {
     		// w.WriteHeader(http.StatusUnauthorized)
        http.Error(w, err.Error(), http.StatusUnauthorized)
        return
       }
     }else {
       fmt.Fprint(w, "Token is not valid")
       // w.WriteHeader(http.StatusUnauthorized)
       http.Error(w, err.Error(), http.StatusUnauthorized)
       return
     }
  })
}
