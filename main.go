package main

import (
	"bioskop-app/config"
	"bioskop-app/routers"
	"fmt"
)

func main() {
	config.InitDB()

	PORT := ":8080"

	routers.StartSever().Run(PORT)

	fmt.Println("Server running at http://localhost:8080")

}
