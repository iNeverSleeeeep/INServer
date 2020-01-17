cd ../proto

# data
echo "--------- data ---------"
cd  ./data
dir=`ls`
for file in $dir 
do
echo $file
protoc --csharp_out=../../clientproto  --proto_path=. $file
done
cd ..

# msg
echo "--------- msg ---------"
cd  ./msg
dir=`ls`
for file in $dir 
do
echo $file
cd ..
protoc --csharp_out=../clientproto --proto_path=./msg --proto_path=./data --proto_path=./etc $file
cd ./msg
done
cd ..

sleep 1s
