package routers

import (
	"bioskop-app/controllers"

	"github.com/gin-gonic/gin"
)

func StartSever() *gin.Engine {
	router := gin.Default()
	//book
	router.POST("/book", controllers.CreateBook)
	router.PUT("/book/:id", controllers.UpdateBook)
	router.DELETE("/book/:id", controllers.DeleteBook)
	router.GET("/book", controllers.GetBookList)
	router.GET("/book/:id", controllers.GetBookByID)
	//category
	router.POST("/category", controllers.CreateCategory)
	router.PUT("/category/:id", controllers.UpdateCategory)
	router.DELETE("/category/:id", controllers.DeleteCategory)
	router.GET("/category", controllers.GetCategoryList)
	router.GET("/category/:id", controllers.GetCategoryByID)

	return router
}
