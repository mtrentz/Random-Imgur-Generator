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

func DownloadImage(filepath string, url string) (err error) {

	// Pega o HTTP HEAD da página
	head, err := http.Head(url)
	if err != nil {
		return nil
	}

	// // Testa se tem content type, se sim, confere se é valido, caso não, sai da função
	// if _, ok := head.Header["Content-Type"]; ok {
	// 	headType := head.Header["Content-Type"][0]
	// 	if !validTypes[headType] {
	// 		fmt.Printf("\tSkip: content %s\n", headType)
	// 		return nil
	// 	}
	// }

	if !ValidContentType(head.Header) {
		// Caso não seja valido, sai da função
		return nil
	}

	if !ValidEtag(head.Header) {
		// Caso não seja valido, sai da função
		return nil
	}

	// // Faz o mesmo pro Etag e checa contra os unavailables
	// if _, ok := head.Header["Etag"]; ok {
	// 	headEtag := head.Header["Etag"][0]
	// 	// A string vem como "abc", com as aspas mesmo, dai removo elas.
	// 	headEtag = headEtag[1 : len(headEtag)-1]
	// 	if unavailableEtags[headEtag] {
	// 		fmt.Printf("\tSkip: Etag unavailable %s\n", headEtag)
	// 		return nil
	// 	}
	// }

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// // Testa se tem content type, se sim, confere se é valido, caso não, sai da função
	// if _, ok := resp.Header["Content-Type"]; ok {
	// 	respType := resp.Header["Content-Type"][0]
	// 	if !validTypes[respType] {
	// 		fmt.Println("Content no HEAD diferente do Response")
	// 		return nil
	// 	}
	// }

	if !ValidContentType(resp.Header) {
		// Caso não seja valido, sai da função
		return nil
	}

	// // Faz o mesmo pro Etag e checa contra os unavailables
	// if _, ok := resp.Header["Etag"]; ok {
	// 	respEtag := resp.Header["Etag"][0]
	// 	// A string vem como "abc", com as aspas mesmo, dai removo elas.
	// 	respEtag = respEtag[1 : len(respEtag)-1]
	// 	if unavailableEtags[respEtag] {
	// 		fmt.Println("Etag no HEAD diferente do Response")
	// 		return nil
	// 	}
	// }

	if !ValidEtag(resp.Header) {
		// Caso não seja valido, sai da função
		return nil
	}

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	fmt.Println("\tSalvo.")

	return nil
}
