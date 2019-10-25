package captcha

import (
	"os"
	"path"
	"reflect"
	"sync"

	"golang.org/x/image/font"

	"gopkg.in/fogleman/gg.v1"
)

const (
	fontFilename     = "CutiveMono-Regular.ttf"
	fontRelativePath = "assets/" + fontFilename
)

var (
	prevFontSize  float64
	prevFont      font.Face
	prevFontMutex sync.RWMutex

	fontPath string
)

func init() {
	fontPath = generateFontPath()
}

func loadFont(size float64) (font font.Face, err error) {
	prevFontMutex.RLock()
	if prevFont != nil && prevFontSize == size {
		font = prevFont
	}
	prevFontMutex.RUnlock()

	if font == nil {
		font, err = gg.LoadFontFace(fontPath, size)
		if err == nil {
			prevFontMutex.Lock()
			prevFont = font
			prevFontMutex.Unlock()
		}
	}

	return
}

func generateFontPath() string {
	return path.Join(goPath(), "src", packagePath(), fontRelativePath)
}

func packagePath() string {
	c := Captcha{}
	return reflect.TypeOf(c).PkgPath()
}

func goPath() string {
	value := os.Getenv("GOPATH")

	if value == "" {
		bin := os.Getenv("GOBIN")

		if bin != "" {
			value = path.Join(bin, "..")
		}
	}

	return value
}
