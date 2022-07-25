package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {

	err := ConnectDatabase()
	checkErr(err)

	r := gin.Default()

	// API v1
	v1 := r.Group("/api/v1")
	{
		v1.GET("measurement", getMeasurements)
		v1.GET("measurement/:id", getMeasurementById)
		v1.POST("measurement", addMeasurement)
		v1.PUT("measurement/:id", updateMeasurement)
		v1.DELETE("measurement/:id", deleteMeasurement)
		v1.OPTIONS("measurement", options)
	}

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	r.Run()
}

func getMeasurements(c *gin.Context) {

	measurements, err := GetMeasurements(10)

	checkErr(err)

	if measurements == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No Records Found"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"data": measurements})
	}
}

func getMeasurementById(c *gin.Context) {

	// grab the Id of the record want to retrieve
	id := c.Param("id")

	measurement, err := GetMeasurementById(id)

	checkErr(err)
	// if the timestamp is blank we can assume nothing is found
	if measurement.Timestamp == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No Records Found"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"data": measurement})
	}
}

func addMeasurement(c *gin.Context) {

	var json Measurement

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	success, err := AddMeasurement(json)

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
}

func updateMeasurement(c *gin.Context) {

	var json Measurement

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	measurementId, err := strconv.Atoi(c.Param("id"))

	fmt.Printf("Updating id %d", measurementId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
	}

	success, err := UpdateMeasurement(json, measurementId)

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
}

func deleteMeasurement(c *gin.Context) {

	measurementId, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
	}

	success, err := DeleteMeasurement(measurementId)

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
}

func options(c *gin.Context) {

	ourOptions := "HTTP/1.1 200 OK\n" +
		"Allow: GET,POST,PUT,DELETE,OPTIONS\n" +
		"Access-Control-Allow-Origin: http://locahost:8080\n" +
		"Access-Control-Allow-Methods: GET,POST,PUT,DELETE,OPTIONS\n" +
		"Access-Control-Allow-Headers: Content-Type\n"

	c.String(200, ourOptions)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
