package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// Valid types of imgur files.
var validTypes = map[string]bool{
	"image/png":  true,
	"image/jpeg": true,
}

// Some imgur links come with an etag that identifies the usual type of unavailable pages
var unavailableEtags = map[string]bool{
	// Some images were removed from imgur, if you try to download them, you'll get a .png saying the img was removed
	"d835884373f4d6c8f24742ceabe74946": true,
	// 404 on imgurs are not actually 404. It comes with a 200 status. This etag represents true 404, without images found on the url
	"a7cb396d0db6af2e63870985cb086fa1": true,
}

func ValidContentType(header map[string][]string) (valid bool) {
	// Check if the headed actually has a content-type field
	if _, ok := header["Content-Type"]; ok {
		// If so, get the value
		contentType := header["Content-Type"][0]
		// Check if type is actually an image
		if !validTypes[contentType] {
			// fmt.Printf("Skip: content %s\n", contentType)
			return false
		} else {
			return true
		}
	} else {
		// Usually, headers without content-type are STILL valid images.
		return true
	}
}

func ValidEtag(header map[string][]string) (valid bool) {
	// Do a very similar check for the etags.
	if _, ok := header["Etag"]; ok {
		headEtag := header["Etag"][0]
		// Remove doublequotes from the string
		headEtag = headEtag[1 : len(headEtag)-1]
		if unavailableEtags[headEtag] {
			// fmt.Printf("Skip: Etag unavailable %s\n", headEtag)
			return false
		} else {
			return true
		}
	} else {
		return true
	}
}

func FindWorkingUrl(codeLen int, urlChan chan<- string, quitChannel <-chan bool) {
	// This function goes into an infinite loop, that will only break when quitChannel is closed by the sender
	// It will test for valid imgur urls containing either a png or jpeg. When found, send the valid url via the url channel.
	select {
	case <-quitChannel:
		return
	default:
		baseUrl := "https://i.imgur.com/"

		for {
			code := ImgurCodeGenerator(codeLen)

			requestUrl := baseUrl + code + ".png"

			// Request only the page header
			head, err := http.Head(requestUrl)
			if err != nil {
				fmt.Println("Error requesting the HEAD", err)
				continue
			}

			if ValidContentType(head.Header) {
				if ValidEtag(head.Header) {
					urlChan <- requestUrl
				}
			}
		}
	}
}

func GetContentType(header map[string][]string) (contentType string) {
	// The point of this function is to get the content type (jpeg, png)
	// so the file can be saved with proper extension (abc.jpeg).
	// If it does not find a content type (sometimes its missing on header)
	// it will default to a string 'png'.

	// Check if header exists
	if _, ok := header["Content-Type"]; ok {
		// If so, get the value
		contentType := header["Content-Type"][0]
		// Content Types are strings like "image/jpeg" or "image/png". So i'll get only the second part
		imageContentType := strings.Split(contentType, "/")[1]
		// Return the content type
		return imageContentType
	} else {
		// It shouldn't get to here, but if it does, defaults to png
		return "png"
	}
}

func getDateString() string {
	// This function returns a string with the current date and time.
	// Which is gonna be used for the file path.
	return time.Now().Format("2006-01-02")
}

func guaranteeDirExists(dir string) {
	// Case not exsits, create folder to store images
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0700)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func GetImage(imageDir string, imgUrl string) {
	// Starting from a valid imgur url, containing an image
	// this function will download that image and save to a directory

	// Get the data
	resp, err := http.Get(imgUrl)
	if err != nil {
		fmt.Println("Error requesting the image", err)
		return
	}
	defer resp.Body.Close()

	// Here tests for valid content type of the header again, just in case something went wrong
	// on the header-only request.
	if !ValidContentType(resp.Header) {
		// If not valid, gives up on downloading
		return
	}

	if !ValidEtag(resp.Header) {
		// If not valid, gives up on downloading
		return
	}

	// Check server response
	if resp.StatusCode != http.StatusOK {
		// If not valid, gives up on downloading
		return
	}

	// Getting the image name, for example: www.imgur.com/aBc123.png -> aBc123.png
	u, err := url.Parse(imgUrl)
	if err != nil {
		fmt.Println("Failed to parse URL", err)
	}
	imageName := u.Path[1:]

	// Fixing the filename. By my default choice all images are with .png extension, since it helps when requesting images from imgur
	// because when imgur links are like i.imgur.com/abc.png the return will be a webpage with only the image in it.
	// Links like i.imgur.com/abc would return a page with the rest of the imgur website interface.
	// Some images are not png though, and now I will fix the content type (getting the correct one from the header)
	// to save the file with the proper extension.
	imageExtension := GetContentType(resp.Header)
	imageNameNoExtension := strings.Split(imageName, ".")[0]
	newImageName := imageNameNoExtension + "." + imageExtension

	// Add the image name to the final image path
	// imagePath := imageDir + "/" + newImageName
	folderDir := imageDir + "/" + getDateString()
	guaranteeDirExists(folderDir)
	imagePath := folderDir + "/" + newImageName

	// Create the file
	out, err := os.Create(imagePath)
	if err != nil {
		fmt.Println("Failed to create file: ", err)
		return
	}
	defer out.Close()

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Println("Error trying to write file", err)
		return
	}
	fmt.Println("Saved.")

}
