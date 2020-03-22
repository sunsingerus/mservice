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

package transiever

import (
	log "github.com/golang/glog"
	pb "github.com/sunsingerus/mservice/pkg/api/mservice"
	"io"
)

type CommandSendReceiveInterface interface {
	Send(*pb.Command) error
	Recv() (*pb.Command, error)
}

func CommandsExchangeEndlessLoop(i CommandSendReceiveInterface) {
	waitIncoming := make(chan bool)
	waitOutgoing := make(chan bool)

	go func() {
		for {
			msg, err := i.Recv()
			if msg != nil {
				log.Infof("CommandsExchangeEndlessLoop.Recv() got msg")
				GetIncomingQueue() <- msg
			}
			if err == nil {
				// All went well, ready to receive more data
			} else if err == io.EOF {
				// Correct EOF
				log.Infof("CommandsExchangeEndlessLoop.Recv() get EOF")

				close(waitIncoming)
				return
			} else {
				// Stream broken
				log.Infof("CommandsExchangeEndlessLoop.Recv() got err: %v", err)

				close(waitIncoming)
				return
			}
		}
	}()

	go func() {
		for {
			select {
			case <-waitIncoming:
				// Incoming stream from this client is closed/broken, no need to wait commands for it
				close(waitOutgoing)
				return
			case command := <-GetOutgoingQueue():
				log.Infof("got command to send")
				err := i.Send(command)
				if err == nil {
					// All went well
					log.Infof("CommandsExchangeEndlessLoop.Send() OK")
				} else if err == io.EOF {
					log.Infof("CommandsExchangeEndlessLoop.Send() got EOF")

					close(waitOutgoing)
					return
				} else {
					log.Fatalf("CommandsExchangeEndlessLoop.Send() got err: %v", err)

					close(waitOutgoing)
					return
				}
			}
		}
	}()

	<-waitIncoming
	<-waitOutgoing
}

var (
	maxIncomingOutstanding int32 = 100
	incomingQueue          chan *pb.Command
	maxOutgoingOutstanding int32 = 100
	outgoingQueue          chan *pb.Command
)

func Init() {
	incomingQueue = make(chan *pb.Command, maxIncomingOutstanding)
	outgoingQueue = make(chan *pb.Command, maxOutgoingOutstanding)
}

func GetOutgoingQueue() chan *pb.Command {
	return outgoingQueue
}

func GetIncomingQueue() chan *pb.Command {
	return incomingQueue
}
