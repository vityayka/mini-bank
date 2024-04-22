package api

import (
	mockdb "bank/db/mock"
	db "bank/db/sqlc"
	"bank/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateUser(t *testing.T) {
	password := "password"
	user := randomUser(password)

	testCases := []struct {
		name          string
		params        gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			params: gin.H{
				"username":  user.Username,
				"password":  password,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				args := db.CreateUserParams{
					Username:       user.Username,
					FullName:       user.FullName,
					HashedPassword: user.HashedPassword,
					Email:          user.Email,
				}

				matcher := func(x any) bool {
					if params, isOk := x.(db.CreateUserParams); isOk {
						if err := utils.CompareHashAndPassword(params.HashedPassword, password); err != nil {
							return false
						}
						return params.Username == args.Username &&
							params.FullName == args.FullName &&
							params.Email == args.Email
					}
					return false
				}

				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Cond(matcher)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)

				//check response body
				body := recorder.Body.Bytes()
				var gotUser db.User
				err := json.Unmarshal(body, &gotUser)
				require.NoError(t, err)

				require.Equal(t, user.Email, gotUser.Email)
				require.Equal(t, user.Username, gotUser.Username)
				require.Equal(t, user.FullName, gotUser.FullName)
				require.Empty(t, gotUser.HashedPassword)
			},
		},
	}

	for _, tc := range testCases {
		ctrl := gomock.NewController(t)

		store := mockdb.NewMockStore(ctrl)

		tc.buildStubs(store)

		server := newTestServer(t, store)
		recorder := httptest.NewRecorder()

		payload, _ := json.Marshal(tc.params)
		request, err := http.NewRequest(http.MethodPost, "/users", bytes.NewReader(payload))

		require.NoError(t, err)

		server.router.ServeHTTP(recorder, request)

		tc.checkResponse(t, recorder)
	}
}

func TestLoginUserAPI(t *testing.T) {
	password := "password"
	user := randomUser(password)

	testCases := []struct {
		name          string
		params        gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			params: gin.H{
				"password": password,
				"email":    user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Eq(user.Email)).
					Times(1).
					Return(user, nil)

				store.EXPECT().
					CreateSession(gomock.Any(), gomock.Any()).
					Times(1)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				//check response body
				body := recorder.Body.Bytes()
				var response struct {
					User db.User
				}
				err := json.Unmarshal(body, &response)
				require.NoError(t, err)

				require.Equal(t, user.Email, response.User.Email)
				require.Equal(t, user.Username, response.User.Username)
				require.Equal(t, user.FullName, response.User.FullName)
				require.Empty(t, response.User.HashedPassword)
			},
		},
	}

	for _, tc := range testCases {
		ctrl := gomock.NewController(t)

		store := mockdb.NewMockStore(ctrl)

		tc.buildStubs(store)

		server := newTestServer(t, store)
		recorder := httptest.NewRecorder()

		payload, _ := json.Marshal(tc.params)
		request, err := http.NewRequest(http.MethodPost, "/users/login", bytes.NewReader(payload))

		require.NoError(t, err)

		server.router.ServeHTTP(recorder, request)

		tc.checkResponse(t, recorder)
	}
}

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
