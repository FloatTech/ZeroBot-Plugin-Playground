package klala

import (
	"image"
	"image/color"
	"math"
	"strconv"
	"sync"

	"github.com/FloatTech/gg"
	img "github.com/FloatTech/imgfactory"
)

const (
	NameFont    = "data/klala/kkk/font/NZBZ.ttf"                    // NameFont 名字字体
	FontFile    = "data/klala/kkk/font/SourceHanMonoSC-HeavyIt.ttf" // FontFile 汉字字体
	FiFile      = "data/klala/kkk/font/tttgbnumber.ttf"             // FiFile 其余字体(数字英文)
	BaFile      = "data/klala/kkk/font/STLITI.TTF"                  // BaFile 华文隶书版本版本号字体
	windowsPath = "data/klala/kkk/sund/冰.jpg"
	refinePath  = "data/klala/kkk/sund/refine.png"
	skillSdPic  = "data/klala/kkk/sund/mz.png"
	lightPath   = "data/klala/kkk/icon/light_cone/"
	liHuiPath   = "data/klala/kkk/lihui/"
	remainPath  = "data/klala/kkk/icon/relic/"
	tPicPath    = "data/klala/kkk/icon/skill/"
)

var skillList = []string{"_rank1.png", "_rank2.png", "_ultimate.png", "_rank4.png", "_skill.png", "_rank6.png", "_basic_atk.png", "_talent.png"} // 0-5为星魂,6-7为普攻+天赋

func (t *thisdata) drawcard(n int) (string, error) {
	var wg sync.WaitGroup
	wg.Add(5)
	yinyinBlack127 := color.NRGBA{R: 0, G: 0, B: 0, A: 127}
	dc := gg.NewContext(1080, 1680)
	dc.SetRGB(1, 1, 1)
	if err := dc.LoadFontFace(FontFile, 40); err != nil {
		panic(err)
	}
	beijing, err := gg.LoadImage(windowsPath)
	if err != nil {
		return "", err
	}
	beijing = img.Size(beijing, 0, 1680).Image()
	dc.DrawImageAnchored(beijing, 540, 0, 0.5, 0)
	lihui, err := gg.LoadPNG(liHuiPath + strconv.Itoa(t.RoleData[n].ID) + ".png")
	if err != nil {
		return "", err
	}
	lihui = img.Size(lihui, 0, 880).Image()
	sxx := lihui.Bounds().Size().X
	dc.DrawImage(lihui, int(300-float64(sxx)/2), 0)
	// 昵称框图
	go func() {
		defer wg.Done()
		zero := gg.NewContext(540, 200)
		zero.SetRGB(1, 1, 1) // 白色
		// 角色名字
		if err := zero.LoadFontFace(NameFont, 80); err != nil {
			panic(err)
		}
		zero.DrawStringAnchored(t.RoleData[n].Name, 505, 130, 1, 0)
		if err := zero.LoadFontFace(FontFile, 30); err != nil {
			panic(err)
		}
		zero.DrawStringAnchored("昵称:"+t.Nickname, 505, 40, 1, 0)
		if err := zero.LoadFontFace(FiFile, 30); err != nil {
			panic(err)
		}
		zero.DrawStringAnchored("UID:"+t.UID+"-LV."+strconv.Itoa(t.Level), 505, 180, 1, 0)
		newying := Yinying(540, 200, 16, color.NRGBA{R: 0, G: 0, B: 0, A: 106})
		dc.DrawImage(newying, 505, 20)
		dc.DrawImage(zero.Image(), 505, 20)
	}()
	// 星魂
	go func() {
		defer wg.Done()
		ten := gg.NewContext(470, 80)
		sdPic, err := gg.LoadImage(skillSdPic)
		if err != nil {
			panic(err)
		}
		sdPic = img.Size(sdPic, 0, 80).Image()
		sdPicBlack := adjustOpacity(sdPic, 0.5)
		var sd image.Image
		for a := 0; a < 6; a++ {
			if skillpic, err := gg.LoadImage(tPicPath + strconv.Itoa(t.RoleData[n].ID) + skillList[a]); err == nil {
				skillpic = img.Size(skillpic, 0, 60).Image()
				if a >= t.RoleData[n].Rank {
					skillpic = adjustOpacity(skillpic, 0.5)
					sd = sdPicBlack
				} else {
					sd = sdPic
				}
				ten.DrawImageAnchored(sd, a*80, 40, 0, 0.5)
				ten.DrawImageAnchored(skillpic, 12+a*80, 40, 0, 0.5)
			}
		}
		dc.DrawImage(ten.Image(), 20, 630)
	}()
	// 属性列表
	go func() {
		defer wg.Done()
		one := gg.NewContext(540, 470)
		if err := one.LoadFontFace(FontFile, 30); err != nil {
			panic(err)
		}
		one.SetRGB(1, 1, 1) // 白色
		one.DrawString("角色等级:", 70, 45)
		one.DrawString("生命值:", 70, 96.25)
		one.DrawString("攻击力:", 70, 147.5)
		one.DrawString("防御力:", 70, 198.75)
		one.DrawString("速度:", 70, 250)
		one.DrawString("暴击率:", 70, 301.25)
		one.DrawString("暴击伤害:", 70, 352.5)
		one.DrawString("效果命中:", 70, 403.75)
		one.DrawString("效果抵抗:", 70, 455)
		// 值,一一对f应
		if err := one.LoadFontFace(FiFile, 30); err != nil {
			panic(err)
		}
		// 属性540*460,字30,间距15,60
		one.SetRGB(1, 1, 1)                                                                // 白色
		one.DrawStringAnchored("Lv"+strconv.Itoa(t.RoleData[n].List.Level), 470, 45, 1, 0) // Lv
		one.DrawStringAnchored("(+"+Ftoone(t.RoleData[n].List.HpFinal-t.RoleData[n].List.HpBase)+")_"+
			Ftoone(t.RoleData[n].List.HpFinal), 470, 96.25, 1, 0) // 生命
		one.DrawStringAnchored("(+"+Ftoone(t.RoleData[n].List.AttackFinal-t.RoleData[n].List.AttackBase)+")_"+
			Ftoone(t.RoleData[n].List.AttackFinal), 470, 147.5, 1, 0) // 攻击
		one.DrawStringAnchored("(+"+Ftoone(t.RoleData[n].List.DefenseFinal-t.RoleData[n].List.DefenseBase)+")_"+
			Ftoone(t.RoleData[n].List.DefenseFinal), 470, 198.75, 1, 0) // 防御
		one.DrawStringAnchored("(+"+Ftoone(t.RoleData[n].List.SpeedFinal-float64(t.RoleData[n].List.SpeedBase))+")_"+
			Ftoone(t.RoleData[n].List.SpeedFinal), 470, 250, 1, 0) // 速度
		one.DrawStringAnchored(Ftoone(t.RoleData[n].List.CriticalChance*100)+"%", 470, 301.25, 1, 0)    // 暴击
		one.DrawStringAnchored(Ftoone(t.RoleData[n].List.CriticalDamage*100)+"%", 470, 352.5, 1, 0)     // 爆伤
		one.DrawStringAnchored(Ftoone(t.RoleData[n].List.StatusProbability*100)+"%", 470, 403.75, 1, 0) // 效果命中
		one.DrawStringAnchored(Ftoone(t.RoleData[n].List.StatusResistance*100)+"%", 470, 455, 1, 0)     // 效果抵抗
		dc.DrawImage(Yinying(540, 470, 16, yinyinBlack127), 505, 240)                                   // 背景
		dc.DrawImage(one.Image(), 505, 240)
	}()
	// 光锥
	go func() {
		defer wg.Done()
		yinlight := Yinying(1040, 180, 16, yinyinBlack127)
		two := gg.NewContext(1040, 180)
		two.SetRGB(1, 1, 1) // 白色
		if err := two.LoadFontFace(FontFile, 30); err != nil {
			panic(err)
		}
		// 天赋
		for ii := 0; ii < 4; ii++ {
			var talentname string
			switch ii {
			case 0:
				talentname = skillList[6]
			case 1:
				talentname = skillList[4]
			case 2:
				talentname = skillList[2]
			default:
				talentname = skillList[7]
			}
			if tpic, err := gg.LoadImage(tPicPath + strconv.Itoa(t.RoleData[n].ID) + talentname); err == nil {
				tpic = img.Size(tpic, 0, 80).Image()
				two.DrawImage(tpic, 10+ii%2*300, 10+ii/2*80)
			}
		}
		two.DrawString("普攻 ", 90, 60)
		two.DrawString("战技 ", 390, 60)
		two.DrawString("终结 ", 90, 140)
		two.DrawString("天赋 ", 390, 140)
		if t.RoleData[n].Light.ID != 0 {
			// 图标
			if lpic, err := gg.LoadImage(lightPath + strconv.Itoa(t.RoleData[n].Light.ID) + ".png"); err == nil {
				lpic = img.Size(lpic, 0, 160).Image()
				two.DrawImage(lpic, 670, 20)
			}
			two.DrawString(t.RoleData[n].Light.Name, 830, 60)
			// 精炼
			if err := two.LoadFontFace(FiFile, 30); err != nil {
				panic(err)
			}
			if refpic, err := gg.LoadImage(refinePath); err == nil {
				refpic = adjustOpacity(refpic, 0.8)
				refpic = img.Size(refpic, 140, 0).Image()
				two.DrawImageAnchored(refpic, 970, 140, 0.5, 0.5)
				two.DrawStringAnchored("ref:"+strconv.Itoa(t.RoleData[n].Light.Rank), 970, 140, 0.5, 0.5)
			}
			// 星级
			two.DrawImageAnchored(img.Size(drawStars("#FFCC00", "#FFE43A", t.RoleData[n].Light.Star), 0, 30).Image(), 1020, 80, 1, 0)
			two.DrawString("LV."+strconv.Itoa(t.RoleData[n].Light.Level), 830, 150)
		}
		if err := two.LoadFontFace(FiFile, 30); err != nil {
			panic(err)
		}
		two.DrawString("LV."+strconv.Itoa(t.RoleData[n].Skill.A), 160, 60)
		two.DrawString("LV."+strconv.Itoa(t.RoleData[n].Skill.E), 460, 60)
		two.DrawString("LV."+strconv.Itoa(t.RoleData[n].Skill.Q), 160, 140)
		two.DrawString("LV."+strconv.Itoa(t.RoleData[n].Skill.T), 460, 140)
		dc.DrawImage(yinlight, 20, 720)
		dc.DrawImage(two.Image(), 20, 720)
	}()
	// 遗物
	go func() {
		defer wg.Done()
		yinsyw := Yinying(340, 350, 16, yinyinBlack127)
		var yw relicsdata
		for i := 0; i < 6; i++ {
			switch i {
			case 0:
				yw = t.RoleData[n].Relics.Head
			case 1:
				yw = t.RoleData[n].Relics.Hand
			case 2:
				yw = t.RoleData[n].Relics.Body
			case 3:
				yw = t.RoleData[n].Relics.Foot
			case 4:
				yw = t.RoleData[n].Relics.Neck
			default:
				yw = t.RoleData[n].Relics.Object
			}
			if yw.SetID == 0 {
				continue
			}
			// 字图层
			three := gg.NewContext(340, 350)
			if err := three.LoadFontFace(FontFile, 30); err != nil {
				panic(err)
			}
			// 字号30,间距50
			three.SetRGB(1, 1, 1) // 白色
			// 画线
			for c := 0; c < 4; c++ {
				three.DrawLine(0, 157+float64(c)*45, 350, 157+float64(c)*45) // 横线条分割
			}
			three.Stroke()
			if tuyw, err := gg.LoadImage(remainPath + strconv.Itoa(yw.SetID) + "_" + strconv.Itoa(i%4) + ".png"); err == nil {
				tuyw = img.Size(tuyw, 0, 90).Image()
				three.DrawImage(tuyw, 15, 15)
			}
			// 星级
			three.DrawImage(img.Size(drawStars("#FFCC00", "#FFE43A", yw.Star), 0, 20).Image(), 145, 60)
			// 遗物name
			three.DrawStringAnchored(yw.Name, 325, 50, 1, 0)
			// 圣遗物属性 主词条
			// 间隔45,初始145
			var xx, yy float64 // xx,yy词条相对位置,x,y文本框在全图位置
			xx = 15
			yy = 145
			// 主词条
			if err := three.LoadFontFace(FontFile, 30); err != nil {
				panic(err)
			}
			three.DrawString(yw.MainV.Name, xx, yy) // "主"
			if err := three.LoadFontFace(FiFile, 30); err != nil {
				panic(err)
			}
			// 主词条属性
			three.DrawStringAnchored("+"+yw.MainV.Value+stofen(yw.MainV.Name), 325, yy, 1, 0) // 主词条属性
			three.DrawString("+"+strconv.Itoa(int(yw.Level)), 85, 90)                         // LV
			three.SetHexColor("#98F5FF")                                                      // 蓝色
			for k := 0; k < len(yw.Vlist); k++ {
				yy += 45
				if err := three.LoadFontFace(FontFile, 30); err != nil {
					panic(err)
				}
				three.DrawString(yw.Vlist[k].Name+func(i int) (s string) { // 副词条名
					for p := 0; p < i; p++ {
						s += "↑"
					}
					return s
				}(yw.Vlist[k].Adds-1), xx, yy)
				if err := three.LoadFontFace(FiFile, 30); err != nil {
					panic(err)
				}
				three.DrawStringAnchored("+"+yw.Vlist[k].Value+stofen(yw.Vlist[k].Name), 325, yy, 1, 0)
			}
			x, y := i%3*350+20, i/3*360+920
			dc.DrawImage(yinsyw, x, y)
			dc.DrawImage(three.Image(), x, y)
		}
	}()
	if err := dc.LoadFontFace(BaFile, 30); err != nil {
		panic(err)
	}
	dc.DrawStringAnchored("Created By Zerobot-Plugin & Klala || Data From LuLuApi", 540, 1655, 0.5, 0.5)
	wg.Wait()
	err = dc.SavePNG("data/klala/user/cache/" + t.UID + t.RoleData[n].Name + ".png")
	if err != nil {
		return "", err
	}
	return "data/klala/user/cache/" + t.UID + t.RoleData[n].Name + ".png", nil
}

// Yinying 绘制阴影 圆角矩形
func Yinying(x int, y int, r float64, c color.Color) image.Image {
	ctx := gg.NewContext(x, y)
	ctx.SetColor(c)
	ctx.DrawRoundedRectangle(0, 0, float64(x), float64(y), r)
	ctx.Fill()
	return ctx.Image()
}

// Polygon 画多边形
func Polygon(n int) []gg.Point {
	result := make([]gg.Point, n)
	for i := 0; i < n; i++ {
		a := float64(i)*2*math.Pi/float64(n) - math.Pi/2
		result[i] = gg.Point{X: math.Cos(a), Y: math.Sin(a)}
	}
	return result
}

// drawStars 画星星
func drawStars(side, all string, num int) image.Image {
	dc := gg.NewContext(500, 80)
	n := 5
	points := Polygon(n)
	for x, i := 40, 0; i < num; x += 80 {
		dc.Push()
		// s := rand.Float64()*S/4 + S/4
		dc.Translate(float64(x), 45)
		//	dc.Rotate(rand.Float64() * 1.5 * math.Pi) //旋转
		dc.Scale(30, 30) // 大小
		for i := 0; i < n+1; i++ {
			index := (i * 2) % n
			p := points[index]
			dc.LineTo(p.X, p.Y)
		}
		dc.SetLineWidth(10)
		dc.SetHexColor(side) // 线
		dc.StrokePreserve()
		dc.SetHexColor(all)
		dc.Fill()
		dc.Pop()
		i++
	}
	return dc.Image()
}

// adjustOpacity 更改透明度
func adjustOpacity(m image.Image, percentage float64) image.Image {
	bounds := m.Bounds()
	dx, dy := bounds.Dx(), bounds.Dy()
	nimg := image.NewRGBA64(bounds)
	for i := 0; i < dx; i++ {
		for j := 0; j < dy; j++ {
			r, g, b, a := m.At(i, j).RGBA()
			opacity := uint16(float64(a) * percentage)
			r, g, b, a = nimg.ColorModel().Convert(color.NRGBA64{R: uint16(r), G: uint16(g), B: uint16(b), A: opacity}).RGBA()
			nimg.SetRGBA64(i, j, color.RGBA64{R: uint16(r), G: uint16(g), B: uint16(b), A: uint16(a)})
		}
	}
	return nimg
}
