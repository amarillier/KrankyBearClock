package main

import (
	"image/color"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"fyne.io/systray"
)

const (
	clockName      = "Tanium Clock"
	clockVersion   = "0.3" // see FyneApp.toml
	clockCopyright = "(c) Tanium, 2024"
	clockAuthor    = "Allan Marillier"
)

var imgDir string
var clockbg string // future optional clock background image

var sndDir string
var debug int = 0
var abt fyne.Window
var hlp fyne.Window
var bg fyne.Canvas

var showseconds int
var showtimezone int
var showdate int
var showutc int
var showhr12 int
var hourchime int
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

// preferences stored via fyne preferences API land in
// ~/Library/Preferences/fyne/com.tanium.taniumclock/preferences.json
// ~\AppData\Roaming\fyne\com.tanium.taniumclock\preferences.json
// {"bgcolor.default":"0,143,251,255","color_recents":"#eee53a,#83de4a,#f44336,#ffffff,#9c27b0,#8bc34a,#ff9800","datecolor.default":"131,222,74,255","datefont.default":"arial","datesize.default":24,"hourchime.default":1,"hourchimesound.default":"cuckoo.mp3","showdate.default":1,"showhr12.default":1,"showseconds.default":1,"showtimezone.default":1,"showutc.default":1,"timecolor.default":"255,123,31,255","timefont.default":"arial","timesize.default":48,"utccolor.default":"238,229,58,255","utcfont.default":"arial","utcsize.default":18}

func main() {
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}

	launchDir := filepath.Dir(exePath)
	if runtime.GOOS == "darwin" {
		if strings.HasPrefix(launchDir, "/Applications/TaniumClock") {
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

	a := app.NewWithID("com.tanium.TaniumClock")
	c := a.NewWindow(clockName)

	a.Settings().SetTheme(&appTheme{Theme: theme.DefaultTheme()})
	c.SetPadded(false)
	c.SetCloseIntercept(func() {
		a.Quit() // force quit, normal when somebody hits "x" to close
	})
	c.SetMaster() // this sets this as master and closes all child windows
	// c.CenterOnScreen() // run centered on primary (laptop) display

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
			r, _ := os.Open("TaniumClock0.txt")
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

	if desk, ok := a.(desktop.App); ok {
		desk.SetSystemTrayIcon(resourceTaniumClockSvg)
		systray.SetTooltip(clockName)
		// systray.SetTitle(clockName)
		show := fyne.NewMenuItem("Show", func() {
			c.Show()
			c.Canvas().Focused()
		})
		hide := fyne.NewMenuItem("Hide", c.Hide)
		about := fyne.NewMenuItem("About", func() {
			aboutText := clockName + " v " + clockVersion
			aboutText += "\n" + clockCopyright
			aboutText += "\n\n" + clockAuthor + ", using Go and fyne GUI"

			if abt == nil || !abt.Content().Visible() {
				abt = a.NewWindow(clockName + ": About")
				abt.Resize(fyne.NewSize(50, 100))
				abt.SetContent(widget.NewLabel(aboutText))
				abt.SetCloseIntercept(func() {
					abt.Close()
					abt = nil
				})
				// abt.CenterOnScreen() // run centered on primary (laptop) display
				abt.Show()
			} else {
				easterEgg(a, abt)
			}
		})
		help := fyne.NewMenuItem("Help", func() {
			// if hlp != nil { // &&  !hlp.Content().Visible() {
			if hlp == nil || !hlp.Content().Visible() {
				hlp = a.NewWindow(clockName + ": Help")

				hlp.SetCloseIntercept(func() {
					hlp.Close()
					hlp = nil
				})
				//}
				hlpText := `More help will be added later
For now we're adding as we go:
- This is a basic clock that currently shows
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

- See Settings Info tab for more detail on settings / preferences

- Default settings will be created on first run if they don't exist
`
				hlpText += "\n" + clockName + " v " + clockVersion
				hlpText += "\n" + clockCopyright
				hlpText += "\n\n" + clockAuthor + ", using Go and fyne GUI"

				plnText := `- Open with clock window focused
	- this is currently MacOS LaunchPad behavior, but only allows one app
	- To run more than one simultaneously, in terminal: open -n -a TaniumClock 
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
- Settings changes to background and clock default times are saved immediately.
	- but clock time format size and color, date size and color and background do
	not currently refresh to new settings - exit and rerun for now
	`
				link, err := url.Parse("https://www.tanium.com/end-user-license-agreement-policy")
				if err != nil {
					fyne.LogError("Could not parse URL", err)
				}
				hyperlink := widget.NewHyperlink("https://www.tanium.com/end-user-license-agreement-policy", link)
				hyperlink.Alignment = fyne.TextAlignLeading
				licText := `TaniumClock is “Beta Software” as defined in the license agreement found at the link below. 
Please take a moment to read the license agreement:
 
In addition, please note that:
TaniumClock is intended for internal Tanium use, however no proprietary
information or features are included, so pending Tanium legal and other 
approvals this application may be made available to others. TaniumClock
provides no guarantees as to stability of operations or suitability for any
purpose, but every attempt has been made to make this application reliable.

Using this application (and reading this text) is considered acceptance of
the terms of the License Agreement, and acknowledgement that this is Beta
Software and the additional terms above
`

				settingsText := `Settings are a separate tray menu item
Settings contains defaults as below, which can be modified, and also reset to defaults:
{"bgcolor.default":"0,143,251,255",
"color_recents":"#eee53a,#83de4a,#f44336,#ffffff,#9c27b0,#8bc34a,#ff9800",
"datecolor.default":"131,222,74,255","datefont.default":"arial","datesize.default":24,
"hourchime.default":1,"hourchimesound.default":"cuckoo.mp3","showdate.default":1,
"showhr12.default":1,"showseconds.default":1,"showtimezone.default":1,"showutc.default":1,
"timecolor.default":"255,123,31,255","timefont.default":"arial","timesize.default":48,
"utccolor.default":"238,229,58,255","utcfont.default":"arial","utcsize.default":18}

TaniumClock looks for directories named Resources/Images and Resources/Sounds,
containing images and sounds.

IMAGES:
Background blue refers to a compiled in resource with Tanium blue background. 
Other supported compiled in backgrounds are: stone, almond, converge24 and converge24a
Future additions will allow selecting images of your choice, png, SVG,
	jpg maybe and specifying size - height / width. Manual window resizing
	is already possible

SOUNDS:
Built in tones include 'ding', 'down', 'up', and 'updown'. These are always available
	and will be listed first in sound selectors
The sounds directory as distributed also contains a number of other .mp3 files
including baseball.mp3, grandfatherclock.mp3, hero.mp3, pinball.mp3, sosumi.mp3
When selecting sounds, the sound will be played as a preview when possible.
When selected sounds are not present (removed from Sounds), TaniumClock defaults
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
				easterEgg(a, hlp)
			}
		})
		settings := fyne.NewMenuItem("Settings", func() {
			makeSettings(a, c, bg)
		})
		menu := fyne.NewMenu(a.Metadata().Name, show, hide, fyne.NewMenuItemSeparator(), about, help, settings)
		desk.SetSystemTrayMenu(menu)
		systray.SetTooltip(clockName)
		// systray.SetTitle(clockName)
		//}
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

	utcFormat := `(UTC 3:04 PM Z07)`
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
		if now.Minute() == 0 && now.Second() == 0 {
			if hourchime == 1 {
				playMp3(sndDir + "/" + hourchimesound)
			}
		}
		nowtime.Text = now.Format(timeFormat)
		nowtime.Refresh()
		nowdate.Refresh()
		nowdate.Text = now.Format(dateFormat)
		if showutc == 1 {
			utctime.Text = now.Format(utcFormat)
			utctime.Refresh()
		}
	}

	updateClock()
	go func() {
		for range time.Tick(time.Second) {
			updateClock()
		}
	}()

	c.SetContent(content)
	c.Resize(fyne.NewSize(content.MinSize().Width*1.2, content.MinSize().Height*1.1))
	// c.Resize(fyne.NewSize(300, 200))
	c.ShowAndRun()
}

// "Now this is not the end. It is not even the beginning of the end. But it is, perhaps, the end of the beginning." Winston Churchill, November 10, 1942

// To-do:
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
