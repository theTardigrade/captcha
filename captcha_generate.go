package captcha

import (
	"bytes"
	"encoding/base64"
	"html/template"
	"strconv"
	"strings"
	"time"

	"gopkg.in/fogleman/gg.v1"
)

var (
	letters = [...]byte{
		'B', 'C', 'D', 'F', 'G', 'H', 'J', 'K', 'L', 'M',
		'N', 'P', 'Q', 'R', 'S', 'T', 'W', 'X', 'Y', 'Z',
	}
	numbers = [...]byte{
		'1', '2', '3', '4', '5', '6', '7', '8', '9',
	}
)

func (c *Captcha) generateImage(opts *Options) error {
	width := float64(opts.Width)
	height := float64(opts.Height)
	area := width * height
	halfWidth := width / 2
	halfHeight := height / 2
	backgroundColor := opts.BackgroundColor
	textColor := opts.TextColor
	fontSize := opts.FontSize
	characterCount := opts.CharacterCount

	dc := gg.NewContext(opts.Width, opts.Height)

	r, g, b := float64(backgroundColor.R)/255, float64(backgroundColor.G)/255, float64(backgroundColor.B)/255

	switch opts.BackgroundType {
	case BackgroundCirclesType:
		{
			dc.SetRGB(1, 1, 1)
			dc.Clear()

			for x, l := float64(0)-width*0.125, width*1.25; x < l; x += float64(c.random.Intn(int(width/5)+1)) + width/80 {
				alpha := 255 - c.random.Intn(129)
				dc.SetRGBA255(int(backgroundColor.R), int(backgroundColor.G), int(backgroundColor.B), alpha)
				r := float64(c.random.Intn(int(area/1e3)+1)) + area/600
				y := (float64(c.random.Intn(21)-10)*DefaultHeight)/height + halfHeight
				dc.DrawCircle(x, y, r)
				dc.Fill()
			}
		}
	case BackgroundFillType:
		{
			dc.SetRGB(r, g, b)
			dc.Clear()
		}
	}

	font, err := loadFont(fontSize)
	if err != nil {
		return err
	}
	dc.SetFontFace(font)

	var builder strings.Builder

	for i := 0; i < characterCount; i++ {
		var b byte

		if f := c.random.Float64(); f <= opts.LetterProportion {
			i := c.random.Intn(len(letters))
			b = letters[i]
		} else {
			i := c.random.Intn(len(numbers))
			b = numbers[i]
		}

		builder.WriteByte(b)

		alpha := 255 - c.random.Intn(65)
		dc.SetRGBA255(int(textColor.R), int(textColor.G), int(textColor.B), alpha)

		s := string(b)
		w, h := dc.MeasureString(s)
		x := width/float64(characterCount)*(float64(i)+0.5) - w/2
		y := halfHeight + h/4
		r := float64(c.random.Intn(65)-32) / 384

		dc.RotateAbout(r, halfWidth, halfHeight)
		dc.DrawString(s, x, y)
		dc.RotateAbout(-r, halfWidth, halfHeight)
	}

	c.Value = builder.String()

	buffer := bytes.NewBuffer(nil)
	err = dc.EncodePNG(buffer)
	if err != nil {
		return err
	}

	builder.Reset()

	builder.WriteString("data:image/png;base64,")
	builder.WriteString(base64.StdEncoding.EncodeToString(buffer.Bytes()))

	c.ImageURL = template.URL(builder.String())

	return nil
}

const (
	identifierSegmentMaxLength = 13
	identifierSegmentCount     = 7
	identifierSeparatorByte    = '-'
	identifierMaxLength        = identifierSegmentMaxLength*identifierSegmentCount + (identifierSegmentCount - 1)
)

func (c *Captcha) generateIdentifier() {
	var builder strings.Builder

	builder.Grow(identifierMaxLength)

	l := identifierSegmentCount / 2

	for i := 0; i < l; i++ {
		builder.WriteString(strconv.FormatInt(c.random.Int63(), 36))
		builder.WriteByte(identifierSeparatorByte)
	}

	builder.WriteString(strconv.FormatInt(int64(time.Now().UTC().UnixNano()), 36))

	if identifierSegmentCount%2 == 0 {
		l--
	}

	for i := 0; i < l; i++ {
		builder.WriteByte(identifierSeparatorByte)
		builder.WriteString(strconv.FormatInt(c.random.Int63(), 36))
	}

	c.Identifier = builder.String()
}
