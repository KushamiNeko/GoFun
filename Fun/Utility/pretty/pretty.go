package pretty

import (
	"fmt"
	"image/color"
	"math"
	"regexp"
	"strconv"
	"strings"
)

func checkHexColor(s string) string {
	re := regexp.MustCompile(`^#?([0-9a-f]{6})$`)

	match := re.FindAllStringSubmatch(s, -1)

	if match == nil {
		panic(fmt.Sprintf("invalid hex color string: %s", s))
	}

	return match[0][1]
}

func ColorCaps(color string, message string) {
	ColorPrintln(color, strings.ToUpper(message))
}

func ColorPrintln(color string, message string) {
	r, g, b := HexToRGB8(color)

	fmt.Printf("\033[38;2;%d;%d;%dm%s\033[0m\n", r, g, b, message)
}

func HexToRGB8(color string) (r, g, b uint8) {
	color = checkHexColor(color)

	ri, err := strconv.ParseUint(color[0:2], 16, 8)
	if err != nil {
		panic(err)
	}

	gi, err := strconv.ParseUint(color[2:4], 16, 8)
	if err != nil {
		panic(err)
	}

	bi, err := strconv.ParseUint(color[4:6], 16, 8)
	if err != nil {
		panic(err)
	}

	return uint8(ri), uint8(gi), uint8(bi)
}

func HexToHSV(color string) (h, s, v float64) {
	color = checkHexColor(color)

	r, g, b := HexToRGB8(color)

	return RGB32ToHSV32(
		float64(r)/float64(math.MaxUint8),
		float64(g)/float64(math.MaxUint8),
		float64(b)/float64(math.MaxUint8),
	)
}

func RGB8ToHSV32(r, g, b uint8) (float64, float64, float64) {
	var fr, fg, fb float64

	fr = float64(r) / float64(math.MaxUint8)
	fg = float64(g) / float64(math.MaxUint8)
	fb = float64(b) / float64(math.MaxUint8)

	return RGB32ToHSV32(fr, fg, fb)
}

func RGB32ToHSV32(fr, fg, fb float64) (float64, float64, float64) {

	minRGB := math.Min(fr, math.Min(fg, fb))
	maxRGB := math.Max(fr, math.Max(fg, fb))

	if minRGB == maxRGB {
		return 0.0, 0.0, minRGB
	}

	var d float64
	switch {
	case fr == minRGB:
		d = fg - fb
	case fb == minRGB:
		d = fr - fg
	default:
		d = fb - fr
	}

	var dh float64
	switch {
	case fr == minRGB:
		dh = 3.0
	case fb == minRGB:
		dh = 1.0
	default:
		dh = 5.0
	}

	var h, s, v float64
	h = (60.0 * (dh - d/(maxRGB-minRGB))) / 360.0
	s = (maxRGB - minRGB) / maxRGB
	v = maxRGB

	return h, s, v
}

func HexToColor(hc string, alpha float64) color.Color {
	r, g, b := HexToRGB8(hc)

	var a float64
	a = math.Max(0.0, alpha)
	a = math.Min(1.0, a)
	a = math.Floor(a * 255.0)

	return color.NRGBA{R: r, G: g, B: b, A: uint8(a)}
}

const (
	resetAll = "\033[0m"

	bold       = "\033[1m"
	dim        = "\033[2m"
	underlined = "\033[4m"
	blink      = "\033[5m"
	reverse    = "\033[7m"
	hidden     = "\033[8m"

	GoogleRed100    = "f4c7c3"
	GoogleRed300    = "e67c73"
	GoogleRed500    = "db4437"
	GoogleRed700    = "c53929"
	GoogleBlue100   = "c6dafc"
	GoogleBlue300   = "7baaf7"
	GoogleBlue500   = "4285f4"
	GoogleBlue700   = "3367d6"
	GoogleGreen100  = "b7e1cd"
	GoogleGreen300  = "57bb8a"
	GoogleGreen500  = "0f9d58"
	GoogleGreen700  = "0b8043"
	GoogleYellow100 = "fce8b2"
	GoogleYellow300 = "f7cb4d"
	GoogleYellow500 = "f4b400"
	GoogleYellow700 = "f09300"
	GoogleGrey100   = "f5f5f5"
	GoogleGrey300   = "e0e0e0"
	GoogleGrey500   = "9e9e9e"
	GoogleGrey700   = "616161"

	PaperRed50          = "ffebee"
	PaperRed100         = "ffcdd2"
	PaperRed200         = "ef9a9a"
	PaperRed300         = "e57373"
	PaperRed400         = "ef5350"
	PaperRed500         = "f44336"
	PaperRed600         = "e53935"
	PaperRed700         = "d32f2f"
	PaperRed800         = "c62828"
	PaperRed900         = "b71c1c"
	PaperRedA100        = "ff8a80"
	PaperRedA200        = "ff5252"
	PaperRedA400        = "ff1744"
	PaperRedA700        = "d50000"
	PaperPink50         = "fce4ec"
	PaperPink100        = "f8bbd0"
	PaperPink200        = "f48fb1"
	PaperPink300        = "f06292"
	PaperPink400        = "ec407a"
	PaperPink500        = "e91e63"
	PaperPink600        = "d81b60"
	PaperPink700        = "c2185b"
	PaperPink800        = "ad1457"
	PaperPink900        = "880e4f"
	PaperPinkA100       = "ff80ab"
	PaperPinkA200       = "ff4081"
	PaperPinkA400       = "f50057"
	PaperPinkA700       = "c51162"
	PaperPurple50       = "f3e5f5"
	PaperPurple100      = "e1bee7"
	PaperPurple200      = "ce93d8"
	PaperPurple300      = "ba68c8"
	PaperPurple400      = "ab47bc"
	PaperPurple500      = "9c27b0"
	PaperPurple600      = "8e24aa"
	PaperPurple700      = "7b1fa2"
	PaperPurple800      = "6a1b9a"
	PaperPurple900      = "4a148c"
	PaperPurpleA100     = "ea80fc"
	PaperPurpleA200     = "e040fb"
	PaperPurpleA400     = "d500f9"
	PaperPurpleA700     = "aa00ff"
	PaperDeepPurple50   = "ede7f6"
	PaperDeepPurple100  = "d1c4e9"
	PaperDeepPurple200  = "b39ddb"
	PaperDeepPurple300  = "9575cd"
	PaperDeepPurple400  = "7e57c2"
	PaperDeepPurple500  = "673ab7"
	PaperDeepPurple600  = "5e35b1"
	PaperDeepPurple700  = "512da8"
	PaperDeepPurple800  = "4527a0"
	PaperDeepPurple900  = "311b92"
	PaperDeepPurpleA100 = "b388ff"
	PaperDeepPurpleA200 = "7c4dff"
	PaperDeepPurpleA400 = "651fff"
	PaperDeepPurpleA700 = "6200ea"
	PaperIndigo50       = "e8eaf6"
	PaperIndigo100      = "c5cae9"
	PaperIndigo200      = "9fa8da"
	PaperIndigo300      = "7986cb"
	PaperIndigo400      = "5c6bc0"
	PaperIndigo500      = "3f51b5"
	PaperIndigo600      = "3949ab"
	PaperIndigo700      = "303f9f"
	PaperIndigo800      = "283593"
	PaperIndigo900      = "1a237e"
	PaperIndigoA100     = "8c9eff"
	PaperIndigoA200     = "536dfe"
	PaperIndigoA400     = "3d5afe"
	PaperIndigoA700     = "304ffe"
	PaperBlue50         = "e3f2fd"
	PaperBlue100        = "bbdefb"
	PaperBlue200        = "90caf9"
	PaperBlue300        = "64b5f6"
	PaperBlue400        = "42a5f5"
	PaperBlue500        = "2196f3"
	PaperBlue600        = "1e88e5"
	PaperBlue700        = "1976d2"
	PaperBlue800        = "1565c0"
	PaperBlue900        = "0d47a1"
	PaperBlueA100       = "82b1ff"
	PaperBlueA200       = "448aff"
	PaperBlueA400       = "2979ff"
	PaperBlueA700       = "2962ff"
	PaperLightBlue50    = "e1f5fe"
	PaperLightBlue100   = "b3e5fc"
	PaperLightBlue200   = "81d4fa"
	PaperLightBlue300   = "4fc3f7"
	PaperLightBlue400   = "29b6f6"
	PaperLightBlue500   = "03a9f4"
	PaperLightBlue600   = "039be5"
	PaperLightBlue700   = "0288d1"
	PaperLightBlue800   = "0277bd"
	PaperLightBlue900   = "01579b"
	PaperLightBlueA100  = "80d8ff"
	PaperLightBlueA200  = "40c4ff"
	PaperLightBlueA400  = "00b0ff"
	PaperLightBlueA700  = "0091ea"
	PaperCyan50         = "e0f7fa"
	PaperCyan100        = "b2ebf2"
	PaperCyan200        = "80deea"
	PaperCyan300        = "4dd0e1"
	PaperCyan400        = "26c6da"
	PaperCyan500        = "00bcd4"
	PaperCyan600        = "00acc1"
	PaperCyan700        = "0097a7"
	PaperCyan800        = "00838f"
	PaperCyan900        = "006064"
	PaperCyanA100       = "84ffff"
	PaperCyanA200       = "18ffff"
	PaperCyanA400       = "00e5ff"
	PaperCyanA700       = "00b8d4"
	PaperTeal50         = "e0f2f1"
	PaperTeal100        = "b2dfdb"
	PaperTeal200        = "80cbc4"
	PaperTeal300        = "4db6ac"
	PaperTeal400        = "26a69a"
	PaperTeal500        = "009688"
	PaperTeal600        = "00897b"
	PaperTeal700        = "00796b"
	PaperTeal800        = "00695c"
	PaperTeal900        = "004d40"
	PaperTealA100       = "a7ffeb"
	PaperTealA200       = "64ffda"
	PaperTealA400       = "1de9b6"
	PaperTealA700       = "00bfa5"
	PaperGreen50        = "e8f5e9"
	PaperGreen100       = "c8e6c9"
	PaperGreen200       = "a5d6a7"
	PaperGreen300       = "81c784"
	PaperGreen400       = "66bb6a"
	PaperGreen500       = "4caf50"
	PaperGreen600       = "43a047"
	PaperGreen700       = "388e3c"
	PaperGreen800       = "2e7d32"
	PaperGreen900       = "1b5e20"
	PaperGreenA100      = "b9f6ca"
	PaperGreenA200      = "69f0ae"
	PaperGreenA400      = "00e676"
	PaperGreenA700      = "00c853"
	PaperLightGreen50   = "f1f8e9"
	PaperLightGreen100  = "dcedc8"
	PaperLightGreen200  = "c5e1a5"
	PaperLightGreen300  = "aed581"
	PaperLightGreen400  = "9ccc65"
	PaperLightGreen500  = "8bc34a"
	PaperLightGreen600  = "7cb342"
	PaperLightGreen700  = "689f38"
	PaperLightGreen800  = "558b2f"
	PaperLightGreen900  = "33691e"
	PaperLightGreenA100 = "ccff90"
	PaperLightGreenA200 = "b2ff59"
	PaperLightGreenA400 = "76ff03"
	PaperLightGreenA700 = "64dd17"
	PaperLime50         = "f9fbe7"
	PaperLime100        = "f0f4c3"
	PaperLime200        = "e6ee9c"
	PaperLime300        = "dce775"
	PaperLime400        = "d4e157"
	PaperLime500        = "cddc39"
	PaperLime600        = "c0ca33"
	PaperLime700        = "afb42b"
	PaperLime800        = "9e9d24"
	PaperLime900        = "827717"
	PaperLimeA100       = "f4ff81"
	PaperLimeA200       = "eeff41"
	PaperLimeA400       = "c6ff00"
	PaperLimeA700       = "aeea00"
	PaperYellow50       = "fffde7"
	PaperYellow100      = "fff9c4"
	PaperYellow200      = "fff59d"
	PaperYellow300      = "fff176"
	PaperYellow400      = "ffee58"
	PaperYellow500      = "ffeb3b"
	PaperYellow600      = "fdd835"
	PaperYellow700      = "fbc02d"
	PaperYellow800      = "f9a825"
	PaperYellow900      = "f57f17"
	PaperYellowA100     = "ffff8d"
	PaperYellowA200     = "ffff00"
	PaperYellowA400     = "ffea00"
	PaperYellowA700     = "ffd600"
	PaperAmber50        = "fff8e1"
	PaperAmber100       = "ffecb3"
	PaperAmber200       = "ffe082"
	PaperAmber300       = "ffd54f"
	PaperAmber400       = "ffca28"
	PaperAmber500       = "ffc107"
	PaperAmber600       = "ffb300"
	PaperAmber700       = "ffa000"
	PaperAmber800       = "ff8f00"
	PaperAmber900       = "ff6f00"
	PaperAmberA100      = "ffe57f"
	PaperAmberA200      = "ffd740"
	PaperAmberA400      = "ffc400"
	PaperAmberA700      = "ffab00"
	PaperBrown50        = "efebe9"
	PaperBrown100       = "d7ccc8"
	PaperBrown200       = "bcaaa4"
	PaperBrown300       = "a1887f"
	PaperBrown400       = "8d6e63"
	PaperBrown500       = "795548"
	PaperBrown600       = "6d4c41"
	PaperBrown700       = "5d4037"
	PaperBrown800       = "4e342e"
	PaperBrown900       = "3e2723"
	PaperOrange50       = "fff3e0"
	PaperOrange100      = "ffe0b2"
	PaperOrange200      = "ffcc80"
	PaperOrange300      = "ffb74d"
	PaperOrange400      = "ffa726"
	PaperOrange500      = "ff9800"
	PaperOrange600      = "fb8c00"
	PaperOrange700      = "f57c00"
	PaperOrange800      = "ef6c00"
	PaperOrange900      = "e65100"
	PaperOrangeA100     = "ffd180"
	PaperOrangeA200     = "ffab40"
	PaperOrangeA400     = "ff9100"
	PaperOrangeA700     = "ff6500"
	PaperDeepOrange50   = "fbe9e7"
	PaperDeepOrange100  = "ffccbc"
	PaperDeepOrange200  = "ffab91"
	PaperDeepOrange300  = "ff8a65"
	PaperDeepOrange400  = "ff7043"
	PaperDeepOrange500  = "ff5722"
	PaperDeepOrange600  = "f4511e"
	PaperDeepOrange700  = "e64a19"
	PaperDeepOrange800  = "d84315"
	PaperDeepOrange900  = "bf360c"
	PaperDeepOrangeA100 = "ff9e80"
	PaperDeepOrangeA200 = "ff6e40"
	PaperDeepOrangeA400 = "ff3d00"
	PaperDeepOrangeA700 = "dd2c00"
	PaperGrey50         = "fafafa"
	PaperGrey100        = "f5f5f5"
	PaperGrey200        = "eeeeee"
	PaperGrey300        = "e0e0e0"
	PaperGrey400        = "bdbdbd"
	PaperGrey500        = "9e9e9e"
	PaperGrey600        = "757575"
	PaperGrey700        = "616161"
	PaperGrey800        = "424242"
	PaperGrey900        = "212121"
	PaperBlueGrey50     = "eceff1"
	PaperBlueGrey100    = "cfd8dc"
	PaperBlueGrey200    = "b0bec5"
	PaperBlueGrey300    = "90a4ae"
	PaperBlueGrey400    = "78909c"
	PaperBlueGrey500    = "607d8b"
	PaperBlueGrey600    = "546e7a"
	PaperBlueGrey700    = "455a64"
	PaperBlueGrey800    = "37474f"
	PaperBlueGrey900    = "263238"

	DarkDividerOpacity    = 0.12
	DarkDisabledOpacity   = 0.38
	DarkSecondaryOpacity  = 0.54
	DarkPrimaryOpacity    = 0.87
	LightDividerOpacity   = 0.12
	LightDisabledOpacity  = 0.3
	LightSecondaryOpacity = 0.7
	LightPrimaryOpacity   = 1.0
)
