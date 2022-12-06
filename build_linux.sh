#!/bin/bash
cd src/run && GOPATH=`pwd`;GOARM=7;GOOS=linux;GOARCH=amd64 go build -v -o ../../smartmeter_reader_linux