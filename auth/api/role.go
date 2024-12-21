package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type createRoleRequest struct {
	Role string `json:"role" binding:"required"`
}

func (server *Server) CreateRole(ctx *gin.Context) {
	var req createRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	role, err := server.store.CreateRole(ctx, req.Role)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusCreated, role)
}

type deleteRoleRequest struct {
	Role string `json:"role" binding:"required"`
}

func (server *Server) DeleteRole(ctx *gin.Context) {
	var req deleteRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	err := server.store.DeleteRole(ctx, req.Role)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, req.Role)
}
