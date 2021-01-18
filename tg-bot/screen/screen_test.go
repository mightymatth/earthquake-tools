package screen_test

import (
	"github.com/mightymatth/earthquake-tools/tg-bot/screen"
	"testing"
)

func TestEncodeDecode(t *testing.T) {
	scr := screen.NewSubscriptionScreen("123", screen.ResetInput)
	e := scr.Encode()
	_, err := screen.Decode(e)
	if err != nil {
		t.Fatal(err)
	}
}
