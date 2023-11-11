package middlewares

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jws"
)

type ParsedUserClaim struct {
    UserId      string  `json:"user_id"`
    Email       string  `json:"email"`
}

func DecodeJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.Request.Header.Get("X-IDTOKEN")
		if auth == "" {
			c.String(http.StatusForbidden, "No Authorization header provided")
			c.Abort()
			return
		}
		//
		// tokenString := strings.TrimPrefix(auth, "Bearer ")
		// if tokenString == auth {
		// 	c.String(http.StatusForbidden, "Could not find bearer token in Authorization header")
		// 	c.Abort()
		// 	return
		// }
		var keysJWK = os.Getenv("JWT_SECRET")
		setOfKeys, err := jwk.ParseString(keysJWK)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to create JWK: %s", err))
			c.Abort()
			return
		}

        // change this for id_token or access_token
		privKey, success := setOfKeys.Get(0)
		if !success {
			c.String(http.StatusInternalServerError, "Could not find key at given index")
			c.Abort()
			return
		}

		pubkey, err := jwk.PublicKeyOf(privKey)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to get public key: %s", err.Error()))
			c.Abort()
			return
		}
  
		verifiedToken, err := jws.Verify([]byte(auth), jwa.RS256, pubkey)
		if err != nil {
			c.String(http.StatusForbidden, fmt.Sprintf("Failed to verify token from HTTP request: %s", err.Error()))
			c.Abort()
			return
		}

		if err != nil {
			log.Fatal(err)
		}

		var parsedUser ParsedUserClaim
        err = json.Unmarshal([]byte(verifiedToken), &parsedUser)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to unmarshal data: %s", err.Error()))
			c.Abort()
			return
		}

        parsedUserBytes, err := json.Marshal(parsedUser)
        if err != nil {
            panic(err)
        }

        var userDetails map[string]interface{}
        err = json.Unmarshal(parsedUserBytes, &userDetails)
        if err != nil {
            panic(err)
        }

        c.Set("userDetails", userDetails)
        c.Next()
	}
}
