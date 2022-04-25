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
	"flag"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/parnurzeal/gorequest"
	"github.com/tcnksm/go-latest"
	"github.com/tidwall/gjson"
)

const INFO = 1
const WARNING = 2
const ERR = 3

// START_NODE and END_NODE get from (steamworkshopdownloader.io)
const DEFAULT_NODE = 8
const START_NODE = 4  
const END_NODE = 8

func GetENDPOINT(node int) string {
	var ENDPOINT string = "https://node0" + strconv.Itoa(node) + ".steamworkshopdownloader.io/prod//api/"
	return ENDPOINT
}

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

func getUUID(api string, publishedFileId string, downloadFormat string) string {
	logger("CHEKING IF THE GAME IS AVAILABLE FOR STEAM WORKSHOP DOWNLOADS . . .", INFO)
	request := gorequest.New()
	resp, body, errs := request.Post(api).
		Set("Content-Type", "application/json").
		Send(`{"publishedFileId":` + publishedFileId + `, "collectionId":null, "hidden":false, "downloadFormat":"` + downloadFormat + `", "autodownload":true}`).
		End()

	if errs != nil {
		logger("CAN'T CONNECT TO THE SERVER", WARNING)
		return "0"
	} else {
		if resp.StatusCode != 200 {
			logger("GAME NOT AVAILABLE OR SERVER IS DOWN, CODE RESPONSE: "+strconv.Itoa(resp.StatusCode), WARNING)
			return "0"
		} else {
			logger("GAME IS AVAILABLE FOR STEAM WORKSHOP DOWNLOADS", INFO)
			return body
		}
	}
}


func main() {

	// Get Args //
	var fileUrl string
    flag.StringVar(&fileUrl, "url", "", "Url of file in steam workshop")
	var fileId string
	flag.StringVar(&fileId, "id", "", "Published file id of the file in steam workshop")
	var downloadFormat string
	flag.StringVar(&downloadFormat, "format", "raw", "Download format")
	var node int
	flag.IntVar(&node, "node", DEFAULT_NODE, "Server node (default: 8)")
	flag.Parse()


	// Validation Args //
	if fileUrl == "" && fileId == "" {
		logger("NEED A FILE URL OR PUBLISHED FILE ID, -help for usage", ERR)
	}

	if fileId == "" {
		fileUrl, err := url.ParseRequestURI(fileUrl)
		if err != nil {
			logger("URL NOT VALID (Example: swd https://steamcommunity.com/sharedfiles/filedetails/?id=1111111111)", ERR)
		}
		fileId = fileUrl.Query().Get("id")
	}

	if node < START_NODE || node > END_NODE {
		logger("NODE NOT VALID (Node must be between 4 and 8)", ERR)
	}
		
	logger("FileId: " + fileId, INFO)
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
	var initResponse string
	var ENDPOINT string

	for i := node; i >= START_NODE; i-- {  // Node can be 4 to 8
		ENDPOINT = GetENDPOINT(i)
		initResponse = getUUID(ENDPOINT + "download/request", fileId, downloadFormat)
		logger("REQUESTING DOWNLOAD FROM NODE " + strconv.Itoa(i), INFO)
		if initResponse != "0" {
			break
		} else {
			logger("TRYING TO CONNECT TO NODE " + strconv.Itoa(i), INFO)
		}
	}

	// Download request //
	uid := gjson.Get(initResponse, "uuid").String()
	var readyFile = false
	var storageNode = ""
	var storagepath = ""
	request := gorequest.New()
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
		err := DownloadFile("https://"+storageNode+"/prod//storage/"+storagepath+"?uuid="+uid, dir+string(os.PathSeparator)+fileId+".zip")

		if err != nil {
			panic(err)
		} else {
			logger("✔️  DOWNLOAD FINISHED IN "+(dir+string(os.PathSeparator)+fileId+".zip"), INFO)
		}

	} else {
		logger("FAIL THE SERVER IS BUSY", ERR)
	}

}

