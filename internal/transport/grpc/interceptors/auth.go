package interceptors

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/Jeskay/micsvc/internal/auth"
	"github.com/Jeskay/micsvc/internal/dto"
	proto "github.com/Jeskay/micsvc/protos"
)

func NewAuthUnary(svc *auth.Service) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		if _, ok := req.(*proto.RegisterUser); ok {
			resp, err = handler(ctx, req)
			return
		}
		if _, ok := req.(*proto.LoginRequest); ok {
			resp, err = handler(ctx, req)
			return
		}
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.InvalidArgument, "missing metadata")
		}
		jwtHeader, ok := md["jwt"]
		if !ok {
			return nil, status.Error(codes.InvalidArgument, "missing jwt")
		}
		token := jwtHeader[0]
		userID, err := svc.VerifyToken(token)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
		}
		resp, err = handler(context.WithValue(ctx, dto.ID, userID), req)
		return
	}
}
