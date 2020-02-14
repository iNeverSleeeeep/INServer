package web

import (
	"INServer/src/common/global"
	"INServer/src/common/logger"
	"INServer/src/proto/msg"
	"INServer/src/services/etcmgr"
	"INServer/src/services/node"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"
)

//定义全局的模板变量
var webhtml *template.Template
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
	http.HandleFunc("/", w.root)
	http.HandleFunc("/zones", w.zones)
	http.HandleFunc("/reloadetc", w.reloadetc)
	go http.ListenAndServe(":"+strconv.Itoa(int(global.CurrentServerConfig.WebConfig.Port)), nil)
}

func (w *Web) root(writer http.ResponseWriter, req *http.Request) {
	err := webhtml.Execute(writer, nil)
	if err != nil {
		logger.Error(err)
	}
}

func (w *Web) zones(writer http.ResponseWriter, req *http.Request) {
	bytes, err := json.Marshal(etcmgr.Instance.Zones())
	if err != nil {
		logger.Debug(err)
	} else {
		io.WriteString(writer, string(bytes))
	}
}

func (w *Web) reloadetc(writer http.ResponseWriter, req *http.Request) {
	_, err := node.Net.RequestBytes(msg.CMD_RELOAD_ETC_REQ, make([]byte, 1))
	if err != nil {
		io.WriteString(writer, err.Error())
	} else {
		io.WriteString(writer, "success")
	}
}

func init() {
	dir, err := os.Getwd()
	if err != nil {
		logger.Fatal(err)
	}
	webhtml, err = template.ParseFiles(fmt.Sprintf("%s/web/web.html", dir))
	if err != nil {
		logger.Fatal(err)
	}
}
