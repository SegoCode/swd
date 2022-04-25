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
go run swd.go -url="https://steamcommunity.com/sharedfiles/filedetails/?id=1111111111"
go run swd.go -id=1111111111
```
Or better [donwload a binary](https://github.com/SegoCode/swd/releases).

## Parameters

Use -h or --help flags to get  help for the program.
```shell
swd -help
```

Use -url for input the url of the steam workshop article you want to download.
```shell
swd -url="https://steamcommunity.com/sharedfiles/filedetails/?id=1111111111"
```

Use -id for input the id of the steam workshop article you want to download.
```shell
swd -id=1111111111
```

-format (optional) for choosing the format
```shell
swd -id=1111111111 -format=gmaextract
```

-node (optional) for choosing the node server to download ( between 4 and 8 )
```shell
swd -id=1111111111 -node=6
```

## Downloads

https://github.com/SegoCode/swd/releases/
