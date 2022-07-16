package api

import (
	mockdb "SimpleBank/db/mock"
	db "SimpleBank/db/sqlc"
	"SimpleBank/token"
	"SimpleBank/util"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewBusinessAPI(t *testing.T) {
	user, _ := randomUser(t)
	account := randomAccount(user.Username)
	amount := util.RandomBalance()

	result := db.BusinessTxResult{
		Account: db.Account{
			ID:       account.ID,
			Owner:    account.Owner,
			Balance:  account.Balance - amount,
			Currency: account.Currency,
		},
		Entry: db.Entry{
			AccountID: account.ID,
			Amount:    amount,
		},
	}

	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, maker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recoder *httptest.ResponseRecorder)
	}{
		{
			name: "WithdrawOK",
			body: gin.H{
				"business":   "Withdraw",
				"account_id": account.ID,
				"amount":     amount,
			},
			setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
				addAuthorization(t, request, maker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
				arg := db.BusinessTxParms{
					AccountID: account.ID,
					Amount:    amount,
				}
				store.EXPECT().
					WithdrawTx(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(result, nil)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recoder.Code)
				requireBodyMatchBusinessResult(t, recoder.Body, result)
			},
		},
		{
			name: "DepositOK",
			body: gin.H{
				"business":   "Deposit",
				"account_id": account.ID,
				"amount":     amount,
			},
			setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
				addAuthorization(t, request, maker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
				arg := db.BusinessTxParms{
					AccountID: account.ID,
					Amount:    amount,
				}
				store.EXPECT().
					DepositTx(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(result, nil)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recoder.Code)
				requireBodyMatchBusinessResult(t, recoder.Body, result)
			},
		},
		{
			name: "InvalidBusiness",
			body: gin.H{
				"business":   "invalid-business",
				"account_id": account.ID,
				"amount":     amount,
			},
			setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
				addAuthorization(t, request, maker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recoder.Code)
			},
		},
		{
			name: "NoAuthorization",
			body: gin.H{},
			setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recoder.Code)
			},
		},
		{
			name: "UnauthorizedUser",
			body: gin.H{
				"business":   "Withdraw",
				"account_id": account.ID,
				"amount":     amount,
			},
			setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
				addAuthorization(t, request, maker, authorizationTypeBearer, "unauthorized_user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recoder.Code)
			},
		},
		{
			name: "NotFound",
			body: gin.H{
				"business":   "Withdraw",
				"account_id": account.ID,
				"amount":     amount,
			},
			setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
				addAuthorization(t, request, maker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recoder.Code)
			},
		},
		{
			name: "WithdrawInsufficientBalance",
			body: gin.H{
				"business":   "Withdraw",
				"account_id": account.ID,
				"amount":     amount,
			},
			setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
				addAuthorization(t, request, maker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
				arg := db.BusinessTxParms{
					AccountID: account.ID,
					Amount:    amount,
				}
				store.EXPECT().
					WithdrawTx(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.BusinessTxResult{}, db.ErrInsufficientBalance)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recoder.Code)
			},
		},
		{
			name: "InternalGetAcountServerError",
			body: gin.H{
				"business":   "Withdraw",
				"account_id": account.ID,
				"amount":     amount,
			},
			setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
				addAuthorization(t, request, maker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recoder.Code)
			},
		},
		{
			name: "WithdrawInternalServerError",
			body: gin.H{
				"business":   "Withdraw",
				"account_id": account.ID,
				"amount":     amount,
			},
			setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
				addAuthorization(t, request, maker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
				arg := db.BusinessTxParms{
					AccountID: account.ID,
					Amount:    amount,
				}
				store.EXPECT().
					WithdrawTx(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.BusinessTxResult{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recoder.Code)
			},
		},
		{
			name: "DepositInternalServerError",
			body: gin.H{
				"business":   "Deposit",
				"account_id": account.ID,
				"amount":     amount,
			},
			setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
				addAuthorization(t, request, maker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
				arg := db.BusinessTxParms{
					AccountID: account.ID,
					Amount:    amount,
				}
				store.EXPECT().
					DepositTx(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.BusinessTxResult{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recoder.Code)
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

			url := "/business"
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

func TestListEntriesAPI(t *testing.T) {
	user, _ := randomUser(t)
	account := randomAccount(user.Username)

	n := 5
	var entries []db.Entry
	for i := 0; i < n; i++ {
		entries = append(entries, randomEntry(account.ID))
	}

	type Query struct {
		pageID   int32
		pageSize int32
	}

	testCases := []struct {
		name          string
		query         Query
		accountID     int64
		setupAuth     func(t *testing.T, request *http.Request, maker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recoder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			query:     Query{pageID: 1, pageSize: 10},
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
				addAuthorization(t, request, maker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
				query := Query{pageID: 1, pageSize: 10}
				arg := db.GetEntryByAccountParams{
					AccountID: account.ID,
					Limit:     query.pageSize,
					Offset:    (query.pageID - 1) * query.pageSize,
				}
				store.EXPECT().
					GetEntryByAccount(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(entries, nil)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recoder.Code)
				requireBodyMatchEntries(t, recoder.Body, entries)
			},
		},
		{
			name:      "InvalidAccountID",
			query:     Query{pageID: 1, pageSize: 10},
			accountID: 0,
			setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
				addAuthorization(t, request, maker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recoder.Code)
			},
		},
		{
			name:      "NotFound",
			query:     Query{pageID: 1, pageSize: 10},
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
				addAuthorization(t, request, maker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recoder.Code)
			},
		},
		{
			name:      "InternalGetAccountServerError",
			query:     Query{pageID: 1, pageSize: 10},
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
				addAuthorization(t, request, maker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recoder.Code)
			},
		},
		{
			name:      "InternalGetEntriesByAccountServerError",
			query:     Query{pageID: 1, pageSize: 10},
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
				addAuthorization(t, request, maker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
				query := Query{pageID: 1, pageSize: 10}
				arg := db.GetEntryByAccountParams{
					AccountID: account.ID,
					Limit:     query.pageSize,
					Offset:    (query.pageID - 1) * query.pageSize,
				}
				store.EXPECT().
					GetEntryByAccount(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return([]db.Entry{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recoder.Code)
			},
		},
		{
			name:      "UnAuthorizedUser",
			query:     Query{pageID: 1, pageSize: 10},
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
				addAuthorization(t, request, maker, authorizationTypeBearer, "unauthorized_user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recoder.Code)
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

			url := fmt.Sprintf("/business?page_id=%d&page_size=%d&account_id=%d", tc.query.pageID, tc.query.pageSize, tc.accountID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func requireBodyMatchBusinessResult(t *testing.T, body *bytes.Buffer, result db.BusinessTxResult) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotResult db.BusinessTxResult
	err = json.Unmarshal(data, &gotResult)
	require.NoError(t, err)

	require.Equal(t, result, gotResult)
}

func requireBodyMatchEntries(t *testing.T, body *bytes.Buffer, entries []db.Entry) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotEntries []db.Entry
	err = json.Unmarshal(data, &gotEntries)
	require.NoError(t, err)

	for i := range entries {
		require.Equal(t, entries[i], gotEntries[i])
	}
}

func randomEntry(accountID int64) db.Entry {
	return db.Entry{
		ID:        util.RandomID(),
		AccountID: accountID,
		Amount:    util.RandomAmount(),
	}
}
