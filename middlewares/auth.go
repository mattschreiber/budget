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

type Response struct {
	Data string `json:"data"`
}

type Token struct {
	Token string `json:"token"`
}

// jwtCustomClaims are custom claims extending default ones.
type jwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.StandardClaims
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
	claims := &jwtCustomClaims{
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



	// token := jwt.New(jwt.SigningMethodHS256)
	// claims := make(jwt.MapClaims)
	// claims["exp"] = time.Now().Add(time.Hour * time.Duration(1)).Unix()
	// claims["iat"] = time.Now().Unix()
	// token.Claims = claims

	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	fmt.Fprintln(w, "Error extracting the key")
	// 	fatal(err)
	// }
  // var signingKey = []byte("signing-key")
	// tokenString, err := token.SignedString(signingKey)
  //
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	fmt.Fprintln(w, "Error while signing the token")
	// 	fatal(err)
	// }
  //
	// response := Token{tokenString}
	// JsonResponse(response, w)

}

func ValidateToken(w http.ResponseWriter, r *http.Request) {
  signingKey := []byte("secret")

  token, err := request.ParseFromRequestWithClaims(r, request.AuthorizationHeaderExtractor, &jwtCustomClaims{},
    func(token *jwt.Token) (interface{}, error) {
     return signingKey, nil
   })
   if claims, ok := token.Claims.(*jwtCustomClaims); ok && token.Valid {
       fmt.Printf("%v %v", claims.Name, claims.StandardClaims.ExpiresAt)
   }
	if err == nil {
		if token.Valid {
      x := map[string]string{"data": "token valid"}
      w.Header().Set("Content-Type", "application/json")
      json.NewEncoder(w).Encode(x)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Token is not valid")
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Unauthorized access to this resource")
  }
}


// func JsonResponse(response interface{}, w http.ResponseWriter) {
//
// 	json, err := json.Marshal(response)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
//
// 	w.WriteHeader(http.StatusOK)
// 	w.Header().Set("Content-Type", "application/json")
// 	w.Write(json)
// }
