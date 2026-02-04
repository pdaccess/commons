package mock

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"

	"github.com/rs/zerolog/log"
)

type SeviceRequst struct {
	PerPage     int    `json:"perPage"`
	CurrentPage int    `json:"currentPage"`
	Filter      string `json:"filter"`
	Category    string `json:"category"`
}

type ServiceResponse struct {
	ServiceId     string   `json:"inventoryId"`
	IPAddress     string   `json:"ipAddress"`
	Port          int      `json:"port"`
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	LeaseDuration int64    `json:"lease_duration"`
	LeaseId       string   `json:"lease_id"`
	RequestId     string   `json:"request_id"`
	Renewable     bool     `json:"renewable"`
	Errors        []string `json:"errors"`
	ServiceType   string   `json:"serviceTypeLogo"`
}

type BreakResponse struct {
	Credential Credential `json:"credential"`
	RequestID  string     `json:"requestId"`
}

type Credential struct {
	CredentialType string     `json:"credentialType"`
	KeyValue       []KeyValue `json:"keyValue"`
	Name           string     `json:"name"`
	ServiceID      string     `json:"serviceId"`
}

type KeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type SearchCredentialResponse struct {
	SearchCredentials []SearchCredential `json:"credentials"`
}

type SearchCredential struct {
	CredentailId string `json:"credentialId"`
	ServiceId    string `json:"serviceId"`
}

func WSEndpointCreate() *httptest.Server {

	serviceSearchMatcher, _ := regexp.Compile("/api/v1/service/sort/name")

	credentialMatcher, _ := regexp.Compile("/api/v1/cmanager/break/(.*)")

	credentialSearchMatcher, _ := regexp.Compile("/api/v1/cmanager/credential")

	serviceIdMatcher, _ := regexp.Compile("/api/v1/service/id/(.*)")

	return httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := []byte(r.URL.String())
		w.Header().Set("Content-Type", "application/json")

		switch {
		case serviceSearchMatcher.Match(url):
			service := []ServiceResponse{
				{
					ServiceId: MockServiceId,
					IPAddress: MockServiceIp,
					Port:      DummySshServerPort,
				},
			}

			serviceJson, _ := json.Marshal(service)
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, string(serviceJson))

			return

		case serviceIdMatcher.Match(url):
			//serviceId := serviceIdMatcher.FindStringSubmatch(string(url))
			response := ServiceResponse{
				IPAddress: MockServiceIp,
				Port:      DummySshServerPort,
			}

			log.Debug().Interface("dump", response).Msg("dump matcher")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(&response)
			return
		case credentialMatcher.Match(url):
			//serviceId := serviceIdMatcher.FindStringSubmatch(string(url))
			response := BreakResponse{
				Credential: Credential{
					KeyValue: []KeyValue{
						{Key: "password", Value: ValidPass},
					},
					Name: ValidUser,
				},
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(&response)
			return
		case credentialSearchMatcher.Match(url):

			response := SearchCredentialResponse{
				SearchCredentials: []SearchCredential{
					{
						CredentailId: MockCredentialId,
						ServiceId:    MockServiceId,
					},
				},
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(&response)
			return
		}

		w.WriteHeader(http.StatusBadRequest)

	}))
}
