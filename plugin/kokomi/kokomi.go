// Package kokomi  原神面板v2.1
package kokomi

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"

	"bytes"
	"image"

	"github.com/Coloured-glaze/gg"
	"github.com/FloatTech/floatbox/img/writer"
	"github.com/FloatTech/floatbox/web"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"

	//"github.com/FloatTech/zbputils/img"
	"github.com/nfnt/resize"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

const (
	url      = "https://enka.minigg.cn/u/%v/__data.json"
	edition  = "Created By ZeroBot-Plugin v1.6.1 & kokomi v2.1"
	tu       = "https://api.yimian.xyz/img?type=moe&size=1920x1080"
	NameFont = "plugin/kokomi/data/font/NZBZ.ttf"        // 名字字体
	FontFile = "plugin/kokomi/data/font/HYWH-65W.ttf"    // 汉字字体
	FiFile   = "plugin/kokomi/data/font/tttgbnumber.ttf" // 其余字体(数字英文)
)

func init() { // 主函数
	en := control.Register("kokomi", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Brief:            "原神相关功能",
		Help: "原神面板执行方法,第一次需要依次执行\n" +
			"- 绑定......(uid)\n" +
			"- 更新面板\n" +
			"- 全部面板\n" +
			"- XX面板",
	})
	en.OnSuffix("面板").SetBlock(true).Handle(func(ctx *zero.Ctx) {
		str := ctx.State["args"].(string) // 获取key
		var wifeid int64
		var allfen = 0.00
		qquid := ctx.Event.UserID
		// 获取uid
		uid := Getuid(qquid)
		// uid := 113781666 //测试用
		suid := strconv.Itoa(uid)
		if uid == 0 {
			ctx.SendChain(message.Text("-未绑定uid"))
			return
		}
		//############################################################判断数据更新,逻辑原因不能合并进switch
		if str == "更新" || str == "#更新" {
			es, err := web.GetData(fmt.Sprintf(url, uid)) // 网站返回结果
			if err != nil {
				ctx.SendChain(message.Text("网站获取信息失败", err))
				return
			}
			// 创建存储文件,路径plugin/kokomi/data/js
			file, _ := os.OpenFile("plugin/kokomi/data/js/"+suid+".kokomi", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
			_, _ = file.Write(es)
			ctx.SendChain(message.Text("-角色面板更新成功喵~"))
			file.Close()
			return
		}
		//##########################################################
		// 获取本地缓存数据
		txt, err := os.ReadFile("plugin/kokomi/data/js/" + suid + ".kokomi")
		if err != nil {
			ctx.SendChain(message.Text("-本地未找到账号信息, 请更新面板"))
			return
		}

		// 解析
		var alldata Data
		err = json.Unmarshal(txt, &alldata)
		if err != nil {
			ctx.SendChain(message.Text("出现错误捏：", err))
			return
		}
		switch str {
		case "全部", "全部角色", "#全部":
			if len(alldata.PlayerInfo.ShowAvatarInfoList) == 0 {
				ctx.SendChain(message.Text("-请在游戏中打开角色面板展示后再尝试"))
				return
			}
			var msg strings.Builder
			msg.WriteString("您的展示角色为:\n")
			for i := 0; i < len(alldata.PlayerInfo.ShowAvatarInfoList); i++ {
				mmm := Uidmap[int64(alldata.PlayerInfo.ShowAvatarInfoList[i].AvatarID)]
				msg.WriteString(mmm)
				if i < len(alldata.PlayerInfo.ShowAvatarInfoList)-1 {
					msg.WriteByte('\n')
				}
			}
			ctx.SendChain(message.Text(msg.String()))
			return
		default: // 角色名解析为id
			//排除#
			if str[0:1] == "#" {
				str = str[1:]
			}
			//匹配简称/外号
			str = FindName(str)
			var flag bool
			wifeid, flag = Namemap[str]
			if !flag {
				ctx.SendChain(message.Text("请输入角色全名"))
				return
			}
		}
		var t = -1
		// 匹配角色
		for i := 0; i < len(alldata.PlayerInfo.ShowAvatarInfoList); i++ {
			if wifeid == int64(alldata.PlayerInfo.ShowAvatarInfoList[i].AvatarID) {
				t = i
			}
		}
		if t == -1 { // 在返回数据中未找到想要的角色
			ctx.SendChain(message.Text("该角色未展示"))
			return
		}

		// 画图
		dc := gg.NewContext(1080, 2400) // 画布大小
		dc.SetHexColor("#98F5FF")
		dc.Clear() // 背景
		pro, flg := Promap[wifeid]
		if !flg {
			ctx.SendChain(message.Text("匹配角色元素失败"))
			return
		}
		beijing, err := gg.LoadImage("plugin/kokomi/data/pro/" + pro + ".jpg")
		if err != nil {
			ctx.SendChain(message.Text("获取背景失败", err))
			return
		}
		dc.Scale(5/3.0, 5/3.0)
		dc.DrawImage(beijing, -792, 0)
		dc.Scale(3/5.0, 3/5.0)
		dc.SetRGB(1, 1, 1) // 换白色
		// 角色立绘565*935
		lihui, err := gg.LoadImage("plugin/kokomi/data/character/" + str + "/imgs/splash.webp")
		if err != nil {
			ctx.SendChain(message.Text("获取立绘失败", err))
			return
		}
		//立绘参数
		sxx := lihui.Bounds().Size().X
		syy := lihui.Bounds().Size().Y
		rate := 800 / float64(syy)
		dc.Scale(rate, rate)
		dc.DrawImage(lihui, int(270-float64(sxx)/2), 90)
		dc.Scale(1/rate, 1/rate)
		//角色名字
		if err := dc.LoadFontFace(NameFont, 80); err != nil {
			panic(err)
		}
		namelen := utf8.RuneCountInString(str)
		dc.DrawString(str, float64(1050-namelen*90), float64(130))
		// 好感度,uid
		if err := dc.LoadFontFace(FontFile, 30); err != nil {
			panic(err)
		}
		// 版本号
		dc.DrawString(edition, 180, 2380)
		ming := len(alldata.AvatarInfoList[t].TalentIDList)
		dc.DrawString("好感度"+strconv.Itoa(alldata.AvatarInfoList[t].FetterInfo.ExpLevel), 0, 40)
		dc.DrawString(alldata.PlayerInfo.Nickname, 700, 40)
		if err := dc.LoadFontFace(FiFile, 30); err != nil {
			panic(err)
		}
		dc.DrawString("UID:"+suid+"---LV"+strconv.Itoa(alldata.PlayerInfo.ShowAvatarInfoList[t].Level)+"---"+strconv.Itoa(ming), 600, 180)
		if err := dc.LoadFontFace(FontFile, 30); err != nil {
			panic(err)
		}
		dc.DrawString("命", 976, 180)
		// 角色等级,命之座(合并上程序)
		//dc.DrawString("LV"+strconv.Itoa(alldata.PlayerInfo.ShowAvatarInfoList[t].Level), 630, 130) // 角色等级
		//dc.DrawString(strconv.Itoa(ming)+"命", 765, 130)

		//新建图层,实现阴影
		bg := Yinying(540, 470, 16)
		//字图层
		one := gg.NewContext(540, 470)
		if err := one.LoadFontFace(FontFile, 30); err != nil {
			panic(err)
		}
		// 属性540*460,字30,间距15,60
		one.SetRGB(1, 1, 1) //白色
		one.DrawString("生命值:", 70, 40)
		one.DrawString("攻击力:", 70, 100)
		one.DrawString("防御力:", 70, 160)
		one.DrawString("元素精通:", 70, 220)
		one.DrawString("暴击率:", 70, 280)
		one.DrawString("暴击伤害:", 70, 340)
		one.DrawString("元素充能:", 70, 400)
		// 元素加伤判断
		e1, e2 := 70, 460
		switch {
		case alldata.AvatarInfoList[t].FightPropMap.Num30*100 > 0:
			one.DrawString("物理加伤:", float64(e1), float64(e2))
		case alldata.AvatarInfoList[t].FightPropMap.Num40*100 > 0:
			one.DrawString("火元素加伤:", float64(e1), float64(e2))
		case alldata.AvatarInfoList[t].FightPropMap.Num41*100 > 0:
			one.DrawString("雷元素加伤:", float64(e1), float64(e2))
		case alldata.AvatarInfoList[t].FightPropMap.Num42*100 > 0:
			one.DrawString("水元素加伤:", float64(e1), float64(e2))
		case alldata.AvatarInfoList[t].FightPropMap.Num44*100 > 0:
			one.DrawString("风元素加伤:", float64(e1), float64(e2))
		case alldata.AvatarInfoList[t].FightPropMap.Num45*100 > 0:
			one.DrawString("岩元素加伤:", float64(e1), float64(e2))
		case alldata.AvatarInfoList[t].FightPropMap.Num46*100 > 0:
			one.DrawString("冰元素加伤:", float64(e1), float64(e2))
		default: //草或者无
			one.DrawString("元素加伤:", float64(e1), float64(e2))
		}

		//值,一一对应
		if err := one.LoadFontFace(FiFile, 30); err != nil {
			panic(err)
		}
		// 属性540*460,字30,间距15,60
		one.SetRGB(1, 1, 1)                                                                       //白色
		one.DrawString(Ftoone(alldata.AvatarInfoList[t].FightPropMap.Num2000), 335, 40)           //生命
		one.DrawString(Ftoone(alldata.AvatarInfoList[t].FightPropMap.Num2001), 335, 100)          //攻击
		one.DrawString(Ftoone(alldata.AvatarInfoList[t].FightPropMap.Num2002), 335, 160)          //防御
		one.DrawString(strconv.Itoa(int(alldata.AvatarInfoList[t].FightPropMap.Num28)), 335, 220) //精通
		one.DrawString(Ftoone(alldata.AvatarInfoList[t].FightPropMap.Num20*100)+"%", 335, 280)    //暴击
		one.DrawString(Ftoone(alldata.AvatarInfoList[t].FightPropMap.Num22*100)+"%", 335, 340)    //爆伤
		one.DrawString(Ftoone(alldata.AvatarInfoList[t].FightPropMap.Num23*100)+"%", 335, 400)    //充能
		// 元素加伤判断
		e1, e2 = 335, 460
		switch {
		case alldata.AvatarInfoList[t].FightPropMap.Num30*100 > 0: //物理
			one.DrawString(Ftoone(alldata.AvatarInfoList[t].FightPropMap.Num30*100)+"%", float64(e1), float64(e2))
		case alldata.AvatarInfoList[t].FightPropMap.Num40*100 > 0: //火
			one.DrawString(Ftoone(alldata.AvatarInfoList[t].FightPropMap.Num40*100)+"%", float64(e1), float64(e2))
		case alldata.AvatarInfoList[t].FightPropMap.Num41*100 > 0: //雷
			one.DrawString(Ftoone(alldata.AvatarInfoList[t].FightPropMap.Num41*100)+"%", float64(e1), float64(e2))
		case alldata.AvatarInfoList[t].FightPropMap.Num42*100 > 0: //水
			one.DrawString(Ftoone(alldata.AvatarInfoList[t].FightPropMap.Num42*100)+"%", float64(e1), float64(e2))
		case alldata.AvatarInfoList[t].FightPropMap.Num44*100 > 0: //风
			one.DrawString(Ftoone(alldata.AvatarInfoList[t].FightPropMap.Num44*100)+"%", float64(e1), float64(e2))
		case alldata.AvatarInfoList[t].FightPropMap.Num45*100 > 0: //岩
			one.DrawString(Ftoone(alldata.AvatarInfoList[t].FightPropMap.Num45*100)+"%", float64(e1), float64(e2))
		case alldata.AvatarInfoList[t].FightPropMap.Num46*100 > 0: //冰
			one.DrawString(Ftoone(alldata.AvatarInfoList[t].FightPropMap.Num46*100)+"%", float64(e1), float64(e2))
		default: //草或者无
			one.DrawString(Ftoone(alldata.AvatarInfoList[t].FightPropMap.Num43*100)+"%", float64(e1), float64(e2))
		}
		dc.DrawImage(bg, 505, 420)
		dc.DrawImage(one.Image(), 505, 420)

		// 天赋等级
		if err := dc.LoadFontFace(FiFile, 30); err != nil { // 字体大小
			panic(err)
		}
		talentid := IdtoTalent[wifeid]
		lin1 := alldata.AvatarInfoList[t].SkillLevelMap[talentid[0]]
		lin2 := alldata.AvatarInfoList[t].SkillLevelMap[talentid[1]]
		lin3 := alldata.AvatarInfoList[t].SkillLevelMap[talentid[2]]
		//v1版本dc.DrawString("天赋等级:"+strconv.Itoa(lin1)+"--"+strconv.Itoa(lin2)+"--"+strconv.Itoa(lin3), 630, 900)
		//贴图
		tulin1, err := gg.LoadImage("plugin/kokomi/data/character/" + str + "/icons/talent-a.webp")
		tulin1 = resize.Resize(80, 0, tulin1, resize.Bilinear)
		if err != nil {
			ctx.SendChain(message.Text("获取天赋图标失败", err))
			return
		}
		tulin2, err := gg.LoadImage("plugin/kokomi/data/character/" + str + "/icons/talent-e.webp")
		tulin2 = resize.Resize(80, 0, tulin2, resize.Bilinear)
		if err != nil {
			ctx.SendChain(message.Text("获取天赋图标失败", err))
			return
		}
		tulin3, err := gg.LoadImage("plugin/kokomi/data/character/" + str + "/icons/talent-q.webp")
		tulin3 = resize.Resize(80, 0, tulin3, resize.Bilinear)
		if err != nil {
			ctx.SendChain(message.Text("获取天赋图标失败", err))
			return
		}
		//边框间隔180
		kuang, err := gg.LoadPNG("plugin/kokomi/data/pro/" + pro + ".png")
		if err != nil {
			ctx.SendChain(message.Text("获取天赋边框失败", err))
			return
		}
		dc.DrawImage(kuang, 520, 220)
		dc.DrawImage(kuang, 700, 220)
		dc.DrawImage(kuang, 880, 220)

		//贴图间隔214
		dc.DrawImage(tulin1, 550, 260)
		//纠正素材问题
		bb := Tianfujiuzhen(str)
		dc.DrawImage(tulin2, 733, bb)
		dc.DrawImage(tulin3, 910, 260)

		//Lv间隔180
		dc.DrawString(strconv.Itoa(lin1), 580, 380)
		dc.DrawString(strconv.Itoa(lin2), 760, 380)
		dc.DrawString(strconv.Itoa(lin3), 940, 380)
		//皇冠
		tuguan, err := gg.LoadImage("plugin/kokomi/data/zawu/crown.png")
		if err != nil {
			ctx.SendChain(message.Text("获取皇冠图标失败", err))
			return
		}
		tuguan = resize.Resize(0, 55, tuguan, resize.Bilinear)
		if lin1 == 10 {
			dc.DrawImage(tuguan, 568, 215)
		}
		if lin2 == 10 {
			dc.DrawImage(tuguan, 748, 215)
		}
		if lin3 == 10 {
			dc.DrawImage(tuguan, 928, 215)
		}

		//武器图层
		//新建图层,实现阴影
		yinwq := Yinying(340, 180, 16)
		// 字图层
		two := gg.NewContext(340, 180)
		if err := two.LoadFontFace(FontFile, 30); err != nil {
			panic(err)
		}
		two.SetRGB(1, 1, 1) //白色
		//武器名
		//纠正圣遗物空缺报错的无返回情况
		if len(alldata.AvatarInfoList[t].EquipList) != 6 {
			ctx.SendChain(message.Text("角色圣遗物有空缺,暂时无法查看"))
			return
		}
		wq := IdforNamemap[alldata.AvatarInfoList[t].EquipList[5].Flat.NameTextHash]
		two.DrawString(wq, 150, 50)

		//详细
		two.DrawString("攻击力:", 150, 130)
		two.DrawString("精炼:", 245, 90)
		if err := two.LoadFontFace(FiFile, 30); err != nil { // 字体大小
			panic(err)
		}
		two.DrawString(strconv.FormatFloat(alldata.AvatarInfoList[t].EquipList[5].Flat.WeaponStat[0].Value, 'f', 1, 32), 250, 130)
		//Lv,精炼
		var wqjl int
		for m := range alldata.AvatarInfoList[t].EquipList[5].Weapon.AffixMap {
			wqjl = m
		}
		two.DrawString("Lv."+strconv.Itoa(alldata.AvatarInfoList[t].EquipList[5].Weapon.Level), 150, 90)
		two.DrawString(strconv.Itoa(alldata.AvatarInfoList[t].EquipList[5].Weapon.AffixMap[wqjl]+1), 316, 90)
		/*副词条,放不下
		fucitiao, _ := IdforNamemap[alldata.AvatarInfoList[t].EquipList[5].Flat.WeaponStat[1].SubPropId] //名称
		var baifen = "%"
		if fucitiao == "元素精通" {
			baifen = ""
		}
		dc.DrawString(fucitiao+":"+strconv.Itoa(int(alldata.AvatarInfoList[t].EquipList[5].Flat.WeaponStat[1].Value))+baifen, 820, 270)
		*/
		//图片
		tuwq, err := gg.LoadPNG("plugin/kokomi/data/wq/" + wq + ".png")
		if err != nil {
			ctx.SendChain(message.Text("获取武器图标失败", err))
			return
		}
		tuwq = resize.Resize(130, 0, tuwq, resize.Bilinear)
		two.DrawImage(tuwq, 10, 10)
		dc.DrawImage(yinwq, 20, 920)
		dc.DrawImage(two.Image(), 20, 920)

		//圣遗物
		//缩小
		yinsyw := Yinying(340, 350, 16)
		for i := 0; i < 5; i++ {
			// 字图层
			three := gg.NewContext(340, 350)
			if err := three.LoadFontFace(FontFile, 30); err != nil {
				panic(err)
			}
			//字号30,间距50
			three.SetRGB(1, 1, 1) //白色
			sywname := IdforNamemap[alldata.AvatarInfoList[t].EquipList[i].Flat.SetNameTextHash]
			tusyw, err := gg.LoadImage("plugin/kokomi/data/syw/" + sywname + "/" + strconv.Itoa(i+1) + ".webp")
			if err != nil {
				ctx.SendChain(message.Text("获取圣遗物图标失败", err))
				return
			}
			tusyw = resize.Resize(80, 0, tusyw, resize.Bilinear) //缩小
			three.DrawImage(tusyw, 15, 15)
			//圣遗物name
			sywallname := SywNamemap[sywname]
			three.DrawString(sywallname[i], 110, 50)
			//圣遗物属性
			zhuci := StoS(alldata.AvatarInfoList[t].EquipList[i].Flat.ReliquaryMainStat.MainPropID) //主词条
			zhucitiao := strconv.Itoa(int(alldata.AvatarInfoList[t].EquipList[i].Flat.ReliquaryMainStat.Value))
			//间隔45,初始145
			var xx, yy, pingfeng float64 //xx,yy词条相对位置,x,y文本框在全图位置
			var x, y int
			xx = 15
			yy = 145
			pingfeng = 0
			//主词条
			three.DrawString("主:"+zhuci, xx, yy)
			if err := three.LoadFontFace(FiFile, 30); err != nil {
				panic(err)
			} //主词条名字
			three.DrawString("+"+zhucitiao+Stofen(alldata.AvatarInfoList[t].EquipList[i].Flat.ReliquaryMainStat.MainPropID), 200, yy) //主词条属性
			//算分
			if i > 1 { //不算前两主词条属性
				pingfeng += Countcitiao(str, zhuci, alldata.AvatarInfoList[t].EquipList[i].Flat.ReliquaryMainStat.Value/4)
			}
			//副词条
			for k := 0; k < 4; k++ {
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
				three.DrawString(StoS(alldata.AvatarInfoList[t].EquipList[i].Flat.ReliquarySubStats[k].SubPropID), xx, yy)
				if err := three.LoadFontFace(FiFile, 30); err != nil {
					panic(err)
				}
				three.DrawString("+"+strconv.FormatFloat(alldata.AvatarInfoList[t].EquipList[i].Flat.ReliquarySubStats[k].Value, 'f', 1, 64)+Stofen(alldata.AvatarInfoList[t].EquipList[i].Flat.ReliquarySubStats[k].SubPropID), 200, yy)
				pingfeng += Countcitiao(str, StoS(alldata.AvatarInfoList[t].EquipList[i].Flat.ReliquarySubStats[k].SubPropID), alldata.AvatarInfoList[t].EquipList[i].Flat.ReliquarySubStats[k].Value)
			}
			//评分处理,对齐
			if i == 2 {
				pingfeng *= 0.90
			} else if i > 2 {
				pingfeng *= 0.85
			}
			allfen += pingfeng

			//圣遗物单个评分
			if err := three.LoadFontFace(FiFile, 30); err != nil {
				panic(err)
			}
			three.DrawString(Ftoone(pingfeng), 110, 85)
			three.DrawString("-"+Pingji(pingfeng), 222, 85) //评级
			if err := three.LoadFontFace(FontFile, 30); err != nil {
				panic(err)
			}
			three.DrawString("分", 175, 85)

			switch i {
			case 0:
				x = 370
				y = 920
			case 1:
				x = 720
				y = 920
			case 2:
				x = 20
				y = 1280
			case 3:
				x = 370
				y = 1280
			case 4:
				x = 720
				y = 1280
			}
			dc.DrawImage(yinsyw, x, y)
			dc.DrawImage(three.Image(), x, y)
		}

		//总评分框
		yinping := Yinying(340, 160, 16)
		// 字图层
		four := gg.NewContext(340, 160)
		if err := four.LoadFontFace(FontFile, 25); err != nil {
			panic(err)
		}
		four.SetRGB(1, 1, 1) //白色
		four.DrawString("评分规则:通用评分规则", 50, 35)

		if err := four.LoadFontFace(FiFile, 50); err != nil {
			panic(err)
		}
		four.DrawString(Ftoone(allfen), 50, 100)

		if err := four.LoadFontFace(FontFile, 25); err != nil {
			panic(err)
		}
		four.DrawString("圣遗物总分", 50, 150)
		dc.DrawImage(yinping, 20, 1110)
		dc.DrawImage(four.Image(), 20, 1110)

		//伤害显示区,暂时展示图片
		pic, err := web.GetData(tu)
		if err != nil {
			ctx.SendChain(message.Text("错误：获取插图失败", err))
			return
		}
		dst, _, err := image.Decode(bytes.NewReader(pic))
		if err != nil {
			ctx.SendChain(message.Text("错误：获取插图失败", err))
			return
		}
		sx := float64(1080) / float64(dst.Bounds().Size().X) // 计算缩放倍率（宽）
		dc.Scale(sx, sx)                                     // 使画笔按倍率缩放
		dc.DrawImage(dst, 0, int(1700*(1/sx)))               // 贴图（会受上述缩放倍率影响）
		dc.Scale(1/sx, 1/sx)
		// 输出图片
		ff, cl := writer.ToBytes(dc.Image())  // 图片放入缓存
		ctx.SendChain(message.ImageBytes(ff)) // 输出
		cl()
	})

	// 获取json,转移位置
	/*en.OnFullMatch("更新面板").SetBlock(true).Handle(func(ctx *zero.Ctx) {
		qquid := ctx.Event.UserID
		uid := Getuid(qquid)
		// uid := 113781666
		suid := strconv.Itoa(uid)
		ctx.SendChain(message.Text(uid))
		es, err := web.GetData(fmt.Sprintf(url, uid)) // 网站返回结果
		if err != nil {
			ctx.SendChain(message.Text("网站获取信息失败", err))
			return
		}
		// 创建存储文件,路径plugin/kokomi/data/js
		file, _ := os.OpenFile("plugin/kokomi/data/js/"+suid+".kokomi", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		_, _ = file.Write(es)
		ctx.SendChain(message.Text("喵~更新成功"))
		file.Close()
	})*/
	// 绑定uid
	en.OnRegex(`^(#)?绑定\s*(uid)?(\d+)?`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		suid := ctx.State["regex_matched"].([]string)[3] // 获取uid
		int64uid, err := strconv.ParseInt(suid, 10, 64)
		if suid == "" || int64uid < 100000000 || int64uid > 1000000000 || err != nil {
			ctx.SendChain(message.Text("-请输入正确的uid"))
			return
		}
		sqquid := strconv.Itoa(int(ctx.Event.UserID))
		file, _ := os.OpenFile("plugin/kokomi/data/uid/"+sqquid+".kokomi", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		_, _ = file.Write([]byte(suid))
		file.Close()
		ctx.SendChain(message.Text("-绑定uid" + suid + "成功喵~\n" + "-尝试获取角色面板信息~"))

		//更新面板程序
		es, err := web.GetData(fmt.Sprintf(url, suid)) // 网站返回结果
		if err != nil {
			ctx.SendChain(message.Text("-网站获取角色信息失败", err))
			return
		}
		// 创建存储文件,路径plugin/kokomi/data/js
		file1, _ := os.OpenFile("plugin/kokomi/data/js/"+suid+".kokomi", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		_, _ = file1.Write(es)
		ctx.SendChain(message.Text("-角色面板更新成功喵~"))
		file1.Close()
	})
}
