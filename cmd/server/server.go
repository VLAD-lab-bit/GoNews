package main

import (
	"GoNews/pkg/api"
	"GoNews/pkg/storage"
	"GoNews/pkg/storage/memdb"
	"GoNews/pkg/storage/mongo"
	"GoNews/pkg/storage/postgres"
	"log"
	"net/http"
)

type server struct {
	db  storage.Interface
	api *api.API
}

func main() {
	var srv server

	storageType := "postgres" // Измените на "postgres" или "mongo" для работы с соответствующими БД

	switch storageType {
	case "postgres":
		connStr := "user=postgres password=vlad5043 dbname=gonnews sslmode=disable"
		db, err := postgres.New(connStr)
		if err != nil {
			log.Fatalf("Ошибка подключения к PostgreSQL: %v", err)
		}
		srv.db = db

	case "mongo":
		mongoURI := "mongodb://localhost:27017"
		dbName := "gonnews" // Имя базы данных
		db, err := mongo.New(mongoURI, dbName)
		if err != nil {
			log.Fatalf("Ошибка подключения к MongoDB: %v", err)
		}
		srv.db = db

	default:
		log.Println("Используется хранилище в памяти (memdb)")
		srv.db = memdb.New()
	}

	srv.api = api.New(srv.db)

	http.ListenAndServe(":8080", srv.api.Router())
}
