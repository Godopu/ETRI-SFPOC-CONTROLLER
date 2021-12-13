package main

import (
	"bytes"
	"encoding/json"
	"etrismartfarmpoccontroller/constants"
	"etrismartfarmpoccontroller/devicemanager"
	"etrismartfarmpoccontroller/router"
	"fmt"
	"io/ioutil"
	"net/http"

	manager "git.godopu.net/lab/etri-smartfarm-poc-controller-serial"
)

// func runBootstrap() {
// 	l, err := net.NewListenUDP("udp4", "", net.WithHeartBeat(time.Second*5))
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	defer l.Close()

// 	init := false
// 	minTimeout := time.Second * 5
// 	timeout := minTimeout

// 	d := func() {
// 		s := udp.NewServer(udp.WithTransmission(time.Second, timeout/2, 2))
// 		var wg sync.WaitGroup
// 		defer wg.Wait()
// 		defer s.Stop()
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()
// 			s.Serve(l)
// 		}()

// 		ctx, cancel := context.WithTimeout(context.Background(), timeout)
// 		defer cancel()

// 		req, err := client.NewGetRequest(ctx, "/bs") /* msg.Option{
// 			ID:    msg.URIQuery,
// 			Value: []byte("rt=oic.wk.d"),
// 		}*/
// 		if err != nil {
// 			panic(fmt.Errorf("cannot create discover request: %w", err))
// 		}

// 		req.SetMessageID(udpmessage.GetMID())
// 		req.SetType(udpmessage.NonConfirmable)
// 		defer pool.ReleaseMessage(req)

// 		err = s.DiscoveryRequest(req, bootstrapAddr, func(cc *client.ClientConn, resp *pool.Message) {
// 			b, err := ioutil.ReadAll(resp.Body())
// 			if err != nil {
// 				panic(err)
// 			}
// 			fmt.Println(string(b))
// 			init = true
// 		})
// 		if err != nil {
// 			panic(err)
// 		}
// 	}
// 	for {
// 		d()
// 		if init {
// 			break
// 		}
// 	}
// }

func register() error {
	// Controller 이름을 읽어옴
	payload := map[string]string{}
	payload["cname"] = constants.Config["cname"]
	fmt.Println(constants.Config["cname"])
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// Controller 등록 메시지 송신
	resp, err := http.Post(
		fmt.Sprintf("http://%s/%s", constants.Config["serverAddr"], "controllers"),
		"application/json",
		bytes.NewReader(b),
	)

	if err != nil {
		return err
	}

	// 응답 메시지 수신
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	json.Unmarshal(b, &payload)

	// 등록 후 생성된 Controller ID 저장
	constants.Set("cid", payload["cid"])

	return nil
}

func main() {

	if constants.Config["cid"] == "" {
		err := register()
		if err != nil {
			panic(err)
		}
	}

	run, _ := devicemanager.NewManager()
	run()

	manager.AddRecvListener(devicemanager.NewRecvHandler())
	// handler := devicemanager.NewRecvHandler()
	manager.AddRegisterHandleFunc(devicemanager.RegisterHandler)
	manager.AddRemoveHandleFunc(devicemanager.RemovedHandler)

	go manager.Run()

	http.ListenAndServe(":4000", router.NewRouter())

	// for {
	// 	fmt.Println("> ")
	// 	var cmd string
	// 	fmt.Scanln(&cmd)
	// 	if cmd == "exit" {
	// 		return
	// 	}

	// 	handler.Handle(&Temp{})
	// }
}

// type Temp struct{}

// func (*Temp) Key() interface{} {
// 	return &struct{}{}
// }

// func (*Temp) Params() map[string]interface{} {
// 	return map[string]interface{}{
// 		"uuid": "DEVICE-A-UUID",
// 	}
// }
