#! /bin/sh

if [ $# -ge 1 ]
then
    ver=$1
else
    echo "Enter a version number"
    cur=$(cat main.go | grep -i "Version =")
    echo "    current: $cur"
    read ver
    if [ -z "$ver" ]
    then
        echo "Enter a version!"
        exit
    else
        echo "Version: $ver"
        # exit
    fi
fi

echo "version: $ver"
echo "main.go"
sed -i '' "s/Version = \".*\"/Version = \"$ver\"/" main.go

echo "FyneApp.toml"
sed -i '' "s/Version = \".*\"/Version = \"$ver\"/" FyneApp.toml

echo "Inno Setup Inno/KrankyBearClock.iss"
sed -i '' "s/MyAppVersion \".*\"/MyAppVersion \"$ver\"/" ./Inno/KrankyBearClock.iss

echo "Inno Setup winres/winres.json"
sed -i '' "s/file_version\":.*/file_version\": \"$ver\",/" ./winres/winres.json
sed -i '' "s/product_version\":.*/product_version\": \"$ver\"/" ./winres/winres.json
sed -i '' "s/FileVersion\":.*/FileVersion\": \"$ver\",/" ./winres/winres.json
sed -i '' "s/ProductVersion\":.*/ProductVersion\": \"$ver\",/" ./winres/winres.json
