// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package control_plane

import (
	"context"
	"fmt"
	"github.com/sunsingerus/mservice/pkg/controller/client"
	"log"
	"net"
	"strings"
	"testing"

	td "github.com/maxatome/go-testdeep"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	pb "github.com/sunsingerus/mservice/pkg/api/mservice"
)

type testCaseDataServer struct {
	name string
	send []*pb.DataChunk
	want []*pb.DataChunk
}

var (
	_0              uint64 = 0
	_1              uint64 = 1
	_2              uint64 = 2
	testsDataServer        = []testCaseDataServer{
		{
			name: "testCaseDataServer EOF in-bound",
			send: []*pb.DataChunk{
				pb.NewDataChunk(pb.NewMetadata("testfile.txt"), &_0, false, []byte("a")),
				pb.NewDataChunk(nil, &_1, false, []byte("b")),
				pb.NewDataChunk(nil, &_2, true, []byte("c")),
			},
			want: []*pb.DataChunk{
				pb.NewDataChunk(pb.NewMetadata("returnback.file"), &_0, false, []byte("ABC")),
				pb.NewDataChunk(nil, nil, true, nil),
			},
		},
		{
			name: "testCaseDataServer EOF out-of-bound",
			send: []*pb.DataChunk{
				pb.NewDataChunk(pb.NewMetadata("testfile.txt"), &_0, false, []byte("d")),
				pb.NewDataChunk(nil, &_1, false, []byte("e")),
				pb.NewDataChunk(nil, &_2, false, []byte("f")),
				pb.NewDataChunk(nil, nil, true, nil),
			},
			want: []*pb.DataChunk{
				pb.NewDataChunk(pb.NewMetadata("returnback.file"), &_0, false, []byte("DEF")),
				pb.NewDataChunk(nil, nil, true, nil),
			},
		},
	}
)

// testDataServer emulates client-server communication with prepared stream of `in` chunks coming from client to server
// and accumulated `out` chunks, which is accumulated replies from server to client
type testDataServer struct {
	pb.MServiceControlPlane_DataServer
	in  []*pb.DataChunk
	cur int
	out []*pb.DataChunk
}

// Send writes chunk to accumulated `out` chunks set, which can be compared with wanted results
func (t *testDataServer) Send(chunk *pb.DataChunk) error {
	t.out = append(t.out, chunk)
	return nil
}

// Recv reads next chunk from `in` chunks set
func (t *testDataServer) Recv() (*pb.DataChunk, error) {
	if t.cur >= len(t.in) {
		return nil, fmt.Errorf("no more 'in' chunks available for Recv()")
	}
	defer func(t *testDataServer) {
		t.cur++
	}(t)
	return t.in[t.cur], nil
}

// TestServer_Data testsDataServer Data() function without any client.
// Incoming stream is simulated, Data() function is called in bare form.
func TestServer_Data(t *testing.T) {
	for n, test := range testsDataServer {
		t.Logf("testCaseDataServer case: %d\n", n)
		t.Run(test.name, func(t *testing.T) {
			tds := &testDataServer{
				in: test.send,
			}
			srv := &Server{}
			err := srv.Data(tds)
			if td.CmpNoError(t, err) {
				if len(tds.out) != len(test.want) {
					t.Fatalf("unexpected number of chunks")
				}

				for i := range tds.out {
					got := tds.out[i]
					want := test.want[i]
					t.Logf("got chunk %d\nlast: %t\nbytes: %s\n", i, got.GetLast(), string(got.GetBytes()))
					t.Logf("got: %v\n", got)
					t.Logf("exp: %v\n", want)
					td.CmpStruct(t, got, want, td.StructFields{}, test.name)
				}
			}
		})
	}
}

// Test the whole round-trip with simulated network - via custom dialer and listener

const bufSize = 1024 * 1024

var (
	// lis is a custom network-less in-memory listener
	lis *bufconn.Listener
)

func init() {
	// Prepare server-side
	// Start the whole service based on network-less in-memory listener
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterMServiceControlPlaneServer(s, &Server{})
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

type testCaseRoundTrip struct {
	name string
	send string
	want string
}

var (
	testsRoundTrip = []testCaseRoundTrip{
		{
			name: "testCaseRoundTrip 1",
			send: "first round trip test",
			want: "FIRST ROUND TRIP TEST",
		},
		{
			name: "testCaseRoundTrip 2",
			send: "second round trip test",
			want: "SECOND ROUND TRIP TEST",
		},
	}
)

// bufDialer returns the client half of the connection network-less in-memory connection based on custom listener.
func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

// Test_RoundTrip tests the whole round-trip via custom listener/dialer
func Test_RoundTrip(t *testing.T) {

	// Prepare client side
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := pb.NewMServiceControlPlaneClient(conn)

	for n, test := range testsRoundTrip {
		t.Logf("testCaseRoundTrip case: %d\n", n)
		t.Run(test.name, func(t *testing.T) {
			metadata := pb.NewMetadata("testfile.txt")
			buf, err := controller_client.Process(client, metadata, strings.NewReader(test.send))
			if td.CmpNoError(t, err) {
				if buf == nil {
					t.Fatalf("Got no buf")
				}
				td.CmpString(t, buf.String(), test.want, test.name)
			}
		})
	}
}
