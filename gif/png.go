package gif

import (
	"math/rand"
	"os"
	"strconv"

	"github.com/FloatTech/zbputils/file"
	"github.com/FloatTech/zbputils/img"
	"github.com/FloatTech/zbputils/img/writer"
)

// Pa 爬
func (cc *context) Pa() (string, error) {
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
func (cc *context) Si() (string, error) {
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
func (cc *context) FlipV() (string, error) {
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
func (cc *context) FlipH() (string, error) {
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
func (cc *context) Invert() (string, error) {
	name := cc.usrdir + `Invert.png`
	// 加载图片
	im, err := img.LoadFirstFrame(cc.headimgsdir[0], 0, 0)
	if err != nil {
		return "", err
	}
	imgnrgba := im.Invert().Im
	return "file:///" + name, writer.SavePNG2Path(name, imgnrgba)
}

// Grayscale 灰度
func (cc *context) Grayscale() (string, error) {
	name := cc.usrdir + `Grayscale.png`
	// 加载图片
	im, err := img.LoadFirstFrame(cc.headimgsdir[0], 0, 0)
	if err != nil {
		return "", err
	}
	imgnrgba := im.Grayscale().Im
	return "file:///" + name, writer.SavePNG2Path(name, imgnrgba)
}

// InsertAndGrayscale 负片
func (cc *context) InsertAndGrayscale() (string, error) {
	name := cc.usrdir + `InsertAndGrayscale.png`
	// 加载图片
	im, err := img.LoadFirstFrame(cc.headimgsdir[0], 0, 0)
	if err != nil {
		return "", err
	}
	imgnrgba := im.Invert().Grayscale().Im
	return "file:///" + name, writer.SavePNG2Path(name, imgnrgba)
}

// Convolve3x3 浮雕
func (cc *context) Convolve3x3() (string, error) {
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
