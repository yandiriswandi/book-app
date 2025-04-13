package routers

import (
	"bioskop-app/controllers"

	"github.com/gin-gonic/gin"
)

func StartSever() *gin.Engine {
	router := gin.Default()

	router.POST("/bioskop", controllers.CreateBook)
	router.PUT("/bioskop/:id", controllers.UpdateBook)
	router.DELETE("/bioskop/:id", controllers.DeleteBook)
	router.GET("/bioskop", controllers.GetBookList)
	router.GET("/bioskop/:id", controllers.GetBookByID)

	return router
}
