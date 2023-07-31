package server

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hiamthach/micro-chat/db"
	"github.com/hiamthach/micro-chat/pb"
	"github.com/hiamthach/micro-chat/util"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
)

func RunGRPCServer(config util.Config, store *mongo.Client, cache util.RedisUtil, conn *grpc.ClientConn) {
	grpcServer := grpc.NewServer()

	// socket server
	socketServer, err := db.NewSocketServer()
	if err != nil {
		log.Fatalf("Failed to create socket server: %v", err)
	}

	// Register gRPC server
	roomServer, err := NewRoomServer(config, cache, store)
	if err != nil {
		log.Fatalf("Failed to create room server: %v", err)
	}
	pb.RegisterRoomServiceServer(grpcServer, roomServer)

	chatServer, err := NewChatServer(config, cache, store, conn, *socketServer)
	if err != nil {
		log.Fatalf("Failed to create chat server: %v", err)
	}
	pb.RegisterChatServiceServer(grpcServer, chatServer)

	// Start gRPC server
	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("gRPC server listening at %v", listener.Addr())

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func RunGatewayServer(config util.Config, store *mongo.Client, cache util.RedisUtil, conn *grpc.ClientConn) {
	// initialize socket server
	socketServer, err := db.NewSocketServer()
	if err != nil {
		log.Fatalf("Failed to create socket server: %v", err)
	}

	// initialize grpc server
	roomServer, err := NewRoomServer(config, cache, store)
	if err != nil {
		log.Fatalf("Failed to create room server: %v", err)
	}

	chatServer, err := NewChatServer(config, cache, store, conn, *socketServer)
	if err != nil {
		log.Fatalf("Failed to create chat server: %v", err)
	}

	// initialize json option
	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	// initialize gRPC gateway mux
	grpcMux := runtime.NewServeMux(jsonOption)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// register gRPC server endpoint
	if err := pb.RegisterRoomServiceHandlerServer(ctx, grpcMux, roomServer); err != nil {
		log.Fatalf("Failed to register gateway: %v", err)
	}

	if err := pb.RegisterChatServiceHandlerServer(ctx, grpcMux, chatServer); err != nil {
		log.Fatalf("Failed to register gateway: %v", err)
	}

	// initialize http server
	mux := http.NewServeMux()
	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", grpcMux))
	mux.Handle("/socket.io/", socketServer)

	listener, err := net.Listen("tcp", config.ServerAddress)
	if err != nil {
		log.Fatal("Can not start server: ", err)
	}

	log.Println("Starting gateway server on", config.ServerAddress)

	err = http.Serve(listener, enableCors(mux))
	if err != nil {
		log.Fatal("Can not start server: ", err)
	}
}

// Enable CORS
func enableCors(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE, HEAD, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, X-CSRF-Token")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		h.ServeHTTP(w, r)
	})
}
