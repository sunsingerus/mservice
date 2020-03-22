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

import (
	"github.com/golang/protobuf/ptypes/timestamp"
)

func NewHeader(
	_type uint32,
	name string,
	version uint32,
	uuid string,
	uuid_reference string,
	seconds int64,
	nanos int32,
	description string,
) *Header {

	h := new(Header)

	h.SetTypeOrName(_type, name)

	if version > 0 {
		h.SetVersion(version)
	}

	if uuid == "" {
		h.SetUUID(CreateNewUUID())
	} else {
		h.SetUUID(uuid)
	}

	if uuid_reference != "" {
		h.SetUUIDReference(uuid_reference)
	}

	if seconds > 0 {
		h.SetTimestamp(seconds, nanos)
	}

	if description != "" {
		h.SetDescription(description)
	}

	return h
}

func (h *Header) ensureTypeName() {
	if h.TypeName == nil {
		h.TypeName = new(TypeName)
	}
}

func (h *Header) SetTypeOrName(_type uint32, name string) {
	if _type > 0 {
		h.SetType(_type)
	} else {
		h.SetName(name)
	}
}

func (h *Header) SetType(_type uint32) {
	h.ensureTypeName()
	if h.TypeName.TypeOptional == nil {
		h.TypeName.TypeOptional = new(TypeName_Type)
	}
	h.TypeName.TypeOptional.(*TypeName_Type).Type = _type
}

func (h *Header) SetName(name string) {
	h.ensureTypeName()
	if h.TypeName.NameOptional == nil {
		h.TypeName.NameOptional = new(TypeName_Name)
	}
	h.TypeName.NameOptional.(*TypeName_Name).Name = name
}

func (h *Header) SetVersion(version uint32) {
	if h.VersionOptional == nil {
		h.VersionOptional = new(Header_Version)
		h.VersionOptional.(*Header_Version).Version = version
	}
}

func (h *Header) SetUUID(uuid string) {
	if h.UuidOptional == nil {
		h.UuidOptional = new(Header_Uuid)
		if h.UuidOptional.(*Header_Uuid).Uuid == nil {
			h.UuidOptional.(*Header_Uuid).Uuid = NewUUID(uuid)
		}

	}
}

func (h *Header) SetUUIDReference(uuid string) {
	if h.UuidReferenceOptional == nil {
		h.UuidReferenceOptional = new(Header_UuidReference)
		if h.UuidReferenceOptional.(*Header_UuidReference).UuidReference == nil {
			h.UuidReferenceOptional.(*Header_UuidReference).UuidReference = NewUUID(uuid)
		}
	}
}

func (h *Header) SetTimestamp(seconds int64, nanos int32) {
	if h.TimestampOptional == nil {
		h.TimestampOptional = new(Header_Ts)
		if h.TimestampOptional.(*Header_Ts).Ts == nil {
			h.TimestampOptional.(*Header_Ts).Ts = new(timestamp.Timestamp)
			h.TimestampOptional.(*Header_Ts).Ts.Seconds = seconds
			h.TimestampOptional.(*Header_Ts).Ts.Nanos = nanos
		}
	}
}

func (h *Header) SetDescription(description string) {
	if h.DescriptionOptional == nil {
		h.DescriptionOptional = new(Header_Description)
		h.DescriptionOptional.(*Header_Description).Description = description
	}
}
