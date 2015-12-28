# mediamock

Mediamock provides mocking of folders and (image) files.

Current use case: The media/assets folder of an online store or content 
management system, containing all images, pdf, etc, can have a pretty huge
total byte size, up to several GB or even TB. Copying these files to your 
development environment takes a long time and consumes lots of precious 
disk space.

A full detailed description with examples on my blog [http://cyrillschumacher.com/projects/2015-12-28-mediamock/](http://cyrillschumacher.com/projects/2015-12-28-mediamock/).

## How to use server mode with my Online Store or CMS?

### Magento

Magento 1: Please install [https://github.com/SchumacherFM/mediamock-magento](https://github.com/SchumacherFM/mediamock-magento).

Magento 2: Please install [https://github.com/SchumacherFM/mediamock-magento2](https://github.com/SchumacherFM/mediamock-magento2) todo.

These modules disable the HDD file access for reading images. Writing still possible.

### TYPO3 / NEOS

Please help

### Drupal

Please help

### Hybris

Please help

### Shopware

Please help

### OXID

Please help

## Install

Download binaries for windows, linux and darwin (OSX) in the [release section](https://github.com/SchumacherFM/mediamock/releases).

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

Copyright (c) 2015-2016 Cyrill (at) Schumacher dot fm. All rights reserved. See LICENSE file.

[Cyrill Schumacher](https://github.com/SchumacherFM) - [My pgp public key](http://www.schumacher.fm/cyrill.asc)

Identicon code by: Copyright (c) 2013, Damian Gryski <damian@gryski.com>
