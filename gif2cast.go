package gif2cast

import (
	"bytes"
	"image"
	"image/gif"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/mkishere/gif2cast/asciicast"
	"github.com/mkishere/gif2cast/imgutil"
)

const (
	endpoint = "https://asciinema.org/api/asciicasts"
)

// Gif2Cast helps with image conversion(resize) and controlling the virtual screen
type Gif2Cast struct {
	oriImg *gif.GIF
	cw, ch int
	rgbImg []*image.RGBA
	title  string
}

// NewGif2Cast creates a new Gif2Cast object
func NewGif2Cast(gif *gif.GIF, w, h int, title string) *Gif2Cast {
	return &Gif2Cast{
		oriImg: gif,
		cw:     w,
		ch:     h,
		title:  title,
	}
}

// Write writes file content to writer
func (gc *Gif2Cast) Write(w io.Writer) (err error) {
	gc.rgbImg, err = imgutil.Resize(gc.oriImg, gc.cw, gc.ch*2)
	if err != nil {
		return err
	}
	vs := asciicast.NewVirtualScreen(gc.cw, gc.ch, gc.title)
	for i := range gc.rgbImg {
		vs.WriteFrame(gc.rgbImg[i], gc.oriImg.Delay[i])
	}

	return vs.Write(w)
}

// Upload the written file to asciinema server
func (gc *Gif2Cast) Upload(apiKey string) (string, error) {
	buf := &bytes.Buffer{}
	writer := multipart.NewWriter(buf)
	filePart, _ := writer.CreateFormFile("asciicast", "ascii.cast")
	_ = gc.Write(filePart)
	writer.Close()
	req, _ := http.NewRequest("POST", endpoint, buf)
	req.SetBasicAuth("gif2cast", apiKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Add("User-Agent", "gif2cast/1.0.0")
	client := &http.Client{}
	rsp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	body := &bytes.Buffer{}
	_, err = body.ReadFrom(rsp.Body)
	rsp.Body.Close()
	return string(body.Bytes()), err
}
