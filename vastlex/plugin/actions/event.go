package actions

import (
	"encoding/json"
	"github.com/VastleLLC/VastleX/vastlex/plugin/actions/internalevents"
	"github.com/VastleLLC/VastleX/vastlex/plugin/actions/protobuf"
)

type Event struct {
	*protobuf.EventAction
}

func (i Event) New() *Event {
	return &Event{
		EventAction: &protobuf.EventAction{},
	}
}

func (i *Event) ID() int16 {
	return IDEvent
}

func (i *Event) Decode() (internalevents.Event, bool) {
	if internalevents.Pool[i.EventId] != nil {
		event := internalevents.Pool[i.EventId]()
		err := json.Unmarshal(i.Data, event)
		if err != nil {
			return nil, false
		} else {
			return event, true
		}
	} else {
		return nil, false
	}
}

func (i *Event) Encode(event internalevents.Event) *Event {
	i.EventId = event.ID()
	i.Data, _ = json.Marshal(event)
	return i
}
