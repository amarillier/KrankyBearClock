#! /bin/sh

# fast update fyne and compile
go get fyne.io/fyne/v2@latest # or a specific version like @v2.4.0
go mod tidy
go mod vendor

go build -trimpath -ldflags="-w -s" -o KrankyBearClock .
./setIcon.sh Resources/Images/KrankyBearBeret.png KrankyBearClock

killall clock
rm ~/bin/clock
cp KrankyBearClock ~/bin/clock