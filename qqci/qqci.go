// Package qqci 简易cicd
package qqci

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

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
		Help: "简易cicd\n- /qqci -a zbp -r git@github.com:FloatTech/ZeroBot-Plugin -dir /usr/local/service -cmd \"zpb\" -make data/Qqci/zbp/Makefile -load data/Qqci/zbp/load.sh -act insert\n" +
			"- /qqci -a zbp -dir D:/test -act update\n" +
			"- /qqci -a zbp -act select\n" +
			"- /qqci -a zbp -b master -act ci\n" +
			"- /qqci -a zbp -act restart\n" +
			"- /qqci -a zbp -act install\n" +
			"- /qqci -a zbp -act start\n" +
			"- /qqci -a zbp -act stop",
		PublicDataFolder: "Qqci",
	}).ApplySingle(ctxext.DefaultSingle)
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
		_ = os.MkdirAll(app.Directory, 0755)
		flagapp := ctx.State["flag"].(*application)
		switch flagapp.Action {
		case "select":
			res, err := adb.get(flagapp)
			if err != nil {
				ctx.SendChain(message.Text("ERROR:", err))
				return
			}
			ctx.SendChain(message.Text(fmt.Sprintf("执行成功,数据库参数: %+v \n\n", res)))
		case "insert":
			err = adb.insert(flagapp)
			if err != nil {
				ctx.SendChain(message.Text("ERROR:", err))
				return
			}
			ctx.SendChain(message.Text(fmt.Sprintf("执行成功,命令参数: %+v \n\n", flagapp)))
		case "update":
			err = adb.update(flagapp)
			if err != nil {
				ctx.SendChain(message.Text("ERROR:", err))
				return
			}
			ctx.SendChain(message.Text(fmt.Sprintf("执行成功,命令参数: %+v \n\n", flagapp)))
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
				ctx.SendChain(message.Text("ERROR:", flagapp.Action, ",", err))
				return
			}
		case "ci":
			ctx.Send("少女祈祷中......")
			logtext := "运行日志: \n\n"
			logtext += fmt.Sprintf("命令参数: %+v \n\n", app)
			app, err = adb.getApp(flagapp)
			if err != nil {
				logtext += fmt.Sprintf("getApp 错误: %v\n\n", err)
				ctx.SendChain(message.Text(logtext))
				return
			}
			logtext += fmt.Sprintf("数据库参数: %+v \n\n", app)
			index := strings.LastIndex(app.Gitrepo, "/")
			if index == -1 {
				logtext += "git的地址错误\n\n"
				ctx.SendChain(message.Text(logtext))
				return
			}
			workdir := cachePath + app.Gitrepo[index:]
			if flagapp.Gitbranch == "" {
				cmd = exec.Command("git", "clone", app.Gitrepo)
				app.Gitbranch = "latest"
			} else {
				cmd = exec.Command("git", "clone", "-b", flagapp.Gitbranch, app.Gitrepo)
			}
			cmd.Dir = cachePath
			err = cmd.Run()
			if err != nil {
				_ = os.RemoveAll(workdir)
				logtext += fmt.Sprintf("执行命令 %v 错误: %v,\n\n", cmd.Args, err)
				ctx.SendChain(message.Text(logtext))
				return
			}
			logtext += fmt.Sprintf("执行命令 %v 成功\n\n", cmd.Args)
			makefileworkdir := filepath.Join(workdir, "Makefile")
			path := filepath.Join(file.BOTPATH, app.MakefilePath)
			err = getConfigFile(path, makefileworkdir, engine.DataFolder()+"Makefile.tpl", app)
			if err != nil {
				_ = os.RemoveAll(workdir)
				logtext += fmt.Sprintf("加载 %v 错误: %v\n\n", makefileworkdir, err)
				ctx.SendChain(message.Text(logtext))
				return
			}
			logtext += fmt.Sprintf("加载 %v 成功\n\n", makefileworkdir)
			cmd = exec.Command("make")
			cmd.Dir = workdir
			err = cmd.Run()
			if err != nil {
				_ = os.RemoveAll(workdir)
				logtext += fmt.Sprintf("执行命令 %v 错误: %v\n\n", cmd.Args, err)
				ctx.SendChain(message.Text(logtext))
				return
			}
			logtext += fmt.Sprintf("执行命令 %v 成功\n\n", cmd.Args)
			tarPath := filepath.Join(file.BOTPATH, workdir, "_output", app.Appname+".tar.gz")
			if flagapp.Upload {
				if app.GroupID > 0 {
					ctx.UploadGroupFile(int64(app.GroupID), tarPath, app.Gitrepo+"@"+app.Gitbranch+".tar.gz", app.Folder)
				} else {
					ctx.UploadThisGroupFile(tarPath, app.Gitrepo+"@"+app.Gitbranch+".tar.gz", app.Folder)
				}
			}
			err = deCompress(tarPath, app.Directory)
			if err != nil {
				_ = os.RemoveAll(workdir)
				logtext += fmt.Sprintf("解压 %v 到 %v 错误: %V\n\n", tarPath, app.Directory, err)
				ctx.SendChain(message.Text(logtext))
				return
			}
			logtext += fmt.Sprintf("解压 %v 到 %v 成功\n\n", tarPath, app.Directory)
			_ = os.RemoveAll(workdir)
			loadfileworkdir := filepath.Join(app.Directory, app.Appname, "load.sh")
			path = filepath.Join(file.BOTPATH, app.LoadfilePath)
			err = getConfigFile(path, loadfileworkdir, engine.DataFolder()+"load.tpl", app)
			if err != nil {
				logtext += fmt.Sprintf("加载 %v 文件错误: %v\n\n", loadfileworkdir, err)
				ctx.SendChain(message.Text(logtext))
				return
			}
			logtext += fmt.Sprintf("加载 %v 文件成功\n\n", loadfileworkdir)
			_ = os.Chmod(loadfileworkdir, 0777)
			cmd = exec.Command("./load.sh", "install")
			cmd.Dir = filepath.Join(app.Directory, app.Appname)
			err = cmd.Run()
			if err != nil {
				logtext += fmt.Sprintf("执行命令 %v 错误: %v \n\n", cmd.Args, err)
				ctx.SendChain(message.Text(logtext))
				return
			}
			logtext += fmt.Sprintf("执行命令 %v 成功\n\n", cmd.Args)
			cmd = exec.Command("./load.sh", "start")
			cmd.Dir = filepath.Join(app.Directory, app.Appname)
			err = cmd.Run()
			if err != nil {
				logtext += fmt.Sprintf("执行命令 %v 错误: %v \n\n", cmd.Args, err)
				ctx.SendChain(message.Text(logtext))
				return
			}
			logtext += fmt.Sprintf("执行命令 %v 成功\n\n", cmd.Args)
			ctx.SendChain(message.Text(logtext))
		default:
			ctx.SendChain(message.Text("无效的动作:", flagapp.Action))
		}
	})
}

func getConfigFile(path, executePath, templatePath string, app application) error {
	_ = os.MkdirAll(filepath.Dir(path), 0755)
	f, err := os.Create(executePath)
	if err != nil {
		return err
	}
	if file.IsNotExist(path) {
		nf, err := os.Create(path)
		if err != nil {
			return err
		}
		tmpl, err := template.ParseFiles(templatePath)
		if err != nil {
			return err
		}
		err = tmpl.Execute(nf, app)
		if err != nil {
			return err
		}
		_ = nf.Close()
	}
	pf, err := os.Open(path)
	if err != nil {
		return err
	}
	_, _ = io.Copy(f, pf)
	_ = f.Close()
	_ = pf.Close()
	return nil
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
		filename := filepath.Join(dest, hdr.Name)
		_ = os.MkdirAll(filepath.Dir(filename), 0755)
		if hdr.Size == 0 {
			_ = os.MkdirAll(filename, 0755)
			continue
		}
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		_, _ = io.Copy(file, tr)
	}
	return nil
}
