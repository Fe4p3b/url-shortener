package main

import (
	"context"
	"log"

	pb "github.com/Fe4p3b/url-shortener/internal/handlers/grpc/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial(`:3200`, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	// получаем переменную интерфейсного типа UsersClient,
	// через которую будем отправлять сообщения
	c := pb.NewShortenerClient(conn)
	getURLResp, err := c.GetURL(context.Background(), &pb.GetURLRequest{ShortUrl: "6j2EuZQ7R"})
	if err != nil {
		log.Printf("Error:%s", err)
	}

	log.Println(getURLResp)

	postURLResp, err := c.PostURL(context.Background(), &pb.PostURLRequest{OriginalUrl: "yandex.com", User: "3c93af39-fa7e-4ce9-90da-cac5c0861129"})
	if err != nil {
		log.Printf("Error:%s", err)
	}

	log.Println(postURLResp)

	getUserURLsResp, err := c.GetUserURLs(context.Background(), &pb.GetUserURLsRequest{User: "ec16638d-8346-4746-b009-a846d19f4862"})
	if err != nil {
		log.Printf("Error:%s", err)
	}

	log.Println(getUserURLsResp)

	delUserURLsResp, err := c.DelUserURLs(context.Background(), &pb.DelUserURLsRequest{User: "ec16638d-8346-4746-b009-a846d19f4862", Urls: []string{"6j2EuZQ7R", "hgbPuZw7g"}})
	if err != nil {
		log.Printf("Error:%s", err)
	}

	log.Printf("delUserURLsResp %s\n", delUserURLsResp)

	pingResp, err := c.Ping(context.Background(), &empty.Empty{})
	if err != nil {
		log.Printf("Error:%s", err)
	}

	log.Printf("pingResp %s\n", pingResp)

	getStatsResp, err := c.GetStats(context.Background(), &empty.Empty{})
	if err != nil {
		log.Printf("Error:%s", err)
	}

	log.Printf("getStatsResp %s\n", getStatsResp)
}
