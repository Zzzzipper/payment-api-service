package gateway

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/fs"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"

	pbPayment "payment-api-service/proto"

	"payment-api-service/third_party"

	"payment-api-service/insecure"
)

// getOpenAPIHandler serves an OpenAPI UI.
func getOpenAPIHandler() http.Handler {
	mime.AddExtensionType(".svg", "image/svg+xml")
	// Use subdirectory in embedded files
	subFS, err := fs.Sub(third_party.OpenAPI, "OpenAPI")
	if err != nil {
		panic("couldn't create sub filesystem: " + err.Error())
	}
	return http.FileServer(http.FS(subFS))
}

// getTestAPIHandler serves an test api UI.
func getTestAPIHandler() http.Handler {
	mime.AddExtensionType(".svg", "image/svg+xml")
	// Use subdirectory in embedded files
	subFS, err := fs.Sub(third_party.TestAPI, "TestAPI")
	if err != nil {
		panic("couldn't create sub filesystem: " + err.Error())
	}
	return http.FileServer(http.FS(subFS))
}

// Run runs the gRPC-Gateway, dialling the provided address.
func Run(dialAddr string) error {
	// Adds gRPC internal logs. This is quite verbose, so adjust as desired!
	log := grpclog.NewLoggerV2(os.Stdout, ioutil.Discard, ioutil.Discard)
	grpclog.SetLoggerV2(log)

	// Create a client connection to the gRPC Server we just started.
	// This is where the gRPC-Gateway proxies the requests.

	var clientConn *grpc.ClientConn
	var err error

	if strings.ToLower(os.Getenv("SERVE_HTTP")) == "true" {
		clientConn, err = grpc.DialContext(
			context.Background(),
			dialAddr,
			grpc.WithInsecure(),
			grpc.WithBlock(),
		)
	} else {
		clientConn, err = grpc.DialContext(
			context.Background(),
			dialAddr,
			grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(insecure.CertPool, "")),
			grpc.WithBlock(),
		)
	}
	if err != nil {
		return fmt.Errorf("failed to dial server: %w", err)
	}

	gwmux := runtime.NewServeMux()
	err = pbPayment.RegisterPaymentServiceHandler(context.Background(), gwmux, clientConn)
	if err != nil {
		return fmt.Errorf("failed to register gateway: %w", err)
	}

	oa := getOpenAPIHandler()
	// ta := getTestAPIHandler()

	httpPort := os.Getenv("SWAGGER_PORT")
	if httpPort == "" {
		httpPort = "11000"
	}
	webAddr := "0.0.0.0:" + httpPort

	gwServer := &http.Server{
		Addr: webAddr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// if strings.HasPrefix(r.URL.Path, "/api") {
			// 	gwmux.ServeHTTP(w, r)
			// 	return
			// }
			// if strings.HasPrefix(r.URL.Path, "/api") {
			// 	fmt.Println(r.URL.Path)
			// 	ta.ServeHTTP(w, r)
			// 	return
			// }
			oa.ServeHTTP(w, r)
		}),
	}
	// Empty parameters mean use the TLS Config specified with the server.
	if strings.ToLower(os.Getenv("SERVE_HTTP")) == "true" {
		log.Info("Serving gRPC-Gateway and OpenAPI Documentation on http://", webAddr)
		return fmt.Errorf("serving gRPC-Gateway server: %w", gwServer.ListenAndServe())
	}

	gwServer.TLSConfig = &tls.Config{
		Certificates: []tls.Certificate{insecure.Cert},
	}
	log.Info("Serving gRPC-Gateway and OpenAPI Documentation on https://", webAddr)
	return fmt.Errorf("serving gRPC-Gateway server: %w", gwServer.ListenAndServeTLS("", ""))
}
