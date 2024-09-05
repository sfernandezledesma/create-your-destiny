package handlers

import (
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/sfernandezledesma/create-your-destiny/internal/auth"
	"github.com/sfernandezledesma/create-your-destiny/internal/game"
)

func GameOwnerMiddleware(c *gin.Context) {
	gameName := c.Param("gameName")

	// Check if the user is logged in and retrieve username
	username := auth.GetUsernameFromContext(c)

	if username == "" {
		c.HTML(http.StatusUnauthorized, "errorPage", "Unauthorized")
		c.Abort()
		return
	}

	// Check if the user is the owner of the game, gameName is unique
	if !slices.Contains(game.GamesByUser[username], gameName) {
		c.HTML(http.StatusForbidden, "errorPage", "Forbidden")
		c.Abort()
		return
	}

	// If everything is fine, proceed to the next handler
	c.Next()
}
