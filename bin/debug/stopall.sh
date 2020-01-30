ps aux|grep "main -id"|grep -v "grep"|awk '{print $2}'|xargs kill -2
