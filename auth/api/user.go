package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "shivesh-ranjan.github.io/m/db/sqlc"
	"shivesh-ranjan.github.io/m/utils"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required,min=5"`
	About    string `json:"about" binding:"required"`
	Photo    string `json:"photo" binding:"required"`
}

type UserResponse struct {
	Username  string    `json:"username"`
	Name      string    `json:"name"`
	About     string    `json:"about"`
	Photo     string    `json:"photo"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}
	hashedPassword, perror := utils.HashPassword(req.Password)
	if perror != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("There was a problem while parsing the password")))
	}
	arg := db.CreateUserParams{
		Name:     req.Name,
		Username: req.Username,
		Role:     "user",
		About:    req.About,
		Password: hashedPassword,
		Photo:    req.Photo,
	}
	result, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(errors.New("username already taken.")))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := UserResponse{
		Name:      result.Name,
		Username:  result.Username,
		Role:      result.Role,
		About:     result.About,
		Photo:     result.Photo,
		CreatedAt: result.CreatedAt,
	}
	ctx.JSON(http.StatusOK, res)
}

type getUserRequest struct {
	Username string `form:"username" binding:"required,alphanum"`
}

func (server *Server) getUser(ctx *gin.Context) {
	var req getUserRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := UserResponse{
		Name:      user.Name,
		Username:  user.Username,
		Role:      user.Role,
		About:     user.About,
		Photo:     user.Photo,
		CreatedAt: user.CreatedAt,
	}
	ctx.JSON(http.StatusOK, res)
}

type updatePasswordRequest struct {
	Username string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (server *Server) UpdatePassword(ctx *gin.Context) {
	var req updatePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := server.getUserFromPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	if user.Username != req.Username || user.Role != "admin" {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("You can only update your password.")))
		return
	}
	arg := db.UpdatePasswordParams{
		Username: req.Username,
	}
	arg.Password, _ = utils.HashPassword(req.Password)
	user, err = server.store.UpdatePassword(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := UserResponse{
		Name:      user.Name,
		Username:  user.Username,
		Role:      user.Role,
		About:     user.About,
		Photo:     user.Photo,
		CreatedAt: user.CreatedAt,
	}
	ctx.JSON(http.StatusOK, res)
}

type updateRoleRequest struct {
	Username string `json:"username" binding:"required"`
	Role     string `json:"role" binding:"required"`
}

func (server *Server) UpdateRole(ctx *gin.Context) {
	var req updateRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := server.getUserFromPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	if user.Role != "admin" {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("Only admins can update Role.")))
	}
	arg := db.UpdateRoleParams{
		Username: req.Username,
		Role:     req.Role,
	}
	user, err = server.store.UpdateRole(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := UserResponse{
		Name:      user.Name,
		Username:  user.Username,
		Role:      user.Role,
		About:     user.About,
		Photo:     user.Photo,
		CreatedAt: user.CreatedAt,
	}
	ctx.JSON(http.StatusOK, res)
}

type updateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Name     string `json:"name" binding:"required"`
	About    string `json:"about" binding:"required"`
	Photo    string `json:"photo" binding:"required"`
}

func (server *Server) UpdateUser(ctx *gin.Context) {
	var req updateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := server.getUserFromPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	arg := db.UpdateUserParams{
		Username: req.Username,
		Name:     req.Name,
		About:    req.About,
		Photo:    req.Photo,
	}
	user, err = server.store.UpdateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := UserResponse{
		Name:      user.Name,
		Username:  user.Username,
		Role:      user.Role,
		About:     user.About,
		Photo:     user.Photo,
		CreatedAt: user.CreatedAt,
	}
	ctx.JSON(http.StatusOK, res)
}

type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=5"`
}

type loginUserResponse struct {
	AccessToken string       `json:"access_token"`
	User        UserResponse `json:"user"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	err = utils.CheckPassword(req.Password, user.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	accessToken, err := server.tokenMaker.CreateToken(user.Username, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := UserResponse{
		Name:      user.Name,
		Username:  user.Username,
		Role:      user.Role,
		About:     user.About,
		Photo:     user.Photo,
		CreatedAt: user.CreatedAt,
	}
	rsp := loginUserResponse{
		AccessToken: accessToken,
		User:        res,
	}
	ctx.JSON(http.StatusOK, rsp)
}
