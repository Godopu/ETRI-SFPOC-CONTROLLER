package devicemanager

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"etrismartfarmpoccontroller/constants"
	"etrismartfarmpoccontroller/model"
	"io/ioutil"
	"net/http"
)

func RegisterDevice(payload *map[string]interface{}, cancelCh <-chan struct{}) ([]byte, error) {
	(*payload)["cid"] = constants.Config["cid"]
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		"http://"+constants.Config["serverAddr"]+"/devices",
		bytes.NewReader(b),
	)
	req.Header.Set("Content-Type", "application/json")
	var resp *http.Response
	done := make(chan bool)
	go func() {
		resp, err = http.DefaultClient.Do(req)
		if resp == nil {
			return
		}
		done <- true
	}()

	select {
	case <-done:
		if err != nil {
			return nil, err
		}
	case <-cancelCh:
		return nil, errors.New("cancel error")
	}

	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	db, err := model.GetDBHandler("sqlite", "./dump.db")
	if err != nil {
		return nil, err
	}

	var device model.Device
	json.Unmarshal(b, &device)
	err = db.AddDevice(&device)
	if err != nil {
		return nil, err
	}
	// db := model.GetDBHandler("sqlite", "./dump.db")
	return b, nil
}
