package captcha

import (
	"bytes"
	"encoding/base64"
	"math/rand"
	"strconv"
	"strings"
	"sync"
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

	var waitGroup sync.WaitGroup

	if opts.UseIdentifier {
		waitGroup.Add(1)
		go func() {
			c.generateIdentifier()
			waitGroup.Done()
		}()
	}

	errChan := make(chan error)

	go func(errChan chan<- error, opts *Options) {
		if err := c.generateImage(opts); err != nil {
			errChan <- err
		}
		waitGroup.Done()
	}(errChan, &opts)

	waitGroup.Wait()

	select {
	case err := <-errChan:
		return nil, err
	default:
		return &c, nil
	}
}

func (c *Captcha) CheckValue(value string) bool {
	return strings.ToUpper(value) == c.Value
}

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

	buffer := bytes.NewBuffer(nil)

	dc := gg.NewContext(opts.Width, opts.Height)

	r, g, b := float64(backgroundColor.R)/255, float64(backgroundColor.G)/255, float64(backgroundColor.B)/255

	switch opts.BackgroundType {
	case BackgroundCirclesType:
		{
			dc.SetRGB(1, 1, 1)
			dc.Clear()

			for x := float64(0); x < width; x += float64(rand.Intn(int(width/5)+1)) + width/80 {
				alpha := 255 - rand.Intn(129)
				dc.SetRGBA255(int(backgroundColor.R), int(backgroundColor.G), int(backgroundColor.B), alpha)
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
		return err
	}
	dc.SetFontFace(font)

	for i := 0; i < characterCount; i++ {
		var s string

		if rand.Float64() < float64(1)/3 {
			s = string('A' + rand.Intn('H'-'A'+1))
		} else {
			s = string('1' + rand.Intn('9'-'1'+1))
		}

		c.Value += s

		alpha := 255 - rand.Intn(65)
		dc.SetRGBA255(int(textColor.R), int(textColor.G), int(textColor.B), alpha)

		w, h := dc.MeasureString(s)
		r := float64(rand.Intn(65)-32) / 384
		dc.RotateAbout(r, halfWidth, halfHeight)
		dc.DrawString(s, width/float64(characterCount)*(float64(i)+0.5)-w/2, halfHeight+h/4)
		dc.RotateAbout(-r, halfWidth, halfHeight)
	}

	err = dc.EncodePNG(buffer)
	if err != nil {
		return err
	}

	{
		imageBuffer := bytes.NewBuffer(nil)
		imageBuffer.WriteString("data:image/png;base64,")
		imageBuffer.WriteString(base64.StdEncoding.EncodeToString(buffer.Bytes()))
		c.Image = imageBuffer.String()
	}

	return nil
}

func (c *Captcha) generateIdentifier() {
	buffer := bytes.NewBuffer(nil)

	for i := 0; i < 4; i++ {
		buffer.WriteString(strconv.FormatInt(rand.Int63(), 36))
		buffer.WriteByte('-')
	}

	buffer.WriteString(strconv.FormatInt(int64(time.Now().UTC().UnixNano()), 36))

	for i := 0; i < 5; i++ {
		buffer.WriteByte('-')
		buffer.WriteString(strconv.FormatInt(rand.Int63(), 36))
	}

	c.Identifier = buffer.String()
}

func CheckValues(expectedValue, receivedValue string) bool {
	c := Captcha{Value: expectedValue}
	return c.CheckValue(receivedValue)
}
