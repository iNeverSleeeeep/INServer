ps aux|grep "@in-"|grep -v "grep"|awk '{print $2}'|xargs kill -2
