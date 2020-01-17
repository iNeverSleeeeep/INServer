package web

import (
	"INServer/src/common/global"
	"INServer/src/common/logger"
	"INServer/src/modules/node"
	"INServer/src/proto/msg"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

var Instance *Web

type (
	Web struct {
	}
)

func New() *Web {
	w := new(Web)
	return w
}

func (w *Web) Start() {
	http.HandleFunc("/zones", w.zones)
	http.HandleFunc("/reloadetc", w.reloadetc)
	go http.ListenAndServe(":"+strconv.Itoa(int(global.ServerConfig.WebConfig.Port)), nil)
}

func (w *Web) zones(writer http.ResponseWriter, req *http.Request) {
	bytes, err := json.Marshal(global.Zones)
	if err != nil {
		logger.Debug(err)
	} else {
		io.WriteString(writer, string(bytes))
	}
}

func (w *Web) reloadetc(writer http.ResponseWriter, req *http.Request) {
	_, err := node.Instance.Net.RequestBytes(msg.Command_RELOAD_ETC_REQ, make([]byte, 1))
	if err != nil {
		io.WriteString(writer, err.Error())
	} else {
		io.WriteString(writer, "success")
	}
}
