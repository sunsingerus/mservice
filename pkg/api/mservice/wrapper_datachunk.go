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

func NewDataChunk(metadata *Metadata, offset *uint64, last bool, data []byte) *DataChunk {
	chunk := &DataChunk{
		Header: NewHeader(
			uint32(DataChunkType_DATA_CHUNK_DATA),
			"",
			0,
			"123",
			"",
			0,
			0,
			"data chunk",
		),
		Bytes: data,
	}
	if metadata != nil {
		chunk.SetMetadata(metadata)
	}

	if offset != nil {
		chunk.SetOffset(*offset)
	}

	if last {
		chunk.SetLast(last)
	}

	return chunk
}

func (dc *DataChunk) ensureMetadata() {
	if dc.MetadataOptional == nil {
		dc.MetadataOptional = new(DataChunk_Metadata)
		dc.MetadataOptional.(*DataChunk_Metadata).Metadata = new(Metadata)
	}
}

func (dc *DataChunk) SetMetadata(metadata *Metadata) {
	if dc.MetadataOptional == nil {
		dc.MetadataOptional = new(DataChunk_Metadata)
	}
	dc.MetadataOptional.(*DataChunk_Metadata).Metadata = metadata
}

func (dc *DataChunk) SetLen(len uint64) {
	if dc.LenOptional == nil {
		dc.LenOptional = new(DataChunk_Len)
		dc.LenOptional.(*DataChunk_Len).Len = len
	}
}

func (dc *DataChunk) GetOffsetWithAvailabilityReport() (uint64, bool) {
	if x, ok := dc.GetOffsetOptional().(*DataChunk_Offset); ok {
		return x.Offset, true
	}
	return 0, false
}

func (dc *DataChunk) SetOffset(offset uint64) {
	if dc.OffsetOptional == nil {
		dc.OffsetOptional = new(DataChunk_Offset)
		dc.OffsetOptional.(*DataChunk_Offset).Offset = offset
	}
}

func (dc *DataChunk) SetLast(last bool) {
	if dc.LastOptional == nil {
		dc.LastOptional = new(DataChunk_Last)
		dc.LastOptional.(*DataChunk_Last).Last = last
	}
}
