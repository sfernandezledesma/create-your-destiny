package router

import (
	"github.com/gin-gonic/gin"
	"github.com/sfernandezledesma/create-your-destiny/internal/auth"
	"github.com/sfernandezledesma/create-your-destiny/internal/game"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	r.SetTrustedProxies(nil)
	r.LoadHTMLGlob("templates/*.html")
	r.GET("/", game.RootHandler)
	r.GET("/register", auth.RegisterFormHandler)
	r.POST("/register", auth.RegisterHandler)
	r.GET("/login", auth.LoginFormHandler)
	r.POST("/login", auth.LoginHandler)
	r.GET("/play/:gameName/:pageNumber", game.PlayHandler)
	r.GET("/edit/:gameName", auth.GameOwnerMiddleware, game.EditGameHandler)
	r.NoMethod(game.BadRouteHandler)
	r.NoRoute(game.BadRouteHandler)

	return r
}
