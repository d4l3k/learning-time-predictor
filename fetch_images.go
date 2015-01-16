package main

import (
	"crypto/md5"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

const filechunk = 8192 // we settle for 8KB

func hash(path string) string {
	file, err := os.Open(path)
	if err != nil {
		panic(err.Error())
	}
	defer file.Close()

	// calculate the file size
	info, _ := file.Stat()

	filesize := info.Size()

	blocks := uint64(math.Ceil(float64(filesize) / float64(filechunk)))

	hash := md5.New()

	for i := uint64(0); i < blocks; i++ {
		blocksize := int(math.Min(filechunk, float64(filesize-int64(i*filechunk))))
		buf := make([]byte, blocksize)

		file.Read(buf)
		io.WriteString(hash, string(buf)) // append into the hash
	}
	return string(hash.Sum(nil))
}

func hashBytes(bytes []byte) string {
	hash := md5.New()
	io.WriteString(hash, string(bytes))
	return string(hash.Sum(nil))
}

func main() {
	nameCleaningRegex := regexp.MustCompile(`[:/]`)
	for {
		feeds := []string{"http://images.drivebc.ca/bchighwaycam/pub/cameras/13.jpg"}
		for _, feed := range feeds {
			cleanName := nameCleaningRegex.ReplaceAll([]byte(feed), []byte("-"))
			dir := "./output/" + string(cleanName)
			os.MkdirAll(dir, 0755)
			files, _ := filepath.Glob(dir + "/*")

			lastHash := ""
			if len(files) > 0 {
				lastHash = hash(files[len(files)-1])
			}

			log.Println("Checking", feed)

			resp, err := http.Get(feed)
			if err != nil {
				log.Println("Get err", err)
				continue
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)

			if hashBytes(body) != lastHash {
				log.Println(" - Found new version.")
				file := dir + "/" + time.Now().String() + ".jpg"
				err = ioutil.WriteFile(file, body, 0644)
				if err != nil {
					log.Println("File write error!", err)
					continue
				}
			} else {
				log.Println(" - No new version.")
			}
		}
		time.Sleep(time.Minute)
	}
}
