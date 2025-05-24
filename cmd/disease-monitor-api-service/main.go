package main

import (
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/smoleniak/disease-monitor-webapi/api"
	"github.com/smoleniak/disease-monitor-webapi/internal/disease_monitor"
)

func main() {
	log.Printf("Server started")
	port := os.Getenv("DISEASE_MONITOR_API_PORT")
	if port == "" {
		port = "8080"
	}
	environment := os.Getenv("DISEASE_MONITOR_API_ENVIRONMENT")
	if !strings.EqualFold(environment, "production") { // case insensitive comparison
		gin.SetMode(gin.DebugMode)
	}
	engine := gin.New()
	engine.Use(gin.Recovery())
	// request routings
	handleFunctions := &disease_monitor.ApiHandleFunctions{
		DiseaseTypesAPI:        disease_monitor.NewDiseaseTypesApi(),
		DiseaseMonitorCasesAPI: disease_monitor.NewDiseaseMonitorCasesApi(),
	}
	disease_monitor.NewRouterWithGinEngine(engine, *handleFunctions)
	engine.GET("/openapi", api.HandleOpenApi)
	engine.Run(":" + port)
}
