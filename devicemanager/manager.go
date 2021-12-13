package devicemanager

import (
	"context"
	"etrismartfarmpoccontroller/model"
)

type task struct {
	Event int
	Ctx   context.Context
}

var taskQueue = make(chan *task, 100)

type managerKey int

const (
	DISCOVERY int = iota
	DISCONNECTED
	STATUSREPORT
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
				b, err := RegisterDevice(p.(map[string]interface{}), t.Ctx.Done())
				if err != nil {
					continue
				}

				respCh, ok := t.Ctx.Value(managerKey(waitResponseKey)).(chan []byte)
				if !ok {
					return
				}
				respCh <- b

			case DISCONNECTED:
				p := t.Ctx.Value(managerKey(parameterKey))
				b, err := DeleteDevice(p.(map[string]interface{}))
				if err != nil {
					continue
				}

				respCh, ok := t.Ctx.Value(managerKey(waitResponseKey)).(chan []byte)
				if !ok {
					return
				}
				respCh <- b
				// p := t.Ctx.Value(managerKey(parameterKey))

			case STATUSREPORT:
				p := t.Ctx.Value(managerKey(parameterKey))
				params, _ := p.(map[string]interface{})["params"].(map[string]interface{})
				device, _ := p.(map[string]interface{})["device"].(*model.Device)

				respCh, ok := t.Ctx.Value(managerKey(waitResponseKey)).(chan []byte)
				if !ok {
					return
				}

				b, err := ForwardMessage(device.DID, device.SID, params)
				if err != nil {
					respCh <- []byte(err.Error())
				}

				respCh <- b
			}
		}
	}
}
