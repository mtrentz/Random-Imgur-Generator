package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func DownloadImage(filepath string, url string) (err error) {
	// No HEADER da pra checar pelo tipo do conteudo
	validTypes := map[string]bool{
		"image/png":  true,
		"image/jpeg": true,
	}

	// acho que esse etag é pro 404 a7cb396d0db6af2e63870985cb086fa1
	// Os HEAD do imgur tem um Etag de uma imagem que é invalida, mas mesmo assim é uma imagem.
	unavailableETag := "d835884373f4d6c8f24742ceabe74946"

	// Pega o HTTP HEAD da página
	head, err := http.Head(url)
	if err != nil {
		return nil
	}

	// Testa primeiro se o header tem so campos necessarios, caso sim, utiliza eles
	if _, ok := head.Header["Content-Type"]; !ok {
		return nil
	}
	headType := head.Header["Content-Type"][0]
	if _, ok := head.Header["Etag"]; !ok {
		return nil
	}
	headEtag := head.Header["Etag"][0]
	headEtag = headEtag[1 : len(headEtag)-1]

	// Se o header não é de imagem, eu retorno nil e saio da função
	if !validTypes[headType] {
		return nil
	}

	if headEtag == unavailableETag {
		return nil
	}

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if _, ok := resp.Header["Content-Type"]; !ok {
		return nil
	}
	respType := resp.Header["Content-Type"][0]

	if _, ok := resp.Header["Etag"]; !ok {
		return nil
	}
	respEtag := resp.Header["Etag"][0]
	respEtag = respEtag[1 : len(respEtag)-1]

	fmt.Println("----")
	fmt.Println(url)
	// fmt.Println(resp.Header)
	fmt.Println(respEtag)
	fmt.Println(unavailableETag)

	// Confere mesmo assim se a resposta é uma imagem mesmo e não é aquele Etag
	if !validTypes[respType] {
		fmt.Println("HEAD request deu diferente que o HEAD do GET")
		return nil
	}

	if respEtag == unavailableETag {
		fmt.Println("HEAD request deu diferente que o HEAD do GET")
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

	return nil
}
