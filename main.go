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

func makeRequest(codeLen int, c chan int) {

	baseUrl := "https://i.imgur.com/"

	// Gera o codigo do imgur, o numero é referete a quantidade de digitos random no fim do link
	// links com 5 digitos random geralmente são imagens mais antigas. Com 6 já são mais novas.
	// Porém, com 6 digitos, a grande maioria dos chutes serão 404.
	code := ImgurCodeGenerator(codeLen)

	// Coloca um .png pra garantir que o site vai abrir só a imagem caso existir
	// Isso funciona mesmo se o original no site seja jpeg, e ele salva .png mesmo sendo jpeg e pelo visto funciona
	imageName := code + ".png"

	// Gera o url do imgur para a imagem
	url := baseUrl + imageName

	// Salva o arquivo caso seja mesmo uma img. Mesmo se for um jpeg vai salvar
	filePath := imageDir + "/" + imageName

	complete, err := DownloadImage(filePath, url)

	if err != nil {
		fmt.Println(err)
		c <- 0
	}

	if complete {
		c <- 1
	} else {
		c <- 0
	}

}

func main() {
	// Tamanho do codigo aleatorio
	codeLen := 6
	// Quantidade de imagens que quero baixar
	imgsWanted := 500
	counter := 0
	c := make(chan int)

	// Roda por n vezes,
	for {
		go makeRequest(codeLen, c)

		counter += <-c

		if counter >= imgsWanted {
			break
		}
	}
}
