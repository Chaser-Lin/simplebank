package api

import (
	db "SimpleBank/db/sqlc"
	"SimpleBank/token"
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"net/http"
)

type createAccountRequest struct {
	//Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,currency"`
}

func (s *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.CreateAccountParams{
		Owner:    payload.Username,
		Balance:  0.00,
		Currency: req.Currency,
	}

	err := s.store.CreateAccount(ctx, arg)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok { //转换为mysql类型的error进一步判断错误类型
			//log.Println(mysqlErr.Number)
			switch mysqlErr.Number {
			case 1062, 1452: //1062:Duplicate (owner, currency), 1452:Unexist owner references to username
				ctx.JSON(http.StatusForbidden, errResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	account, err := s.store.GetLastAccount(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (s *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	account, err := s.store.GetAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	if account.Owner != payload.Username {
		err := errors.New("account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, account)
}

type listAccountRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (s *Server) listAccount(ctx *gin.Context) {
	var req listAccountRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.ListAccountsParams{
		Owner:  payload.Username,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}
	accounts, err := s.store.ListAccounts(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, accounts)
}

type updateAccountRequest struct {
	ID     int64   `json:"id" binding:"required,min=1"`
	Amount float64 `json:"amount" binding:"required"`
}

func (s *Server) updateAccount(ctx *gin.Context) {
	var req updateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	account, err := s.store.GetAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	if account.Owner != payload.Username {
		err := errors.New("account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errResponse(err))
		return
	}

	arg := db.AddAccountBalanceParams{
		ID:     req.ID,
		Amount: req.Amount,
	}
	err = s.store.AddAccountBalance(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	account.Balance += req.Amount
	ctx.JSON(http.StatusOK, account)

	//ctx.JSON(http.StatusOK, gin.H{
	//	"msg":             "update account success",
	//	"updated_account": account,
	//})

	//ctx.Redirect(http.StatusFound, "/accounts/"+strconv.Itoa(int(arg.ID)))
}

type deleteAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (s *Server) deleteAccount(ctx *gin.Context) {
	var req deleteAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	account, err := s.store.GetAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	if account.Owner != payload.Username {
		err := errors.New("account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errResponse(err))
		return
	}

	err = s.store.DeleteAccount(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "delete account success",
	})
}
