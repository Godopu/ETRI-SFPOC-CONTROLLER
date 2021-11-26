package router

import (
	"etrismartfarmpoccontroller/constants"
	"etrismartfarmpoccontroller/devicemanager"
	"etrismartfarmpoccontroller/model"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

var dbHandler model.DBHandler

func init() {
	var err error
	dbHandler, err = model.GetDBHandler("sqlite", "./dump.db")

	if err != nil {
		panic(err)
	}

}
func NewRouter() http.Handler {
	m := mux.NewRouter()

	m.HandleFunc("/echo", Echo).Methods("GET")
	m.HandleFunc("/devices", devicemanager.PostDevice).Methods("POST")

	// sub := mux.NewRouter()
	// sub.PathPrefix("/{did}/").HandlerFunc(EchoPath)
	m.PathPrefix("/devices/{did}/").HandlerFunc(EchoPath)

	n := negroni.Classic() // 파일 서버 및 로그기능을 제공함
	n.UseHandler(m)
	return n
}

func EchoPath(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	path := r.URL.Path[(len([]rune("/devices"))+len([]rune(vars["did"])))+1:]

	sname, err := dbHandler.GetServiceForDevice(vars["did"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sid, err := dbHandler.GetSID(sname)
	if err != nil || len(sid) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// w.Write([]byte("http://" + constants.Config["serverAddr"] + "/services/" + sid + path))
	req, err := http.NewRequest(
		r.Method,
		"http://"+constants.Config["serverAddr"]+"/services/"+sid+path,
		r.Body,
	)

	if err != nil {
		// 잘못된 메시지 포맷이 전달된 경우
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		// 잘못된 메시지 포맷이 전달된 경우
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}

	io.Copy(w, resp.Body)
}

func Forward(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	sname, err := dbHandler.GetServiceForDevice(vars["did"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sid, err := dbHandler.GetSID(sname)
	if err != nil || len(sid) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Write([]byte(sid))
}

func Echo(w http.ResponseWriter, r *http.Request) {
	io.Copy(w, r.Body)
}
