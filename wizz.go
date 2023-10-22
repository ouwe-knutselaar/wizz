package wizz

import (
	"fmt"
	"log"
	"net"
	"net/netip"

	"github.com/ouwe-knutselaar/wizz/connection"
	"github.com/ouwe-knutselaar/wizz/models"
	"github.com/ouwe-knutselaar/wizz/utils"
)

func GetState(bulbIp string) (*models.ResponsePayload, error) {
	payload := &models.RequestPayload{
		Method: "getPilot",
	}
	return connection.SendUdpMessage(bulbIp, payload)
}

func GetConfig(bulbIp string) (*models.ResponsePayload, error) {
	payload := &models.RequestPayload{
		Method: "getSystemConfig",
	}
	return connection.SendUdpMessage(bulbIp, payload)
}

func TurnOnLight(bulbIp string) (*models.ResponsePayload, error) {
	payload := &models.RequestPayload{
		Method: "setPilot",
		Params: models.ParamPayload{
			State: true,
			Speed: 50, // must between 0 - 100
		},
	}
	return connection.SendUdpMessage(bulbIp, payload)
}

func TurnOffLight(bulbIp string) (*models.ResponsePayload, error) {
	payload := &models.RequestPayload{
		Method: "setPilot",
		Params: models.ParamPayload{
			State: false,
			Speed: 50, // must between 0 - 100
		},
	}
	return connection.SendUdpMessage(bulbIp, payload)
}

func SetColorTemp(bulbIp string, value float64) (*models.ResponsePayload, error) {
	payload := &models.RequestPayload{
		Method: "setPilot",
	}
	// normalize the kelvin values - should be removed
	if value < 2500 {
		value = 2500
	}
	if value > 6500 {
		value = 6500
	}
	payload.Params = models.ParamPayload{
		ColorTemp: value,
		State:     true,
		Speed:     50, // must between 0 - 100
	}
	return connection.SendUdpMessage(bulbIp, payload)
}

func SetBrightness(bulbIp string, value float64) (*models.ResponsePayload, error) {
	payload := &models.RequestPayload{
		Method: "setPilot",
	}
	brightnessPercent := utils.HexToPercent(value)
	if brightnessPercent < 10 {
		brightnessPercent = 10
	}
	payload.Params = models.ParamPayload{
		Dimming: int64(brightnessPercent),
		State:   true,
		Speed:   50, // must between 0 - 100
	}
	return connection.SendUdpMessage(bulbIp, payload)
}

func SetColorRGB(bulbIp string, r, g, b float64) (*models.ResponsePayload, error) {

	payload := &models.RequestPayload{
		Method: "setPilot",
	}
	r = valuecorrection(r)
	g = valuecorrection(g)
	b = valuecorrection(b)
	payload.Params = models.ParamPayload{
		R:     r,
		G:     g,
		B:     b,
		State: true,
		Speed: 50, // must between 0 - 100
	}
	return connection.SendUdpMessage(bulbIp, payload)
}

func SetColorScene(bulbIp string, sceneId int64) (*models.ResponsePayload, error) {
	var exists bool
	payload := &models.RequestPayload{
		Method: "setPilot",
		Params: models.ParamPayload{
			SceneId: 1,
			State:   true,
			Speed:   50, // must between 0 - 100
		},
	}
	if _, exists = models.SceneModel[sceneId]; exists == true {
		payload.Params.SceneId = sceneId
	}
	return connection.SendUdpMessage(bulbIp, payload)
}

func SetColorWarmWhite(bulbIp string, value float64) (*models.ResponsePayload, error) {
	payload := &models.RequestPayload{
		Method: "setPilot",
		Params: models.ParamPayload{
			WarmWhite: 0,
			State:     true,
			Speed:     50, // must between 0 - 100
		},
	}
	value = valuecorrection(value)
	payload.Params.WarmWhite = value
	return connection.SendUdpMessage(bulbIp, payload)
}

func SetColorColdWhite(bulbIp string, value float64) (*models.ResponsePayload, error) {
	payload := &models.RequestPayload{
		Method: "setPilot",
		Params: models.ParamPayload{
			ColdWhite: 0,
			State:     true,
			Speed:     50, // must between 0 - 100
		},
	}
	value = valuecorrection(value)
	payload.Params.ColdWhite = value
	return connection.SendUdpMessage(bulbIp, payload)
}

func getSourceAddress() ([4]byte, error) {
	addlist, _ := net.InterfaceAddrs()

	for _, addr := range addlist {
		fmt.Println(addr.Network(), addr.String())
		netAddr, err := netip.ParsePrefix(addr.String())
		if err == nil {
			if netAddr.Addr().Is4() {
				if netAddr.Addr().As4()[0] == 192 {
					log.Println("Found", netAddr.Addr().String())
					return netAddr.Addr().As4(), nil
				}
			}
		}
	}
	var none [4]byte
	return none, nil
}

func valuecorrection(value float64) float64 {
	if value > 255 {
		value = 255
	}
	if value < 0 {
		value = 0
	}
	return value
}
