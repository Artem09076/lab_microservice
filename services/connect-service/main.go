package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"s0709-22/internal/proxyproto"
	"s0709-22/services/connect-service/internal/service"

	"google.golang.org/grpc"
)

const (
	// ConnString ...
	ConnString = "postgres://appuser:apppass@127.0.0.1:5432/userdb?sslmode=disable"
)

func main() {
	listener, err := net.Listen("tcp4", "127.0.0.1:10000")
	if err != nil {
		log.Fatalln(err)
	}

	errChan := make(chan error)

	srv := grpc.NewServer()

	svc, err := service.New(ConnString)

	proxyproto.RegisterCentrifugoProxyServer(srv, svc)

	exitCtx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}

		cancel()

		srv.GracefulStop()

		close(errChan)

		if err := listener.Close(); err != nil {
			log.Println(err)
		}
	}()

	go func() {
		errChan <- srv.Serve(listener)
	}()

	select {
	case err := <-errChan:
		log.Fatalln(err)
	case <-exitCtx.Done():
		log.Println("exit")
	}

}
