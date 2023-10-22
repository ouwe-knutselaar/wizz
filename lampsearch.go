package wizz

import (
	"fmt"
	"log"
	"sync"

	probing "github.com/prometheus-community/pro-bing"
)

var workers int = 15

func SearchLamp(a, b, c byte) {
	adrChannel := make(chan string)
	var wg sync.WaitGroup
	wg.Add(1)

	for tel := 0; tel < workers; tel++ {
		go searchWorker(adrChannel, wg)
	}

	for tel := 0; tel < 256; tel++ {
		adrStr := fmt.Sprintf("%d.%d.%d.%d", a, b, c, tel)
		adrChannel <- adrStr

	}
	wg.Done()
	wg.Wait()
	close(adrChannel)
}

func searchWorker(adrChannel chan string, wg sync.WaitGroup) {
	log.Println("Add worker")
	wg.Add(1)
	for adrString := range adrChannel {
		//log

		pinger, _ := probing.NewPinger(adrString)

		pinger.SetPrivileged(true)
		pinger.Count = 3
		pinger.Size = 24

		pinger.Run()
		stats := pinger.Statistics()
		if stats.PacketsRecv > 0 {
			log.Println("Check", adrString, "stats", stats.PacketsRecv)
			response, err := GetState(adrString)
			if err == nil {
				log.Println("Found lamp at", adrString, "  ", response)
			}
		}
	}
	wg.Done()
}

func getLocalAddress() {
	var temp []string

	localAddr, _ := getSourceAddress()

	var tel byte
	for tel = 1; tel < 255; tel++ {
		localAddr[3] = tel
		addrStr := fmt.Sprintf("%d.%d.%d.%d", localAddr[0], localAddr[1], localAddr[2], localAddr[3])
		log.Println("Search", addrStr)
		_, err := GetState(addrStr)
		if err == nil {
			temp = append(temp, addrStr)
			log.Println("Found lamp at", addrStr)
		}
	}

}
