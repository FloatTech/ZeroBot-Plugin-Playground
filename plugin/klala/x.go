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
	affixMainFile   = "data/klala/kkk/json/RelicMainAffixConfig.json" //主词条属性
	affixFile       = "data/klala/kkk/json/RelicSubAffixConfig.json"  //副词条属性
	yiWuPath        = "data/klala/kkk/json/relics.json"               //遗物介绍
	lightJSONPath   = "data/klala/kkk/json/light_cones.json"          //光锥详情
	wifesPath       = "data/klala/kkk/json/nickname.json"             //别名
	relicConfigPath = "data/klala/kkk/json/RelicConfig.json"          //遗物对应属性
	uidPath         = "data/klala/kkk/uid/"
	wifeDataPath    = "data/klala/kkk/json/AvatarPromotionConfig.json" //角色基础属性
	lightsPath      = "data/klala/kkk/json/light_cone_promotions.json" //光锥属性
	lightAffixPath  = "data/klala/kkk/json/light_cone_ranks.json"      //光锥副词条
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
	es, err := json.Marshal(&t)
	if err != nil {
		return "", errors.New("ERROR: " + err.Error())
	}
	file, _ := os.OpenFile("data/klala/kkk/js/"+t.UID+".klala", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	_, _ = file.Write(es)
	file.Close()
	var msg strings.Builder
	msg.WriteString("-更新成功,您展示的角色为: ")
	for _, v := range t.RoleData {
		msg.WriteString("\n ")
		msg.WriteString(v.Name)
	}
	m = msg.String()
	return
}

func (r info) convertData() thisdata {
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
	t.UID = strconv.Itoa(r.PlayerDetailInfo.UID)
	t.Nickname = r.PlayerDetailInfo.NickName
	t.Level = r.PlayerDetailInfo.Level
	for k, v := range r.PlayerDetailInfo.DisplayAvatarList {
		t.RoleData = append(t.RoleData, ro{
			ID:   v.AvatarID,
			Star: v.EquipmentID.Rank + 4,
			Name: wife.idmap("wife", strconv.Itoa(v.AvatarID)),
			Rank: v.Rank,
		})
		//给基础值
		thisWifeData := wifeData[strconv.Itoa(v.AvatarID)][strconv.Itoa(v.Promotion)]
		t.RoleData[k].List = combat{
			AvatarID:          v.AvatarID,
			Level:             v.Level,
			Promotion:         v.Promotion,
			HpBase:            thisWifeData.HPBase.Value + thisWifeData.HPAdd.Value*float64(v.Level-1),
			HpFinal:           thisWifeData.HPBase.Value + thisWifeData.HPAdd.Value*float64(v.Level-1),
			AttackBase:        thisWifeData.AttackBase.Value + thisWifeData.AttackAdd.Value*float64(v.Level-1),
			AttackFinal:       thisWifeData.AttackBase.Value + thisWifeData.AttackAdd.Value*float64(v.Level-1),
			DefenseBase:       thisWifeData.DefenceBase.Value + thisWifeData.DefenceAdd.Value*float64(v.Level-1),
			DefenseFinal:      thisWifeData.DefenceBase.Value + thisWifeData.DefenceAdd.Value*float64(v.Level-1),
			SpeedBase:         thisWifeData.SpeedBase.Value,
			SpeedFinal:        float64(thisWifeData.SpeedBase.Value),
			CriticalChance:    thisWifeData.CriticalChance.Value,
			CriticalDamage:    thisWifeData.CriticalDamage.Value,
			StatusProbability: 0,
			StatusResistance:  0,
		}
		t.RoleData[k].Light = light{
			Name:      wife.idmap("light", strconv.Itoa(v.EquipmentID.ID)),
			ID:        v.EquipmentID.ID,
			Star:      lights[strconv.Itoa(v.EquipmentID.ID)].Rarity,
			Level:     v.EquipmentID.Level,
			Promotion: v.EquipmentID.Promotion,
			Rank:      v.EquipmentID.Rank,
		}
		w := &t.RoleData[k].List
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
		t.RoleData[k].Skill = skill{
			A: v.BehaviorList[0].Level,
			E: v.BehaviorList[1].Level,
			Q: v.BehaviorList[2].Level,
			T: v.BehaviorList[3].Level,
			F: v.BehaviorList[4].Level,
		}
		for i := 0; i < len(v.RelicList); i++ {
			affixID := strconv.Itoa(v.RelicList[i].ID - 10000)
			mainSetID := relicConfig[strconv.Itoa(v.RelicList[i].ID)].SetID
			mainData := affixMain[strconv.Itoa(relicConfig[strconv.Itoa(v.RelicList[i].ID)].MainAffixGroup)][strconv.Itoa(v.RelicList[i].MainAffixID)]
			na := typeMap[mainData.Property]
			//属性计算
			{
				w.addList(na, v.RelicList[i].Level*mainData.LevelAdd.Value+mainData.BaseValue.Value)
				for _, vv := range v.RelicList[i].RelicSubAffix {
					nnn := typeMap[affix[affixID[0:1]][strconv.Itoa(vv.SubAffixID)].Property]
					w.addList(nnn, float64(vv.Cnt)*affix[affixID[0:1]][strconv.Itoa(vv.SubAffixID)].BaseValue.Value+float64(vv.Step)*affix[affixID[0:1]][strconv.Itoa(vv.SubAffixID)].StepValue.Value)
				}
			}
			switch v.RelicList[i].Type {
			case 1:
				t.RoleData[k].Relics.Head = relicsdata{
					Name:  yi[strconv.Itoa(mainSetID)].Pieces.Head.Name,
					SetID: mainSetID,
					Type:  1,
					Star:  v.RelicList[i].ID/10000 - 1,
					Level: v.RelicList[i].Level,
				}
				t.RoleData[k].Relics.Head.MainV = vlist{
					Name:  na,
					Value: Ftoone((v.RelicList[i].Level*mainData.LevelAdd.Value + mainData.BaseValue.Value) * sto100(na)),
				}
				for _, vv := range v.RelicList[i].RelicSubAffix {
					na = typeMap[affix[affixID[0:1]][strconv.Itoa(vv.SubAffixID)].Property]
					t.RoleData[k].Relics.Head.Vlist = append(t.RoleData[k].Relics.Head.Vlist, vlist{
						Name:  na,
						Value: Ftoone((float64(vv.Cnt)*affix[affixID[0:1]][strconv.Itoa(vv.SubAffixID)].BaseValue.Value + float64(vv.Step)*affix[affixID[0:1]][strconv.Itoa(vv.SubAffixID)].StepValue.Value) * sto100(na)),
						Adds:  vv.Cnt,
					})
				}
			case 2:
				t.RoleData[k].Relics.Hand = relicsdata{
					Name:  yi[strconv.Itoa(mainSetID)].Pieces.Hands.Name,
					SetID: mainSetID,
					Type:  2,
					Star:  v.RelicList[i].ID/10000 - 1,
					Level: v.RelicList[i].Level,
				}
				t.RoleData[k].Relics.Hand.MainV = vlist{
					Name:  na,
					Value: Ftoone((v.RelicList[i].Level*mainData.LevelAdd.Value + mainData.BaseValue.Value) * sto100(na)),
				}
				for _, vv := range v.RelicList[i].RelicSubAffix {
					na = typeMap[affix[affixID[0:1]][strconv.Itoa(vv.SubAffixID)].Property]
					t.RoleData[k].Relics.Hand.Vlist = append(t.RoleData[k].Relics.Hand.Vlist, vlist{
						Name:  na,
						Value: Ftoone((float64(vv.Cnt)*affix[affixID[0:1]][strconv.Itoa(vv.SubAffixID)].BaseValue.Value + float64(vv.Step)*affix[affixID[0:1]][strconv.Itoa(vv.SubAffixID)].StepValue.Value) * sto100(na)),
						Adds:  vv.Cnt,
					})
				}
			case 3:
				t.RoleData[k].Relics.Body = relicsdata{
					Name:  yi[strconv.Itoa(mainSetID)].Pieces.Body.Name,
					SetID: mainSetID,
					Type:  3,
					Star:  v.RelicList[i].ID/10000 - 1,
					Level: v.RelicList[i].Level,
				}
				t.RoleData[k].Relics.Body.MainV = vlist{
					Name:  na,
					Value: Ftoone((v.RelicList[i].Level*mainData.LevelAdd.Value + mainData.BaseValue.Value) * sto100(na)),
				}
				for _, vv := range v.RelicList[i].RelicSubAffix {
					na = typeMap[affix[affixID[0:1]][strconv.Itoa(vv.SubAffixID)].Property]
					t.RoleData[k].Relics.Body.Vlist = append(t.RoleData[k].Relics.Body.Vlist, vlist{
						Name:  na,
						Value: Ftoone((float64(vv.Cnt)*affix[affixID[0:1]][strconv.Itoa(vv.SubAffixID)].BaseValue.Value + float64(vv.Step)*affix[affixID[0:1]][strconv.Itoa(vv.SubAffixID)].StepValue.Value) * sto100(na)),
						Adds:  vv.Cnt,
					})
				}
			case 4:
				t.RoleData[k].Relics.Foot = relicsdata{
					Name:  yi[strconv.Itoa(mainSetID)].Pieces.Feet.Name,
					SetID: mainSetID,
					Type:  4,
					Star:  v.RelicList[i].ID/10000 - 1,
					Level: v.RelicList[i].Level,
				}
				t.RoleData[k].Relics.Foot.MainV = vlist{
					Name:  na,
					Value: Ftoone((v.RelicList[i].Level*mainData.LevelAdd.Value + mainData.BaseValue.Value) * sto100(na)),
				}
				for _, vv := range v.RelicList[i].RelicSubAffix {
					na = typeMap[affix[affixID[0:1]][strconv.Itoa(vv.SubAffixID)].Property]
					t.RoleData[k].Relics.Foot.Vlist = append(t.RoleData[k].Relics.Foot.Vlist, vlist{
						Name:  na,
						Value: Ftoone((float64(vv.Cnt)*affix[affixID[0:1]][strconv.Itoa(vv.SubAffixID)].BaseValue.Value + float64(vv.Step)*affix[affixID[0:1]][strconv.Itoa(vv.SubAffixID)].StepValue.Value) * sto100(na)),
						Adds:  vv.Cnt,
					})
				}
			case 5:
				t.RoleData[k].Relics.Neck = relicsdata{
					Name:  yi[affixID[1:4]].Pieces.PlanarSphere.Name,
					SetID: relicConfig[strconv.Itoa(v.RelicList[i].ID)].SetID,
					Type:  5,
					Star:  v.RelicList[i].ID/10000 - 1,
					Level: v.RelicList[i].Level,
				}
				t.RoleData[k].Relics.Neck.MainV = vlist{
					Name:  na,
					Value: Ftoone((v.RelicList[i].Level*mainData.LevelAdd.Value + mainData.BaseValue.Value) * sto100(na)),
				}
				for _, vv := range v.RelicList[i].RelicSubAffix {
					na = typeMap[affix[affixID[0:1]][strconv.Itoa(vv.SubAffixID)].Property]
					t.RoleData[k].Relics.Neck.Vlist = append(t.RoleData[k].Relics.Neck.Vlist, vlist{
						Name:  na,
						Value: Ftoone((float64(vv.Cnt)*affix[affixID[0:1]][strconv.Itoa(vv.SubAffixID)].BaseValue.Value + float64(vv.Step)*affix[affixID[0:1]][strconv.Itoa(vv.SubAffixID)].StepValue.Value) * sto100(na)),
						Adds:  vv.Cnt,
					})
				}
			case 6:
				t.RoleData[k].Relics.Object = relicsdata{
					Name:  yi[affixID[1:4]].Pieces.LinkRope.Name,
					SetID: relicConfig[strconv.Itoa(v.RelicList[i].ID)].SetID,
					Type:  6,
					Star:  v.RelicList[i].ID/10000 - 1,
					Level: v.RelicList[i].Level,
				}
				t.RoleData[k].Relics.Object.MainV = vlist{
					Name:  na,
					Value: Ftoone((v.RelicList[i].Level*mainData.LevelAdd.Value + mainData.BaseValue.Value) * sto100(na)),
				}
				for _, vv := range v.RelicList[i].RelicSubAffix {
					na = typeMap[affix[affixID[0:1]][strconv.Itoa(vv.SubAffixID)].Property]
					t.RoleData[k].Relics.Object.Vlist = append(t.RoleData[k].Relics.Object.Vlist, vlist{
						Name:  na,
						Value: Ftoone((float64(vv.Cnt)*affix[affixID[0:1]][strconv.Itoa(vv.SubAffixID)].BaseValue.Value + float64(vv.Step)*affix[affixID[0:1]][strconv.Itoa(vv.SubAffixID)].StepValue.Value) * sto100(na)),
						Adds:  vv.Cnt,
					})
				}
			}
		}

	}

	return *t
}
func downdata(ctx *zero.Ctx) bool {
	if file.IsNotExist("data/klala/kkk") {
		ctx.SendChain(message.Text("-开始下载资源文件..."))
		cmd := exec.Command("git", "clone", "https://gitee.com/lianhong2758/kkk.git")
		cmd.Dir = file.BOTPATH + "/data/klala"
		_, err := cmd.CombinedOutput()
		if err != nil {
			return false
		}
		ctx.SendChain(message.Text("-下载资源文件成功..."))
	}
	if file.IsNotExist("data/klala/kkk/uid") {
		err := os.MkdirAll("data/klala/kkk/uid", 0755) // 递归创建文件夹
		if err != nil {
			return false
		}
	}
	if file.IsNotExist("data/klala/kkk/js") {
		err := os.MkdirAll("data/klala/kkk/js", 0755) // 递归创建文件夹
		if err != nil {
			return false
		}
	}
	return true
}