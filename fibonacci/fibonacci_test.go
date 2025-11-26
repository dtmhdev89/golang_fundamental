package fibonacci

import (
	"testing"
)

// func TestFibonacci(t *testing.T) {
// 	if Fibonacci(0) != 0 {
// 		t.Errorf("Expected %v, Actual %v", 0, Fibonacci(0))
// 	}

// 	if Fibonacci(1) != 1 {
// 		t.Errorf("Expected %v, Actual %v", 1, Fibonacci(1))
// 	}

// 	if Fibonacci(6) != 8 {
// 		t.Errorf("Expected %v, Actual %v", 8, Fibonacci(6))
// 	}

// 	if Fibonacci(21) != -1 {
// 		t.Errorf("Expected %v, Actual %v", -1, Fibonacci(21))
// 	}

// 	if Fibonacci(-1) != -1 {
// 		t.Errorf("Expected %v, Actual %v", -1, Fibonacci(-1))
// 	}
// }

func TestFibonacci(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		input int
		want  int
	}{
		{
			name:  "negative value",
			input: -1,
			want:  -1,
		},
		{
			name:  "input too large",
			input: 21,
			want:  -1,
		},
		{
			name:  "zeroth value",
			input: 0,
			want:  0,
		},
		{
			name:  "input is 1",
			input: 1,
			want:  1,
		},
		{
			name:  "input is 1",
			input: 2,
			want:  1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Fibonacci(tt.input)
			if got != tt.want {
				t.Errorf("Fibonacci() = %v, want %v", got, tt.want)
			}
		})
	}
}
