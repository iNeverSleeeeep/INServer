cd "$(dirname "$0")"
go build -o ../../inserver-temp.exe ../../main.go
cd ../..
./inserver-temp.exe -id 0 &
./inserver-temp.exe -id 1 &
./inserver-temp.exe -id 2 &
./inserver-temp.exe -id 3 &
./inserver-temp.exe -id 4 &
./inserver-temp.exe -id 5 &
./inserver-temp.exe -id 6 &
./inserver-temp.exe -id 7 &
./inserver-temp.exe -id 8 &
./inserver-temp.exe -id 9 &
