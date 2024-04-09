package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"ticketing-api/api"
	"ticketing-api/data"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gocql/gocql"
	"github.com/joho/godotenv"
)

func main() {
	if _, err := os.Stat(".env"); err == nil {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	postgres, err := sql.Open("postgres", os.Getenv("POSTGRES_DSN"))
	if err != nil {
		log.Fatal("failed to open postgres db connection:", err)
	}

	cluster := gocql.NewCluster(strings.Split(os.Getenv("SCYLLA_HOSTS"), ",")...)
	// clusterPort, err := strconv.Atoi(os.Getenv("SCYLLA_PORT"))
	// if err != nil {
	// 	log.Fatal("failed to parse scylla port:", err)
	// }
	// cluster.Port = clusterPort
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: os.Getenv("SCYLLA_USERNAME"),
		Password: os.Getenv("SCYLLA_PASSWORD"),
	}
	cluster.Keyspace = os.Getenv("SCYLLA_KEYSPACE")

	scylla, err := cluster.CreateSession()
	if err != nil {
		log.Fatal("failed to open scylla db connection:", err)
	}

	dataAdapter := data.CreateDataAdapter(
		data.CreateAccountAdapter(postgres),
		data.CreateTicketAdapter(postgres),
		data.CreateMessageAdapter(scylla),
	)

	server := api.CreateAPIServer(fmt.Sprintf(":%s", os.Getenv("PORT")), dataAdapter)
	log.Fatal(server.Start())
}
