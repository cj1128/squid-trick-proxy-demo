package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
	"sync"
)

const serverPrefix = "http://127.0.0.1:7777"
const rootDir = "/tmp/squid"
const userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/68.0.3440.106 Safari/537.36"

var logFile, _ = os.Create(path.Join(rootDir, "rewrite.log"))

var logMutex = &sync.Mutex{}
var printMutex = &sync.Mutex{}

func log(prefix, msg string) {
	logMutex.Lock()
	logFile.WriteString(fmt.Sprintf("[%s] %s\n", prefix, msg))
	logMutex.Unlock()
}

func main() {
	defer logFile.Close()

	scanner := bufio.NewScanner(os.Stdin)

	wg := &sync.WaitGroup{}
	workChan := make(chan string, 10000)

	// workers

	for i := 0; i < 50; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			for line := range workChan {
				log("squid", line)
				result, err := process(line)

				if err != nil {
					log("err", err.Error())
				} else {
					log("result", result)
				}

				printMutex.Lock()
				fmt.Println(result)
				printMutex.Unlock()
			}
		}()
	}

	for scanner.Scan() {
		workChan <- scanner.Text()
	}

	close(workChan)
	wg.Wait()
}

var isImgRegexp = regexp.MustCompile(`\.(jpg|png|jpeg)$`)

func process(line string) (string, error) {
	parts := strings.Split(line, " ")
	channelID := parts[0]
	urlStr := parts[1]

	url, _ := url.Parse(urlStr)

	if !isImgRegexp.MatchString(url.Path) {
		return channelID + " OK", nil
	}

	_, filename := path.Split(url.Path)
	filename = cleanFilename(filename)

	dstPath := path.Join(rootDir, "cache", filename)

	if !fileExists(dstPath) {
		if err := downloadURL(url, dstPath); err != nil {
			return channelID + " ERR", err
		}

		if err := mogrifyImage(dstPath); err != nil {
			return channelID + " ERR", err
		}
	}

	targetURL := serverPrefix + "/" + filename
	return fmt.Sprintf(`%s OK rewrite-url="%s"`, channelID, targetURL), nil
}

var invalidCharacter = regexp.MustCompile(`[^\w.]`)

func cleanFilename(name string) string {
	return invalidCharacter.ReplaceAllString(name, "")
}

func downloadURL(url *url.URL, dst string) error {
	req, _ := http.NewRequest("GET", url.String(), nil)
	req.Header.Set("User-Agent", userAgent)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(dst, buf, 0644)
}

func mogrifyImage(path string) error {
	cmd := exec.Command("/usr/local/bin/mogrify", "-flip", path)
	return cmd.Run()
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
