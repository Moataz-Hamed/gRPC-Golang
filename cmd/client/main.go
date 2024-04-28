package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/moataz-hamed/client"
	"github.com/moataz-hamed/pb/pb"
	"github.com/moataz-hamed/sample"
	"google.golang.org/grpc"
)

func testUploadImage(laptopClient client.LaptopClient) {
	laptop := sample.NewLaptop()
	laptopClient.CreateLaptop(laptop)
	laptopClient.UploadImage(laptop.GetId(), "tmp/angry.png")
}

func testCreateLaptop(laptopClient client.LaptopClient) {
	laptopClient.CreateLaptop(sample.NewLaptop())
}

func testSeachLaptop(laptopClient client.LaptopClient) {
	for i := 0; i < 10; i++ {
		laptopClient.CreateLaptop(sample.NewLaptop())
	}

	filter := &pb.Filter{
		MaxPriceUsd: 3000,
		MinCpuCores: 4,
		MinCpuGhz:   2.5,
		MinRam:      &pb.Memory{Value: 8, Unit: pb.Memory_GIGABYTE},
	}

	laptopClient.SerachLaptop(filter)
}

func testRateLaptop(laptopClient client.LaptopClient) {
	n := 3
	laptopIDS := make([]string, n)

	for i := 0; i < n; i++ {
		laptop := sample.NewLaptop()
		laptopIDS[i] = laptop.Id
		laptopClient.CreateLaptop(laptop)
	}

	scores := make([]float64, n)
	for {
		fmt.Print("rate laptop? (y/n)	")
		var answer string
		fmt.Scan(&answer)

		if strings.ToLower(answer) != "y" {
			break
		}
		for i := 0; i < n; i++ {
			scores[i] = sample.RandomLaptopScore()
		}

		err := laptopClient.RateLaptop(laptopIDS, scores)
		if err != nil {
			log.Fatal(err)
		}
	}
}

const (
	username        = "Moataz"
	password        = "password"
	refreshDuration = 30 * time.Second
)

func authMethods() map[string]bool {
	const laptopServicePath = "/mypackage.LaptopService/"
	return map[string]bool{
		laptopServicePath + "CreateLaptop": true,
		laptopServicePath + "UploadImage":  true,
		laptopServicePath + "RateLaptop":   true,
	}
}

func main() {
	serverAddress := flag.String("address", "", "The server address")
	flag.Parse()
	log.Printf("dial server %s", *serverAddress)

	cc1, err := grpc.Dial(*serverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatal("cannot dial server", err)
	}

	authClient := client.NewAuthClient(cc1, username, password)
	interceptor, err := client.NewAuthInterceptor(authClient, authMethods(), refreshDuration)
	if err != nil {
		log.Fatal("Can't create auth interceptor:", err)
	}

	cc2, err := grpc.Dial(
		*serverAddress,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(interceptor.Unary()),
		grpc.WithStreamInterceptor(interceptor.Stream()))
	if err != nil {
		log.Fatal("cannot dial server", err)
	}

	laptopClient := client.NewLaptopClient(cc2)
	testRateLaptop(*laptopClient)
}
