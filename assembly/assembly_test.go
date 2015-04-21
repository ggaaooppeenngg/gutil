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

func TestMethod(t *testing.T) {
	var tt T = 1
	if tt.Get() != 1 {
		t.Fail()
	}
}
func TestIndex(t *testing.T) {
	t.Skip()
	var a = []int64{100, 1000, 10000}
	if Index(a, 1000) != 1 {
		t.Logf("%d", Index(a, 1000))
		t.Fail()
	}
}
