package main

import (
	"github.com/fxsjy/gonn/gonn"
	"github.com/gographics/imagick/imagick"

	"bytes"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"math"
	"strings"
	"time"
)

// NN Parameters
const trainImageCount = 2000
const testImageCount = 500
const hiddenNodes = 800
const trainingIterations = 2
const series = false
const width = 120
const height = 80

func ImageFiles(p string) []string {
	filePaths := make([]string, 0)

	dirs, err := ioutil.ReadDir(p)
	if err != nil {
		panic(err)
	}

	for _, dir := range dirs {
		path := p + "/" + dir.Name()
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
			filePaths = append(filePaths, filePath)
		}
	}
	return filePaths
}

func ConvertMWToImage(mw *imagick.MagickWand) (image.Image, error) {
	mw.SetImageFormat("PNG")
	blob := mw.GetImageBlob()

	return png.Decode(bytes.NewReader(blob))
}

func main() {

	imagick.Initialize()
	defer imagick.Terminate()

	inputs := make([][]float64, 0)
	outputs := make([][]float64, 0)
	prev := make([]float64, 0)
	count := 0
	for _, filePath := range ImageFiles("output") {
		//log.Println("File", count, filePath)
		if count%100 == 0 {
			log.Println("File", count)
		}
		mw := imagick.NewMagickWand()
		if err := mw.ReadImage(filePath); err != nil {
			log.Println(err)
			continue
		}
		//mw.CropImage(width, height, x, y)
		//mw.DisplayImage(os.Getenv("DISPLAY"))
		mw.AdaptiveResizeImage(width, height)
		img, err := ConvertMWToImage(mw)
		mw.Destroy()

		if err != nil {
			log.Println(err)
		}

		b := img.Bounds()

		width := b.Max.X - b.Min.X
		height := b.Max.Y - b.Min.Y

		input := make([]float64, width*height*3)

		for y := b.Min.Y; y < b.Max.Y; y++ {
			arrY := y - b.Min.Y
			for x := b.Min.X; x < b.Max.X; x++ {
				arrX := x - b.Min.X
				r, g, b, _ := img.At(x, y).RGBA()

				pos := arrY*width + arrX
				input[pos] = (float64)(r) / 0xFFFF
				input[pos+1] = (float64)(g) / 0xFFFF
				input[pos+2] = (float64)(b) / 0XFFFF
			}
		}
		if series {
			final := make([]float64, width*height*3*2)
			copy(final, input)
			copy(final[len(input):], prev)
			if len(prev) != 0 {
				inputs = append(inputs, final)
				prev = input
			} else {
				prev = input
				continue
			}
		} else {
			inputs = append(inputs, input)
		}

		filePieces := strings.Split(filePath, "/")
		fileName := filePieces[len(filePieces)-1]
		t, err := time.Parse("2006-01-02 15:04:05 -0700 MST.jpg", fileName)
		if err != nil {
			log.Fatal(err)
		}

		timeNorm := (float64)(t.Hour()*60+t.Minute()) / (24.0 * 60.0)
		outputs = append(outputs, []float64{timeNorm})
		count += 1
		if count > (trainImageCount + testImageCount) {
			break
		}
	}
	sample := trainImageCount

	log.Println("Training network")
	nn := gonn.NewNetwork(len(inputs[0]), hiddenNodes, 1, false, 0.25, 0.1)
	nn.Train(inputs[:sample], outputs[:sample], trainingIterations)

	testCount := 0.0
	testSum := 0.0

	timeBlockSum := make([]float64, 24)
	timeBlockCount := make([]float64, 24)

	for i, input := range inputs[sample:] {
		predicted := nn.Forward(input)
		actual := outputs[i+sample]
		predictedTime := predicted[0] * 24
		actualTime := actual[0] * 24

		offset := math.Abs(predictedTime - actualTime)

		if 24-offset < offset {
			offset = 24 - offset
		}

		log.Printf("Test Offset: %.2f, Predicted: %.2f, Actual: %.2f", offset, predictedTime, actualTime)

		block := (int)(math.Floor(actualTime))
		timeBlockSum[block] += offset
		timeBlockCount[block] += 1

		testCount += 1
		testSum += offset
	}
	log.Printf("Estimate: average error is %.2f hours off", testSum/testCount)
	for i, count := range timeBlockCount {
		avg := timeBlockSum[i] / count
		log.Printf("Hour: %d - %.2f", i, avg)
	}
}
