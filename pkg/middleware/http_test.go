package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pdaccess/commons/pkg/domain"
	"github.com/pdaccess/commons/pkg/middleware"
)

type mockHandler struct{}

func (*mockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func TestAuthMiddleware(t *testing.T) {
	hand := middleware.AuthzAdapter(&mockHandler{})

	server := httptest.NewServer(hand)
	defer server.Close()

	resp, err := http.Get(server.URL)

	if err != nil {
		t.Fatalf("error in http: %v", err)
	}

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("wrong response status code: %d", resp.StatusCode)
	}

	client := &http.Client{}

	req, _ := http.NewRequest(http.MethodGet, server.URL, nil)
	req.Header.Add("Authorization", "Bearer "+"wrong token")

	resp, err = client.Do(req)

	if err != nil {
		t.Fatalf("error in http: %v", err)
	}

	if resp.StatusCode == http.StatusOK {
		t.Fatalf("wrong response status code: %d", resp.StatusCode)
	}

	req, _ = http.NewRequest(http.MethodGet, server.URL, nil)
	req.Header.Add("Authorization", "Bearer "+getToken())

	resp, err = client.Do(req)

	if err != nil {
		t.Fatalf("error in http: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("wrong response status code: %d", resp.StatusCode)
	}
}

func getToken() string {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, &domain.PdaccessClaims{
		UserId: "testUser",
		AuthId: "456",
	})

	tokenString, _ := claims.SignedString([]byte("secret"))

	return tokenString
}
