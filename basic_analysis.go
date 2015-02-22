package main

import (
	"github.com/gographics/imagick/imagick"

	"code.google.com/p/plotinum/plot"
	"code.google.com/p/plotinum/plotter"
	"code.google.com/p/plotinum/plotutil"

	"io/ioutil"
	"log"
	//"sort"
	"strings"
	"time"
)

type XY struct {
	X, Y float64
}

// ByX implements sort.Interface for XYs based on the X field.
type ByX plotter.XYs

func (a ByX) Len() int           { return len(a) }
func (a ByX) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByX) Less(i, j int) bool { return a[i].X < a[j].X }

func main() {

	imagick.Initialize()
	defer imagick.Terminate()

	dirs, err := ioutil.ReadDir("output")
	if err != nil {
		panic(err)
	}
	inputs := make([][]float64, 0)
	outputs := make([][]float64, 0)

	redPoints := make(plotter.XYs, 0)
	greenPoints := make(plotter.XYs, 0)
	bluePoints := make(plotter.XYs, 0)

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
			t, err := time.Parse("2006-01-02 15:04:05 -0700 MST.jpg", file.Name())
			if err != nil {
				log.Fatal(err)
			}
			timeHours := (float64)(t.Hour()*60+t.Minute()) / 60.0
			//timeNorm := timeHours / (24.0 * 60.0)
			outputs = append(outputs, []float64{timeHours})
			fileInput := []float64{meanRed, stdevRed, meanBlue, stdevBlue, meanGreen, stdevGreen}
			inputs = append(inputs, fileInput)

			redPoints = append(redPoints, XY{timeHours, meanRed})
			greenPoints = append(greenPoints, XY{timeHours, meanGreen})
			bluePoints = append(bluePoints, XY{timeHours, meanBlue})

			count += 1
			mw.Destroy()
			if count > 50000 {
				break
			}
		}
	}

	/*sort.Sort(ByX(redPoints))
	sort.Sort(ByX(greenPoints))
	sort.Sort(ByX(bluePoints))*/

	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	err = plotutil.AddScatters(p,
		"Red", redPoints,
		"Green", greenPoints,
		"Blue", bluePoints)

	if err != nil {
		panic(err)
	}

	// Save the plot to a PNG file.
	if err := p.Save(11, 8.5, "points.png"); err != nil {
		panic(err)
	}
}
