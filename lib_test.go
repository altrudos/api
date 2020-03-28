package charityhonor

import "testing"

func TestAmountToString(t *testing.T) {
	was := AmountToString(1234)
	if was != "12.34" {
		t.Error("12.34 should be 12.34 not", was)
	}

	was = AmountToString(7500)
	if was != "75.00" {
		t.Error("75 should be 75.00 not", was)
	}

	was = AmountToString(235)
	if was != "2.35" {
		t.Error("235 should be 2.35 not", was)
	}

	was = AmountToString(2350)
	if was != "23.50" {
		t.Error("23.5 should be 23.50 not", was)
	}

	was = AmountToString(001)
	if was != "0.01" {
		t.Error("0.01 should be 0.01 not", was)
	}

	was = AmountToString(1)
	if was != "0.01" {
		t.Error("1 should be 0.01 not", was)
	}

	was = AmountToString(100)
	if was != "1.00" {
		t.Error("1 should be 1.00 not", was)
	}
}
