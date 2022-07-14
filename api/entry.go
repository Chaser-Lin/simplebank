package api

import (
	db "SimpleBank/db/sqlc"
	"SimpleBank/token"
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

//var ErrInsufficientBalance = errors.New("account's balance is insufficient")

type newBusinessRequest struct {
	Business  string  `json:"business" binding:"required,business"` //使用ShouldBind，注意定义为包级
	AccountID int64   `json:"account_id" binding:"required,min=1"`
	Amount    float64 `json:"amount" binding:"required,gt=0"`
}

func (s Server) NewBusiness(ctx *gin.Context) {
	var req newBusinessRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}
	account, err := s.store.GetAccount(ctx, req.AccountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	var result db.BusinessTxResult
	if req.Business == "Withdraw" {
		//s.createEntry(ctx, req.AccountID, -req.Amount)
		if payload.Username != account.Owner { //取钱需要是号主
			err := errors.New("account doesn't belong to the authenticated user")
			ctx.JSON(http.StatusUnauthorized, errResponse(err))
			return
		}

		result, err = s.store.WithdrawTx(ctx, db.BusinessTxParms{
			AccountID: req.AccountID,
			Amount:    req.Amount,
		})
		if err != nil {
			if err == db.ErrInsufficientBalance { //不可以直接比较error对象是否相等
				ctx.JSON(http.StatusBadRequest, errResponse(err))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errResponse(err))
			return
		}
	} else if req.Business == "Deposit" { //存钱不需要号主
		//s.createEntry(ctx, req.AccountID, req.Amount)
		result, err = s.store.DepositTx(ctx, db.BusinessTxParms{
			AccountID: req.AccountID,
			Amount:    req.Amount,
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errResponse(err))
			return
		}
	}
	ctx.JSON(http.StatusOK, result)
}

type listEntryRequest struct {
	PageID    int32 `form:"page_id" binding:"required,min=1"`
	PageSize  int32 `form:"page_size" binding:"required,min=5,max=10"`
	AccountID int64 `form:"account_id" binding:"required,min=1"`
}

func (s *Server) listEntries(ctx *gin.Context) {
	var req listEntryRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	account, err := s.store.GetAccount(ctx, req.AccountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if payload.Username != account.Owner {
		err := errors.New("account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errResponse(err))
		return
	}

	entries, err := s.store.GetEntryByAccount(ctx, db.GetEntryByAccountParams{
		AccountID: req.AccountID,
		Limit:     req.PageSize,
		Offset:    (req.PageID - 1) * req.PageSize,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, entries)
}

//func (s *Server) createEntry(ctx *gin.Context, accountID int64, amount float64) {
//	account, err := s.store.GetAccount(ctx, accountID)
//	if err != nil {
//		if err == sql.ErrNoRows {
//			ctx.JSON(http.StatusNotFound, errResponse(err))
//			return
//		}
//		ctx.JSON(http.StatusInternalServerError, errResponse(err))
//		return
//	}
//
//	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
//	if payload.Username != account.Owner {
//		err := errors.New("account doesn't belong to the authenticated user")
//		ctx.JSON(http.StatusBadRequest, errResponse(err))
//		return
//	}
//
//	err = s.store.CreateEntry(ctx, db.CreateEntryParams{
//		AccountID: accountID,
//		Amount:    amount,
//	})
//	if err != nil {
//		ctx.JSON(http.StatusInternalServerError, errResponse(err))
//		return
//	}
//
//	entry, err := s.store.GetLastEntry(ctx)
//
//	if err != nil {
//		ctx.JSON(http.StatusInternalServerError, errResponse(err))
//		return
//	}
//
//	ctx.JSON(http.StatusOK, entry)
//}
