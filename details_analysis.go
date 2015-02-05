package main

import (
	"github.com/gographics/imagick/imagick"

	"io/ioutil"
	"log"
)

func main() {

	imagick.Initialize()
	defer imagick.Terminate()

	dirs, err := ioutil.ReadDir("output")
	if err != nil {
		panic(err)
	}
	for _, dir := range dirs {
		path := "output/" + dir.Name()
		log.Println("Processing", path)
		files, err := ioutil.ReadDir(path)
		log.Println("Found", len(files), "files")
		if err != nil {
			log.Panic("Error", err)
			continue
		}
		for _, file := range files {
			filePath := path + "/" + file.Name()
			log.Println("File", filePath)
			mw := imagick.NewMagickWand()
			defer mw.Destroy()
			if err := mw.ReadImage(filePath); err != nil {
				panic(err)
			}
			mean, stdev, err := mw.GetImageChannelMean(imagick.CHANNEL_RED)
			if err != nil {
				log.Panic(err)
			}
			log.Println("Mean", mean, "Stdev", stdev)
			break
		}
	}
}
