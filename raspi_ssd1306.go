package main

import (
	"fmt"
	"image"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	_ "github.com/hnw/gobot-ssd1306-example/statik" // register embedded data to statik/fs on init func

	"github.com/golang/freetype/truetype"
	"github.com/rakyll/statik/fs"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

func main() {
	width := 128
	height := 32
	r := raspi.NewAdaptor()
	oled := i2c.NewSSD1306Driver(r, i2c.WithSSD1306DisplayWidth(width), i2c.WithSSD1306DisplayHeight(height))

	face, err := newFace()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	work := func() {
		i := 0
		gobot.Every(1*time.Second, func() {
			oled.Clear()
			img := textToImage(strconv.Itoa(i), face, width, height)
			oled.ShowImage(img)
			oled.Display()
			i++
		})
	}

	robot := gobot.NewRobot("ssd1306Robot",
		[]gobot.Connection{r},
		[]gobot.Device{oled},
		work,
	)

	robot.Start()
}

func newFace() (face font.Face, err error) {
	statikFs, err := fs.New()
	if err != nil {
		return
	}
	f, err := statikFs.Open("/Koruri-Regular.ttf")
	if err != nil {
		return
	}
	defer f.Close()

	ftBinary, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}
	ft, err := truetype.Parse(ftBinary)
	if err != nil {
		return
	}
	opt := truetype.Options{
		Size:              24,
		DPI:               0,
		Hinting:           0,
		GlyphCacheEntries: 0,
		SubPixelsX:        0,
		SubPixelsY:        0,
	}
	face = truetype.NewFace(ft, &opt)
	return
}

func textToImage(text string, face font.Face, width int, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	dr := &font.Drawer{
		Dst:  img,
		Src:  image.White,
		Face: face,
		Dot:  fixed.Point26_6{},
	}

	dr.Dot.X = (fixed.I(width) - dr.MeasureString(text)) / 2
	dr.Dot.Y = fixed.I(24)
	dr.DrawString(text)

	return img
}
