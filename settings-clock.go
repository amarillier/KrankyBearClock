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
	Clock display colors and sizes now update automatically.
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
		jiggler := widget.NewRadioGroup([]string{"0", "5", "10", "15", "20", "30"}, func(value string) {
			if debug == 1 {
				log.Println("jiggler set to", value)
			}
			jiggle, _ = strconv.Atoi(value)
			a.Preferences().SetInt("jiggle.default", jiggle)
			lastJiggleMinute = -1
		})
		jiggler.Horizontal = true
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
			// Invalidate file cache when sound changes
			hourChimeFileChecked = false
			hourChimeCachedFile = ""
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
				updateClockDisplayColors()
			}
		}
		// Create buttons for increase and decrease
		tincrease := widget.NewButton("▲", func() {
			value, _ := strconv.Atoi(tsz.Text)
			if value < 200 {
				tsz.SetText(fmt.Sprintf("%d", value+1))
				timesize = value + 1
				a.Preferences().SetInt("timesize.default", timesize)
				updateClockDisplayColors()
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
				updateClockDisplayColors()
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
				updateClockDisplayColors()
			}
		}
		// Create buttons for increase and decrease
		dincrease := widget.NewButton("▲", func() {
			value, _ := strconv.Atoi(dsz.Text)
			if value < 200 {
				dsz.SetText(fmt.Sprintf("%d", value+1))
				datesize = value + 1
				a.Preferences().SetInt("datesize.default", datesize)
				updateClockDisplayColors()
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
				updateClockDisplayColors()
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
				updateClockDisplayColors()
			}
		}
		// Create buttons for increase and decrease
		uincrease := widget.NewButton("▲", func() {
			value, _ := strconv.Atoi(usz.Text)
			if value < 200 {
				usz.SetText(fmt.Sprintf("%d", value+1))
				utcsize = value + 1
				a.Preferences().SetInt("utcsize.default", utcsize)
				updateClockDisplayColors()
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
				updateClockDisplayColors()
			} else {
				playBeep("ding")
			}
		})

		/*
			timefont
			datefont
			utcfont
		*/

		// Reset button will be created after timezone widgets are defined
		var reset *widget.Button
		close := widget.NewButton("Close settings", func() {
			settingsc.Close()
			settingsc = nil
		})
		close.Importance = widget.WarningImportance // orange

		// buttonRow will be created after reset button is initialized
		var buttonRow *fyne.Container

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
		switch jiggle {
		case 0:
			jiggler.SetSelected("0")
		case 5:
			jiggler.SetSelected("5")
		case 10:
			jiggler.SetSelected("10")
		case 15:
			jiggler.SetSelected("15")
		case 20:
			jiggler.SetSelected("20")
		case 30:
			jiggler.SetSelected("30")
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
		// Common timezones list for selection
		commonTimezones := []string{
			"", "America/New_York", "America/Chicago", "America/Denver", "America/Los_Angeles",
			"America/Phoenix", "America/Anchorage", "America/Honolulu", "Europe/London",
			"Europe/Paris", "Europe/Berlin", "Europe/Rome", "Europe/Madrid", "Asia/Tokyo",
			"Asia/Shanghai", "Asia/Hong_Kong", "Asia/Dubai", "Asia/Kolkata", "Australia/Sydney",
			"Australia/Melbourne", "Pacific/Auckland", "America/Toronto", "America/Vancouver",
			"America/Mexico_City", "America/Sao_Paulo", "Africa/Johannesburg",
		}

		// Store references to timezone widgets for reset functionality
		var tz1Check, tz2Check, tz3Check, tz4Check, tz5Check *widget.Check
		var tz1Select, tz2Select, tz3Select, tz4Select, tz5Select *widget.Select
		var tz1OffsetEntry, tz2OffsetEntry, tz3OffsetEntry, tz4OffsetEntry, tz5OffsetEntry *widget.Entry

		// Timezone selector and enable checkbox functions
		createTimezoneRow := func(num int, enabledVar *int, nameVar *string, offsetVar *string, enabledCheck **widget.Check, timezoneSelect **widget.Select, offsetEntry **widget.Entry) *widget.FormItem {
			*enabledCheck = widget.NewCheck("", func(value bool) {
				if value {
					*enabledVar = 1
				} else {
					*enabledVar = 0
				}
				a.Preferences().SetInt(fmt.Sprintf("timezone%d.enabled", num), *enabledVar)
				// Reinitialize clock display to show/hide timezone
				if clock.Content() != nil {
					content := InitializeClockDisplay(a, clock)
					clock.SetContent(content)
					clock.Resize(fyne.NewSize(content.MinSize().Width*1.2, content.MinSize().Height*1.1))
				}
			})

			*timezoneSelect = widget.NewSelect(commonTimezones, func(value string) {
				*nameVar = value
				a.Preferences().SetString(fmt.Sprintf("timezone%d.name", num), *nameVar)
				// Reinitialize clock display to update timezone
				if clock.Content() != nil {
					content := InitializeClockDisplay(a, clock)
					clock.SetContent(content)
					clock.Resize(fyne.NewSize(content.MinSize().Width*1.2, content.MinSize().Height*1.1))
				}
			})

			// UTC offset entry field
			entry := widget.NewEntry()
			entry.SetPlaceHolder("UTC offset (e.g., +5, -3.5)")
			entry.SetText(*offsetVar)
			entry.OnChanged = func(value string) {
				*offsetVar = strings.TrimSpace(value)
				a.Preferences().SetString(fmt.Sprintf("timezone%d.offset", num), *offsetVar)
				// Reinitialize clock display to update timezone
				if clock.Content() != nil {
					content := InitializeClockDisplay(a, clock)
					clock.SetContent(content)
					clock.Resize(fyne.NewSize(content.MinSize().Width*1.2, content.MinSize().Height*1.1))
				}
			}
			// Make the entry wider
			entryContainer := container.NewWithoutLayout(entry)
			entryContainer.Resize(fyne.NewSize(200, entry.MinSize().Height))
			entry.Resize(fyne.NewSize(200, entry.MinSize().Height))
			*offsetEntry = entry

			// Set initial values
			if *enabledVar == 1 {
				(*enabledCheck).SetChecked(true)
			}
			if *nameVar != "" {
				(*timezoneSelect).SetSelected(*nameVar)
			} else {
				(*timezoneSelect).SetSelected("")
			}

			// Layout: checkbox, timezone select, UTC offset entry all on same row
			offsetLabel := widget.NewLabel("UTC offset:")
			timezoneWidget := container.NewHBox(
				*enabledCheck,
				*timezoneSelect,
				offsetLabel,
				entryContainer,
			)
			return widget.NewFormItem(fmt.Sprintf("Timezone %d", num), timezoneWidget)
		}

		tz1Row := createTimezoneRow(1, &timezone1Enabled, &timezone1Name, &timezone1Offset, &tz1Check, &tz1Select, &tz1OffsetEntry)
		tz2Row := createTimezoneRow(2, &timezone2Enabled, &timezone2Name, &timezone2Offset, &tz2Check, &tz2Select, &tz2OffsetEntry)
		tz3Row := createTimezoneRow(3, &timezone3Enabled, &timezone3Name, &timezone3Offset, &tz3Check, &tz3Select, &tz3OffsetEntry)
		tz4Row := createTimezoneRow(4, &timezone4Enabled, &timezone4Name, &timezone4Offset, &tz4Check, &tz4Select, &tz4OffsetEntry)
		tz5Row := createTimezoneRow(5, &timezone5Enabled, &timezone5Name, &timezone5Offset, &tz5Check, &tz5Select, &tz5OffsetEntry)

		// Now create reset button after timezone widgets are defined
		reset = widget.NewButton("Reset defaults", func() {
			if debug == 1 {
				log.Println("preferences reset to defaults")
			}
			writeDefaultSettings(a)
			showsec.SetChecked(false)
			showtz.SetChecked(true)
			showdt.SetChecked(true)
			showut.SetChecked(true)
			showhr1224.SetSelected("12")
			jiggler.SetSelected("0")
			lastJiggleMinute = -1
			lockmute.SetChecked(false)
			mute.SetChecked(false)
			muteonhr = 18
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

			// Reset additional timezones
			timezone1Enabled = 0
			timezone1Name = ""
			timezone1Offset = ""
			timezone2Enabled = 0
			timezone2Name = ""
			timezone2Offset = ""
			timezone3Enabled = 0
			timezone3Name = ""
			timezone3Offset = ""
			timezone4Enabled = 0
			timezone4Name = ""
			timezone4Offset = ""
			timezone5Enabled = 0
			timezone5Name = ""
			timezone5Offset = ""

			// Update timezone UI elements
			if tz1Check != nil {
				tz1Check.SetChecked(false)
				tz1Check.Refresh()
			}
			if tz1Select != nil {
				tz1Select.SetSelected("")
				tz1Select.Refresh()
			}
			if tz1OffsetEntry != nil {
				tz1OffsetEntry.SetText("")
				tz1OffsetEntry.Refresh()
			}
			if tz2Check != nil {
				tz2Check.SetChecked(false)
				tz2Check.Refresh()
			}
			if tz2Select != nil {
				tz2Select.SetSelected("")
				tz2Select.Refresh()
			}
			if tz2OffsetEntry != nil {
				tz2OffsetEntry.SetText("")
				tz2OffsetEntry.Refresh()
			}
			if tz3Check != nil {
				tz3Check.SetChecked(false)
				tz3Check.Refresh()
			}
			if tz3Select != nil {
				tz3Select.SetSelected("")
				tz3Select.Refresh()
			}
			if tz3OffsetEntry != nil {
				tz3OffsetEntry.SetText("")
				tz3OffsetEntry.Refresh()
			}
			if tz4Check != nil {
				tz4Check.SetChecked(false)
				tz4Check.Refresh()
			}
			if tz4Select != nil {
				tz4Select.SetSelected("")
				tz4Select.Refresh()
			}
			if tz4OffsetEntry != nil {
				tz4OffsetEntry.SetText("")
				tz4OffsetEntry.Refresh()
			}
			if tz5Check != nil {
				tz5Check.SetChecked(false)
				tz5Check.Refresh()
			}
			if tz5Select != nil {
				tz5Select.SetSelected("")
				tz5Select.Refresh()
			}
			if tz5OffsetEntry != nil {
				tz5OffsetEntry.SetText("")
				tz5OffsetEntry.Refresh()
			}

			showsec.Refresh()
			showtz.Refresh()
			showut.Refresh()
			showhr1224.Refresh()
			jiggler.Refresh()
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

			// Reload global color variables from preferences after reset
			bgcolor = a.Preferences().StringWithFallback("bgcolor.default", "0,143,251,255")
			timecolor = a.Preferences().StringWithFallback("timecolor.default", "255,123,31,255")
			datecolor = a.Preferences().StringWithFallback("datecolor.default", "131,222,74,255")
			utccolor = a.Preferences().StringWithFallback("utccolor.default", "238,229,58,255")

			// Reinitialize clock display to remove timezones
			if clock.Content() != nil {
				content := InitializeClockDisplay(a, clock)
				clock.SetContent(content)
				clock.Resize(fyne.NewSize(content.MinSize().Width*1.2, content.MinSize().Height*1.1))
			}

			// Update display immediately
			updateClockDisplayColors()
		})
		reset.Importance = widget.SuccessImportance // green

		// Now create buttonRow after reset button is initialized
		buttonRow = container.NewCenter(container.NewHBox(container.NewCenter(reset), container.NewCenter(close)))

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
		mwidget := container.NewHBox(muteonbutton, muteoffbutton)
		autoMuteSection := container.NewVBox(mute, mwidget)
		setform := widget.NewForm(
			widget.NewFormItem("Show Seconds", showsec),
			widget.NewFormItem("Show Timezone", showtz),
			widget.NewFormItem("Show Date", showdt),
			widget.NewFormItem("Show UTC", showut),
			widget.NewFormItem("Show 12/24 Hour Time", showhr1224),
			widget.NewFormItem("Mouse jiggler", jiggler),
			widget.NewFormItem("Auto Start at Boot", startatboot),
			widget.NewFormItem("Hourly Chime", chime),
			widget.NewFormItem("Hourly Chime Sound", chimesound),
			widget.NewFormItem("Lock Mute Volume", lockmute),
			widget.NewFormItem("Auto Mute Volume", autoMuteSection),
			widget.NewFormItem("Additional Timezones", widget.NewLabel("")),
			tz1Row,
			tz2Row,
			tz3Row,
			tz4Row,
			tz5Row,
		)
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
			widget.NewFormItem("Time size", twidget),
			widget.NewFormItem("Date size", dwidget),
			widget.NewFormItem("UTC size", uwidget),
		)

		// Resize to minimum required width to show all content (especially wider UTC offset fields)
		settingsc.Resize(fyne.NewSize(650, 300))
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
	a.Preferences().SetInt("jiggle.default", 0)
	a.Preferences().SetInt("hourchime.default", 1)
	a.Preferences().SetInt("slockmute.default", 0)
	a.Preferences().SetInt("automute.default", 0)
	a.Preferences().SetInt("muteonhr.default", 18)
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

	// Default timezone settings (all disabled)
	a.Preferences().SetInt("timezone1.enabled", 0)
	a.Preferences().SetString("timezone1.name", "")
	a.Preferences().SetString("timezone1.offset", "")
	a.Preferences().SetInt("timezone2.enabled", 0)
	a.Preferences().SetString("timezone2.name", "")
	a.Preferences().SetString("timezone2.offset", "")
	a.Preferences().SetInt("timezone3.enabled", 0)
	a.Preferences().SetString("timezone3.name", "")
	a.Preferences().SetString("timezone3.offset", "")
	a.Preferences().SetInt("timezone4.enabled", 0)
	a.Preferences().SetString("timezone4.name", "")
	a.Preferences().SetString("timezone4.offset", "")
	a.Preferences().SetInt("timezone5.enabled", 0)
	a.Preferences().SetString("timezone5.name", "")
	a.Preferences().SetString("timezone5.offset", "")

	// Reload global color variables from preferences after setting defaults
	bgcolor = a.Preferences().StringWithFallback("bgcolor.default", "0,143,251,255")
	timecolor = a.Preferences().StringWithFallback("timecolor.default", "255,123,31,255")
	datecolor = a.Preferences().StringWithFallback("datecolor.default", "131,222,74,255")
	utccolor = a.Preferences().StringWithFallback("utccolor.default", "238,229,58,255")

	// Update display immediately
	updateClockDisplayColors()
}

func writeSettings(a fyne.App) {
	// write current settings to global prefs
	a.Preferences().SetInt("showseconds.default", showseconds)
	a.Preferences().SetInt("showtimezone.default", showtimezone)
	a.Preferences().SetInt("showutc.default", showutc)
	a.Preferences().SetInt("showdate.default", showdate)
	a.Preferences().SetInt("showhr12.default", showhr12)
	a.Preferences().SetInt("jiggle.default", jiggle)
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

	// Save timezone settings
	a.Preferences().SetInt("timezone1.enabled", timezone1Enabled)
	a.Preferences().SetString("timezone1.name", timezone1Name)
	a.Preferences().SetString("timezone1.offset", timezone1Offset)
	a.Preferences().SetInt("timezone2.enabled", timezone2Enabled)
	a.Preferences().SetString("timezone2.name", timezone2Name)
	a.Preferences().SetString("timezone2.offset", timezone2Offset)
	a.Preferences().SetInt("timezone3.enabled", timezone3Enabled)
	a.Preferences().SetString("timezone3.name", timezone3Name)
	a.Preferences().SetString("timezone3.offset", timezone3Offset)
	a.Preferences().SetInt("timezone4.enabled", timezone4Enabled)
	a.Preferences().SetString("timezone4.name", timezone4Name)
	a.Preferences().SetString("timezone4.offset", timezone4Offset)
	a.Preferences().SetInt("timezone5.enabled", timezone5Enabled)
	a.Preferences().SetString("timezone5.name", timezone5Name)
	a.Preferences().SetString("timezone5.offset", timezone5Offset)
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
		timecolor = mycolor
	case "background":
		a.Preferences().SetString("bgcolor.default", mycolor)
		bgcolor = mycolor
	case "date":
		a.Preferences().SetString("datecolor.default", mycolor)
		datecolor = mycolor
	case "utc":
		a.Preferences().SetString("utccolor.default", mycolor)
		utccolor = mycolor
	}
	// Update clock display immediately with new colors
	updateClockDisplayColors()
}

// updateClockDisplayColors updates the clock display with current color and size settings
func updateClockDisplayColors() {
	if nowtime == nil || nowdate == nil || utctime == nil || background == nil {
		// Clock display not initialized yet
		return
	}

	// Parse time color
	var tre, tgr, tbl, ta uint8
	colors := strings.Split(timecolor, ",")
	if len(colors) == 4 {
		col, _ := strconv.ParseUint(colors[0], 10, 8)
		tre = uint8(col)
		col, _ = strconv.ParseUint(colors[1], 10, 8)
		tgr = uint8(col)
		col, _ = strconv.ParseUint(colors[2], 10, 8)
		tbl = uint8(col)
		col, _ = strconv.ParseUint(colors[3], 10, 8)
		ta = uint8(col)
		nowtime.Color = color.RGBA{R: tre, G: tgr, B: tbl, A: ta}
		nowtime.TextSize = float32(timesize)
	}

	// Parse date color
	var dre, dgr, dbl, da uint8
	colors = strings.Split(datecolor, ",")
	if len(colors) == 4 {
		col, _ := strconv.ParseUint(colors[0], 10, 8)
		dre = uint8(col)
		col, _ = strconv.ParseUint(colors[1], 10, 8)
		dgr = uint8(col)
		col, _ = strconv.ParseUint(colors[2], 10, 8)
		dbl = uint8(col)
		col, _ = strconv.ParseUint(colors[3], 10, 8)
		da = uint8(col)
		nowdate.Color = color.RGBA{R: dre, G: dgr, B: dbl, A: da}
		nowdate.TextSize = float32(datesize)
	}

	// Parse UTC color
	var ure, ugr, ubl, ua uint8
	colors = strings.Split(utccolor, ",")
	if len(colors) == 4 {
		col, _ := strconv.ParseUint(colors[0], 10, 8)
		ure = uint8(col)
		col, _ = strconv.ParseUint(colors[1], 10, 8)
		ugr = uint8(col)
		col, _ = strconv.ParseUint(colors[2], 10, 8)
		ubl = uint8(col)
		col, _ = strconv.ParseUint(colors[3], 10, 8)
		ua = uint8(col)
		utctime.Color = color.RGBA{R: ure, G: ugr, B: ubl, A: ua}
		utctime.TextSize = float32(utcsize)
	}

	// Parse background color
	var bre, bgr, bbl, ba uint8
	colors = strings.Split(bgcolor, ",")
	if len(colors) == 4 {
		col, _ := strconv.ParseUint(colors[0], 10, 8)
		bre = uint8(col)
		col, _ = strconv.ParseUint(colors[1], 10, 8)
		bgr = uint8(col)
		col, _ = strconv.ParseUint(colors[2], 10, 8)
		bbl = uint8(col)
		col, _ = strconv.ParseUint(colors[3], 10, 8)
		ba = uint8(col)
		background.FillColor = color.RGBA{R: bre, G: bgr, B: bbl, A: ba}
	}

	// Refresh all elements
	fyne.Do(func() {
		nowtime.Refresh()
		nowdate.Refresh()
		utctime.Refresh()
		background.Refresh()
	})

	// Also update countdown window if it's open
	updateCountdownColors()
	// Also update alarm window if it's open
	updateAlarmColors()
	// Also update weather window if it's open
	updateWeatherColors()
}

// updateCountdownColors updates the countdown window with current color settings
func updateCountdownColors() {
	if countdown == nil || countdownBackground == nil || countdownTitleText == nil || countdownHelpText == nil {
		// Countdown window not initialized yet
		return
	}

	// Parse time color
	var tre, tgr, tbl, ta uint8
	colors := strings.Split(timecolor, ",")
	if len(colors) == 4 {
		col, _ := strconv.ParseUint(colors[0], 10, 8)
		tre = uint8(col)
		col, _ = strconv.ParseUint(colors[1], 10, 8)
		tgr = uint8(col)
		col, _ = strconv.ParseUint(colors[2], 10, 8)
		tbl = uint8(col)
		col, _ = strconv.ParseUint(colors[3], 10, 8)
		ta = uint8(col)
		countdownTitleText.Color = color.RGBA{R: tre, G: tgr, B: tbl, A: ta}
		countdownTitleText.TextSize = float32(timesize)
	}

	// Parse date color
	var dre, dgr, dbl, da uint8
	colors = strings.Split(datecolor, ",")
	if len(colors) == 4 {
		col, _ := strconv.ParseUint(colors[0], 10, 8)
		dre = uint8(col)
		col, _ = strconv.ParseUint(colors[1], 10, 8)
		dgr = uint8(col)
		col, _ = strconv.ParseUint(colors[2], 10, 8)
		dbl = uint8(col)
		col, _ = strconv.ParseUint(colors[3], 10, 8)
		da = uint8(col)
		countdownHelpText.Color = color.RGBA{R: dre, G: dgr, B: dbl, A: da}
		countdownHelpText.TextSize = float32(datesize) * 0.8

		// Update days text colors (only if not in special state - updateDays handles special cases)
		// For normal future dates, use date color
		if countdownDaysText1 != nil && countdownDaysText1.Text != "" &&
			!strings.Contains(countdownDaysText1.Text, "ago") &&
			countdownDaysText1.Text != "Today!" && countdownDaysText1.Text != "Tomorrow!" {
			countdownDaysText1.Color = color.RGBA{R: dre, G: dgr, B: dbl, A: da}
		}
		if countdownDaysText2 != nil && countdownDaysText2.Text != "" &&
			!strings.Contains(countdownDaysText2.Text, "ago") &&
			countdownDaysText2.Text != "Today!" && countdownDaysText2.Text != "Tomorrow!" {
			countdownDaysText2.Color = color.RGBA{R: dre, G: dgr, B: dbl, A: da}
		}
		if countdownDaysText3 != nil && countdownDaysText3.Text != "" &&
			!strings.Contains(countdownDaysText3.Text, "ago") &&
			countdownDaysText3.Text != "Today!" && countdownDaysText3.Text != "Tomorrow!" {
			countdownDaysText3.Color = color.RGBA{R: dre, G: dgr, B: dbl, A: da}
		}
	}

	// Parse background color
	var bre, bgr, bbl, ba uint8
	colors = strings.Split(bgcolor, ",")
	if len(colors) == 4 {
		col, _ := strconv.ParseUint(colors[0], 10, 8)
		bre = uint8(col)
		col, _ = strconv.ParseUint(colors[1], 10, 8)
		bgr = uint8(col)
		col, _ = strconv.ParseUint(colors[2], 10, 8)
		bbl = uint8(col)
		col, _ = strconv.ParseUint(colors[3], 10, 8)
		ba = uint8(col)
		countdownBackground.FillColor = color.RGBA{R: bre, G: bgr, B: bbl, A: ba}
	}

	fyne.Do(func() {
		countdownTitleText.Refresh()
		countdownHelpText.Refresh()
		countdownBackground.Refresh()
		if countdownDaysText1 != nil {
			countdownDaysText1.Refresh()
		}
		if countdownDaysText2 != nil {
			countdownDaysText2.Refresh()
		}
		if countdownDaysText3 != nil {
			countdownDaysText3.Refresh()
		}
	})
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
	// If window already exists and is visible, just bring it to front
	if tselect != nil && tselect.Content().Visible() {
		tselect.Show()
		return ""
	}

	// If window exists but not visible, close it first
	if tselect != nil {
		tselect.Close()
		tselect = nil
	}

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

	tselect = a.NewWindow("Select " + caller + " Time")
	tselect.Resize(fyne.NewSize(250, 100))
	current = fmt.Sprintf("%02d:%02d", hr, min)
	timeEntry := widget.NewEntry()
	timeEntry.SetPlaceHolder(current)
	timeEntry.SetText(current)
	messageLabel := widget.NewLabel("")
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
