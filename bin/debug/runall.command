cd "$(dirname "$0")"
for file in `ls|grep command`
do
    if [[ $file != "runall.command" ]]
    then
        open "$file"
    fi
done