package main

import (
	"fmt"
	"os"
)

const imageDir string = "imgs"

func init() {
	if _, err := os.Stat(imageDir); os.IsNotExist(err) {
		err := os.Mkdir(imageDir, 0700)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func main() {

	baseUrl := "https://i.imgur.com/"

	// Tenta N vezes gerar um url válido. Pode ser que nao consiga nenhum.
	for i := 0; i < 10; i++ {
		// Gera o codigo do imgur, o numero é referete a quantidade de digitos random no fim do link
		// links com 5 digitos random geralmente são imagens mais antigas. Com 6 já são mais novas.
		// Porém, com 6 digitos, a grande maioria dos chutes serão 404.
		code := ImgurCodeGenerator(5)
		// Coloca um .png pra garantir que o site vai abrir só a imagem caso existir
		// Isso funciona mesmo se o original no site seja jpeg, e ele salva .png mesmo sendo jpeg e pelo visto funciona
		imageName := code + ".png"
		// Gera o url do imgur para a imagem
		url := baseUrl + imageName
		fmt.Printf("%d - %s \n", i+1, url)
		// Salva o arquivo caso seja mesmo uma img. Mesmo se for um jpeg vai salvar
		filePath := imageDir + "/" + imageName
		DownloadImage(filePath, url)
	}
}
