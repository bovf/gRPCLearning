package main

import (
  "context"
  "flag"
  "fmt"
  "log"
  "strings"
  "time"
  
  pb "github.com/bovf/gRPCLearning/proto"
  "google.golang.org/grpc"
)

func main() {
  showTime := flag.Bool("time", false, "Show time delta")
  addLDAP := flag.Bool("add-ldap", false, "Add an entry to LDAP server")
  searchLDAP := flag.Bool("search-ldap", false, "Search LDAP server")
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

  if *addLDAP {
    if flag.NArg() < 2 {
      log.Fatal("Usage: --add-ldap <dn> <attr1=value1> <attr2=value2> ...")
    }
    dn := flag.Arg(0)
    attrs := make(map[string]*pb.LDAPAttribute)
    
    // Add objectClass attribute automatically
    attrs["objectClass"] = &pb.LDAPAttribute{Values: []string{"top", "person", "organizationalPerson", "inetOrgPerson"}}
    
    for _, arg := range flag.Args()[1:] {
      parts := strings.SplitN(arg, "=", 2)
      if len(parts) != 2 {
          log.Fatalf("Invalid attribute format: %s", arg)
      }
      attrs[parts[0]] = &pb.LDAPAttribute{Values: []string{parts[1]}}
    }
    req := &pb.LDAPAddRequest{DN: dn, Attributes: attrs}
    reply, err := c.AddLDAP(ctx, req)
    if err != nil {
      log.Fatalf("Could not add LDAP entry: %v", err)
    }
    if reply.Success {
      fmt.Println("LDAP entry added successfully")
    } else {
      fmt.Printf("Failed to add LDAP entry: %s\n", reply.Error)
    }
  }
  if *searchLDAP {
    if flag.NArg() != 3 {
      log.Fatal("Usage: --search-ldap <baseDN> <filter> <attr1,attr2,...>")
    }
    baseDN := flag.Arg(0)
    filter := flag.Arg(1)
    attributes := strings.Split(flag.Arg(2), ",")
    req := &pb.LDAPSearchRequest{BaseDN: baseDN, Filter: filter, Attributes: attributes}
    reply, err := c.SearchLDAP(ctx, req)
    if err != nil {
      log.Fatalf("Could not search LDAP: %v", err)
    }
    for _, entry := range reply.Entries {
      fmt.Printf("DN: %s\n", entry.DN)
      for attrName, attrValues := range entry.Attributes {
          fmt.Printf("  %s: %v\n", attrName, attrValues.Values)
      }
      fmt.Println()
    }
  }
  r, err := c.SayHello(ctx, &pb.HelloRequest{Name: "gRPC"})
  if err != nil {
    log.Fatalf("could not greet: %v", err)
  }
  log.Printf("Greetings: %s", r.GetMessage())
}
