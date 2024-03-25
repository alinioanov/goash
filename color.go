package goash

var colors = map[string]string{
	"red":    "\033[31m",
	"green":  "\033[32m",
	"yellow": "\033[33m",
	"bold":   "\033[1m",
	"white":  "\033[37m",
	"end":    "\033[m",
}

// color returns a color wrapped string.
func color(s, color string) string {
	return colors[color] + s + colors["end"]
}
