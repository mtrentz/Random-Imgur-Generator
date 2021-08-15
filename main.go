package main

import (
	"fmt"
	"os"
	"runtime"
)

const imageDir string = "imgs"

func init() {
	// Case not exsits, create folder to store images
	if _, err := os.Stat(imageDir); os.IsNotExist(err) {
		err := os.Mkdir(imageDir, 0700)
		if err != nil {
			fmt.Println(err)
		}
	}

	fmt.Println("OS:", runtime.GOOS)
	fmt.Println("CPUs:", runtime.NumCPU())
}

func main() {
	// Size of imgur code, for exemple i.imgur.com/x123xD -> 6 random characters
	// Codes with 5 characters are older images uploaded to imgur.
	// Codes with 6 are usually newer, but its harder to find working urls
	codeLen := 6
	// Amount of images to download
	imgsWanted := 10000
	counter := 0
	urlChannel := make(chan string)
	quitChannel := make(chan bool)

	// Number of goroutines running in the background
	// Its ok to add more than num of CPU cores since most of time is spent waiting for http requests
	for i := 0; i <= 10; i++ {
		go FindWorkingUrl(codeLen, urlChannel, quitChannel)
	}

	for val := range urlChannel {
		GetImage(imageDir, val)
		counter++
		if counter >= imgsWanted {
			// Close channel and stop all goroutines
			close(quitChannel)
			fmt.Println("Salvo", imgsWanted, "imagens.")
			// Has to break out of loop, else code will be stuck waiting to read from urlChannel
			break
		}
	}

}
