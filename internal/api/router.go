package api

import (
	"github.com/gin-gonic/gin"
	"github.com/sfernandezledesma/create-your-destiny/internal/api/handlers"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	r.SetTrustedProxies(nil)
	r.LoadHTMLGlob("templates/*.html")

	r.GET("/", handlers.RootHandler)

	r.GET("/register", handlers.RegisterFormHandler)
	r.POST("/register", handlers.RegisterHandler)
	r.GET("/login", handlers.LoginFormHandler)
	r.POST("/login", handlers.LoginHandler)
	r.GET("logout", handlers.LogoutHandler)

	r.GET("/play/:gameId/:sceneNumber", handlers.GamePublicMiddleware, handlers.PlayHandler)

	r.GET("/createForm", handlers.LoggedInMiddleware, handlers.CreateFormHandler)
	r.POST("/create", handlers.LoggedInMiddleware, handlers.CreateGameHandler)
	r.GET("/edit/:gameId", handlers.GameOwnerMiddleware, handlers.EditGameHandler)

	r.NoMethod(handlers.BadRouteHandler)
	r.NoRoute(handlers.BadRouteHandler)

	return r
}
