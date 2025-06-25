#! /bin/sh

version="0.4.2"
cp README.md installers/KrankyBearClock/Resources
cp ReleaseNotes.txt installers/KrankyBearClock/Resources
cd installers || exit
cp ../bin/WinAMD64/KrankyBearClock.exe KrankyBearClock
zip -r KrankyBearClockWinAMD.zip KrankyBearClock
rm KrankyBearClock/KrankyBearClock.exe

cp ../bin/MacOSAMD64/KrankyBearClock KrankyBearClock
zip -r KrankyBearClockMacOSAMD.zip KrankyBearClock
rm KrankyBearClock/KrankyBearClock

cp ../bin/MacOSARM64/KrankyBearClock KrankyBearClock
zip -r KrankyBearClockMacOSARM.zip KrankyBearClock
rm KrankyBearClock/KrankyBearClock

# see gh docs: https://cli.github.com/manual/gh_release_create
awk '/0.4.2/{flag=1}/^$/{flag=0}flag' ../ReleaseNotes.txt > latestReleaseNotes.txt
gh release create --title v"$version" v"$version" --draft --notes-file latestReleaseNotes.txt --prerelease KrankyBearClockWinAMD.zip KrankyBearClockMacOSAMD.zip KrankyBearClockMacOSARM.zip KrankyBearClockSetup.exe KrankyBearClockARM.dmg KrankyBearClockIntel.dmg

echo "Created draft release $version"
echo "Remember to publish when ready"
echo "gh release edit v$version --draft=false --prerelease=false"