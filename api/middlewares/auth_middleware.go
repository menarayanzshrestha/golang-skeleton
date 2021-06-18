package middlewares

import (
	"net/http"
	"strings"

	"github.com/menarayanzshrestha/trello/api/responses"
	"github.com/menarayanzshrestha/trello/constants"
	"github.com/menarayanzshrestha/trello/services"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
)

// FirebaseAuthMiddleware structure
type FirebaseAuthMiddleware struct {
	service services.FirebaseService
}

// NewFirebaseAuthMiddleware creates new firebase authentication
func NewFirebaseAuthMiddleware(service services.FirebaseService) FirebaseAuthMiddleware {
	return FirebaseAuthMiddleware{
		service: service,
	}
}

// Handle handles auth requests
func (m FirebaseAuthMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := m.getTokenFromHeader(c)

		if err != nil {
			responses.ErrorJSON(c, http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}

		c.Set(constants.Claims, token.Claims)
		c.Set(constants.UID, token.UID)

		c.Next()
	}
}

// HandleAdminOnly handles middleware for admin role only
func (m FirebaseAuthMiddleware) HandleAdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := m.getTokenFromHeader(c)

		if err != nil {
			responses.ErrorJSON(c, http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}

		if !m.isAdmin(token.Claims) {
			responses.ErrorJSON(c, http.StatusUnauthorized, "un-authorized request")
			c.Abort()
			return
		}

		c.Set(constants.Claims, token.Claims)
		c.Set(constants.UID, token.UID)

		c.Next()
	}
}

// getTokenFromHeader gets token from header
func (m FirebaseAuthMiddleware) getTokenFromHeader(c *gin.Context) (*auth.Token, error) {
	header := c.GetHeader("Authorization")
	idToken := strings.TrimSpace(strings.Replace(header, "Bearer", "", 1))

	token, err := m.service.VerifyToken(idToken)
	if err != nil {
		return nil, err
	}

	return token, nil
}

// isAdmin check if cliams is admin
func (M FirebaseAuthMiddleware) isAdmin(claims map[string]interface{}) bool {

	role := claims["role"]
	isAdmin := false
	if role != nil {
		isAdmin = role.(string) == "admin"
	}

	return isAdmin

}
