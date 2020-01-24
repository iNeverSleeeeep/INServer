ps aux|grep "go run main.go"|grep -v "grep"|awk '{print $2}'|xargs kill -9
ps aux|grep "main -id"|grep -v "grep"|awk '{print $2}'|xargs kill -1
