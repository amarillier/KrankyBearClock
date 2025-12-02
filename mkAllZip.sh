#! /bin/sh

version="0.4.4"
cp README.md installers/KrankyBearClock/Resources
cp ReleaseNotes.txt installers/KrankyBearClock/Resources
cd installers || exit
if [ ! -d KrankyBearClock ]
then
    mkdir -p KrankyBearClock
fi
rm KrankyBearClock/KrankyBearClock*

rm KrankyBearClock/clock.exe
cp ../bin/KrankyBearClock-windows.exe KrankyBearClock/clock.exe
zip -r KrankyBearClockWinAMD64.zip KrankyBearClock
rm KrankyBearClock/clock.exe


cp ../bin/KrankyBearClock-macos-amd64 KrankyBearClock/clock
zip -r KrankyBearClockMacOSAMD64.zip KrankyBearClock
rm KrankyBearClock/clock

cp ../bin/KrankyBearClock-macos-arm64 KrankyBearClock/clock
zip -r KrankyBearClockMacOSARM64.zip KrankyBearClock
rm KrankyBearClock/clock

cp ../bin/KrankyBearClock-linux-amd64 KrankyBearClock/clock
zip -r KrankyBearClockLinuxAMD64.zip KrankyBearClock
rm KrankyBearClock/clock

cp ../bin/KrankyBearClock-linux-arm64 KrankyBearClock/clock
zip -r KrankyBearClockLinuxARM64.zip KrankyBearClock
rm KrankyBearClock/clock

exit
# see gh docs: https://cli.github.com/manual/gh_release_create
awk '/0.4.4/{flag=1}/^$/{flag=0}flag' ../ReleaseNotes.txt > latestReleaseNotes.txt
gh release create --title v"$version" v"$version" --draft --notes-file latestReleaseNotes.txt --prerelease KrankyBearClockLinuxAMD64.zip KrankyBearClock_0.4.4-1_aarch64.rpm KrankyBearClockLinuxARM64.zip KrankyBearClock_0.4.4-1_amd64.deb KrankyBearClockMacOSAMD64.zip KrankyBearClock_0.4.4-1_amd64.pkg KrankyBearClockMacOSARM64.zip KrankyBearClock_0.4.4-1_arm64.deb KrankyBearClockSetup.exe KrankyBearClock_0.4.4-1_arm64.pkg KrankyBearClockWinAMD.zip KrankyBearClock_0.4.4-1_x86_64.rpm

echo "Created draft release $version"
echo "Remember to publish when ready"
echo "gh release edit v$version --draft=false --prerelease=false"