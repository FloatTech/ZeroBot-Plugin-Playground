// Package partygame 轮盘赌游戏
package partygame

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"os"
	"sync"
	"time"

	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

// Session 会话操作
type Session struct {
	GroupID    int64   // 群id
	Creator    int64   // 创建者
	Users      []int64 // 参与者
	Max        int64   // 最大人数
	Cartridges []int   // 弹夹
	IsValid    bool    // 是否有效
	ExpireTime int64   // 过期时间
	CreateTime int64   // 创建时间
}

var dataPath string

var rlmu sync.RWMutex

var dieMsg = []string{
	"很不幸, 你死了......",
	"砰...很不幸, 你死了......",
	"你死了...",
	"很不幸, 你死了......",
	"你扣下了扳机\n你死了...",
	"你拿着手枪掂了掂, 你赌枪里没有子弹\n然后很不幸, 你死了...",
	"你是一个有故事的人, 但是子弹并不想知道这些, 它只看见了白花花的脑浆\n你死了",
	"你没有想太多, 扣下了扳机。你感觉到有什么东西从你的旁边飞过, 然后意识陷入了黑暗\n你死了",
	"大多数人对自己活着并不心存感激, 但你不再是了\n你死了...",
	"你举起了枪又放下, 然后又举了起来, 你的内心在挣扎, 但是你还是扣下了扳机, 你死了...",
	"你开枪之前先去吃了杯泡面\n然后很不幸, 你死了...",
	"你对此胸有成竹, 你曾经在精神病院向一个老汉学习了用手指夹住子弹的功夫\n然后很不幸你没夹住手滑了, 死了...",
	"今天的风儿很喧嚣, 沙尘能让眼睛感到不适。你去揉眼睛的时候手枪走火, 贯穿了你的小腹。然后很不幸, 你死了...",
	"我会死吗？我死了吗？你正这样想着\n然后很不幸, 你死了...",
	"漆黑的眩晕中, 心脏渐渐窒息无力, 彻骨的寒冷将你包围\n很不幸, 你死了...",
}

var aliveMsg = []string{
	"你活了下来, 下一位",
	"你扣动扳机, 无事发生\n你活了下来",
	"你自信的扣动了扳机, 正如你所想象的那样\n你活了下来, 下一位",
	"你感觉命运女神在向你招手\n然后, 你活了下来, 下一位",
	"你吃了杯泡面发现没有调料, 你觉得不幸的你恐怕是死定了\n然后, 你活了下来, 下一位",
	"人和人的体质不能一概而论, 你在极度愤怒下, 扣下了扳机。利用扳机扣下和触发子弹的时间差, 手指一个加速硬生生扣断了它。\n然后, 你活了下来, 下一位",
	"你曾经在精神病院向一个老汉学习了用手指夹住子弹的功夫\n然后, 子弹并没有射出, 你活了下来, 下一位",
	"你曾经在精神病院向一个老汉学习过用手指夹住射出子弹的功夫, 在子弹射出的一瞬间, 你把他塞了回去\n你活了下来, 下一位",
}

func init() { // 插件主体
	engine := control.Register("roulette", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Help: "轮盘赌\n" +
			"- 创建轮盘赌\n- 加入轮盘赌\n- 开始轮盘赌\n- 开火\n- 终止轮盘赌",
		PrivateDataFolder: "roulette",
	})
	dataPath := engine.DataFolder() + "rate.json"
	_ = os.Remove(dataPath)
	checkFile(dataPath)

	checkSession := func(ctx *zero.Ctx) bool {
		ss := getSession(ctx.Event.GroupID, dataPath)
		switch ctx.Event.RawMessage {
		case "创建轮盘赌":
			if ss.GroupID == 0 {
				return true
			}
			if ss.IsValid {
				if ss.isExpire() {
					ss.close()
					return true
				}
				ctx.SendChain(message.Text("轮盘赌游戏已经开始了"))
				return false
			}
		default:
			if ss.GroupID != ctx.Event.GroupID {
				return false
			}
			if ss.IsValid {
				if ss.isExpire() {
					ctx.SendChain(message.Text("轮盘赌游戏已过期, 请重新开始"))
					ss.close()
					return false
				}
				ctx.SendChain(message.Text("轮盘赌游戏已经开始了"))
				return false
			}
		}
		return true
	}
	// 创建轮盘赌
	engine.OnFullMatch(`创建轮盘赌`, zero.OnlyGroup, checkSession).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			gid := ctx.Event.GroupID
			uid := ctx.Event.UserID

			// 创建会话
			addSession(gid, uid, dataPath)
			ctx.SendChain(message.Text("游戏开始, 目前有1位玩家, 最多还能再加入2名玩家, 发送\"加入轮盘赌\"加入游戏"))
		})

	// 加入轮盘赌
	engine.OnFullMatch("加入轮盘赌", zero.OnlyGroup, checkSession).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			gid := ctx.Event.GroupID
			uid := ctx.Event.UserID
			ss := getSession(gid, dataPath)
			if ss.checkJoin(uid) {
				ctx.SendChain(message.Text("已经在游戏中, 无法再次加入或者加入其它互动"))
				return
			}

			if ss.countUser() >= int(ss.Max) {
				ctx.SendChain(message.Text("目前已有", ss.countUser(), "位玩家, 已达人数上限, 发送\"开始轮盘赌\"进行游戏"))
				return
			}
			ss.addUser(uid)
			ctx.SendChain(message.Text("成功加入,目前已有", ss.countUser()+1, "位玩家,发送\"开始轮盘赌\"进行游戏"))
		})

	// 开始轮盘赌
	engine.OnFullMatch("开始轮盘赌", zero.OnlyGroup, checkSession).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			gid := ctx.Event.GroupID
			uid := ctx.Event.UserID

			ss := getSession(gid, dataPath)
			// 未参与不处理
			if !ss.checkJoin(uid) {
				ctx.SendChain(message.ReplyWithMessage(ctx.Event.MessageID, message.Text("你未参与游戏"))...)
				return
			}

			if ss.countUser() <= 1 {
				ctx.SendChain(message.ReplyWithMessage(ctx.Event.MessageID, message.Text("人数不足"))...)
				return
			}

			ss.IsValid = true
			rule := "轮盘容量为6, 但只填充了一发子弹, 请参与游戏的双方轮流发送`开火`, 枪响结束后"
			// 打乱参与人
			ss.rotateUser()
			// 发送游戏开始消息
			ctx.SendChain(message.Text("游戏开始,", rule, "现在请"), message.At(ss.Users[0]), message.Text("开火"))

			// 游戏进行
			stop, cancelStop := zero.NewFutureEvent("message", 8, true,
				zero.FullMatchRule("终止轮盘赌"),
				zero.AdminPermission).
				Repeat()
			defer cancelStop()
			next := zero.NewFutureEvent("message", 999, false, zero.FullMatchRule("开火"),
				zero.OnlyGroup, zero.CheckGroup(ctx.Event.GroupID))
			recv, cancel := next.Repeat()
			defer cancel()
			tick := time.NewTimer(105 * time.Second)
			after := time.NewTimer(120 * time.Second)
			for {
				select {
				case <-tick.C:
					ctx.SendChain(message.Text("轮盘赌, 还有15s过期"))
				case <-after.C:
					ctx.Send(
						message.ReplyWithMessage(ctx.Event.MessageID,
							message.Text("轮盘赌超时, 游戏结束..."),
						),
					)
					return
				case <-stop:
					ss := getSession(ctx.Event.GroupID, dataPath)
					ss.close()
					ctx.Send("轮盘赌已终止")
					return
				case c := <-recv:
					tick.Reset(105 * time.Second)
					after.Reset(120 * time.Second)
					s := getSession(gid, dataPath)
					u := c.Event.UserID
					if !s.checkJoin(u) {
						ctx.SendChain(message.ReplyWithMessage(c.Event.MessageID, message.Text("你未参与游戏"))...)
						continue
					}
					if !s.checkTurn(u) {
						ctx.SendChain(message.ReplyWithMessage(c.Event.MessageID, message.Text("未轮到你开火"))...)
						continue
					}
					if s.cartridgesLeft() == 1 {
						s.close()
						ctx.SendChain(message.Text("你长舒了一口气, 并反手击毙了"), message.At(s.Users[1]))
						c.Event.UserID = s.Users[1]
						_ = getTruthOrDare(c)
						return
					}
					if s.openFire() {
						s.close()
						ctx.SendChain(message.Text(dieMsg[rand.Intn(len(dieMsg))]))
						_ = getTruthOrDare(c)
						return
					}
					ctx.SendChain(message.Text(aliveMsg[rand.Intn(len(aliveMsg))]), message.Text(",轮到"), message.At(s.Users[1]), message.Text("开火"))
				}
			}
		})
}

func checkFile(path string) {
	rlmu.Lock()
	defer rlmu.Unlock()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		_, err := os.Create(path)
		if err != nil {
			return
		}
	}
	dataPath = path
}

func saveItem(dataPath string, info Session) {
	interact := loadSessions(dataPath)
	rlmu.Lock()
	defer rlmu.Unlock()
	if len(interact) == 0 {
		interact = append(interact, info)
	} else {
		for i, v := range interact {
			if v.GroupID == info.GroupID {
				interact[i] = info
				break
			}
		}
	}
	bytes, err := json.Marshal(&interact)
	if err != nil {
		panic(err)
	}
	// 将数据data写入文件filePath中，并且修改文件权限为755
	if err = ioutil.WriteFile(dataPath, bytes, 0644); err != nil {
		panic(err)
	}
}

func loadSessions(dataPath string) []Session {
	// 读取指定文件内容，返回的data是[]byte类型数据
	rlmu.RLock()
	defer rlmu.RUnlock()
	ss := make([]Session, 0)
	data, err := ioutil.ReadFile(dataPath)
	if err != nil {
		return ss
	}
	if err = json.Unmarshal(data, &ss); err != nil {
		return ss
	}
	return ss
}

func getSession(gid int64, dataPath string) Session {
	interact := loadSessions(dataPath)
	for _, v := range interact {
		if v.GroupID == gid {
			return v
		}
	}
	return Session{}
}

// 添加会话
func addSession(gid, uid int64, dataPath string) {
	cls := Session{}
	cls.GroupID = gid
	cls.Creator = uid
	cls.Users = append(cls.Users, uid)
	cls.IsValid = false
	cls.Max = 3
	cls.Cartridges = cls.rotateRoulette()
	cls.ExpireTime = 300
	cls.CreateTime = time.Now().Unix()

	saveItem(dataPath, cls)
}

// 获取参与人数
func (cls Session) countUser() int {
	return len(cls.Users)
}

// 加入会话
func (cls Session) addUser(userID int64) {
	cls.Users = append(cls.Users, userID)
	saveItem(dataPath, cls)
}

// 关闭
func (cls Session) close() {
	interact := loadSessions(dataPath)

	run := make([]Session, 0)
	for _, v := range interact {
		if v.GroupID == cls.GroupID {
			continue
		}
		run = append(run, v)
	}

	bytes, err := json.Marshal(&run)
	if err != nil {
		panic(err)
	}
	// 将数据data写入文件filePath中，并且修改文件权限为755
	if err = ioutil.WriteFile(dataPath, bytes, 0644); err != nil {
		panic(err)
	}
}

// 判断会话是否过期
func (cls Session) isExpire() bool {
	// 当前时间
	now := time.Now().Unix()
	// 创建时间加存活时间
	return cls.CreateTime+cls.ExpireTime < now
}

// 判断是否在队伍中
func (cls Session) checkJoin(uid int64) bool {
	// 判断是否在参与者列表中
	for _, j := range cls.Users {
		if j == uid {
			return true
		}
	}
	return false
}

// 判断是否轮到用户
func (cls Session) checkTurn(uid int64) bool {
	return cls.Users[0] == uid
}

// 剩余子弹数
func (cls Session) cartridgesLeft() int {
	return len(cls.Cartridges)
}

// 开火
func (cls Session) openFire() bool {
	// 压出头部
	bullet := cls.Cartridges[0]
	cls.Cartridges = cls.Cartridges[1:]
	if bullet == 1 {
		return true
	}
	// 获取开枪人
	user := cls.Users[0]
	// 人员轮转
	cls.Users = cls.Users[1:]
	cls.Users = append(cls.Users, user)

	saveItem(dataPath, cls)
	return false
}

// 打乱参与人顺序
func (cls Session) rotateUser() {
	// 随机打乱数组
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(cls.Users), func(i, j int) { cls.Users[i], cls.Users[j] = cls.Users[j], cls.Users[i] })
	saveItem(dataPath, cls)
}

// 旋转轮盘
func (cls Session) rotateRoulette() []int {
	// 创建6个仓位的左轮弹夹
	cartridges := []int{1, 0, 0, 0, 0, 0}
	// 随机打乱数组
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(cartridges), func(i, j int) { cartridges[i], cartridges[j] = cartridges[j], cartridges[i] })
	return cartridges
}
