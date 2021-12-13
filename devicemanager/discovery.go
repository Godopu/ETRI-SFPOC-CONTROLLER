package devicemanager

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	manager "git.godopu.net/lab/etri-smartfarm-poc-controller-serial"
)

func PostDevice(w http.ResponseWriter, r *http.Request) {
	defer log.Println("Post Device End")
	// 장치로 부터 전달된 데이터 처리
	parameter := map[string]interface{}{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&parameter)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	// 탐색 이벤트를 처리 큐에 입력
	respCh := make(chan []byte)
	ctx := context.WithValue(context.Background(), managerKey(parameterKey), parameter)
	ctx = context.WithValue(ctx, managerKey(waitResponseKey), respCh)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	taskQueue <- &task{Event: DISCOVERY, Ctx: ctx}
	select {
	case resp := <-respCh:
		if len(resp) == 0 {
			w.WriteHeader(http.StatusBadRequest)
		}
		w.WriteHeader(http.StatusCreated)
		w.Write(resp)
	case <-r.Context().Done():
		w.WriteHeader(http.StatusNotAcceptable)
	}
}

func RegisterHandler(e manager.Event) {
	param := e.Params()
	payload := map[string]interface{}{}
	payload["dname"] = param["uuid"]
	payload["type"] = "device"
	payload["sname"] = param["sname"]

	// b, err := json.Marshal(payload)
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	// resp, err := http.Post("http://localhost:4000/devices", "application/json", bytes.NewReader(b))
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }

	respCh := make(chan []byte)
	ctx := context.WithValue(context.Background(), managerKey(parameterKey), payload)
	ctx = context.WithValue(ctx, managerKey(waitResponseKey), respCh)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	taskQueue <- &task{Event: DISCOVERY, Ctx: ctx}

	fmt.Println(string(<-respCh))
}
