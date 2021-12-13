package devicemanager

import (
	"context"
	"encoding/json"
	"etrismartfarmpoccontroller/model"
	"fmt"
	"log"

	manager "git.godopu.net/lab/etri-smartfarm-poc-controller-serial"
)

func NewRecvHandler() manager.EventHandler {
	return manager.NewEventHandler(func(e manager.Event) {
		fmt.Println("RECV: ", e.Params())
		parmas := e.Params()
		payload := map[string]interface{}{}

		var ok bool
		payload["uuid"], ok = parmas["uuid"]
		if !ok {
			return
		}

		db, err := model.GetDBHandler("sqlite", "./dump.db")
		if err != nil {
			return
		}

		device, err := db.GetDeviceID(payload["uuid"].(string))
		if err != nil {
			log.Println(err)
			return
		}

		device.DName = payload["uuid"].(string)
		device.SID, err = db.GetSID(device.SName)
		if err != nil {
			log.Println(err)
			return
		}

		parameter := map[string]interface{}{
			"params": parmas,
			"device": device,
		}
		respCh := make(chan []byte)
		ctx := context.WithValue(context.Background(), managerKey(parameterKey), parameter)
		ctx = context.WithValue(ctx, managerKey(waitResponseKey), respCh)
		taskQueue <- &task{Event: STATUSREPORT, Ctx: ctx}

		resp := map[string]interface{}{}
		b := <-respCh
		err = json.Unmarshal(b, &resp)

		if err != nil {
			fmt.Println(string(b))
		} else {

			changed := db.StatusCheck(device.DID, resp)

			if changed {
				fmt.Println("change property: ", resp)
				manager.Sync(device.DName, resp)
			}
		}
		// ctx, cancel := context.WithCancel(ctx)
	})
}

func RemovedHandler(e manager.Event) {
	param := e.Params()
	payload := map[string]interface{}{}
	payload["dname"] = param["uuid"]
	respCh := make(chan []byte)
	ctx := context.WithValue(context.Background(), managerKey(parameterKey), payload)
	ctx = context.WithValue(ctx, managerKey(waitResponseKey), respCh)
	taskQueue <- &task{Event: DISCONNECTED, Ctx: ctx}

	fmt.Println(string(<-respCh))
}
