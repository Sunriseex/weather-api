package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-redis/redis"
)

type Config struct {
	WeatherAPIKey   string `json:"weather_api_key"`
	RedisURL        string `json:"redisUrl"`
	CacheExpiration int    `json:"cacheExpiration"`
}

var (
	config          Config
	rdb             *redis.Client
	logger          *log.Logger
	cacheExpiration time.Duration
)

type WeatherResp struct {
	LocationName string `json:"locationName"`
	Temperature  string `json:"temperature"`
	Condition    string `json:"condition"`
}

func init() {
	logFile, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	logger = log.New(logFile, "WeatherAPI: ", log.Ldate|log.Ltime|log.Lshortfile)

	file, err := os.ReadFile("config.json")
	if err != nil {
		logger.Fatalf("Failed to read config file: %v", err)
	}
	err = json.Unmarshal(file, &config)
	if err != nil {
		logger.Fatalf("Could not parse config: %v", err)
	}
	rdb = redis.NewClient(&redis.Options{
		Addr: config.RedisURL,
	})
	logger.Println("Initializing Redis client...")
}

func main() {

	cacheExpiration = time.Duration(config.CacheExpiration) * time.Second

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/weather/{city}", getWeatherHandler)

	logger.Println("Starting server on :8080...")
	log.Println("Starting server on :8080...")

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		logger.Fatalf("Failed to start server: %v", err)
	}

}

func getWeatherHandler(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Path[len("/weather/"):]
	logger.Printf("Received request for city: %s", city)
	if city == "" {
		http.Error(w, "City parameter is required", http.StatusBadRequest)
		return
	}

	cachedData, err := rdb.Get(city).Result()
	if err == nil {
		logger.Printf("Cache miss for city: %s", city)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cachedData))
		return
	}
	weatherData, err := fetchWeatherData(city)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching weather data: %v", err), http.StatusInternalServerError)
		return
	}
	dataToCache, _ := json.Marshal(weatherData)
	rdb.Set(city, string(dataToCache), cacheExpiration)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(weatherData)

}
func fetchWeatherData(city string) (*WeatherResp, error) {
	apiURL := fmt.Sprintf(
		"https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline/%s?key=%s",
		city, config.WeatherAPIKey,
	)
	resp, err := http.Get(apiURL)
	if err != nil {
		logger.Printf("Error calling weather API: %v", err)
		return nil, fmt.Errorf("Failed to call Weather API: %v", err)

	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		logger.Printf("Weather API returned status code %d", resp.StatusCode)
		return nil, fmt.Errorf("Weather API returned status code %d", resp.StatusCode)
	}

	var apiResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		logger.Printf("Failed to decode API response: %v", err)
		return nil, fmt.Errorf("Failed to decode API response: %v", err)

	}
	weather := &WeatherResp{
		LocationName: city,
		Temperature:  fmt.Sprintf("%.2f", apiResponse["currentConditions"].(map[string]interface{})["temp"].(float64)),
		Condition:    apiResponse["currentConditions"].(map[string]interface{})["conditions"].(string),
	}
	logger.Printf("Weather data fetched for city: %s", city)
	return weather, nil
}
