package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
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
	// audio "github.com/amarillier/KrankyBearModule/audio"
	// audio "github.com/amarillier/KrankyBearModule/util"
)

const (
	// appName    = "Kranky Bear Clock"
	appVersion = "0.4.5" // see FyneApp.toml
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
var countdown fyne.Window

// alarmWindow is declared in alarm.go

// Countdown window elements for dynamic color updates
var countdownTitleText *canvas.Text
var countdownHelpText *canvas.Text
var countdownBackground *canvas.Rectangle
var countdownDaysText1 *canvas.Text
var countdownDaysText2 *canvas.Text
var countdownDaysText3 *canvas.Text

// var egg fyne.Window
var bg fyne.Canvas

// Clock display elements are declared in clock.go

var showseconds int
var showtimezone int
var showdate int
var showutc int
var showhr12 int
var hourchime int
var slockmute int
var clockmutedvol int
var automute int
var jiggle int
var jiggleconf int
var lastJiggleMinute = -1
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
var lastChimeHour int = -1                      // Track last hour when chime played to prevent double playback
var clockUpdateLoopRunning bool = false         // Prevent multiple update loops from running
var clockUpdateLoopStop chan bool               // Channel to stop the update loop
var timezoneLocations map[string]*time.Location // Cache timezone locations
var hourChimeFileExists bool = false            // Cache file existence check
var hourChimeFileChecked bool = false           // Track if we've checked the file
var hourChimeCachedFile string = ""             // Track which file was cached
var startclock int
var processName string
var prefs string

// Countdown dates (up to 3)
var countdownDate1 string
var countdownDesc1 string
var countdownDate2 string
var countdownDesc2 string
var countdownDate3 string
var countdownDesc3 string

// Additional timezones (up to 5)
var timezone1Enabled int
var timezone1Name string
var timezone1Offset string // UTC offset (e.g., "+5", "-3.5")
var timezone2Enabled int
var timezone2Name string
var timezone2Offset string
var timezone3Enabled int
var timezone3Name string
var timezone3Offset string
var timezone4Enabled int
var timezone4Name string
var timezone4Offset string
var timezone5Enabled int
var timezone5Name string
var timezone5Offset string

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
	// Initialize speaker at startup for faster audio
	initSpeaker()
	a.Settings().SetTheme(&appTheme{Theme: theme.DefaultTheme()})
	clock = a.NewWindow(appName)
	_, month, _ := time.Now().Date()
	if month == time.December {
		clock.SetIcon(resourceKrankyBearChristmasGrinchPng)
	} else {
		clock.SetIcon(resourceKrankyBearBeretPng)
	}
	clock.SetPadded(false)
	//clock.SetCloseIntercept(func() {
	//	a.Quit() // force quit, normal when somebody hits "x" to close
	//})
	clock.SetMaster() // this sets this as master and closes all child windows
	// clock.CenterOnScreen() // run centered on primary (laptop) display

	prefs = strings.ReplaceAll((a.Storage().RootURI()).String(), "file://", "") + "/preferences.json"
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
	jiggle = a.Preferences().IntWithFallback("jiggle.default", 0)
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

	// Load countdown dates
	countdownDate1 = a.Preferences().StringWithFallback("countdown.date1", "")
	countdownDesc1 = a.Preferences().StringWithFallback("countdown.desc1", "")
	countdownDate2 = a.Preferences().StringWithFallback("countdown.date2", "")
	countdownDesc2 = a.Preferences().StringWithFallback("countdown.desc2", "")
	countdownDate3 = a.Preferences().StringWithFallback("countdown.date3", "")
	countdownDesc3 = a.Preferences().StringWithFallback("countdown.desc3", "")

	// Load additional timezones
	timezone1Enabled = a.Preferences().IntWithFallback("timezone1.enabled", 0)
	timezone1Name = a.Preferences().StringWithFallback("timezone1.name", "")
	timezone1Offset = a.Preferences().StringWithFallback("timezone1.offset", "")
	timezone2Enabled = a.Preferences().IntWithFallback("timezone2.enabled", 0)
	timezone2Name = a.Preferences().StringWithFallback("timezone2.name", "")
	timezone2Offset = a.Preferences().StringWithFallback("timezone2.offset", "")
	timezone3Enabled = a.Preferences().IntWithFallback("timezone3.enabled", 0)
	timezone3Name = a.Preferences().StringWithFallback("timezone3.name", "")
	timezone3Offset = a.Preferences().StringWithFallback("timezone3.offset", "")
	timezone4Enabled = a.Preferences().IntWithFallback("timezone4.enabled", 0)
	timezone4Name = a.Preferences().StringWithFallback("timezone4.name", "")
	timezone4Offset = a.Preferences().StringWithFallback("timezone4.offset", "")
	timezone5Enabled = a.Preferences().IntWithFallback("timezone5.enabled", 0)
	timezone5Name = a.Preferences().StringWithFallback("timezone5.name", "")
	timezone5Offset = a.Preferences().StringWithFallback("timezone5.offset", "")

	// Load alarms and start alarm checker
	loadAlarms(a)
	startAlarmChecker(a)

	// Load weather settings and start weather refresh if enabled
	loadWeatherSettings(a)
	startWeatherRefresh(a)

	writeSettings(a)

	clockmutedvol = 0

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
		updateAlert(a, updtmsg)
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

			kb := canvas.NewImageFromResource(resourceKrankyBearBeretPng)
			_, month, _ := time.Now().Date()
			if month == time.December {
				kb = canvas.NewImageFromResource(resourceKrankyBearChristmasGrinchPng)
			}
			text := widget.NewLabel(aboutText)
			kb.FillMode = canvas.ImageFillOriginal
			content := container.NewHBox(kb, text)

			if abt == nil || !abt.Content().Visible() {
				abt = a.NewWindow(appName + ": About")
				_, month, _ := time.Now().Date()
				if month == time.December {
					abt.SetIcon(resourceKrankyBearChristmasGrinchPng)
				} else {
					abt.SetIcon(resourceKrankyBearBeretPng)
				}
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
				_, month, _ := time.Now().Date()
				if month == time.December {
					hlp.SetIcon(resourceKrankyBearChristmasGrinchPng)
				} else {
					hlp.SetIcon(resourceKrankyBearBeretPng)
				}

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
- optional additional time zones (up to 5)
- optional hourly chime with user selectable sound
	- ideally chose a very short chime sound, less than half second
- customizable font sizes for each of time, date and UTC time
- customizable font color for each of background, time, date and UTC time
- clock display window resizes automatically to suit selected font sizes
- optional setting to enable auto starting at boot
- countdown days to future dates - up to 3 future dates

- Note: Displaying seconds can be quite resource intensive with clock display updates every second. 
  The app can be substantially less CPU intensive when seconds are not displayed, allowing the app to
  refresh the display every minute rather than every second

- See Settings Info tab for more detail on settings / preferences

- Default settings will be created on first run if they don't exist
`
				hlpText += "\n" + appName + " v " + appVersion
				hlpText += "\n" + appCopyright
				hlpText += "\n\n" + appAuthor + ", using Go and fyne GUI"

				plnText := `- Allow settings set/save window locations to open clock, 
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
- Settings changes are saved immediately and clock display colors and sizes now update automatically.
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
		countdownDates := fyne.NewMenuItem("Countdown Dates", func() {
			makeCountdownDates(a, clock, bg)
		})
		alarmsMenu := fyne.NewMenuItem("Alarms", func() {
			makeAlarmsWindow(a, clock, bg)
		})
		weatherMenu := fyne.NewMenuItem("Weather", func() {
			makeWeatherWindow(a, clock, bg)
		})
		prefsEdit := fyne.NewMenuItem("Preferences manual edit", func() {
			var cmd *exec.Cmd

			switch runtime.GOOS {
			case "windows":
				cmd = exec.Command("cmd", "/d", "/c", "start", prefs)
			case "darwin": // macOS
				cmd = exec.Command("open", prefs)
			case "linux":
				cmd = exec.Command("xdg-open", prefs)
			default:
				fmt.Printf("Unsupported operating system: %s\n", runtime.GOOS)
				return
			}
			err := cmd.Run()
			if err != nil {
				playBeep("down")
			}
		})
		updtchk := fyne.NewMenuItem("Check for update", func() {
			// throw away updateAvail here, use _, unneeded for manual check
			updtmsg, _ := updateChecker("amarillier", "KrankyBearClock", "Kranky Bear Clock", "https://github.com/amarillier/KrankyBearClock/releases/latest")
			if updt == nil {
				updateAlert(a, updtmsg)
			} else {
				updt.RequestFocus()
			}
		})
		menu := fyne.NewMenu(a.Metadata().Name, show, hide, alarmsMenu, countdownDates, weatherMenu, fyne.NewMenuItemSeparator(), about, updtchk, help, settingsClock, settingsTheme, prefsEdit)
		desk.SetSystemTrayMenu(menu)
		_, month, _ := time.Now().Date()
		if month == time.December {
			desk.SetSystemTrayIcon(resourceKrankyBearChristmasGrinchPng)
		} else {
			desk.SetSystemTrayIcon(resourceKrankyBearBeretPng)
		}
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
		newMenuOps := fyne.NewMenu("Operations", show, hide, alarmsMenu, countdownDates, weatherMenu, fyne.NewMenuItemSeparator(), quit)
		newMenuHelp := fyne.NewMenu("Help", about, updtchk, help)
		newMenuSettings := fyne.NewMenu("Settings", settingsClock, settingsTheme, prefsEdit)
		// New main menu
		cmenu := fyne.NewMainMenu(newMenuOps, newMenuHelp, newMenuSettings)
		// setup main menu
		clock.SetMainMenu(cmenu)
		// cmenu.Refresh()
	}

	// Initialize clock display
	content := InitializeClockDisplay(a, clock)
	clock.SetContent(content)
	clock.Resize(fyne.NewSize(content.MinSize().Width*1.2, content.MinSize().Height*1.1))
	// clock.Resize(fyne.NewSize(300, 200))
	clock.ShowAndRun()
	// clock.Show() // for func inside KrankyBearTimer
}

func updateAlert(a fyne.App, updtmsg string) {
	// open a window to show the update message
	// no need to test for updt window open at first start
	var kbimg *canvas.Image
	releaselink, rerr := url.Parse("https://github.com/amarillier/KrankyBearClock/releases/latest")
	if rerr != nil {
		fyne.LogError("Could not parse URL", rerr)
	}
	myreleaselink := widget.NewHyperlink("https://github.com/amarillier/KrankyBearClock/releases/latest", releaselink)
	myreleaselink.Alignment = fyne.TextAlignLeading

	releasenoteslink, rnerr := url.Parse("https://github.com/amarillier/KrankyBearClock/blob/main/ReleaseNotes.txt")
	if rnerr != nil {
		fyne.LogError("Could not parse URL", rnerr)
	}
	myreleasenoteslink := widget.NewHyperlink("https://github.com/amarillier/KrankyBearClock/blob/main/ReleaseNotes.txt", releasenoteslink)
	myreleasenoteslink.Alignment = fyne.TextAlignLeading

	if strings.Contains(updtmsg, "newer version") {
		kbimg = canvas.NewImageFromResource(resourceKrankyBearHardHatPng)
		kbimg.FillMode = canvas.ImageFillOriginal
	} else if strings.Contains(updtmsg, "running the latest") {
		kbimg = canvas.NewImageFromResource(resourceKrankyBearBeretPng)
		kbimg.FillMode = canvas.ImageFillOriginal
	} else {
		alert := sndDir + "/KrankyBearGrowl.mp3"
		alert = sndDir + "/uhOh.mp3"
		if !checkFileExists(alert) {
			playBeep("up")
		} else {
			playMp3(alert) // Basso, Blow, Hero, Funk, Glass, Ping, Purr, Sosumi, Submarine,
		}
		kbimg = canvas.NewImageFromResource(resourceKrankyBearVikingHelmetPng)
		kbimg.FillMode = canvas.ImageFillOriginal
	}

	text := widget.NewLabel(updtmsg)
	content := container.NewVBox(kbimg, text, myreleaselink, myreleasenoteslink)
	updt = a.NewWindow(appName + ": Update Check")
	_, month, _ := time.Now().Date()
	if month == time.December {
		updt.SetIcon(resourceKrankyBearChristmasGrinchPng)
	} else {
		updt.SetIcon(resourceKrankyBearBeretPng)
	}
	updt.Resize(fyne.NewSize(50, 100))
	updt.SetContent(content)
	updt.SetCloseIntercept(func() {
		updt.Close()
		updt = nil
	})
	// updt.CenterOnScreen() // run centered on primary (laptop) display
	updt.Show()
}

// "Now this is not the end. It is not even the beginning of the end. But it is, perhaps, the end of the beginning." Winston Churchill, November 10, 1942
