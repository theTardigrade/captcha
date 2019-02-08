package captcha

import (
	"bytes"
	"encoding/base64"
	"image/color"
	"math/rand"
	"strconv"
	"time"

	"gopkg.in/fogleman/gg.v1"
)

const (
	characterCount = 7
	imageWidth     = 800
	imageHeight    = 200
)

type Options struct {
	BackgroundColor color.RGBA
}

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

	dc := gg.NewContext(imageWidth, imageHeight)
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	r, g, b := float64(opts.backgroundColor.R)/255, float64(opts.backgroundColor.G)/255, float64(opts.backgroundColor.B)/255

	for x := float64(0); x < imageWidth; x += float64(rand.Intn(81) + 16) {
		a := float64(rand.Intn(49)+16) / 64
		dc.SetRGBA(r, g, b, a)
		r := float64(rand.Intn(41) + 60)
		y := float64(rand.Intn(21)-10) + imageHeight/2
		dc.DrawCircle(x, y, r)
		dc.Fill()
	}

	font, err := gg.LoadFontFace("assets/CutiveMono-Regular.ttf", 128)
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

		c.value += s

		w, h := dc.MeasureString(s)
		a := float64(rand.Intn(65)-32) / 384
		dc.RotateAbout(a, imageWidth/2, imageHeight/2)
		dc.DrawString(s, float64(imageWidth)/float64(characterCount)*(float64(i)+0.5)-w/2, float64(imageHeight)/2+h/4)
		dc.RotateAbout(-a, imageWidth/2, imageHeight/2)
	}

	buffer := bytes.NewBuffer(nil)
	err = dc.EncodePNG(buffer)
	if err != nil {
		return nil, err
	}

	image := base64.StdEncoding.EncodeToString(buffer.Bytes())
	c.image = "data:image/png;base64," + image

	buffer.Reset()
	for i := 0; i < 11; i++ {
		buffer.WriteString(strconv.FormatInt(rand.Int63(), 36))
		buffer.WriteByte('-')
	}
	buffer.WriteString(strconv.FormatInt(int64(time.Now().UTC().UnixNano()), 36))
	c.identifier = buffer.String()

	return &c, nil
}
