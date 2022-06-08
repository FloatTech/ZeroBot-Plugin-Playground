// Package midicreate 简易midi音乐制作
package midicreate

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
	"github.com/FloatTech/zbputils/file"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"gitlab.com/gomidi/midi/gm"
	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/smf"
)

func init() {
	engine := control.Register("midicreate", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault:  false,
		Help:              "midi音乐制作\n- midi制作 CCGGAAGR FFEEDDCR GGFFEEDR GGFFEEDR CCGGAAGR FFEEDDCR",
		PrivateDataFolder: "midicreate",
	})
	cachePath := engine.DataFolder() + "cache/"
	go func() {
		_ = os.RemoveAll(cachePath)
		err := os.MkdirAll(cachePath, 0755)
		if err != nil {
			panic(err)
		}
	}()
	engine.OnRegex(`^midi制作\s?(.{1,1000})$`).SetBlock(true).Limit(ctxext.LimitByUser).
		Handle(func(ctx *zero.Ctx) {
			uid := ctx.Event.UserID
			input := ctx.State["regex_matched"].([]string)[1]
			midiFile := cachePath + strconv.FormatInt(uid, 10) + time.Now().Format("20060102150405") + "_midicreate.mid"
			err := mkMidi(midiFile, input)
			if err != nil {
				ctx.SendChain(message.Text("ERROR:", err))
				return
			}
			cmidiFile := strings.ReplaceAll(midiFile, ".mid", ".wav")
			result, err := command("timidity " + file.BOTPATH + "/" + midiFile + " -Ow -o " + file.BOTPATH + "/" + cmidiFile)
			if err != nil {
				ctx.SendChain(message.Text("ERROR:", err))
				_ = ctx.CallAction("upload_group_file", zero.Params{"group_id:": ctx.Event.GroupID, "file": file.BOTPATH + "/" + midiFile, "name": filepath.Base(midiFile)})
				return
			}
			if !strings.Contains(result, "Notes lost totally: 0") {
				ctx.SendChain(message.Text("也许你需要安装timidity,用于midi转wav格式,result:", result))
				_ = ctx.CallAction("upload_group_file", zero.Params{"group_id:": ctx.Event.GroupID, "file": file.BOTPATH + "/" + midiFile, "name": filepath.Base(midiFile)})
				return
			}
			ctx.SendChain(message.Record("file:///" + file.BOTPATH + "/" + cmidiFile))
		})
}

var (
	noteMap = map[string]uint8{
		"C":  60,
		"Db": 61,
		"D":  62,
		"Eb": 63,
		"E":  64,
		"F":  65,
		"Gb": 66,
		"G":  67,
		"Ab": 68,
		"A":  69,
		"Bb": 70,
		"B":  71,
	}
)

func mkMidi(filePath, input string) error {
	if file.IsExist(filePath) {
		return nil
	}
	var (
		bf    bytes.Buffer
		clock = smf.MetricTicks(96)
		tr    smf.Track
	)

	tr.Add(0, smf.MetaMeter(4, 4))
	tr.Add(0, smf.MetaTempo(60))
	tr.Add(0, smf.MetaInstrument("Violin"))
	tr.Add(0, midi.ProgramChange(0, gm.Instr_Violin.Value()))

	k := strings.ReplaceAll(input, " ", "")

	var (
		base        uint8
		level       uint8
		delay       uint32
		sleepFlag   bool
		lengthBytes = make([]byte, 0)
	)

	for i := 0; i < len(k); {
		base = 0
		level = 0
		sleepFlag = false
		lengthBytes = lengthBytes[:0]
		for {
			switch {
			case k[i] == 'R':
				sleepFlag = true
				i++
			case k[i] >= 'A' && k[i] <= 'G':
				base = noteMap[k[i:i+1]] % 12
				i++
			case k[i] == 'b':
				base--
				i++
			case k[i] == '#':
				base++
				i++
			case k[i] >= '0' && k[i] <= '9':
				level = level*10 + k[i] - '0'
				i++
			case k[i] == '<':
				i++
				for i < len(k) && (k[i] == '-' || (k[i] >= '0' && k[i] <= '9')) {
					lengthBytes = append(lengthBytes, k[i])
					i++
				}
			default:
				return fmt.Errorf("无法解析第%d个位置的%c字符", i, k[i])
			}
			if i >= len(k) || (k[i] >= 'A' && k[i] <= 'G') || k[i] == 'R' {
				break
			}
		}
		length, _ := strconv.Atoi(string(lengthBytes))
		if sleepFlag {
			if length >= 0 {
				delay = clock.Ticks4th() * (1 << length)
			} else {
				delay = clock.Ticks4th() / (1 << -length)
			}
			continue
		}
		if level == 0 {
			level = 5
		}
		tr.Add(delay, midi.NoteOn(0, o(base, level), 120))
		if length >= 0 {
			tr.Add(clock.Ticks4th()*(1<<length), midi.NoteOff(0, o(base, level)))
		} else {
			tr.Add(clock.Ticks4th()/(1<<-length), midi.NoteOff(0, o(base, level)))
		}
		delay = 0
	}
	tr.Close(0)

	s := smf.New()
	s.TimeFormat = clock
	err := s.Add(tr)
	if err != nil {
		return err
	}
	_, err = s.WriteTo(&bf)
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, bf.Bytes(), 0666)
}

func o(base uint8, oct uint8) uint8 {
	if oct > 10 {
		oct = 10
	}

	if oct == 0 {
		return base
	}

	res := base + 12*oct
	if res > 127 {
		res -= 12
	}

	return res
}
