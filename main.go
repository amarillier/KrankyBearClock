package main

import (
	"fmt"
	"image/color"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/IamFaizanKhalid/lock"
	"github.com/itchyny/volume-go"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"fyne.io/systray"
	// audio "github.com/amarillier/KrankyBearModule/audio"
	// util "github.com/amarillier/KrankyBearModule/util"
)

const (
	// appName    = "Kranky Bear Clock"
	appVersion = "0.4.1" // see FyneApp.toml
	appAuthor  = "Allan Marillier"
)

var appName = "Kranky Bear Clock"
var appNameCustom = ""
var appCopyright = "Copyright (c) Allan Marillier, 2024-" + strconv.Itoa(time.Now().Year())
var imgDir string
var clockbg string // future optional clock background image

var sndDir string
var debug int = 0
var clock fyne.Window
var settingsc fyne.Window
var settingsth fyne.Window
var abt fyne.Window
var updt fyne.Window
var hlp fyne.Window

// var egg fyne.Window
var bg fyne.Canvas

var showseconds int
var showtimezone int
var showdate int
var showutc int
var showhr12 int
var hourchime int
var slockmute int
var clockmutedvol int
var automute int
var currentvolume int
var muteonhr int
var muteonmin int
var muteoffhr int
var muteoffmin int
var bgcolor string
var timecolor string
var datecolor string
var utccolor string
var timefont string
var datefont string
var utcfont string
var timesize int
var datesize int
var utcsize int
var hourchimesound string
var startclock int
var processName string

// preferences stored via fyne preferences API land in
// ~/Library/Preferences/fyne/com.github.amarillier.KrankyBearClock/preferences.json
// ~\AppData\Roaming\fyne\com.github.amarillier.KrankyBearClock\preferences.json
// {"bgcolor.default":"0,143,251,255","color_recents":"#eee53a,#83de4a,#f44336,#ffffff,#9c27b0,#8bc34a,#ff9800","datecolor.default":"131,222,74,255","datefont.default":"arial","datesize.default":24,"hourchime.default":1,"hourchimesound.default":"cuckoo.mp3","showdate.default":1,"showhr12.default":1,"showseconds.default":0,"showtimezone.default":1,"showutc.default":1,"startclock.default":0,"timecolor.default":"255,123,31,255","timefont.default":"arial","timesize.default":48,"utccolor.default":"238,229,58,255","utcfont.default":"arial","utcsize.default":18}

func main() {
	exePath, err := os.Executable()
	processName = filepath.Base(os.Args[0])
	if err != nil {
		panic(err)
	}

	launchDir := filepath.Dir(exePath)
	if runtime.GOOS == "darwin" {
		if strings.HasPrefix(launchDir, "/Applications/KrankyBearClock") {
			sndDir = launchDir + "/../Resources/Sounds"
			imgDir = launchDir + "/../Resources/Images"
		} else {
			sndDir = launchDir + "/Resources/Sounds"
			imgDir = launchDir + "/Resources/Images"
		}
	} else if runtime.GOOS == "windows" {
		sndDir = launchDir + "/Resources/Sounds"
		imgDir = launchDir + "/Resources/Images"
	}

	a := app.NewWithID("com.github.amarillier.KrankyBearClock")
	a.Settings().SetTheme(&appTheme{Theme: theme.DefaultTheme()})
	clock = a.NewWindow(appName)
	clock.SetIcon(resourceKrankyBearClockPng)
	clock.SetPadded(false)
	//clock.SetCloseIntercept(func() {
	//	a.Quit() // force quit, normal when somebody hits "x" to close
	//})
	clock.SetMaster() // this sets this as master and closes all child windows
	// clock.CenterOnScreen() // run centered on primary (laptop) display

	prefs := strings.ReplaceAll((a.Storage().RootURI()).String(), "file://", "") + "/preferences.json"
	if !checkFileExists(prefs) {
		if debug == 1 {
			log.Println("prefs file does not exist")
		}
		// add some default prefs that can be modified via settings
		writeDefaultSettings(a)
	}
	// get default settings from preferences
	showseconds = a.Preferences().IntWithFallback("showseconds.default", 1)
	showtimezone = a.Preferences().IntWithFallback("showtimezone.default", 1)
	showdate = a.Preferences().IntWithFallback("showdate.default", 1)
	showutc = a.Preferences().IntWithFallback("showutc.default", 1)
	showhr12 = a.Preferences().IntWithFallback("showhr12.default", 1)
	slockmute = a.Preferences().IntWithFallback("slockmute.default", 0)
	automute = a.Preferences().IntWithFallback("automute.default", 0)
	muteonhr = a.Preferences().IntWithFallback("muteonhr.default", 20)
	muteonmin = a.Preferences().IntWithFallback("muteonmin.default", 0)
	muteoffhr = a.Preferences().IntWithFallback("muteoffhr.default", 8)
	muteoffmin = a.Preferences().IntWithFallback("muteoffmin.default", 0)
	hourchime = a.Preferences().IntWithFallback("hourchime.default", 1)
	bgcolor = a.Preferences().StringWithFallback("bgcolor.default", "0,143,251,255")      // blue
	timecolor = a.Preferences().StringWithFallback("timecolor.default", "255,123,31,255") // orange
	datecolor = a.Preferences().StringWithFallback("datecolor.default", "131,222,74,255") // red
	utccolor = a.Preferences().StringWithFallback("utccolor.default", "238,229,58.255")   // yellow
	timefont = a.Preferences().StringWithFallback("timefont.default", "arial")            // not yet!
	datefont = a.Preferences().StringWithFallback("datefont.default", "arial")            // not yet!
	utcfont = a.Preferences().StringWithFallback("utcfont.default", "arial")              // not yet!
	timesize = a.Preferences().IntWithFallback("timesize.default", 36)
	datesize = a.Preferences().IntWithFallback("datesize.default", 24)
	utcsize = a.Preferences().IntWithFallback("utcsize.default", 18)
	hourchimesound = a.Preferences().StringWithFallback("hourchimesound.default", "hero.mp3")
	startclock = a.Preferences().IntWithFallback("startclock.default", 0)
	writeSettings(a)

	clockmutedvol = 0
	var tre, tgr, tbl, ta uint8
	colors := strings.Split(timecolor, ",")
	col, _ := strconv.ParseUint(colors[0], 10, 8)
	tre = uint8(col)
	col, _ = strconv.ParseUint(colors[1], 10, 8)
	tgr = uint8(col)
	col, _ = strconv.ParseUint(colors[2], 10, 8)
	tbl = uint8(col)
	col, _ = strconv.ParseUint(colors[3], 10, 8)
	ta = uint8(col)

	var bre, bgr, bbl, ba uint8
	colors = strings.Split(bgcolor, ",")
	col, _ = strconv.ParseUint(colors[0], 10, 8)
	bre = uint8(col)
	col, _ = strconv.ParseUint(colors[1], 10, 8)
	bgr = uint8(col)
	col, _ = strconv.ParseUint(colors[2], 10, 8)
	bbl = uint8(col)
	col, _ = strconv.ParseUint(colors[3], 10, 8)
	ba = uint8(col)

	var dre, dgr, dbl, da uint8
	colors = strings.Split(datecolor, ",")
	col, _ = strconv.ParseUint(colors[0], 10, 8)
	dre = uint8(col)
	col, _ = strconv.ParseUint(colors[1], 10, 8)
	dgr = uint8(col)
	col, _ = strconv.ParseUint(colors[2], 10, 8)
	dbl = uint8(col)
	col, _ = strconv.ParseUint(colors[3], 10, 8)
	da = uint8(col)

	var ure, ugr, ubl, ua uint8
	colors = strings.Split(utccolor, ",")
	col, _ = strconv.ParseUint(colors[0], 10, 8)
	ure = uint8(col)
	col, _ = strconv.ParseUint(colors[1], 10, 8)
	ugr = uint8(col)
	col, _ = strconv.ParseUint(colors[2], 10, 8)
	ubl = uint8(col)
	col, _ = strconv.ParseUint(colors[3], 10, 8)
	ua = uint8(col)

	if len(os.Args) >= 2 {
		log.Println("arg count:", len(os.Args))
		if os.Args[1] == "debug" || os.Args[1] == "d" {
			debug = 1
			logInit()
			r, _ := os.Open("KrankyBearClock0.txt")
			logLines, _ := lineCounter(r)
			r.Close()
			InfoLog.Println("logLines:", logLines)
			if logLines >= 100 {
				logRotate()
			}
			logInit()
			InfoLog.Println("Opening the application...")
			InfoLog.Println("Something has occurred...")
			WarningLog.Println("WARNING!!!..")
			ErrorLog.Println("Some error has occurred...")

			log.Println("debug mode:", debug)
			log.Println("exepath:", exePath)
			log.Println("launchdir:", launchDir)
			log.Println("Images:", imgDir)
			log.Println("Sounds:", sndDir)
			log.Println("showseconds:", showseconds)
			log.Println("showtimezone:", showtimezone)
			log.Println("showutc:", showutc)
			log.Println("showhr12:", showhr12)
			log.Println("hourchime:", hourchime)
			log.Println("slockmute:", slockmute)
			log.Println("bgcolor:", bgcolor)
			log.Println("timecolor:", timecolor)
			log.Println("datecolor:", datecolor)
			log.Println("utccolor:", utccolor)
			log.Println("timefont:", timefont)
			log.Println("datefont:", datefont)
			log.Println("utcfont:", utcfont)
			log.Println("timesize:", timesize)
			log.Println("datesize:", datesize)
			log.Println("utcsize:", utcsize)
			log.Println("hourchimesound:", hourchimesound)
			log.Println("startclock:", startclock)
		}
	}

	// check update first
	updtmsg, updateAvail := updateChecker("amarillier", "KrankyBearClock", "Kranky Bear Clock", "https://github.com/amarillier/KrankyBearClock/releases/latest")
	if updateAvail {
		// open a window to show the update message
		// no need to test for updt window open at first start
		kb := canvas.NewImageFromResource(resourceKrankyBearPng)
		kb.FillMode = canvas.ImageFillOriginal
		text := widget.NewLabel(updtmsg)
		content := container.NewHBox(kb, text)
		updt = a.NewWindow(appName + ": Update Check")
		updt.SetIcon(resourceKrankyBearPng)
		updt.Resize(fyne.NewSize(50, 100))
		updt.SetContent(content)
		updt.SetCloseIntercept(func() {
			updt.Close()
			updt = nil
		})
		updt.CenterOnScreen() // run centered on primary (laptop) display
		updt.Show()
	}

	if desk, ok := a.(desktop.App); ok {
		show := fyne.NewMenuItem("Show", func() {
			clock.Show()
			clock.Canvas().Focused()
		})
		hide := fyne.NewMenuItem("Hide", clock.Hide)
		about := fyne.NewMenuItem("About", func() {
			aboutText := appName + " v " + appVersion
			aboutText += "\n" + appCopyright
			aboutText += "\n\nCreated by " + appAuthor + ", using Go and fyne GUI"
			aboutText += "\n\nNo obligation, it's rewarding to hear if use this app."
			aboutText += "\n\nAnd looking about about and help or settings too much might expose an easter egg!"

			kb := canvas.NewImageFromResource(resourceKrankyBearPng)
			text := widget.NewLabel(aboutText)
			kb.FillMode = canvas.ImageFillOriginal
			content := container.NewHBox(kb, text)

			if abt == nil || !abt.Content().Visible() {
				abt = a.NewWindow(appName + ": About")
				abt.SetIcon(resourceKrankyBearClockPng)
				abt.Resize(fyne.NewSize(50, 100))
				// abt.SetContent(widget.NewLabel(aboutText))
				abt.SetContent(content)
				abt.SetCloseIntercept(func() {
					abt.Close()
					abt = nil
				})
				// abt.CenterOnScreen() // run centered on primary (laptop) display
				abt.Show()
			} else {
				abt.Show()
				easterEgg(a, abt)
			}
		})
		help := fyne.NewMenuItem("Help", func() {
			// if hlp != nil { // &&  !hlp.Content().Visible() {
			if hlp == nil || !hlp.Content().Visible() {
				hlp = a.NewWindow(appName + ": Help")
				hlp.SetIcon(resourceKrankyBearClockPng)

				hlp.SetCloseIntercept(func() {
					hlp.Close()
					hlp = nil
				})
				//}
				hlpText := `This is a basic desktop clock that currently shows:

- time in 12 /24 hour format 
- optional seconds
- optional timezone
- optional date in full day name, month date #, 4 digit year
- optional UTC time and time zone offset in hours
- optional hourly chime with user selectable sound
	- ideally chose a very short chime sound, less than half second
- customizable font sizes for each of time, date and UTC time
- customizable font color for each of background, time, date and UTC time
- clock display window resizes automatically to suit selected font sizes
- optional setting to enable auto starting at boot

- Note: Displaying seconds can be quite resource intensive with clock display updates every second. 
  The app can be substantially less CPU intensive when seconds are not displayed, allowing the app to
  refresh the display every minute rather than every second

- See Settings Info tab for more detail on settings / preferences

- Default settings will be created on first run if they don't exist
`
				hlpText += "\n" + appName + " v " + appVersion
				hlpText += "\n" + appCopyright
				hlpText += "\n\n" + appAuthor + ", using Go and fyne GUI"

				plnText := `- Allow multiple time zones for clock, hh:mm only + offset
- Allow multiple alarm times with user selectable tones for each
- Allow settings set/save window locations to open clock, 
	unfortunately not implemented in the fyne library yet
- Open with clock window focused
	- this is currently MacOS LaunchPad behavior, but only allows one app
	- To run more than one simultaneously, in terminal: open -n -a KrankyBearClock 
- Add settings to allow:
	- clock text font
	- day text font
	- UTC time font
	- possible alarm(s), if so optional sounds to play
	- possible tie to Outlook calendar for alarms? Probably not
	- possible choice of svg, png, jpg background image
`

				bugText := `
- Activating tray menus causes running clock display to not show updates
	until Help, About, Settings etc are selected
	- But clock does continue to run, fix to run systray, settings etc in parallel
- Font type settings in preferences are currently ignored, the app uses system theme defaults. (Future planned update)
- Settings changes to background and clock default times are saved immediately.
	- but clock time format size and color, date size and color and background do
	not currently refresh to new settings - exit and rerun for now
`
				link, err := url.Parse("https://github.com/amarillier/KrankyBearClock/blob/main/license.txt")
				if err != nil {
					fyne.LogError("Could not parse URL", err)
				}
				hyperlink := widget.NewHyperlink("https://github.com/amarillier/KrankyBearClock/blob/main/license.txt", link)
				hyperlink.Alignment = fyne.TextAlignLeading

				licText := `Kranky Bear Clock is FREE Software” as defined in the license agreement below. 
 
This application is "FREE Software". 

This application is intended for any use by any individual, in any organization.

This application provides no guarantees as to stability of operations or suitability 
for any purpose, but every attempt has been made to make this application reliable.

This application may not be sold, no money may be asked by anyone for provision of, or any services related to this application.

Using this application (and reading this text) is considered acceptance of
the terms of the License Agreement, and acknowledgement that this is FREE
Software and the additional terms above.

See https://github.com/amarillier/KrankyBearClock/
`

				settingsText := `Settings are a separate tray menu item
Settings contains defaults as below, which can be modified, and also reset to defaults:
{"bgcolor.default":"0,143,251,255",
"color_recents":"#eee53a,#83de4a,#f44336,#ffffff,#9c27b0,#8bc34a,#ff9800",
"datecolor.default":"131,222,74,255","datefont.default":"arial",
"datesize.default":24,"hourchime.default":1,
"hourchimesound.default":"cuckoo.mp3","showdate.default":1,
"showhr12.default":1,"showseconds.default":0,"showtimezone.default":1,
"showutc.default":1,"startclock.default":0,"timecolor.default":"255,123,31,255",
"timefont.default":"arial","timesize.default":48,"utccolor.default":"238,229,58,255",
"utcfont.default":"arial","utcsize.default":18}

KrankyBearClock looks for directories named Resources/Images and Resources/Sounds,
containing images and sounds.

IMAGES:
Future additions will allow selecting background images of your choice, png, SVG,
	jpg maybe and specifying size - height / width. Manual window resizing
	is already possible

SOUNDS:
Built in tones include 'ding', 'down', 'up', and 'updown'. These are always available
	and will be listed first in sound selectors
The sounds directory as distributed also contains a number of other .mp3 files
including baseball.mp3, grandfatherclock.mp3, hero.mp3, pinball.mp3, sosumi.mp3
When selecting sounds, the sound will be played as a preview when possible.
When selected sounds are not present (removed from Sounds), KrankyBearClock defaults
	to playing built in tones ding, down, up or updown
Future additions will allow also choosing from any .mid or .wav sound files of your
	choice if located in the Sounds directory
`
				lic := widget.NewLabel(licText)
				tabs := container.NewDocTabs(
					container.NewTabItem("Help", widget.NewLabel(hlpText)),
					container.NewTabItem("Known Issues", widget.NewLabel(bugText)),
					container.NewTabItem("Planned Updates", widget.NewLabel(plnText)),
					container.NewTabItem("Settings Info", widget.NewLabel(settingsText)),
					container.NewTabItem("License", container.NewVBox(lic, hyperlink)),
				)
				tabs.SetTabLocation(container.TabLocationTop)
				tabs.Show()
				hlp.Resize(fyne.NewSize(800, 300))
				hlp.SetContent(tabs)
				// hlp.CenterOnScreen() // run centered on primary (laptop) display
				hlp.Show()
			} else {
				hlp.Show()
				easterEgg(a, hlp)
			}
		})
		settingsClock := fyne.NewMenuItem("Settings (Clock)", func() {
			makeSettingsClock(a, clock, bg)
		})
		settingsTheme := fyne.NewMenuItem("Settings (Theme)", func() {
			makeSettingsTheme(a, clock, bg)
		})
		updtchk := fyne.NewMenuItem("Check for update", func() {
			// throw away updateAvail here, use _, unneeded for manual check
			updtmsg, _ := updateChecker("amarillier", "KrankyBearClock", "Kranky Bear Clock", "https://github.com/amarillier/KrankyBearClock/releases/latest")
			if updt == nil {
				kb := canvas.NewImageFromResource(resourceKrankyBearPng)
				kb.FillMode = canvas.ImageFillOriginal
				text := widget.NewLabel(updtmsg)
				content := container.NewHBox(kb, text)
				updt = a.NewWindow(appName + ": Update Check")
				updt.SetIcon(resourceKrankyBearPng)
				updt.Resize(fyne.NewSize(50, 100))
				// updt.SetContent(widget.NewLabel(updtmsg))
				updt.SetContent(content)
				updt.SetCloseIntercept(func() {
					updt.Close()
					updt = nil
				})
				updt.CenterOnScreen() // run centered on pr1imary (laptop) display
				updt.Show()
				// if !strings.Contains(updtmsg, "You are running the latest") {
				if updateAvail {
					if !checkFileExists(sndDir + "/KrankyBearGrowl.mp3") {
						playBeep("up")
					} else {
						playMp3(sndDir + "//KrankyBearGrowl.mp3") // Basso, Blow, Hero, Funk, Glass, Ping, Purr, Sosumi, Submarine,
					}
				}
			} else {
				updt.RequestFocus()
			}
		})
		menu := fyne.NewMenu(a.Metadata().Name, show, hide, fyne.NewMenuItemSeparator(), about, updtchk, help, settingsClock, settingsTheme)
		desk.SetSystemTrayMenu(menu)
		desk.SetSystemTrayIcon(resourceKrankyBearClockPng)
		systray.SetTooltip(appName)
		// systray.SetTitle(clockName)

		// Menu items
		// compile / run with syntax below to force Mac to do menus like Windows
		// otherwise menus will be at the top of the display
		// https://github.com/fyne-io/fyne/issues/3988
		// go build -tags no_native_menus .
		// go run -tags no_native_menus .
		quit := fyne.NewMenuItem("Quit", func() {
			a.Quit()
		})
		newMenuOps := fyne.NewMenu("Operations", show, hide, fyne.NewMenuItemSeparator(), quit)
		newMenuHelp := fyne.NewMenu("Help", about, help)
		newMenuSettings := fyne.NewMenu("Settings", settingsClock, settingsTheme)
		// New main menu
		cmenu := fyne.NewMainMenu(newMenuOps, newMenuHelp, newMenuSettings)
		// setup main menu
		clock.SetMainMenu(cmenu)
		// cmenu.Refresh()
	}

	now := time.Now()
	// timeFormat := `15:04:05`
	// timeFormat := `3:04:05 PM (MST)`
	timeFormat := ``
	if showhr12 == 1 {
		timeFormat += `3:04`
	} else {
		timeFormat += `15:04`
	}
	if showseconds == 1 {
		timeFormat += `:05`
	}
	if showhr12 == 1 {
		timeFormat += ` PM` // this needs to be added AFTER seconds if 12 hour
	}
	if showtimezone == 1 {
		timeFormat += ` (MST)`
	}

	// Get the local time zone and offset
	_, offset := now.Zone()
	offsetHours := offset / 3600
	offsetMinutes := (offset % 3600) / 60
	offsetString := fmt.Sprintf(" (local is  %+02d:%02d)", offsetHours, offsetMinutes)
	// utcFormat := `(UTC 3:04 PM Z07)`
	utcFormat := `UTC 3:04 PM`
	dateFormat := ` Monday, January 2, 2006 `

	// nowtime := canvas.NewText(now.Format(timeFormat), color.RGBA{R: 255, G: 123, B: 31, A: 255})
	nowtime := canvas.NewText(now.Format(timeFormat), color.RGBA{R: tre, G: tgr, B: tbl, A: ta})
	nowtime.TextStyle = fyne.TextStyle{Bold: true}
	// nowtime.TextStyle = fyne.TextStyle{Monospace: true} // EXAMPLE FONT TYPE
	nowtime.Alignment = fyne.TextAlignCenter
	nowtime.TextSize = float32(timesize)

	// utctime := canvas.NewText(now.Format(utcFormat), color.RGBA{R: 255, G: 123, B: 31, A: 255})
	utctime := canvas.NewText(now.Format(utcFormat), color.RGBA{R: ure, G: ugr, B: ubl, A: ua})
	utctime.TextStyle = fyne.TextStyle{Bold: true}
	utctime.Alignment = fyne.TextAlignCenter
	utctime.TextSize = float32(utcsize)

	// nowdate := canvas.NewText(now.Format(dateFormat), color.RGBA{R: 208, G: 145, B: 38, A: 255})
	nowdate := canvas.NewText(now.Format(dateFormat), color.RGBA{R: dre, G: dgr, B: dbl, A: da})
	nowdate.TextStyle = fyne.TextStyle{Bold: true}
	nowdate.Alignment = fyne.TextAlignCenter
	nowdate.TextSize = float32(datesize)

	//background := canvas.NewRectangle(color.RGBA{R: 0, G: 0, B: 255, A: 255})
	// bgcolor := color.RGBA{R: 0, G: 143, B: 251, A: 255}
	bgcolor := color.RGBA{R: bre, G: bgr, B: bbl, A: ba}
	background := canvas.NewRectangle(bgcolor)

	vbox := container.NewVBox()
	if showutc == 1 {
		if showdate == 1 {
			vbox = container.NewVBox(nowtime, nowdate, utctime)
		} else {
			vbox = container.NewVBox(nowtime, utctime)
		}
	} else {
		if showdate == 1 {
			vbox = container.NewVBox(nowtime, nowdate)
		} else {
			vbox = container.NewVBox(nowtime)
		}
	}
	content := container.NewStack(background, vbox)
	// content := container.NewStack(background, container.NewVBox(nowtime, nowdate, utctime))

	updateClock := func() {
		now = time.Now()
		if now.Hour() == muteonhr && now.Minute() == muteonmin && now.Second() == 0 {
			if automute == 1 {
				muted, _ := volume.GetMuted()
				if !muted {
					currentvolume, _ = volume.GetVolume()
					volume.Mute()
				}
			}
		} else if now.Hour() == muteoffhr && now.Minute() == muteoffmin && now.Second() == 0 {
			if automute == 1 {
				muted, _ := volume.GetMuted()
				if muted {
					volume.Unmute()
					// volume.SetVolume(20)
					volume.SetVolume(currentvolume)
				}
			}
		}
		if now.Minute() == 0 && now.Second() == 0 {
			if hourchime == 1 {
				if !checkFileExists(sndDir + "/" + hourchimesound) {
					playBeep("updown")
				} else {
					playMp3(sndDir + "/" + hourchimesound)
				}
			}
		}

		nowtime.Text = now.Format(timeFormat)
		fyne.Do(func() {
			nowtime.Refresh()
			nowdate.Refresh()
		})
		nowdate.Text = now.Format(dateFormat)
		if showutc == 1 {
			utc := now.UTC()
			utctime.Text = utc.Format(utcFormat) + offsetString
			fyne.Do(func() {
				utctime.Refresh()
			})
		}
	}

	updateClock()
	go func() {
		for range time.Tick(time.Second) {
			// updating frequently is something of a resource hog (CPU)
			// check here if seconds are displayed, update
			// if seconds are not displayed, check for seconds == 0
			// at the minute change, and only update the clock then
			now = time.Now()
			if showseconds == 1 || now.Second() == 0 {
				updateClock()
			}
			// lock screen / mute volume event handler, but only if enabled
			// and only unmute if we auto muted. If user had already muted, don't
			if slockmute == 1 {
				if lock.IsScreenLocked() {
					muted, _ := volume.GetMuted()
					if !muted {
						clockmutedvol = 1
						volume.Mute()
					}
				} else {
					lockmuted, _ := volume.GetMuted()
					if lockmuted && clockmutedvol == 1 {
						clockmutedvol = 0
						volume.Unmute()
					}
				}
			}
		}
	}()

	clock.SetContent(content)
	clock.Resize(fyne.NewSize(content.MinSize().Width*1.2, content.MinSize().Height*1.1))
	// clock.Resize(fyne.NewSize(300, 200))
	clock.ShowAndRun()
	// clock.Show() // for func inside KrankyBearTimer
}

// "Now this is not the end. It is not even the beginning of the end. But it is, perhaps, the end of the beginning." Winston Churchill, November 10, 1942

// To-do:

// a few notes, format specific
// timeFormat := `3:04:05 PM (MST)`
// clock.SetText(now.Format("Mon Jan 2 15:04:05 2006"))
// clock.SetText(now.Format("15:04:05`nMonday, January 2, 2006"))

// show seconds
// clockFormat := `15:04:05
//Monday, January 2, 2006`

// no show seconds - not always valuable when we update every second
// anyway, but still - user preference ...
// clockFormat := `15:04`
//clockFormat := `15:04
//   Monday, January 2, 2006`
//clock.SetText(now.Format(clockFormat))
//clock.Alignment = fyne.TextAlignCenter
