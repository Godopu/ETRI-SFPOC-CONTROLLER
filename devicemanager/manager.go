package devicemanager

import (
	"context"
)

type task struct {
	Event int
	Ctx   context.Context
}

var taskQueue = make(chan *task, 100)

type managerKey int

const (
	DISCOVERY int = iota
)

const (
	waitResponseKey managerKey = iota
	parameterKey
)

func NewManager() (func(), func()) {
	ctx, cancel := context.WithCancel(context.Background())
	run := func() {
		go run(ctx)
	}

	return run, cancel
}

func run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case t := <-taskQueue:

			switch t.Event {
			case DISCOVERY:
				p := t.Ctx.Value(managerKey(parameterKey))
				b, err := RegisterDevice(p.(*map[string]interface{}), t.Ctx.Done())

				if err != nil {
					continue
				}

				respCh := t.Ctx.Value(managerKey(waitResponseKey)).(chan []byte)
				if respCh != nil {
					// , _ := json.Marshal(map[string]interface{}{"hello": "World"})
					respCh <- b
				}
			}
		}
	}
}
