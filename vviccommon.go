package vviccommon

import (
	"bytes"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"regexp"
	"strings"

	"github.com/nfnt/resize"
)

var tidyWordsPattern = regexp.MustCompile(strings.Join([]string{
	"[0-9]+",
	"#",
	"[0-9]+年",
	"春秋季",
	"真丝",
	"-",
	"春夏",
	"夏装",
	"春装",
	"春季",
	"夏季",
	`\s+`,
	"[0-9]+斤",
	"特价",
	"不退现",
	"新款",
	`\(.+\)`,
	`【.+】`,
	`（.+）`,
	"实拍",
	"代发",
	"模特",
	"新品",
	"官网",
	"超模",
	"新品",
	"官方图",
	"现货",
}, "|"))

func TidyTitle(in string) string {
	return tidyWordsPattern.ReplaceAllString(in, "")
}

var logoImage, WatermarkImage image.Image

func init() {
	var err error
	logoBytes, _ := logoNolightPngBytes()
	logoImage, _, err = image.Decode(bytes.NewReader(logoBytes))
	if err != nil {
		panic("decode logo image")
	}
	watermarkBytes, _ := watermarkPngBytes()
	WatermarkImage, _, err = image.Decode(bytes.NewReader(watermarkBytes))
	if err != nil {
		panic("decode watermark image")
	}
}

func CompositeLogo(r io.Reader, w io.Writer) (err error) {
	defer ct(&err)
	// try to resize and add logo
	img, what, err := image.Decode(r)
	ce(err, "decode image")
	// resize
	img = resize.Resize(800, 800, img, resize.Bicubic)
	// composite
	dst := image.NewRGBA(image.Rect(0, 0, 800, 800))
	draw.Draw(dst, img.Bounds(), img, image.Pt(0, 0), draw.Over)
	draw.Draw(dst, image.Rect(90, 0, logoImage.Bounds().Max.X+90, logoImage.Bounds().Max.Y),
		logoImage, image.Pt(0, 0), draw.Over)
	switch what {
	case "jpeg":
		ce(jpeg.Encode(w, dst, &jpeg.Options{
			Quality: 90,
		}), "encode image")
	case "png":
		ce(png.Encode(w, dst), "encode image")
	case "gif":
		ce(gif.Encode(w, dst, &gif.Options{
			NumColors: 256,
		}), "encode image")
	default:
		panic("image file not supported")
	}
	return
}

func CompositeWatermark(r io.Reader, w io.Writer) (err error) {
	defer ct(&err)
	img, what, err := image.Decode(r)
	ce(err, "decode image")
	dst := image.NewRGBA(img.Bounds())
	imageRect := img.Bounds()
	draw.Draw(dst, imageRect, img, image.Pt(0, 0), draw.Over)
	watermarkRect := WatermarkImage.Bounds()
	draw.Draw(dst, image.Rect(img.Bounds().Max.X-watermarkRect.Max.X, 0,
		imageRect.Max.X, watermarkRect.Max.Y), WatermarkImage,
		image.Pt(0, 0), draw.Over)
	switch what {
	case "jpeg":
		ce(jpeg.Encode(w, dst, &jpeg.Options{
			Quality: 90,
		}), "encode image")
	case "png":
		ce(png.Encode(w, dst), "encode image")
	case "gif":
		ce(gif.Encode(w, dst, &gif.Options{
			NumColors: 256,
		}), "encode image")
	default:
		panic("image file not supported")
	}
	return
}
