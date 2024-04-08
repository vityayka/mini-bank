package gapi

import (
	"bank/token"
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc/metadata"
)

const (
	authHeader = "authorization"
	authBearer = "bearer"
)

func (server *Server) authorizeUser(ctx context.Context) (*token.Payload, error) {
	md, isOK := metadata.FromIncomingContext(ctx)
	if !isOK {
		return nil, fmt.Errorf("failed to fetch metadata")
	}

	auth := md.Get(authHeader)
	if len(auth) == 0 {
		return nil, fmt.Errorf("auth header is missing")
	}

	authFields := strings.Fields(auth[0])
	if len(authFields) != 2 {
		return nil, fmt.Errorf("auth header is bad")
	}

	if strings.ToLower(authFields[0]) != authBearer {
		return nil, fmt.Errorf("unsupported auth scheme")
	}

	payload, err := server.tokenMaker.VerifyToken(authFields[1])
	if err != nil {
		return nil, fmt.Errorf("failed to auth: %s", err.Error())
	}
	return payload, nil
}
