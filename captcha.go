package captcha

import (
	"bytes"
	"encoding/base64"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"gopkg.in/fogleman/gg.v1"
)

type Captcha struct {
	Image      string
	Value      string
	Identifier string
}

func init() {
	rand.Seed(int64(time.Now().UTC().UnixNano()))
}

func New(opts Options) (*Captcha, error) {
	c := Captcha{}

	opts.SetDefaults()

	width := float64(opts.Width)
	height := float64(opts.Height)
	area := width * height
	halfWidth := width / 2
	halfHeight := height / 2
	backgroundColor := opts.BackgroundColor
	fontSize := opts.FontSize
	characterCount := opts.CharacterCount

	dc := gg.NewContext(opts.Width, opts.Height)

	r, g, b := float64(backgroundColor.R)/255, float64(backgroundColor.G)/255, float64(backgroundColor.B)/255

	switch opts.BackgroundType {
	case BackgroundCirclesType:
		{
			dc.SetRGB(1, 1, 1)
			dc.Clear()

			for x := float64(0); x < width; x += float64(rand.Intn(int(width/5)+1)) + width/80 {
				a := float64(rand.Intn(49)+16) / 64
				dc.SetRGBA(r, g, b, a)
				r := float64(rand.Intn(int(area/1e3)+1)) + area/600
				y := (float64(rand.Intn(21)-10)*DefaultHeight)/height + halfHeight
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
		return nil, err
	}
	dc.SetFontFace(font)
	dc.SetRGBA(1, 1, 1, 1)

	for i := 0; i < characterCount; i++ {
		var s string

		if rand.Float64() < float64(1)/3 {
			s = string('A' + rand.Intn('H'-'A'+1))
		} else {
			s = string('1' + rand.Intn('9'-'1'+1))
		}

		c.Value += s

		w, h := dc.MeasureString(s)
		a := float64(rand.Intn(65)-32) / 384
		dc.RotateAbout(a, halfWidth, halfHeight)
		dc.DrawString(s, width/float64(characterCount)*(float64(i)+0.5)-w/2, halfHeight+h/4)
		dc.RotateAbout(-a, halfWidth, halfHeight)
	}

	buffer := bytes.NewBuffer(nil)
	err = dc.EncodePNG(buffer)
	if err != nil {
		return nil, err
	}

	{
		imageBuffer := bytes.NewBuffer(nil)
		imageBuffer.WriteString("data:image/png;base64,")
		imageBuffer.WriteString(base64.StdEncoding.EncodeToString(buffer.Bytes()))
		c.Image = imageBuffer.String()
	}

	buffer.Reset()
	for i := 0; i < 11; i++ {
		buffer.WriteString(strconv.FormatInt(rand.Int63(), 36))
		buffer.WriteByte('-')
	}
	buffer.WriteString(strconv.FormatInt(int64(time.Now().UTC().UnixNano()), 36))
	for i := 0; i < 4; i++ {
		buffer.WriteByte('-')
		buffer.WriteString(strconv.FormatInt(rand.Int63(), 36))
	}
	c.Identifier = buffer.String()

	return &c, nil
}

func (c *Captcha) CheckValue(value string) bool {
	return strings.ToUpper(value) == c.Value
}

func CheckValues(expectedValue, receivedValue string) bool {
	c := Captcha{Value: expectedValue}
	return c.CheckValue(receivedValue)
}
