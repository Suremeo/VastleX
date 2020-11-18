package actions

import (
	"encoding/binary"
	"encoding/gob"
	"errors"
	"io"
	"sync"
)

type Decoder struct {
	io io.Reader
	gob *gob.Decoder
	mutex sync.Mutex
}

func NewDecoder(Io io.Reader) *Decoder {
	return &Decoder{
		io:  Io,
		gob: gob.NewDecoder(Io),
	}
}

func (decoder *Decoder) Read() (Action, error) {
	decoder.mutex.Lock()
	defer decoder.mutex.Unlock()

	var id int16
	err := binary.Read(decoder.io, binary.LittleEndian, &id)
	if err != nil {
		return nil, err
	}
	pkf := Pool[id]
	if pkf == nil {
		return nil, errors.New("invalid packet recieved")
	}
	pk := pkf()
	return pk, decoder.gob.Decode(&pk)
}