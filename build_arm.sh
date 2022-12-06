#!/bin/bash
cd src/run && GOPATH=`pwd`;GOARM=7;GOOS=linux;GOARCH=arm go build -v -o ../../smartmeter_reader_arm