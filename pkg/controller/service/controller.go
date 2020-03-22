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

package controller_service

import (
	log "github.com/golang/glog"

	pb "github.com/sunsingerus/mservice/pkg/api/mservice"
)

func IncomingCommandsHandler(incomingQueue, outgoingQueue chan *pb.Command) {
	log.Infof("Start IncomingCommandsHandler()")
	defer log.Infof("Exit IncomingCommandsHandler()")

	for {
		cmd := <-incomingQueue
		log.Infof("Got cmd %v", cmd)
		if cmd.GetType() == pb.CommandType_COMMAND_ECHO_REQUEST {
			command := pb.NewCommand(
				pb.CommandType_COMMAND_ECHO_REPLY,
				"",
				0,
				pb.CreateNewUUID(),
				"reference: "+cmd.GetHeader().GetUuid().StringValue,
				0,
				0,
				"desc",
			)
			outgoingQueue <- command
		}
	}
}
