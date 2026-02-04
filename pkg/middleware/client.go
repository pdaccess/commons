package middleware

import (
	"context"

	"git.h2hsecure.com/pda/commons/pkg/domain"
	"github.com/google/uuid"
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
