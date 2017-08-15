package funcs

import (
	"encoding/json"
	"errors"
	"fmt"
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

func DeepCopy(from interface{}, to interface{}) error {
	if from == nil {
		return errors.New("Unable to find 'from' interface")
	}
	if to == nil {
		return errors.New("Unable to find 'to' interface")
	}
	bytes, err := json.Marshal(from)
	if err != nil {
		return fmt.Errorf("Unable to marshal src: %s", err)
	}
	err = json.Unmarshal(bytes, to)
	if err != nil {
		return fmt.Errorf("Unable to unmarshal into dst: %s", err)
	}
	return nil
}
