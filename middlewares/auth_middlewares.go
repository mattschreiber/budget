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
  "budget/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
  "context"
  "os"
	// "strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
  "golang.org/x/crypto/bcrypt"
)


func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type TokenContextKey struct {
  Name string
}

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

// Checks user credentials and if valid issues a new jwt token
func LoginHandler(w http.ResponseWriter, r *http.Request) {

	var user models.UserCredentials
  // get user credentials from the body of http request
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, "Error in request")
		return
	}

  // query db for username and password information
  var passwordHash string
  passwordHash, err = models.VerifyLogin(user)
  if err != nil {
    w.WriteHeader(http.StatusForbidden)
    fmt.Fprint(w, "Invalid Username or Password")
    return
  }
  // verify password is correct
  err = checkPasswordHash(passwordHash, user.Password)
  if err != nil {
    w.WriteHeader(http.StatusForbidden)
    fmt.Fprint(w, "Invalid credentials")
    return
  }

  response, err := generateToken(user)
  if err != nil {
    // http.Error(w, err.Error(), http.StatusInternalServerError)
    w.WriteHeader(http.StatusForbidden)
    fmt.Fprint(w, "Unable to generate Token")
    return
  }
  json, err := json.Marshal(response)
  if err != nil {
    // http.Error(w, err.Error(), http.StatusInternalServerError)
    w.WriteHeader(http.StatusForbidden)
    fmt.Fprint(w, "Unable to marshal token to json")
    return
  }
  w.Header().Set("Content-Type", "application/json")
  w.Write(json)
}

func checkPasswordHash(passwordHash, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
  if err != nil {
    return err
  }
	return nil
}

func generateToken(user models.UserCredentials) (map[string]string, error) {
  // Set custom claims
  claims := &JwtCustomClaims{
    user.Username,
    true,
    jwt.StandardClaims{
      ExpiresAt: time.Now().Add(time.Hour * 720).Unix(),
    },
  }

  token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
  t, err := token.SignedString([]byte(os.Getenv("AWS_SIGNING_KEY")))
  if err != nil {
    return map[string]string{}, err
  }
  response := map[string]string{"token": t}
  return response, err
}

func ValidateToken(next http.HandlerFunc) http.HandlerFunc {
  signingKey := []byte(os.Getenv("AWS_SIGNING_KEY"))
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
