package middleware

import (
	"context"

	"github.com/google/uuid"
	"github.com/pdaccess/commons/pkg/domain"
)

const ClientRequestStr = "requestId"

type ClientRequestType string

type ClientRequest struct {
	User      domain.PdaccessClaims
	RequestId string
}

func ClientFromCtx(ctx context.Context) *ClientRequest {
	clientRequest := ctx.Value(ClientRequestType(ClientRequestStr))

	if clientRequest == nil {
		clientRequest = &ClientRequest{
			RequestId: uuid.NewString(),
		}
	}

	return clientRequest.(*ClientRequest)
}
