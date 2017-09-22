package state

import (
	"log"
	"testing"
)

func TestSecret(t *testing.T) {
	s := &Secret{Value: "abc"}
	key, err := s.Encrypt()
	if err != nil {
		log.Print(err)
		t.Fail()
	}
	err = s.Decrypt(key)
	if err != nil {
		log.Print(err)
		t.Fail()
	}
	if s.Value != "abc" {
		log.Print("failed...")
		t.Fail()
	}
}
