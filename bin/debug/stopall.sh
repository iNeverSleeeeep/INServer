ps aux|grep "go run main.go"|grep -v "grep"|awk '{print $2}'|xargs kill -9
