audio.go:func playMp3(name string) {
audio.go:func playMid(name string) {
audio.go:func playWav(name string) {
audio.go:func playBeep(style string) {
main.go:func main() {
main.go:// anyway, but still - user preference ...
settings.go:func makeSettingsClock(a fyne.App, w fyne.Window, bg fyne.Canvas) {
settings.go:func makeSettingsTheme(a fyne.App, w fyne.Window, bg fyne.Canvas) {
settings.go:func writeDefaultSettings(a fyne.App) {
settings.go:func writeSettings(a fyne.App) {
settings.go:func colorPicker(parent fyne.Window, s string, a fyne.App) color.Color {
settings.go:func colorSelected(c color.Color, w fyne.Window, s string, a fyne.App) {
settings.go:func ColorToString(c color.Color) string {
settings.go:func showFilePicker(w fyne.Window) {
settings.go:func selectTime(a fyne.App, w fyne.Window, bg fyne.Canvas, caller string, hr int, min int) string {
settings.go:func isValidTime(t string) bool {
theme.go:func (a *appTheme) Size(n fyne.ThemeSizeName) float32 {
theme.go:func themeTimer(text *widget.RichText, time int) {
util.go:func checkFileExists(filePath string) bool {
util.go:func daysUntil(targetDate string) (int, error) {
util.go:func logInit() {
util.go:func lineCounter(r io.Reader) (int, error) {
util.go:func logRotate() {
util.go:func easterEgg(a fyne.App, w fyne.Window) {
util.go:func teapot(a fyne.App, w fyne.Window) {
util.go:func listMatchingFiles(directory, pattern string) ([]string, error) {
util.go:func dadjoke() string {
util.go:func isProcessRunning(processRegex string) bool {
util.go:func updateChecker(repoOwner string, repo string, repoName string, repodl string) (string, bool) {
