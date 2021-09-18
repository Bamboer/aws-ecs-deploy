package main

import (
    "os"
    "bytes"
    "fmt"
    "image"
    "image/draw"
    "image/color"
    "image/gif"
    "image/jpeg"
    "image/png"
    "io/ioutil"
    "net/http"
    "github.com/golang/freetype"
    "github.com/golang/freetype/truetype"
    "golang.org/x/image/font"
)


type DrawTextInfo struct {
    Text string 
    X   int 
    Y   int 
}


type TextBrush struct {
    FontType      *truetype.Font  
    FontSize      float64 
    FontColor     *image.Uniform 
    TextWidth     int 
}

func DrawImage()error{
    const width, height = 150, 150

    // Create a colored image of the given width and height.
    img := image.NewNRGBA(image.Rect(0, 0, width, height))

    for y := 0; y < height; y++ {
        for x := 0; x < width; x++ {
            img.Set(x, y, color.NRGBA{
                R: uint8((x + y) & 255),
                G: uint8((x + y) << 1 & 255),
                B: uint8((x + y) << 2 & 255),
                A: 255,
            })
        }
    }

    f, err := os.Create("image.png")
    if err != nil {
        log.Fatal(err)
    }

    if err := png.Encode(f, img); err != nil {
        f.Close()
        log.Fatal(err)
    }

    if err := f.Close(); err != nil {
        log.Fatal(err)
    }
}

func main(){
 file,_ := os.Open("test.png")
 data,_ :=  ioutil.ReadAll(file)
 texts := []*DrawTextInfo{&DrawTextInfo{Text:"你 好",X:15,Y:60},&DrawTextInfo{Text:"世 界",X:15,Y:135}}
 if err := DrawStringOnImageAndSave("test2.png",data,texts);err != nil{
   fmt.Println(err)
  }
}

func NewTextBrush(FontFilePath string,FontSize float64,FontColor *image.Uniform,textWidth int)(*TextBrush,error){
    fontFile,err := ioutil.ReadFile(FontFilePath)
    if err != nil{
        return nil,err
    }
    fontType,err := truetype.Parse(fontFile)
    if err != nil {
        return nil, err
    }
    if textWidth <= 0 {
        textWidth = 42
    }
    return &TextBrush{FontType:fontType,FontSize:FontSize,FontColor:FontColor,TextWidth:textWidth},nil
}



func DrawStringOnImageAndSave(imagePath string,imageData []byte,infos []*DrawTextInfo)(err error){
    var background image.Image 
    filetype := http.DetectContentType(imageData)
    switch filetype{
    case "image/jpeg","image/jpg":
        background,err = jpeg.Decode(bytes.NewReader(imageData))
        if err != nil{
            fmt.Println("jpeg error")
            return
        }
    case "image/gif":
        background,err = gif.Decode(bytes.NewReader(imageData))
        if err != nil{
            return  
        }
    case "image/png":
        background,err = png.Decode(bytes.NewReader(imageData))
        if err != nil{
            return
        }
    default:
        return err 
    }

    des := image.NewRGBA(background.Bounds()) 
    textBrush,_ := NewTextBrush("/Users/mac/PingFang.ttf",60,image.White,60)

    c := freetype.NewContext()
    c.SetDPI(72)
    c.SetFont(textBrush.FontType)
    c.SetHinting(font.HintingFull)
    c.SetFontSize(textBrush.FontSize)
    c.SetClip(des.Bounds()) 
    c.SetDst(des)
    textBrush.FontColor = image.NewUniform(color.RGBA{
        R: 0XFF,
        G: 0XFF,
        B: 0XFF,
        A: 255,
    })

    c.SetSrc(textBrush.FontColor)

    for _, info := range infos{
        c.DrawString(info.Text,freetype.Pt(info.X,info.Y))
    }

    fSave, err := os.Create(imagePath)
    if err != nil{
        return err
    }
    defer fSave.Close()
    err = png.Encode(fSave,des)
    if err != nil{
        return err
    }
    return nil 

}
