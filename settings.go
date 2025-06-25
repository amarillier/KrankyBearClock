package main

import (
	"fmt"
	"image/color"
	"log"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/cmd/fyne_settings/settings"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/spiretechnology/go-autostart/v2"
)

var fileButton *widget.Button
var selectedFile *widget.Label
var fileURI fyne.URI
var tselect fyne.Window
var mycolor color.Color
var muteonbutton *widget.Button
var muteoffbutton *widget.Button
var muteonlabel string
var muteofflabel string

func makeSettingsClock(a fyne.App, w fyne.Window, bg fyne.Canvas) {
	// settings window
	if settingsc != nil { // &&  !settings.Content().Visible() {
		settingsc.Show()
		teapot(a, settingsc)
	} else {
		settingsc = a.NewWindow(appName + ": Settings")
		settingsc.SetIcon(resourceKrankyBearBeretPng)
		settingsText := `All updates are applied / saved immediately.
	Note: clock display settings do not currently auto refresh, restart is required.
	Displaying clock seconds can be much more CPU intensive than not!`
		setText := widget.NewLabel(settingsText)
		setText.TextStyle = fyne.TextStyle{Bold: true}

		todoText := `Still to be added: 
	font type selection
	allow .mid and .wav sounds
	background color or selectable background images in addition to built in images`
		doText := widget.NewLabel(todoText)
		doText.TextStyle = fyne.TextStyle{Italic: true, Bold: true}

		mp3files, err := listMatchingFiles(sndDir, "*.mp3")
		if err != nil {
			log.Fatal(err)
		}
		mp3 := []string{"ding", "down", "up", "updown"}
		for _, file := range mp3files {
			mp3 = append(mp3, file)
		}

		showsec := widget.NewCheck("", func(value bool) {
			if debug == 1 {
				log.Println("showseconds set to", value)
			}
			switch value {
			case true:
				showseconds = 1
				// clock.Content().Refresh()
				if debug == 1 {
					fmt.Println("showseconds on")
				}
			case false:
				showseconds = 0
				// clock.Content().Refresh()
				if debug == 1 {
					fmt.Println("showseconds off")
				}
			}
			a.Preferences().SetInt("showseconds.default", showseconds)
		})
		showdt := widget.NewCheck("", func(value bool) {
			if debug == 1 {
				log.Println("show date set to", value)
			}
			switch value {
			case true:
				showdate = 1
			case false:
				showdate = 0
			}
			a.Preferences().SetInt("showdate.default", showdate)
		})
		showtz := widget.NewCheck("", func(value bool) {
			if debug == 1 {
				log.Println("showtimezone set to", value)
			}
			switch value {
			case true:
				showtimezone = 1
			case false:
				showtimezone = 0
			}
			a.Preferences().SetInt("showtimezone.default", showtimezone)
		})
		showut := widget.NewCheck("", func(value bool) {
			if debug == 1 {
				log.Println("showutc set to", value)
			}
			switch value {
			case true:
				showutc = 1
			case false:
				showutc = 0
			}
			a.Preferences().SetInt("showutc.default", showutc)
		})
		showhr1224 := widget.NewRadioGroup([]string{"12", "24"}, func(value string) {
			if debug == 1 {
				log.Println("12 / 24 time set to", value)
			}
			switch value {
			case "12":
				showhr12 = 1
			case "24":
				showhr12 = 0
			}
			a.Preferences().SetInt("showhr12.default", showhr12)
		})
		showhr1224.Horizontal = true
		mute := widget.NewCheck("", func(value bool) {
			if debug == 1 {
				log.Println("automute set to", value)
			}
			switch value {
			case true:
				automute = 1
			case false:
				automute = 0
			}
			a.Preferences().SetInt("automute.default", automute)
		})
		chime := widget.NewCheck("", func(value bool) {
			if debug == 1 {
				log.Println("hourchime set to", value)
			}
			switch value {
			case true:
				hourchime = 1
			case false:
				hourchime = 0
			}
			a.Preferences().SetInt("hourchime.default", hourchime)
		})
		chimesound := widget.NewSelect(mp3, func(value string) {
			if debug == 1 {
				log.Println("chimesound set to", value)
			}
			hourchimesound = value // strings.Replace(value, "builtin ", "", 1)
			switch hourchimesound {
			case "up", "down", "updown", "ding":
				playBeep(hourchimesound) // built in sounds
			default:
				playMp3(sndDir + "/" + hourchimesound)
			}
			a.Preferences().SetString("hourchimesound.default", hourchimesound)
		})
		startatboot := widget.NewCheck("", func(value bool) {
			if debug == 1 {
				log.Println("startatboot set to", value)
			}
			autoClock := autostart.New(autostart.Options{
				Label:       "com.github.amarillier.KrankyBearClock",
				Name:        "KrankyBearClock",
				Description: "Kranky Bear Clock",
				Mode:        autostart.ModeUser,
				Arguments:   []string{},
			})
			switch value {
			case true:
				startclock = 1
				autoClock.Enable()
			case false:
				startclock = 0
				autoClock.Disable()
			}
			a.Preferences().SetInt("startclock.default", startclock)
		})
		lockmute := widget.NewCheck("", func(value bool) {
			if debug == 1 {
				log.Println("slockmute set to", value)
			}
			switch value {
			case true:
				slockmute = 1
			case false:
				slockmute = 0
			}
			a.Preferences().SetInt("slockmute.default", slockmute)
		})

		tsz := widget.NewEntry()
		tsz.SetText(strconv.Itoa(timesize))
		tsz.OnChanged = func(value string) {
			if debug == 1 {
				log.Println("time font size set to", value)
			}
			timesize, err = strconv.Atoi(value)
			if err != nil {
				playBeep("ding")
				tsz.SetText(strconv.Itoa(48))
			} else {
				switch {
				case timesize < 10:
					timesize = 10
					value = strconv.Itoa(10)
				case timesize > 200:
					timesize = 200
					value = strconv.Itoa(200)
				}
				tsz.SetText(strconv.Itoa(timesize))
				a.Preferences().SetInt("timesize.default", timesize)
			}
		}
		// Create buttons for increase and decrease
		tincrease := widget.NewButton("▲", func() {
			value, _ := strconv.Atoi(tsz.Text)
			if value < 200 {
				tsz.SetText(fmt.Sprintf("%d", value+1))
				timesize = value + 1
				a.Preferences().SetInt("timesize.default", timesize)
			} else {
				playBeep("ding")
			}
		})
		tdecrease := widget.NewButton("▼", func() {
			value, _ := strconv.Atoi(tsz.Text)
			if value > 10 {
				tsz.SetText(fmt.Sprintf("%d", value-1))
				timesize = value - 1
				a.Preferences().SetInt("timesize.default", timesize)
			} else {
				playBeep("ding")
			}
		})

		dsz := widget.NewEntry()
		dsz.SetText(strconv.Itoa(datesize))
		dsz.OnChanged = func(value string) {
			if debug == 1 {
				log.Println("date font size set to", value)
			}
			datesize, err = strconv.Atoi(value)
			if err != nil {
				playBeep("ding")
				tsz.SetText(strconv.Itoa(24))
			} else {
				switch {
				case datesize < 10:
					datesize = 10
					value = strconv.Itoa(10)
				case datesize > 200:
					datesize = 200
					value = strconv.Itoa(200)
				}
				dsz.SetText(strconv.Itoa(datesize))
				a.Preferences().SetInt("datesize.default", datesize)
			}
		}
		// Create buttons for increase and decrease
		dincrease := widget.NewButton("▲", func() {
			value, _ := strconv.Atoi(dsz.Text)
			if value < 200 {
				dsz.SetText(fmt.Sprintf("%d", value+1))
				datesize = value + 1
				a.Preferences().SetInt("datesize.default", datesize)
			} else {
				playBeep("ding")
			}
		})
		ddecrease := widget.NewButton("▼", func() {
			value, _ := strconv.Atoi(dsz.Text)
			if value > 10 {
				dsz.SetText(fmt.Sprintf("%d", value-1))
				datesize = value - 1
				a.Preferences().SetInt("datesize.default", datesize)
			} else {
				playBeep("ding")
			}
		})

		usz := widget.NewEntry()
		usz.SetText(strconv.Itoa(utcsize))
		usz.OnChanged = func(value string) {
			if debug == 1 {
				log.Println("utc font size set to", value)
			}
			utcsize, err = strconv.Atoi(value)
			if err != nil {
				playBeep("ding")
				usz.SetText(strconv.Itoa(18))
			} else {
				switch {
				case utcsize < 10:
					utcsize = 10
					value = strconv.Itoa(10)
				case utcsize > 200:
					utcsize = 200
					value = strconv.Itoa(200)
				}
				usz.SetText(strconv.Itoa(utcsize))
				a.Preferences().SetInt("utcsize.default", utcsize)
			}
		}
		// Create buttons for increase and decrease
		uincrease := widget.NewButton("▲", func() {
			value, _ := strconv.Atoi(usz.Text)
			if value < 200 {
				usz.SetText(fmt.Sprintf("%d", value+1))
				utcsize = value + 1
				a.Preferences().SetInt("utcsize.default", utcsize)
			} else {
				playBeep("ding")
			}
		})
		udecrease := widget.NewButton("▼", func() {
			value, _ := strconv.Atoi(usz.Text)
			if value > 10 {
				usz.SetText(fmt.Sprintf("%d", value-1))
				utcsize = value - 1
				a.Preferences().SetInt("utcsize.default", utcsize)
			} else {
				playBeep("ding")
			}
		})

		/*
			timefont
			datefont
			utcfont
		*/

		reset := widget.NewButton("Reset defaults", func() {
			if debug == 1 {
				log.Println("preferences reset to defaults")
			}
			writeDefaultSettings(a)
			showsec.SetChecked(false)
			showtz.SetChecked(true)
			showdt.SetChecked(true)
			showut.SetChecked(true)
			showhr1224.SetSelected("12")
			lockmute.SetChecked(false)
			mute.SetChecked(false)
			muteonhr = 20
			muteonmin = 0
			muteoffhr = 8
			muteoffmin = 0
			muteonlabel = fmt.Sprintf("%02d:%02d", muteonhr, muteonmin)
			muteofflabel = fmt.Sprintf("%02d:%02d", muteoffhr, muteoffmin)
			muteonbutton.SetText("Mute: " + muteonlabel)
			muteoffbutton.SetText("Unmute: " + muteofflabel)
			muteonbutton.Refresh()
			chime.SetChecked(true)
			hourchimesound = "cuckoo.mp3"
			chimesound.Selected = hourchimesound
			startatboot.SetChecked(false)
			showsec.Refresh()
			showtz.Refresh()
			showut.Refresh()
			showhr1224.Refresh()
			startatboot.Refresh()
			lockmute.Refresh()
			mute.Refresh()
			chime.Refresh()
			chimesound.Refresh()
			timesize = 48
			datesize = 24
			utcsize = 18
			tsz.SetText(strconv.Itoa(timesize))
			dsz.SetText(strconv.Itoa(datesize))
			usz.SetText(strconv.Itoa(utcsize))
		})
		reset.Importance = widget.SuccessImportance // green
		// reset.Resize(fyne.NewSize(reset.MinSize().Width, reset.MinSize().Height))
		close := widget.NewButton("Close settings", func() {
			settingsc.Close()
			settingsc = nil
		})
		close.Importance = widget.WarningImportance // orange
		buttonRow := container.NewCenter(container.NewHBox(container.NewCenter(reset), container.NewCenter(close)))

		if showseconds == 1 {
			showsec.SetChecked(true)
		} else {
			showsec.SetChecked(false)
		}
		if showtimezone == 1 {
			showtz.SetChecked(true)
		} else {
			showtz.SetChecked(false)
		}
		if showdate == 1 {
			showdt.SetChecked(true)
		} else {
			showdt.SetChecked(false)
		}
		if showutc == 1 {
			showut.SetChecked(true)
		} else {
			showut.SetChecked(false)
		}
		switch showhr12 {
		case 1:
			showhr1224.SetSelected("12")
		case 0:
			showhr1224.SetSelected("24")
		}
		if automute == 1 {
			mute.SetChecked(true)
		} else {
			mute.SetChecked(false)
		}
		if hourchime == 1 {
			chime.SetChecked(true)
		} else {
			chime.SetChecked(false)
		}
		chimesound.Selected = hourchimesound
		if startclock == 1 {
			startatboot.SetChecked(true)
		} else {
			startatboot.SetChecked(false)
		}
		if slockmute == 1 {
			lockmute.SetChecked(true)
		} else {
			lockmute.SetChecked(false)
		}

		/*
			background.Selected = timerbg
		*/
		setform := widget.NewForm(
			widget.NewFormItem("Show Seconds", showsec),
			widget.NewFormItem("Show Timezone", showtz),
			widget.NewFormItem("Show Date", showdt),
			widget.NewFormItem("Show UTC", showut),
			widget.NewFormItem("Show 12/24 Hour Time", showhr1224),
			widget.NewFormItem("Auto Start at Boot", startatboot),
			widget.NewFormItem("Hourly Chime", chime),
			widget.NewFormItem("Hourly Chime Sound", chimesound),
			widget.NewFormItem("Lock Mute Volume", lockmute),
			widget.NewFormItem("Auto Mute Volume", mute),
		)
		muteonlabel = fmt.Sprintf("%02d:%02d", muteonhr, muteonmin)
		muteonbutton = widget.NewButton("Mute: "+muteonlabel, func() {
			muteon := selectTime(a, w, bg, "muteon", muteonhr, muteonmin)
			if debug == 1 {
				log.Println("muteon set to", muteon)
			}
			muteonlabel = fmt.Sprintf("%02d:%02d", muteonhr, muteonmin)
			// muteonbutton.SetText("Mute: " + muteonlabel)
			// muteonbutton.Refresh()
		})
		muteofflabel := fmt.Sprintf("%02d:%02d", muteoffhr, muteoffmin)
		muteoffbutton = widget.NewButton("Unmute: "+muteofflabel, func() {
			muteoff := selectTime(a, w, bg, "muteoff", muteoffhr, muteoffmin)
			if debug == 1 {
				log.Println("muteoff set to", muteoff)
			}
			// a.Preferences().SetInt("muteoffhr.default", muteoffhr)
			// a.Preferences().SetInt("muteoffmin.default", muteoffmin)
		})
		mwidget := container.NewHBox(
			muteonbutton, muteoffbutton)
		tcbutton := widget.NewButton("Time Color", func() {
			tcolor := colorPicker(settingsc, "time", a)
			if debug == 1 {
				fmt.Println("tcolor:", tcolor)
			}
		})
		bgbutton := widget.NewButton("Background Color", func() {
			bcolor := colorPicker(settingsc, "background", a)
			if debug == 1 {
				fmt.Println("bcolor:", bcolor)
			}
		})
		twidget := container.NewHBox(
			tdecrease,
			tsz,
			tincrease,
			tcbutton,
			bgbutton)
		dcbutton := widget.NewButton("Date Color", func() {
			dcolor := colorPicker(settingsc, "date", a)
			if debug == 1 {
				fmt.Println("dcolor:", dcolor)
			}
		})
		dwidget := container.NewHBox(
			ddecrease,
			dsz,
			dincrease,
			dcbutton)
		ucbutton := widget.NewButton("UTC Time Color", func() {
			ucolor := colorPicker(settingsc, "utc", a)
			if debug == 1 {
				fmt.Println("ucolor:", ucolor)
			}
		})
		uwidget := container.NewHBox(
			udecrease,
			usz,
			uincrease,
			ucbutton)

		display := widget.NewForm(
			widget.NewFormItem("", mwidget),
			widget.NewFormItem("Time size", twidget),
			widget.NewFormItem("Date size", dwidget),
			widget.NewFormItem("UTC size", uwidget),
		)

		settingsc.Resize(fyne.NewSize(500, 300))
		// settings.CenterOnScreen() // run centered on primary (laptop) display
		settingsc.SetContent(container.NewVBox(setText, setform, display, buttonRow, doText))
		// reset.Resize(fyne.NewSize(reset.MinSize().Width, reset.MinSize().Height))
		settingsc.SetCloseIntercept(func() {
			if tselect != nil {
				tselect.Close()
				tselect = nil
			}
			settingsc.Close()
			settingsc = nil
		})
		settingsc.Show()
	}
}

func makeSettingsTheme(a fyne.App, w fyne.Window, bg fyne.Canvas) {
	// allow modifying the fyne theme
	// this is dependent on fyne_settings in ~/go/pkg/mod/fyne.io/fyne/v2/cmd/fyne_settings/settings
	// but here I use a customized version to add a button 'Apply & Close'
	// modify as shown below
	if settingsth != nil { // &&  !settingsc.Content().Visible() {
		settingsth.Show()
		teapot(a, settingsth)
	} else {
		s := settings.NewSettings()
		settingsth = a.NewWindow(appName + ": Theme Settings")
		settingsth.SetIcon(resourceKrankyBearBeretPng)

		appearance := s.LoadAppearanceScreen(w)
		tabs := container.NewAppTabs(
			&container.TabItem{Text: "Theme Appearance - affects all fyne based apps", Icon: s.AppearanceIcon(), Content: appearance})
		tabs.SetTabLocation(container.TabLocationLeading)
		settingsth.SetContent(tabs)

		settingsth.Resize(fyne.NewSize(520, 520))
		settingsth.CenterOnScreen() // run centered on primary (laptop) display
		settingsth.Show()
		settingsth.SetCloseIntercept(func() {
			settingsth.Close()
			settingsth = nil
		})
	}
}

// modify the latest ~/go/pkg/mod/fyne.io/fyne/v2/cmd/fyne_settings/settings/appearance.go

// add to function LoadAppearanceScreen last part with Apply & Close button:
/*
bottom := container.NewHBox(layout.NewSpacer(),
		&widget.Button{Text: "Apply", Importance: widget.HighImportance, OnTapped: func() {
			if s.fyneSettings.Scale == 0.0 {
				s.chooseScale(1.0)
			}
			err := s.save()
			if err != nil {
				fyne.LogError("Failed on saving", err)
			}

			s.appliedScale(s.fyneSettings.Scale)
		}},
		&widget.Button{Text: "Apply & Close", Importance: widget.WarningImportance, OnTapped: func() {
			if s.fyneSettings.Scale == 0.0 {
				s.chooseScale(1.0)
			}
			err := s.save()
			if err != nil {
				fyne.LogError("Failed on saving", err)
			}

			s.appliedScale(s.fyneSettings.Scale)
			w.Close()
		}},
	)
*/

func writeDefaultSettings(a fyne.App) {
	// write default prefs that can be modified via settings
	a.Preferences().SetInt("showseconds.default", 0)
	a.Preferences().SetInt("showtimezone.default", 1)
	a.Preferences().SetInt("showutc.default", 1)
	a.Preferences().SetInt("showdate.default", 1)
	a.Preferences().SetInt("showhr12.default", 1)
	a.Preferences().SetInt("hourchime.default", 1)
	a.Preferences().SetInt("slockmute.default", 0)
	a.Preferences().SetInt("automute.default", 0)
	a.Preferences().SetInt("muteonhr.default", 20)
	a.Preferences().SetInt("muteonmin.default", 0)
	a.Preferences().SetInt("muteoffhr.default", 8)
	a.Preferences().SetInt("muteoffmin.default", 0)
	a.Preferences().SetString("bgcolor.default", "0,143,251,255")
	a.Preferences().SetString("timecolor.default", "255,123,31,255")
	a.Preferences().SetString("datecolor.default", "131,222,74,255")
	a.Preferences().SetString("utccolor.default", "238,229,58,255")
	a.Preferences().SetString("timefont.default", "arial")
	a.Preferences().SetString("datefont.default", "arial")
	a.Preferences().SetString("utcfont.default", "arial")
	a.Preferences().SetInt("timesize.default", 48)
	a.Preferences().SetInt("datesize.default", 24)
	a.Preferences().SetInt("utcsize.default", 18)
	a.Preferences().SetString("hourchimesound.default", "cuckoo.mp3")
	a.Preferences().SetInt("startclock.default", startclock)
}

func writeSettings(a fyne.App) {
	// write current settings to global prefs
	a.Preferences().SetInt("showseconds.default", showseconds)
	a.Preferences().SetInt("showtimezone.default", showtimezone)
	a.Preferences().SetInt("showutc.default", showutc)
	a.Preferences().SetInt("showdate.default", showdate)
	a.Preferences().SetInt("showhr12.default", showhr12)
	a.Preferences().SetInt("hourchime.default", hourchime)
	a.Preferences().SetInt("slockmute.default", slockmute)
	a.Preferences().SetInt("automute.default", automute)
	a.Preferences().SetInt("muteonhr.default", muteonhr)
	a.Preferences().SetInt("muteonmin.default", muteonmin)
	a.Preferences().SetInt("muteoffhr.default", muteoffhr)
	a.Preferences().SetInt("muteoffmin.default", muteoffmin)
	a.Preferences().SetString("bgcolor.default", bgcolor)
	a.Preferences().SetString("timecolor.default", timecolor)
	a.Preferences().SetString("datecolor.default", datecolor)
	a.Preferences().SetString("utccolor.default", utccolor)
	a.Preferences().SetString("timefont.default", timefont)
	a.Preferences().SetString("datefont.default", datefont)
	a.Preferences().SetString("utcfont.default", utcfont)
	a.Preferences().SetInt("timesize.default", timesize)
	a.Preferences().SetInt("datesize.default", datesize)
	a.Preferences().SetInt("utcsize.default", utcsize)
	a.Preferences().SetString("hourchimesound.default", hourchimesound)
	a.Preferences().SetInt("startclock.default", startclock)
}

// func colorPicker(parent fyne.Window, colorDisplay *canvas.Rectangle) color.Color {
func colorPicker(parent fyne.Window, s string, a fyne.App) color.Color {
	// dialog.ShowCustom("Pick a Color", "Close", colorPicker, parent)
	picker := dialog.NewColorPicker("Select a color", "Choose your favorite color", func(c color.Color) {
		colorSelected(c, parent, s, a)
		mycolor = c
	}, parent)
	picker.Advanced = true
	picker.Show()
	return mycolor
}

func colorSelected(c color.Color, w fyne.Window, s string, a fyne.App) {
	rectangle := canvas.NewRectangle(c)
	size := 2 * theme.IconInlineSize()
	rectangle.SetMinSize(fyne.NewSize(size, size*1.8))
	mycolor := ColorToString(c)
	cmsg := "Color selected: " + mycolor
	dialog.ShowCustom(cmsg, "Ok", rectangle, w)
	switch s {
	case "time":
		a.Preferences().SetString("timecolor.default", mycolor)
	case "background":
		a.Preferences().SetString("bgcolor.default", mycolor)
	case "date":
		a.Preferences().SetString("datecolor.default", mycolor)
	case "utc":
		a.Preferences().SetString("utccolor.default", mycolor)
	}
}

// ColorToString converts a color.Color to a string in "rgba(r,g,b,a)" format.
func ColorToString(c color.Color) string {
	r, g, b, a := c.RGBA()
	// RGBA() method returns 16 bit values, need to divide by 257 to get 8 bit values
	// return fmt.Sprintf("rgba(%d,%d,%d,%.2f)", r/257, g/257, b/257, float64(a)/65535)
	// return fmt.Sprintf("rgba(%d,%d,%d,%d)", r/257, g/257, b/257, a/257)
	return fmt.Sprintf("%d,%d,%d,%d", r/257, g/257, b/257, a/257)
}

func showFilePicker(w fyne.Window) {
	// Show file picker and return selected file
	// https://dev.to/cjr29/learning-go-building-a-file-picker-using-fyneio-33le
	dialog.ShowFileOpen(func(f fyne.URIReadCloser, err error) {
		saveFile := "NoFileYet"
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		if f == nil {
			return
		}
		saveFile = f.URI().Path()
		fileURI = f.URI()
		selectedFile.SetText(saveFile)
	}, w)
}

func selectTime(a fyne.App, w fyne.Window, bg fyne.Canvas, caller string, hr int, min int) string {
	// var selectedTime time.Time
	var current string
	var myTime string

	switch caller {
	case "muteon":
		hr = muteonhr
		min = muteonmin
	case "muteoff":
		hr = muteoffhr
		min = muteoffmin
	default:
		hr = time.Now().Hour()
		min = time.Now().Minute()
	}

	tselect = a.NewWindow("Select Time")
	// Set window size to fit the input prompt
	tselect.Resize(fyne.NewSize(250, 100))

	current = fmt.Sprintf("%02d:%02d", hr, min)

	// Create a time entry widget
	timeEntry := widget.NewEntry()
	//timeEntry.SetPlaceHolder("Enter time (HH:MM)" + current)
	timeEntry.SetPlaceHolder(current)
	timeEntry.SetText(current)

	// Create a label to display messages
	messageLabel := widget.NewLabel("")

	// Create a button to submit the time
	submitButton := widget.NewButton("Set", func() {
		selectedTime := timeEntry.Text
		if isValidTime(selectedTime) {
			endTime, _ := time.Parse("15:04", selectedTime)
			messageLabel.SetText("Entered time: " + endTime.Format("15:04"))
			tselect.Close()
			tselect = nil
			parts := strings.Split(selectedTime, ":")
			hour, _ := strconv.Atoi(parts[0])
			min, _ := strconv.Atoi(parts[1])

			switch caller {
			case "muteon":
				muteonhr = hour
				muteonmin = min
				muteonbutton.SetText(fmt.Sprintf("Mute: %02d:%02d", muteonhr, muteonmin))
				muteonbutton.Refresh()
				a.Preferences().SetInt("muteonhr.default", muteonhr)
				a.Preferences().SetInt("muteonmin.default", muteonmin)
			case "muteoff":
				muteoffhr = hour
				muteoffmin = min
				muteoffbutton.SetText(fmt.Sprintf("Mute: %02d:%02d", muteoffhr, muteoffmin))
				muteoffbutton.Refresh()
				a.Preferences().SetInt("muteoffhr.default", muteoffhr)
				a.Preferences().SetInt("muteoffmin.default", muteoffmin)
			default:
				hour = time.Now().Hour()
				min = time.Now().Minute()
			}
		} else {
			messageLabel.SetText("Enter a valid time 00:00 to 23:59 (HH:MM)")
		}
	})

	// Arrange the widgets in a vertical box
	content := container.NewVBox(
		timeEntry,
		submitButton,
		messageLabel,
	)

	tselect.SetContent(content)
	// tselect.CenterOnScreen() // run centered on primary (laptop) display
	tselect.Show()
	return myTime
}

// isValidTime checks if the entered time is valid in 24-hour format hh:mm
func isValidTime(t string) bool {
	parts := strings.Split(t, ":")
	if len(parts) != 2 {
		return false
	}

	hours, err1 := strconv.Atoi(parts[0])
	minutes, err2 := strconv.Atoi(parts[1])

	if err1 != nil || err2 != nil {
		return false
	}

	if hours < 0 || hours > 23 || minutes < 0 || minutes > 59 {
		return false
	}
	return true
}

// "Now this is not the end. It is not even the beginning of the end. But it is, perhaps, the end of the beginning." Winston Churchill, November 10, 1942

// sample settings
// {"bgcolor.default":"0,143,251,255","color_recents":"#eee53a,#83de4a,#f44336,#ffffff,#9c27b0,#8bc34a,#ff9800","datecolor.default":"131,222,74,255","datefont.default":"arial","datesize.default":24,"hourchime.default":1,"hourchimesound.default":"cuckoo.mp3","showdate.default":1,"showhr12.default":1,"showseconds.default":1,"showtimezone.default":1,"showutc.default":1,"timecolor.default":"255,123,31,255","timefont.default":"arial","timesize.default":48,"utccolor.default":"238,229,58,255","utcfont.default":"arial","utcsize.default":18}
