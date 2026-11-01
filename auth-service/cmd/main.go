package main

import (

	"github.com/gin-gonic/gin"
)

func main() {
    router := gin.Default()

	if err := router.Run(":8081"); err != nil {
		panic(err)
	}
}
