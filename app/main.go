package main

import (
	"io/ioutil"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"

	pbPayment "payment-api-service/proto"
	"payment-api-service/server"

	"payment-api-service/gateway"

	"payment-api-service/insecure"
)

// version: 0.3.0
func main() {

	// Adds gRPC internal logs. This is quite verbose, so adjust as desired!
	log := grpclog.NewLoggerV2(os.Stdout, ioutil.Discard, ioutil.Discard)
	grpclog.SetLoggerV2(log)

	grpcPort := os.Getenv("API_GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "10000"
	}

	grpcAddr := "0.0.0.0:" + grpcPort

	// Initialize listener
	listener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}
	defer listener.Close()

	grpcServer := grpc.NewServer(
		// TODO: Replace with your own certificate!
		grpc.Creds(credentials.NewServerTLSFromCert(&insecure.Cert)),
	)
	pbPayment.RegisterPaymentServiceServer(grpcServer, server.New())

	log.Info("Serving gRPC on https://", grpcAddr)

	// Serve gRPC Server
	go func() {
		log.Fatal(grpcServer.Serve(listener))
	}()

	err = gateway.Run("dns:///" + grpcAddr)

	log.Fatalln(err)

}
