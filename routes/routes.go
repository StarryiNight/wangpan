package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wangpan/controllers"
	"wangpan/middlewares"
)

func Setup() *gin.Engine {
	r := gin.New()



	//注册业务路由
	r.POST("/signup",controllers.SignUpHandler)
	//登陆业务路由
	r.POST("/login",controllers.LoginHandler)

	r.GET("/", middlewares.JWTAuthMiddleware(),func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
	r.GET("/logout",middlewares.JWTAuthMiddleware(),controllers.Logout)
	//上传业务路由
	r.POST("/upload",middlewares.JWTAuthMiddleware(),controllers.UploadHandler)
	//下载业务路由
	r.GET("/download",middlewares.JWTAuthMiddleware(),controllers.DownloadHandler)
	return r
}
