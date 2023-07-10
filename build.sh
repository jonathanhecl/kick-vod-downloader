GOOS=windows GOARCH=amd64 go build -ldflags "-w -s" -o kvd-mac
GOOS=linux GOARCH=amd64 go build -ldflags "-w -s" -o kvd-linux
GOOS=windows GOARCH=amd64 go build -ldflags "-w -s" -o kvd-win.exe