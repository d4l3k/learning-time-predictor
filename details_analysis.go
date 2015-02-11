package main

import (
	"github.com/fxsjy/gonn/gonn"
	"github.com/gographics/imagick/imagick"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

func main() {

	imagick.Initialize()
	defer imagick.Terminate()

	dirs, err := ioutil.ReadDir("output")
	if err != nil {
		panic(err)
	}
	inputs := make([][]float64, 0)
	outputs := make([][]float64, 0)
	count := 0
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
			if !strings.HasSuffix(filePath, ".jpg") {
				continue
			}
			//log.Println("File", count, filePath)
			if count%10 == 0 {
				log.Println("File", count)
			}
			mw := imagick.NewMagickWand()
			if err := mw.ReadImage(filePath); err != nil {
				log.Println(err)
				continue
			}
			meanRed, stdevRed, _ := mw.GetImageChannelMean(imagick.CHANNEL_RED)
			meanGreen, stdevGreen, _ := mw.GetImageChannelMean(imagick.CHANNEL_GREEN)
			meanBlue, stdevBlue, _ := mw.GetImageChannelMean(imagick.CHANNEL_BLUE)
			fileInput := []float64{meanRed, stdevRed, meanBlue, stdevBlue, meanGreen, stdevGreen}
			inputs = append(inputs, fileInput)
			t, err := time.Parse("2006-01-02 15:04:05 -0700 MST.jpg", file.Name())
			if err != nil {
				log.Fatal(err)
			}
			timeNorm := (float64)(t.Hour()*60+t.Minute()) / (24.0 * 60.0)
			outputs = append(outputs, []float64{timeNorm})
			count += 1
			mw.Destroy()
			if count > 8500 {
				break
			}
		}
	}
	sample := 8000

	log.Println("Training network")
	nn := gonn.NewNetwork(len(inputs[0]), 100, 1, false, 0.25, 0.1)
	nn.Train(inputs[:sample], outputs[:sample], 10)

	testCount := 0.0
	testSum := 0.0

	for i, input := range inputs[sample:] {
		predicted := nn.Forward(input)
		actual := outputs[i+sample]
		predictedTime := predicted[0] * 24
		actualTime := actual[0] * 24
		log.Println("Test", predictedTime-actualTime, predictedTime, actualTime)
		testCount += 1
		testSum += predictedTime - actualTime
	}
	log.Println("Estimate:", testSum/testCount, testSum)
}
