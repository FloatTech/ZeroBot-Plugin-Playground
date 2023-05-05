package klala

import (
	"io/ioutil"
	"net/http"
)

const (
	upURL   = "https://xt.qyinter.com/update/"
	caseURL = "https://xt.qyinter.com/case/"
	roleURL = "https://xt.qyinter.com/overview/"
)

var (
	client = &http.Client{}
)

func Updata(uid string) (body []byte, err error) {
	req, err := http.NewRequest(http.MethodGet, upURL+uid, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Host", "xt.qyinter.com")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("referer", "https://servicewechat.com/wx06d9c99ccf3356e5/5/page-frame.html")
	req.Header.Add("xweb_xhr", "1")
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.102 Safari/537.36 MicroMessenger/7.0.20.1781(0x6700143B) NetType/WIFI MiniProgramEnv/Windows WindowsWechat/WMPF XWEB/6763")
	req.Header.Add("Accept", "*/*")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return
}

func CaseList(uid string) (body []byte, err error) {
	req, err := http.NewRequest(http.MethodGet, caseURL+uid, nil)

	if err != nil {
		return nil, err
	}
	req.Header.Add("Host", "xt.qyinter.com")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("referer", "https://servicewechat.com/wx06d9c99ccf3356e5/5/page-frame.html")
	req.Header.Add("xweb_xhr", "1")
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.102 Safari/537.36 MicroMessenger/7.0.20.1781(0x6700143B) NetType/WIFI MiniProgramEnv/Windows WindowsWechat/WMPF XWEB/6763")
	req.Header.Add("Accept", "*/*")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return
}

func GetRole(uid, n string) (body []byte, err error) {
	req, err := http.NewRequest(http.MethodGet, roleURL+uid+"/"+n, nil)

	if err != nil {
		return nil, err
	}
	req.Header.Add("Host", "xt.qyinter.com")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("referer", "https://servicewechat.com/wx06d9c99ccf3356e5/5/page-frame.html")
	req.Header.Add("xweb_xhr", "1")
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.102 Safari/537.36 MicroMessenger/7.0.20.1781(0x6700143B) NetType/WIFI MiniProgramEnv/Windows WindowsWechat/WMPF XWEB/6763")
	req.Header.Add("Accept", "*/*")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return
}

type roleList struct {
	Code int `json:"code"`
	Data struct {
		Characters []characters `json:"characters"`
		NickName   string       `json:"nickName"`
	} `json:"data"`
	Msg string `json:"msg"`
}

type characters struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Star         int    `json:"star"`
	Type         string `json:"type"`
	Element      string `json:"element"`
	Desc         string `json:"desc"`
	Icon         string `json:"icon"`
	CardIcon     string `json:"card_icon"`
	DrawcardIcon string `json:"drawcard_icon"`
	FilletIcon   string `json:"fillet_icon"`
	LongIcon     string `json:"long_icon"`
	BgIcon       string `json:"bg_icon"`
}

type roles struct {
	Code int `json:"code"`
	Data struct {
		AvatarCombat struct {
			ID                int     `json:"id"`
			UID               int     `json:"uid"`
			AvatarID          int     `json:"avatarId"` //角色序号
			Level             int     `json:"level"`
			Promotion         int     `json:"promotion"`   //突破次数
			HpBase            float64 `json:"hpBase"`      //hp白值
			HpFinal           float64 `json:"hpFinal"`     //总和
			AttackBase        float64 `json:"attackBase"`  //akt白值
			AttackFinal       float64 `json:"attackFinal"` //akt总和
			DefenseBase       float64 `json:"defenseBase"` //防御
			DefenseFinal      float64 `json:"defenseFinal"`
			SpeedBase         int     `json:"speedBase"` //速度
			SpeedFinal        float64 `json:"speedFinal"`
			CriticalChance    float64 `json:"criticalChance"` //暴击率
			CriticalDamage    float64 `json:"criticalDamage"` //暴击伤害
			StanceBreakRatio  int     `json:"stanceBreakRatio"`
			HealRatio         int     `json:"healRatio"`
			SpRatio           int     `json:"spRatio"`
			StatusProbability float64 `json:"statusProbability"` //效果命中
			StatusResistance  float64 `json:"statusResistance"`  //效果抵抗
			PhysicalAddHurt   float64 `json:"physicalAddHurt"`   //物理加伤
			FireAddHurt       int     `json:"fireAddHurt"`
			IceAddHurt        int     `json:"iceAddHurt"`
			ElecAddHurt       int     `json:"elecAddHurt"`
			WindAddHurt       int     `json:"windAddHurt"`
			QuantumAddHurt    int     `json:"quantumAddHurt"`
			ImaginaryAddHurt  int     `json:"imaginaryAddHurt"`
			WeaponID          int     `json:"weapon_id"`
			WeaponLevel       int     `json:"weapon_level"`
			WeaponPromotion   int     `json:"weapon_promotion"`
			WeaponRank        int     `json:"weapon_rank"`
			SkillList         []struct {
				SkillID int `json:"skill_id"`
				Level   int `json:"level"`
			} `json:"skill_list"`
		} `json:"AvatarCombat"`
		AvatarSkills []struct {
			ID               int         `json:"id"`
			CharacterID      int         `json:"character_id"`
			SkillName        string      `json:"skill_name"`
			SkillTag         string      `json:"skill_tag"`
			MaxLevel         int         `json:"max_level"`
			TreeMaxLevel     int         `json:"tree_max_level"`
			Icon             string      `json:"icon"`
			SkillDesc        string      `json:"skill_desc"`
			SkillSimpleDesc  string      `json:"skill_simple_desc"`
			ParamList        interface{} `json:"param_list"`
			StanceDamageType string      `json:"stance_damage_type"`
			AtkType          string      `json:"atk_type"`
			SkillEffect      string      `json:"skill_effect"`
			MaterialList     interface{} `json:"material_list"`
			PromotionList    interface{} `json:"promotion_list"`
		} `json:"AvatarSkills"`
		AvatarRelics []struct {
			ID             int     `json:"Id"`
			UID            int     `json:"Uid"`
			RelicID        int     `json:"RelicId"`
			AvatarID       int     `json:"AvatarId"`
			Level          int     `json:"Level"`
			Type           string  `json:"Type"`
			MainAffixType  int     `json:"MainAffixType"`
			MainAffixValue float64 `json:"MainAffixValue"`
			Sub1Type       int     `json:"Sub1Type"`
			Sub1Value      float64 `json:"Sub1Value"`
			Sub1Count      int     `json:"Sub1Count"`
			Sub1Step       int     `json:"Sub1Step"`
			Sub2Type       int     `json:"Sub2Type"`
			Sub2Value      float64 `json:"Sub2Value"`
			Sub2Count      int     `json:"Sub2Count"`
			Sub2Step       int     `json:"Sub2Step"`
			Sub3Type       int     `json:"Sub3Type"`
			Sub3Value      float64 `json:"Sub3Value"`
			Sub3Count      int     `json:"Sub3Count"`
			Sub3Step       int     `json:"Sub3Step"`
			Sub4Type       int     `json:"Sub4Type"`
			Sub4Value      float64 `json:"Sub4Value"`
			Sub4Count      int     `json:"Sub4Count"`
			Sub4Step       int     `json:"Sub4Step"`
		} `json:"AvatarRelics"`
		ItemRelic []struct {
			ID      int    `json:"Id"`
			RelicID int    `json:"RelicId"`
			Name    string `json:"Name"`
			Desc    string `json:"Desc"`
			Icon    string `json:"Icon"`
			Type    string `json:"Type"`
			Rarity  int    `json:"Rarity"`
		} `json:"ItemRelic"`
		AvatarWeapon struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			Star int    `json:"star"`
			Icon string `json:"icon"`
		} `json:"AvatarWeapon"`
	} `json:"data"`
	Msg string `json:"msg"`
}

// 本地数据
type thisdata struct {
	UID      string `json:"uid"`
	Nickname string `json:"nickname"`
	RoleData ro     `json:"data"`
}
type ro struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Star    int    `json:"star"`
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
	Promotion         int     `json:"promotion"`   //突破次数
	HpBase            float64 `json:"hpBase"`      //hp白值
	HpFinal           float64 `json:"hpFinal"`     //总和
	AttackBase        float64 `json:"attackBase"`  //akt白值
	AttackFinal       float64 `json:"attackFinal"` //akt总和
	DefenseBase       float64 `json:"defenseBase"` //防御
	DefenseFinal      float64 `json:"defenseFinal"`
	SpeedBase         int     `json:"speedBase"` //速度
	SpeedFinal        float64 `json:"speedFinal"`
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
	RelicId int     `json:"relicId"`
	Name    string  `json:"name"`
	Type    string  `json:"type"` //位置"OBJECT","HAND","BODY","FOOT","HEAD", "NECK"
	MainV   vlist   `json:"main"`
	Vlist   []vlist `json:"vlist"`
}
type vlist struct {
	Name  string `json:"valname"`
	Value string `json:"value"`
}
