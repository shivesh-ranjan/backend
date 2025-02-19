package api

import (
	"fmt"

	"github.com/gin-gonic/gin"

	db "shivesh-ranjan.github.io/m/db/sqlc"
	"shivesh-ranjan.github.io/m/token"
	"shivesh-ranjan.github.io/m/utils"
)

// Server serves HTTP requests for our auth service
type Server struct {
	config     utils.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config utils.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}
	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()
	router.Use(CORSMiddleware())

	router.POST("/", server.createUser)
	router.POST("/login", server.loginUser)
	router.GET("/", server.getUser)

	// ======================================================
	router.GET("/blog/*proxyPath", func(ctx *gin.Context) {
		targetURL := server.config.BlogMicroURL
		proxyRequest(ctx, targetURL, "")
	})
	// ======================================================

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	authRoutes.POST("/role", server.CreateRole)
	authRoutes.DELETE("/role", server.DeleteRole)
	authRoutes.PUT("/", server.UpdateUser)
	authRoutes.PUT("/role", server.UpdateRole)
	authRoutes.PUT("/login", server.UpdatePassword)

	// ======================================================
	authRoutes.POST("/blog/*proxyPath", func(ctx *gin.Context) {
		targetURL := server.config.BlogMicroURL
		user, _ := server.getUserFromPayload(ctx)
		proxyRequest(ctx, targetURL, user.Username)
	})
	authRoutes.PUT("/blog/*proxyPath", func(ctx *gin.Context) {
		targetURL := server.config.BlogMicroURL
		user, _ := server.getUserFromPayload(ctx)
		proxyRequest(ctx, targetURL, user.Username)
	})
	authRoutes.DELETE("/blog/*proxyPath", func(ctx *gin.Context) {
		targetURL := server.config.BlogMicroURL
		user, _ := server.getUserFromPayload(ctx)
		proxyRequest(ctx, targetURL, user.Username)
	})
	// ======================================================

	server.router = router
}

// Start runs the HTTP Server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error EDIT": err.Error()}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
