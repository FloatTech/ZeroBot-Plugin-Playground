// Package qqci 简易cicd
package qqci

import (
	"archive/tar"
	"compress/gzip"
	"html/template"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
	"github.com/FloatTech/zbputils/file"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

func init() {
	engine := control.Register("qqci", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Help: "简易cicd\n- /qqci -act insert -a zbp -r git@github.com:FloatTech/ZeroBot-Plugin -dir /usr/local/service -cmd \"zpb\" -make data/Qqci/zbp/Makefile -load data/Qqci/zbp/load.sh -n FloatTech\n" +
			"- /qqci -a zbp -b master\n" +
			"- /qqci -a zbp -act restart\n" +
			"- /qqci -a zbp -act install\n" +
			"- /qqci -a zbp -act start\n" +
			"- /qqci -a zbp -act stop",
		PublicDataFolder: "Qqci",
	})
	cachePath := engine.DataFolder() + "cache/"
	_ = os.RemoveAll(cachePath)
	_ = os.MkdirAll(cachePath, 0755)
	adb = initialize(engine.DataFolder() + "qqci.db")
	getdb := ctxext.DoOnceOnSuccess(func(ctx *zero.Ctx) bool {
		_, _ = engine.GetLazyData("makefile.tpl", true)
		_, _ = engine.GetLazyData("load.tpl", true)
		return true
	})
	engine.OnShell("qqci", application{}, getdb, zero.SuperUserPermission).Handle(func(ctx *zero.Ctx) {
		var (
			app application
			err error
			cmd *exec.Cmd
		)
		flagapp := ctx.State["flag"].(*application)
		switch flagapp.Action {
		case "insert":
			err = adb.insert(flagapp)
			if err != nil {
				ctx.SendChain(message.Text("ERROR:", err))
				return
			}
			ctx.SendChain(message.Text("成功,参数:", flagapp))
		case "update":
			err = adb.update(flagapp)
			if err != nil {
				ctx.SendChain(message.Text("ERROR:", err))
				return
			}
			ctx.SendChain(message.Text("成功,参数:", flagapp))
		case "install", "start", "restart", "stop":
			app, err = adb.getApp(flagapp)
			if err != nil {
				ctx.SendChain(message.Text("ERROR:", err))
				return
			}
			cmd := exec.Command("load.sh", flagapp.Action)
			cmd.Dir = filepath.Join(app.Directory, app.Appname)
			err = cmd.Run()
			if err != nil {
				ctx.SendChain(message.Text("ERROR:", err))
				return
			}
		default:
			app, err = adb.getApp(flagapp)
			if err != nil {
				ctx.SendChain(message.Text("ERROR:", err))
				return
			}
			if flagapp.Gitbranch == "" {
				cmd = exec.Command("git", "clone", app.Gitrepo)
			} else {
				cmd = exec.Command("git", "clone", "-b", flagapp.Gitbranch, app.Gitrepo)
			}
			cmd.Dir = cachePath
			err = cmd.Run()
			if err != nil {
				_ = os.RemoveAll(cachePath)
				ctx.SendChain(message.Text("ERROR:", err))
				return
			}
			index := strings.LastIndex(app.Gitrepo, "/")
			if index == -1 {
				_ = os.RemoveAll(cachePath)
				ctx.SendChain(message.Text("ERROR: git的地址错误"))
				return
			}
			workdir := cachePath + app.Gitrepo[index:]
			makefileworkdir := filepath.Join(workdir, "Makefile")
			err = getConfigFile(makefileworkdir, engine.DataFolder()+"Makefile.tpl", app)
			if err != nil {
				_ = os.RemoveAll(cachePath)
				ctx.SendChain(message.Text("ERROR:", err))
				return
			}
			cmd = exec.Command("make", "tar")
			cmd.Dir = workdir
			err = cmd.Run()
			if err != nil {
				_ = os.RemoveAll(cachePath)
				ctx.SendChain(message.Text("ERROR:", err))
				return
			}
			tarPath := filepath.Join(file.BOTPATH, "_output", app.Appname+".tar.gz")
			ctx.UploadThisGroupFile(tarPath, app.Appname+"@"+app.Gitbranch, "")
			err = deCompress(tarPath, app.Directory)
			if err != nil {
				_ = os.RemoveAll(cachePath)
				ctx.SendChain(message.Text("ERROR:", err))
				return
			}
			_ = os.RemoveAll(cachePath)
			loadfileworkdir := filepath.Join(app.Directory, app.Appname, "load.sh")
			err = getConfigFile(loadfileworkdir, engine.DataFolder()+"load.tpl", app)
			if err != nil {
				_ = os.RemoveAll(cachePath)
				ctx.SendChain(message.Text("ERROR:", err))
				return
			}
			cmd = exec.Command("load.sh", "install")
			cmd.Dir = filepath.Join(app.Directory, app.Appname)
			err = cmd.Run()
			if err != nil {
				ctx.SendChain(message.Text("ERROR:", err))
				return
			}
			cmd = exec.Command("load.sh", "start")
			cmd.Dir = filepath.Join(app.Directory, app.Appname)
			err = cmd.Run()
			if err != nil {
				ctx.SendChain(message.Text("ERROR:", err))
				return
			}
		}
	})
}

func getConfigFile(executePath, templatePath string, app application) error {
	path := filepath.Join(file.BOTPATH, app.MakefilePath)
	f, err := os.Create(executePath)
	if err != nil {
		return err
	}
	if file.IsNotExist(path) {
		tmpl, err := template.ParseFiles(templatePath)
		if err != nil {
			return err
		}
		err = tmpl.Execute(f, app)
		if err != nil {
			return err
		}
		return f.Close()
	}
	data, err := os.Open(path)
	if err != nil {
		return err
	}
	_, _ = io.Copy(f, data)
	return f.Close()
}

func deCompress(tarFile, dest string) error {
	srcFile, err := os.Open(tarFile)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	gr, err := gzip.NewReader(srcFile)
	if err != nil {
		return err
	}
	defer gr.Close()
	tr := tar.NewReader(gr)
	for {
		hdr, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}
		filename := dest + hdr.Name
		err = os.MkdirAll(string([]rune(filename)[0:strings.LastIndex(filename, "/")]), 0755)
		if err != nil {
			return err
		}
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		_, _ = io.Copy(file, tr)
	}
	return nil
}
