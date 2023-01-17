package main

import (
	"context"
	"database/sql"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hmoodallahma/123bank/api"
	db "github.com/hmoodallahma/123bank/db/sqlc"
	"github.com/hmoodallahma/123bank/gapi"
	"github.com/hmoodallahma/123bank/pb"
	"github.com/hmoodallahma/123bank/util"
	_ "github.com/lib/pq"
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

	store := db.NewStore(conn)
	go runGatewayServer(config, store)
	runGrpcServer(config, store)

}
func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server")
	}
	err = server.Start(config.HttpServerAddress)
	if err != nil {
		log.Fatal("cannot start server", err)
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

	fs := http.FileServer(http.Dir("./doc/swagger"))
	http.Handle("/swagger/", http.StripPrefix("/swagger", fs))

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
