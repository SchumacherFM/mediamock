package record

import (
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/SchumacherFM/mediamock/common"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/ugorji/go/codec"

	"github.com/golang/freetype/truetype"
)

const DirPerm os.FileMode = 0755
const CSVSep = "|"
const CSVNewLine = "\n"
const ps = string(os.PathSeparator)

var codecJSON codec.Handle = new(codec.JsonHandle)

type Record struct {
	Path    string    // idx 0
	ModTime time.Time // idx 1
	Width   int       // idx 2
	Height  int       // idx 3

	pattern       string
	allowDrawText bool
	ext           string
}

func NewRecord(pattern string, csv ...string) (Record, error) {
	if len(csv) != 4 {
		return Record{}, fmt.Errorf("Incorrect number %d of %d of CSV columns: %#v", len(csv), 4, csv)
	}

	w, _ := strconv.Atoi(csv[2])
	h, _ := strconv.Atoi(csv[3])
	t, _ := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", csv[1])

	isText := strings.Contains(pattern, "text")
	if isText {
		pattern = strings.Replace(pattern, "text", "", -1)
	}

	return Record{
		Path:          csv[0],
		ModTime:       t,
		Width:         w,
		Height:        h,
		pattern:       pattern,
		allowDrawText: isText,
		ext:           filepath.Ext(csv[0]),
	}, nil
}

func NewRecordFields(pattern, path string, width, height int) Record {
	isText := strings.Contains(pattern, "text")
	if isText {
		pattern = strings.Replace(pattern, "text", "", -1)
	}
	return Record{
		Path:          path,
		ModTime:       time.Now(),
		Width:         width,
		Height:        height,
		pattern:       pattern,
		allowDrawText: isText,
		ext:           filepath.Ext(path),
	}
}

func (r Record) FileExt() string {
	return r.ext
}

func (r Record) ToJSON(w io.Writer) error {
	var enc *codec.Encoder = codec.NewEncoder(w, codecJSON)
	return enc.Encode(r)
}

func (r Record) Write(w io.Writer) error {
	_, err := fmt.Fprintf(w, "%s%s%s%s%d%s%d%s", r.Path, CSVSep, r.ModTime, CSVSep, r.Width, CSVSep, r.Height, CSVNewLine)
	return err
}

func (r Record) CreateFile(basePath string) error {

	d, f := r.getDirFile(basePath)
	if err := os.MkdirAll(d, DirPerm); err != nil {
		common.UsageAndExit("Failed to create directory: %s with error: %s", d, err)
	}

	file, err := os.Create(d + f)
	if err != nil {
		return fmt.Errorf("Failed to create file %s%s", d, f)
	}
	defer func() {
		if err := file.Close(); err != nil {
			common.InfoErr("Failed to close file %s with error: %s\n", d+f, err)
		}
	}()

	if err := os.Chtimes(d+f, r.ModTime, r.ModTime); err != nil {
		common.InfoErr("Failed to change time for file %s with error: %s\n", d+f, err)
	}
	r.CreateContent(f, file)

	return nil
}

func (r Record) CreateContent(f string, w io.Writer) {
	switch r.ext {
	case ".png":
		if err := png.Encode(w, r.generateImage()); err != nil {
			common.InfoErr("Failed to create PNG file %s with error: %s\n", f, err)
		}
	case ".jpg", ".jpeg":
		// big file size? reason why is here: https://www.reddit.com/r/golang/comments/3kn1zp/filesize_of_jpegencode/
		if err := jpeg.Encode(w, r.generateImage(), &jpeg.Options{Quality: 75}); err != nil {
			common.InfoErr("Failed to create JPEG file %s with error: %s\n", f, err)
		}
	case ".gif", ".ico":
		if err := gif.Encode(w, r.generateImage(), nil); err != nil {
			common.InfoErr("Failed to create GIF file %s with error: %s\n", f, err)
		}
	default:
		if _, err := w.Write(nil); err != nil {
			common.InfoErr("Failed to write file %s with error: %s\n", f, err)
		}
	}
}

func (r Record) generateImage() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, r.Width, r.Height))

	var src image.Image
	switch {
	case "warm" == r.pattern:
		src = &image.Uniform{colorful.WarmColor()}
	case "happy" == r.pattern:
		src = &image.Uniform{colorful.HappyColor()}
	case "rand" == r.pattern:
		src = &image.Uniform{colorful.LinearRgb(rand.Float64(), rand.Float64(), rand.Float64())}
	case r.isHex():
		hc, _ := colorful.Hex(r.pattern)
		src = &image.Uniform{hc}
	default:
		src = &image.Uniform{colorful.FastWarmColor()}
	}
	draw.Draw(img, img.Bounds(), src, image.ZP, draw.Src)

	if r.allowDrawText {
		_, f := r.getDirFile(r.Path)
		gc := draw2dimg.NewGraphicContext(img)
		drawText(gc, f, 2, float64(r.Height))
	}
	return img
}

func (r Record) getDirFile(base string) (dir, file string) {
	if false == strings.HasSuffix(base, ps) && false == strings.HasPrefix(r.Path, ps) {
		base = base + ps
	}
	dir, file = filepath.Split(base + r.Path)
	return
}

func (r Record) isHex() bool {
	_, err := colorful.Hex(r.pattern)
	return err == nil
}

func init() {
	luximrTTF, err := truetype.Parse(luximr)
	if err != nil {
		common.UsageAndExit("failed to parse luximir font: %s", err)
	}
	luximr = nil // kill 72Kb of font data

	draw2d.RegisterFont(
		draw2d.FontData{Name: "luxi", Family: draw2d.FontFamilyMono, Style: draw2d.FontStyleNormal},
		luximrTTF,
	)
}

func drawText(gc draw2d.GraphicContext, text string, x, y float64) {
	var fontSize float64 = 14
	fontSizeHeight := fontSize + 14
	if y < fontSizeHeight {
		return
	}
	gc.SetFontData(draw2d.FontData{Name: "luxi", Family: draw2d.FontFamilyMono, Style: draw2d.FontStyleNormal})
	// Set the fill text color to black

	gc.SetFillColor(image.Black)
	gc.SetFontSize(fontSize)
	// 9px width each letter and 2px letter spacing at font-size 14

	var yPos float64
	for ; yPos < y; yPos = yPos + fontSizeHeight {
		gc.FillStringAt(text, x, yPos)
	}
}
