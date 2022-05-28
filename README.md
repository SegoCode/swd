### ⚠️ 28/05/2022 - doesn't work anymore. **This repo is no longer maintained**, [API has been shut down.](https://www.reddit.com/r/swd_io/comments/uy55qg/we_are_no_longer_serving_any_files_through_our/)

<details>
  <summary>We are no longer serving any files through our network</summary> 
  <img  src="https://raw.githubusercontent.com/SegoCode/swd/main/media/We are no longer serving any files through our network. swd_io.png">
</details>


# swd 
<img  src="https://raw.githubusercontent.com/SegoCode/swd/main/media/demo1.1.gif">

Easy to use, simply wrapper for [steamworkshopdownloader.io](https://steamworkshopdownloader.io/) API.

## Usage & info

swd downloads the files compressed in a zip file with the steam workshop article id.

```shell
id_steamworkshop.zip
```

run from source code (Golang installation required).

```shell
git clone https://github.com/SegoCode/swd
cd swd
go get -d ./...
go run swd.go https://steamcommunity.com/sharedfiles/filedetails/?id=1111111111
```
Or better [donwload a binary](https://github.com/SegoCode/swd/releases).

## Parameters

It's simple, there is only one parameter, the url of the steam workshop article you want to download.
```shell
swd https://steamcommunity.com/sharedfiles/filedetails/?id=1111111111
```

But... steamworkshopdownloader.io has an optional parameter for download, if you know any you can specify 
```shell
swd https://steamcommunity.com/sharedfiles/filedetails/?id=1111111111 --downloadFormat gmaextract
```

## Downloads

https://github.com/SegoCode/swd/releases/
