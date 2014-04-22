package util

import (
	"testing"
)

func TestGcd(t *testing.T) {
	if Gcd(2, 4) != 2 {
		t.Fatal("gcd should be 2")
	}

	if Gcd(1, 2, 3, 4, 5) != 1 {
		t.Fatal("gcd should be 1")
	}

	if Gcd(4, 60, 4, 10) != 2 {
		t.Fatal("gcd should be 2")
	}

	if Gcd(3, 6, 9) != 3 {
		t.Fatal("gcd should be 3")
	}

	if Gcd(1, 2, 3, 4, 0) != 0 {
		t.Fatal("Tolerant failed, gcd should be 0")
	}

}
