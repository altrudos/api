package charityhonor

import "testing"

func TestAmountToString(t *testing.T) {
	was := AmountToString(1234)
	if was != "12.34" {
		t.Error("200 should be 2.00 not", was)
	}

	was = AmountToString(75)
	if was != "0.75" {
		t.Error("75 should be 0.75 not", was)
	}

	was = AmountToString(235)
	if was != "2.35" {
		t.Error("235 should be 2.35 not", was)
	}

	was = AmountToString(1)
	if was != "0.01" {
		t.Error("1 should be 0.01 not", was)
	}
}
