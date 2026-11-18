//go:generate fyne bundle -o bundled.go -a assets/images/KrankyBear.png
//go:generate fyne bundle -o bundled.go -a assets/images/KrankyBearBeret.png
//go:generate fyne bundle -o bundled.go -a assets/images/KrankyBearBeanieMultiColor.png
//go:generate fyne bundle -o bundled.go -a assets/images/KrankyBearChristmas.png
//go:generate fyne bundle -o bundled.go -a assets/images/KrankyBearChristmasGrinch.png
//go:generate fyne bundle -o bundled.go -a assets/images/KrankyBearCowboyBrown.png
//go:generate fyne bundle -o bundled.go -a assets/images/KrankyBearFedoraRed.png
//go:generate fyne bundle -o bundled.go -a assets/images/KrankyBearHardHat.png
//go:generate fyne bundle -o bundled.go -a assets/images/KrankyBearHogwartsSorting.png
//go:generate fyne bundle -o bundled.go -a assets/images/KrankyBearTrapperRedPlaid.png
//go:generate fyne bundle -o bundled.go -a assets/images/KrankyBearVikingHelmet.png
//go:generate fyne bundle -o bundled.go -a assets/images/http418.png

///go:generate fyne bundle -o bundled.go -a assets/images/Tcn256.png
///go:generate fyne bundle -o bundled.go -a assets/images/Tcc256.png
///go:generate fyne bundle -o bundled.go -a assets/images/Tcbe256.png

///go:generate fyne bundle -o bundled.go -a assets/sounds/boing.mp3
///go:generate fyne bundle -o bundled.go -a assets/sounds/Basso.mp3
///go:generate fyne bundle -o bundled.go -a assets/sounds/Sosumi.mp3
///go:generate fyne bundle -o bundled.go -a assets/sounds/Submarine.mp3
///go:generate fyne bundle -o bundled.go -a assets/sounds/baseball.mp3
///go:generate fyne bundle -o bundled.go -a assets/sounds/grandfatherClock.mp3
///go:generate fyne bundle -o bundled.go -a assets/sounds/pinball.mp3

package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type appTheme struct {
	fyne.Theme
}

func (a *appTheme) Size(n fyne.ThemeSizeName) float32 {
	if n == theme.SizeNameHeadingText {
		return a.Theme.Size(n) * 1.5
	}

	return a.Theme.Size(n)
}

func themeTimer(text *widget.RichText, time int) {
	seg := text.Segments[0].(*widget.TextSegment)
	if time <= 30 {
		seg.Style.ColorName = theme.ColorNameError
	} else if time < 150 {
		seg.Style.ColorName = theme.ColorNameWarning
	} else {
		seg.Style.ColorName = theme.ColorNameSuccess
	}
	text.Refresh()
}

// "Now this is not the end. It is not even the beginning of the end. But it is, perhaps, the end of the beginning." Winston Churchill, November 10, 1942
