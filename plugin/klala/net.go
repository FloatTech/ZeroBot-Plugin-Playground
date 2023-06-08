package klala

const (
	// nets = "https://mhy.fuckmys.tk/"
	nets = "http://api.mihomo.me/"
	path = "sr_info/"
)

type info struct {
	PlayerDetailInfo struct {
		AssistAvatar        role   `json:"assistAvatarDetail"`
		IsDisplayAvatarList bool   `json:"isDisplayAvatar"`
		DisplayAvatarList   []role `json:"avatarDetailList"`
		UID                 int    `json:"uid"`
		CurFriendCount      int    `json:"friendCount"`
		WorldLevel          int    `json:"worldLevel"`
		NickName            string `json:"nickname"`
		Level               int    `json:"level"`
		PlayerSpaceInfo     struct {
			ChallengeData struct {
				PreMazeGroupIndex int `json:"scheduleMaxLevel"`
			} `json:"challengeInfo"`
			PassAreaProgress int `json:"maxRogueChallengeScore"`
			LightConeCount   int `json:"equipmentCount"`
			AvatarCount      int `json:"avatarCount"`
			AchievementCount int `json:"achievementCount"`
		} `json:"recordInfo"`
		HeadIconID int `json:"headIcon"`
	} `json:"detailInfo"`
}

/*
	type role struct {
		Promotion    int `json:"promotion"` //角色晋阶
		AvatarId     int `json:"avatarId"`
		BehaviorList []struct {
			BehaviorID int `json:"avatarId"`
			Level      int `json:"promotion"`
		} `json:"BehaviorList"`
		AvatarID    int `json:"AvatarID"`
		Level       int `json:"Level"`
		EquipmentID struct {
			Level     int `json:"Level"`
			ID        int `json:"ID"`
			Promotion int `json:"Promotion"`
			Rank      int `json:"Rank"`
		} `json:"EquipmentID"`
		RelicList []struct {
			RelicSubAffix []struct {
				SubAffixID int `json:"SubAffixID"`
				Cnt        int `json:"Cnt"`
				Step       int `json:"Step"`
			} `json:"RelicSubAffix"`
			ID          int     `json:"ID"`
			MainAffixID int     `json:"MainAffixID"`
			Level       float64 `json:"Level,omitempty"`
			Type        int     `json:"Type"`
			EXP         int     `json:"EXP,omitempty"`
		} `json:"RelicList"`
	}
*/
type role struct {
	AvatarID      int `json:"avatarId"`
	Promotion     int `json:"promotion"`
	Exp           int `json:"exp,omitempt"`
	Level         int `json:"level,omitempty"`
	Rank          int `json:"rank,omitempt"`
	SkillTreeList []struct {
		PointID int `json:"pointId"`
		Level   int `json:"level"`
	} `json:"skillTreeList"`
	RelicList []struct {
		Exp          int `json:"exp,omitempty"`
		Type         int `json:"type"`
		Tid          int `json:"tid"`
		SubAffixList []struct {
			AffixID int `json:"affixId"`
			Cnt     int `json:"cnt"`
			Step    int `json:"step,omitempty"`
		} `json:"subAffixList"`
		Level       float64 `json:"level,omitempty"`
		MainAffixID int     `json:"mainAffixId"`
	} `json:"relicList"`
	Equipment struct {
		Level     int `json:"level"`
		Tid       int `json:"tid"`
		Promotion int `json:"promotion"`
		Rank      int `json:"rank"`
	} `json:"equipment"`
}

// 本地数据
type thisdata struct {
	UID      string `json:"uid"`
	Nickname string `json:"nickname"`
	Level    int    `json:"level"`
	RoleData []ro   `json:"data"`
}
type ro struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Star    string  `json:"star"`
	Rank    int     `json:"rank"`
	Path    string  `json:"path"`    //命途
	Element string  `json:"element"` //元素属性
	List    combat  `json:"combat"`  //属性列表
	Light   light   `json:"light"`
	Skill   skill   `json:"skill"`  //技能
	Relics  relics  `json:"rolics"` //遗物
	Scores  float64 `json:"scores"`
}
type combat struct {
	AvatarID          int     `json:"avatarId"` //角色序号
	Level             int     `json:"level"`
	Promotion         int     `json:"promotion"`         //突破次数
	HpBase            float64 `json:"hpBase"`            //hp白值
	HpFinal           float64 `json:"hpFinal"`           //总和
	AttackBase        float64 `json:"attackBase"`        //akt白值
	AttackFinal       float64 `json:"attackFinal"`       //akt总和
	DefenseBase       float64 `json:"defenseBase"`       //防御白值
	DefenseFinal      float64 `json:"defenseFinal"`      //防御
	SpeedBase         int     `json:"speedBase"`         //速度白值
	SpeedFinal        float64 `json:"speedFinal"`        //速度
	CriticalChance    float64 `json:"criticalChance"`    //暴击率
	CriticalDamage    float64 `json:"criticalDamage"`    //暴击伤害
	HealRatio         int     `json:"healRatio"`         //治疗
	StatusProbability float64 `json:"statusProbability"` //效果命中
	StatusResistance  float64 `json:"statusResistance"`  //效果抵抗
}

type light struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Star      int    `json:"star"`
	Level     int    `json:"level"`
	Promotion int    `json:"promotion"`
	Rank      int    `json:"rank"`
	Vice      string `json:"vice"`
}

type skill struct {
	A int `json:"a"`
	E int `json:"e"`
	Q int `json:"q"`
	T int `json:"t"`
	F int `json:"f"`
}
type relics struct {
	Head   relicsdata `json:"head"`
	Hand   relicsdata `json:"hand"`
	Body   relicsdata `json:"body"`
	Foot   relicsdata `json:"foot"`
	Neck   relicsdata `json:"neck"`
	Object relicsdata `json:"object"`
}
type relicsdata struct {
	SetID int     `json:"relicId"` //套装id
	Name  string  `json:"name"`
	Star  int     `json:"star"`
	Type  int     `json:"type"`
	Level float64 `json:"level"`
	MainV vlist   `json:"main"`
	Vlist []vlist `json:"vlist"`
}
type vlist struct {
	Name  string  `json:"valname"`
	Value string  `json:"value"`
	Score float64 `json:"score"`
	Adds  int     `json:"adds"`
}
type lightmap map[string]struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Rarity        int      `json:"rarity"`
	Path          string   `json:"path"`
	Desc          string   `json:"desc"`
	Icon          string   `json:"icon"`
	Preview       string   `json:"preview"`
	Portrait      string   `json:"portrait"`
	GuideOverview []string `json:"guide_overview"`
}

// FindMap 各种简称map查询
type FindMap struct {
	Characters map[string][]string `json:"characters"`
	LightCones map[string][]string `json:"light_cones"`
}

// 词条信息

type affixStarMap map[string]affixTypeMap

type affixTypeMap map[string]struct {
	GroupID   int    `json:"GroupID"`
	AffixID   int    `json:"AffixID"`
	Property  string `json:"Property"`
	BaseValue struct {
		Value float64 `json:"Value"`
	} `json:"BaseValue"`
	StepValue struct {
		Value float64 `json:"Value"`
	} `json:"StepValue"`
	StepNum int `json:"StepNum"`
}
type affixStarMainMap map[string]affixMainMap
type affixMainMap map[string]struct {
	GroupID   int    `json:"GroupID"`
	AffixID   int    `json:"AffixID"`
	Property  string `json:"Property"`
	BaseValue struct {
		Value float64 `json:"Value"`
	} `json:"BaseValue"`
	LevelAdd struct {
		Value float64 `json:"Value"`
	} `json:"LevelAdd"`
	IsAvailable bool `json:"IsAvailable"`
}
type relicConfigMap map[string]struct {
	ID             int    `json:"ID"`
	SetID          int    `json:"SetID"`
	Type           string `json:"Type"`
	Rarity         string `json:"Rarity"`
	MainAffixGroup int    `json:"MainAffixGroup"`
	SubAffixGroup  int    `json:"SubAffixGroup"`
	MaxLevel       int    `json:"MaxLevel"`
	ExpType        int    `json:"ExpType"`
	ExpProvide     int    `json:"ExpProvide"`
	CoinCost       int    `json:"CoinCost"`
}

type yiwumap map[string]struct {
	ID          string `json:"id"`
	SetID       string `json:"set_id"`
	Name        string `json:"name"`
	Rarity      int    `json:"rarity"`
	Type        string `json:"type"`
	MaxLevel    int    `json:"max_level"`
	MainAffixID string `json:"main_affix_id"`
	SubAffixID  string `json:"sub_affix_id"`
	Icon        string `json:"icon"`
}

type ropeData struct {
	RelicType     string `json:"relic_type"`
	RelicTypeText string `json:"relic_type_text"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Backstory     string `json:"backstory"`
	Icon          string `json:"icon"`
}
type wifeData map[string]struct {
	ID     string `json:"id"`
	Values []struct {
		Hp struct {
			Base float64 `json:"base"`
			Step float64 `json:"step"`
		} `json:"hp"`
		Atk struct {
			Base float64 `json:"base"`
			Step float64 `json:"step"`
		} `json:"atk"`
		Def struct {
			Base float64 `json:"base"`
			Step float64 `json:"step"`
		} `json:"def"`
		Spd struct {
			Base float64 `json:"base"`
			Step float64 `json:"step"`
		} `json:"spd"`
		Taunt struct {
			Base float64 `json:"base"`
			Step float64 `json:"step"`
		} `json:"taunt"`
		CritRate struct {
			Base float64 `json:"base"`
			Step float64 `json:"step"`
		} `json:"crit_rate"`
		CritDmg struct {
			Base float64 `json:"base"`
			Step float64 `json:"step"`
		} `json:"crit_dmg"`
	} `json:"values"`
	Materials [][]struct {
		ID  string `json:"id"`
		Num int    `json:"num"`
	} `json:"materials"`
}

type lightData map[string]struct {
	ID     string `json:"id"`
	Values []struct {
		Hp struct {
			Base float64 `json:"base"`
			Step float64 `json:"step"`
		} `json:"hp"`
		Atk struct {
			Base float64 `json:"base"`
			Step float64 `json:"step"`
		} `json:"atk"`
		Def struct {
			Base float64 `json:"base"`
			Step float64 `json:"step"`
		} `json:"def"`
	} `json:"values"`
	Materials [][]struct {
		ID  string `json:"id"`
		Num int    `json:"num"`
	} `json:"materials"`
}

type lightAffix map[string]struct {
	ID         string      `json:"id"`
	Skill      string      `json:"skill"`
	Desc       string      `json:"desc"`
	Params     [][]float64 `json:"params"`
	Properties [][]struct {
		Type  string  `json:"type"`
		Value float64 `json:"value"`
	} `json:"properties"`
}

type ywSetData map[string]struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Properties [][]struct {
		Type  string  `json:"type"`
		Value float64 `json:"value"`
	} `json:"properties"`
	Icon string `json:"icon"`
}

type wifeTrees map[string]struct {
	ID       string `json:"id"`
	MaxLevel int    `json:"max_level"`
	Anchor   string `json:"anchor"`
	//	PrePoints     []any  `json:"pre_points"`
	//	LevelUpSkills []any  `json:"level_up_skills"`
	Levels []struct {
		Promotion  int `json:"promotion"`
		Properties []struct {
			Type  string  `json:"type"`
			Value float64 `json:"value"`
		} `json:"properties"`
		Materials []struct {
			ID  string `json:"id"`
			Num int    `json:"num"`
		} `json:"materials"`
	} `json:"levels"`
}

type wifeIntrod map[string]struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Rarity   string `json:"rarity"`
	Element  string `json:"element"`
	Path     string `json:"path"`
	Eidolons []struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Effect string `json:"effect"`
		Icon   string `json:"icon"`
	} `json:"eidolons"`
	Skills struct {
		BasicAtk  introdData `json:"basic_atk"`
		Skill     introdData `json:"skill"`
		Ultimate  introdData `json:"ultimate"`
		Talent    introdData `json:"talent"`
		Technique introdData `json:"technique"`
	} `json:"skills"`
	VersionAdded      string   `json:"version_added"`
	Icon              string   `json:"icon"`
	ElementIcon       string   `json:"element_icon"`
	PathIcon          string   `json:"path_icon"`
	Preview           string   `json:"preview"`
	Portrait          string   `json:"portrait"`
	CharacterOverview []string `json:"character_overview"`
	CharacterMaterial string   `json:"character_material"`
}
type introdData struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Tag         string `json:"tag"`
	MaxLevel    int    `json:"max_level"`
	Effect      string `json:"effect"`
	ElementType string `json:"element_type"`
	Icon        string `json:"icon"`
}

type weightData map[string]map[string]float64

// 词条英文对应中文
var typeMap = map[string]string{
	"MaxHP":                     "生命值",
	"Attack":                    "攻击力",
	"Defence":                   "防御力",
	"Speed":                     "速度",
	"CriticalChance":            "暴击率",
	"CriticalDamage":            "暴击伤害",
	"BreakDamageAddedRatio":     "击破特攻",
	"BreakDamageAddedRatioBase": "击破特攻",
	"HealRatio":                 "治疗量加成",
	"MaxSP":                     "能量上限",
	"SPRatio":                   "能量恢复效率",
	"StatusProbability":         "效果命中",
	"StatusResistance":          "效果抵抗",
	"CriticalChanceBase":        "暴击率",
	"CriticalDamageBase":        "暴击伤害",
	"HealRatioBase":             "治疗量加成",
	"StanceBreakAddedRatio":     "dev_失效字段",
	"SPRatioBase":               "能量恢复效率",
	"StatusProbabilityBase":     "效果命中",
	"StatusResistanceBase":      "效果抵抗",
	"PhysicalAddedRatio":        "物理属性伤害",
	"PhysicalResistance":        "物理属性抗性",
	"FireAddedRatio":            "火属性伤害",
	"FireResistance":            "火属性抗性",
	"IceAddedRatio":             "冰属性伤害",
	"IceResistance":             "冰属性抗性",
	"ThunderAddedRatio":         "雷属性伤害",
	"ThunderResistance":         "雷属性抗性",
	"WindAddedRatio":            "风属性伤害",
	"WindResistance":            "风属性抗性",
	"QuantumAddedRatio":         "量子属性伤害",
	"QuantumResistance":         "量子属性抗性",
	"ImaginaryAddedRatio":       "虚数属性伤害",
	"ImaginaryResistance":       "虚数属性抗性",
	"BaseHP":                    "基础生命值提高<unbreak>#1[i]</unbreak>",
	"HPDelta":                   "生命值",
	"HPAddedRatio":              "大生命",
	"BaseAttack":                "基础攻击力提高<unbreak>#1[i]</unbreak>",
	"AttackDelta":               "攻击力",
	"AttackAddedRatio":          "大攻击",
	"BaseDefence":               "基础防御力提高<unbreak>#1[i]</unbreak>",
	"DefenceDelta":              "防御力",
	"DefenceAddedRatio":         "大防御",
	"BaseSpeed":                 "速度",
	"HealTakenRatio":            "治疗量加成",
	"PhysicalResistanceDelta":   "物理属性抗性",
	"FireResistanceDelta":       "火属性抗性",
	"IceResistanceDelta":        "冰属性抗性",
	"ThunderResistanceDelta":    "雷属性抗性",
	"WindResistanceDelta":       "风属性抗性",
	"QuantumResistanceDelta":    "量子属性抗性",
	"ImaginaryResistanceDelta":  "虚数属性抗性",
	"SpeedDelta":                "速度",
}

// 词条权重
var weiMap = map[string]float64{
	"大攻击":    1.35,
	"大生命":    1.80,
	"大防御":    0.72,
	"暴击率":    2.4,
	"暴击伤害":   1.2,
	"能量恢复效率": 1.25,
	"速度":     3,
	"效果抵抗":   1.35,
	"效果命中":   1.35,
	"击破特攻":   1.35,
	"量子属性伤害": 1.25,
	"虚数属性伤害": 1.25,
	"物理属性伤害": 1.25,
	"火属性伤害":  1.25,
	"冰属性伤害":  1.25,
	"雷属性伤害":  1.25,
	"风属性伤害":  1.25,
	"治疗量加成":  1.25,
}
