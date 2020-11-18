package actions

type Action interface {
	ID() int16
}

var Pool = map[int16]func() Action{
	IDInit: func() Action { return &Init{} },
}

const (
	IDInit = int16(iota)
)
