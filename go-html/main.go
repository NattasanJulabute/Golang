package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type JsonData struct {
	Hourly struct {
		Time           []string
		Temperature_2m []float32
	}
}

func fetchTemperature() (JsonData, error) {
	if err := godotenv.Load("local.env"); err != nil {
		log.Panicln("no env file")
	}

	fmt.Println(os.Getenv("URL"))
	url := os.Getenv("URL")
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	var jsonData JsonData
	if err := json.NewDecoder(resp.Body).Decode(&jsonData); err != nil {
		panic(err)
	}

	for i, v := range jsonData.Hourly.Time {
		s := strings.Replace(v, "T", " ", 1)
		jsonData.Hourly.Time[i] = s
	}

	return jsonData, nil
}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	r.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Main website",
		})
	})

	r.GET("/", func(c *gin.Context) {

		data, _ := fetchTemperature()

		// labels := data.Hourly.Time
		// temps := data.Hourly.Temperature_2m
		labelBytes, _ := json.Marshal(data.Hourly.Time)
		tempBytes, _ := json.Marshal(data.Hourly.Temperature_2m)

		c.HTML(http.StatusOK, "index.html", gin.H{
			"Labels": template.JS(labelBytes),
			"Temps":  template.JS(tempBytes),
		})
	})

	r.Run(":8080")
}
