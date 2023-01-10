package kokomi // 导入yuan-shen模块
import (
	"os"
	"strconv"
)

const (
// url = "https://enka.microgg.cn/u/%v/__data.json"
)

// 评分权重结构
type wifequan struct {
	Hp       int //生命
	Atk      int //攻击力
	Def      int //防御力
	Cpct     int //暴击率
	Cdmg     int //暴击伤害
	Mastery  int //元素精通
	Dmg      int //元素伤害
	Phy      int //物理伤害
	Recharge int //元素充能
	Heal     int //治疗加成
}

// Uidmap wifeid->wifename
var Uidmap = map[int64]string{ //
	10000036: "重云",
	10000050: "托马",
	10000051: "优菈",
	10000066: "神里绫人",
	10000067: "柯莱",
	10000016: "迪卢克",
	10000025: "行秋",
	10000030: "钟离",
	10000053: "早柚",
	10000071: "赛诺",
	10000002: "神里绫华",
	10000003: "琴",
	10000005: "空",
	10000068: "多莉",
	10000070: "妮露",
	10000072: "坎蒂丝",
	10000001: "凯特",
	10000055: "五郎",
	10000060: "夜兰",
	10000023: "香菱",
	10000042: "刻晴",
	10000039: "迪奥娜",
	10000057: "荒泷一斗",
	10000069: "提纳里",
	10000075: "流浪者",
	10000078: "艾尔海森",
	10000006: "丽莎",
	10000014: "芭芭拉",
	10000049: "宵宫",
	10000056: "九条裟罗",
	10000058: "八重神子",
	10000065: "久岐忍",
	10000044: "辛焱",
	10000047: "枫原万叶",
	10000046: "胡桃",
	10000048: "烟绯",
	10000063: "申鹤",
	10000076: "珐露珊",
	10000015: "凯亚",
	10000020: "雷泽",
	10000074: "莱依拉",
	10000022: "温迪",
	10000064: "云堇",
	10000034: "诺艾尔",
	10000077: "瑶瑶",
	10000021: "安柏",
	10000033: "达达利亚",
	10000043: "砂糖",
	10000029: "可莉",
	10000045: "罗莎莉亚",
	10000062: "埃洛伊",
	10000035: "七七",
	10000052: "雷电将军",
	10000031: "菲谢尔",
	10000037: "甘雨",
	10000038: "阿贝多",
	10000041: "莫娜",
	10000007: "荧",
	10000024: "北斗",
	10000054: "珊瑚宫心海",
	10000026: "魈",
	10000027: "凝光",
	10000073: "纳西妲",
	10000032: "班尼特",
	10000059: "鹿野院平藏",
}

// Data 从网站获取的数据
type Data struct {
	PlayerInfo struct {
		Nickname             string `json:"nickname"`
		Level                int    `json:"level"`
		Signature            string `json:"signature"`
		WorldLevel           int    `json:"worldLevel"`
		NameCardID           int    `json:"nameCardId"`
		FinishAchievementNum int    `json:"finishAchievementNum"`
		TowerFloorIndex      int    `json:"towerFloorIndex"`
		TowerLevelIndex      int    `json:"towerLevelIndex"`
		ShowAvatarInfoList   []struct {
			AvatarID  int `json:"avatarId"`
			Level     int `json:"level"`
			CostumeID int `json:"costumeId,omitempty"`
		} `json:"showAvatarInfoList"`
		ShowNameCardIDList []int `json:"showNameCardIdList"`
		ProfilePicture     struct {
			AvatarID int `json:"avatarId"`
		} `json:"profilePicture"`
	} `json:"playerInfo"`
	AvatarInfoList []struct {
		AvatarID int `json:"avatarId"`
		PropMap  struct {
			Num1001 struct {
				Type int    `json:"type"`
				Ival string `json:"ival"`
			} `json:"1001"`
			Num1002 struct {
				Type int    `json:"type"`
				Ival string `json:"ival"`
				Val  string `json:"val"`
			} `json:"1002"`
			Num1003 struct {
				Type int    `json:"type"`
				Ival string `json:"ival"`
			} `json:"1003"`
			Num1004 struct {
				Type int    `json:"type"`
				Ival string `json:"ival"`
			} `json:"1004"`
			Num4001 struct {
				Type int    `json:"type"`
				Ival string `json:"ival"`
				Val  string `json:"val"`
			} `json:"4001"`
			Num10010 struct {
				Type int    `json:"type"`
				Ival string `json:"ival"`
				Val  string `json:"val"`
			} `json:"10010"`
		} `json:"propMap"`
		FightPropMap struct {
			Num1    float64 `json:"1"`
			Num2    float64 `json:"2"`
			Num3    float64 `json:"3"`
			Num4    float64 `json:"4"`
			Num5    float64 `json:"5"`
			Num6    float64 `json:"6"`
			Num7    float64 `json:"7"`
			Num8    float64 `json:"8"`
			Num20   float64 `json:"20"`
			Num21   float64 `json:"21"`
			Num22   float64 `json:"22"`
			Num23   float64 `json:"23"`
			Num26   float64 `json:"26"`
			Num27   float64 `json:"27"`
			Num28   float64 `json:"28"`
			Num29   float64 `json:"29"`
			Num30   float64 `json:"30"`
			Num40   float64 `json:"40"`
			Num41   float64 `json:"41"`
			Num42   float64 `json:"42"`
			Num43   float64 `json:"43"`
			Num44   float64 `json:"44"`
			Num45   float64 `json:"45"`
			Num46   float64 `json:"46"`
			Num50   float64 `json:"50"`
			Num51   float64 `json:"51"`
			Num52   float64 `json:"52"`
			Num53   float64 `json:"53"`
			Num54   float64 `json:"54"`
			Num55   float64 `json:"55"`
			Num56   float64 `json:"56"`
			Num70   float64 `json:"70"`
			Num80   float64 `json:"80"`
			Num1000 float64 `json:"1000"`
			Num1010 float64 `json:"1010"`
			Num2000 float64 `json:"2000"`
			Num2001 float64 `json:"2001"`
			Num2002 float64 `json:"2002"`
			Num2003 float64 `json:"2003"`
			Num3007 float64 `json:"3007"`
			Num3008 float64 `json:"3008"`
			Num3015 float64 `json:"3015"`
			Num3016 float64 `json:"3016"`
			Num3017 float64 `json:"3017"`
			Num3018 float64 `json:"3018"`
			Num3019 float64 `json:"3019"`
			Num3020 float64 `json:"3020"`
			Num3021 float64 `json:"3021"`
			Num3022 float64 `json:"3022"`
			Num3045 float64 `json:"3045"`
			Num3046 float64 `json:"3046"`
		} `json:"fightPropMap"`
		SkillDepotID           int         `json:"skillDepotId"`
		InherentProudSkillList []int       `json:"inherentProudSkillList"`
		SkillLevelMap          map[int]int `json:"skillLevelMap"`
		EquipList              []struct {
			ItemID    int `json:"itemId"`
			Reliquary struct {
				Level            int   `json:"level"`
				MainPropID       int   `json:"mainPropId"`
				AppendPropIDList []int `json:"appendPropIdList"`
			} `json:"reliquary,omitempty"`
			Flat   Flat `json:"flat"` //标记
			Weapon struct {
				Level        int         `json:"level"`
				PromoteLevel int         `json:"promoteLevel"`
				AffixMap     map[int]int `json:"affixMap"`
			} `json:"weapon,omitempty"`
		} `json:"equipList"`
		FetterInfo struct {
			ExpLevel int `json:"expLevel"`
		} `json:"fetterInfo"`
		TalentIDList            []int `json:"talentIdList,omitempty"`
		ProudSkillExtraLevelMap struct {
			Num4239 int `json:"4239"`
		} `json:"proudSkillExtraLevelMap,omitempty"`
		CostumeID int `json:"costumeId,omitempty"`
	} `json:"avatarInfoList"`
	TTL int    `json:"ttl"`
	UID string `json:"uid"`
}

// Flat ... 详细数据
type Flat struct {
	// l10n
	NameTextHash    string `json:"nameTextMapHash"`
	SetNameTextHash string `json:"setNameTextMapHash,omitempty"`

	// artifact
	ReliquaryMainStat Stat   `json:"reliquaryMainstat,omitempty"`
	ReliquarySubStats []Stat `json:"reliquarySubstats,omitempty"`
	EquipType         string `json:"equipType,omitempty"`

	// weapon
	WeaponStat []Stat `json:"weaponStats,omitempty"`

	RankLevel uint8  `json:"rankLevel"` // 3, 4 or 5
	ItemType  string `json:"itemType"`  // ITEM_WEAPON or ITEM_RELIQUARY
	Icon      string `json:"icon"`      // You can get the icon from https://enka.network/ui/{Icon}.png
}

// Stat ...  属性对
type Stat struct {
	MainPropID string  `json:"mainPropId,omitempty"`
	SubPropID  string  `json:"appendPropId,omitempty"`
	Value      float64 `json:"statValue"`
}

// Getuid qquid->uid
func Getuid(qquid int64) (uid int) { // 获取对应游戏uid
	sqquid := strconv.Itoa(int(qquid))
	// 获取本地缓存数据
	txt, err := os.ReadFile("plugin/kokomi/data/uid/" + sqquid + ".kokomi")
	if err != nil {
		return 0
	}
	sss, _ := strconv.Atoi(string(txt))
	return sss
}

// StoS 圣遗物词条简单描述
func StoS(val string) string {
	var vv string
	switch val {
	case "FIGHT_PROP_HP":
		vv = "小生命"
	case "FIGHT_PROP_HP_PERCENT":
		vv = "大生命"
	case "FIGHT_PROP_ATTACK":
		vv = "小攻击"
	case "FIGHT_PROP_ATTACK_PERCENT":
		vv = "大攻击"
	case "FIGHT_PROP_DEFENSE":
		vv = "小防御"
	case "FIGHT_PROP_DEFENSE_PERCENT":
		vv = "大防御"
	case "FIGHT_PROP_CRITICAL":
		vv = "暴击率"
	case "FIGHT_PROP_CRITICAL_HURT":
		vv = "暴击伤害"
	case "FIGHT_PROP_CHARGE_EFFICIENCY":
		vv = "元素充能"
	case "FIGHT_PROP_HEAL_ADD":
		vv = "治疗加成"
	case "FIGHT_PROP_ELEMENT_MASTERY":
		vv = "元素精通"
	case "FIGHT_PROP_PHYSICAL_ADD_HURT":
		vv = "物理加伤"
	case "FIGHT_PROP_FIRE_ADD_HURT":
		vv = "火元素加伤"
	case "FIGHT_PROP_ELEC_ADD_HURT":
		vv = "雷元素加伤"
	case "FIGHT_PROP_WATER_ADD_HURT":
		vv = "水元素加伤"
	case "FIGHT_PROP_GRASS_ADD_HURT":
		vv = "草元素加伤"
	case "FIGHT_PROP_WIND_ADD_HURT":
		vv = "风元素加伤"
	case "FIGHT_PROP_ROCK_ADD_HURT":
		vv = "岩元素加伤"
	case "FIGHT_PROP_ICE_ADD_HURT":
		vv = "冰元素加伤"
	}
	return vv
}

// Stofen 判断词条分号
func Stofen(val string) string {
	var vv = "%"
	switch val {
	case "FIGHT_PROP_HP":
		vv = ""
	case "FIGHT_PROP_HP_PERCENT":
	case "FIGHT_PROP_ATTACK":
		vv = ""
	case "FIGHT_PROP_ATTACK_PERCENT":
	case "FIGHT_PROP_DEFENSE":
		vv = ""
	case "FIGHT_PROP_DEFENSE_PERCENT":
	case "FIGHT_PROP_CRITICAL":
	case "FIGHT_PROP_CRITICAL_HURT":
	case "FIGHT_PROP_CHARGE_EFFICIENCY":
	case "FIGHT_PROP_HEAL_ADD":
	case "FIGHT_PROP_ELEMENT_MASTERY":
		vv = ""
	case "FIGHT_PROP_PHYSICAL_ADD_HURT":
	case "FIGHT_PROP_FIRE_ADD_HURT":
	case "FIGHT_PROP_ELEC_ADD_HURT":
	case "FIGHT_PROP_WATER_ADD_HURT":
	case "FIGHT_PROP_GRASS_ADD_HURT":
	case "FIGHT_PROP_WIND_ADD_HURT":
	case "FIGHT_PROP_ROCK_ADD_HURT":
	case "FIGHT_PROP_ICE_ADD_HURT":
	}
	return vv
}

// Tianfujiuzhen 修复部分贴图大小错误
func Tianfujiuzhen(val string) int {
	var bb = 257 //280
	switch val {
	case "芭芭拉", "北斗", "多莉", "甘雨", "胡桃", "科莱", "雷电将军", "罗莎莉亚", "凝光", "赛诺", "魈", "行秋", "烟绯", "夜兰", "早柚":
		bb = 280
	}
	return bb
}

// Countcitiao 计算圣遗物单词条分
func Countcitiao(wifename string, funame string, figure float64) float64 {
	var grade float64
	var ti wifequan
	ti = Wifequanmap[wifename]
	switch funame {
	case "大生命":
		grade = figure * 1.33 * float64(ti.Hp) / 100
	case "大攻击":
		grade = figure * 1.33 * float64(ti.Atk) / 100
	case "大防御":
		grade = figure * 1.33 * float64(ti.Def) / 100
	case "暴击率":
		grade = figure * 2.0 * float64(ti.Cpct) / 100
	case "暴击伤害":
		grade = figure * 1.0 * float64(ti.Cdmg) / 100
	case "元素精通":
		grade = figure * 0.33 * float64(ti.Mastery) / 100
	case "雷元素加伤", "水元素加伤", "火元素加伤", "风元素加伤", "草元素加伤", "岩元素加伤", "冰元素加伤":
		grade = figure * 1.33 * float64(ti.Dmg) / 100
	case "物理加伤":
		grade = figure * 1.33 * float64(ti.Phy) / 100
	case "元素充能":
		grade = figure * 1.2 * float64(ti.Recharge) / 100
	case "治疗加成":
		grade = figure * 1.73 * float64(ti.Heal) / 100
	default:
		grade = 0
	}
	return grade
}

// Pingji 词条评级
func Pingji(val float64) string {
	var fff string
	switch {
	case val < 18:
		fff = "C"
	case val < 24:
		fff = "B"
	case val < 29.7:
		fff = "A"
	case val < 36.3:
		fff = "S"
	case val < 42.9:
		fff = "SS"
	case val < 49.5:
		fff = "SSS"
	case val < 56.1:
		fff = "ACE"
	case val >= 56.1:
		fff = "ACE^2"
	}
	return fff
}

// Ftoone 保留一位小数并转化string
func Ftoone(f float64) string {
	return strconv.FormatFloat((float64(int(f*10)))/10.0, 'f', 1, 64)
}
