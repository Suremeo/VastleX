package actions

import (
	"bytes"
	"encoding/binary"
	"github.com/golang/protobuf/proto"
	"io"
	"sync"
)

type Decoder struct {
	io         io.Reader
	mutex      sync.Mutex
	readBuffer *bytes.Buffer
}

func NewDecoder(Io io.Reader) *Decoder {
	buf := bytes.NewBuffer(nil)
	return &Decoder{
		io:         Io,
		readBuffer: buf,
	}
}

func (decoder *Decoder) Read() (Action, error) {
	defer decoder.readBuffer.Reset()
	decoder.mutex.Lock()
	defer decoder.mutex.Unlock()
	var length int64
	err := binary.Read(decoder.io, binary.LittleEndian, &length)
	if err != nil {
		return nil, err
	}
	temp := make([]byte, length)
	_, err = decoder.io.Read(temp)
	if err != nil {
		return nil, err
	}
	decoder.readBuffer.Write(temp)
	var id int16
	err = binary.Read(decoder.readBuffer, binary.LittleEndian, &id)
	if err != nil {
		return nil, err
	}
	pkf := Pool[id]
	if pkf == nil {
		return decoder.Read()
	}
	pk := pkf()
	err = proto.Unmarshal(decoder.readBuffer.Bytes(), pk)
	return pk, err
}
