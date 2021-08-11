package main

import (
	"fmt"
	"io"
	"net/http"
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
			fmt.Printf("\tSkip: content %s\n", contentType)
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
			fmt.Printf("\tSkip: Etag unavailable %s\n", headEtag)
			return false
		} else {
			return true
		}
	} else {
		return true
	}
}

func DownloadImage(filepath string, url string) (finished bool, err error) {

	// Pega o HTTP HEAD da página
	head, err := http.Head(url)
	if err != nil {
		return false, nil
	}

	if !ValidContentType(head.Header) {
		// Caso não seja valido, sai da função
		return false, nil
	}

	if !ValidEtag(head.Header) {
		// Caso não seja valido, sai da função
		return false, nil
	}

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if !ValidContentType(resp.Header) {
		// Caso não seja valido, sai da função
		return false, nil
	}

	if !ValidEtag(resp.Header) {
		// Caso não seja valido, sai da função
		return false, nil
	}

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return false, nil
	}

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return false, err
	}
	defer out.Close()

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return false, err
	}

	fmt.Println("\tSalvo.")

	return true, nil
}
