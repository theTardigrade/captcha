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
	fontRelativePath = "assets/CutiveMono-Regular.ttf"
)

var (
	prevFontSize  float64
	prevFont      font.Face
	prevFontMutex sync.RWMutex
)

func loadFont(size float64) (font font.Face, err error) {
	prevFontMutex.RLock()
	defer prevFontMutex.RUnlock()

	if prevFont != nil && size == prevFontSize {
		font = prevFont
		return
	}

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
	s := struct{}{}
	return reflect.TypeOf(s).PkgPath()
}

func goPath() string {
	return os.Getenv("GOPATH")
}
