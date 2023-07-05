package utils

import (
	"github.com/FloatTech/ZeroBot-Plugin-Playground/plugin/miaox/vars"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func StringToMessageSegment(msg string) []message.MessageSegment {
	// 转换消息对象Chain
	pattern := `\[@[0-9]{5,}\]`
	r, _ := regexp.Compile(pattern)
	matches := r.FindStringSubmatch(msg)

	pos := 0
	var slice []message.MessageSegment
	for _, mat := range matches {
		qq := strings.TrimPrefix(strings.TrimSuffix(mat, "]"), "[@")
		index := strings.Index(msg[pos:], mat)
		if index < 0 {
			continue
		}

		slice = append(slice, message.Text(msg[pos:index]))
		pos = index + len(qq) + 3
		at, err := strconv.ParseInt(strings.TrimSpace(qq), 10, 64)
		if err != nil {
			continue
		}

		slice = append(slice, message.At(at))
	}

	if len(msg)-1 > pos {
		slice = append(slice, message.Text(msg[pos:]))
	}

	return slice
}

func NewDelay(ctx *zero.Ctx) *Delay {
	d := Delay{t: time.Now().Add(3 * time.Second), closed: false, ctx: ctx}
	go d.run()
	return &d
}

// 续时器
type Delay struct {
	t         time.Time
	next      bool
	closed    bool
	ctx       *zero.Ctx
	messageId *message.MessageID
}

func (d *Delay) Defer() {
	d.t = time.Now().Add(3 * time.Second)
	// 需要执行
	d.next = true
}

func (d *Delay) Close() {
	d.closed = true
	d.next = false
}

func (d *Delay) Send() {
	if d.messageId != nil {
		d.ctx.DeleteMessage(*d.messageId)
	}
	time.Sleep(500 * time.Millisecond)
	messageId := d.ctx.Send(message.ImageBytes(vars.Loading))
	d.messageId = &messageId
}

func (d *Delay) run() {
	for {
		if d.closed {
			if d.messageId != nil {
				d.ctx.DeleteMessage(*d.messageId)
			}
			return
		}

		if d.next && time.Now().After(d.t) {
			d.Send()
			d.next = false
		}
		time.Sleep(500 * time.Millisecond)
	}
}
