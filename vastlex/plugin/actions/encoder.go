package actions

import (
	"encoding/binary"
	"encoding/gob"
	"io"
	"sync"
)

type Encoder struct {
	io    io.Writer
	gob   *gob.Encoder
	mutex sync.Mutex
}

func NewEncoder(Io io.Writer) *Encoder {
	return &Encoder{
		io:  Io,
		gob: gob.NewEncoder(Io),
	}
}

func (encoder *Encoder) Encode(action Action) error {
	encoder.mutex.Lock()
	defer encoder.mutex.Unlock()
	err := binary.Write(encoder.io, binary.LittleEndian, action.ID())
	if err != nil {
		return err
	}
	err = encoder.gob.Encode(action)
	return err
}
