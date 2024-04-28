package client

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/moataz-hamed/pb/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LaptopClient struct {
	service pb.LaptopServiceClient
}

func NewLaptopClient(cc *grpc.ClientConn) *LaptopClient {
	service := pb.NewLaptopServiceClient(cc)
	return &LaptopClient{service: service}

}

func (laptopClient LaptopClient) CreateLaptop(laptop *pb.Laptop) {

	req := &pb.CreateLaptopRequest{
		Laptop: laptop,
	}

	res, err := laptopClient.service.CreateLaptop(context.Background(), req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.AlreadyExists {
			log.Println("Already exists")
		} else {
			log.Fatal("Can't Create laptop", err)
		}
		return
	}

	log.Printf("Laptop is created with id:%s", res.Id)
}

func (laptopClient *LaptopClient) SerachLaptop(filter *pb.Filter) {
	log.Println("search filter:", filter)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.SearchLaptopRequest{Filter: filter}
	stream, err := laptopClient.service.SearchLaptop(ctx, req)
	if err != nil {
		log.Fatal("can't search", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			return
		}

		if err != nil {
			log.Fatal("can't receive response: ", err)
		}

		laptop := res.GetLaptop()
		log.Print("-found:", laptop)
	}
}

func (laptopClient *LaptopClient) UploadImage(laptopID string, path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal("can't open image file", err)
	}

	defer file.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stream, err := laptopClient.service.UploadImage(ctx)
	if err != nil {
		log.Fatal("can't upload image:", err)
	}

	req := &pb.UploadImageRequest{
		Data: &pb.UploadImageRequest_Into{
			Into: &pb.ImageInfo{
				LaptopId:   laptopID,
				ImageTypes: filepath.Ext(path), //get the extension of the file
			},
		},
	}

	err = stream.Send(req)
	if err != nil {
		log.Fatal("can't send image info:", err)
	}

	reader := bufio.NewReader(file)
	buffer := make([]byte, 1024)

	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("can't read chunk to buffer")
		}

		req := &pb.UploadImageRequest{
			Data: &pb.UploadImageRequest_ChunkData{
				ChunkData: buffer[:n],
			},
		}

		err = stream.Send(req)
		if err != nil {
			log.Fatal("Can't send chunk to server", err)
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		err2 := stream.RecvMsg(nil)
		log.Fatal("Can't receive and close server connection:", err, err2)
	}

	log.Printf("Image uploaded with id: %s and size: %d", res.GetId(), res.GetSize())
}

func (laptopClient *LaptopClient) RateLaptop(laptopIDs []string, scores []float64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stream, err := laptopClient.service.RateLaptop(ctx)
	if err != nil {
		return fmt.Errorf("can't rate laptop %v", err)
	}

	waitReponse := make(chan error)

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				log.Println("no more responses")
				waitReponse <- nil
				return
			}
			if err != nil {
				waitReponse <- fmt.Errorf("can't receive stream response %v", err)
				return
			}
			log.Println("Received response:", res)
		}
	}()

	for i, laptopID := range laptopIDs {
		req := &pb.RateLaptopRequest{
			LaptopId: laptopID,
			Score:    scores[i],
		}

		err = stream.Send(req)
		if err != nil {
			return fmt.Errorf("can't send stream request %v - %v", err, stream.RecvMsg(nil))
		}

		log.Println("sent request", req)
	}

	err = stream.CloseSend() // tell the server that we won't send anymore data
	if err != nil {
		return fmt.Errorf("can't close send:%v", err)
	}

	err = <-waitReponse

	return err
}
