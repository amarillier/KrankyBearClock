package main

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"
	"time"

	"github.com/IamFaizanKhalid/lock"
	"github.com/go-vgo/robotgo"
	"github.com/itchyny/volume-go"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

// Clock display elements for dynamic updates
var nowtime *canvas.Text
var nowdate *canvas.Text
var utctime *canvas.Text
var background *canvas.Rectangle

// Additional timezone display elements
var timezone1Text *canvas.Text
var timezone2Text *canvas.Text
var timezone3Text *canvas.Text
var timezone4Text *canvas.Text
var timezone5Text *canvas.Text

// Weather display element
var weatherText *canvas.Text

// InitializeClockDisplay creates and initializes the clock display elements
// Returns the container content for the clock window
func InitializeClockDisplay(a fyne.App, clockWindow fyne.Window) fyne.CanvasObject {
	now := time.Now()

	// Build time format string based on settings
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
	utcFormat := `UTC 3:04 PM`
	nonutcFormat := `3:04 PM`
	dateFormat := ` Monday, January 2, 2006 `

	// Parse colors
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
	}

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
	}

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
	}

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
	}

	// Create clock display elements
	nowtime = canvas.NewText(now.Format(timeFormat), color.RGBA{R: tre, G: tgr, B: tbl, A: ta})
	nowtime.TextStyle = fyne.TextStyle{Bold: true}
	nowtime.Alignment = fyne.TextAlignCenter
	nowtime.TextSize = float32(timesize)

	utctime = canvas.NewText(now.Format(utcFormat), color.RGBA{R: ure, G: ugr, B: ubl, A: ua})
	utctime.TextStyle = fyne.TextStyle{Bold: true}
	utctime.Alignment = fyne.TextAlignCenter
	utctime.TextSize = float32(utcsize)

	nowdate = canvas.NewText(now.Format(dateFormat), color.RGBA{R: dre, G: dgr, B: dbl, A: da})
	nowdate.TextStyle = fyne.TextStyle{Bold: true}
	nowdate.Alignment = fyne.TextAlignCenter
	nowdate.TextSize = float32(datesize)

	bgcolorRect := color.RGBA{R: bre, G: bgr, B: bbl, A: ba}
	background = canvas.NewRectangle(bgcolorRect)

	// Initialize additional timezone text elements
	timezone1Text = nil
	timezone2Text = nil
	timezone3Text = nil
	timezone4Text = nil
	timezone5Text = nil
	weatherText = nil

	// Build container based on what's shown
	vbox := container.NewVBox()
	vbox.Add(nowtime)
	if showdate == 1 {
		vbox.Add(nowdate)
	}
	if showutc == 1 {
		vbox.Add(utctime)
	}

	// Add enabled additional timezones (either by name or UTC offset)
	if timezone1Enabled == 1 && (timezone1Name != "" || timezone1Offset != "") {
		timezone1Text = canvas.NewText("", color.RGBA{R: ure, G: ugr, B: ubl, A: ua})
		timezone1Text.TextStyle = fyne.TextStyle{Bold: true}
		timezone1Text.Alignment = fyne.TextAlignCenter
		timezone1Text.TextSize = float32(utcsize)
		vbox.Add(timezone1Text)
	}
	if timezone2Enabled == 1 && (timezone2Name != "" || timezone2Offset != "") {
		timezone2Text = canvas.NewText("", color.RGBA{R: ure, G: ugr, B: ubl, A: ua})
		timezone2Text.TextStyle = fyne.TextStyle{Bold: true}
		timezone2Text.Alignment = fyne.TextAlignCenter
		timezone2Text.TextSize = float32(utcsize)
		vbox.Add(timezone2Text)
	}
	if timezone3Enabled == 1 && (timezone3Name != "" || timezone3Offset != "") {
		timezone3Text = canvas.NewText("", color.RGBA{R: ure, G: ugr, B: ubl, A: ua})
		timezone3Text.TextStyle = fyne.TextStyle{Bold: true}
		timezone3Text.Alignment = fyne.TextAlignCenter
		timezone3Text.TextSize = float32(utcsize)
		vbox.Add(timezone3Text)
	}
	if timezone4Enabled == 1 && (timezone4Name != "" || timezone4Offset != "") {
		timezone4Text = canvas.NewText("", color.RGBA{R: ure, G: ugr, B: ubl, A: ua})
		timezone4Text.TextStyle = fyne.TextStyle{Bold: true}
		timezone4Text.Alignment = fyne.TextAlignCenter
		timezone4Text.TextSize = float32(utcsize)
		vbox.Add(timezone4Text)
	}
	if timezone5Enabled == 1 && (timezone5Name != "" || timezone5Offset != "") {
		timezone5Text = canvas.NewText("", color.RGBA{R: ure, G: ugr, B: ubl, A: ua})
		timezone5Text.TextStyle = fyne.TextStyle{Bold: true}
		timezone5Text.Alignment = fyne.TextAlignCenter
		timezone5Text.TextSize = float32(utcsize)
		vbox.Add(timezone5Text)
	}

	// Add weather temperature display if weather is enabled
	if weatherEnabled {
		weatherText = canvas.NewText("", color.RGBA{R: tre, G: tgr, B: tbl, A: ta})
		weatherText.TextStyle = fyne.TextStyle{Bold: true}
		weatherText.Alignment = fyne.TextAlignCenter
		weatherText.TextSize = float32(utcsize)
		vbox.Add(weatherText)
	}

	content := container.NewStack(background, vbox)

	// Start clock update loop
	StartClockUpdateLoop(offsetString, timeFormat, utcFormat, nonutcFormat, dateFormat)

	return content
}

// Helper functions moved outside updateClock to avoid recreation
var parseUTCOffset = func(offsetStr string) (hours int, minutes int, valid bool) {
	offsetStr = strings.TrimSpace(offsetStr)
	if offsetStr == "" {
		return 0, 0, false
	}
	// Check sign first
	negative := strings.HasPrefix(offsetStr, "-")
	// Remove leading + or - if present
	if strings.HasPrefix(offsetStr, "+") || strings.HasPrefix(offsetStr, "-") {
		offsetStr = offsetStr[1:]
	}
	// Parse the offset
	var offsetHours float64
	var err error
	if strings.Contains(offsetStr, ".") || strings.Contains(offsetStr, ",") {
		offsetStr = strings.ReplaceAll(offsetStr, ",", ".")
		offsetHours, err = strconv.ParseFloat(offsetStr, 64)
	} else {
		var h int
		h, err = strconv.Atoi(offsetStr)
		offsetHours = float64(h)
	}
	if err != nil {
		return 0, 0, false
	}
	// Apply sign
	if negative {
		offsetHours = -offsetHours
	}
	// Convert to hours and minutes
	hours = int(offsetHours)
	minutes = int((offsetHours - float64(hours)) * 60)
	if minutes < 0 {
		minutes = -minutes
	}
	return hours, minutes, true
}

var formatTimezoneTime = func(tzTime time.Time, tzName string, nonutcFormat string) string {
	_, tzOffset := tzTime.Zone()
	tzOffsetHours := tzOffset / 3600
	tzOffsetMinutes := (tzOffset % 3600) / 60
	tzOffsetString := fmt.Sprintf(" (%s %+02d:%02d)", tzName, tzOffsetHours, tzOffsetMinutes)
	return tzTime.Format(nonutcFormat) + tzOffsetString
}

var formatTimezoneTimeFromOffset = func(tzTime time.Time, offsetHours int, offsetMinutes int, offsetLabel string, nonutcFormat string) string {
	tzOffsetString := fmt.Sprintf(" (UTC%+02d:%02d)", offsetHours, offsetMinutes)
	if offsetLabel != "" {
		tzOffsetString = fmt.Sprintf(" (%s UTC%+02d:%02d)", offsetLabel, offsetHours, offsetMinutes)
	}
	return tzTime.Format(nonutcFormat) + tzOffsetString
}

// getTimezoneLocation gets a timezone location, using cache if available
func getTimezoneLocation(tzName string) *time.Location {
	if timezoneLocations == nil {
		timezoneLocations = make(map[string]*time.Location)
	}
	if loc, exists := timezoneLocations[tzName]; exists {
		return loc
	}
	if loc, err := time.LoadLocation(tzName); err == nil {
		timezoneLocations[tzName] = loc
		return loc
	}
	return nil
}

// StartClockUpdateLoop starts the background goroutine that updates the clock display
func StartClockUpdateLoop(offsetString, timeFormat, utcFormat, nonutcFormat, dateFormat string) {
	// Stop existing loop if running
	if clockUpdateLoopRunning && clockUpdateLoopStop != nil {
		close(clockUpdateLoopStop)
		clockUpdateLoopStop = nil
		clockUpdateLoopRunning = false
	}

	// Create new stop channel
	clockUpdateLoopStop = make(chan bool)
	clockUpdateLoopRunning = true

	updateClock := func() {
		now := time.Now()
		if automute == 1 {
			if now.Hour() == muteonhr && now.Minute() == muteonmin {
				k := wallClockMinuteKey(now)
				if k != lastAutomuteOnKey {
					lastAutomuteOnKey = k
					muted, _ := volume.GetMuted()
					jiggleconf = jiggle
					jiggle = 0 // disable jiggle while muted
					lastJiggleMinute = -1
					if !muted {
						currentvolume, _ = volume.GetVolume()
						volume.Mute()
					}
				}
			} else if now.Hour() == muteoffhr && now.Minute() == muteoffmin {
				k := wallClockMinuteKey(now)
				if k != lastAutomuteOffKey {
					lastAutomuteOffKey = k
					muted, _ := volume.GetMuted()
					jiggle = jiggleconf // restore jiggle value
					lastJiggleMinute = -1
					jiggleconf = 0
					if muted {
						volume.Unmute()
						volume.SetVolume(currentvolume)
					}
				}
			}
		}
		if now.Minute() == 0 {
			if hourchime == 1 {
				// Top of the hour: play once per local hour (no second==0; update may be late)
				currentHour := now.Hour()
				if lastChimeHour != currentHour {
					lastChimeHour = currentHour
					// Cache file existence check - only recheck if file changed
					chimeFilePath := sndDir + "/" + hourchimesound
					if !hourChimeFileChecked || hourChimeCachedFile != hourchimesound {
						hourChimeFileExists = checkFileExists(chimeFilePath)
						hourChimeFileChecked = true
						hourChimeCachedFile = hourchimesound
					}
					if !hourChimeFileExists {
						playBeep("updown")
					} else {
						playMp3(sndDir + "/" + hourchimesound)
					}
				}
			}
		}

		nowtime.Text = now.Format(timeFormat)
		fyne.Do(func() {
			nowtime.Refresh()
			nowdate.Refresh()
			// add here to also override when mute turns on/off
			// if screen is not locked and jiggle is on and minute modulo jiggle ...
			if jiggle == 0 {
				lastJiggleMinute = -1
			}
			if !lock.IsScreenLocked() && jiggle > 0 && now.Minute()%jiggle == 0 {
				if now.Minute() != lastJiggleMinute {
					robotgo.MoveRelative(1, 0)  // MoveSmoothRelative(200, 0)
					robotgo.MoveRelative(0, 1)  // MoveSmoothRelative(0, 200)
					robotgo.MoveRelative(-1, 0) // MoveSmoothRelative(-200, 0)
					robotgo.MoveRelative(0, -1) // MoveSmoothRelative(0, -200)
					lastJiggleMinute = now.Minute()
				}
			}
		})
		nowdate.Text = now.Format(dateFormat)
		utcTime := now.UTC()
		if showutc == 1 {
			utctime.Text = utcTime.Format(utcFormat) + offsetString
			fyne.Do(func() {
				utctime.Refresh()
			})
		}

		// Update additional timezones
		// Priority: UTC offset > timezone name
		if timezone1Enabled == 1 && timezone1Text != nil {
			if timezone1Offset != "" {
				// Use UTC offset
				hours, minutes, valid := parseUTCOffset(timezone1Offset)
				if valid {
					offsetDuration := time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute
					tzTime := utcTime.Add(offsetDuration)
					timezone1Text.Text = formatTimezoneTimeFromOffset(tzTime, hours, minutes, "Custom", nonutcFormat)
					fyne.Do(func() {
						timezone1Text.Refresh()
					})
				}
			} else if timezone1Name != "" {
				// Use timezone name (with caching)
				if loc := getTimezoneLocation(timezone1Name); loc != nil {
					tzTime := now.In(loc)
					timezone1Text.Text = formatTimezoneTime(tzTime, timezone1Name, nonutcFormat)
					fyne.Do(func() {
						timezone1Text.Refresh()
					})
				}
			}
		}
		if timezone2Enabled == 1 && timezone2Text != nil {
			if timezone2Offset != "" {
				hours, minutes, valid := parseUTCOffset(timezone2Offset)
				if valid {
					offsetDuration := time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute
					tzTime := utcTime.Add(offsetDuration)
					timezone2Text.Text = formatTimezoneTimeFromOffset(tzTime, hours, minutes, "Custom", nonutcFormat)
					fyne.Do(func() {
						timezone2Text.Refresh()
					})
				}
			} else if timezone2Name != "" {
				if loc := getTimezoneLocation(timezone2Name); loc != nil {
					tzTime := now.In(loc)
					timezone2Text.Text = formatTimezoneTime(tzTime, timezone2Name, nonutcFormat)
					fyne.Do(func() {
						timezone2Text.Refresh()
					})
				}
			}
		}
		if timezone3Enabled == 1 && timezone3Text != nil {
			if timezone3Offset != "" {
				hours, minutes, valid := parseUTCOffset(timezone3Offset)
				if valid {
					offsetDuration := time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute
					tzTime := utcTime.Add(offsetDuration)
					timezone3Text.Text = formatTimezoneTimeFromOffset(tzTime, hours, minutes, "Custom", nonutcFormat)
					fyne.Do(func() {
						timezone3Text.Refresh()
					})
				}
			} else if timezone3Name != "" {
				if loc := getTimezoneLocation(timezone3Name); loc != nil {
					tzTime := now.In(loc)
					timezone3Text.Text = formatTimezoneTime(tzTime, timezone3Name, nonutcFormat)
					fyne.Do(func() {
						timezone3Text.Refresh()
					})
				}
			}
		}
		if timezone4Enabled == 1 && timezone4Text != nil {
			if timezone4Offset != "" {
				hours, minutes, valid := parseUTCOffset(timezone4Offset)
				if valid {
					offsetDuration := time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute
					tzTime := utcTime.Add(offsetDuration)
					timezone4Text.Text = formatTimezoneTimeFromOffset(tzTime, hours, minutes, "Custom", nonutcFormat)
					fyne.Do(func() {
						timezone4Text.Refresh()
					})
				}
			} else if timezone4Name != "" {
				if loc := getTimezoneLocation(timezone4Name); loc != nil {
					tzTime := now.In(loc)
					timezone4Text.Text = formatTimezoneTime(tzTime, timezone4Name, nonutcFormat)
					fyne.Do(func() {
						timezone4Text.Refresh()
					})
				}
			}
		}
		if timezone5Enabled == 1 && timezone5Text != nil {
			if timezone5Offset != "" {
				hours, minutes, valid := parseUTCOffset(timezone5Offset)
				if valid {
					offsetDuration := time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute
					tzTime := utcTime.Add(offsetDuration)
					timezone5Text.Text = formatTimezoneTimeFromOffset(tzTime, hours, minutes, "Custom", nonutcFormat)
					fyne.Do(func() {
						timezone5Text.Refresh()
					})
				}
			} else if timezone5Name != "" {
				if loc := getTimezoneLocation(timezone5Name); loc != nil {
					tzTime := now.In(loc)
					timezone5Text.Text = formatTimezoneTime(tzTime, timezone5Name, nonutcFormat)
					fyne.Do(func() {
						timezone5Text.Refresh()
					})
				}
			}
		}

		// Update weather temperature display if enabled
		if weatherEnabled && weatherText != nil {
			if currentWeatherData != nil {
				// Format temperature: Fahrenheit first, then Celsius
				tempF := celsiusToFahrenheit(currentWeatherData.CurrentTemp)
				tempC := currentWeatherData.CurrentTemp
				weatherText.Text = fmt.Sprintf("Weather: %.1f°F / %.1f°C", tempF, tempC)
			} else {
				weatherText.Text = "Weather: Loading..."
			}
			fyne.Do(func() {
				weatherText.Refresh()
			})
		}
	}

	updateClock()
	go func() {
		// When seconds are hidden, run the full clock update once per wall-clock minute
		// (first 1s tick that sees the new minute). Relying only on local second==0 misses
		// updates if the goroutine wakes late (sleep, scheduling), which breaks jiggle/chimes.
		lastSlowClockMinuteEpoch := time.Now().Unix() / 60
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				now := time.Now()
				if showseconds == 1 {
					updateClock()
				} else {
					minuteEpoch := now.Unix() / 60
					if minuteEpoch != lastSlowClockMinuteEpoch {
						lastSlowClockMinuteEpoch = minuteEpoch
						updateClock()
					}
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
			case <-clockUpdateLoopStop:
				return
			}
		}
	}()
}

// "Now this is not the end. It is not even the beginning of the end. But it is, perhaps, the end of the beginning." Winston Churchill, November 10, 1942
