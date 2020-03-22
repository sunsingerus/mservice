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

package mservice

func NewCommand(
	_type CommandType,
	name string,
	version uint32,
	uuid string,
	uuid_reference string,
	seconds int64,
	nanos int32,
	description string,
) *Command {
	return &Command{
		Header: NewHeader(
			uint32(_type),
			name,
			version,
			uuid,
			uuid_reference,
			seconds,
			nanos,
			description,
		),
	}
}

func (m *Command) GetType() CommandType {
	return CommandType(m.GetHeader().GetTypeName().GetType())
}
