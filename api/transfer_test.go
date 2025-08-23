package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	mockdb "github.com/hykura1501/simple_bank/db/mock"
	db "github.com/hykura1501/simple_bank/db/sqlc"
	"github.com/hykura1501/simple_bank/ulti"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func TestCreateTransfer(t *testing.T) {
	acc1 := randomAccount()
	acc2 := randomAccount()

	fromAcc := acc1
	toAcc := acc2
	amount := int64(10)
	currency := ulti.USD

	fromAcc.Currency = ulti.USD
	toAcc.Currency = ulti.USD

	fromAcc.Balance -= amount
	toAcc.Balance += amount

	toEntry := db.Entry{
		AccountID: toAcc.ID,
		Amount:    amount,
	}

	fromEntry := db.Entry{
		AccountID: fromAcc.ID,
		Amount:    -amount,
	}

	transfer := db.Transfer{
		FromAccountID: fromAcc.ID,
		ToAccountID:   toAcc.ID,
		Amount:        amount,
	}

	req := transferRequest{
		FromAccountID: fromAcc.ID,
		ToAccountID:   toAcc.ID,
		Amount:        amount,
		Currency:      currency,
	}

	invalidIDReq := transferRequest{
		FromAccountID: 0,
		ToAccountID:   0,
		Amount:        amount,
		Currency:      currency,
	}

	invalidAmountReq := transferRequest{
		FromAccountID: fromAcc.ID,
		ToAccountID:   toAcc.ID,
		Amount:        -amount,
		Currency:      currency,
	}

	invalidCurrencyReq := transferRequest{
		FromAccountID: fromAcc.ID,
		ToAccountID:   toAcc.ID,
		Amount:        amount,
		Currency:      ulti.VND,
	}

	argTransfer := db.CreateTransferParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	result := db.TransferTxResult{
		Transfer:    transfer,
		FromAccount: fromAcc,
		ToAccount:   toAcc,
		ToEntry:     toEntry,
		FromEntry:   fromEntry,
	}

	testCases := []struct {
		name            string
		transferRequest transferRequest
		buildStubs      func(store *mockdb.MockStore)
		checkResponse   func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:            "OK",
			transferRequest: req,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), fromAcc.ID).Times(1).Return(fromAcc, nil)
				store.EXPECT().GetAccount(gomock.Any(), toAcc.ID).Times(1).Return(toAcc, nil)

				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Eq(argTransfer)).
					Times(1).
					Return(result, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchTransferTxResult(t, recorder.Body, result)
			},
		},
		{
			name:            "InvalidID",
			transferRequest: invalidIDReq,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)

				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:            "InvalidAmount",
			transferRequest: invalidAmountReq,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)

				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:            "NotFoundFromAcc",
			transferRequest: req,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), fromAcc.ID).Times(1).Return(db.Account{}, pgx.ErrNoRows)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:            "NotFoundToAcc",
			transferRequest: req,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), fromAcc.ID).Times(1).Return(fromAcc, nil)
				store.EXPECT().GetAccount(gomock.Any(), toAcc.ID).Times(1).Return(db.Account{}, pgx.ErrNoRows)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:            "InternalError_1",
			transferRequest: req,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), fromAcc.ID).Times(1).Return(db.Account{}, sql.ErrConnDone)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:            "InternalError_2",
			transferRequest: req,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), fromAcc.ID).Times(1).Return(fromAcc, nil)
				store.EXPECT().GetAccount(gomock.Any(), toAcc.ID).Times(1).Return(toAcc, nil)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Eq(argTransfer)).Times(1).Return(db.TransferTxResult{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:            "InvalidCurrency",
			transferRequest: invalidCurrencyReq,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), fromAcc.ID).Times(1).Return(fromAcc, nil)
				// store.EXPECT().GetAccount(gomock.Any(), toAcc.ID).Times(0)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for _, _tc := range testCases {
		t.Run(_tc.name, func(t *testing.T) {
			tc := _tc
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			// build stubs
			tc.buildStubs(store)
			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := "/transfers"
			bodyData, err := json.Marshal(tc.transferRequest)
			require.NoError(t, err)
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(bodyData))

			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func requireBodyMatchTransferTxResult(t *testing.T, body *bytes.Buffer, result db.TransferTxResult) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotTransferResult db.TransferTxResult
	err = json.Unmarshal(data, &gotTransferResult)
	require.NoError(t, err)

	require.Equal(t, gotTransferResult, result)
}
