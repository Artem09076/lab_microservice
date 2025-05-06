package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/Artem09076/lab_microservice.git/internal/proxyproto"
	"github.com/Artem09076/lab_microservice.git/services/connect-service/internal/config"
	"github.com/Artem09076/lab_microservice.git/services/connect-service/internal/service"
	"google.golang.org/grpc"
)

func serve() error {
	conf, err := config.Load()
	if err != nil {
		return err
	}
	listener, err := net.Listen("tcp4", ":"+conf.Port)
	if err != nil {
		return err
	}

	errChan := make(chan error)

	srv := grpc.NewServer()

	svc, err := service.New(conf)
	if err != nil {
		return err
	}

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
		return err
	case <-exitCtx.Done():
		log.Println("exit")
	}
	return nil
}

func main() {
	log.Println("serve start")
	if err := serve(); err != nil {
		log.Fatalln(err)
	}
}
