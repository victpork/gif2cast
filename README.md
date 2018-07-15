# gif2cast - Convert animated gif into asciicast session

## What it does
Basically does the opposite of [asciicast2gif](https://github.com/asciinema/asciicast2gif) - turning an animated gif into asciicast file, which also uploading to the [asciinema](http://asciinema.org) site.

## Preview
[![asciicast](https://asciinema.org/a/191865.png)](https://asciinema.org/a/191865)

[Original file](https://en.wikipedia.org/wiki/GIF#/media/File:Rotating_earth_(large).gif)

## Usage
### Bundled tool
```
> AC_APIKEY=abcdef-123445 gif2cast animated.gif
or
> gif2cast -o out.cast animated.gif
```
### Library
```go
fp, err := os.Open(fileName)
...
gifImage, err := gif.DecodeAll(fp)
...
width := 80
height := 24
gc := gif2cast.NewGif2Cast(gifImage, 80, 24, "title")
out, err := os.OpenFile(outFilename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
gc.Write(out)
//or
url, err := gc.Upload(APIKEY)
if err != nil {
    fmt.Println("Cannot upload:", err)
    os.Exit(1)
}
fmt.Println(url)
```

## Note

This library is some hobby project I wrote over the weekend so it has a lot of rough edges, e.g.
- No support on partial frame refresh, it takes every frame in the gif as a full refresh.
- Does not support disposal methods
- Does not support loop
- Bundled tool does not find API automatically
- Rudimentary optimization to reduce asciicast file size