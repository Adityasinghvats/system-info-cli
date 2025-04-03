# Golang System Profiler
## This project using Tview and Gopsutil to show you the system information of your machine.
- This tool is made to work like the `htop` utility of linux and due to fact that it is built in Golang, it can be used across various platform.
### How to use ?
- make directory 
```bash
  mkdir systeminfo
  cd systeminfo
```
- get the code
```bash
git clone https://github.com/Adityasinghvats/system-info-cli.git
```
- get dependencies
```bash
go mod tidy
```
- run the code
```bash
go build main.go
./main
```

