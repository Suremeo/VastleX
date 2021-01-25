package actions

import (
	"bytes"
	"encoding/binary"
	"github.com/golang/protobuf/proto"
	"io"
	"sync"
)

type Encoder struct {
	io          io.Writer
	mutex       sync.Mutex
	writeBuffer *bytes.Buffer
}

func NewEncoder(Io io.Writer) *Encoder {
	buf := bytes.NewBuffer(nil)
	return &Encoder{
		io:          Io,
		writeBuffer: buf,
	}
}

func (encoder *Encoder) Encode(action Action) error {
	encoder.mutex.Lock()
	defer encoder.writeBuffer.Reset()
	defer encoder.mutex.Unlock()
	err := binary.Write(encoder.writeBuffer, binary.LittleEndian, action.ID())
	if err != nil {
		return err
	}
	b, err := proto.Marshal(action)
	if err != nil {
		return err
	}
	encoder.writeBuffer.Write(b)
	err = binary.Write(encoder.io, binary.LittleEndian, int64(encoder.writeBuffer.Len()))
	if err != nil {
		return err
	}
	_, err = encoder.io.Write(encoder.writeBuffer.Bytes())
	if err != nil {
		return err
	}
	return err
}
