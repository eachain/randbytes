package randbytes

import "testing"

func TestNew(t *testing.T) {
	const N = 40
	for i := 0; i < N; i++ {
		b := New(i)
		if len(b) != i {
			t.Fatalf("size of b: %v, should be %v", len(b), i)
		}
	}
}
