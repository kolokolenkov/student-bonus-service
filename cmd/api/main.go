package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"student-bonus-service/internal/bonus"
	"student-bonus-service/internal/collector"
	"student-bonus-service/internal/database"
	"student-bonus-service/internal/student"
)

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	return value
}

func main() {
	db, err := database.NewPostgres(database.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		Name:     getEnv("DB_NAME", "student_bonus"),
	})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	defer db.Close()

	studentRepository := student.NewRepository(db)
	studentService := student.NewService(studentRepository)

	bonusRepository := bonus.NewRepository(db)

	bonusService := bonus.NewService(bonusRepository)

	calculationRepository := bonus.NewCalculationRepository(db)

	collectorConfig := collector.Config{
		BaseURL:  os.Getenv("COLLECTOR_API_URL"),
		SendPath: os.Getenv("COLLECTOR_METHOD_SEND"),
		Token:    os.Getenv("COLLECTOR_BEARER_TOKEN"),
	}

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	httpCollector := collector.NewHTTPCollector(
		collectorConfig,
		httpClient,
	)

	calculationService := bonus.NewCalculationService(
		studentService,
		bonusService,
		calculationRepository,
		httpCollector,
	)

	bonusHandler := bonus.NewHandler(calculationService)

	mux := http.NewServeMux()

	mux.HandleFunc(
		"/api/v1/bonuses/calculate",
		bonusHandler.Calculate,
	)

	mux.HandleFunc(
		"/api/v1/bonuses",
		bonusHandler.GetCalculations,
	)

	mux.HandleFunc(
		"/api/v1/bonuses/send",
		bonusHandler.SendToCollector,
	)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Println("server started on :8080")

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
