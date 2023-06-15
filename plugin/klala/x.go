package klala

import (
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/FloatTech/floatbox/file"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

const (
	jsPath          = "data/klala/user/js/"
	uidPath         = "data/klala/user/uid/"
	affixMainFile   = "data/klala/kkk/json/RelicMainAffixConfig.json"  //主词条属性
	affixFile       = "data/klala/kkk/json/RelicSubAffixConfig.json"   //副词条属性
	relicConfigPath = "data/klala/kkk/json/RelicConfig.json"           //遗物对应属性
	ywSetPath       = "data/klala/kkk/json/relic_sets.json"            //遗物Set属性
	yiWuPath        = "data/klala/kkk/json/relics.json"                //遗物介绍
	wifesPath       = "data/klala/kkk/json/nickname.json"              //别名
	wifeDataPath    = "data/klala/kkk/json/character_promotions.json"  //角色基础属性
	wifeTreePath    = "data/klala/kkk/json/character_skill_trees.json" //角色行迹属性
	wifeIntrodPath  = "data/klala/kkk/json/characters.json"            //角色介绍
	lightsPath      = "data/klala/kkk/json/light_cone_promotions.json" //光锥属性
	lightAffixPath  = "data/klala/kkk/json/light_cone_ranks.json"      //光锥副词条
	lightJSONPath   = "data/klala/kkk/json/light_cones.json"           //光锥详情
)

func getuid(sqquid string) (uid int) { // 获取对应游戏uid
	// 获取本地缓存数据
	txt, err := os.ReadFile(uidPath + sqquid + ".klala")
	if err != nil {
		return 0
	}
	uid, _ = strconv.Atoi(string(txt))
	return
}

func getWifeOrWq() (m FindMap) {
	txt, _ := os.ReadFile(wifesPath)
	_ = json.Unmarshal(txt, &m)
	return
}

// Findnames 遍历寻找匹配昵称
func (m FindMap) findnames(typ, val string) string {
	if typ == "wife" {
		for k, v := range m.Characters {
			for _, vv := range v {
				if vv == val {
					return k
				}
			}
		}
	} else {
		for k, v := range m.LightCones {
			for _, vv := range v {
				if vv == val {
					return k
				}
			}
		}
	}
	return ""
}

// Idmap wifeid->wifename
func (m FindMap) idmap(typ, val string) string {
	if typ == "wife" {
		f, b := m.Characters[val]
		if !b {
			return ""
		}
		return f[0]
	}
	f, b := m.LightCones[val]
	if !b {
		return ""
	}
	return f[0]
}

func getLights() (m lightmap) {
	txt, _ := os.ReadFile(lightJSONPath)
	_ = json.Unmarshal(txt, &m)
	return
}

func getYiWu() (m yiwumap) {
	txt, _ := os.ReadFile(yiWuPath)
	_ = json.Unmarshal(txt, &m)
	return
}

func getAffix() (m affixStarMap) {
	txt, _ := os.ReadFile(affixFile)
	_ = json.Unmarshal(txt, &m)
	return
}

func getAffixMain() (m affixStarMainMap) {
	txt, _ := os.ReadFile(affixMainFile)
	_ = json.Unmarshal(txt, &m)
	return
}

func getRelicConfig() (m relicConfigMap) {
	txt, _ := os.ReadFile(relicConfigPath)
	_ = json.Unmarshal(txt, &m)
	return
}

func getWifeData() (m wifeData) {
	txt, _ := os.ReadFile(wifeDataPath)
	_ = json.Unmarshal(txt, &m)
	return
}

func getWifeIntrod() (m wifeIntrod) {
	txt, _ := os.ReadFile(wifeIntrodPath)
	_ = json.Unmarshal(txt, &m)
	return
}

func getLightsData() (m lightData) {
	txt, _ := os.ReadFile(lightsPath)
	_ = json.Unmarshal(txt, &m)
	return
}

func getLightAffix() (m lightAffix) {
	txt, _ := os.ReadFile(lightAffixPath)
	_ = json.Unmarshal(txt, &m)
	return
}

func getYiwuSet() (m ywSetData) {
	txt, _ := os.ReadFile(ywSetPath)
	_ = json.Unmarshal(txt, &m)
	return
}

func getWifeTree() (m wifeTrees) {
	txt, _ := os.ReadFile(wifeTreePath)
	_ = json.Unmarshal(txt, &m)
	return
}

// Ftoone 保留一位小数并转化string
func Ftoone(f float64) string {
	// return strconv.FormatFloat(f, 'f', 1, 64)
	if f == 0 {
		return "0"
	}
	return strconv.FormatFloat(f, 'f', 1, 64)
}

// Stofen 判断词条分号
func stofen(val string) string {
	switch val {
	case "攻击力", "防御力", "生命值", "速度":
		return ""
	}
	return "%"
}

func sto100(val string) float64 {
	switch val {
	case "攻击力", "防御力", "生命值", "速度":
		return 1
	}
	return 100
}

// 计算总属性
func (w *combat) addList(str string, val float64) {
	switch str {
	case "生命值":
		w.HpFinal += val
	case "大生命":
		w.HpFinal += val * w.HpBase
	case "攻击力":
		w.AttackFinal += val
	case "大攻击":
		w.AttackFinal += val * w.AttackBase
	case "防御力":
		w.DefenseFinal += val
	case "大防御":
		w.DefenseFinal += val * w.DefenseBase
	case "速度":
		w.SpeedFinal += val
	case "效果命中":
		w.StatusProbability += val
	case "效果抵抗":
		w.StatusResistance += val
	case "暴击率":
		w.CriticalChance += val
	case "暴击伤害":
		w.CriticalDamage += val
	default:
	}
}

// ywtz 遗物套装判断
func ywtz(syws []int) map[int]int {
	ywMap := make(map[int]int)
	for _, v := range syws {
		i := ywMap[v]
		ywMap[v] = i + 1
	}
	ywMap[0] = 0
	return ywMap
}
func saveRoel(uid string) (m string, err error) {
	data, err := getRole(uid)
	if err != nil {
		return "", err
	}
	var r info
	err = json.Unmarshal(data, &r)
	if err != nil {
		return "", errors.New("ERROR: " + err.Error())
	}
	// 映射本地结构
	t := r.convertData()
	if len(t.RoleData) < 1 {
		return "", errors.New("ERROR: 展柜无展示角色")
	}
	es, err := json.Marshal(t)
	if err != nil {
		return "", errors.New("ERROR: " + err.Error())
	}
	file, _ := os.OpenFile(jsPath+t.UID+".klala", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	_, _ = file.Write(es)
	file.Close()
	var msg strings.Builder
	msg.WriteString("-更新成功,您展示的角色为: ")
	for _, v := range t.RoleData {
		if v.Name != "" {
			msg.WriteString("\n ")
			msg.WriteString(v.Name)
		}
	}
	m = msg.String()
	return
}

func (r *info) mergeRole() {
	for _, v := range r.PlayerDetailInfo.DisplayAvatarList {
		if v.AvatarID == r.PlayerDetailInfo.AssistAvatar.AvatarID {
			return
		}
	}
	//未找到相同
	r.PlayerDetailInfo.DisplayAvatarList = append(r.PlayerDetailInfo.DisplayAvatarList, r.PlayerDetailInfo.AssistAvatar)
}
func (r *info) convertData() *thisdata {
	t := new(thisdata)
	wife := getWifeOrWq()
	lights := getLights()
	yi := getYiWu()
	affix := getAffix()
	affixMain := getAffixMain()
	relicConfig := getRelicConfig()
	wifeData := getWifeData()
	lightsData := getLightsData()
	lightAffix := getLightAffix()
	ywSetData := getYiwuSet()
	wifeTree := getWifeTree()
	wifeIntrods := getWifeIntrod()
	t.UID = strconv.Itoa(r.PlayerDetailInfo.UID)
	t.Nickname = r.PlayerDetailInfo.NickName
	t.Level = r.PlayerDetailInfo.Level
	//合并助战角色
	r.mergeRole()
	for k, v := range r.PlayerDetailInfo.DisplayAvatarList {
		ywtzs := []int{}
		introd := wifeIntrods[strconv.Itoa(v.AvatarID)]
		t.RoleData = append(t.RoleData, ro{
			ID:      v.AvatarID,
			Star:    v.EquipmentID.Rank + 4,
			Name:    wife.idmap("wife", strconv.Itoa(v.AvatarID)),
			Rank:    v.Rank,
			Path:    introd.Path,
			Element: introd.Element,
		})
		//给基础值
		thisWifeData := wifeData[strconv.Itoa(v.AvatarID)].Values[v.Promotion]
		t.RoleData[k].List = combat{
			AvatarID:          v.AvatarID,
			Level:             v.Level,
			Promotion:         v.Promotion,
			HpBase:            thisWifeData.Hp.Base + thisWifeData.Hp.Step*float64(v.Level-1),
			HpFinal:           thisWifeData.Hp.Base + thisWifeData.Hp.Step*float64(v.Level-1),
			AttackBase:        thisWifeData.Atk.Base + thisWifeData.Atk.Step*float64(v.Level-1),
			AttackFinal:       thisWifeData.Atk.Base + thisWifeData.Atk.Step*float64(v.Level-1),
			DefenseBase:       thisWifeData.Def.Base + thisWifeData.Def.Step*float64(v.Level-1),
			DefenseFinal:      thisWifeData.Def.Base + thisWifeData.Def.Step*float64(v.Level-1),
			SpeedBase:         int(thisWifeData.Spd.Base),
			SpeedFinal:        thisWifeData.Spd.Base,
			CriticalChance:    thisWifeData.CritRate.Base,
			CriticalDamage:    thisWifeData.CritDmg.Base,
			StatusProbability: 0,
			StatusResistance:  0,
		}
		w := &t.RoleData[k].List
		if v.EquipmentID.ID != 0 {
			t.RoleData[k].Light = light{
				Name:      wife.idmap("light", strconv.Itoa(v.EquipmentID.ID)),
				ID:        v.EquipmentID.ID,
				Star:      lights[strconv.Itoa(v.EquipmentID.ID)].Rarity,
				Level:     v.EquipmentID.Level,
				Promotion: v.EquipmentID.Promotion,
				Rank:      v.EquipmentID.Rank,
			}
			lD := lightsData[strconv.Itoa(v.EquipmentID.ID)].Values[v.EquipmentID.Promotion]
			{
				//光锥基础属性
				w.HpFinal += lD.Hp.Base + lD.Hp.Step*float64(v.EquipmentID.Level-1)
				w.AttackFinal += lD.Atk.Base + lD.Atk.Step*float64(v.EquipmentID.Level-1)
				w.DefenseFinal += lD.Def.Base + lD.Def.Step*float64(v.EquipmentID.Level-1)
				w.HpBase += lD.Hp.Base + lD.Hp.Step*float64(v.EquipmentID.Level-1)
				w.AttackBase += lD.Atk.Base + lD.Atk.Step*float64(v.EquipmentID.Level-1)
				w.DefenseBase += lD.Def.Base + lD.Def.Step*float64(v.EquipmentID.Level-1)
				//副词条
				if b := lightAffix[strconv.Itoa(v.EquipmentID.ID)].Properties; len(b) > 0 {
					for _, bb := range b[v.EquipmentID.Rank-1] {
						w.addList(typeMap[bb.Type], bb.Value)
					}
				}
			}
		}
		t.RoleData[k].Skill = skill{
			A: v.BehaviorList[0].Level,
			E: v.BehaviorList[1].Level,
			Q: v.BehaviorList[2].Level,
			T: v.BehaviorList[3].Level,
			F: v.BehaviorList[4].Level,
		}
		//遗迹属性加成
		for _, vv := range v.BehaviorList {
			if vv.BehaviorID%1000 > 200 {
				for _, vvv := range wifeTree[strconv.Itoa(vv.BehaviorID)].Levels[0].Properties {
					w.addList(typeMap[vvv.Type], vvv.Value)
				}
			}
		}
		for i := 0; i < len(v.RelicList); i++ {
			affixID := strconv.Itoa(v.RelicList[i].ID - 10000)
			mainSetID := relicConfig[strconv.Itoa(v.RelicList[i].ID)].SetID
			mainData := affixMain[strconv.Itoa(relicConfig[strconv.Itoa(v.RelicList[i].ID)].MainAffixGroup)][strconv.Itoa(v.RelicList[i].MainAffixID)]
			na := typeMap[mainData.Property]
			//遗物套装加成
			ywtzs = append(ywtzs, mainSetID)
			//属性计算
			{
				w.addList(na, v.RelicList[i].Level*mainData.LevelAdd.Value+mainData.BaseValue.Value)
				for _, vv := range v.RelicList[i].RelicSubAffix {
					nnn := typeMap[affix[affixID[0:1]][strconv.Itoa(vv.SubAffixID)].Property]
					w.addList(nnn, float64(vv.Cnt)*affix[affixID[0:1]][strconv.Itoa(vv.SubAffixID)].BaseValue.Value+float64(vv.Step)*affix[affixID[0:1]][strconv.Itoa(vv.SubAffixID)].StepValue.Value)
				}
			}
			tRelicsdata := relicsdata{
				SetID: mainSetID,
				Type:  v.RelicList[i].Type,
				Star:  v.RelicList[i].ID/10000 - 1,
				Level: v.RelicList[i].Level,
			}
			tVlist := vlist{
				Name:  na,
				Value: Ftoone((v.RelicList[i].Level*mainData.LevelAdd.Value + mainData.BaseValue.Value) * sto100(na)),
			}
			var tAffixVlist = []vlist{}
			for _, vv := range v.RelicList[i].RelicSubAffix {
				nb := typeMap[affix[affixID[0:1]][strconv.Itoa(vv.SubAffixID)].Property]
				tAffixVlist = append(tAffixVlist, vlist{
					Name:  nb,
					Value: Ftoone((float64(vv.Cnt)*affix[affixID[0:1]][strconv.Itoa(vv.SubAffixID)].BaseValue.Value + float64(vv.Step)*affix[affixID[0:1]][strconv.Itoa(vv.SubAffixID)].StepValue.Value) * sto100(nb)),
					Adds:  vv.Cnt,
				})
			}
			switch v.RelicList[i].Type {
			case 1:
				t.RoleData[k].Relics.Head = tRelicsdata
				t.RoleData[k].Relics.Head.Name = yi[strconv.Itoa(mainSetID)].Pieces.Head.Name
				t.RoleData[k].Relics.Head.MainV = tVlist
				t.RoleData[k].Relics.Head.Vlist = append(t.RoleData[k].Relics.Head.Vlist, tAffixVlist...)
			case 2:
				t.RoleData[k].Relics.Hand = tRelicsdata
				t.RoleData[k].Relics.Hand.Name = yi[strconv.Itoa(mainSetID)].Pieces.Hands.Name
				t.RoleData[k].Relics.Hand.MainV = tVlist
				t.RoleData[k].Relics.Hand.Vlist = append(t.RoleData[k].Relics.Hand.Vlist, tAffixVlist...)
			case 3:
				t.RoleData[k].Relics.Body = tRelicsdata
				t.RoleData[k].Relics.Body.Name = yi[strconv.Itoa(mainSetID)].Pieces.Body.Name
				t.RoleData[k].Relics.Body.MainV = tVlist
				t.RoleData[k].Relics.Body.Vlist = append(t.RoleData[k].Relics.Body.Vlist, tAffixVlist...)
			case 4:
				t.RoleData[k].Relics.Foot = tRelicsdata
				t.RoleData[k].Relics.Foot.Name = yi[strconv.Itoa(mainSetID)].Pieces.Feet.Name
				t.RoleData[k].Relics.Foot.MainV = tVlist
				t.RoleData[k].Relics.Foot.Vlist = append(t.RoleData[k].Relics.Foot.Vlist, tAffixVlist...)
			case 5:
				t.RoleData[k].Relics.Neck = tRelicsdata
				t.RoleData[k].Relics.Neck.Name = yi[strconv.Itoa(mainSetID)].Pieces.PlanarSphere.Name
				t.RoleData[k].Relics.Neck.MainV = tVlist
				t.RoleData[k].Relics.Neck.Vlist = append(t.RoleData[k].Relics.Neck.Vlist, tAffixVlist...)
			case 6:
				t.RoleData[k].Relics.Object = tRelicsdata
				t.RoleData[k].Relics.Object.Name = yi[strconv.Itoa(mainSetID)].Pieces.LinkRope.Name
				t.RoleData[k].Relics.Object.MainV = tVlist
				t.RoleData[k].Relics.Object.Vlist = append(t.RoleData[k].Relics.Object.Vlist, tAffixVlist...)
			}
		}
		//套装属性
		{
			for kk, vv := range ywtz(ywtzs) {
				if vv > 1 {
					for _, vvv := range ywSetData[strconv.Itoa(kk)].Properties {
						for _, vvvv := range vvv {
							w.addList(typeMap[vvvv.Type], vvvv.Value)
						}
						if vv < 3 {
							break
						}
					}
				}
			}
		}

	}

	return t
}
func downdata(ctx *zero.Ctx) bool {
	if file.IsNotExist("data/klala/kkk") {
		ctx.SendChain(message.Text("-开始下载资源文件..."))
		cmd := exec.Command("git", "clone", "https://gitee.com/lianhong2758/kkk.git")
		cmd.Dir = file.BOTPATH + "/data/klala"
		_, err := cmd.CombinedOutput()
		if err != nil {
			ctx.SendChain(message.Text("-下载资源文件失败...", err))
			return false
		}
		ctx.SendChain(message.Text("-下载资源文件成功..."))
	}
	if file.IsNotExist("data/klala/user/uid") {
		err := os.MkdirAll("data/klala/user/uid", 0755) // 递归创建文件夹
		if err != nil {
			ctx.SendChain(message.Text("-ERROR: ", err))
			return false
		}
	}
	if file.IsNotExist("data/klala/user/js") {
		err := os.MkdirAll("data/klala/user/js", 0755) // 递归创建文件夹
		if err != nil {
			ctx.SendChain(message.Text("-ERROR: ", err))
			return false
		}
	}
	if file.IsNotExist("data/klala/user/cache") {
		err := os.MkdirAll("data/klala/user/cache", 0755) // 递归创建文件夹
		if err != nil {
			ctx.SendChain(message.Text("-ERROR: ", err))
			return false
		}
	}
	return true
}
