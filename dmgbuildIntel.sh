#! /bin/sh
# compile, then create a dmg package
# https://github.com/create-dmg/create-dmg

# go build .
GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -ldflags="-w -s" -o bin/MacOSAMD64/

cp bin/MacOSAMD64/TaniumClock TaniumClock.app/Contents/MacOS

test -f TaniumClockIntel.dmg && rm TaniumClockIntel.dmg
#   --volicon "TaniumClock.icns" \
create-dmg \
  --volname "TaniumClock" \
  --window-pos 200 120 \
  --window-size 800 400 \
  --icon-size 100 \
  --icon "TaniumClock.app" 200 190 \
  --hide-extension "TaniumClock.app" \
  --app-drop-link 600 185 \
  --eula license.txt \
  "TaniumClockIntel.dmg" \
  "TaniumClock.app"
  # --add-file TaniumClock.app ./TaniumClock.app
  # "./"

./setIcon.sh TaniumClock.png TaniumClockIntel.dmg
cp TaniumClockIntel.dmg installers