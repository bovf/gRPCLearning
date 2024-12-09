package main

import (
  "context"
  "time"
  "net"

  pb "github.com/bovf/gRPCLearning/proto"
  "github.com/bovf/gRPCLearning/logging"
  "github.com/bovf/gRPCLearning/ldap"
  "google.golang.org/grpc"
)

type server struct {
  pb.UnimplementedGreeterServer
  logger *logging.Logger
  ldapClient *ldap.LDAPClient
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

func (s *server) SearchLDAP(ctx context.Context, in *pb.LDAPSearchRequest) (*pb.LDAPSearchReply, error) {
  s.logger.LogRPC("SearchLDAP", "")
  entries, err := s.ldapClient.Search(in.BaseDN, in.Filter, in.Attributes)
  if err != nil {
    return nil, err
  }

  var results []*pb.LDAPEntry
  for _,entry := range entries {
    attrs := make(map[string]*pb.LDAPAttribute)
    for _, attr := range entry.Attributes {
      attrs[attr.Name] = &pb.LDAPAttribute{Values: attr. Values}
    }
    results = append(results, &pb.LDAPEntry{
      DN: entry.DN,
      Attributes: attrs, 
    })
  }

  return &pb.LDAPSearchReply{Entries: results}, nil
}

func (s *server) AddLDAP(ctx context.Context, in *pb.LDAPAddRequest) (*pb.LDAPAddReply, error) {
  s.logger.LogRPC("AddLDAP", "")
  attrs := make(map[string][]string)
  for key, attr := range in.Attributes {
    attrs[key] = attr.Values
  }
  err := s.ldapClient.Add(in.DN, attrs)
  if err != nil {
    return &pb.LDAPAddReply{Success: false, Error: err.Error()}, nil
  }
  return &pb.LDAPAddReply{Success: true}, nil
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
  ldapClient, err := ldap.NewLDAPClient("localhost:389")
  if err != nil {
    logger.Fatalf("failed to connect to LDAP: %v", err)
  }
  defer ldapClient.Close()

  lis, err := net.Listen("tcp", ":50051")
  if err!= nil {
    logger.Fatalf("failed to listen: %v", err)
  }
  s := grpc.NewServer(
    grpc.UnaryInterceptor(loggingInterceptor(logger)),
  )

  pb.RegisterGreeterServer(s, &server{
    logger: logger,
    ldapClient: ldapClient,
  })
  logger.Printf("server listening at %v", lis.Addr())
  if err := s.Serve(lis); err != nil {
    logger.Fatalf("failed to serve :%v", err)
  }
}
