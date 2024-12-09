package main

import (
  "context"
  "time"
  "net"

  pb "github.com/bovf/gRPCLearning/proto"
  "github.com/bovf/gRPCLearning/logging"
  "google.golang.org/grpc"
)

type server struct {
  pb.UnimplementedGreeterServer
  logger *logging.Logger
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
  start := time.Now()
  reply := &pb.HelloReply{Message: "Hello " + in.Name}
  duration := time.Since(start)
  s.logger.LogRPC("SayHello", duration.String())
  return reply, nil
}

func (s *server) GetTimeDelta(ctx context.Context, in *pb.TimeRequest) (*pb.TimeReply, error) {
  serverTime := time.Now().UnixNano()
  clientTime := in.ClientTime
  delta := serverTime - clientTime

  s.logger.LogRPC("GetTimeDelta", "")
  return &pb.TimeReply{
    ServerTime: serverTime,
    ClientTime: clientTime,
    Delta: delta,
  }, nil
}

func loggingInterceptor(logger *logging.Logger) grpc.UnaryServerInterceptor {
  return func (ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface {}, error) {
    start := time.Now()
    h, err := handler (ctx, req)
    duration := time.Since(start)
    logger.LogRPC(info.FullMethod, duration.String())
    return h, err 
  }
}

func main () {
  logger := logging.NewLogger()
  lis, err := net.Listen("tcp", ":50051")
  if err!= nil {
    logger.Fatalf("failed to listen: %v", err)
  }
  s := grpc.NewServer(
    grpc.UnaryInterceptor(loggingInterceptor(logger)),
  )

  pb.RegisterGreeterServer(s, &server{logger: logger})
  logger.Printf("server listening at %v", lis.Addr())
  if err := s.Serve(lis); err != nil {
    logger.Fatalf("failed to serve :%v", err)
  }
}
