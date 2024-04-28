package service

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"

	"github.com/google/uuid"
	"github.com/moataz-hamed/pb/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LaptopServer struct {
	LaptopStore LaptopStore
	ImageStore  ImageStore
	ratingStore RatingStore
	pb.UnimplementedLaptopServiceServer
}

// 1 MegaByte max sizes
const maxImageSize = 1 << 20

// mustEmbedUnimplementedLaptopServiceServer implements pb.LaptopServiceServer.
func (*LaptopServer) mustEmbedUnimplementedLaptopServiceServer() {
	panic("unimplemented")
}

func (server *LaptopServer) RateLaptop(stream pb.LaptopService_RateLaptopServer) error {
	for {
		err := contextError(stream.Context())
		if err != nil {
			return err
		}

		req, err := stream.Recv()
		if err == io.EOF {
			log.Println("No More data")
			break
		}
		if err != nil {
			return logError(status.Errorf(codes.Unknown, "Can't receive stream request:%v", err))
		}

		laptopID := req.GetLaptopId()
		score := req.GetScore()

		log.Printf("received a rate-laptop request: ID:%v,score=%.2f", laptopID, score)

		found, err := server.LaptopStore.Find(laptopID)
		if err != nil {
			return logError(status.Errorf(codes.Internal, "can't find laptop", err))
		}

		if found == nil {
			return logError(status.Errorf(codes.NotFound, "laptop is not found in our database %v", laptopID))
		}

		rating, err := server.ratingStore.Add(laptopID, score)
		if err != nil {
			return logError(status.Errorf(codes.Internal, "can't add rate to the laptop", err))
		}

		res := &pb.RateLaptopResponse{
			LaptopId:     laptopID,
			RatedCount:   rating.Count,
			AverageScore: rating.Sum / float64(rating.Count),
		}

		err = stream.Send(res)
		if err != nil {
			return logError(status.Errorf(codes.Unknown, "Can't send response to the client", err))
		}
	}
	return nil
}

func (server *LaptopServer) UploadImage(stream pb.LaptopService_UploadImageServer) error {
	req, err := stream.Recv()
	if err != nil {

		return logError(status.Errorf(codes.Unknown, "Can't receive image info"))
	}

	laptopID := req.GetInto().GetLaptopId()
	imageType := req.GetInto().GetImageTypes()
	log.Printf("receive an upload image request for laptop %s with image type %s", laptopID, imageType)

	laptop, err := server.LaptopStore.Find(laptopID)
	if err != nil {
		return logError(status.Errorf(codes.Unknown, "Can't findlaptop"))
	}

	if laptop == nil {
		return logError(status.Errorf(codes.InvalidArgument, "laptop %s does not exist", laptopID))
	}

	imageData := bytes.Buffer{}
	imageSize := 0

	for {

		if err := contextError(stream.Context()); err != nil {
			return err
		}

		log.Print("Waiting to receive more data")

		req, err := stream.Recv()
		if err == io.EOF {
			log.Print("No more data")
			break
		}
		if err != nil {
			return logError(status.Errorf(codes.Unknown, "can't receive chunk data %v", err))
		}

		chunk := req.GetChunkData()
		size := len(chunk)
		imageSize += size

		if imageSize > maxImageSize {
			return logError(status.Errorf(codes.InvalidArgument, "Image is too large, Max image size is:%d", maxImageSize))
		}

		_, err = imageData.Write(chunk)
		if err != nil {
			return logError(status.Errorf(codes.Internal, "can't write chunk data %v", err))
		}

	}

	imageID, err := server.ImageStore.Save(laptopID, imageType, imageData)
	if err != nil {
		return logError(status.Errorf(codes.Internal, "can't save image to the store", err))
	}

	res := &pb.UploadImageResponse{
		Id:   imageID,
		Size: uint32(imageSize),
	}

	err = stream.SendAndClose(res)
	if err != nil {
		return logError(status.Errorf(codes.Unknown, "Can't send response %v", err))
	}

	log.Printf("saved image with id:%v and size: %d", imageID, imageSize)
	return nil
}

func logError(err error) error {
	if err != nil {
		log.Print(err)
	}
	return err
}

func NewLaptopServer(store LaptopStore, imageStore ImageStore, ratingStore RatingStore) *LaptopServer {

	return &LaptopServer{LaptopStore: store, ImageStore: imageStore, ratingStore: ratingStore}
}

func (server *LaptopServer) CreateLaptop(ctx context.Context, in *pb.CreateLaptopRequest) (*pb.CreateLaptopResponse, error) {
	laptop := in.GetLaptop()
	log.Println("Received Create Laptop Request with this id:", laptop.Id)
	if len(laptop.Id) > 0 {
		// check if it is a valid uuid
		_, err := uuid.Parse(laptop.Id)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "laptop ID is invalid: %v", err)
		}
	} else {
		id, err := uuid.NewRandom()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Can not generate a new laptop ID: %v", err)
		}
		laptop.Id = id.String()
	}

	err := server.LaptopStore.Save(laptop)
	if err != nil {
		code := codes.Internal
		if errors.Is(err, ErrAlreadyExists) {
			code = codes.AlreadyExists
		}
		return nil, status.Errorf(code, "Can not save laptop to the store:%v", err)
	}

	log.Printf("Saved laptop with id: %s", laptop.Id)
	return &pb.CreateLaptopResponse{Id: laptop.Id}, nil
}

func (server *LaptopServer) SearchLaptop(in *pb.SearchLaptopRequest, stream pb.LaptopService_SearchLaptopServer) error {
	filter := in.GetFilter()
	log.Printf("receive a search-laptop request with filter:%v", filter)

	err := server.LaptopStore.Search(
		filter,
		func(laptop *pb.Laptop) error {
			res := &pb.SearchLaptopResponse{Laptop: laptop}

			err := stream.Send(res)
			if err != nil {
				return err
			}
			log.Printf("Send laptop with id:%s", laptop.GetId())
			return nil
		},
	)

	if err != nil {
		return status.Errorf(codes.Internal, "unexpected error,%w", err)
	}
	return nil
}

func contextError(ctx context.Context) error {
	switch ctx.Err() {
	case context.Canceled:
		return logError(status.Error(codes.Canceled, "request is canceled"))
	case context.DeadlineExceeded:
		return logError(status.Error(codes.DeadlineExceeded, "request Timeout"))
	default:
		return nil

	}
}
