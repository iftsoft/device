
mkdir build 2> NUL
mkdir build\windows 2> NUL

set GOARCH=amd64
set GOOS=windows

cd example\client
go build -o ..\..\build\windows\client.exe .
copy client.yml ..\..\build\windows\client.yml
cd ..\..

cd example\server
go build -o ..\..\build\windows\server.exe .
copy server.yml ..\..\build\windows\server.yml
cd ..\..

cd example\cashvalidator
go build -o ..\..\build\windows\cashvalidator.exe .
copy cashvalidator.yml ..\..\build\windows\cashvalidator.yml
cd ..\..
