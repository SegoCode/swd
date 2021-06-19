# swd
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

It's simple, there is only one parameter, the url of the thread you want to download.
```shell
swd https://steamcommunity.com/sharedfiles/filedetails/?id=......
```

## Downloads

https://github.com/SegoCode/swd/releases/
