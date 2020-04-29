#!/bin/sh
echo "Building 'Moodler Online Bot x86' for Windows..."
GOOS=windows GOARCH=386 go build -o online.moodler/bin/windows/86/online.moodler.exe online.moodler/source/online.moodler.go
echo "Building 'Moodler Online Bot x86' for Linux..."
GOOS=linux GOARCH=386 go build -o online.moodler/bin/linux/86/online.moodler.run online.moodler/source/online.moodler.go
if [ ! $((0xffffffff)) -eq -1 ]
then
    echo "Building 'Moodler Online Bot x64' for Windows..."
    GOOS=windows GOARCH=amd64 go build -o online.moodler/bin/windows/64/online.moodler.exe online.moodler/source/online.moodler.go
    echo "Building 'Moodler Online Bot x64' for Linux..."
    GOOS=linux GOARCH=amd64 go build -o online.moodler/bin/linux/64/online.moodler.run online.moodler/source/online.moodler.go
fi