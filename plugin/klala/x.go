package klala

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"

	"github.com/FloatTech/floatbox/file"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

func Getuid(sqquid string) (uid int) { // 获取对应游戏uid
	// 获取本地缓存数据
	txt, err := os.ReadFile("data/klala/kkk/uid/" + sqquid + ".klala")
	if err != nil {
		return 0
	}
	uid, _ = strconv.Atoi(string(txt))
	return
}

// FindMap 各种简称map查询
type FindMap map[string][]string

func GetWifeOrWq(val string) FindMap {
	var txt []byte
	switch val {
	case "wife":
		txt, _ = os.ReadFile("data/klala/kkk/json/wife.json")
	}
	var m FindMap = make(map[string][]string)
	if nil == json.Unmarshal(txt, &m) {
		return m
	}
	return nil
}

// Findnames 遍历寻找匹配昵称
func (m FindMap) Findnames(val string) string {
	for k, v := range m {
		for _, vv := range v {
			if vv == val {
				return k
			}
		}
	}
	return ""
}

// Idmap wifeid->wifename
func (m FindMap) Idmap(val string) string {
	f, b := m[val]
	if !b {
		return ""
	}
	return f[0]
}

// 下载立绘文件
func downcard(url, id string) error {
	if file.IsExist("data/klala/kkk/lihui/" + id + ".png") {
		return nil
	}
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 创建文件
	f, err := os.Create(id + ".png")
	if err != nil {
		return err
	}
	defer f.Close()

	// 将响应内容写入文件
	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

// Ftoone 保留一位小数并转化string
func Ftoone(f float64) string {
	// return strconv.FormatFloat(f, 'f', 1, 64)
	if f == 0 {
		return "0"
	}
	return strconv.FormatFloat(f, 'f', 1, 64)
}
func typeToName(i int) string {
	switch i {
	case 0:
		return ""
	case 27:
		return "生命值"
	case 32:
		return "大生命"
	case 29:
		return "攻击力"
	case 33:
		return "大攻击"
	case 51:
		return "速度"
	case 59:
		return "击破特攻"
	case 34:
		return "大防御"
	case 31:
		return "防御力"
	case 57:
		return "效果抵抗"
	case 56:
		return "效果命中"
	case 52:
		return "暴击率"
	case 53:
		return "暴击伤害"
	case 12:
		return "物理伤害"
	case 55:
		return "治疗加成"
	}
	return "属性加伤"
}

// Stofen 判断词条分号
func Stofen(val string) string {
	switch val {
	case "攻击力", "防御力", "生命值", "速度":
		return ""
	}
	return "%"
}

func Sto100(val string) float64 {
	switch val {
	case "攻击力", "防御力", "生命值", "速度":
		return 1
	}
	return 100
}

func (r roles) convertData(nickname string, v characters) (thisdata, error) {
	t := new(thisdata)
	t.UID = strconv.Itoa(r.Data.AvatarCombat.UID)
	t.Nickname = nickname
	t.RoleData = ro{
		ID:      v.ID,
		Name:    v.Name,
		Star:    v.Star,
		Type:    v.Type,
		Element: v.Element,
	}
	t.RoleData.List = combat{
		AvatarID:          r.Data.AvatarCombat.AvatarID,
		Level:             r.Data.AvatarCombat.Level,
		Promotion:         r.Data.AvatarCombat.Promotion,
		HpBase:            r.Data.AvatarCombat.HpBase,
		HpFinal:           r.Data.AvatarCombat.HpFinal,
		AttackBase:        r.Data.AvatarCombat.AttackBase,
		AttackFinal:       r.Data.AvatarCombat.AttackFinal,
		DefenseBase:       r.Data.AvatarCombat.DefenseBase,
		DefenseFinal:      r.Data.AvatarCombat.DefenseFinal,
		SpeedBase:         r.Data.AvatarCombat.SpeedBase,
		SpeedFinal:        r.Data.AvatarCombat.SpeedFinal,
		CriticalChance:    r.Data.AvatarCombat.CriticalChance,
		CriticalDamage:    r.Data.AvatarCombat.CriticalDamage,
		HealRatio:         r.Data.AvatarCombat.HealRatio,
		StatusProbability: r.Data.AvatarCombat.StatusProbability,
		StatusResistance:  r.Data.AvatarCombat.StatusResistance,
	}
	t.RoleData.Light = light{
		ID:        r.Data.AvatarCombat.WeaponID,
		Name:      r.Data.AvatarWeapon.Name,
		Star:      r.Data.AvatarWeapon.Star,
		Level:     r.Data.AvatarCombat.WeaponLevel,
		Promotion: r.Data.AvatarCombat.WeaponPromotion,
		Rank:      r.Data.AvatarCombat.WeaponRank,
	}
	t.RoleData.Skill = skill{
		A: r.Data.AvatarCombat.SkillList[2].Level,
		E: r.Data.AvatarCombat.SkillList[1].Level,
		Q: r.Data.AvatarCombat.SkillList[0].Level,
		T: r.Data.AvatarCombat.SkillList[3].Level,
		F: r.Data.AvatarCombat.SkillList[4].Level,
	}
	for i := 0; i < len(r.Data.ItemRelic); i++ {
		switch r.Data.ItemRelic[i].Type {
		case "HEAD":
			t.RoleData.Relics.Head = relicsdata{
				RelicId: r.Data.ItemRelic[i].RelicID,
				Name:    r.Data.ItemRelic[i].Name,
				Type:    "HEAD",
			}
			for _, vv := range r.Data.AvatarRelics {
				if vv.Type == "HEAD" {
					t.RoleData.Relics.Head.MainV = vlist{
						Name:  typeToName(vv.MainAffixType),
						Value: Ftoone(vv.MainAffixValue * Sto100(typeToName(vv.MainAffixType))),
					}
					z := vlist{}
					for l := 0; l < 4; l++ {
						switch l {
						case 0:
							if vv.Sub1Type == 0 {
								continue
							}
							z.Name = typeToName(vv.Sub1Type)
							z.Value = Ftoone(vv.Sub1Value * Sto100(z.Name))
						case 1:
							if vv.Sub2Type == 0 {
								continue
							}
							z.Name = typeToName(vv.Sub2Type)
							z.Value = Ftoone(vv.Sub2Value * Sto100(z.Name))
						case 2:
							if vv.Sub3Type == 0 {
								continue
							}
							z.Name = typeToName(vv.Sub3Type)
							z.Value = Ftoone(vv.Sub3Value * Sto100(z.Name))
						case 3:
							if vv.Sub4Type == 0 {
								continue
							}
							z.Name = typeToName(vv.Sub4Type)
							z.Value = Ftoone(vv.Sub4Value * Sto100(z.Name))
						}
						t.RoleData.Relics.Head.Vlist = append(t.RoleData.Relics.Head.Vlist, z)
					}
					break
				}
			}
		case "HAND":
			t.RoleData.Relics.Hand = relicsdata{
				RelicId: r.Data.ItemRelic[i].RelicID,
				Name:    r.Data.ItemRelic[i].Name,
				Type:    "HAND",
			}
			for _, vv := range r.Data.AvatarRelics {
				if vv.Type == "HAND" {
					t.RoleData.Relics.Hand.MainV = vlist{
						Name:  typeToName(vv.MainAffixType),
						Value: Ftoone(vv.MainAffixValue * Sto100(typeToName(vv.MainAffixType))),
					}
					z := vlist{}
					for l := 0; l < 4; l++ {
						switch l {
						case 0:
							if vv.Sub1Type == 0 {
								continue
							}
							z.Name = typeToName(vv.Sub1Type)
							z.Value = Ftoone(vv.Sub1Value * Sto100(z.Name))
						case 1:
							if vv.Sub2Type == 0 {
								continue
							}
							z.Name = typeToName(vv.Sub2Type)
							z.Value = Ftoone(vv.Sub2Value * Sto100(z.Name))
						case 2:
							if vv.Sub3Type == 0 {
								continue
							}
							z.Name = typeToName(vv.Sub3Type)
							z.Value = Ftoone(vv.Sub3Value * Sto100(z.Name))
						case 3:
							if vv.Sub4Type == 0 {
								continue
							}
							z.Name = typeToName(vv.Sub4Type)
							z.Value = Ftoone(vv.Sub4Value * Sto100(z.Name))
						}
						t.RoleData.Relics.Hand.Vlist = append(t.RoleData.Relics.Hand.Vlist, z)
					}
					break
				}
			}
		case "BODY":
			t.RoleData.Relics.Body = relicsdata{
				RelicId: r.Data.ItemRelic[i].RelicID,
				Name:    r.Data.ItemRelic[i].Name,
				Type:    "BODY",
			}
			for _, vv := range r.Data.AvatarRelics {
				if vv.Type == "BODY" {
					t.RoleData.Relics.Body.MainV = vlist{
						Name:  typeToName(vv.MainAffixType),
						Value: Ftoone(vv.MainAffixValue * Sto100(typeToName(vv.MainAffixType))),
					}
					z := vlist{}
					for l := 0; l < 4; l++ {
						switch l {
						case 0:
							if vv.Sub1Type == 0 {
								continue
							}
							z.Name = typeToName(vv.Sub1Type)
							z.Value = Ftoone(vv.Sub1Value * Sto100(z.Name))
						case 1:
							if vv.Sub2Type == 0 {
								continue
							}
							z.Name = typeToName(vv.Sub2Type)
							z.Value = Ftoone(vv.Sub2Value * Sto100(z.Name))
						case 2:
							if vv.Sub3Type == 0 {
								continue
							}
							z.Name = typeToName(vv.Sub3Type)
							z.Value = Ftoone(vv.Sub3Value * Sto100(z.Name))
						case 3:
							if vv.Sub4Type == 0 {
								continue
							}
							z.Name = typeToName(vv.Sub4Type)
							z.Value = Ftoone(vv.Sub4Value * Sto100(z.Name))
						}
						t.RoleData.Relics.Body.Vlist = append(t.RoleData.Relics.Body.Vlist, z)
					}
					break
				}
			}
		case "FOOT":
			t.RoleData.Relics.Foot = relicsdata{
				RelicId: r.Data.ItemRelic[i].RelicID,
				Name:    r.Data.ItemRelic[i].Name,
				Type:    "FOOT",
			}
			for _, vv := range r.Data.AvatarRelics {
				if vv.Type == "FOOT" {
					t.RoleData.Relics.Foot.MainV = vlist{
						Name:  typeToName(vv.MainAffixType),
						Value: Ftoone(vv.MainAffixValue * Sto100(typeToName(vv.MainAffixType))),
					}
					z := vlist{}
					for l := 0; l < 4; l++ {
						switch l {
						case 0:
							if vv.Sub1Type == 0 {
								continue
							}
							z.Name = typeToName(vv.Sub1Type)
							z.Value = Ftoone(vv.Sub1Value * Sto100(z.Name))
						case 1:
							if vv.Sub2Type == 0 {
								continue
							}
							z.Name = typeToName(vv.Sub2Type)
							z.Value = Ftoone(vv.Sub2Value * Sto100(z.Name))
						case 2:
							if vv.Sub3Type == 0 {
								continue
							}
							z.Name = typeToName(vv.Sub3Type)
							z.Value = Ftoone(vv.Sub3Value * Sto100(z.Name))
						case 3:
							if vv.Sub4Type == 0 {
								continue
							}
							z.Name = typeToName(vv.Sub4Type)
							z.Value = Ftoone(vv.Sub4Value * Sto100(z.Name))
						}
						t.RoleData.Relics.Foot.Vlist = append(t.RoleData.Relics.Foot.Vlist, z)
					}
					break
				}
			}
		case "NECK":
			t.RoleData.Relics.Neck = relicsdata{
				RelicId: r.Data.ItemRelic[i].RelicID,
				Name:    r.Data.ItemRelic[i].Name,
				Type:    "NECK",
			}
			for _, vv := range r.Data.AvatarRelics {
				if vv.Type == "NECK" {
					t.RoleData.Relics.Neck.MainV = vlist{
						Name:  typeToName(vv.MainAffixType),
						Value: Ftoone(vv.MainAffixValue * Sto100(typeToName(vv.MainAffixType))),
					}
					z := vlist{}
					for l := 0; l < 4; l++ {
						switch l {
						case 0:
							if vv.Sub1Type == 0 {
								continue
							}
							z.Name = typeToName(vv.Sub1Type)
							z.Value = Ftoone(vv.Sub1Value * Sto100(z.Name))
						case 1:
							if vv.Sub2Type == 0 {
								continue
							}
							z.Name = typeToName(vv.Sub2Type)
							z.Value = Ftoone(vv.Sub2Value * Sto100(z.Name))
						case 2:
							if vv.Sub3Type == 0 {
								continue
							}
							z.Name = typeToName(vv.Sub3Type)
							z.Value = Ftoone(vv.Sub3Value * Sto100(z.Name))
						case 3:
							if vv.Sub4Type == 0 {
								continue
							}
							z.Name = typeToName(vv.Sub4Type)
							z.Value = Ftoone(vv.Sub4Value * Sto100(z.Name))
						}
						t.RoleData.Relics.Neck.Vlist = append(t.RoleData.Relics.Neck.Vlist, z)
					}
					break
				}
			}
		case "OBJECT":
			t.RoleData.Relics.Object = relicsdata{
				RelicId: r.Data.ItemRelic[i].RelicID,
				Name:    r.Data.ItemRelic[i].Name,
				Type:    "OBJECT",
			}
			for _, vv := range r.Data.AvatarRelics {
				if vv.Type == "OBJECT" {
					t.RoleData.Relics.Object.MainV = vlist{
						Name:  typeToName(vv.MainAffixType),
						Value: Ftoone(vv.MainAffixValue * Sto100(typeToName(vv.MainAffixType))),
					}
					z := vlist{}
					for l := 0; l < 4; l++ {
						switch l {
						case 0:
							if vv.Sub1Type == 0 {
								continue
							}
							z.Name = typeToName(vv.Sub1Type)
							z.Value = Ftoone(vv.Sub1Value * Sto100(z.Name))
						case 1:
							if vv.Sub2Type == 0 {
								continue
							}
							z.Name = typeToName(vv.Sub2Type)
							z.Value = Ftoone(vv.Sub2Value * Sto100(z.Name))
						case 2:
							if vv.Sub3Type == 0 {
								continue
							}
							z.Name = typeToName(vv.Sub3Type)
							z.Value = Ftoone(vv.Sub3Value * Sto100(z.Name))
						case 3:
							if vv.Sub4Type == 0 {
								continue
							}
							z.Name = typeToName(vv.Sub4Type)
							z.Value = Ftoone(vv.Sub4Value * Sto100(z.Name))
						}
						t.RoleData.Relics.Object.Vlist = append(t.RoleData.Relics.Object.Vlist, z)
					}
					break
				}
			}
		}
	}

	return *t, nil
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
	}
	return true
}
