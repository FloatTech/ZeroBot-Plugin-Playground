package klala

import (
	"image"
	"image/color"
	"math"
	"strconv"

	"github.com/FloatTech/gg"
	img "github.com/FloatTech/imgfactory"
)

const (
	NameFont = "data/klala/kkk/font/NZBZ.ttf"                    //NameFont 名字字体
	FontFile = "data/klala/kkk/font/SourceHanMonoSC-HeavyIt.ttf" //FontFile  汉字字体
	FiFile   = "data/klala/kkk/font/tttgbnumber.ttf"             //FiFile  其余字体(数字英文)
	BaFile   = "data/klala/kkk/font/STLITI.TTF"                  //	BaFile   华文隶书版本版本号字体
)

func (t *thisdata) drawcard() (image.Image, error) {
	dc := gg.NewContext(1080, 1680)
	dc.SetRGB(1, 1, 1)
	if err := dc.LoadFontFace(FontFile, 40); err != nil {
		panic(err)
	}
	beijing, err := gg.LoadImage("data/klala/kkk/pro/冰.jpg")
	if err != nil {
		return nil, err
	}
	beijing = img.Size(beijing, 0, 1680).Image()
	dc.DrawImageAnchored(beijing, 540, 0, 0.5, 0)
	lihui, err := gg.LoadPNG("data/klala/kkk/lihui/" + strconv.Itoa(t.RoleData.ID) + ".png")
	if err != nil {
		return nil, err
	}
	lihui = img.Size(lihui, 0, 880).Image()
	sxx := lihui.Bounds().Size().X
	dc.DrawImage(lihui, int(300-float64(sxx)/2), 0)
	//昵称框图
	{
		zero := gg.NewContext(540, 200)
		zero.SetRGB(1, 1, 1) //白色
		//角色名字
		if err := zero.LoadFontFace(NameFont, 80); err != nil {
			panic(err)
		}
		zero.DrawStringAnchored(t.RoleData.Name, 505, 130, 1, 0)
		if err := zero.LoadFontFace(FontFile, 30); err != nil {
			panic(err)
		}
		zero.DrawStringAnchored("昵称:"+t.Nickname, 505, 40, 1, 0)
		if err := zero.LoadFontFace(FiFile, 30); err != nil {
			panic(err)
		}
		zero.DrawStringAnchored("UID:"+t.UID, 505, 180, 1, 0)
		newying := Yinying(540, 200, 16, color.NRGBA{R: 0, G: 0, B: 0, A: 106})
		dc.DrawImage(newying, 505, 20)
		dc.DrawImage(zero.Image(), 505, 20)
	}
	//属性列表
	yinyinBlack127 := color.NRGBA{R: 0, G: 0, B: 0, A: 127}
	{
		one := gg.NewContext(540, 470)
		if err := one.LoadFontFace(FontFile, 30); err != nil {
			panic(err)
		}
		one.SetRGB(1, 1, 1) //白色
		one.DrawString("角色等级:", 70, 45)
		one.DrawString("生命值:", 70, 96.25)
		one.DrawString("攻击力:", 70, 147.5)
		one.DrawString("防御力:", 70, 198.75)
		one.DrawString("速度:", 70, 250)
		one.DrawString("暴击率:", 70, 301.25)
		one.DrawString("暴击伤害:", 70, 352.5)
		one.DrawString("治疗加成:", 70, 403.75)
		one.DrawString("效果命中:", 70, 455)
		//值,一一对f应
		if err := one.LoadFontFace(FiFile, 30); err != nil {
			panic(err)
		}
		// 属性540*460,字30,间距15,60
		one.SetRGB(1, 1, 1)                                                                           //白色
		one.DrawStringAnchored("Lv"+strconv.Itoa(t.RoleData.List.Level), 470, 45, 1, 0)               //Lv
		one.DrawStringAnchored(Ftoone(t.RoleData.List.HpFinal), 470, 96.25, 1, 0)                     //生命
		one.DrawStringAnchored(Ftoone(t.RoleData.List.AttackFinal), 470, 147.5, 1, 0)                 //攻击
		one.DrawStringAnchored(Ftoone(t.RoleData.List.DefenseFinal), 470, 198.75, 1, 0)               //防御
		one.DrawStringAnchored(Ftoone(t.RoleData.List.SpeedFinal), 470, 250, 1, 0)                    //速度
		one.DrawStringAnchored(Ftoone(t.RoleData.List.CriticalChance*100)+"%", 470, 301.25, 1, 0)     //暴击
		one.DrawStringAnchored(Ftoone(t.RoleData.List.CriticalDamage*100)+"%", 470, 352.5, 1, 0)      //爆伤
		one.DrawStringAnchored(Ftoone(float64(t.RoleData.List.HealRatio*100))+"%", 470, 403.75, 1, 0) //治疗加成
		one.DrawStringAnchored(Ftoone(t.RoleData.List.StatusProbability*100)+"%", 470, 455, 1, 0)     //效果命中
		dc.DrawImage(Yinying(540, 470, 16, yinyinBlack127), 505, 240)                                 // 背景
		dc.DrawImage(one.Image(), 505, 240)
	}
	// 光锥
	{
		yinlight := Yinying(1040, 180, 16, yinyinBlack127)
		two := gg.NewContext(1040, 180)
		two.SetRGB(1, 1, 1) //白色
		if err := two.LoadFontFace(FontFile, 30); err != nil {
			panic(err)
		}
		//天赋
		two.DrawString("普攻 ", 40, 80)
		two.DrawString("战技 ", 240, 80)
		two.DrawString("终结 ", 40, 150)
		two.DrawString("天赋 ", 240, 150)
		//图标
		lpic, err := gg.LoadImage("data/klala/kkk/lights/" + strconv.Itoa(t.RoleData.Light.ID) + ".png")
		if err == nil {
			lpic = img.Size(lpic, 0, 140).Image()
			two.DrawImage(lpic, 700, 30)
		}
		two.DrawString(t.RoleData.Light.Name, 830, 60)
		//星级
		two.DrawImageAnchored(img.Size(Drawstars("#FFCC00", "#FFE43A", t.RoleData.Light.Star), 0, 30).Image(), 1020, 80, 1, 0)
		if err := two.LoadFontFace(FiFile, 30); err != nil {
			panic(err)
		}
		two.DrawString("LV"+strconv.Itoa(t.RoleData.Light.Level), 830, 150)
		two.DrawString("LV"+strconv.Itoa(t.RoleData.Skill.A), 110, 80)
		two.DrawString("LV"+strconv.Itoa(t.RoleData.Skill.E), 310, 80)
		two.DrawString("LV"+strconv.Itoa(t.RoleData.Skill.Q), 110, 150)
		two.DrawString("LV"+strconv.Itoa(t.RoleData.Skill.T), 310, 150)
		dc.DrawImage(yinlight, 20, 720)
		dc.DrawImage(two.Image(), 20, 720)
	}
	//遗物
	{
		yinsyw := Yinying(340, 350, 16, yinyinBlack127)
		var yw relicsdata
		for i := 0; i < 6; i++ {
			switch i {
			case 0:
				yw = t.RoleData.Relics.Head
			case 1:
				yw = t.RoleData.Relics.Hand
			case 2:
				yw = t.RoleData.Relics.Body
			case 3:
				yw = t.RoleData.Relics.Foot
			case 4:
				yw = t.RoleData.Relics.Neck
			default:
				yw = t.RoleData.Relics.Object
			}
			if yw.Name == "" {
				continue
			}
			// 字图层
			three := gg.NewContext(340, 350)
			if err := three.LoadFontFace(FontFile, 30); err != nil {
				panic(err)
			}
			//字号30,间距50
			three.SetRGB(1, 1, 1) //白色
			//画线
			for c := 0; c < 4; c++ {
				three.DrawLine(0, 157+float64(c)*45, 350, 157+float64(c)*45) //横线条分割
			}
			three.Stroke()
			/*tusyw, err := gg.LoadImage("data/klala/kkk/")
			if err == nil {
				tusyw = resize.Resize(80, 0, tusyw, resize.Bilinear) //缩小
				three.DrawImage(tusyw, 15, 15)
			}*/
			//遗物name
			three.DrawStringAnchored(yw.Name, 325, 50, 1, 0)
			//圣遗物属性 主词条
			//间隔45,初始145
			var xx, yy float64 //xx,yy词条相对位置,x,y文本框在全图位置
			var x, y int
			xx = 15
			yy = 145
			//主词条
			if err := three.LoadFontFace(FontFile, 30); err != nil {
				panic(err)
			}
			three.DrawString(yw.MainV.Name, xx, yy) //"主:"
			if err := three.LoadFontFace(FiFile, 30); err != nil {
				panic(err)
			}
			//主词条属性
			three.DrawStringAnchored("+"+yw.MainV.Value+stofen(yw.MainV.Name), 325, yy, 1, 0) //主词条属性
			three.SetHexColor("#98F5FF")                                                      //蓝色
			for k := 0; k < len(yw.Vlist); k++ {
				switch k {
				case 0:
					yy = 190
				case 1:
					yy = 235
				case 2:
					yy = 280
				case 3:
					yy = 325
				}
				if err := three.LoadFontFace(FontFile, 30); err != nil {
					panic(err)
				}
				three.DrawString(yw.Vlist[k].Name, xx, yy)
				if err := three.LoadFontFace(FiFile, 30); err != nil {
					panic(err)
				}
				three.DrawStringAnchored("+"+yw.Vlist[k].Value+stofen(yw.Vlist[k].Name), 325, yy, 1, 0)
			}
			switch i {
			case 0:
				x, y = 20, 920
			case 1:
				x, y = 370, 920
			case 2:
				x, y = 720, 920
			case 3:
				x, y = 20, 1280
			case 4:
				x, y = 370, 1280
			case 5:
				x, y = 720, 1280
			}
			dc.DrawImage(yinsyw, x, y)
			dc.DrawImage(three.Image(), x, y)
		}
	}
	if err := dc.LoadFontFace(BaFile, 30); err != nil {
		panic(err)
	}
	dc.DrawStringAnchored("Created By Klala", 540, 1655, 0.5, 0.5)
	return dc.Image(), nil
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

// Drawstars 画星星
func Drawstars(side, all string, num int) image.Image {
	dc := gg.NewContext(500, 80)
	n := 5
	points := Polygon(n)
	for x, i := 40, 0; i < num; x += 80 {
		dc.Push()
		//s := rand.Float64()*S/4 + S/4
		dc.Translate(float64(x), 45)
		//	dc.Rotate(rand.Float64() * 1.5 * math.Pi) //旋转
		dc.Scale(30, 30) //大小
		for i := 0; i < n+1; i++ {
			index := (i * 2) % n
			p := points[index]
			dc.LineTo(p.X, p.Y)
		}
		dc.SetLineWidth(10)
		dc.SetHexColor(side) //线
		dc.StrokePreserve()
		dc.SetHexColor(all)
		dc.Fill()
		dc.Pop()
		i++
	}
	return dc.Image()
}
