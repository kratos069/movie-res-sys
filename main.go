package main

import (
	"context"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kratos69/movie-app/api"
	db "github.com/kratos69/movie-app/db/sqlc"
	"github.com/kratos69/movie-app/util"
)

func main() {
	// load ENV variables
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatalln("cannot load config file:", err)
	}

	// conn to database
	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatalln("cannot connect to the db:", err)
	}

	// run db migrations
	runDBMigration(config.MigrationURL, config.DBSource)

	store := db.NewStore(connPool)

	runGinServer(config, store)
}

func runDBMigration(migrationURL, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatalln("cannot create new migration instance:", err)
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalln("failed to run migrate up:", err)
	}

	fmt.Println("db migrated successfully")
}

// run Gin Server for HTTP requests
func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatalln("cannot create server:", err)
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatalln("cannot start server:", err)
	}
}
