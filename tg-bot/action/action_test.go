package action_test

import (
	"github.com/mightymatth/earthquake-tools/tg-bot/action"
	"testing"
)

func TestEncodeDecode(t *testing.T) {
	a := action.NewSubscriptionAction("123", action.ResetInput)
	e := a.Encode()
	_, err := action.Decode(e)
	if err != nil {
		t.Fatal(err)
	}
}
