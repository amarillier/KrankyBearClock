#! /bin/sh
# compile, then create a dmg package
# https://github.com/create-dmg/create-dmg

# go build .
# GOOS=linux GOARCH=arm64 go build -ldflags="-w -s" -o bin/MacOSARM64/
GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 go build -ldflags="-w -s" -o bin/MacOSARM64/

# cp TaniumClock TaniumClock.app/Contents/MacOS
cp bin/MacOSARM64/TaniumClock TaniumClock.app/Contents/MacOS
test -f TaniumClockARM.dmg && rm TaniumClockARM.dmg
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
  "TaniumClockARM.dmg" \
  "TaniumClock.app"
  # --add-file TaniumClock.app ./TaniumClock.app
  # "./"

  ./setIcon.sh TaniumClock.png TaniumClockARM.dmg
   cp TaniumClockARM.dmg installers