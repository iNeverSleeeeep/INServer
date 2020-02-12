cd ../proto

# data
echo "--------- data ---------"
for file in `ls data` 
do
echo $file
protoc --csharp_out=../clientproto  --proto_path=./data --proto_path=./engine $file
done

# msg
echo "--------- msg ---------"
for file in `ls msg` 
do
echo $file
protoc --csharp_out=../clientproto --proto_path=./msg --proto_path=./data --proto_path=./etc --proto_path=./engine $file
done

# msg
echo "--------- engine ---------"
for file in `ls engine` 
do
echo $file
protoc --csharp_out=../clientproto --proto_path=./engine $file
done

# msg
echo "--------- etc ---------"
for file in `ls etc` 
do
echo $file
protoc --csharp_out=../clientproto --proto_path=./etc $file
done
