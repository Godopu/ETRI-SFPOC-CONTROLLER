package devicemanager

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
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
	ctx := context.WithValue(context.Background(), managerKey(parameterKey), &parameter)
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
