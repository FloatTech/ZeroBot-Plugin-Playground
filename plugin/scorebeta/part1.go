package main

import (
	"fmt"
	"github.com/FloatTech/gg"
	"github.com/FloatTech/imgfactory"
	"image"
	"net/http"
)

var exampleName = "name"
var dayinfo = "下午好"
var coin = 100
var nowcoin = 2000
var level = 3
var date = "2023年3月22日"
var avatarlink = "link"

func main() {
	// load picture that's example.
	// background
	background, err := gg.LoadImage("./original.jpg")
	if err != nil {
		panic(err)
		return
	}
	// resize background
	back := imgfactory.Limit(background, 1280, 720)
	backWidth := back.Bounds().Dx()  // get width
	backHeight := back.Bounds().Dy() // get height
	mainCanvas := gg.NewContext(backWidth, backHeight)
	mainCanvas.DrawImage(background, 0, 0) // draw background
	// draw rounded rectangle (three parts)
	// name,avatar part

	// draw rounded rectangle
	mainCanvas.DrawRoundedRectangle(100, 100, float64(backWidth-200), float64(200), 16)
	mainCanvas.SetRGB255(255, 255, 255)
	mainCanvas.SetLineWidth(3)
	mainCanvas.Stroke()
	mainCanvas.SetRGBA(0, 0, 0, 0.5)
	mainCanvas.DrawRoundedRectangle(100, 100, float64(backWidth-200), float64(200), 16)
	mainCanvas.Fill()
	// draw avatar
	avatarByte, err := http.Get(avatarlink)
	if err != nil {
		panic(err)
		return
	}
	avatarByteUni, _, _ := image.Decode(avatarByte.Body)
	avatarFormat := imgfactory.Size(avatarByteUni, 200, 200)
	mainCanvas.DrawImage(avatarFormat.Circle(0).Image(), 130, 30)
	// combine and draw name,and other info
	mainCanvas.LoadFontFace("./font.ttf", 50)
	mainCanvas.SetRGB255(255, 255, 255)
	mainCanvas.DrawStringAnchored(fmt.Sprintf("%s, %s", dayinfo, exampleName), 400, 200, 0, 0)
	// draw the second part
	mainCanvas.DrawRoundedRectangle(100, 350, float64(((backWidth-200)/2)-50), float64(200), 16)
	mainCanvas.SetRGB255(255, 255, 255)
	mainCanvas.SetLineWidth(3)
	mainCanvas.Stroke()
	mainCanvas.SetRGBA(0, 0, 0, 0.5)
	mainCanvas.DrawRoundedRectangle(100, 350, float64(((backWidth-200)/2)-50), float64(200), 16)
	mainCanvas.Fill()
	mainCanvas.LoadFontFace("./font.ttf", 35)
	mainCanvas.SetRGB255(255, 255, 255)
	mainCanvas.DrawStringAnchored(fmt.Sprintf("ATRI币 +%d", nowcoin), 130, 420, 0, 0)
	mainCanvas.DrawStringAnchored(fmt.Sprintf("当前ATRI币 %d", coin), 130, 470, 0, 0)
	mainCanvas.DrawStringAnchored(fmt.Sprintf("Level: %d", level), 130, 520, 0, 0)
	mainCanvas.Fill()
	// draw the third part
	mainCanvas.DrawRoundedRectangle(100+float64(((backWidth-200)/2)+50), 350, float64(((backWidth-200)/2)-50), float64(200), 16)
	mainCanvas.SetRGB255(255, 255, 255)
	mainCanvas.SetLineWidth(3)
	mainCanvas.Stroke()
	mainCanvas.SetRGBA(0, 0, 0, 0.5)
	mainCanvas.DrawRoundedRectangle(100+float64(((backWidth-200)/2)+50), 350, float64(((backWidth-200)/2)-50), float64(200), 16)
	mainCanvas.Fill()
	mainCanvas.LoadFontFace("./font.ttf", 45)
	mainCanvas.SetRGB255(255, 255, 255)
	mainCanvas.DrawStringAnchored(fmt.Sprintf("%s", date), 150+float64(((backWidth-200)/2)+50), 420, 0, 0)
	mainCanvas.LoadFontFace("./font.ttf", 20)
	mainCanvas.DrawStringAnchored(fmt.Sprintf("Created By Zerobot-Plugin v1.70-beta5"), 150+float64(((backWidth-200)/2)+20), 520, 0, 0)
	mainCanvas.Fill()
	mainCanvas.SaveJPG("./result.jpg", 100)
}
