package api

import (
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	r.SetTrustedProxies(nil)
	r.LoadHTMLGlob("templates/*.html")
	r.GET("/", rootHandler)
	r.GET("/register", registerFormHandler)
	r.POST("/register", registerHandler)
	r.GET("/login", loginFormHandler)
	r.POST("/login", loginHandler)
	r.GET("/play/:gameName/:pageNumber", playHandler)
	r.GET("/edit/:gameName", gameOwnerMiddleware, editGameHandler)
	r.NoMethod(badRouteHandler)
	r.NoRoute(badRouteHandler)

	return r
}
