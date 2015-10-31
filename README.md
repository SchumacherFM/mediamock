# mediamock

Mediamock provides mocking of folders and (image) files.

Current use case: The media/assets folder of an online store or content 
management system, containing all images, pdf, etc, can have a pretty huge
total byte size, up to several GB or even TB. Copying these files to your 
development environment takes a long time and consumes lots of precious 
disk space.

## Mediamock has two run modes:

1. First analyze and store the folder structure on the server.
2. Download the stored structure onto your development machine
and recreate the folder and files.

### Run mode: Analyze

The program will recursively walk through the folders and stores
each file including path, image width + height and modification 
date in a simple CSV file.

### Run mode: Mock

The mock mode will read the CSV file from your hard drive or via HTTP and 
creates all the folders and files including correct modification date for the
files. For images it creates an uniform colored image with the correct width 
and height.

The image contains an uniform color in a random, warm or happy tone.

Supported image formats: png, jpeg and gif.

The mocked images should be as small as possible. All other non-image
files are of size 0kb.

### Run mode: Server

Same as Mock, but generates the images on-the-fly.

## How to use server mode with my Online Store or CMS?

### Magento

Please install [https://github.com/SchumacherFM/mediamock-magento](https://github.com/SchumacherFM/mediamock-magento).

This module disables HDD file access in Magento.

### TYPO3

Please help

### Drupal

Please help

### Hybris

Please help

### Shopware

Please help

### OXID

Please help

## Command line

```
$ ./mediamock
NAME:
   mediamock - reads your assets/media directory on your server and
               replicates that structure on your development machine.

               $ mediamock help analyze|mock|server for more options!


USAGE:
   mediamock [global options] command [command options] [arguments...]

VERSION:
   v0.1.0 by @SchumacherFM

COMMANDS:
   analyze, a	Analyze the directory structure on you production server and write into a
		csv.gz file.
   mock, m	Mock reads the csv.gz file and recreates the files and folders. If a file represents
		an image, it will be created with a tiny file size and correct width x height.
   server, s	Server reads the csv.gz file and creates the assets/media structure on the fly
		as a HTTP server. Does not write anything to your hard disk. Open URL / on the server
		to retrieve a list of all files and folders.
   help, h	Shows a list of commands or help for one command

GLOBAL OPTIONS:
   -p "happy"		Image pattern: happy, warm, rand, happytext, warmtext, HTML hex value or icon
   --help, -h		show help
   --version, -v	print the version
```

If the option `p` contains the word `text` like h`appytext` or `warmtext` then the image file name
will be printed all over the generate image. This is useful if you need to inspect a zooming
effect on the frontend.

Pattern *icon* generates identicons like introduced here [https://github.com/blog/1586-identicons](https://github.com/blog/1586-identicons).
The pattern will be generated from the file name.

### Run analyze

```
$ ./mediamock help a
NAME:
   analyze - Analyze the directory structure on you production server and write into a
	csv.gz file.

USAGE:
   command analyze [command options] [arguments...]

OPTIONS:
   -d "."						Read this directory recursively and write into -o
   -o "/tmp/mediamock.csv.gz"	Write to this output file.
```

The following is an example output:

```
$ ./mediamock analyze -d ~/Sites/magento19-data/media
Image ~/Sites/magento19-data/media/catalog/product/6/4/64406_66803218048_1831204_n.jpg decoding error: image: unknown format
Image ~/Sites/magento19-data/media/catalog/product/6/4/64406_66803218048_1813204_n_1.jpg decoding error: image: unknown format
Image ~/Sites/magento19-data/media/catalog/product/i/m/IMG_1658_4.png decoding error: image: unknown format
Image ~/Sites/magento19-data/media/catalog/product/i/m/IMG_7445.JPG decoding error: image: unknown format
Image ~/Sites/magento19-data/media/catalog/product/i/m/IMG_7450_1.JPG decoding error: image: unknown format
Image ~/Sites/magento19-data/media/catalog/product/i/m/IMG_7450_2.JPG decoding error: image: unknown format
Image ~/Sites/magento19-data/media/catalog/product/p/h/photo1_4.JPG decoding error: image: unknown format
Wrote to file: /var/folders/bp/b3k4vgcd716fgxqtm_rlvn5m0100gn/T/mediamock.csv.gz
```

Concerns: You must download the mediamock binary onto your server and execute 
it from there on the command line. 

### Run mock

```
$ ./mediamock help m
NAME:
   mock - Mock reads the csv.gz file and recreates the files and folders. If a file represents
	an image, it will be created with a tiny file size and correct width x height.

USAGE:
   command mock [command options] [arguments...]

OPTIONS:
   -i 		Read csv.gz from this input file or input URL.
   -d "."	Create structure in this directory.
```

```
$ ./mediamock m -d pathToMyDevFolder -i https://www.onlineshop.com.au/mediamock.csv.gz -p warm
Directory pathToMyDevFolder created
241 / 9957 [=====>------------------------------------------------------------------] 2.42 % 9m32s
```

### Run server

```
$ ./mediamock help s
NAME:
   server - Server reads the csv.gz file and creates the assets/media structure on the fly
	as a HTTP server. Does not write anything to your hard disk. Open URL / on the server
	to retrieve a list of all files and folders.

USAGE:
   command server [command options] [arguments...]

OPTIONS:
   --urlPrefix 		       Prefix in the URL path
   -i 				       Read csv.gz from this input file or input URL.
   --host "127.0.0.1:4711" IP address or host name
```

```
$ mediamock s -i /tmp/mediamock.csv.gz -urlPrefix media/
```

For detailed statistic about the running mediamock server internals you can visit:

- `http://localhost:4711/` Tiny directory index
- `http://localhost:4711/debug/charts/` Garbage Collection and Memory Usage statistics
- `http://localhost:4711/json` all files as a JSON stream
- `http://localhost:4711/html` all files as a HTML table

![ScreenShot](/debugCharts.png)

You can retrieve a list of all served files by navigating to `http://127.0.0.1:4711`.

If an image doesn't exists in the CSV file but is requested from the front end
mediamock will try to generated an appropriate image if it can detect width and height
information within the URL.

E.g.: `http://127.0.0.1:4711/media/catalog/product/cache/2/small_image/218x258/9df78eab33525d08d6e5fb8d27136e95/detail/myImage.jpg`
mediamock can detect that this image is 218px x 258px in size because URL mentions 218x258.

## TODO

- Use concurrent file walk: [https://github.com/MichaelTJones/walk](https://github.com/MichaelTJones/walk)

## Install

Download binaries for windows, linux and darwin in the [release section](https://github.com/SchumacherFM/mediamock/releases).

## Contribute

`GO15VENDOREXPERIMENT` introduces reproduceable builds. 

```
$ go get -u -v github.com/SchumacherFM/mediamock/...
$ cd $GOPATH/src/github.com/SchumacherFM/mediamock
$ git remote rm origin
$ git remote add origin git@github.com:username/CloneOfMediaMock.git
$ git submodule init
$ git submodule update
hack hack hack ...
$ GO15VENDOREXPERIMENT=1 go run *.go
hack hack hack ...
$ GO15VENDOREXPERIMENT=1 go run *.go
$ gofmt -w *.go common/*.go record/*.go
$ git commit -a -m 'Add feature X including tests'
$ git push -u origin master
create pull request to github.com/SchumacherFM/mediamock
```

If you introduce a new dependency this is how to add it:

```
$ cd $GOPATH/src/github.com/SchumacherFM/mediamock
$ git submodule add git@github.com:username/GoLangRep.git vendor/github.com/username/GoLangRep
```

How do I know all dependencies?

```
$ go list -json github.com/SchumacherFM/mediamock/...
```

## License

Copyright (c) 2015 Cyrill (at) Schumacher dot fm. All rights reserved.

[Cyrill Schumacher](https://github.com/SchumacherFM) - [My pgp public key](http://www.schumacher.fm/cyrill.asc)

Identicon code by: Copyright (c) 2013, Damian Gryski <damian@gryski.com>
