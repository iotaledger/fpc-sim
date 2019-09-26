#!/bin/bash

echo "Building executables..."

GOOS=darwin GOARCH=386 go build -o sim
mv sim ./bin/sim_osx
echo "MAC OSX version created"
GOOS=windows GOARCH=386 go build -o sim
mv sim ./bin/sim_win.exe
echo "Windows version created"
GOOS=linux GOARCH=386 go build -o sim
mv sim ./bin/sim_linux
echo "Linux version created"

echo "All done!"
