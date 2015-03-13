package main

import (
	"github.com/gographics/imagick/imagick"
	"github.com/golang/protobuf/proto"
	"github.com/jmhodges/levigo"

	"./caffe"

	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"
)

func main() {
	opts := levigo.NewOptions()
	opts.SetCache(levigo.NewLRUCache(3 << 30))
	opts.SetCreateIfMissing(true)
	db, err := levigo.Open("picture_train.db", opts)
	wo := levigo.NewWriteOptions()
	db2, err := levigo.Open("picture_test.db", opts)

	path := "images"
	files, err := ioutil.ReadDir(path)
	log.Println("Found", len(files), "files")
	if err != nil {
		log.Panic("Error", err)
	}
	count := 0
	countTest := 0
	for _, file := range files {
		filePath := path + "/" + file.Name()
		if !strings.HasSuffix(filePath, ".jpg") {
			continue
		}
		if count%10 == 0 {
			log.Println("File", count)
		}
		count += 1
		mw := imagick.NewMagickWand()
		if err := mw.ReadImage(filePath); err != nil {
			log.Println(err)
			continue
		}

		t, err := time.Parse("2006-01-02_15:04:05_-0700_MST.jpg", file.Name())
		if err != nil {
			log.Fatal(err)
		}
		timeNorm := (float64)(t.Hour()*60+t.Minute()) / (24.0 * 60.0)
		height := int32(mw.GetImageHeight())
		width := int32(mw.GetImageWidth())
		channels := int32(3)
		encoded := true
		bytes, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Println(err)
			continue
		}
		datum := caffe.Datum{
			Height:    &height,
			Width:     &width,
			Channels:  &channels,
			Encoded:   &encoded,
			Data:      bytes,
			FloatData: []float32{float32(timeNorm)},
		}
		data, err := proto.Marshal(&datum)
		if err != nil {
			log.Println("marshaling error: ", err)
			continue
		}
		if count > 9800 {
			err = db2.Put(wo, []byte(strconv.Itoa(countTest)), data)
			countTest += 1
		} else {
			err = db.Put(wo, []byte(strconv.Itoa(count-1)), data)
		}
		if err != nil {
			log.Println(err)
			continue
		}
		mw.Destroy()
	}
}
