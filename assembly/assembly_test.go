package assembly

import (
	"testing"
)

func TestAdd(t *testing.T) {
	if Add(1, 2) != 3 {
		t.Fail()
	}
}

func TestAbs(t *testing.T) {
	if Abs(-1) != 1 {
		t.Logf("%d", Abs(-1))
		t.Fail()
	}
}
