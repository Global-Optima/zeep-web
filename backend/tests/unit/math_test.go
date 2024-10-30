// backend/main_test.go
package unit

import "testing"

// add function adds two integers.
func add(a int, b int) int {
	return a + b
}

// subtract function subtracts two integers.
func subtract(a int, b int) int {
	return a - b
}

// multiply function multiplies two integers.
func multiply(a int, b int) int {
	return a * b
}

// divide function divides two integers and returns the quotient and remainder.
func divide(a int, b int) (int, int) {
	if b == 0 {
		return 0, 0 // Handle division by zero
	}
	return a / b, a % b
}

func TestAdd(t *testing.T) {
	result := add(2, 3)
	expected := 5
	if result != expected {
		t.Errorf("Add(2, 3) = %d; want %d", result, expected)
	}
}

func TestSubtract(t *testing.T) {
	result := subtract(5, 3)
	expected := 2
	if result != expected {
		t.Errorf("Subtract(5, 3) = %d; want %d", result, expected)
	}
}

func TestMultiply(t *testing.T) {
	result := multiply(2, 3)
	expected := 6
	if result != expected {
		t.Errorf("Multiply(2, 3) = %d; want %d", result, expected)
	}
}

func TestDivide(t *testing.T) {
	result, remainder := divide(7, 2)
	expectedQuotient := 3
	expectedRemainder := 1
	if result != expectedQuotient || remainder != expectedRemainder {
		t.Errorf("Divide(7, 2) = %d, remainder %d; want quotient %d, remainder %d", result, remainder, expectedQuotient, expectedRemainder)
	}
}

func TestDivideByZero(t *testing.T) {
	result, remainder := divide(7, 0)
	expectedQuotient := 0
	expectedRemainder := 0
	if result != expectedQuotient || remainder != expectedRemainder {
		t.Errorf("Divide(7, 0) = %d, remainder %d; want quotient %d, remainder %d", result, remainder, expectedQuotient, expectedRemainder)
	}
}
