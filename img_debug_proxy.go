package main

import (
    "net/http"
    "github.com/elazarl/goproxy"
    "github.com/elazarl/goproxy/ext/image"
    "fmt"
    "log"
    "flag"
    "io/ioutil"
    "image"
    "golang.org/x/image/draw"
    "code.google.com/p/freetype-go/freetype"
    "code.google.com/p/freetype-go/freetype/truetype"
    "code.google.com/p/freetype-go/freetype/raster"
)

var (
    dpi      = flag.Float64("dpi", 72, "screen resolution in Dots Per Inch")
    fontFile = flag.String("fontfile", "Lato-Regular.ttf", "filename of the ttf font")
    hinting  = flag.String("hinting", "none", "none | full")
    size     = flag.Float64("size", 12, "font size in points")
    verbose  = flag.Bool("verbose", false, "proxy verbose mode, boolean")
    port     = flag.String("port", "8080", "proxy port number")
    indent   = flag.Int("indent", 5, "label indentation from the top-left image's corner in pixels")
)


func main() {
    flag.Parse()

    font, err := loadFont()
    if err != nil {
        log.Println(err)
        return
    }

    fontHeight := int(createContext(font).PointToFix32(*size)>>8)
    // two points to output the text and its shadow
    ptA := freetype.Pt(*indent + 1, *indent + 1 + fontHeight)
    ptB := freetype.Pt(*indent, *indent + fontHeight)

    proxy := goproxy.NewProxyHttpServer()
    proxy.OnResponse().Do(goproxy_image.HandleImage(func(img image.Image, ctx *goproxy.ProxyCtx) image.Image {
        outImage := image.NewRGBA(img.Bounds())
        draw.Copy(outImage, image.ZP, img, img.Bounds(), nil)
        text := fmt.Sprintf("%dx%d", img.Bounds().Dx(), img.Bounds().Dy())
        fontContext := createContext(font)
        fontContext.SetClip(img.Bounds())
        fontContext.SetDst(outImage)

        drawString(image.White, fontContext, ptA, text)
        drawString(image.Black, fontContext, ptB, text)

        return outImage
    }))
    proxy.Verbose = *verbose
    log.Fatal(http.ListenAndServe(":" + *port, proxy))
}

func loadFont() (*truetype.Font, error){
    fontBytes, err := ioutil.ReadFile(*fontFile)
    if err != nil {
        log.Println(err)
        return nil, err
    }
    font, err := freetype.ParseFont(fontBytes)
    if err != nil {
        log.Println(err)
        return nil, err
    }
    return font, nil
}

func drawString(src image.Image, fontContext *freetype.Context, position raster.Point, text string){
    fontContext.SetSrc(src)
    _, err := fontContext.DrawString(text, position)
    if err != nil {
        log.Println(err)
    }
}

func createContext(font *truetype.Font) *freetype.Context {
    ctx := freetype.NewContext()
    ctx.SetDPI(*dpi)
    ctx.SetFont(font)
    ctx.SetFontSize(*size)
    switch *hinting {
        default:
        ctx.SetHinting(freetype.NoHinting)
        case "full":
        ctx.SetHinting(freetype.FullHinting)
    }
    return ctx
}