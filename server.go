package serviceCommunicatorServer

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"

	GoEnvTools "github.com/kaizer666/goenvtools"
	"github.com/kaizer666/gologger"
)

var locker sync.Mutex

type ServerStruct struct {
	Commands       []CommandStruct
	logger         *gologger.Logger
	environment    *GoEnvTools.GoEnv
	handlers       map[string]func(http.ResponseWriter, *http.Request)
	loggerLoaded   bool
	FileDescriptor *int
	ExitListener   *chan int
	listener       net.Listener
	StopChannels   []*chan int
	StopFunctions  []func()
	isStopped      bool
	file           *os.File
	address        string
}

func (mainServer *ServerStruct) SetAddress(address string) {
	mainServer.address = address
}

func (mainServer *ServerStruct) SetEnvironment(env *GoEnvTools.GoEnv) {
	mainServer.environment = env
}
func (mainServer *ServerStruct) SetLogger(logger *gologger.Logger) {
	mainServer.logger = logger
	mainServer.loggerLoaded = true
}

func (mainServer *ServerStruct) GetCommands(w http.ResponseWriter, _ *http.Request) {
	data, err := json.Marshal(mainServer.Commands)
	if err != nil {
		if mainServer.loggerLoaded {
			mainServer.logger.Error("error: %v", err)
		}
		_, _ = io.WriteString(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()))
		return
	}
	_, _ = io.WriteString(w, string(data))
}

func (mainServer *ServerStruct) SetHandlers(handlers map[string]func(http.ResponseWriter, *http.Request)) {
	mainServer.handlers = make(map[string]func(http.ResponseWriter, *http.Request))
	for path, handler := range handlers {
		mainServer.handlers[path] = handler
	}
	mainServer.handlers["/getCommands"] = mainServer.GetCommands
}

func (mainServer *ServerStruct) StartServer() {
	locker.Lock()
	mainServer.isStopped = false
	s := &http.Server{
		Addr:         mainServer.address,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	var err error
	if *mainServer.FileDescriptor != 0 {
		if mainServer.loggerLoaded {
			mainServer.logger.Debug("Starting with fileDescriptor %v", *mainServer.FileDescriptor)
		}
		mainServer.file = os.NewFile(uintptr(*mainServer.FileDescriptor), "parent socket")
		mainServer.listener, err = net.FileListener(mainServer.file)
		if err != nil {
			if mainServer.loggerLoaded {
				mainServer.logger.Error("fd listener failed: %v", err)
			}
		}
	} else {
		if mainServer.loggerLoaded {
			mainServer.logger.Debug("Virgin start")
		}
		mainServer.listener, err = net.Listen("tcp", s.Addr)
		if err != nil {
			if mainServer.loggerLoaded {
				mainServer.logger.Error("listener failed: %v", err)
			}
		}
	}
	for path, handler := range mainServer.handlers {
		http.HandleFunc(path, handler)
	}
	locker.Unlock()
	err = s.Serve(mainServer.listener)
	if err != nil {
		if mainServer.loggerLoaded {
			_ = mainServer.logger.Critical("error: %v", err)
		}
	}
	<-*mainServer.ExitListener
}

func (mainServer *ServerStruct) GraceStop() {
	locker.Lock()
	for _, function := range mainServer.StopFunctions {
		function()
	}
	if mainServer.loggerLoaded {
		mainServer.loggerLoaded = false
		_ = mainServer.logger.Close()
	}
	if mainServer.isStopped {
		locker.Lock()
		return
	}
	mainServer.isStopped = true
	defer func() { fmt.Println("GoodBye") }()
	listener2 := mainServer.listener.(*net.TCPListener)
	file2, err := listener2.File()
	if err != nil {
		locker.Lock()
		panic(err)
	}
	fd1 := int(file2.Fd())
	_, err = syscall.Dup(fd1)
	if err != nil {
		fmt.Println("dup error: ", err)
	}
	err = mainServer.listener.Close()
	if err != nil {
		fmt.Println(err)
	}
	if mainServer.file != nil {
		err := mainServer.file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}

	for _, channel := range mainServer.StopChannels {
		*channel <- 1
		time.Sleep(100 * time.Millisecond)
	}
	locker.Unlock()
	time.Sleep(2 * time.Second)
	*mainServer.ExitListener <- 1
}

func (mainServer *ServerStruct) GraceHandler() {
	for _, function := range mainServer.StopFunctions {
		function()
	}
	programName := os.Args[0]
	args := os.Args[1:]
	var cleanArgs []string
	for _, arg := range args {
		if arg[0:4] != "-fd=" {
			cleanArgs = append(cleanArgs, arg)
		}
	}
	if mainServer.loggerLoaded {
		mainServer.loggerLoaded = false
		_ = mainServer.logger.Close()
	}
	if mainServer.isStopped {
		return
	}
	mainServer.isStopped = true
	listener2 := mainServer.listener.(*net.TCPListener)
	file2, err := listener2.File()
	if err != nil {
		panic(err)
	}
	fd1 := int(file2.Fd())
	fd2, err := syscall.Dup(fd1)
	if err != nil {
		fmt.Println("dup error: ", err)
	}
	err = mainServer.listener.Close()
	if err != nil {
		fmt.Println(err)
	}
	if mainServer.file != nil {
		err = mainServer.file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}

	for _, channel := range mainServer.StopChannels {
		*channel <- 1
		time.Sleep(100 * time.Millisecond)
	}

	cleanArgs = append(cleanArgs, fmt.Sprint("-fd=", fd2))
	e := GoEnvTools.GoEnv{}
	_ = e.InitEnv()
	cmd := exec.Command(programName, cleanArgs...)
	cmd.Env = append(cmd.Env, os.Environ()...)
	err = cmd.Run()
	if err != nil {
		panic(fmt.Sprintf("grace starting error: %s", err))
	}
	time.Sleep(2 * time.Second)
	*mainServer.ExitListener <- 1
}
