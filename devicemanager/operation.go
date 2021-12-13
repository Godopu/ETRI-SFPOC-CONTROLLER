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

func RegisterDevice(payload map[string]interface{}, cancelCh <-chan struct{}) ([]byte, error) {
	payload["cid"] = constants.Config["cid"]
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

	if resp.StatusCode == http.StatusCreated {
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
	}

	// db := model.GetDBHandler("sqlite", "./dump.db")
	return b, nil
}

func DeleteDevice(payload map[string]interface{}) ([]byte, error) {
	payload["cid"] = constants.Config["cid"]
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		"DELETE",
		"http://"+constants.Config["serverAddr"]+"/devices",
		bytes.NewReader(b),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func ForwardMessage(did, sid string, payload map[string]interface{}) ([]byte, error) {

	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("PUT", "http://localhost:3000/services/"+sid+"/"+did, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}
