package web

import (
	"INServer/src/common/global"
	"INServer/src/common/logger"
	"INServer/src/proto/msg"
	"INServer/src/services/etcmgr"
	"INServer/src/services/node"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
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
	dir, err := os.Getwd()
	if err != nil {
		logger.Fatal(err)
	}
	certFile := fmt.Sprintf("%s/web/server.crt", dir)
	keyFile := fmt.Sprintf("%s/web/server.key", dir)
	go http.ListenAndServeTLS(":"+strconv.Itoa(int(global.CurrentServerConfig.WebConfig.Port)), certFile, keyFile, nil)
	go http.ListenAndServe(":"+strconv.Itoa(int(global.CurrentServerConfig.WebConfig.Port)), nil)
}

func (w *Web) checkauth(writer http.ResponseWriter, req *http.Request) bool {
	auth := req.Header.Get("Authorization")
	if auth == "" {
		writer.Header().Set("WWW-Authenticate", `Basic realm="Dotcoo User Login"`)
		writer.WriteHeader(http.StatusUnauthorized)
		return false
	}
	auths := strings.SplitN(auth, " ", 2)
	if len(auths) != 2 {
		fmt.Println("error")
		return false
	}
	authMethod := auths[0]
	authB64 := auths[1]
	switch authMethod {
	case "Basic":
		authstr, err := base64.StdEncoding.DecodeString(authB64)
		if err != nil {
			logger.Error(err)
			io.WriteString(writer, "Unauthorized!\n")
			return false
		}
		user := strings.SplitN(string(authstr), ":", 2)
		if len(user) != 2 {
			logger.Error(user)
			return false
		}
		account := user[0]
		password := user[1]
		config := global.CurrentServerConfig.WebConfig
		if account != config.Account || password != config.Password {
			io.WriteString(writer, "账号或密码错误!\n")
			return false
		}
	default:
		logger.Error(authMethod)
		return false
	}
	return true
}

func (w *Web) root(writer http.ResponseWriter, req *http.Request) {
	logger.Info("HTTP请求", req.RemoteAddr)
	if w.checkauth(writer, req) == false {
		return
	}
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
