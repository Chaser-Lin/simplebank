package api

import (
	mockdb "SimpleBank/db/mock"
	db "SimpleBank/db/sqlc"
	"SimpleBank/db/util"
	"SimpleBank/token"
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCreateTransferAPI(t *testing.T) {
	user1, _ := randomUser(t)
	user2, _ := randomUser(t)
	fromAccount := randomAccount(user1.Username)
	toAccount1 := randomAccount(user2.Username) //测试currency正确样例
	toAccount2 := randomAccount(user2.Username) //测试currency错误样例
	toAccount1.Currency = fromAccount.Currency
	for toAccount2.Currency == fromAccount.Currency {
		toAccount2.Currency = util.RandomCurrency()
	}
	amount := util.RandomBalance()

	result := db.TransferTxResult{
		FromAccount: db.Account{
			ID:       fromAccount.ID,
			Owner:    fromAccount.Owner,
			Balance:  fromAccount.Balance - amount,
			Currency: fromAccount.Currency,
		},
		ToAccount: db.Account{
			ID:       toAccount1.ID,
			Owner:    toAccount1.Owner,
			Balance:  toAccount1.Balance + amount,
			Currency: toAccount1.Currency,
		},
		FromEntry: db.Entry{
			AccountID: fromAccount.ID,
			Amount:    -amount,
		},
		ToEntry: db.Entry{
			AccountID: toAccount1.ID,
			Amount:    amount,
		},
		Transfer: db.Transfer{
			FromAccountID: fromAccount.ID,
			ToAccountID:   toAccount1.ID,
			Amount:        amount,
		},
	}

	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, maker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"from_account_id": fromAccount.ID,
				"to_account_id":   toAccount1.ID,
				"amount":          amount,
				"currency":        fromAccount.Currency,
			},
			setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
				addAuthorization(t, request, maker, authorizationTypeBearer, fromAccount.Owner, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(fromAccount.ID)).
					Times(1).
					Return(fromAccount, nil)
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(toAccount1.ID)).
					Times(1).
					Return(toAccount1, nil)
				arg := db.TransferTxParms{
					FromAccountID: fromAccount.ID,
					ToAccountID:   toAccount1.ID,
					Amount:        amount,
				}
				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(result, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchResult(t, recorder.Body, result)
			},
		},
		{
			name: "UnAuthorizedUser",
			body: gin.H{
				"from_account_id": fromAccount.ID,
				"to_account_id":   toAccount1.ID,
				"amount":          amount,
				"currency":        fromAccount.Currency,
			},
			setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
				addAuthorization(t, request, maker, authorizationTypeBearer, "unauthorized_user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(fromAccount.ID)).
					Times(1).
					Return(fromAccount, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InternalServerError",
			body: gin.H{
				"from_account_id": fromAccount.ID,
				"to_account_id":   toAccount1.ID,
				"amount":          amount,
				"currency":        fromAccount.Currency,
			},
			setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
				addAuthorization(t, request, maker, authorizationTypeBearer, fromAccount.Owner, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(fromAccount.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidCurrency",
			body: gin.H{
				"from_account_id": fromAccount.ID,
				"to_account_id":   toAccount1.ID,
				"amount":          amount,
				"currency":        toAccount2.Currency,
			},
			setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
				addAuthorization(t, request, maker, authorizationTypeBearer, fromAccount.Owner, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(fromAccount.ID)).
					Times(1).
					Return(fromAccount, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidFromAccountID",
			body: gin.H{
				"from_account_id": 0,
				"to_account_id":   toAccount1.ID,
				"amount":          amount,
				"currency":        fromAccount.Currency,
			},
			setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
				addAuthorization(t, request, maker, authorizationTypeBearer, fromAccount.Owner, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidToAccountID",
			body: gin.H{
				"from_account_id": fromAccount.ID,
				"to_account_id":   0,
				"amount":          amount,
				"currency":        fromAccount.Currency,
			},
			setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
				addAuthorization(t, request, maker, authorizationTypeBearer, fromAccount.Owner, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidAmount",
			body: gin.H{
				"from_account_id": fromAccount.ID,
				"to_account_id":   toAccount1.ID,
				"amount":          -amount,
				"currency":        fromAccount.Currency,
			},
			setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
				addAuthorization(t, request, maker, authorizationTypeBearer, fromAccount.Owner, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "NotFound",
			body: gin.H{
				"from_account_id": fromAccount.ID,
				"to_account_id":   toAccount1.ID,
				"amount":          amount,
				"currency":        fromAccount.Currency,
			},
			setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
				addAuthorization(t, request, maker, authorizationTypeBearer, fromAccount.Owner, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := NewTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := "/transfers"
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func requireBodyMatchResult(t *testing.T, body *bytes.Buffer, result db.TransferTxResult) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotResult db.TransferTxResult
	err = json.Unmarshal(data, &gotResult)
	require.NoError(t, err)

	//result.FromEntry.ID = gotResult.FromEntry.ID
	//result.FromEntry.CreatedAt = gotResult.FromEntry.CreatedAt
	//result.ToEntry.ID = gotResult.ToEntry.ID
	//result.ToEntry.CreatedAt = gotResult.ToEntry.CreatedAt
	//result.Transfer.ID = gotResult.Transfer.ID
	//result.Transfer.CreatedAt = gotResult.Transfer.CreatedAt
	require.Equal(t, result, gotResult)
}
