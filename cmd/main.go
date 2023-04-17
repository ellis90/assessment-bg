package main

import (
	"fmt"
	"github.com/ellis90/assessment-bg/router"
	"github.com/ellis90/assessment-bg/service"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
)

// init gets called before the main function
func init() {
	if err := godotenv.Load(); err != nil {
		log.Error("No .env file found create")
		os.Exit(1)
	}
}

func main() {

	host := os.Getenv("HOST")
	dbUSER := os.Getenv("DB_USERNAME")
	pass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	src := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUSER, pass, host, dbPort, dbName)
	log.Println(src)

	os, err := service.NewCustomerServices(
		service.WithPGXConfiguration(log.New(), src),
	)
	if err != nil {
		log.Fatal("failed to create service")
	}

	e := router.Router(os)
	if err := e.Start(":9090"); err != nil {
		log.Fatal("failed to start up server")
	}

}
