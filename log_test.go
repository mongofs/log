package log

import (
	"errors"
	"testing"
)

func TestNew(t *testing.T) {
	t.Fatal()
	l := New("testApp",false)
	defer l.Close()

	l.MError(errors.New("connection reset by peer "))
}
