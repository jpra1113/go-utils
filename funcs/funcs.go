package funcs

import (
	"errors"
	"time"
)

func LoopUntil(timeout time.Duration, stepTime time.Duration, loopFunc func() (bool, error)) error {
	doneChan := make(chan bool, 1)
	errChan := make(chan error, 1)
	quitChan := make(chan bool, 1)

	go func() {
		for {
			select {
			case <-quitChan:
				return
			default:
				done, err := loopFunc()
				if err != nil {
					errChan <- err
					return
				} else if done {
					doneChan <- true
					return
				}
				time.Sleep(stepTime)
			}
		}
	}()

	select {
	case <-doneChan:
		return nil
	case err := <-errChan:
		return err
	case <-time.After(timeout):
		quitChan <- true
		return errors.New("Time out waiting for loop to terminate")
	}
}
