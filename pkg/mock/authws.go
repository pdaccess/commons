package mock

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"git.h2hsecure.com/pda/commons/pkg/domain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

type tokenJSON struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    uint64 `json:"expires_in"` // at least PayPal returns string, while most return number
	// error fields
	// https://datatracker.ietf.org/doc/html/rfc6749#section-5.2
	ErrorCode        string `json:"error"`
	ErrorDescription string `json:"error_description"`
	ErrorURI         string `json:"error_uri"`
}

func AuthwsEndpointCreate() *httptest.Server {

	return httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// loginRequest := make(map[string]string)
		// loginData, _ := io.ReadAll(r.Body)
		// decodedText, _ := base64.StdEncoding.DecodeString(string(loginData))
		// json.Unmarshal(loginData, &loginRequest)
		user, password, _ := r.BasicAuth()
		err := r.ParseForm()

		if err != nil {
			w.Write([]byte(err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
		}

		log.Info().
			Str("path", r.RequestURI).
			Str("user", user).
			Str("password", password).
			Interface("body", r.Form).
			Msg("auth request came")

		var loginResponse tokenJSON
		if r.Form.Get("username") == ValidUser && r.Form.Get("password") == ValidPass {
			claims := jwt.NewWithClaims(jwt.SigningMethodHS256, &domain.PdaccessClaims{
				UserId: MockUserId,
				AuthId: "456",
			})

			tokenString, err := claims.SignedString([]byte("secret"))
			if err != nil {
				w.Write([]byte(err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
			}

			loginResponse = tokenJSON{
				AccessToken:  tokenString,
				TokenType:    "Berear",
				RefreshToken: "ok",
				ExpiresIn:    uint64(time.Now().Add(1 * time.Hour).Unix()),
			}
		} else {
			loginResponse = tokenJSON{
				ErrorDescription: fmt.Sprintf("auth failed: %s", r.Form.Get("username")),
			}
		}

		log.Info().Interface("res", loginResponse).Send()

		lrJson, _ := json.Marshal(loginResponse)

		w.Header().Add("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(lrJson)

	}))
}
