package bottle // zbp driftbottle 魔改 更干净更适用于多个群的bot使用

import (
	"fmt"
	"hash/crc64"
	"strconv"
	"sync"
	"time"

	"github.com/FloatTech/floatbox/binary"
	sql "github.com/FloatTech/sqlite"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

type sea struct {
	ID   int64  `db:"id"`   // ID qq_grp_name_msg 的 crc64 hashCheck.
	QQ   int64  `db:"qq"`   // Get current user(Who sends this)
	Name string `db:"Name"` //  his or her name at that time:P
	Msg  string `db:"msg"`  // What he or she sent to bot?
	Grp  int64  `db:"grp"`  // which group sends this msg?
	Time string `db:"time"` // we need to know the current time,master>
}

var seaSide = &sql.Sqlite{}
var seaLocker sync.RWMutex

// We need a container to inject what we need :(

func init() {
	engine := control.Register("bottle", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault:  false,
		Help:              "基于driftbottle的魔改版漂流瓶\n",
		PrivateDataFolder: "bottle",
	})
	seaSide.DBPath = engine.DataFolder() + "sea.db"
	err := seaSide.Open(time.Hour * 24)
	if err != nil {
		panic(err)
	}

	_ = CreateChannel(seaSide)
	engine.OnFullMatch("pick", zero.OnlyToMe, zero.OnlyGroup).Limit(ctxext.LimitByGroup).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		be, err := fetchBottle(seaSide)
		if err != nil {
			ctx.SendChain(message.Text("ERR:", err))
		}
		IDStr := strconv.Itoa(int(be.ID))
		QQStr := strconv.Itoa(int(be.QQ))
		GrpStr := strconv.Itoa(int(be.Grp))
		botName := zero.BotConfig.NickName[0]
		msg := make(message.Message, 0, 10)
		msg = append(msg, message.CustomNode(botName, ctx.Event.SelfID, botName+"试着帮你捞出来了这个~\nID:"+IDStr+"\n投递人: "+be.Name+"("+QQStr+")"+"\n群号: "+GrpStr+"\n时间: "+be.Time+"\n内容: \n"+be.Msg))
		ctx.SendGroupForwardMessage(ctx.Event.GroupID, msg)
	})

	engine.OnRegex(`throw.*?(.*)`, zero.OnlyToMe, zero.OnlyGroup).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		senderFormatTime := time.Unix(ctx.Event.Time, 0).Format("2006-01-02 15:04:05")
		rawSenderMessage := ctx.State["regex_matched"].([]string)[1]
		rawMessageCallBack := message.UnescapeCQCodeText(rawSenderMessage)
		// check current needs and prepare to throw bottle.
		err = globalbottle(
			ctx.Event.UserID,
			ctx.Event.GroupID,
			senderFormatTime,
			ctx.CardOrNickName(ctx.Event.UserID),
			rawMessageCallBack,
		).throw(seaSide)
		if err != nil {
			ctx.SendChain(message.Text("ERROR: ", err))
			return
		}
		ctx.Send(message.ReplyWithMessage(ctx.Event.MessageID, message.Text("已经帮你丢出去了哦~")))
	})
}

func globalbottle(qq, grp int64, time, name, msg string) *sea { // Check as if the User is available and collect information to store.
	id := int64(crc64.Checksum(binary.StringToBytes(fmt.Sprintf("%d_%d_%s_%s_%s", grp, qq, time, name, msg)), crc64.MakeTable(crc64.ISO)))
	return &sea{ID: id, Grp: grp, Time: time, QQ: qq, Name: name, Msg: msg}
}

func (be *sea) throw(db *sql.Sqlite) error {
	seaLocker.Lock()
	defer seaLocker.Unlock()
	return db.Insert("global", be)
}

func (be *sea) destory(db *sql.Sqlite) error {
	seaLocker.Lock()
	defer seaLocker.Unlock()
	return db.Del("global", "WHERE id="+strconv.FormatInt(be.ID, 10))
}

func fetchBottle(db *sql.Sqlite) (*sea, error) {
	seaLocker.Lock()
	defer seaLocker.Unlock()
	be := new(sea)
	return be, db.Pick("global", be)
}

func CreateChannel(db *sql.Sqlite) error {
	seaLocker.Lock()
	defer seaLocker.Unlock()
	return db.Create("global", &sea{})
}
