package id3v2

import (
	"bytes"
	"math/big"
	"testing"
)

// Make sure that counter of popularimeter is at least 4 bytes even if it's small number.
func TestPopularimeterFrameSmallCounter(t *testing.T) {
	popmFrame := PopularimeterFrame{
		Email:   "foo@bar.com",
		Rating:  1,
		Counter: big.NewInt(1),
	}

	expectedBodyLength := len(popmFrame.Email) + 1 + 1 + 4

	buf := new(bytes.Buffer)
	written, err := popmFrame.WriteTo(buf)
	if err != nil {
		t.Fatalf("Error by writing: %v", err)
	}

	if written != int64(expectedBodyLength) {
		t.Fatalf("Expected popularimeter frame body length: %v, got: %v", expectedBodyLength, written)
	}

	expectedCounter := []byte{0, 0, 0, 1}
	gotCounter := buf.Bytes()[expectedBodyLength-4:]
	if !bytes.Equal(expectedCounter, gotCounter) {
		t.Fatalf("Expected popularimeter counter: %v, got: %v", expectedCounter, gotCounter)
	}
}
