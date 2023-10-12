package connection

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/ouwe-knutselaar/wizz/models"
)

const (
	Port        = "38899"
	ReadTimeOut = 5
)

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
		return nil, fmt.Errorf(`Unable resolve udp: %s`, err)
	}
	if conn, err = net.DialUDP("udp", nil, remoteAddr); err != nil {
		return nil, fmt.Errorf(`Unable to dial up to udp: %s`, err)
	}
	defer conn.Close()
	// marshall payload to json string
	if payload, err = json.Marshal(message); err != nil {
		return nil, fmt.Errorf(`Unable to marshal payload: %s`, err)
	}
	payloadString := string(payload)
	log.Printf(`Payload string: %s`, payloadString)
	// send payload to bulb
	if _, err = conn.Write(payload); err != nil {
		return nil, fmt.Errorf(`Unable to send message to UDP: %s`, err)
	}
	// read response from bulb
	conn.SetReadDeadline(time.Now().Add(ReadTimeOut * time.Second))

	if _, err = bufio.NewReader(conn).Read(response); err != nil {
		return nil, fmt.Errorf(`Unable to read message from UDP: %s`, err)
	}
	result := []byte(strings.Trim(string(response), "\x00'"))
	// convert string result to struct again
	if err = json.Unmarshal(result, responsePayload); err != nil {
		return nil, fmt.Errorf(`Unable to unmarshal payload: %s`, err)
	}
	return responsePayload, nil
}
