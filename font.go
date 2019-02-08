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
)

func loadFont(size float64) (font font.Face, err error) {
	prevFontMutex.RLock()

	if prevFont != nil && prevFontSize == size {
		font = prevFont
		prevFontMutex.RUnlock()
		return
	}
	prevFontMutex.RUnlock()

	font, err = gg.LoadFontFace(fontPath(), size)
	if err == nil {
		prevFontMutex.Lock()
		prevFont = font
		prevFontMutex.Unlock()
	}

	return
}

func fontPath() string {
	return path.Join(goPath(), "src", packagePath(), fontRelativePath)
}

func packagePath() string {
	c := Captcha{}
	return reflect.TypeOf(c).PkgPath()
}

func goPath() string {
	return os.Getenv("GOPATH")
}
