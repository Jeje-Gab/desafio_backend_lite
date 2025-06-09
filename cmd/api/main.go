package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"

	negdelivery "desafio_backend/internal/negociacoes/delivery/http"
	negrepo "desafio_backend/internal/negociacoes/repository"
	negusecase "desafio_backend/internal/negociacoes/usecase"
)

func mustGetenv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("variável de ambiente %s não definida", key)
	}
	return v
}

func main() {
	// carrega .env (se não achar, prossegue usando variáveis do ambiente)
	if err := godotenv.Load(); err != nil {
		log.Printf("⚠️  .env não encontrado, usando variáveis do ambiente")
	}

	// agora exigimos todas as variáveis — sem fallback
	dbHost := mustGetenv("DB_HOST")
	dbPort := mustGetenv("DB_PORT")
	dbUser := mustGetenv("DB_USER")
	dbPass := mustGetenv("DB_PASS")
	dbName := mustGetenv("DB_NAME")
	apiPort := mustGetenv("API_PORT")
	frontPort := mustGetenv("FRONT_PORT")

	// conexão com PostgreSQL
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName,
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco: %v", err)
	}
	defer db.Close()
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	// repository, service e handlers
	repo := negrepo.NewRepository(db)
	svc := negusecase.NewService(repo)

	// servidor da API
	api := echo.New()
	api.Use(middleware.Logger(), middleware.Recover())
	api.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{fmt.Sprintf("http://localhost:%s", frontPort)},
		AllowMethods: []string{http.MethodGet},
	}))

	api.GET("/api/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "up"})
	})
	negGroup := api.Group("/api/negociacoes")
	negdelivery.RegisterRoutes(negGroup, svc)

	log.Printf("🚀 API Health   -> http://localhost:%s/api/health", apiPort)
	log.Printf("🚀 API Stats    -> http://localhost:%s/api/negociacoes/stats?ticker=WINM25", apiPort)

	// servidor do front
	front := echo.New()
	front.Use(middleware.Logger(), middleware.Recover())
	front.GET("/", func(c echo.Context) error {
		return c.File("public/index.html")
	})
	front.Static("/static", "public")

	if err := api.Start(":" + apiPort); err != nil {
		log.Fatalf("API server erro: %v", err)
	}
}
