package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/moataz-hamed/pb/pb"
	"github.com/moataz-hamed/service"
	"google.golang.org/grpc"
)

const (
	secretKey     = "secret"
	tokenDuration = 15 * time.Minute
)

func seedUser(userStore service.UserStore) error {
	err := createUser(userStore, "Moataz", "password", "admin")
	if err != nil {
		return err
	}
	return createUser(userStore, "user1", "password", "user")
}

func createUser(userStore service.UserStore, username, password, role string) error {
	user, err := service.NewUser(username, password, role)
	if err != nil {
		return err
	}

	return userStore.Save(user)
}

func accessibleRoles() map[string][]string {
	const laptopServicePath = "/mypackage.LaptopService/"
	return map[string][]string{
		laptopServicePath + "CreateLaptop": {"admin"},
		laptopServicePath + "UploadImage":  {"admin"},
		laptopServicePath + "RateLaptop":   {"admin", "user"},
	}
}

func main() {
	port := flag.Int("port", 0, "the server port")
	flag.Parse()
	log.Printf("start server on port %d", *port)

	userStore := service.NewInMemoryUserStore()

	err := seedUser(userStore)
	if err != nil {
		log.Fatal("Error:%v", err)
	}

	jwtManager := service.NewJWTManager(secretKey, tokenDuration)
	authServer := service.NewAuthServer(userStore, jwtManager)

	laptopServer := service.NewLaptopServer(service.NewInMemoryLaptopStore(), service.NewDiskImageStore("img"), service.NewInMemoryRatingStore())

	interceptor := service.NewAuthInterceptor(*jwtManager, accessibleRoles())

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.Unary()),
		grpc.StreamInterceptor(interceptor.Stream()),
	)
	pb.RegisterLaptopServiceServer(grpcServer, laptopServer)
	pb.RegisterAuthServiceServer(grpcServer, authServer)

	address := fmt.Sprintf("0.0.0.0:%d", *port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("Can not start server:", err)
	}

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("Can not start server2:", err)
	}
}
