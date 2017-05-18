package funcs

import (
	"errors"
	"testing"
	"time"
)

func TestLoopUntil(t *testing.T) {
	count := 0
	err := LoopUntil(time.Second*10, time.Millisecond*100, func() (bool, error) {
		count += 1
		if count < 5 {
			return false, nil
		}

		return true, nil
	})

	if err != nil {
		t.Error("Unable to finish loop before timeout: " + err.Error())
	}

	err = LoopUntil(time.Second*10, time.Second*1, func() (bool, error) {
		return false, errors.New("Test error")
	})

	if err == nil {
		t.Error("Expecting error from loop")
	}
}
