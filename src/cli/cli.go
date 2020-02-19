package cli

import (
	"INServer/src/common/global"
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
)

// Run 启动交互命令行
func Run() {
	reader := bufio.NewReader(os.Stdin)

	http.HandleFunc("/cli", handleHTTPCli)
	go http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", 8000+global.CurrentServerID), nil)
	for {
		fmt.Print("$ ")
		cmdString, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		err = runCommand(cmdString)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func handleHTTPCli(w http.ResponseWriter, r *http.Request) {
	err := runCommand(r.FormValue("cmd"))
	if err != nil {
		w.Write([]byte(fmt.Sprintf("%s", err)))
	} else {
		w.Write([]byte("Success"))
	}
}

func runCommand(commandStr string) error {
	commandStr = strings.TrimSuffix(commandStr, "\n")
	arrCommandStr := strings.Fields(commandStr)
	switch arrCommandStr[0] {
	case "exit":
		global.Exit <- true
		break
	case "ping":
		fmt.Println("pong")
		break
	}
	return nil
}
