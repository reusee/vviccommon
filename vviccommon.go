package vviccommon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"net/http"
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
	`\[.+\]`,
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

func ShuffleTitle(in string) (out string, err error) {
	defer ct(&err)
	var data [][][]struct {
		Id        int
		Cont      string
		Pos       string
		Ne        string
		Parent    int
		Relate    string
		Semparent int
		Semrelate string
		Arg       []interface{}
	}
	reqPath := fmt.Sprintf("http://api.ltp-cloud.com/analysis/?api_key=f6m2I6r8pnGqqgZdbldlosMWRQNgxwcBzX7k1Tui&text=%s&pattern=all&format=json",
		in)
	resp, err := http.Get(reqPath)
	ce(err, "get %s", reqPath)
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	ce(err, "read body")
	err = json.Unmarshal(content, &data)
	ce(err, "decode")
	fmt.Printf("%v\n", data)
	return
}

var logoImage, WatermarkImage image.Image

func init() {
	var err error
	logoBytes, err := ioutil.ReadFile("./logo.png")
	if err != nil {
		panic("read logo file")
	}
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
	img = resize.Resize(800, 800, img, resize.Lanczos3)
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

func ScaleTo800x800(r io.Reader, w io.Writer) (err error) {
	defer ct(&err)
	// try to resize and add logo
	img, what, err := image.Decode(r)
	ce(err, "decode image")
	// resize
	img = resize.Resize(800, 800, img, resize.Lanczos3)
	// composite
	switch what {
	case "jpeg":
		ce(jpeg.Encode(w, img, &jpeg.Options{
			Quality: 90,
		}), "encode image")
	case "png":
		ce(png.Encode(w, img), "encode image")
	case "gif":
		ce(gif.Encode(w, img, &gif.Options{
			NumColors: 256,
		}), "encode image")
	default:
		panic("image file not supported")
	}
	return
}

func ScaleForMobile(r io.Reader, w io.Writer) (err error) {
	defer ct(&err)
	// decode
	img, what, err := image.Decode(r)
	ce(err, "decode image")
	// resize
	img = resize.Resize(600, 0, img, resize.Lanczos3)
	switch what {
	case "jpeg":
		ce(jpeg.Encode(w, img, &jpeg.Options{
			Quality: 75,
		}), "encode image")
	case "png":
		ce(png.Encode(w, img), "encode image")
	case "gif":
		ce(gif.Encode(w, img, &gif.Options{
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

func ScaleImageToJpeg(width uint, quality int, r io.Reader, w io.Writer) (err error) {
	defer ct(&err)
	// decode
	img, _, err := image.Decode(r)
	ce(err, "decode image")
	// resize
	img = resize.Resize(width, 0, img, resize.Lanczos3)
	// encode
	ce(jpeg.Encode(w, img, &jpeg.Options{
		Quality: quality,
	}), "encode image")
	return
}
