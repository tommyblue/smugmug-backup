package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

var (
	sockFlag = flag.Bool("sock", false, "use unix socket")
	helpFlag = flag.Bool("help", false, "show help")
)

func main() {
	flag.Parse()

	if *helpFlag {
		flag.PrintDefaults()
		os.Exit(0)
	}

	var ln net.Listener
	if *sockFlag {
		var err error
		ln, err = useSock()
		if err != nil {
			log.Fatalf("failed to create socket: %v", err)
		}
	} else {
		var err error
		ln, err = net.Listen("tcp", "127.0.0.1:")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		fmt.Printf("{\"listen\": \"%s\"}", ln.Addr().String())
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Exit(1)
	}()

	srv := grpc.NewServer()
	grpc_health_v1.RegisterHealthServer(srv, health.NewServer())
	reflection.Register(srv)

	log.Fatal(srv.Serve(ln))

	// mux := http.NewServeMux()
	// mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Header().Add("Content-Type", "application/json")
	// 	ret := `{"message": "Hello, World!"}`
	// 	fmt.Fprintf(w, ret)
	// })

	// log.Info("Server started on port 8088")
	// log.Fatal(http.ListenAndServe("0.0.0.0:8088", mux))
}

func useSock() (net.Listener, error) {
	t, err := os.MkdirTemp("", "grpc")
	if err != nil {
		log.Fatalf("failed to create temp dir: %v", err)
	}
	socket := filepath.Join(t, "grpc.sock")
	if err != nil {
		log.Fatalf("failed to create socket: %v", err)
	}
	defer os.Remove(socket)

	log.Printf("grpc ran on unix socket protocol %s", socket)
	return net.Listen("unix", socket)
}
