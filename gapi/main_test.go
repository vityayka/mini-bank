package gapi

import (
	"bank/async"
	db "bank/db/sqlc"
	"bank/utils"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

func newTestServer(t *testing.T, store db.Store, taskDistributor async.TaskDistributor) *Server {
	srv, err := NewServer(utils.Config{
		TokenSymmetricKey:   utils.RandomString(32),
		AccessTokenDuration: time.Minute,
	}, store, taskDistributor)
	require.NoError(t, err)
	return srv
}

func newContextWithAuthMetadata(
	t *testing.T,
	server *Server,
	userID int64,
	duration time.Duration,
	authHeader,
	authType string,
) context.Context {
	token, _, err := server.tokenMaker.CreateToken(userID, duration)
	require.NoError(t, err)
	return metadata.NewIncomingContext(
		context.Background(),
		metadata.New(map[string]string{authHeader: fmt.Sprintf("%s %s", authType, token)}),
	)
}
