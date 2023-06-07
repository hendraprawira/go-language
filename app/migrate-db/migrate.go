package main

import (
	"fmt"
	"log"
	"os"

	"Remember-Golang/app/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DBs *gorm.DB

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(err)
	}
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PSWD")
	dbname := os.Getenv("DB_NAME")
	psqlInfo := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, port, dbname)
	DB, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})

	if err != nil {
		log.Fatalln(err, "ERR")
	}

	DBs = DB
}

func main() {
	// run syntax " go run migrate/migrate.go" for migrate schema to db
	DBs.AutoMigrate(&models.MasterUser{})
	fmt.Println("? Migration complete")
}
