cd `dirname $0`
open 0-center.command
sleep 2s
for file in `ls|grep command`
do
    if [[ $file != "runall.command" ]]
    then
        open "$file"
    fi
done