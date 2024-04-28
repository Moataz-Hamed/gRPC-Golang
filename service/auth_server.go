package service

import (
	"context"

	"github.com/moataz-hamed/pb/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServer struct {
	userStore  UserStore
	jwtManager *JWTManager
	pb.UnimplementedAuthServiceServer
}

// mustEmbedUnimplementedAuthServiceServer implements pb.AuthServiceServer.
func (*AuthServer) mustEmbedUnimplementedAuthServiceServer() {
	panic("unimplemented")
}

func NewAuthServer(userStore UserStore, jwtManager *JWTManager) *AuthServer {
	return &AuthServer{userStore, jwtManager, pb.UnimplementedAuthServiceServer{}}
}

func (server *AuthServer) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginReponse, error) {
	user, err := server.userStore.Find(in.GetUsername())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Can't find user:%v", err)
	}

	if user == nil || !user.IsCorrectPassword(in.GetPassword()) {
		return nil, status.Errorf(codes.NotFound, "Incorrect username or password")
	}

	token, err := server.jwtManager.Generate(user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Can't create token fot user:%v", user)
	}

	res := &pb.LoginReponse{
		AccessToken: token,
	}

	return res, nil
}
