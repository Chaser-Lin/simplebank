package api

import (
	db "SimpleBank/db/sqlc"
	"SimpleBank/db/util"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"net/http"
	"time"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"email"`
}

type UserResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func (s *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password) // 自定义比较器后可以发现"abc"错误
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	//arg = db.CreateUserParams{} //自定义比较器后可以判断出错

	err = s.store.CreateUser(ctx, arg)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok { //转换为mysql类型的error进一步判断错误类型
			switch mysqlErr.Number {
			case 1062: //1062:Duplicate (username)
				ctx.JSON(http.StatusForbidden, errResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	//ctx.Redirect(http.StatusFound, fmt.Sprintf("users/%s", arg.Username))

	user, err := s.store.GetUser(ctx, arg.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, user)

	//rsp := UserResponse{
	//	Username:          user.Username,
	//	FullName:          user.FullName,
	//	Email:             user.Email,
	//	PasswordChangedAt: user.PasswordChangedAt,
	//	CreatedAt:         user.CreatedAt,
	//}

	//ctx.JSON(http.StatusOK, rsp)
}

type getUserRequest struct {
	Username string `uri:"username" binding:"required,alphanum"`
}

func (s *Server) getUser(ctx *gin.Context) {
	var req getUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}
	user, err := s.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	//rsp := UserResponse{
	//	Username:          user.Username,
	//	FullName:          user.FullName,
	//	Email:             user.Email,
	//	PasswordChangedAt: user.PasswordChangedAt,
	//	CreatedAt:         user.CreatedAt,
	//}

	ctx.JSON(http.StatusOK, user)
}
