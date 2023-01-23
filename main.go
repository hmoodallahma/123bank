package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hmoodallahma/123bank/api"
	db "github.com/hmoodallahma/123bank/db/sqlc"
	_ "github.com/hmoodallahma/123bank/doc/statik"
	"github.com/hmoodallahma/123bank/gapi"
	"github.com/hmoodallahma/123bank/pb"
	"github.com/hmoodallahma/123bank/util"
	_ "github.com/lib/pq"
	"github.com/rakyll/statik/fs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
	"log"
	"net"
	"net/http"
)

func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load configurations", err)
	}

	// connect to db and create a store
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	// run db migrgtion
	runDBMigration(config.MigrationURL, config.DBSource)
	store := db.NewStore(conn)
	go runGatewayServer(config, store)
	runGrpcServer(config, store)
	// runGinServer(config, store)
}

func runDBMigration(migrationUrl string, dbSource string) {
	migration, err := migrate.New(migrationUrl, dbSource)
	if err != nil {
		log.Fatal("cannot create a new migration instance ", err)
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("failed to run migrate up ", err)
	}

	fmt.Println("migration ran successfully")
}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server ")
	}
	err = server.Start(config.HttpServerAddress)
	if err != nil {
		log.Fatal("cannot start server ", err)
	}
}
func runGrpcServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)

	grpcServer := grpc.NewServer()
	pb.RegisterBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GrpcServerAddress)

	if err != nil {
		log.Fatal("cannot create listener", err)
	}

	log.Printf("starting GRPC server on %d", config.GrpcServerAddress)

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start grpc server", err)
	}
}

func runGatewayServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	grpcMux := runtime.NewServeMux(jsonOption)
	err = pb.RegisterBankHandlerServer(ctx, grpcMux, server)

	if err != nil {
		log.Fatal("cannot register handler server", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	statikFS, err := fs.New()

	if err != nil {
		log.Fatal("cannot create statik fs", err)
	}

	swaggerHandler := http.StripPrefix("/swagger", http.FileServer(statikFS))
	mux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", config.HttpServerAddress)

	if err != nil {
		log.Fatal("cannot create listener", err)
	}

	log.Printf("starting HTTP server on %d", config.HttpServerAddress)

	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal("cannot start http server", err)
	}
}
