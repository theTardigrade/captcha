package captcha

import "image/color"

const (
	DefaultWidth                    = 800
	DefaultHeight                   = 200
	DefaultFontSize         float64 = 64
	DefaultCharacterCount           = 7
	DefaultLetterProportion float64 = 0.5

	defaultArea = DefaultWidth * DefaultHeight
)

type BackgroundType uint8

const (
	BackgroundFillType BackgroundType = iota
	BackgroundCirclesType
)

type Options struct {
	BackgroundColor  color.RGBA
	TextColor        color.RGBA
	BackgroundType   BackgroundType
	Width, Height    int
	FontSize         float64
	CharacterCount   int
	UseIdentifier    bool
	UseConcurrency   bool
	LetterProportion float64
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

	if o.LetterProportion == 0 {
		o.LetterProportion = DefaultLetterProportion
	}
}
