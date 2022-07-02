package gif

import (
	"errors"
	"image/color"
	"math/rand"
	"os"
	"strconv"
	"sync"

	"github.com/FloatTech/zbputils/file"
	"github.com/FloatTech/zbputils/img"
	"github.com/FloatTech/zbputils/img/text"
	"github.com/FloatTech/zbputils/img/writer"
	"github.com/fogleman/gg"
)

// Pa 爬
func (cc *context) Pa(value ...string) (string, error) {
	name := cc.usrdir + `爬.png`
	tou, err := cc.getLogo(0, 0)
	if err != nil {
		return "", err
	}
	// 随机爬图序号
	rand := rand.Intn(60) + 1
	if file.IsNotExist(datapath + "materials/pa") {
		err = os.MkdirAll(datapath+"materials/pa", 0755)
		if err != nil {
			return "", err
		}
	}
	f, err := dlblock(`pa/` + strconv.Itoa(rand) + `.png`)
	if err != nil {
		return "", err
	}
	imgf, err := img.LoadFirstFrame(f, 0, 0)
	if err != nil {
		return "", err
	}
	return "file:///" + name, writer.SavePNG2Path(name, imgf.InsertBottom(tou, 100, 100, 0, 400).Im)
}

// Si 撕
func (cc *context) Si(value ...string) (string, error) {
	name := cc.usrdir + `撕.png`
	tou, err := cc.getLogo(0, 0)
	if err != nil {
		return "", err
	}
	im1 := img.Rotate(tou, 20, 380, 380)
	im2 := img.Rotate(tou, -12, 380, 380)
	if file.IsNotExist(datapath + "materials/si") {
		err = os.MkdirAll(datapath+"materials/si", 0755)
		if err != nil {
			return "", err
		}
	}
	f, err := dlblock(`si/0.png`)
	if err != nil {
		return "", err
	}
	imgf, err := img.LoadFirstFrame(f, 0, 0)
	if err != nil {
		return "", err
	}
	return "file:///" + name, writer.SavePNG2Path(name, imgf.InsertBottom(im1.Im, im1.W, im1.H, -3, 370).InsertBottom(im2.Im, im2.W, im2.H, 653, 310).Im)
}

// FlipV 上翻,下翻
func (cc *context) FlipV(value ...string) (string, error) {
	name := cc.usrdir + `FlipV.png`
	// 加载图片
	im, err := img.LoadFirstFrame(cc.headimgsdir[0], 0, 0)
	if err != nil {
		return "", err
	}
	imgnrgba := im.FlipV().Im
	return "file:///" + name, writer.SavePNG2Path(name, imgnrgba)
}

// FlipH 左翻,右翻
func (cc *context) FlipH(value ...string) (string, error) {
	name := cc.usrdir + `FlipH.png`
	// 加载图片
	im, err := img.LoadFirstFrame(cc.headimgsdir[0], 0, 0)
	if err != nil {
		return "", err
	}
	imgnrgba := im.FlipH().Im
	return "file:///" + name, writer.SavePNG2Path(name, imgnrgba)
}

// Invert 反色
func (cc *context) Invert(value ...string) (string, error) {
	name := cc.usrdir + `Invert.png`
	// 加载图片
	im, err := img.LoadFirstFrame(cc.headimgsdir[0], 0, 0)
	if err != nil {
		return "", err
	}
	imgnrgba := im.Invert().Im
	return "file:///" + name, writer.SavePNG2Path(name, imgnrgba)
}

// Blur 反色
func (cc *context) Blur(value ...string) (string, error) {
	name := cc.usrdir + `Blur.png`
	// 加载图片
	im, err := img.LoadFirstFrame(cc.headimgsdir[0], 0, 0)
	if err != nil {
		return "", err
	}
	imgnrgba := im.Blur(10).Im
	return "file:///" + name, writer.SavePNG2Path(name, imgnrgba)
}

// Grayscale 灰度
func (cc *context) Grayscale(value ...string) (string, error) {
	name := cc.usrdir + `Grayscale.png`
	// 加载图片
	im, err := img.LoadFirstFrame(cc.headimgsdir[0], 0, 0)
	if err != nil {
		return "", err
	}
	imgnrgba := im.Grayscale().Im
	return "file:///" + name, writer.SavePNG2Path(name, imgnrgba)
}

// InvertAndGrayscale 负片
func (cc *context) InvertAndGrayscale(value ...string) (string, error) {
	name := cc.usrdir + `InvertAndGrayscale.png`
	// 加载图片
	im, err := img.LoadFirstFrame(cc.headimgsdir[0], 0, 0)
	if err != nil {
		return "", err
	}
	imgnrgba := im.Invert().Grayscale().Im
	return "file:///" + name, writer.SavePNG2Path(name, imgnrgba)
}

// Convolve3x3 浮雕
func (cc *context) Convolve3x3(value ...string) (string, error) {
	name := cc.usrdir + ` Convolve3x3.png`
	// 加载图片
	im, err := img.LoadFirstFrame(cc.headimgsdir[0], 0, 0)
	if err != nil {
		return "", err
	}
	imgnrgba := im.Convolve3x3().Im
	return "file:///" + name, writer.SavePNG2Path(name, imgnrgba)
}

// Rotate 旋转,带参数暂时不用
func (cc *context) Rotate(value ...string) (string, error) {
	name := cc.usrdir + `Rotate.png`
	// 加载图片
	im, err := img.LoadFirstFrame(cc.headimgsdir[0], 0, 0)
	if err != nil {
		return "", err
	}
	r, _ := strconv.ParseFloat(value[0], 64)
	imgnrgba := img.Rotate(im.Im, r, 0, 0).Im
	return "file:///" + name, writer.SavePNG2Path(name, imgnrgba)
}

// Deformation 变形,带参数暂时不用
func (cc *context) Deformation(value ...string) (string, error) {
	name := cc.usrdir + `Deformation.png`
	// 加载图片
	im, err := img.LoadFirstFrame(cc.headimgsdir[0], 0, 0)
	if err != nil {
		return "", err
	}
	w, err := strconv.Atoi(value[0])
	if err != nil {
		return "", err
	}
	h, err := strconv.Atoi(value[1])
	if err != nil {
		return "", err
	}
	imgnrgba := img.Size(im.Im, w, h).Im
	return "file:///" + name, writer.SavePNG2Path(name, imgnrgba)
}

// Alike 你像个xxx一样
func (cc *context) Alike(args ...string) (string, error) {
	var wg sync.WaitGroup
	var m sync.Mutex
	var err error
	c := dlrange("alike", 1, &wg, func(e error) {
		m.Lock()
		err = e
		m.Unlock()
	})
	if err != nil {
		return "", err
	}
	wg.Wait()
	name := cc.usrdir + `Alike.png`
	back, err := gg.LoadImage(c[0])
	if err != nil {
		return "", err
	}
	canvas := gg.NewContextForImage(back)
	canvas.SetColor(color.Black)
	_, _ = file.GetLazyData(text.BoldFontFile, true)
	fontSize := 26.0
	_ = canvas.LoadFontFace(text.BoldFontFile, fontSize)
	l, _ := canvas.MeasureString(args[0])
	if l > fontSize*4 {
		err = errors.New("输入字符太长")
		return "", err
	}
	canvas.DrawString(args[0], 180-l/2, 65)
	return "file:///" + name, canvas.SavePNG(name)
}

// Marriage
func (cc *context) Marriage(args ...string) (string, error) {
	var wg sync.WaitGroup
	var m sync.Mutex
	var err error
	c := dlrange("marriage", 2, &wg, func(e error) {
		m.Lock()
		err = e
		m.Unlock()
	})
	if err != nil {
		return "", err
	}
	wg.Wait()
	imgs, err := loadFirstFrames(c, 2)
	if err != nil {
		return "", err
	}
	name := cc.usrdir + `Marriage.png`
	im, err := img.LoadFirstFrame(cc.headimgsdir[0], 800, 1080)
	if err != nil {
		return "", err
	}
	imgnrgba := im.InsertUp(imgs[0].Im, 0, 0, 0, 0).InsertUp(imgs[1].Im, 0, 0, 520, 0).Im
	return "file:///" + name, writer.SavePNG2Path(name, imgnrgba)
}
