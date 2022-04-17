package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/parnurzeal/gorequest"
	"github.com/tcnksm/go-latest"
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
		fmt.Println("[" + color.YellowString("WARN") + "]  " + text)
	}

	if errorlevel == ERR {
		fmt.Println("[" + color.RedString("ERR") + "]   " + text)
		os.Exit(1)
	}
}

func DownloadFile(url string, filepath string) error {

	loadSp := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	loadSp.Prefix = "[" + color.CyanString("INFO") + "]  " + "RECEIVING DATA: "
	loadSp.FinalMSG = "\033[F"

	loadSp.Start()

	out, err := os.Create(filepath) // Create the file
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url) // Get the data
	if err != nil {
		return err
	}
	defer resp.Body.Close()

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

	var downloadFormat = "raw"
	if len(os.Args) >= 3 && (os.Args[2] == "--downloadFormat") {
		downloadFormat = os.Args[3]
	}
	// End Args validation //

	githubTag := &latest.GithubTag{
		Owner:      "SegoCode",
		Repository: "swd",
	}

	res, err := latest.Check(githubTag, "1.6.0")
	if err == nil {
		if res.Outdated {
			logger("NEW VERSION IS AVAILABLE, CHECK https://github.com/SegoCode/swd/releases", WARNING)
		}
	} else {
		logger("CAN'T CHECK THE LATEST VERSION IN GITHUB, CHECK https://github.com/SegoCode/swd/releases", WARNING)
	}

	// Get initial request //
	logger("CHEKING IF THE GAME IS AVAILABLE FOR STEAM WORKSHOP DOWNLOADS . . .", INFO)
	request := gorequest.New()
	resp, body, errs := request.Post(ENDPOINT+"download/request").
		Set("Content-Type", "application/json").
		Send(`{"publishedFileId":` + idUrl + `, "collectionId":null, "hidden":false, "downloadFormat":"` + downloadFormat + `", "autodownload":true}`).
		End()

	if errs != nil {
		logger("CAN'T CONNECT TO THE SERVER, EXITING . . .", ERR)
	} else {
		if resp.StatusCode != 200 {
			logger("GAME NOT AVAILABLE OR SERVER IS DOWN, CODE RESPONSE: "+strconv.Itoa(resp.StatusCode), ERR)
		} else {
			logger("GAME IS AVAILABLE FOR STEAM WORKSHOP DOWNLOADS", INFO)
		}
	}

	// Download request //
	uid := gjson.Get(body, "uuid").String()
	var readyFile = false
	var storageNode = ""
	var storagepath = ""
	for i := 0; i < 10; i++ { // Try 10 times for 2 seconds of waiting, total 20 seconds of preparation maximum
		_, body, _ := request.Post(ENDPOINT+"download/status").
			Set("Content-Type", "application/json").
			Send(`{"uuids": ["` + uid + `"]}`).
			End()

		logger("WAITING FOR THE SERVER. . . DOWNLOAD STATUS: "+strings.ToUpper(gjson.Get(body, uid+".status").String()), INFO)

		if strings.Contains(body, "prepared") {
			readyFile = true
			storageNode = gjson.Get(body, uid+".storageNode").String()
			storagepath = gjson.Get(body, uid+".storagePath").String()
			logger("INITIATING DOWNLOADING. . . ", INFO)
			break
		}
		time.Sleep(2 * time.Second)
	}

	// File ready, start download //
	if readyFile {
		dir, _ := os.Getwd()
		err := DownloadFile("https://"+storageNode+"/prod//storage/"+storagepath+"?uuid="+uid, dir+string(os.PathSeparator)+idUrl+".zip")

		if err != nil {
			panic(err)
		} else {
			logger("✔️  DOWNLOAD FINISHED IN "+(dir+string(os.PathSeparator)+idUrl+".zip"), INFO)
		}

	} else {
		logger("FAIL THE SERVER IS BUSY", ERR)
	}

}
