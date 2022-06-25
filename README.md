# Random-Imgur-Generator

Generate valid random imgur urls and download them to a directory within the project.

As an example, a valid imgur url looks like this:

```
https://i.imgur.com/2TbCIb.png
```

The image identifier can have 5, 6 or 7 characters. The best way of finding valid urls seems to be by trial and error.

By default this code will try to generate 6 digits identifiers, check if the url contains a valid image and, if so, will download it.

Trying to find valid 7 digits identifiers is way harder and it will take quite a while to download a good amount of images. In the other hand, 5 digit identifiers are very easy to find, but will be mostly pretty old images.

## Getting started
Running it with
```
go run .
```
will start the requets.

## Settings
The settings of:
- How many images will be downloaded before the program terminates
- How many workers will be guessing for urls in the background
- Length of the image identifier

can all be changed within the main function in the main.go file.

## Warning
Since the urls are randomly generated, there **WILL** be a good amount of NSFW images.

Running the program with a large amount of workers is a very good way of getting your IP blocked by imgur.
So don't make too many requests per second for an extended amount of time.
