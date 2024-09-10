package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sfernandezledesma/create-your-destiny/internal/auth"
	"github.com/sfernandezledesma/create-your-destiny/internal/cache"
	"github.com/sfernandezledesma/create-your-destiny/internal/utils"
)

func LoggedInMiddleware(c *gin.Context) {
	if username := auth.GetUsernameFromContext(c); username == "" {
		c.HTML(http.StatusUnauthorized, "errorPage", "Unauthorized")
		c.Abort()
	} else {
		c.Set("username", username)
		c.Next()
	}
}

func GamePublicMiddleware(c *gin.Context) {
	gameId, _ := utils.StringToNat(c.Param("gameId")) // FIXME: Handle error
	gameData := cache.GetGameDataFromId(gameId)       // FIXME: Should check if game is nil
	if gameData.Public {
		c.Next()
	} else { // if the game is private, check if the user is the owner
		GameOwnerMiddleware(c)
	}
}

func GameOwnerMiddleware(c *gin.Context) {
	gameId, _ := utils.StringToNat(c.Param("gameId")) // FIXME: Handle error
	gameData := cache.GetGameDataFromId(gameId)       // FIXME: Should check if game is nil

	// Check if the user is logged in and retrieve username
	username := auth.GetUsernameFromContext(c)
	if username == "" {
		c.HTML(http.StatusUnauthorized, "errorPage", "Unauthorized")
		c.Abort()
		return
	}

	// Check if the user is the owner of the game
	if gameData.Author != username {
		c.HTML(http.StatusForbidden, "errorPage", "Forbidden")
		c.Abort()
		return
	}

	// If everything is fine, proceed to the next handler
	c.Next()
}
