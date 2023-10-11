package connection

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/ouwe-knutselaar/wizz/models"
)

const Port = "38899"

func SendUdpMessage(host string, message *models.RequestPayload) (*models.ResponsePayload, error) {
	var (
		err             error
		response        = make([]byte, 4096)
		responsePayload = new(models.ResponsePayload)
		remoteAddr      = new(net.UDPAddr)
		conn            = new(net.UDPConn)
		payload         []byte
	)
	time.Sleep(time.Second * 1)
	// doing connection to UDP
	if remoteAddr, err = net.ResolveUDPAddr("udp", fmt.Sprintf(`%s:%s`, host, Port)); err != nil {
		err = errors.New(fmt.Sprintf(`Unable resolve udp: %s`, err))
		return nil, err
	}
	if conn, err = net.DialUDP("udp", nil, remoteAddr); err != nil {
		err = errors.New(fmt.Sprintf(`Unable to dial up to udp: %s`, err))
		return nil, err
	}
	defer conn.Close()
	// marshall payload to json string
	if payload, err = json.Marshal(message); err != nil {
		err = errors.New(fmt.Sprintf(`Unable to marshal payload: %s`, err))
		return nil, err
	}
	payloadString := string(payload)
	log.Println(fmt.Sprintf(`Payload string: %s`, payloadString))
	// send payload to bulb
	if _, err = conn.Write(payload); err != nil {
		err = errors.New(fmt.Sprintf(`Unable to send message to UDP: %s`, err))
		return nil, err
	}
	// read response from bulb
	if _, err = bufio.NewReader(conn).Read(response); err != nil {
		err = errors.New(fmt.Sprintf(`Unable to read message from UDP: %s`, err))
		return nil, err
	}
	result := []byte(strings.Trim(string(response), "\x00'"))
	// convert string result to struct again
	if err = json.Unmarshal(result, responsePayload); err != nil {
		err = errors.New(fmt.Sprintf(`Unable to unmarshal payload: %s`, err))
		return nil, err
	}
	return responsePayload, nil
}
