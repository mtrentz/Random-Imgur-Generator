package main

import (
	"fmt"
	"os"
	"runtime"
)

const imageDir string = "imgs"

func init() {
	// Caso não exista, cria a página onde vou guardar as imgs
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
	// Tamanho do codigo aleatorio
	codeLen := 6
	// Quantidade de imagens que quero baixar
	imgsWanted := 50
	counter := 0
	urlChannel := make(chan string)

	for i := 0; i <= 500; i++ {
		go FindWorkingUrl(codeLen, urlChannel)
	}

	for val := range urlChannel {
		GetImage(imageDir, val)
		counter++
		if counter >= imgsWanted {
			os.Exit(3)
		}
	}
}
