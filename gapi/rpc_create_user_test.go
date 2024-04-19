package gapi

import (
	async "bank/async/mock"
	mockdb "bank/db/mock"
	db "bank/db/sqlc"
	"bank/pb"
	"bank/utils"
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func randomUser(password string) db.User {
	hashedPassword, _ := utils.HashedPassword(password)
	return db.User{
		ID:             utils.RandomInt(1, 1000),
		Username:       utils.RandomString(6),
		HashedPassword: hashedPassword,
		FullName:       fmt.Sprintf("%s %s", utils.RandomString(6), utils.RandomString(6)),
		Email:          utils.RandomEmail(),
	}
}

func TestCreateUser(t *testing.T) {
	password := "password"
	user := randomUser(password)

	testCases := []struct {
		name          string
		params        pb.CreateUserRequest
		buildStubs    func(store *mockdb.MockStore, distributor *async.MockTaskDistributor)
		checkResponse func(t *testing.T, res *pb.CreateUserResponse, err error)
	}{
		{
			name: "OK",
			params: pb.CreateUserRequest{
				Username: user.Username,
				Password: password,
				FullName: user.FullName,
				Email:    user.Email,
			},
			buildStubs: func(store *mockdb.MockStore, distributor *async.MockTaskDistributor) {
				args := db.CreateUserParams{
					Username:       user.Username,
					FullName:       user.FullName,
					HashedPassword: user.HashedPassword,
					Email:          user.Email,
				}

				matcher := func(x any) bool {
					if params, isOk := x.(db.CreateUserTxParams); isOk {
						if err := utils.CompareHashAndPassword(params.HashedPassword, password); err != nil {
							return false
						}
						if err := params.AfterCreate(user); err != nil {
							return false
						}

						return params.Username == args.Username &&
							params.FullName == args.FullName &&
							params.Email == args.Email
					}
					return false
				}

				store.EXPECT().
					CreateUserTX(gomock.Any(), gomock.Cond(matcher)).
					Times(1).
					Return(db.CreateUserTxResult{User: user}, nil)

				distributor.EXPECT().
					DistributeTaskVerifyEmail(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, res *pb.CreateUserResponse, err error) {
				require.NoError(t, err)
				gotUser := res.User

				require.Equal(t, user.Email, gotUser.Email)
				require.Equal(t, user.Username, gotUser.Username)
				require.Equal(t, user.FullName, gotUser.FullName)
			},
		},
		{
			name: "Validation fail",
			params: pb.CreateUserRequest{
				Username: "ba",
				Password: "bad",
				FullName: "ba",
				Email:    "bad_email",
			},
			buildStubs: func(store *mockdb.MockStore, distributor *async.MockTaskDistributor) {},
			checkResponse: func(t *testing.T, res *pb.CreateUserResponse, err error) {
				require.Error(t, err)
				require.Nil(t, res)
			},
		},
		{
			name: "Redis malfunction",
			params: pb.CreateUserRequest{
				Username: user.Username,
				Password: password,
				FullName: user.FullName,
				Email:    user.Email,
			},
			buildStubs: func(store *mockdb.MockStore, distributor *async.MockTaskDistributor) {
				args := db.CreateUserParams{
					Username:       user.Username,
					FullName:       user.FullName,
					HashedPassword: user.HashedPassword,
					Email:          user.Email,
				}

				matcher := func(x any) bool {
					if params, isOk := x.(db.CreateUserTxParams); isOk {
						if err := utils.CompareHashAndPassword(params.HashedPassword, password); err != nil {
							return false
						}

						params.AfterCreate(user)

						return params.Username == args.Username &&
							params.FullName == args.FullName &&
							params.Email == args.Email
					}
					return false
				}

				err := status.Errorf(codes.Internal, "couldn't enqueue task")
				store.EXPECT().
					CreateUserTX(gomock.Any(), gomock.Cond(matcher)).
					Times(1).Return(db.CreateUserTxResult{}, err)

				distributor.EXPECT().
					DistributeTaskVerifyEmail(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(err)
			},
			checkResponse: func(t *testing.T, res *pb.CreateUserResponse, err error) {
				require.Error(t, err)
				require.Nil(t, res)
			},
		},
		{
			name: "Already exists",
			params: pb.CreateUserRequest{
				Username: user.Username,
				Password: password,
				FullName: user.FullName,
				Email:    user.Email,
			},
			buildStubs: func(store *mockdb.MockStore, distributor *async.MockTaskDistributor) {
				store.EXPECT().
					CreateUserTX(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.CreateUserTxResult{}, db.ErrUserAlreadyExists)
			},
			checkResponse: func(t *testing.T, res *pb.CreateUserResponse, err error) {
				require.ErrorContains(t, err, "user already exists")
				require.Nil(t, res)
			},
		},
	}

	for _, tc := range testCases {
		ctrl := gomock.NewController(t)

		store := mockdb.NewMockStore(ctrl)

		distrCtrl := gomock.NewController(t)
		distributor := async.NewMockTaskDistributor(distrCtrl)

		tc.buildStubs(store, distributor)

		server := newTestServer(t, store, distributor)

		res, err := server.CreateUser(context.TODO(), &tc.params)

		tc.checkResponse(t, res, err)
	}
}
