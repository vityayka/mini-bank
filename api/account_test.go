package api

import (
	mockdb "bank/db/mock"
	db "bank/db/sqlc"
	"bank/token"
	"bank/utils"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetAccountApi(t *testing.T) {
	user := randomUser("password")
	account := randomAccount(user.ID)

	testCases := []struct {
		name            string
		accountID       int64
		setupAuthHeader func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs      func(store *mockdb.MockStore)
		checkResponse   func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			setupAuthHeader: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authHeaderTypeBearer, account.UserID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserAccount(gomock.Any(), gomock.Eq(db.GetUserAccountParams{UserID: account.UserID, ID: account.ID})).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				//check response body
				body := recorder.Body.Bytes()
				var gotAccount db.Account
				err := json.Unmarshal(body, &gotAccount)
				require.NoError(t, err)

				require.Equal(t, account, gotAccount)
			},
		},
		{
			name:      "Not found",
			accountID: account.ID,
			setupAuthHeader: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authHeaderTypeBearer, account.UserID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserAccount(gomock.Any(), gomock.Eq(db.GetUserAccountParams{UserID: account.UserID, ID: account.ID})).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)

				//check response body
				body := recorder.Body.Bytes()
				var gotAccount db.Account
				err := json.Unmarshal(body, &gotAccount)
				require.NoError(t, err)

				require.Equal(t, db.Account{}, gotAccount)
			},
		},
		{
			name:      "Foreign token",
			accountID: account.ID,
			setupAuthHeader: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authHeaderTypeBearer, account.UserID+100, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "Internal server error",
			accountID: account.ID,
			setupAuthHeader: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authHeaderTypeBearer, account.UserID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserAccount(gomock.Any(), gomock.Eq(db.GetUserAccountParams{UserID: account.UserID, ID: account.ID})).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "Invalid ID",
			accountID: -1,
			setupAuthHeader: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authHeaderTypeBearer, account.UserID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		ctrl := gomock.NewController(t)

		store := mockdb.NewMockStore(ctrl)

		tc.buildStubs(store)

		server := newTestServer(t, store)
		recorder := httptest.NewRecorder()

		url := fmt.Sprintf("/accounts/%d", tc.accountID)
		request, err := http.NewRequest(http.MethodGet, url, nil)
		tc.setupAuthHeader(t, request, server.tokenMaker)

		require.NoError(t, err)

		server.router.ServeHTTP(recorder, request)

		tc.checkResponse(t, recorder)
	}
}

func TestCreateAccount(t *testing.T) {
	user := randomUser("password")
	account := randomAccount(user.ID)

	arg := db.CreateAccountParams{
		Owner:    account.Owner,
		UserID:   account.UserID,
		Balance:  0,
		Currency: account.Currency,
	}

	testCases := []struct {
		name            string
		params          gin.H
		setupAuthHeader func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs      func(store *mockdb.MockStore)
		checkResponse   func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			params: gin.H{
				"owner":    account.Owner,
				"currency": account.Currency,
			},
			setupAuthHeader: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authHeaderTypeBearer, account.UserID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				matcher := func(x any) bool {
					if params, isOk := x.(db.CreateAccountParams); isOk {
						return params.Balance == arg.Balance && params.UserID == arg.UserID
					}
					return false
				}

				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Cond(matcher)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)

				//check response body
				body := recorder.Body.Bytes()
				var gotAccount db.Account
				err := json.Unmarshal(body, &gotAccount)
				require.NoError(t, err)

				require.Equal(t, account, gotAccount)
			},
		},
		{
			name: "DB error",
			params: gin.H{
				"owner":    account.Owner,
				"currency": account.Currency,
			},
			setupAuthHeader: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authHeaderTypeBearer, account.UserID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(account, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Validation fail",
			params: gin.H{
				"owner":    account.Owner,
				"currency": "RUB",
			},
			setupAuthHeader: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authHeaderTypeBearer, account.UserID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
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
		request, err := http.NewRequest(http.MethodPost, "/accounts", bytes.NewReader(payload))

		require.NoError(t, err)

		tc.setupAuthHeader(t, request, server.tokenMaker)

		server.router.ServeHTTP(recorder, request)

		tc.checkResponse(t, recorder)
	}
}

func TestListAccounts(t *testing.T) {
	user := randomUser("password")
	accounts := []db.Account{
		randomAccount(user.ID),
		randomAccount(user.ID),
		randomAccount(user.ID),
		randomAccount(user.ID),
		randomAccount(user.ID),
	}

	testCases := []struct {
		name   string
		params struct {
			PageNum int64
			Count   int64
		}
		setupAuthHeader func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs      func(store *mockdb.MockStore)
		checkResponse   func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			params: struct {
				PageNum int64
				Count   int64
			}{1, 5},
			setupAuthHeader: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authHeaderTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Eq(db.ListAccountsParams{UserID: user.ID, Limit: 5, Offset: 0})).
					Times(1).
					Return(accounts, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				//check response body
				body := recorder.Body.Bytes()
				var gotAccounts []db.Account
				err := json.Unmarshal(body, &gotAccounts)
				require.NoError(t, err)

				require.Equal(t, accounts, gotAccounts)
			},
		},
		{
			name: "Validation fail",
			params: struct {
				PageNum int64
				Count   int64
			}{1, 11},
			setupAuthHeader: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authHeaderTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Any).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "DB error",
			params: struct {
				PageNum int64
				Count   int64
			}{2, 5},
			setupAuthHeader: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authHeaderTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Eq(db.ListAccountsParams{UserID: user.ID, Limit: 5, Offset: 5})).
					Times(1).
					Return([]db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		ctrl := gomock.NewController(t)

		store := mockdb.NewMockStore(ctrl)

		tc.buildStubs(store)

		server := newTestServer(t, store)
		recorder := httptest.NewRecorder()

		url := fmt.Sprintf("/accounts?page_num=%d&count=%d", tc.params.PageNum, tc.params.Count)

		request, err := http.NewRequest(http.MethodGet, url, nil)
		tc.setupAuthHeader(t, request, server.tokenMaker)

		require.NoError(t, err)

		server.router.ServeHTTP(recorder, request)

		tc.checkResponse(t, recorder)
	}
}

func randomAccount(userID int64) db.Account {
	return db.Account{
		ID:       utils.RandomInt(1, 1000),
		Owner:    utils.RandomName(),
		Currency: utils.RandomCurrency(),
		UserID:   userID,
	}
}
