package routers

import (
	"MyGramAPI/app/middleware"
	"MyGramAPI/app/services"

	_ "MyGramAPI/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func StartServer() *gin.Engine {
	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Authorization", "Content-Type"}

	router.Use(cors.New(config))
	router.Use(cors.Default())

	router.MaxMultipartMemory = 10 << 20 // 10 MB
	router.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := router.Group("/api/v1")
	{
		userRouter := v1.Group("/users")
		{
			userRouter.POST("/register", services.UserRegister)
			userRouter.POST("/login", services.UserLogin)
		}

		photoRouter := v1.Group("/photos")
		{
			photoRouter.GET("/", services.GetAllPhoto)
			photoRouter.GET("/:id", services.GetPhoto)
			photoRouter.Use(middleware.Authentication())
			photoRouter.POST("/", services.CreatePhoto)
			photoRouter.PUT("/:id", middleware.Authorization("photo"), services.UpdatePhoto)
			photoRouter.DELETE("/:id", middleware.Authorization("photo"), services.DeletePhoto)
		}

		commentRouter := v1.Group("/comments")
		{
			commentRouter.GET("/", services.GetAllComment)
			commentRouter.GET("/:id", services.GetComment)
			commentRouter.Use(middleware.Authentication())
			commentRouter.POST("/", services.CreateComment)
			commentRouter.PUT("/:id", middleware.Authorization("comment"), services.UpdateComment)
			commentRouter.DELETE("/:id", middleware.Authorization("comment"), services.DeleteComment)
		}

		socialMediaRouter := v1.Group("/social-media")
		{
			socialMediaRouter.GET("/", services.GetAllSocialMedia)
			socialMediaRouter.GET("/:id", services.GetSocialMedia)
			socialMediaRouter.Use(middleware.Authentication())
			socialMediaRouter.POST("/", services.CreateSocialMedia)
			socialMediaRouter.PUT("/:id", middleware.Authorization("socialMedia"), services.UpdateSocialMedia)
			socialMediaRouter.DELETE("/:id", middleware.Authorization("socialMedia"), services.DeleteSocialMedia)
		}
	}

	router.Run(":8082")
	return router
}
