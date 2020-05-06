package serviceCommunicatorServer

// import (
// 	`context`
// 	`fmt`
// 	`net`
// 	`time`
//
//
//

// 	GoEnvTools "github.com/kaizer666/goenvtools"
// 	GoLogger "github.com/kaizer666/gologger"
// )

// type UdpServerStruct struct {
// 	logger       *GoLogger.Logger
// 	environment  *GoEnvTools.GoEnv
// 	address      string
// 	loggerLoaded bool
// 	ctx context.Context
// 	maxBufferSize int
// }

// func (udpServer *UdpServerStruct) SetAddress(address string) {
// 	udpServer.address = address
// }
// func (udpServer *UdpServerStruct) SetContext(ctx context.Context) {
// 	udpServer.ctx = ctx
// }

// func (udpServer *UdpServerStruct) SetEnvironment(env *GoEnvTools.GoEnv) {
// 	udpServer.environment = env
// }
// func (udpServer *UdpServerStruct) SetLogger(logger *GoLogger.Logger) {
// 	udpServer.logger = logger
// 	udpServer.loggerLoaded = true
// }

// func (udpServer *UdpServerStruct) StartServer() {
// 	udpServer.maxBufferSize = 1024
// 	pc, err := net.ListenPacket("udp", udpServer.address)
// 	if err != nil {
// 		return
// 	}

// 	// `Close`ing the packet "connection" means cleaning the data structures
// 	// allocated for holding information about the listening socket.
// 	defer pc.Close()

// 	doneChan := make(chan error, 1)
// 	buffer := make([]byte, udpServer.maxBufferSize)

// 	// Given that waiting for packets to arrive is blocking by nature and we want
// 	// to be able of canceling such action if desired, we do that in a separate
// 	// go routine.
// 	go func() {
// 		for {
// 			// By reading from the connection into the buffer, we block until there's
// 			// new content in the socket that we're listening for new packets.
// 			//
// 			// Whenever new packets arrive, `buffer` gets filled and we can continue
// 			// the execution.
// 			//
// 			// note.: `buffer` is not being reset between runs.
// 			//	  It's expected that only `n` reads are read from it whenever
// 			//	  inspecting its contents.
// 			n, addr, err := pc.ReadFrom(buffer)
// 			if err != nil {
// 				doneChan <- err
// 				return
// 			}

// 			fmt.Printf("packet-received: bytes=%d from=%s\n",
// 				n, addr.String())

// 			// Setting a deadline for the `write` operation allows us to not block
// 			// for longer than a specific timeout.
// 			//
// 			// In the case of a write operation, that'd mean waiting for the send
// 			// queue to be freed enough so that we are able to proceed.
// 			deadline := time.Now().Add(time.Second)
// 			err = pc.SetWriteDeadline(deadline)
// 			if err != nil {
// 				doneChan <- err
// 				return
// 			}

// 			// Write the packet's contents back to the client.
// 			n, err = pc.WriteTo(buffer[:n], addr)
// 			if err != nil {
// 				doneChan <- err
// 				return
// 			}

// 			fmt.Printf("packet-written: bytes=%d to=%s\n", n, addr.String())
// 		}
// 	}()

// 	select {
// 	case <-udpServer.ctx.Done():
// 		fmt.Println("cancelled")
// 		err = udpServer.ctx.Err()
// 	case err = <-doneChan:
// 	}

// 	return
// }
