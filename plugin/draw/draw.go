// Package draw 服务详情
package draw

import (
	"errors"
	"fmt"
	"image"
	"math/rand"
	"os"
	"strings"
	"sync"

	"github.com/FloatTech/floatbox/file"
	"github.com/FloatTech/floatbox/web"
	"github.com/FloatTech/gg"
	"github.com/FloatTech/imgfactory"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/img/text"
	"github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

type mpic struct {
	path       string      // 看板娘图片路径
	w          int         // 看板娘图片宽度
	h          int         // 看板娘图片高度
	isDisplay  bool        // 显示看板娘
	isCustom   bool        // 自定义看板娘
	statusText string      // 状态文本
	enableText bool        // 启用状态
	isDouble   bool        // 双列排版
	pluginName string      // 插件名
	font1      string      // 字体1
	font2      string      // 字体2
	info       []string    // 插件信息
	info2      []string    // 插件信息2
	multiple   float64     // 图片拓展倍数
	fontSize   float64     // 字体大小
	im         image.Image // 图片
}

type titleColor struct {
	r, g, b  int  // 颜色
	isRandom bool // 是否随机
}

type location struct {
	lastH            int     // 上一个高度
	drawX, maxTwidth float64 // 文字边距
	rlineX, rlineY   float64 // 宽高记录
	rtitleW          float64 // 标题位置
}

const (
	kanbanPath = "data/Control/"
)

var (
	customKanban = false // 自定义看板娘
	kanbanEnable = true  // 开关
	roleName     = "kanban.png"
)

func init() {
	if !file.IsExist(kanbanPath + "img") {
		err := os.MkdirAll(kanbanPath+"img", 0755)
		if err != nil {
			panic(err)
		}
	}
}

func init() { // 主函数
	en := control.Register("draw", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Brief:            "服务详情",
		Help:             "- 服务详情\n",
	})
	en.OnCommandGroup([]string{"服务详情", "service_detail"}, zero.UserOrGrpAdmin).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		lenmap := len(ctx.State["manager"].(*ctrl.Control[*zero.Ctx]).Manager.M)
		i, j := 0, 0
		double, fontSize, multiple := true, 40.0, 5.0
		gid := ctx.Event.GroupID
		if gid == 0 {
			gid = -ctx.Event.UserID
		}
		var tmp strings.Builder
		var tmp2 strings.Builder
		tmp.Grow(lenmap * 100)
		tmp2.Grow(lenmap * 100)

		end := lenmap / 2
		if lenmap <= 5 {
			double = false // 单列模式
			fontSize = 40
			multiple = 3
		}
		tab := "\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\n"
		tmp.WriteString("\t\t\t\t  <---------服务详情--------->   \t\t\t\t")
		// managers.ForEach(func(key string, service *ctrl.Control[*zero.Ctx]) bool {
		control.ForEachByPrio(func(p int, service *ctrl.Control[*zero.Ctx]) bool {
			key := service.Service
			i++
			if i > end+1 && lenmap > 5 {
				goto label
			}
			tmp.WriteString(fmt.Sprint("\n", i, ": ", service.EnableMarkIn(gid), " ", key,
				tab, strings.Trim(fmt.Sprint(service), "\n")))
			return true

		label:
			if j > 0 {
				tmp2.WriteString(fmt.Sprint("\n", i, ": ", service.EnableMarkIn(gid), " ", key,
					tab, strings.Trim(fmt.Sprint(service), "\n")))
			} else {
				tmp2.WriteString(fmt.Sprint(i, ": ", service.EnableMarkIn(gid), " ", key,
					tab, strings.Trim(fmt.Sprint(service), "\n")))
			}
			j++
			return true
		})

		msg := strings.Split(tmp.String(), "\n")
		msg2 := strings.Split(tmp2.String(), "\n")
		var menu = mpic{
			path:       kanbanPath + roleName, // 看板娘图片
			isDisplay:  kanbanEnable,          // 显示看板娘
			isCustom:   customKanban,
			statusText: "○ INFO",            // 启用状态
			enableText: true,                // 启用状态
			isDouble:   double,              // 双列排版
			pluginName: "ZeroBot-Plugin",    // 插件名
			font1:      text.BoldFontFile,   // 字体1
			font2:      text.SakuraFontFile, // 字体2
			info:       msg,                 // 插件信息
			info2:      msg2,                // 插件信息
			multiple:   multiple,            // 倍数
			fontSize:   fontSize,            // 字体大小
		}
		err := menu.loadpic()
		if err != nil {
			ctx.SendChain(message.Text("ERROR: ", err))
			return
		}
		pic, err := menu.dyna(&location{
			lastH:     0,
			drawX:     0.0,
			maxTwidth: 1200.0, // 文字边距
			rlineX:    0.0,    // 宽高记录
			rlineY:    140.0,
		})
		if err != nil {
			ctx.SendChain(message.Text("ERROR: ", err))
			return
		}
		data, err := imgfactory.ToBytes(pic) // 图片放入缓存
		if err != nil {
			ctx.SendChain(message.Text("ERROR: ", err))
			return
		}
		if id := ctx.SendChain(message.ImageBytes(data)); id.ID() == 0 {
			ctx.SendChain(message.Text("ERROR: 可能被风控了"))
		}
	})
}

// 返回菜单图片
func (mp *mpic) dyna(lt *location) (image.Image, error) {
	titleW := 1280          // 标题文字
	fontSize := mp.fontSize // 图片宽度和字体大小
	if mp.isDisplay && !mp.isDouble {
		titleW += mp.w // 标题位置
	}
	one := gg.NewContext(titleW, 256+len(mp.info)*15) // 新图像
	if mp.isDouble && mp.isDisplay || mp.isDouble && !mp.isDisplay {
		one = gg.NewContext(titleW*2, 256+len(mp.info)*15) // 新图像
		lt.rtitleW = float64(one.W()) / 2
	} else {
		lt.rtitleW = float64(one.W() - 1280)
	}
	one.SetRGB255(255, 255, 255)
	one.Clear()
	if err := one.LoadFontFace(mp.font2, fontSize*2); err != nil {
		return nil, err
	}
	one.SetRGBA255(55, 55, 55, 255) // 字体颜色
	switch {
	case mp.isDouble && !mp.isDisplay, !mp.isDouble && !mp.isDisplay:
		width, _ := one.MeasureString(mp.pluginName)
		one.DrawString(mp.pluginName, (float64(one.W())-width)/2, fontSize*2) // 绘制插件名在中间
	case mp.isDouble && mp.isDisplay: //
		one.DrawString(mp.pluginName, lt.rtitleW*1.36, fontSize*2) // 绘制插件名在右边
	default:
		width, _ := one.MeasureString(mp.pluginName)
		one.DrawString(mp.pluginName, (1280.0-width)/2+float64(mp.w), fontSize*2) // 绘制插件名在右边
	}
	if err := one.LoadFontFace(mp.font1, fontSize); err != nil { // 加载字体
		return nil, err
	}
	if mp.isDouble && !mp.isDisplay || !mp.isDouble && !mp.isDisplay {
		one.DrawRoundedRectangle(27, fontSize-5, fontSize*4.5, fontSize*1.5, 10) // 创建圆角矩形
	} else if mp.isDouble || mp.isDisplay { //
		one.DrawRoundedRectangle(lt.rtitleW+27, fontSize-5, fontSize*4.5, fontSize*1.5, 10) // 创建圆角矩形
	}
	if mp.enableText { // 如果启用
		one.SetRGBA255(15, 175, 15, 200)   // 设置绿色
		one.Fill()                         // 填充
		one.SetRGBA255(255, 255, 255, 255) // 设置白色
	} else {
		one.SetRGBA255(200, 15, 15, 200) // 设置红色
		one.Fill()
		one.SetRGBA255(255, 255, 255, 255)
	}
	if mp.isDouble && !mp.isDisplay || !mp.isDouble && !mp.isDisplay {
		one.DrawString(mp.statusText, 35, fontSize*2) // 绘制启用状态
	} else if mp.isDouble || mp.isDisplay { //
		one.DrawString(mp.statusText, lt.rtitleW+35, fontSize*2) // 绘制启用状态
	}
	return mp.createPic(one, lt)
}

// 创建图片
func (mp *mpic) createPic(one *gg.Context, lt *location) (image.Image, error) {
	var wg sync.WaitGroup
	if mp.isDouble {
		var imgs [2]image.Image
		var err1, err2 error
		wg.Add(2)
		titlec := titleColor{isRandom: false}
		titlec.randfill()
		go func() {
			imgs[0], err1 = mp.createPic2(lt, &titlec, mp.info)
			wg.Done()
		}()
		go func() {
			imgs[1], err2 = mp.createPic2(lt, &titlec, mp.info2)
			wg.Done()
		}()
		wg.Wait()
		if err1 != nil {
			return nil, err1
		}
		if err2 != nil {
			return nil, err2
		}
		imRY := imgs[0].Bounds().Dy() + 100 // 右边 图像的高度
		imLY := imgs[1].Bounds().Dy() + 5   // 左边 图像的高度
		max := 0
		if mp.isDisplay {
			imLY += mp.h + 5
		}
		if imRY > imLY {
			max = imRY
		} else {
			max = imLY
		}
		if max > one.H() {
			imgtmp := gg.NewContext(one.W(), max+int(mp.fontSize)) // 高度
			imgtmp.SetRGB255(255, 255, 255)
			imgtmp.Clear()
			imgtmp.DrawImage(one.Image(), 0, 0)
			one = gg.NewContextForImage(imgtmp.Image())
		}
		if mp.isDisplay {
			if imLY > one.H() {
				imgtmp := gg.NewContext(one.W(), imLY) // 高度
				imgtmp.SetRGB255(255, 255, 255)
				imgtmp.Clear()
				imgtmp.DrawImage(one.Image(), 0, 0)
				one = gg.NewContextForImage(imgtmp.Image())
			}
			one.DrawImage(mp.im, 0, 0) // 放入看板娘
			one.DrawImage(imgs[1], 0, mp.h+5)
			one.DrawImage(imgs[0], 1280, 0) // 最终的绘制位置
		} else {
			one.DrawImage(imgs[0], 0, 0) // 最终的绘制位置
			one.DrawImage(imgs[1], 1280, 0)
		}
	} else {
		titlec := titleColor{isRandom: true}
		titlec.randfill()
		var img image.Image
		var err1 error
		wg.Add(1)
		go func() {
			img, err1 = mp.createPic2(lt, &titlec, mp.info)
			wg.Done()
		}()
		wg.Wait()
		if err1 != nil {
			return nil, err1
		}
		imY := img.Bounds().Dy()
		if imY+int(mp.fontSize) > one.H() {
			imgtmp := gg.NewContext(one.W(), imY) // 高度
			imgtmp.SetRGB255(255, 255, 255)
			imgtmp.Clear()
			imgtmp.DrawImage(one.Image(), 0, 0)
			one = gg.NewContextForImage(imgtmp.Image())
		}
		if mp.isDisplay {
			if mp.h > one.H() {
				imgtmp := gg.NewContext(one.W(), mp.h) // 宽和高
				imgtmp.SetRGB255(255, 255, 255)
				imgtmp.Clear()
				imgtmp.DrawImage(one.Image(), 0, 0)
				one = gg.NewContextForImage(imgtmp.Image())
			}
			if mp.isCustom { // 如果自定义看板娘
				one.DrawImage(mp.im, 0, (one.H()-mp.h)/2) // 放入看板娘
			} else {
				one.DrawImage(mp.im, 0, 0) // 最终的绘制位置
			}
			one.DrawImage(img, one.W()-1280, 50) // 最终的绘制位置
		} else {
			one.DrawImage(img, 0, 50) // 最终的绘制位置
		}
	}
	return one.Image(), nil
}

// 创建图片
func (mp *mpic) createPic2(lt *location, titlec *titleColor, info []string) (image.Image, error) {
	fontSize := mp.fontSize
	one := gg.NewContext(1280, 256+len(mp.info)*15)
	if err := one.LoadFontFace(mp.font1, fontSize); err != nil { // 加载字体
		return nil, err
	}
	lineTexts := make([]string, 0, 32)
	rlx := lt.rlineX
	rly := lt.rlineY
	lh := lt.lastH
	for i := 0; i < len(info); i++ { // 遍历文字切片
		lineText, textW, textH, tmpw := "", 0.0, 0.0, 0.0
		if mp.isDouble {
			if strings.Contains(info[i], ": ● ") || strings.Contains(info[i], ": ○ ") {
				titlec.randfill() // 随机一次颜色
			}
		}
		for len(info[i]) > 0 {
			lineText, tmpw = truncate(one, info[i], lt.maxTwidth)
			lineTexts = append(lineTexts, lineText)
			if tmpw > textW {
				textW = tmpw
			}
			if len(lineText) >= len(info[i]) {
				break // 如果写入的文本大于等于本次写入的文本 则跳出
			}
			textH += fontSize * 1.3           // 截断一次加一行高度
			info[i] = info[i][len(lineText):] // 丢弃已经写入的文字并重新赋值
		}
		threeW, threeH := textW+fontSize, (textH + (fontSize * 1.2)) // 圆角矩形宽度和高度
		if int(rlx+textW)+int(fontSize*2) > one.W() {                // 越界
			rly += float64(lh) + fontSize/4              // 加一次高度
			rlx = 5                                      // 重置宽度位置
			if threeH+rly+fontSize >= float64(one.H()) { // 超出最大高度则进行加高
				imgtmp := gg.NewContext(one.W(), int(rly+threeH*mp.multiple)) // 高度
				imgtmp.DrawImage(one.Image(), 0, 0)
				one = gg.NewContextForImage(imgtmp.Image())
				if err := one.LoadFontFace(mp.font1, mp.fontSize); err != nil { // 加载字体
					return nil, err
				}
			}
			dx := rlx + 13                                          // 圆角矩形位置宽度
			one.DrawRoundedRectangle(dx, rly, threeW, threeH, 20.0) // 创建圆角矩形
			titlec.drawsc(one, fontSize, dx, rly, lineTexts)
			rlx += threeW + fontSize/2 // 添加后加一次宽度
			lh = int(threeH)
		} else {
			dx := rlx + 13                                          // 圆角矩形位置宽度
			one.DrawRoundedRectangle(dx, rly, threeW, threeH, 20.0) // 创建圆角矩形
			titlec.drawsc(one, fontSize, dx, rly, lineTexts)
			rlx += threeW + fontSize/2 // 添加后加一次宽度
			lh = int(threeH)
		}
		lineTexts = lineTexts[:0]
	}
	return one.Image(), nil
}

// 绘制文字
func (titlec *titleColor) drawsc(one *gg.Context, fontSize, drawX, rlineY float64, lineTexts []string) {
	if titlec.isRandom {
		titlec.randfill()
	}
	one.SetRGBA255(titlec.r, titlec.g, titlec.b, 85)
	one.Fill() // 填充颜色
	one.SetRGBA255(55, 55, 55, 255)
	h := fontSize + rlineY - 3
	for i := range lineTexts { // 逐行绘制文字
		one.DrawString(lineTexts[i], drawX+fontSize/2, h)
		h += fontSize + (fontSize / 4)
	}
}

// 填充颜色
func (titlec *titleColor) randfill() {
	titlec.r = rand.Intn(245) // 随机颜色
	titlec.g = rand.Intn(245)
	titlec.b = rand.Intn(245)
	for titlec.r < 15 || titlec.r > 210 {
		titlec.r = rand.Intn(245)
	}
	for titlec.g < 15 || titlec.g > 210 {
		titlec.g = rand.Intn(245)
	}
	for titlec.b < 15 || titlec.b > 210 {
		titlec.b = rand.Intn(245)
	}
}

// 截断文字
func truncate(one *gg.Context, text string, maxW float64) (string, float64) {
	var tmp strings.Builder
	tmp.Grow(len(text))
	res, w := make([]rune, 0, len(text)), 0.0
	for _, r := range text {
		tmp.WriteRune(r)
		width, _ := one.MeasureString(tmp.String()) // 获取文字宽度
		if width > maxW {                           // 如果宽度大于文字边距
			break // 跳出
		} else {
			w = width
			res = append(res, r) // 写入
		}
	}
	return string(res), w
}

// 编码看板娘图片和加载字体
func (mp *mpic) loadpic() error {
	if !file.IsExist(mp.font1) { // 获取字体
		err := downloadData(text.BoldFontFile, "https://gitcode.net/u011570312/zbpdata/-/raw/main/"+text.BoldFontFile[5:]+"?inline=true")
		if err != nil {
			return err
		}
	}
	if !file.IsExist(mp.font2) { // 获取字体
		err := downloadData(text.SakuraFontFile, "https://gitcode.net/u011570312/zbpdata/-/raw/main/"+text.SakuraFontFile[5:]+"?inline=true")
		if err != nil {
			return err
		}
	}
	if mp.isDisplay {
		f, err := os.Open(mp.path)
		if err != nil {
			return err
		}
		defer f.Close()
		mp.im, _, err = image.Decode(f)
		if err != nil {
			return err
		}
		mp.im = imgfactory.Limit(mp.im, 1280, 1280)
		b := mp.im.Bounds().Size()
		mp.w, mp.h = b.X, b.Y
	}
	return nil
}

// 下载数据,url=dataurl+path[5:]+"?inline=true"
func downloadData(path, url string) error {
	data, err := web.RequestDataWith(web.NewTLS12Client(), url, "GET", "gitcode.net", web.RandUA(), nil)
	if err != nil {
		return errors.New("下载" + url + "失败!")
	}
	logrus.Printf("[file]从镜像下载数据%d字节...", len(data))
	if len(data) == 0 {
		return errors.New("read body len <= 0")
	}
	// 写入数据
	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return err
	}
	logrus.Printf("[file]下载文件成功")
	return nil
}
