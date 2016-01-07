//package gosignaler provides signal handler wrapper
package gosignaler

import (
	"errors"
	"os"
	"os/signal"
	"sync"
	"time"
)

const (
	SignalBufferSize = 32
)

var (
	ErrNoReceiver      = errors.New("no signal receiver") // no receiver
	ErrInvalidReceiver = errors.New("invalid receiver")   //invalid receiver
)

//signal receiver to receive signal, and provide wait group
//for graceful shutdown
type SignalReceiver interface {
	//receive signal and return whether to proceed next signal
	Receive(signal os.Signal) (proceed bool)
	//wait group for graceful shutdown
	WaitGroup() (wg *sync.WaitGroup)
	//list of signal to listen
	Interests() (signals []os.Signal)
}

//listen for signal, and pass signal to current callback
//if error occurred, return error
func Listen(receiver SignalReceiver) (ret error) {
	if nil == receiver ||
		0 == len(receiver.Interests()) {
		ret = ErrInvalidReceiver
		return
	}

	signalChannel := make(chan os.Signal, SignalBufferSize)
	signal.Notify(signalChannel, receiver.Interests()...)
	if nil != receiver.WaitGroup() {
		receiver.WaitGroup().Add(1)
	}

	go listenSignal(receiver, signalChannel)
	return
}

func listenSignal(receiver SignalReceiver, signalChannel chan os.Signal) {
	for {
		select {
		case s := <-signalChannel:
			proceed := receiver.Receive(s)
			if !proceed {
				break
			}
		case <-time.After(time.Second):
		}
	}

	if nil != receiver.WaitGroup() {
		receiver.WaitGroup().Done()
	}
}
