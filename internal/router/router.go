package router

import (
	"net/http"
	"runtime"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/fudio101/fube/docs"
	"github.com/fudio101/fube/internal/handlers"
)

// SetupRouter initializes the router and sets up all routes
func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(CORSMiddleware())
	r.Use(RequestIDMiddleware())
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	setupV1Routes(r)
	setupSwagger(r)
	setupStaticFiles(r)

	return r
}

// setupV1Routes configures all v1 API routes
func setupV1Routes(r *gin.Engine) {
	v1 := r.Group("/v1")
	{
		/*** START USER ***/
		user := new(handlers.UserHandler)

		v1.POST("/user/login", user.Login)
		v1.POST("/user/register", user.Register)
		v1.GET("/user/logout", user.Logout)

		/*** START AUTH ***/
		auth := new(handlers.AuthHandler)

		//Refresh the token when needed to generate new access_token and refresh_token for the user
		v1.POST("/token/refresh", auth.Refresh)

		/*** START Article ***/
		article := new(handlers.ArticleHandler)

		v1.POST("/article", TokenAuthMiddleware(), article.Create)
		v1.GET("/articles", TokenAuthMiddleware(), article.All)
		v1.GET("/article/:id", TokenAuthMiddleware(), article.One)
		v1.PUT("/article/:id", TokenAuthMiddleware(), article.Update)
		v1.DELETE("/article/:id", TokenAuthMiddleware(), article.Delete)
	}
}

// setupSwagger configures the Swagger documentation route
func setupSwagger(r *gin.Engine) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}

// setupStaticFiles sets up static file serving and HTML routes
func setupStaticFiles(r *gin.Engine) {
	r.LoadHTMLGlob("./public/html/*")
	r.Static("/public", "./public")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"ginBoilerplateVersion": "v0.03",
			"goVersion":             runtime.Version(),
		})
	})

	r.NoRoute(func(c *gin.Context) {
		c.HTML(404, "404.html", gin.H{})
	})
}
