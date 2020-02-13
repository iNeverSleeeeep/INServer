cd ../proto

# engine
echo "--------- engine ---------"
for file in `ls engine`
do
echo $file
protoc --gofast_out=../../  --go-json_out=../src/proto/engine --proto_path=./engine $file
done

# config
echo "--------- config ---------"
for file in `ls config`
do
echo $file
protoc --gofast_out=../../  --go-json_out=../src/proto/config --proto_path=./config --proto_path=./engine $file
done

# data
echo "--------- data ---------"
for file in `ls data`
do
echo $file
protoc --gofast_out=../../  --go-json_out=../src/proto/data --proto_path=./data --proto_path=./engine $file
done

# msg
echo "--------- msg ---------"
for file in `ls msg`
do
echo $file
protoc --gofast_out=../../ --proto_path=./msg --proto_path=./data --proto_path=./etc --proto_path=./engine  --proto_path=./gogoproto --proto_path=./protobuf $file
done

# etc
echo "--------- etc ---------"
for file in `ls etc`
do
echo $file
protoc --gofast_out=../../  --go-json_out=../src/proto/etc --proto_path=./etc --proto_path=./gogoproto --proto_path=./protobuf $file
done

# db
echo "--------- db ---------"
for file in `ls db` 
do
echo $file
protoc --gofast_out=../../ --proto_path=./db --proto_path=./gogoproto --proto_path=./protobuf $file
done
