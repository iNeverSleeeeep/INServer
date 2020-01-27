ps ef|grep "main -id"|grep -v "grep"|awk '{print $1}'|xargs kill -2
