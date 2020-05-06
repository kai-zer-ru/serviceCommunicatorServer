package serviceCommunicatorServer

import (
	"fmt"
	"os"
	"testing"
	"time"

	GoEnvTools "github.com/kaizer666/goenvtools"
	"github.com/kaizer666/gologger"
)

var (
	exitListener = make(chan int)
	waitListener = make(chan int)
)

func connect() (ServerStruct, error) {
	os.Create(".env")
	defer os.Remove(".env")
	s := ServerStruct{}
	s.ExitListener = &exitListener
	s.SetAddress(":22222")
	env := GoEnvTools.GoEnv{}
	err := env.InitEnv()
	if err != nil {
		return s, err
	}
	s.SetEnvironment(&env)
	fd := 0
	s.FileDescriptor = &fd
	logger := gologger.Logger{}
	logger.SetLogFileName("main.log")
	logger.SetLogLevel(1)
	err = logger.Init()
	if err != nil {
		return s, err
	}
	s.SetLogger(&logger)
	return s, nil
}

func TestMain(t *testing.T) {
	s, err := connect()
	if err != nil {
		t.Error(err)
	}
	s.StopFunctions = append(s.StopFunctions, stopFunc1)
	s.StopChannels = append(s.StopChannels, &waitListener)
	go waitFunc1()
	go s.StartServer()
	time.Sleep(10 * time.Second)
	s.GraceStop()
}

func stopFunc1() {
	time.Sleep(time.Second)
}

func waitFunc1() {
	<-waitListener
	fmt.Println("exit")
}

func TestRemoveLogFile(_ *testing.T) {
	os.Remove("main.log")
}
