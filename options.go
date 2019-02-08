package captcha

import "image/color"

const (
	DefaultWidth                       = 800
	DefaultHeight                      = 200
	DefaultFontSize            float64 = 64
	DefaultCharacterCount              = 7
	DefaultLetterToNumberRatio float64 = 0.7

	defaultArea = DefaultWidth * DefaultHeight
)

type backgroundType uint8

const (
	BackgroundFillType backgroundType = iota
	BackgroundCirclesType
)

type Options struct {
	BackgroundColor     color.RGBA
	TextColor           color.RGBA
	BackgroundType      backgroundType
	Width, Height       int
	FontSize            float64
	CharacterCount      int
	UseIdentifier       bool
	LetterToNumberRatio float64
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

	if o.LetterToNumberRatio == 0 {
		o.LetterToNumberRatio = DefaultLetterToNumberRatio
	}
}
