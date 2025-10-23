package colors

import "fmt"

// ANSI color codes
const (
	AnsiReset  = "\033[0m"
	AnsiRed    = "\033[31m"
	AnsiGreen  = "\033[32m"
	AnsiYellow = "\033[33m"
	AnsiBlue   = "\033[34m"
	AnsiPurple = "\033[35m"
	AnsiCyan   = "\033[36m"
	AnsiWhite  = "\033[37m"
	AnsiGray   = "\033[90m"

	// Bright colors
	AnsiBrightRed    = "\033[91m"
	AnsiBrightGreen  = "\033[92m"
	AnsiBrightYellow = "\033[93m"
	AnsiBrightBlue   = "\033[94m"
	AnsiBrightPurple = "\033[95m"
	AnsiBrightCyan   = "\033[96m"
	AnsiBrightWhite  = "\033[97m"
)

// Color functions
func Red(s string) string          { return AnsiRed + s + AnsiReset }
func Green(s string) string        { return AnsiGreen + s + AnsiReset }
func Yellow(s string) string       { return AnsiYellow + s + AnsiReset }
func Blue(s string) string         { return AnsiBlue + s + AnsiReset }
func Purple(s string) string       { return AnsiPurple + s + AnsiReset }
func Cyan(s string) string         { return AnsiCyan + s + AnsiReset }
func White(s string) string        { return AnsiWhite + s + AnsiReset }
func Gray(s string) string         { return AnsiGray + s + AnsiReset }
func BrightRed(s string) string    { return AnsiBrightRed + s + AnsiReset }
func BrightGreen(s string) string  { return AnsiBrightGreen + s + AnsiReset }
func BrightYellow(s string) string { return AnsiBrightYellow + s + AnsiReset }
func BrightBlue(s string) string   { return AnsiBrightBlue + s + AnsiReset }
func BrightPurple(s string) string { return AnsiBrightPurple + s + AnsiReset }
func BrightCyan(s string) string   { return AnsiBrightCyan + s + AnsiReset }
func BrightWhite(s string) string  { return AnsiBrightWhite + s + AnsiReset }

// Special formatting
func Bold(s string) string      { return "\033[1m" + s + AnsiReset }
func Italic(s string) string    { return "\033[3m" + s + AnsiReset }
func Underline(s string) string { return "\033[4m" + s + AnsiReset }

// Predefined styled messages
func Success(s string) string { return Green("✅ " + s) }
func Error(s string) string   { return Red("❌ " + s) }
func Warning(s string) string { return Yellow("⚠️  " + s) }
func Info(s string) string    { return Blue("ℹ️  " + s) }
func Rocket(s string) string  { return Cyan("🚀 " + s) }
func Globe(s string) string   { return Purple("🌐 " + s) }
func Stop(s string) string    { return Red("🛑 " + s) }
func Check(s string) string   { return Green("✓ " + s) }
func Cross(s string) string   { return Red("✗ " + s) }

// Print functions
func PrintSuccess(s string) { fmt.Print(Success(s)) }
func PrintError(s string)   { fmt.Print(Error(s)) }
func PrintWarning(s string) { fmt.Print(Warning(s)) }
func PrintInfo(s string)    { fmt.Print(Info(s)) }
func PrintRocket(s string)  { fmt.Print(Rocket(s)) }
func PrintGlobe(s string)   { fmt.Print(Globe(s)) }
func PrintStop(s string)    { fmt.Print(Stop(s)) }
func PrintCheck(s string)   { fmt.Print(Check(s)) }
func PrintCross(s string)   { fmt.Print(Cross(s)) }

// Printf functions
func PrintfSuccess(format string, args ...interface{}) { fmt.Printf(Success(format), args...) }
func PrintfError(format string, args ...interface{})   { fmt.Printf(Error(format), args...) }
func PrintfWarning(format string, args ...interface{}) { fmt.Printf(Warning(format), args...) }
func PrintfInfo(format string, args ...interface{})    { fmt.Printf(Info(format), args...) }
func PrintfRocket(format string, args ...interface{})  { fmt.Printf(Rocket(format), args...) }
func PrintfGlobe(format string, args ...interface{})   { fmt.Printf(Globe(format), args...) }
func PrintfStop(format string, args ...interface{})    { fmt.Printf(Stop(format), args...) }
func PrintfCheck(format string, args ...interface{})   { fmt.Printf(Check(format), args...) }
func PrintfCross(format string, args ...interface{})   { fmt.Printf(Cross(format), args...) }
