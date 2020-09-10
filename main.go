package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"demo-user/config"
	"demo-user/modules/database"
	"demo-user/modules/zookeeper"
	"demo-user/routes"
	grpcnode "demo-user/grpc/node"
)

func init() {
	config.InitENV()
	zookeeper.Connect()
	database.Connect()
}

func main() {
	envVars := config.GetEnv()
	server := echo.New()

	//CORS
	server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentLength, echo.HeaderContentType, echo.HeaderAuthorization},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodOptions},
		MaxAge:           600,
		AllowCredentials: false,
	}))

	// Middleware
	server.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} | ${remote_ip} | ${method} ${uri} - ${status} - ${latency_human}\n",
	}))
	server.Use(middleware.Recover())

	// Route
	routes.Boostrap(server)

	// Start gRPC server
	go func() {
		grpcnode.Start()
	}()

	// Start server
	server.Logger.Fatal(server.Start(envVars.AppPort))
}
