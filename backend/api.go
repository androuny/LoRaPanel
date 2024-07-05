package backend

import (
	"log"

	"github.com/gin-gonic/gin"
)

type Api struct {
	router *gin.Engine
	mhdl *MongoHandler
}

type CreateSensorRouteBody struct {
	Name string `json:"name" binding:"required"`
}

type GetAllSensorDataBody struct {
	Name string `json:"name" binding:"required"`
}

type CreateSensorDataRouteBody struct {
	Temperature float64 `json:"temperature" binding:"required"`
	Humidity float64 `json:"humidity" binding:"required"`
	SensorName string `json:"sensorName" binding:"required"`
}

func (api Api) PingPongRoute(c *gin.Context) {
	c.JSON(200, gin.H {
		"message" : "pong",
	})
}

func (api Api) CreateSensorRoute(c *gin.Context) {
	var reqBody CreateSensorRouteBody
	err := c.BindJSON(&reqBody)

	if err != nil {
		log.Fatalf("Create Sensor API error - %v", err)
	}
	
	sensorName := reqBody.Name
	_, err1 := api.mhdl.NewSensor(sensorName)
	if err1 != nil {
		c.JSON(400, gin.H{"message": err1.Error()})
	} else {
		c.JSON(200, gin.H{"message": "new sensor created successfully"})
	}

}

func (api Api) CreateSensorDataRoute(c *gin.Context) {
	var reqBody CreateSensorDataRouteBody
	err := c.BindJSON(&reqBody)

	if err != nil {
		log.Fatalf("Create SensorData API error - %v", err)
	}
	
	sensorName := reqBody.SensorName
	Temperature := reqBody.Temperature
	Humidity := reqBody.Humidity

	_, err1 := api.mhdl.NewSensorData(sensorName, Temperature, Humidity)
	if err1 != nil {
		c.JSON(400, gin.H{"message": err1.Error()})
	} else {
		c.JSON(200, gin.H{"message": "new sensor data created successfully"})
	}
}

func (api Api) GetAllSensorDataRoute(c *gin.Context) {
	var reqBody GetAllSensorDataBody
	err := c.BindJSON(&reqBody)

	if err != nil {
		log.Fatalf("GetAllSensorData API error - %v", err.Error())
	}

	sensorName := reqBody.Name
	data, _, err1 := api.mhdl.GetAllSensorData(sensorName)
	if err1 != nil {
		c.JSON(400, gin.H{"message" : err.Error()})
	}

	c.JSON(200, gin.H{"results" : data})
}

func InitBackend(db_uri string, database string) (Api) {
	r := gin.Default()
	log.Printf("Connecting to mongodb database...")
	mhdl, _, err := NewMongoHandler(db_uri, database)
	if err != nil {
		log.Fatalf("database error %v", err)
	}
	log.Printf("Connected to mongodb database successfully!")

	api := Api{r, &mhdl}
	api.router.GET("/ping", api.PingPongRoute)
	api.router.POST("/api/createSensor", api.CreateSensorRoute)
	api.router.POST("/api/createSensorData", api.CreateSensorDataRoute)
	api.router.POST("/api/getAllSensorData", api.GetAllSensorDataRoute)

	return api
}

func (api Api) StartBackend() (error) {
	err := api.router.Run()
	return err
}