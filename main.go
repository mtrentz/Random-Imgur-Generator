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
	imgsWanted := 1000
	counter := 0
	urlChannel := make(chan string)
	quitChannel := make(chan bool)

	// Numero de goroutines no background
	for i := 0; i <= 10; i++ {
		go FindWorkingUrl(codeLen, urlChannel, quitChannel)
	}

	for val := range urlChannel {
		GetImage(imageDir, val)
		counter++
		if counter >= imgsWanted {
			// Fecha o canal que faz parar as goroutines
			close(quitChannel)
			fmt.Println("Salvo", imgsWanted, "imagens.")
			break
		}
	}

}
