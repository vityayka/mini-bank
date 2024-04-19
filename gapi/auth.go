package gapi

import (
	"bank/token"
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc/metadata"
)

const (
	authHeader                  = "authorization"
	authBearer                  = "bearer"
	msgErrMetadata              = "failed to fetch metadata"
	msgErrAuthHeaderMissing     = "auth header is missing"
	msgErrAuthHeaderCorrupted   = "auth header is bad"
	msgErrAuthHeaderUnsupported = "unsupported auth scheme"
)

func (server *Server) authorizeUser(ctx context.Context) (*token.Payload, error) {
	md, isOK := metadata.FromIncomingContext(ctx)
	if !isOK {
		return nil, fmt.Errorf(msgErrMetadata)
	}

	auth := md.Get(authHeader)
	if len(auth) == 0 {
		return nil, fmt.Errorf(msgErrAuthHeaderMissing)
	}

	authFields := strings.Fields(auth[0])
	if len(authFields) != 2 {
		return nil, fmt.Errorf(msgErrAuthHeaderCorrupted)
	}

	if strings.ToLower(authFields[0]) != authBearer {
		return nil, fmt.Errorf(msgErrAuthHeaderUnsupported)
	}

	payload, err := server.tokenMaker.VerifyToken(authFields[1])
	if err != nil {
		return nil, fmt.Errorf("failed to auth: %s", err.Error())
	}
	return payload, nil
}
