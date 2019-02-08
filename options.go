package captcha

import "image/color"

const (
	DefaultWidth                  = 800
	DefaultHeight                 = 200
	DefaultFontSize       float64 = 64
	DefaultCharacterCount         = 7

	defaultArea = DefaultWidth * DefaultHeight
)

type Options struct {
	BackgroundColor color.RGBA
	Width, Height   int
	FontSize        float64
	CharacterCount  int
}

func (o *Options) SetDefaults() {
	if o.Width == 0 {
		o.Width = DefaultWidth
	}

	if o.Height == 0 {
		o.Height = DefaultHeight
	}

	if o.FontSize == 0 {
		o.FontSize = DefaultFontSize
	}

	if o.CharacterCount == 0 {
		o.CharacterCount = DefaultCharacterCount
	}
}
