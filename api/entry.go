package api

import (
	db "SimpleBank/db/sqlc"
	"SimpleBank/token"
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type newBusinessRequest struct {
	Business  string  `json:"business" binding:"required,business"` //使用bind，注意定义为包级
	AccountID int64   `json:"account_id" binding:"required,min=1"`
	Amount    float64 `json:"amount" binding:"required,gt=0"`
}

func (s Server) NewBusiness(ctx *gin.Context) {
	var req newBusinessRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}
	if req.Business == "Withdraw" {
		s.createEntry(ctx, req.AccountID, -req.Amount)
	} else if req.Business == "Deposit" {
		s.createEntry(ctx, req.AccountID, req.Amount)
	}
}

func (s *Server) createEntry(ctx *gin.Context, accountID int64, amount float64) {
	account, err := s.store.GetAccount(ctx, accountID)
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
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	err = s.store.CreateEntry(ctx, db.CreateEntryParams{
		AccountID: accountID,
		Amount:    amount,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	entry, err := s.store.GetLastEntry(ctx)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, entry)
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
		ctx.JSON(http.StatusBadRequest, errResponse(err))
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
