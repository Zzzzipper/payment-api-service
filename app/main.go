package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"

	pbPayment "payment-api-service/proto"
	"payment-api-service/server"

	"payment-api-service/gateway"

	"payment-api-service/insecure"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
)

var (
	username = flag.String("username", "user", "The username to authenticate with")
	password = flag.String("password", "1234", "The password to authenticate with")
	token    = flag.String("token", "testtoken", "The token to authenticate with")
)

// version: 0.6.0
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

	auther := server.Authenticator{
		Username: *username,
		Password: *password,
		Token:    *token,
	}

	var grpcServer *grpc.Server
	var protocol string

	if strings.ToLower(os.Getenv("SERVE_HTTP")) == "true" {
		grpcServer = grpc.NewServer(
			grpc.StreamInterceptor(grpc_auth.StreamServerInterceptor(auther.Authenticate)),
			grpc.UnaryInterceptor(grpc_auth.UnaryServerInterceptor(auther.Authenticate)),
		)
		protocol = "http"
	} else {
		grpcServer = grpc.NewServer(
			// TODO: Replace with your own certificate!
			grpc.Creds(credentials.NewServerTLSFromCert(&insecure.Cert)),
			grpc.StreamInterceptor(grpc_auth.StreamServerInterceptor(auther.Authenticate)),
			grpc.UnaryInterceptor(grpc_auth.UnaryServerInterceptor(auther.Authenticate)),
		)
		protocol = "https"
	}

	pbPayment.RegisterPaymentServiceServer(grpcServer, server.New())

	log.Info(fmt.Sprintf("Serving gRPC on %s://%d", protocol, grpcAddr))

	// Serve gRPC Server
	go func() {
		log.Fatal(grpcServer.Serve(listener))
	}()

	err = gateway.Run("dns:///" + grpcAddr)

	log.Fatalln(err)

}
