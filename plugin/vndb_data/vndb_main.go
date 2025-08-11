package vndbData

import (
	"errors"
	"fmt"
	"strings"

	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"

	"context"
	"encoding/json"
	"github.com/tuihub/go-vndb"
)

var tokenAPIStr = "YOUR_TOKEN"

func init() {
	engine := control.Register("vndbData", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		// 插件的简介
		Brief: "vndb",
		Help:  "- vndb搜索 [title]",
	})

	engine.OnPrefix("vndb搜索").SetBlock(true).Handle(func(ctx *zero.Ctx) {
		gameTitle := strings.TrimSpace(ctx.State["args"].(string))
		if len(gameTitle) == 0 {
			ctx.SendChain(message.Text("tokenError: 没看到标题捏~"))
			return
		}
		ctx.SendChain(message.Text("你输入的title为: ", gameTitle, "\n现在开始搜索。"))

		client := vndb.New(
			vndb.WithToken(tokenAPIStr),
		)
		getErr := searchByName(ctx, client, gameTitle)
		if getErr != nil {
			ctx.SendChain(message.Text("vndbError: ", getErr))
			return
		}
		ctx.SendChain(message.Text("获取到的信息如上，请输入【vndbID搜索 vid】\n例如：vndbID搜索 v12" +
			"\n如果上面没有你想要的游戏，请再次搜索，并使用全名搜索。"))
		ctx.SendChain(message.Text("如果需要查找某个游戏的汉化情况，请输入【查询汉化信息 vid】\n例如：查询汉化信息 v12"))

		engine.OnPrefix("vndbID搜索").SetBlock(true).Handle(func(ctx *zero.Ctx) {
			gameId := strings.TrimSpace(ctx.State["args"].(string))
			if len(gameId) == 0 {
				ctx.SendChain(message.Text("IDError: 没看到ID捏~"))
				return
			}
			ctx.SendChain(message.Text("你输入的id为: ", gameId, "\n现在开始搜索。"))

			getErr := searchById(ctx, client, gameId)
			if getErr != nil {
				ctx.SendChain(message.Text("vndbError: ", getErr))
				return
			}
		})

		engine.OnPrefix("查询汉化信息").SetBlock(true).Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(message.Text("信息查询中......"))
			gameId := strings.TrimSpace(ctx.State["args"].(string))
			if len(gameId) == 0 {
				ctx.SendChain(message.Text("IDError: 没看到ID捏~"))
				return
			}
			ctx.SendChain(message.Text("你输入的id为: ", gameId, "\n现在开始搜索。"))

			getErr := searchTranslateInfo(ctx, client, gameId)
			if getErr != nil {
				ctx.SendChain(message.Text("vndbError: ", getErr))
				return
			}
		})
		return
	})
}

func decodeData(data any) (map[string]interface{}, error) {
	jsonData, jsonErr := json.MarshalIndent(data, "", "\t")
	if jsonErr != nil {
		return nil, jsonErr
	}
	var outData map[string]interface{}
	jsonErr2 := json.Unmarshal(jsonData, &outData)
	if jsonErr2 != nil {
		return nil, jsonErr2
	}
	return outData, nil
}

func removeDuplicates(s []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, v := range s {
		if _, ok := seen[v]; !ok {
			seen[v] = true
			result = append(result, v)
		}
	}
	return result
}

func searchByName(ctx *zero.Ctx, client *vndb.Vndb, titleName string) error {
	var ret string
	res, err := client.Vn(context.Background(), vndb.Request{
		Filters: []string{"search", "=", titleName},
		Fields:  "titles.title",
	})
	if err != nil {
		return err
	}
	data, err1 := decodeData(res)
	if err1 != nil {
		return err1
	}
	proTop := data["results"].([]interface{})
	if len(proTop) == 0 {
		errMsg := errors.New("没查找到条目。")
		return errMsg
	}
	for _, p := range proTop {
		proInner := p.(interface{})
		innerData, err2 := decodeData(proInner)
		if err2 != nil {
			return err2
		}
		proId := innerData["id"].(string)
		proTitles := innerData["titles"].([]interface{})
		ret += proId
		ret += "\n"
		for _, t := range proTitles {
			pT := t.(interface{})
			pTitle, err3 := decodeData(pT)
			if err3 != nil {
				return err3
			}
			dTitle := pTitle["title"].(string)
			ret += "\t" + dTitle + "\n"
		}
	}
	ctx.SendChain(message.Text(ret))
	return nil
}

func searchById(ctx *zero.Ctx, client *vndb.Vndb, ID string) error {
	res, err := client.Vn(context.Background(), vndb.Request{
		Filters: []string{"id", "=", ID},
		Fields: "title, alttitle, released, devstatus, platforms, image.url, developers.name, developers.original, " +
			"staff.name, staff.original, va.staff.name, va.staff.original",
	})
	if err != nil {
		return err
	}
	data, err1 := decodeData(res)
	if err1 != nil {
		return err1
	}

	ret := ""
	pImg := ""
	pTop := data["results"].([]interface{})
	for _, p := range pTop {
		pInner := p.(interface{})
		innerData, err2 := decodeData(pInner)
		if err2 != nil {
			return err2
		}
		pTitle := innerData["alttitle"]
		if pTitle == nil {
			pTitle = innerData["title"]
		}
		dTitle := pTitle.(string)
		ret += dTitle + "\n【状态】："
		pDevstatus := innerData["devstatus"].(float64)
		if pDevstatus == 0 {
			ret += "开发完成"
		} else if pDevstatus == 1 {
			ret += "开发中"
		} else if pDevstatus == 2 {
			ret += "停止开发"
		} else {
			ret += "状态不明"
		}
		ret += "\n【发售日期】："

		pRelease := innerData["released"]
		if pRelease == nil {
			err4Msg := errors.New("无效的发布日期")
			return err4Msg
		}
		ret += pRelease.(string) + "\n平台："

		pPlatforms := innerData["platforms"].([]interface{})
		for _, pPlat := range pPlatforms {
			pP := pPlat.(interface{})
			ret += pP.(string) + " | "
		}
		ret += "\n【发售商】："
		pDevData := innerData["developers"].([]interface{})
		for _, pDev := range pDevData {
			dP := pDev.(interface{})
			dDData, err3 := decodeData(dP)
			if err3 != nil {
				return err3
			}
			pName := dDData["name"]
			pOriginal := dDData["original"]
			if pOriginal == nil {
				ret += pName.(string) + " | "
			} else {
				ret += pOriginal.(string) + " | "
			}
		}
		ret += "\n【CV】："
		vaTop := innerData["va"].([]interface{})
		for _, vaTopItem := range vaTop {
			vaP := vaTopItem.(interface{})
			vaPData, err3 := decodeData(vaP)
			if err3 != nil {
				return err3
			}
			vaStaff := vaPData["staff"].(map[string]interface{})
			pName := vaStaff["name"]
			pOriginal := vaStaff["original"]
			if pOriginal == nil {
				ret += pName.(string) + " | "
			} else {
				ret += pOriginal.(string) + " | "
			}
		}
		ret += "\n【Staff】："
		pStaff := innerData["staff"].([]interface{})
		var pStaffMap []string
		for _, pStaffItem := range pStaff {
			pStaffinner := pStaffItem.(interface{})
			pStaffData, err3 := decodeData(pStaffinner)
			if err3 != nil {
				return err3
			}
			pName := pStaffData["name"]
			pOriginal := pStaffData["original"]
			if pOriginal == nil {
				pStaffMap = append(pStaffMap, pName.(string))
			} else {
				pStaffMap = append(pStaffMap, pOriginal.(string))
			}
		}
		dStaffMap := removeDuplicates(pStaffMap)
		for _, dStf := range dStaffMap {
			ret += dStf + " | "
		}
		pImages := innerData["image"]
		if pImages != nil {
			pImgD := pImages.(map[string]interface{})
			pImg = pImgD["url"].(string)
		}
		ret += "【\n详细信息请访问：https://vndb.org/" + ID + "】"
	}
	ctx.SendChain(message.Text(ret))
	if len(pImg) != 0 {
		ctx.SendChain(message.Image(pImg))
	}
	return nil
}

func searchTranslateInfo(ctx *zero.Ctx, client *vndb.Vndb, Id string) error {
	// Get ID
	filtersStr := "[\"and\", [\"or\", [\"lang\", \"=\", \"zh-Hans\"],[\"lang\", \"=\", \"zh-Hant\"]], [\"vn\", \"=\", [\"and\", [\"id\", \"=\", \"%s\"]]]]"
	sentMsg := fmt.Sprintf(filtersStr, Id)
	res, idErr := client.Release(context.Background(), vndb.Request{
		Filters: sentMsg,
		Fields: "languages.lang, released, languages.title , official, languages.mtl, platforms, vns.title, vns.alttitle, vns.image.url, " +
			"vns.released, vns.developers.name, vns.devstatus, producers.name, producers.name, producers.original",
		Results: 100,
	})
	if idErr != nil {
		return idErr
	}
	ret := ""
	data, err1 := decodeData(res)
	if err1 != nil {
		return err1
	}
	pTop := data["results"].([]interface{})
	if len(pTop) == 0 {
		errMsg := errors.New("没有汉化条目。")
		return errMsg
	}
	for _, p := range pTop {
		ret += "\n======================================\n"
		pInner := p.(interface{})
		innerData, err2 := decodeData(pInner)
		if err2 != nil {
			return err2
		}
		pVns := innerData["vns"].([]interface{})
		for _, pV := range pVns {
			//titles
			pVData, err3 := decodeData(pV)
			if err3 != nil {
				return err3
			}
			ret += "【原始标题】："
			pVnsTittle := pVData["alttitle"]
			if pVnsTittle == nil {
				pVnsTittle = pVData["title"]
			}
			dVnsSTitle := pVnsTittle.(string)
			ret += dVnsSTitle + "\n"
			//developers
			ret += "【制作者】："
			pDevelopers := pVData["developers"].([]interface{})
			for _, pDev := range pDevelopers {
				pDevData := pDev.(map[string]interface{})
				dDevData, err4 := decodeData(pDevData)
				if err4 != nil {
					return err4
				}
				developersNames := dDevData["name"].(string)

				ret += developersNames + " | "
			}
			ret += "\n【平台】："
			pPlatforms := innerData["platforms"].([]interface{})
			for _, pPlat := range pPlatforms {
				pP := pPlat.(interface{})
				ret += pP.(string) + " | "
			}
			ret += "\n"
			//devstatus
			ret += "【开发状态】："
			pDevstatus := pVData["devstatus"].(float64)
			if pDevstatus == 0 {
				ret += "开发完成"
			} else if pDevstatus == 1 {
				ret += "开发中"
			} else if pDevstatus == 2 {
				ret += "停止开发"
			} else {
				ret += "状态不明"
			}
			ret += "\n"
			//released
			pReleased := pVData["released"]
			if pReleased == nil {
				err4Msg := errors.New("无效的发布日期")
				return err4Msg
			}
			ret += "【在" + pReleased.(string) + "时发售】"
			ret += "\n"
		}
		// TranslateLanguage
		pTitleName := ""
		pLanguage := innerData["languages"].([]interface{})
		for _, l := range pLanguage {
			pLangInfo := l.(interface{})
			pLangData, err3 := decodeData(pLangInfo)
			if err3 != nil {
				return err3
			}
			pLang := pLangData["lang"].(string)
			if pLang == "zh-Hans" || pLang == "zh-Hant" {
				pLanguageTitle := pLangData["title"].(string)
				pTitleName = pLanguageTitle
			}
		}
		ret += "【翻译标题】：" + pTitleName + "\n"
		ret += "【语言】："
		pLanguages := innerData["languages"].([]interface{})
		for _, pL := range pLanguages {
			pLang := pL.(map[string]interface{})
			dLang := pLang["lang"].(string)
			dLanguage := ""
			if dLang == "zh-Hans" {
				dLanguage = "简体中文"
			} else if dLang == "zh-Hant" {
				dLanguage = "繁体中文"
			} else if dLang == "ja" {
				dLanguage = "日语"
			} else if dLang == "en" {
				dLanguage = "英语"
			} else if dLang == "ko" {
				dLanguage = "韩语"
			} else if dLang == "ru" {
				dLanguage = "俄语"
			} else {
				dLanguage = dLang
			}
			ret += dLanguage + " | "
		}
		ret += "\n"
		ret += "【是否机翻】："
		isMtlStr := ""
		pLang := innerData["languages"].([]interface{})
		for _, pL := range pLang {
			pLangData := pL.(map[string]interface{})
			isMtl := pLangData["mtl"].(bool)
			if isMtl == true {
				isMtlStr = "是"
			} else {
				isMtlStr = "否"
			}
		}
		ret += isMtlStr + "\n【是否官中】："
		pOfficial := innerData["official"].(bool)
		isOffical := ""
		if pOfficial == true {
			isOffical = "是"
		} else {
			isOffical = "否"
		}
		ret += isOffical + "\n【翻译组/个人】："
		pProducers := innerData["producers"].([]interface{})
		for _, pP := range pProducers {
			pPData, err4 := decodeData(pP)
			if err4 != nil {
				return err4
			}
			pName := pPData["name"]
			pOriginal := pPData["original"]
			if pOriginal == nil {
				ret += pName.(string) + " | "
			} else {
				ret += pOriginal.(string) + " | "
			}
		}
		ret += "\n"
		pReleased := innerData["released"].(string)
		if pReleased == "" {
			err4Msg := errors.New("无效的发布日期")
			return err4Msg
		}
		if pReleased == "TBA" {
			ret += "【制作中......】"
		} else {
			ret += "【在" + pReleased + "时制作完成并发布】"
		}

	}
	ctx.SendChain(message.Text(ret))
	return nil
}
