package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/parnurzeal/gorequest"
	"github.com/tidwall/gjson"
)

const INFO = 1
const WARNING = 2
const ERR = 3

const ENDPOINT string = "https://node02.steamworkshopdownloader.io/prod//api/"

func logger(text string, errorlevel int) {

	if errorlevel == INFO {
		fmt.Println("[" + color.CyanString("INFO") + "]  " + text)
	}

	if errorlevel == WARNING {
		fmt.Println("[" + color.YellowString("WARNING") + "]  " + text)
	}

	if errorlevel == ERR {
		fmt.Println("[" + color.RedString("ERR") + "]  " + text)
		os.Exit(1)
	}
}

func DownloadFile(url string, filepath string) error {

	loadSp := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	loadSp.Prefix = "[" + color.CyanString("INFO") + "]  " + "RECEIVING DATA: "
	loadSp.FinalMSG = "\033[F"

	loadSp.Start()

	out, err := os.Create(filepath) // Create the file
	defer out.Close()
	if err != nil {
		return err
	}

	resp, err := http.Get(url) // Get the data
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	_, err = io.Copy(out, resp.Body) // Write the body to file
	if err != nil {
		return err
	}

	loadSp.Stop()

	return nil
}

func main() {

	// Args validation //
	if len(os.Args) <= 1 {
		logger("USAGE: swd https://steamcommunity.com/sharedfiles/filedetails/?id=1111111111", ERR)
	}

	url, err := url.ParseRequestURI(os.Args[1])
	if err != nil {
		logger("URL NOT VALID (Example: swd https://steamcommunity.com/sharedfiles/filedetails/?id=1111111111)", ERR)
	}

	idUrl := url.Query().Get("id")
	if idUrl == "" {
		logger("URL NOT VALID (Example: swd \"https://steamcommunity.com/sharedfiles/filedetails/?id=1111111111\")", ERR)
	}
	// End Args validation //

	// Get initial request //
	logger("CHEKING IF THE GAME IS AVAILABLE FOR STEAM WORKSHOP DOWNLOADS . . .", INFO)
	request := gorequest.New()
	resp, body, _ := request.Post(ENDPOINT+"download/request").
		Set("Content-Type", "application/json").
		Send(`{"publishedFileId":` + idUrl + `, "collectionId":0, "extract":true, "hidden":false, "direct":false, "autodownload":true}`).
		End()

	if resp.StatusCode != 200 {
		logger("GAME NOT AVAILABLE OR SERVER IS DOWN", ERR)
	} else {
		logger("GAME IS AVAILABLE FOR STEAM WORKSHOP DOWNLOADS", INFO)
	}

	// Download request //
	uid := gjson.Get(body, "uuid").String()
	var readyFile = false

	for i := 0; i < 10; i++ { // Try 10 times for 2 seconds of waiting, total 20 seconds of preparation maximum
		_, body, _ := request.Post(ENDPOINT+"download/status").
			Set("Content-Type", "application/json").
			Send(`{"uuids": ["` + uid + `"]}`).
			End()

		logger("WAITING FOR THE SERVER. . . DOWNLOAD STATUS: "+strings.ToUpper(gjson.Get(body, uid+".status").String()), INFO)

		if strings.Contains(body, "prepared") {
			readyFile = true
			logger("INITIATING DOWNLOADING. . . ", INFO)
			break
		}
		time.Sleep(2 * time.Second)
	}

	// File ready, start download //
	if readyFile {
		dir, _ := os.Getwd()
		err := DownloadFile(ENDPOINT+"download/transmit?uuid="+uid, dir+string(os.PathSeparator)+idUrl+".zip")

		if err != nil {
			panic(err)
		} else {
			logger("✔️ DOWNLOAD FINISHED IN "+(dir+string(os.PathSeparator)+idUrl+".zip"), INFO)
		}

	} else {
		logger("FAIL THE SERVER IS BUSY", ERR)
	}

}
