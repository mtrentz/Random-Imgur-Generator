package main

import (
	"fmt"
	"os"
)

const imageDir string = "imgs"

var c chan int

func init() {
	// Caso não exista, cria a página onde vou guardar as imgs
	if _, err := os.Stat(imageDir); os.IsNotExist(err) {
		err := os.Mkdir(imageDir, 0700)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func makeRequest(codeLen int, c chan int) {

	baseUrl := "https://i.imgur.com/"

	// Gera o codigo do imgur, o numero é referete a quantidade de digitos random no fim do link
	// links com 5 digitos random geralmente são imagens mais antigas. Com 6 já são mais novas.
	// Porém, com 6 digitos, a grande maioria dos chutes serão 404.
	code := ImgurCodeGenerator(codeLen)
	fmt.Println(code)

	// Coloca um .png pra garantir que o site vai abrir só a imagem caso existir
	// Isso funciona mesmo se o original no site seja jpeg, e ele salva .png mesmo sendo jpeg e pelo visto funciona
	imageName := code + ".png"

	// Gera o url do imgur para a imagem
	url := baseUrl + imageName

	// Salva o arquivo caso seja mesmo uma img. Mesmo se for um jpeg vai salvar
	filePath := imageDir + "/" + imageName

	// complete, err := DownloadImage(filePath, url)
	fmt.Println(filePath + url)
	complete := false
	// if err != nil {
	// 	fmt.Println(err)
	// 	c <- 0
	// }

	if complete {
		c <- 1
	} else {
		c <- 0
	}

}

func main() {
	// Quantidade de imagens que quero baixar
	imgsWanted := 10
	counter := 0

	// Roda por n vezes,
	for i := 0; i < 500; i++ {
		go makeRequest(5, c)

		counter += <-c

		if counter >= imgsWanted {
			break
		}
	}
}
