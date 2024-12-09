package main

import (
  "context"
  "log"
  "time"
  "flag"
  
  pb "github.com/bovf/gRPCLearning/proto"
  "google.golang.org/grpc"
)

func main() {
  showTime := flag.Bool("time", false, "Show time delta")
  flag.Parse()

  conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
  if err != nil {
    log.Fatalf("did not connect: %v", err)
  }
  defer conn.Close()
  c := pb.NewGreeterClient(conn)

  ctx, cancel := context.WithTimeout(context.Background(), time.Second)
  defer cancel()

  if *showTime {
    clientTime := time.Now().UnixNano()
    r, err := c.GetTimeDelta(ctx, &pb.TimeRequest{ClientTime: clientTime})
    if err != nil {
      log.Fatalf("could not get time delta: %v", err)
    }
    log.Printf("Time Delta: %d nanoseconds", r.Delta)
    log.Printf("ClientTime: %s", time.Unix(0, r.ClientTime).Format(time.RFC3339Nano))
    log.Printf("ServerTime: %s", time.Unix(0, r.ServerTime).Format(time.RFC3339Nano))
  }

  r, err := c.SayHello(ctx, &pb.HelloRequest{Name: "gRPC"})
  if err != nil {
    log.Fatalf("could not greet: %v", err)
  }
  log.Printf("Greetings: %s", r.GetMessage())
}
