package main

import (
	"flag"
	"fmt"
	"image/gif"
	"os"

	"github.com/mkishere/gif2cast"
)

func main() {
	outFilename := ""
	flag.StringVar(&outFilename, "o", "", "Output file name")
	flag.Parse()

	if flag.Arg(0) == "" {
		fmt.Println("Input file expected")
		os.Exit(1)
	}
	inFilename := flag.Arg(0)
	fp, err := os.Open(inFilename)
	if err != nil {
		fmt.Println("Cannot open file:", err)
		os.Exit(1)
	}
	gifImage, err := gif.DecodeAll(fp)
	if err != nil {
		fmt.Println("Cannot open file:", err)
		os.Exit(1)
	}
	gc := gif2cast.NewGif2Cast(gifImage, 80, 24, flag.Arg(0))
	if outFilename != "" {
		out, err := os.OpenFile(outFilename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
		if err != nil {
			fmt.Println("Cannot create output file:", err)
			os.Exit(1)
		}
		defer out.Close()
		gc.Write(out)
	} else {
		//TODO: Get APIkey
		str, err := gc.Upload(os.Getenv("AC_APIKEY"))
		if err != nil {
			fmt.Println("Cannot upload:", err)
			os.Exit(1)
		}
		fmt.Println(str)
	}
}
