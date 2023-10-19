package connection

import (
	"log"
	"testing"

	"github.com/ouwe-knutselaar/wizz/models"
)

func TestUnknownIPaddress(t *testing.T) {
	var (
		payload = &models.RequestPayload{
			Method: "getPilot",
		}
	)

	_, err := SendUdpMessage("1.2.3.4", payload)
	if err == nil {
		t.Fatal("Can connect to unknown IP adress")
	}
	log.Println(err)
}
