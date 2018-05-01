package ethblock

import "testing"

func TestTrimQuote(t *testing.T) {
	str := trimQuote("\"test\"")
	if str != "test" {
		t.Error("Error TestTrimQuote ", str)
	}
}
