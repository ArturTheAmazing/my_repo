package main

import (
	"Projects/SberAuto/service/shared"
	pb "Projects/SberAuto/service/validator"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
)



func main() {
	listen, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("fail to listen: %v", err)
	}

	httpServer := http.Server{
		Handler: promhttp.HandlerFor(shared.Reg, promhttp.HandlerOpts{}),
		Addr: "localhost:8081",
	}

	var opts []grpc.ServerOption

	server := grpc.NewServer(opts...)
	pb.RegisterValidatorServer(server, shared.NewServer())

	go func(){
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatalf("Error on http server: %s", err)
		}
	}()

	server.Serve(listen)
}
