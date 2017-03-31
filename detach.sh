#!/bin/bash/ -e
go build
screen -dmS quikrent-go-server $HOME/go/src/github.com/ruellia/quikrent-go-server/quikrent-go-server
