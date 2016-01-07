package gosignaler

import (
	"os"
	"sync"
	"syscall"
	"testing"
	"time"
)

type testSignalReceiver struct {
	wg        *sync.WaitGroup
	ChkSignal os.Signal
}

func (t *testSignalReceiver) Receive(signal os.Signal) (proceed bool) {
	t.ChkSignal = signal
	return
}

func (t *testSignalReceiver) WaitGroup() (ret *sync.WaitGroup) {
	return
}

func (t *testSignalReceiver) Interests() (ret []os.Signal) {
	return []os.Signal{syscall.SIGHUP}
}

//test listen method
func TestListen(t *testing.T) {
	testReceiver := &testSignalReceiver{}

	Listen(testReceiver)
	select {
	case <-time.After(time.Second):
		syscall.Kill(syscall.Getpid(), syscall.SIGHUP)
	}

	select {
	case <-time.After(time.Second):
		if testReceiver.ChkSignal != syscall.SIGHUP {
			t.Fatalf("unexpected signal: %s", testReceiver.ChkSignal.String())
		}
	}
}
