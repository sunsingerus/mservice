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

package controller_client

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/proto"
	"github.com/sunsingerus/mservice/pkg/controller/client/internal/grpctest"

	datapb "github.com/sunsingerus/mservice/pkg/api/mservice"
)

type s struct {
	grpctest.Tester
}

func Test(t *testing.T) {
	grpctest.RunSubTests(t, s{})
}

var (
	request = datapb.NewDataChunk(datapb.NewMetadata("qwe.txt"), nil, true, []byte("some data goes here"))
	reply   = datapb.NewDataChunk(datapb.NewMetadata("returnback.file"), nil, true, []byte("SOME DATA GOES HERE"))
)

// TestData performs server-less test of client. Server side is mocked by mockgen-generated code.
func (s) TestData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock for the stream returned by Data() function
	stream := NewMockMServiceControlPlane_DataClient(ctrl)

	// Set expectation on sending.
	stream.EXPECT().Send(
		gomock.Any(),
	).Return(nil)

	// Set expectation on receiving.
	stream.EXPECT().Recv().Return(reply, nil)
	stream.EXPECT().CloseSend().Return(nil)

	// Create mock for the client interface.
	dataclient := NewMockMServiceControlPlaneClient(ctrl)
	// Set expectation on Data
	dataclient.EXPECT().Data(
		gomock.Any(),
	).Return(stream, nil)

	if err := testDataFn(dataclient); err != nil {
		t.Fatalf("Test failed: %v", err)
	}
}

// testDataFn makes call from client (real) to mocked server and test results
func testDataFn(client datapb.MServiceControlPlaneClient) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := client.Data(ctx)
	if err != nil {
		return err
	}
	if err := stream.Send(request); err != nil {
		return err
	}
	if err := stream.CloseSend(); err != nil {
		return err
	}
	got, err := stream.Recv()
	if err != nil {
		return err
	}
	if !proto.Equal(got, reply) {
		return fmt.Errorf("stream.Recv() = %v, want %v", got, reply)
	}
	return nil
}
