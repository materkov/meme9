package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/twitchtv/twirp"

	pb "github.com/materkov/meme9/photos/internal/authclient/pb/github.com/materkov/meme9/api/auth"
)

type AuthClient interface {
	VerifyToken(ctx context.Context, req *pb.VerifyTokenRequest) (*pb.VerifyTokenResponse, error)
}

type Service struct {
	client AuthClient
}

func New(client AuthClient) *Service {
	return &Service{
		client: client,
	}
}

var ErrInvalidToken = errors.New("invalid token")

func (s *Service) Auth(ctx context.Context, header string) (string, error) {
	header = strings.TrimSpace(header)
	header = strings.TrimPrefix(header, "Bearer ")
	if header == "" {
		return "", ErrInvalidToken
	}

	req := &pb.VerifyTokenRequest{
		Token: header,
	}
	resp, err := s.client.VerifyToken(ctx, req)
	if err != nil {
		var twerr twirp.Error
		if errors.As(err, &twerr) && twerr.Msg() == "invalid_token" {
			return "", ErrInvalidToken
		} else {
			return "", fmt.Errorf("auth verify failed: %w", err)
		}
	}

	return resp.UserId, nil
}
