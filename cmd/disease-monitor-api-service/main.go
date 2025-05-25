package main

import (
	"log"
	"os"
	"strings"

	"context"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/smoleniak/disease-monitor-webapi/api"
	"github.com/smoleniak/disease-monitor-webapi/internal/db_service"
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
	corsMiddleware := cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "PUT", "POST", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{""},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	})
	engine.Use(corsMiddleware)

	// setup context update  middleware
	dbService := db_service.NewMongoService[disease_monitor.Region](db_service.MongoServiceConfig{})
	defer dbService.Disconnect(context.Background())
	engine.Use(func(ctx *gin.Context) {
		ctx.Set("db_service", dbService)
		ctx.Next()
	})

	// request routings
	handleFunctions := &disease_monitor.ApiHandleFunctions{
		DiseaseTypesAPI:        disease_monitor.NewDiseaseTypesApi(),
		DiseaseMonitorCasesAPI: disease_monitor.NewDiseaseMonitorCasesApi(),
		RegionsAPI:             disease_monitor.NewRegionsApi(),
	}
	disease_monitor.NewRouterWithGinEngine(engine, *handleFunctions)
	engine.GET("/openapi", api.HandleOpenApi)
	engine.Run(":" + port)
}
