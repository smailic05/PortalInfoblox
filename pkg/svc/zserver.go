package svc

import (
	"context"
	"flag"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/my-app/pkg/pb"
	rpb "github.com/my-app/pkg/resp_pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	// version is the current version of the service
	version = "0.0.1"
)

var (
	serverAddr = flag.String("addr", "localhost:9099", "The server address in the format of host:port")
)

// Default implementation of the MyApp server interface
type server struct {
	pb.UnimplementedMyAppServer
	Description string
	Timestamp   time.Time
	Requests    int64
}

func GRPCConnect() (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// GetVersion returns the current version of the service
func (*server) GetVersion(context.Context, *empty.Empty) (*pb.VersionResponse, error) {
	conn, err := GRPCConnect()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := rpb.NewMyResponderClient(conn)
	resp, err := client.GetVersion(context.Background(), &empty.Empty{})
	if err != nil {
		return nil, err
	}
	return &pb.VersionResponse{Version: resp.Version}, nil
}
func (s *server) UpdateDescription(ctx context.Context, req *pb.UpdateDescriptionRequest) (*pb.UpdateDescriptionResponse, error) {
	if req.GetDescription() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Description can't be empty")
	}
	if req.GetService() != 1 {
		conn, err := GRPCConnect()
		if err != nil {
			return nil, err
		}
		defer conn.Close()
		client := rpb.NewMyResponderClient(conn)
		resp, err := client.UpdateDescription(ctx, &rpb.UpdateDescriptionRequest{Description: req.Description, Service: req.Service})
		if err != nil {
			return nil, err
		}
		return &pb.UpdateDescriptionResponse{Description: resp.Description}, nil
	}
	s.Requests += 1
	s.Description = req.Description
	return &pb.UpdateDescriptionResponse{Description: s.Description}, nil
}

func (s *server) GetDescription(ctx context.Context, req *pb.GetDescriptionRequest) (*pb.GetDescriptionResponse, error) {
	if req.GetService() != 1 {
		conn, err := GRPCConnect()
		if err != nil {
			return nil, err
		}
		defer conn.Close()
		client := rpb.NewMyResponderClient(conn)
		resp, err := client.GetDescription(ctx, &rpb.GetDescriptionRequest{Service: req.Service})
		if err != nil {
			return nil, err
		}
		return &pb.GetDescriptionResponse{Description: resp.Description}, nil

	}
	s.Requests += 1
	return &pb.GetDescriptionResponse{Description: s.Description}, nil
}

// Returns uptime in seconds
func (s *server) GetUptime(ctx context.Context, req *pb.GetUptimeRequest) (*pb.GetUptimeResponse, error) {
	if req.GetService() != 1 {
		conn, err := GRPCConnect()
		if err != nil {
			return nil, err
		}
		defer conn.Close()
		client := rpb.NewMyResponderClient(conn)
		resp, err := client.GetUptime(ctx, &rpb.GetUptimeRequest{Service: req.Service})
		if err != nil {
			return nil, err
		}
		return &pb.GetUptimeResponse{Uptime: resp.Uptime}, nil

	}
	s.Requests += 1
	uptime := time.Now().Unix() - s.Timestamp.Unix()
	return &pb.GetUptimeResponse{Uptime: uptime}, nil
}

func (s *server) GetRequests(ctx context.Context, req *pb.GetRequestsRequest) (*pb.GetRequestsResponse, error) {
	if req.GetService() != 1 {
		conn, err := GRPCConnect()
		if err != nil {
			return nil, err
		}
		defer conn.Close()
		client := rpb.NewMyResponderClient(conn)
		resp, err := client.GetRequests(ctx, &rpb.GetRequestsRequest{Service: req.Service})
		if err != nil {
			return nil, err
		}
		return &pb.GetRequestsResponse{Requests: resp.Requests}, nil

	}
	s.Requests += 1
	return &pb.GetRequestsResponse{Requests: s.Requests}, nil
}

// NewBasicServer returns an instance of the default server interface
func NewBasicServer() (pb.MyAppServer, error) {
	return &server{Description: "Portal", Timestamp: time.Now()}, nil
	//TODO use Viper
	//Incapsulate client
}
