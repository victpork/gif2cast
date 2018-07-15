package asciicast

import (
	"encoding/json"
	"fmt"
	"image"
	"io"
	"os"
	"strings"
	"time"
)

type asciiCastHeader struct {
	Version   int               `json:"version"`
	Width     int               `json:"width"`
	Height    int               `json:"height"`
	Timestamp int64             `json:"timestamp"`
	Command   string            `json:"command"`
	Title     string            `json:"title"`
	Env       map[string]string `json:"env"`
}

type VirtualScreen struct {
	header           asciiCastHeader
	accumulatedDelay float32
	textBuffer       strings.Builder
}

const (
	LogTimeFormat string = "20060102-150405"
)

// NewVirtualScreen creates new virtual screen struct, provide
// it with console size and title that would in the asciicast
// file header
func NewVirtualScreen(w, h int, title string) *VirtualScreen {
	header := asciiCastHeader{
		Version:   2,
		Width:     w,
		Height:    h,
		Timestamp: time.Now().Unix(),
		Command:   "gif2cast",
		Title:     title,
		Env:       map[string]string{"TERM": "gif2cast", "ENV": "gif2cast"},
	}
	return &VirtualScreen{
		header:           header,
		accumulatedDelay: 0,
		textBuffer:       strings.Builder{},
	}
}

func ansiString(frame *image.RGBA) string {
	strBuf := strings.Builder{}
	fmt.Fprintf(&strBuf, "\033[1;1H")
	for y := 0; y < frame.Rect.Dy(); y += 2 {
		for x := 0; x < frame.Rect.Dx(); x++ {
			i := frame.PixOffset(x, y)
			//48 is background, upper part of cell
			if i < 3 || i > 3 && (frame.Pix[i] != frame.Pix[i-3] ||
				frame.Pix[i+1] != frame.Pix[i-2] ||
				frame.Pix[i+2] != frame.Pix[i-1]) {
				fmt.Fprintf(&strBuf, "\033[48;2;%d;%d;%dm", frame.Pix[i], frame.Pix[i+1], frame.Pix[i+2])
			}
			lPix := i + frame.Stride
			//38 is foreground, lower part
			if i < 3 || i > 3 && (frame.Pix[lPix] != frame.Pix[lPix-3] ||
				frame.Pix[lPix+1] != frame.Pix[lPix-2] ||
				frame.Pix[lPix+2] != frame.Pix[lPix-1]) {
				fmt.Fprintf(&strBuf, "\033[38;2;%d;%d;%dm",
					frame.Pix[i+frame.Stride], frame.Pix[i+frame.Stride+1], frame.Pix[i+frame.Stride+2])
			}
			//u2584 is the fill lowerblock char
			fmt.Fprint(&strBuf, "\u2584")
		}
		fmt.Fprint(&strBuf, "\r\n")
	}
	return strBuf.String()
}

// WriteFrame write a image.RGBA into the string builder in asciicast format
func (vs *VirtualScreen) WriteFrame(frame *image.RGBA, delay int) {
	escStr, _ := json.Marshal(ansiString(frame))
	fmt.Fprintf(&vs.textBuffer, "[%f, \"o\", %s]\r\n", vs.accumulatedDelay, escStr)
	vs.accumulatedDelay += float32(delay) / 100
}

// Write writes the asciicast file in buffer to the file
func (vs *VirtualScreen) WriteToFile(filename string) error {
	fp, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer fp.Close()
	return vs.Write(fp)
}

// Write writes the asciicast file in buffer to the writer
func (vs *VirtualScreen) Write(w io.Writer) error {
	bheader, err := json.Marshal(vs.header)
	if err != nil {
		return err
	}
	_, err = w.Write(bheader)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte("\r\n"))
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(vs.textBuffer.String()))
	if err != nil {
		return err
	}
	return nil
}
