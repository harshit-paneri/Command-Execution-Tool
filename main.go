package main

import (
	"bytes"
	"embed"
	"io/fs"
	"net/http"
	"os/exec"
	"runtime"
	"github.com/gin-gonic/gin"
)

var frontendFS embed.FS

func main() {
	r := gin.Default()

	//Starting frontend server
	go func() {
		http.Handle("/", http.FileServer(http.Dir("./frontend")))
		http.ListenAndServe(":3000", nil)
	}()

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	frontend, _ := fs.Sub(frontendFS, "frontend")
	r.StaticFS("/", http.FS(frontend))

	r.POST("/", func(c *gin.Context) {
		var request struct {
			Command string `json:"command"`
		}

		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		cmd := exec.Command("")
		if runtime.GOOS=="windows" {
			cmd = exec.Command("cmd", "/C", request.Command)
		} else {
			cmd = exec.Command("sh", "-c", request.Command)
		}

		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Run()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":  err.Error(),
				"stdout": stdout.String(),
				"stderr": stderr.String(),
			})
			println(err.Error())
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"stdout": stdout.String(),
			"stderr": stderr.String(),
		})
	})

	println("Frontend should run on : http://localhost:3000/")

	r.Run(":31337")
}
