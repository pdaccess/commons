package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"runtime"

	"strings"

	common_domain "git.h2hsecure.com/pda/commons/pkg/domain"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

func LoggingAdapter(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientRequest := ClientFromCtx(r.Context())

		event := log.Info().
			Str("RequestId", clientRequest.RequestId).
			Str("url", r.RequestURI).
			Str("method", r.Method)
		for name, values := range r.Header {
			for _, value := range values {
				event.Str(name, value)
			}
		}
		event.Send()

		h.ServeHTTP(w, r)
	})
}

func RequestIdAdapter(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := ClientRequestType(ClientRequestStr)
		clientRequest := r.Context().Value(t)
		if clientRequest == nil {
			clientRequest = ClientFromCtx(r.Context())
		}

		h.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), t, clientRequest)))
	})
}

func AuthzAdapter(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		client := ClientFromCtx(r.Context())

		authStr := r.Header.Get("Authorization")
		tokenStr, found := strings.CutPrefix(authStr, "Bearer ")

		decodedToken, _, err := jwt.NewParser().ParseUnverified(tokenStr, &common_domain.PdaccessClaims{})
		if !found || err != nil {
			response := common_domain.ApiResponse{
				Code:      common_domain.ErrNotAuth.Error(),
				Message:   "Not Authorized",
				RequestId: client.RequestId,
			}
			buf, _ := json.Marshal(response)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(buf)
			return
		}

		client.User = *decodedToken.Claims.(*common_domain.PdaccessClaims)

		h.ServeHTTP(w, r)
	})
}

func PanicAdapter(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				buf := make([]byte, 10<<10)
				n := runtime.Stack(buf, false)
				log.Error().
					Str("context", "panic").
					Msgf("%v\n\n%s", err, buf[:n])

				client := ClientFromCtx(r.Context())

				resp := common_domain.ApiResponse{
					RequestId: client.RequestId,
					Message:   string(buf[:n]),
					Code:      common_domain.ErrInternal.Error(),
				}

				json, _ := json.Marshal(resp)
				http.Error(w, string(json), http.StatusInternalServerError)
			}
		}()

		h.ServeHTTP(w, r)
	})
}
