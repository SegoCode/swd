⚠️ 6/11/2021 - doesn't work any more. This repo is no longer maintained, see this [issue](https://github.com/SegoCode/swd/issues/2)

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
go run swd.go https://steamcommunity.com/sharedfiles/filedetails/?id=...
```
Or better [donwload a binary](https://github.com/SegoCode/swd/releases).

## Parameters

It's simple, there is only one parameter, the url of the steam workshop article you want to download.
```shell
swd https://steamcommunity.com/sharedfiles/filedetails/?id=......
```

## Downloads

https://github.com/SegoCode/swd/releases/
