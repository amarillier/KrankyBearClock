#! /bin/sh
# compile, then create a dmg package
# https://github.com/create-dmg/create-dmg

# go build .
# GOOS=linux GOARCH=arm64 go build -ldflags="-w -s" -o bin/MacOSARM64/
GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 go build -ldflags="-w -s" -o bin/MacOSARM64/
# set executable icon
./setIcon.sh KrankyBear.png bin/MacOSARM64/KrankyBearClock

# cp KrankyBearClock KrankyBearClock.app/Contents/MacOS
cp bin/MacOSARM64/KrankyBearClock KrankyBearClock.app/Contents/MacOS
test -f KrankyBearClockARM.dmg && rm KrankyBearClockARM.dmg
#   --volicon "KrankyBearClock.icns" \
create-dmg \
  --volname "KrankyBearClock" \
  --window-pos 200 120 \
  --window-size 800 400 \
  --icon-size 100 \
  --icon "KrankyBearClock.app" 200 190 \
  --hide-extension "KrankyBearClock.app" \
  --app-drop-link 600 185 \
  --eula license.txt \
  "KrankyBearClockARM.dmg" \
  "KrankyBearClock.app"
  # --add-file KrankyBearClock.app ./KrankyBearClock.app
  # "./"

# set dmg icon
./setIcon.sh KrankyBearClock.png KrankyBearClockARM.dmg
if [ ! -d installers ]
then
  mkdir installers
fi
cp KrankyBearClockARM.dmg installers
# cp KrankyBearClockARM.dmg ~/OneDrive\ -\ KrankyBear\ Inc/Apps/

# "Now this is not the end. It is not even the beginning of the end. But it is, perhaps, the end of the beginning." Winston Churchill, November 10, 1942