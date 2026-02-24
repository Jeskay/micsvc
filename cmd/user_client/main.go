package main

import (
	"context"
	"log"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"

	"github.com/Jeskay/micsvc/config"
	proto "github.com/Jeskay/micsvc/protos"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	var cfg config.ClientConfig
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}
	conn, err := grpc.NewClient(cfg.Address(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to establish connection: %v", err)
	}
	defer conn.Close()
	grpcClient := proto.NewUsersClient(conn)

	generatePayload(grpcClient)
}

func generatePayload(c proto.UsersClient) {
	ctx := context.Background()
	var adminID int32 = 0
	adminPwd := "123"
	if res, err := c.Register(ctx, &proto.RegisterUser{Id: adminID, Password: adminPwd}); err != nil || res.Code != http.StatusCreated {
		log.Println("error Register: ", err)
		return
	}
	res, err := c.Login(ctx, &proto.LoginRequest{Id: adminID, Password: adminPwd})
	if err != nil || res.Status != http.StatusAccepted {
		log.Println(res)
		log.Println(err)
		return
	}
	token := res.Token
	ctx = metadata.AppendToOutgoingContext(ctx, "jwt", token)
	usersRes, err := c.GetAll(ctx, &proto.Empty{})
	if err != nil {
		log.Println(res)
		log.Println(err)
		return
	}
	log.Printf("received users: %v", usersRes.Users)
	createdRes, err := c.Post(ctx, &proto.User{Id: 1, Name: "Ivan", Surname: "Ivanov", Email: "ivan@gmail.com", Password: "321"})
	if err != nil || createdRes.Code != http.StatusCreated {
		log.Println(err)
		return
	}
}
