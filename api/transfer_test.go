package api

import (
	mockdb "bank/db/mock"
	db "bank/db/sqlc"
	"bank/token"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateTransfer(t *testing.T) {
	user1 := randomUser("password")
	user2 := randomUser("password1")
	acc1 := randomAccount(user1.ID)
	acc2 := randomAccount(user2.ID)
	// acc3 := randomAccount()
	acc1.Currency, acc2.Currency = "USD", "USD"

	testCases := []struct {
		name            string
		body            gin.H
		setupAuthHeader func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs      func(store *mockdb.MockStore)
		checkResponse   func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"from_account_id": acc1.ID,
				"to_account_id":   acc2.ID,
				"currency":        "USD",
				"amount":          int64(100),
			},
			setupAuthHeader: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authHeaderTypeBearer, acc1.UserID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserAccount(gomock.Any(), gomock.Eq(db.GetUserAccountParams{acc1.UserID, acc1.ID})).Times(1).Return(acc1, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(acc1.ID)).Times(1).Return(acc1, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(acc2.ID)).Times(1).Return(acc2, nil)

				arg := db.TransferTxParams{
					FromAccountID: acc1.ID,
					ToAccountID:   acc2.ID,
					Amount:        int64(100),
				}

				store.EXPECT().TransferTx(gomock.Any(), gomock.Eq(arg)).Times(1)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		ctrl := gomock.NewController(t)

		store := mockdb.NewMockStore(ctrl)

		tc.buildStubs(store)

		server := newTestServer(t, store)
		recorder := httptest.NewRecorder()

		payload, _ := json.Marshal(tc.body)
		request, err := http.NewRequest(http.MethodPost, "/transfers", bytes.NewReader(payload))
		tc.setupAuthHeader(t, request, server.tokenMaker)

		require.NoError(t, err)

		server.router.ServeHTTP(recorder, request)

		tc.checkResponse(t, recorder)
	}

}
