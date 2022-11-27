package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"inputpoller/poller"
	"inputpoller/stdinput"
)

func main() {
	inputProvider := stdinput.Provider{}
	inputScanPoller := poller.New(5*time.Second, 5, &inputProvider)

	quitChan := make(chan struct{})

	go scanning(inputScanPoller, quitChan)

	interruptScanningChan := make(chan os.Signal, 1)
	signal.Notify(interruptScanningChan, syscall.SIGINT, syscall.SIGTERM)
	<-interruptScanningChan

	inputProvider.StopProviding()
	fmt.Println("input any text to exit")

	<-quitChan
}

func scanning(scanPoller *poller.ScanPoller, quitChan chan struct{}) {
	for {
		input, err := scanPoller.Poll()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			fmt.Println("error: ", err.Error())

			continue
		}

		fmt.Println("scanned:", input)
	}

	quitChan <- struct{}{}
}
