package database

import (
	"fmt"
	"log"
	"shifty-backend/configs"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ConnectQdrant(cfg *configs.Config) *grpc.ClientConn {
	target := fmt.Sprintf("%s:%s", cfg.QdrantHost, cfg.QdrantGrpcPort)

	conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Can not connect to Qdrant: ", err)
	}
	log.Println("Connect to Qdrant successful!")
	return conn

}
