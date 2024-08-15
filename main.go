package toascii

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"math"
	"strconv"
	"strings"

	"golang.org/x/image/draw"
)

const template = "\033[XXXm"

type ANSIString string

func createANSIString(parameters []int) ANSIString {
	var out = template
	for _, param := range parameters {
		out = strings.Replace(out, "XXX", strconv.Itoa(param)+";XXX", 1)
	}
	out = strings.Replace(out, ";XXX", "", 1)
	return ANSIString(out)
}

var RESET = createANSIString([]int{0})

const (
	BOLD              = 1
	FAINT             = 2 // Not widely supported
	ITALIC            = 3 // Not widely supported
	UNDERLINE         = 4
	SLOW_BLINK        = 5
	FAST_BLINK        = 6 // Not widely supported
	REVERSE_VIDEO     = 7 // Switches background and foreground colors
	CONCEAL           = 8 // Not widely supported
	STRIKETHROUGH     = 9 // Not widely supported
	DEFAULT_FONT      = 10
	ALT_FONT1         = 11
	ALT_FONT2         = 12
	ALT_FONT3         = 13
	ALT_FONT4         = 14
	ALT_FONT5         = 15
	ALT_FONT6         = 16
	ALT_FONT7         = 17
	ALT_FONT8         = 18
	ALT_FONT9         = 19
	FRAKTUR           = 20 // Almost never supported
	DOUBLE_UNDERLINE  = 21 // Almost never supported(same as BOLD_OFF)
	BOLD_OFF          = 21 // Not widely supported(same as DOUBLE_UNDERLINE)
	NORMAL_INTENSITY  = 22 // Default intensity(between BOLD and FAINT)
	NO_ITALIC_FRAKTUR = 23
	UNDERLINE_OFF     = 24 // Turns off UNDERLINE and DOUBLE_UNDERLINE
	BLINK_OFF         = 25
	INVERSE_OFF       = 27
	REVEAL            = 28 // CONCEAL off
	STRIKETHROUGH_OFF = 29
	FRAMED            = 51
	ENCIRCLED         = 52
	OVERLINED         = 53
	NO_FRAME_ENCIRCLE = 54 // Turns off FRAMED and ENCIRCLED
	NO_OVERLINED      = 55

	COLOR_BLACK   = 30
	COLOR_RED     = 31
	COLOR_GREEN   = 32
	COLOR_YELLOW  = 33
	COLOR_BLUE    = 34
	COLOR_MAGENTA = 35
	COLOR_CYAN    = 36
	COLOR_WHITE   = 37
	DEFAULT_COLOR = 39
)

func CreateColor(text_color int, background_color int, bright_text bool, bright_background bool) ANSIString {
	background_color += 10
	if bright_text {
		text_color += 60
	}
	if bright_background {
		background_color += 60
	}
	return createANSIString([]int{text_color, background_color})
}

type ConvertConfig struct {
	OutputWidth  int
	OutputHeight int
	Scaler       draw.Interpolator
}

var ascii_density_string = " `.-':_,^=;><+!rc*/z?sLTv)J7(|Fi{C}fI31tlu[neoZ5Yxjya]2ESwqkP6h9d4VpOGbUAKXHm8RD#$Bg0MNWQ%&@"
var density_factor = float64(len(ascii_density_string)) / float64(255)
var scaling_factor = float64(255) / float64(0xffff)

func Convert(imageFile io.Reader, config ConvertConfig) string {
	if config.Scaler == nil {
		config.Scaler = draw.NearestNeighbor
	}
	img, _, _ := image.Decode(imageFile)
	img = Resize(img, config.OutputWidth, config.OutputHeight, config.Scaler)

	var out = ""

	for y := 0; y < img.Bounds().Max.Y; y++ {
		for x := 0; x < img.Bounds().Max.X; x++ {
			var r, g, b, _ = img.At(x, y).RGBA()
			var R = float64(r) * scaling_factor
			var G = float64(g) * scaling_factor
			var B = float64(b) * scaling_factor
			var luminence = math.Sqrt(math.Pow(0.299*float64(R), 2) + math.Pow(0.587*float64(G), 2) + math.Pow(0.114*float64(B), 2))
			luminence *= density_factor
			luminence = math.Round(luminence)
			out += string(createANSIString([]int{38, 2, int(R), int(G), int(B)})) + string(ascii_density_string[int(luminence)])
		}
		out += "\n"
	}
	return out
}

func Resize(src image.Image, outputWidth int, outputHeight int, scaler draw.Interpolator) image.Image {
	dst := image.NewRGBA(image.Rect(0, 0, outputWidth, outputHeight))

	scaler.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)

	return dst
}
