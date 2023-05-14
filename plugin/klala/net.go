package klala

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
)

const (
	// nets = "https://mhy.fuckmys.tk/"
	nets = "http://api.mihomo.me/"
	path = "sr_info/"
)

func getRole(uid string) (body []byte, err error) {
	var client = &http.Client{}
	req, err := http.NewRequest(http.MethodGet, nets+path+uid, nil)

	if err != nil {
		return nil, err
	}
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.102 Safari/537.36 MicroMessenger/7.0.20.1781(0x6700143B) NetType/WIFI MiniProgramEnv/Windows WindowsWechat/WMPF XWEB/6763")
	req.Header.Add("Accept", "*/*")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		return nil, errors.New("获取数据失败, Code: " + strconv.Itoa(res.StatusCode))
	}
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return
}

type info struct {
	PlayerDetailInfo struct {
		AssistAvatar        role   `json:"AssistAvatar"`
		IsDisplayAvatarList bool   `json:"IsDisplayAvatarList"`
		DisplayAvatarList   []role `json:"DisplayAvatarList"`
		UID                 int    `json:"UID"`
		CurFriendCount      int    `json:"CurFriendCount"`
		WorldLevel          int    `json:"WorldLevel"`
		NickName            string `json:"NickName"`
		Birthday            int    `json:"Birthday"`
		Level               int    `json:"Level"`
		PlayerSpaceInfo     struct {
			ChallengeData struct {
				PreMazeGroupIndex int `json:"PreMazeGroupIndex"`
			} `json:"ChallengeData"`
			PassAreaProgress int `json:"PassAreaProgress"`
			LightConeCount   int `json:"LightConeCount"`
			AvatarCount      int `json:"AvatarCount"`
			AchievementCount int `json:"AchievementCount"`
		} `json:"PlayerSpaceInfo"`
		HeadIconID int `json:"HeadIconID"`
	} `json:"PlayerDetailInfo"`
}

type role struct {
	BehaviorList []struct {
		BehaviorID int `json:"BehaviorID"`
		Level      int `json:"Level"`
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
	Promotion int `json:"Promotion"` //角色晋阶
	Rank      int `json:"Rank"`      //星魂
}

// 本地数据
type thisdata struct {
	UID      string `json:"uid"`
	Nickname string `json:"nickname"`
	Level    int    `json:"level"`
	RoleData []ro   `json:"data"`
}
type ro struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Star    int    `json:"star"`
	Rank    int    `json:"rank"`
	Type    string `json:"type"`    //命途
	Element string `json:"element"` //元素属性
	List    combat `json:"combat"`  //属性列表
	Light   light  `json:"light"`
	Skill   skill  `json:"skill"`  //技能
	Relics  relics `json:"rolics"` //遗物
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
	Name  string `json:"valname"`
	Value string `json:"value"`
	Adds  int    `json:"adds"`
}
type lightmap map[string]struct {
	ID                string   `json:"id"`
	Name              string   `json:"name"`
	Rarity            int      `json:"rarity"`
	Path              string   `json:"path"`
	EffectName        string   `json:"effect_name"`
	Effects           []string `json:"effects"`
	VersionAdded      string   `json:"version_added"`
	Icon              string   `json:"icon"`
	LightConeOverview []string `json:"light_cone_overview"`
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
	ID      string `json:"id"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Effects struct {
		Pieces2 string `json:"pieces2"`
		Pieces4 string `json:"pieces4"`
	} `json:"effects"`
	Pieces struct {
		Head struct {
			RelicType     string `json:"relic_type"`
			RelicTypeText string `json:"relic_type_text"`
			Name          string `json:"name"`
			Description   string `json:"description"`
			Backstory     string `json:"backstory"`
			Icon          string `json:"icon"`
		} `json:"head"`
		Hands struct {
			RelicType     string `json:"relic_type"`
			RelicTypeText string `json:"relic_type_text"`
			Name          string `json:"name"`
			Description   string `json:"description"`
			Backstory     string `json:"backstory"`
			Icon          string `json:"icon"`
		} `json:"hands"`
		Body struct {
			RelicType     string `json:"relic_type"`
			RelicTypeText string `json:"relic_type_text"`
			Name          string `json:"name"`
			Description   string `json:"description"`
			Backstory     string `json:"backstory"`
			Icon          string `json:"icon"`
		} `json:"body"`
		Feet struct {
			RelicType     string `json:"relic_type"`
			RelicTypeText string `json:"relic_type_text"`
			Name          string `json:"name"`
			Description   string `json:"description"`
			Backstory     string `json:"backstory"`
			Icon          string `json:"icon"`
		} `json:"feet"`
		PlanarSphere struct {
			RelicType     string `json:"relic_type"`
			RelicTypeText string `json:"relic_type_text"`
			Name          string `json:"name"`
			Description   string `json:"description"`
			Backstory     string `json:"backstory"`
			Icon          string `json:"icon"`
		} `json:"planar_sphere"`
		LinkRope struct {
			RelicType     string `json:"relic_type"`
			RelicTypeText string `json:"relic_type_text"`
			Name          string `json:"name"`
			Description   string `json:"description"`
			Backstory     string `json:"backstory"`
			Icon          string `json:"icon"`
		} `json:"link_rope"`
	} `json:"pieces"`
	VersionAdded string `json:"version_added"`
	Icon         string `json:"icon"`
}

type wifeData map[string]wifeDataMap
type wifeDataMap map[string]struct {
	AvatarID          int `json:"AvatarID"`
	Promotion         int `json:"Promotion"`
	PromotionCostList []struct {
		ItemID  int `json:"ItemID"`
		ItemNum int `json:"ItemNum"`
	} `json:"PromotionCostList"`
	MaxLevel           int `json:"MaxLevel"`
	PlayerLevelRequire int `json:"PlayerLevelRequire"`
	AttackBase         struct {
		Value float64 `json:"Value"`
	} `json:"AttackBase"`
	AttackAdd struct {
		Value float64 `json:"Value"`
	} `json:"AttackAdd"`
	DefenceBase struct {
		Value float64 `json:"Value"`
	} `json:"DefenceBase"`
	DefenceAdd struct {
		Value float64 `json:"Value"`
	} `json:"DefenceAdd"`
	HPBase struct {
		Value float64 `json:"Value"`
	} `json:"HPBase"`
	HPAdd struct {
		Value float64 `json:"Value"`
	} `json:"HPAdd"`
	SpeedBase struct {
		Value int `json:"Value"`
	} `json:"SpeedBase"`
	CriticalChance struct {
		Value float64 `json:"Value"`
	} `json:"CriticalChance"`
	CriticalDamage struct {
		Value float64 `json:"Value"`
	} `json:"CriticalDamage"`
	BaseAggro struct {
		Value float64 `json:"Value"`
	} `json:"BaseAggro"`
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
