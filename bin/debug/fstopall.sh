ps ef|grep "@in-"|grep -v "grep"|awk '{print $1}'|xargs kill -9
