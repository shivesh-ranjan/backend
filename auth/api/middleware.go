package api

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	db "shivesh-ranjan.github.io/m/db/sqlc"
	"shivesh-ranjan.github.io/m/token"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

// AuthMiddleware creates a gin middleware for authorization
func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)

		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}

// Gets User from Payload
func (server *Server) getUserFromPayload(ctx *gin.Context) (db.User, error) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	user, err := server.store.GetUser(ctx, authPayload.Username)
	if err != nil {
		return user, errors.New("User doesn't exist for this session. That's weird.")
	}
	return user, nil
}

// Reverse Proxy logic
func proxyRequest(c *gin.Context, targetURL string, username string) {
	client := &http.Client{Timeout: 10 * time.Second}

	// Creating a new HTTP request based on the incoming Gin request
	req, err := http.NewRequest(c.Request.Method, targetURL+c.Param("proxyPath"), c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// copying headers from the original request
	for k, v := range c.Request.Header {
		req.Header[k] = v
	}

	// To add authenticated user info if required
	if username != "" {
		req.Header.Set("X-Username", username)
	}

	// Performing the request
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, errorResponse(err))
		return
	}

	// Copying the response back to the client
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
}
