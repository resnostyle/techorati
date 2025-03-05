package main

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

func main() {
    r := gin.Default()

    r.GET("/ping", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"message": "pong"})
    })

    if err := r.Run(":8081"); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}