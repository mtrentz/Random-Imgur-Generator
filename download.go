package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

// No HEADER da pra checar pelo tipo do conteudo
var validTypes = map[string]bool{
	"image/png":  true,
	"image/jpeg": true,
}

// Os HEAD do imgur tem um Etag de uma imagem que é invalida, mas mesmo assim é uma imagem.
var unavailableEtags = map[string]bool{
	// Pras imagens que vem com o texto nela dizendo unavailable
	"d835884373f4d6c8f24742ceabe74946": true,
	// Pras páginas que carregam no imgur.com mas são na real um 404
	"a7cb396d0db6af2e63870985cb086fa1": true,
}

func ValidContentType(header map[string][]string) (valid bool) {
	// Confere se tem o Content-Type no header
	if _, ok := header["Content-Type"]; ok {
		// Caso sim, pega o valor do content type
		contentType := header["Content-Type"][0]
		// Confere se o type é valido como imagem
		if !validTypes[contentType] {
			fmt.Printf("Skip: content %s\n", contentType)
			return false
		} else {
			return true
		}
	} else {
		// Caso não tenha o content-type no header a imagem ainda assim pode ser válida
		return true
	}
}

func ValidEtag(header map[string][]string) (valid bool) {
	// Faz o mesmo pro Etag e checa contra os unavailables
	if _, ok := header["Etag"]; ok {
		headEtag := header["Etag"][0]
		// A string vem como "abc", com as aspas mesmo, dai removo elas.
		headEtag = headEtag[1 : len(headEtag)-1]
		if unavailableEtags[headEtag] {
			fmt.Printf("Skip: Etag unavailable %s\n", headEtag)
			return false
		} else {
			return true
		}
	} else {
		return true
	}
}

func FindWorkingUrl(codeLen int, urlChan chan<- string) {
	baseUrl := "https://i.imgur.com/"

	for {
		code := ImgurCodeGenerator(codeLen)

		requestUrl := baseUrl + code + ".png"

		// Pega o HTTP HEAD da página
		head, err := http.Head(requestUrl)
		if err != nil {
			fmt.Println("Erro no request do HEAD", err)
			return
		}

		if ValidContentType(head.Header) {
			if ValidEtag(head.Header) {
				urlChan <- requestUrl
			}
		}
	}
}

func GetImage(imageDir string, imgUrl string, counterChan chan int) {

	// Nome da imagem, por exemplo: www.imgur.com/aBc123.png -> aBc123.png
	u, err := url.Parse(imgUrl)
	if err != nil {
		fmt.Println("Erro ao parsear a url", err)
	}
	imageName := u.Path[1:]
	// Adiciona o nome no diretorio pra salvar a imagem
	imagePath := imageDir + "/" + imageName

	// Get the data
	resp, err := http.Get(imgUrl)
	if err != nil {
		fmt.Println("Erro no request da imagem", err)
		return
	}
	defer resp.Body.Close()

	if !ValidContentType(resp.Header) {
		// Caso não seja valido, sai da função
		return
	}

	if !ValidEtag(resp.Header) {
		// Caso não seja valido, sai da função
		return
	}

	// Check server response
	if resp.StatusCode != http.StatusOK {
		// TODO: Podia gerar um erro pra status code aqui
		return
	}

	// Create the file
	out, err := os.Create(imagePath)
	if err != nil {
		fmt.Println("Erro ao criar arquivo", err)
		return
	}
	defer out.Close()

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Println("Erro ao salvar arquivo", err)
		return
	}

	counterChan <- 1
	fmt.Println("Salvo.")

}
