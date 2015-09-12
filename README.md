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

## Future TODOs

Run as server and allow clients to trigger the analyze steop
once connected to the server.

## Command line

```
Usage: mediamock options...

Options:
  -i  Read CSV data from this input URL/file.
  -d  Read this directory recursively and write into -o. If -i is provided
      generate all mocks in this directory. Default: current directory.
  -o  Write data into out file (optional, default a temp file).
  -p  Image pattern: happy (default), warm or rand
```

### Run analyze

The following is an example output:

```
$ ./mediamock -d ~/Sites/magento19-data/media
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
$ ./mediamock -d pathToMyDevFolder -i https://www.onlineshop.com.au/mediamock.csv.gz -p warm
Directory pathToMyDevFolder created
241 / 9957 [=====>------------------------------------------------------------------] 2.42 % 9m32s
```

## Install

Download binaries in the release section or 
`go get -u github.com/SchumacherFM/mediamock` or
`go install github.com/SchumacherFM/mediamock`

## License

Copyright (c) 2015 Cyrill (at) Schumacher dot fm. All rights reserved.

[Cyrill Schumacher](https://github.com/SchumacherFM) - [My pgp public key](http://www.schumacher.fm/cyrill.asc)
