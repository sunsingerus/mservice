package server_test

import (
	"context"
	"net"
	"testing"

	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	pb "github.com/sunsingerus/mservice/pkg/api/mservice"
	"github.com/sunsingerus/mservice/pkg/transiever/service"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterMServiceControlPlaneServer(s, &transiever_service.MServiceControlPlaneEndpoint{})
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

var (
	request = pb.NewDataChunk(pb.NewMetadata("qwe.txt"), nil, true, []byte("some data goes here"))
	zero    uint64
	reply   = pb.NewDataChunk(pb.NewMetadata("returnback.file"), &zero, false, []byte("SOME DATA GOES HERE"))
)

func TestData(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := pb.NewMServiceControlPlaneClient(conn)
	stream, err := client.Data(ctx)
	if err != nil {
		t.Fatalf("TestData failed: %v", err)
	}
	if err := stream.Send(request); err != nil {
		t.Fatalf("TestData failed: %v", err)
	}
	if err := stream.CloseSend(); err != nil {
		t.Fatalf("TestData failed: %v", err)
	}
	got, err := stream.Recv()
	if err != nil {
		t.Fatalf("TestData failed: %v", err)
	}
	if !proto.Equal(got, reply) {
		t.Fatalf("stream.Recv() =\n%v\nwant=\n%v", got, reply)
	}
}
