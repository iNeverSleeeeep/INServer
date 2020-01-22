cd `dirname $0`
open 0-center.command
sleep 2s
open 3-database.command
sleep 2s
for file in `ls|grep command`
do
    if [[ $file != "runall.command" && $file != "0-center.command" && $file != "3-database.command" ]]
    then
        open "$file"
    fi
done