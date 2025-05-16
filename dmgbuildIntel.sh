#! /bin/sh
# compile, then create a dmg package
# https://github.com/create-dmg/create-dmg

# go build .
GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -ldflags="-w -s" -o bin/MacOSAMD64/
# set executable icon
./setIcon.sh KrankyBearClock.png bin/MacOSAMD64/KrankyBearClock

cp bin/MacOSAMD64/KrankyBearClock KrankyBearClock.app/Contents/MacOS

test -f KrankyBearClockIntel.dmg && rm KrankyBearClockIntel.dmg
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
  "KrankyBearClockIntel.dmg" \
  "KrankyBearClock.app"
  # --add-file KrankyBearClock.app ./KrankyBearClock.app
  # "./"

# set dmg icon
./setIcon.sh KrankyBearClock.png KrankyBearClockIntel.dmg
if [ ! -d installers ]
then
  mkdir installers
fi
cp KrankyBearClockIntel.dmg installers
# cp KrankyBearClockIntel.dmg ~/OneDrive\ -\ KrankyBear\ Inc/Apps/

# "Now this is not the end. It is not even the beginning of the end. But it is, perhaps, the end of the beginning." Winston Churchill, November 10, 1942