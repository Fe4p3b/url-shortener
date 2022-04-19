package grpc

import (
	"context"

	"github.com/Fe4p3b/url-shortener/internal/handlers"
	pb "github.com/Fe4p3b/url-shortener/internal/handlers/grpc/proto"
	"github.com/golang/protobuf/ptypes/empty"
)

type ShortenerServer struct {
	pb.UnimplementedShortenerServer
	h handlers.Handlers
}

func NewShortenerServer(h handlers.Handlers) *ShortenerServer {
	return &ShortenerServer{
		h: h,
	}
}

func (s *ShortenerServer) GetURL(ctx context.Context, in *pb.GetURLRequest) (*pb.GetURLResponse, error) {
	var response pb.GetURLResponse

	u, err := s.h.GetURL(in.ShortUrl)
	if err != nil {
		response.Error = err.Error()
		return &response, err
	}
	response.OriginalUrl = u.URL

	return &response, nil
}

func (s *ShortenerServer) PostURL(ctx context.Context, in *pb.PostURLRequest) (*pb.PostURLResponse, error) {
	var response pb.PostURLResponse

	u, err := s.h.PostURL(in.OriginalUrl, in.User)
	if err != nil {
		response.Error = err.Error()
		return &response, err
	}
	response.ShortUrl = u
	return &response, nil
}

func (s *ShortenerServer) GetUserURLs(ctx context.Context, in *pb.GetUserURLsRequest) (*pb.GetUserURLsResponse, error) {
	var response pb.GetUserURLsResponse

	u, err := s.h.GetUserURLs(in.User)
	if err != nil {
		return &response, err
	}

	for _, v := range u {
		response.Urls = append(response.Urls, &pb.URL{CorrelationId: v.CorrelationID, Url: v.URL, ShortUrl: v.ShortURL, UserId: v.UserID, IsDeleted: v.IsDeleted})
	}
	return &response, nil
}

func (s *ShortenerServer) DelUserURLs(ctx context.Context, in *pb.DelUserURLsRequest) (*pb.DelUserURLsResponse, error) {
	var response pb.DelUserURLsResponse

	s.h.DeleteUserURLs(in.User, in.Urls)

	return &response, nil
}

func (s *ShortenerServer) Ping(ctx context.Context, in *empty.Empty) (*pb.PingResponse, error) {
	var response pb.PingResponse

	if err := s.h.Ping(); err != nil {
		response.Error = err.Error()
		return &response, err
	}

	return &response, nil
}

func (s *ShortenerServer) GetStats(ctx context.Context, in *empty.Empty) (*pb.GetStatsResponse, error) {
	var response pb.GetStatsResponse

	stats, err := s.h.GetStats()
	if err != nil {
		response.Error = err.Error()
		return &response, err
	}
	response.Stats = &pb.Stats{Urls: uint64(stats.URLs), Users: uint64(stats.Users)}

	return &response, nil
}
