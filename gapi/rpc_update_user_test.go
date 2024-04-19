package gapi

import (
	mockdb "bank/db/mock"
	db "bank/db/sqlc"
	"bank/pb"
	"bank/utils"
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestUpdateUser(t *testing.T) {
	password := "password"
	user := randomUser(password)

	newName := utils.RandomString(10)
	newEmail := utils.RandomEmail()

	testCases := []struct {
		name          string
		params        pb.UpdateUserRequest
		buildStubs    func(store *mockdb.MockStore)
		makeContext   func(server *Server) context.Context
		checkResponse func(t *testing.T, res *pb.UpdateUserResponse, err error)
	}{
		{
			name: "OK",
			params: pb.UpdateUserRequest{
				FullName: &newName,
				Email:    &newEmail,
			},
			buildStubs: func(store *mockdb.MockStore) {
				args := db.UpdateUserParams{
					ID:       user.ID,
					FullName: sql.NullString{String: newName, Valid: true},
					Email:    sql.NullString{String: newEmail, Valid: true},
				}

				updatedUser := db.User{
					ID:             user.ID,
					Username:       user.Username,
					HashedPassword: user.HashedPassword,
					FullName:       newName,
					Email:          newEmail,
				}

				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Eq(args)).
					Times(1).
					Return(updatedUser, nil)
			},
			makeContext: func(server *Server) context.Context {
				return newContextWithAuthMetadata(t, server, user.ID, time.Minute, authHeader, authBearer)
			},
			checkResponse: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
				require.NoError(t, err)
				gotUser := res.User

				require.Equal(t, newEmail, gotUser.Email)
				require.Equal(t, user.Username, gotUser.Username)
				require.Equal(t, newName, gotUser.FullName)
			},
		},
		{
			name: "Token expired",
			params: pb.UpdateUserRequest{
				FullName: &newName,
				Email:    &newEmail,
			},
			buildStubs: func(store *mockdb.MockStore) {
			},
			makeContext: func(server *Server) context.Context {
				return newContextWithAuthMetadata(t, server, user.ID, -time.Minute, authHeader, authBearer)
			},
			checkResponse: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
				require.Error(t, err)
			},
		},
		{
			name: "Wrong auth type",
			params: pb.UpdateUserRequest{
				FullName: &newName,
				Email:    &newEmail,
			},
			buildStubs: func(store *mockdb.MockStore) {
			},
			makeContext: func(server *Server) context.Context {
				return newContextWithAuthMetadata(t, server, user.ID, time.Minute, authHeader, "basic")
			},
			checkResponse: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
				require.Error(t, err)
				require.ErrorContains(t, err, msgErrAuthHeaderUnsupported)
			},
		},
		{
			name: "User not found",
			params: pb.UpdateUserRequest{
				FullName: &newName,
				Email:    &newEmail,
			},
			buildStubs: func(store *mockdb.MockStore) {
				args := db.UpdateUserParams{
					ID:       user.ID,
					FullName: sql.NullString{String: newName, Valid: true},
					Email:    sql.NullString{String: newEmail, Valid: true},
				}

				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Eq(args)).
					Times(1).
					Return(db.User{}, sql.ErrNoRows)
			},
			makeContext: func(server *Server) context.Context {
				return newContextWithAuthMetadata(t, server, user.ID, time.Minute, authHeader, authBearer)
			},
			checkResponse: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
				require.Error(t, err)
				require.ErrorContains(t, err, "user not found")
			},
		},
		{
			name: "DB not responding",
			params: pb.UpdateUserRequest{
				FullName: &newName,
				Email:    &newEmail,
			},
			buildStubs: func(store *mockdb.MockStore) {
				args := db.UpdateUserParams{
					ID:       user.ID,
					FullName: sql.NullString{String: newName, Valid: true},
					Email:    sql.NullString{String: newEmail, Valid: true},
				}

				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Eq(args)).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			makeContext: func(server *Server) context.Context {
				return newContextWithAuthMetadata(t, server, user.ID, time.Minute, authHeader, authBearer)
			},
			checkResponse: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
				require.Error(t, err)
				require.ErrorContains(t, err, "failed to update user")
			},
		},
	}

	for _, tc := range testCases {
		ctrl := gomock.NewController(t)

		store := mockdb.NewMockStore(ctrl)

		tc.buildStubs(store)

		server := newTestServer(t, store, nil)

		ctx := tc.makeContext(server)

		res, err := server.UpdateUser(ctx, &tc.params)

		tc.checkResponse(t, res, err)
	}
}
