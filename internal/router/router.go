package router

import (
	"github.com/gin-gonic/gin"
	"github.com/sfernandezledesma/create-your-destiny/internal/api"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	r.SetTrustedProxies(nil)
	r.LoadHTMLGlob("templates/*.html")
	r.GET("/", api.RootHandler)
	r.GET("/register", api.RegisterFormHandler)
	r.POST("/register", api.RegisterHandler)
	r.GET("/login", api.LoginFormHandler)
	r.POST("/login", api.LoginHandler)
	r.GET("/play/:gameName/:sceneNumber", api.PlayHandler)
	r.GET("/edit/:gameName", api.GameOwnerMiddleware, api.EditGameHandler)
	r.NoMethod(api.BadRouteHandler)
	r.NoRoute(api.BadRouteHandler)

	return r
}
