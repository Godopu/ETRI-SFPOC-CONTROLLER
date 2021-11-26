package model

import (
	"bytes"
	"encoding/json"
	"etrismartfarmpoccontroller/constants"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (s *dbHandler) GetSID(sname string) (string, error) {
	sid, ok := s.cache[sname]
	if !ok {
		payload := map[string]string{"sname": sname}
		b, _ := json.Marshal(payload)

		req, err := http.NewRequest("GET",
			fmt.Sprintf("http://%s/%s", constants.Config["serverAddr"], "services"),
			bytes.NewReader(b),
		)
		if err != nil {
			return "", err
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return "", err
		}

		b, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		sid = string(b)
	}

	return sid, nil
}
