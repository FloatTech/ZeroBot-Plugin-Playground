// Package qqci 简易cicd
package qqci

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/binary"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
	"github.com/FloatTech/zbputils/file"
	"github.com/FloatTech/zbputils/img/text"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

func init() {
	engine := control.Register("qqci", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Help: "简易cicd\n- /qqci -a zbp -r git@github.com:FloatTech/ZeroBot-Plugin -dir /usr/local/service -cmd \"zpb\"  -act insert\n" +
			"- /qqci -a zbp -dir D:/test -act update\n" +
			"- /qqci -a zbp -act select\n" +
			"- /qqci -a zbp -act folder -f /23bc19be-6b54-4542-b42e-d97a6bea81fd" +
			"- /qqci -a zbp -b master -act ci\n" +
			"- /qqci -a zbp -act restart\n" +
			"- /qqci -a zbp -act install\n" +
			"- /qqci -a zbp -act start\n" +
			"- /qqci -a zbp -act stop\n" +
			"- /qqci -a zbp -act upload",
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
			cmd := exec.Command("./load.sh", flagapp.Action)
			cmd.Dir = filepath.Join(app.Directory, app.Appname)
			err = cmd.Run()
			if err != nil {
				ctx.SendChain(message.Text("执行命令 ", cmd.Args, " 错误: ", err, "\n\n"))
				return
			}
			ctx.SendChain(message.Text("执行命令 ", cmd.Args, " 成功\n\n"))
		case "folder":
			var folders []folder
			text := "文件夹id: \n"
			if flagapp.Folder == "" {
				_ = json.Unmarshal(binary.StringToBytes(ctx.GetGroupRootFiles(ctx.Event.GroupID).Get("folders|@pretty").Raw), &folders)
			} else {
				_ = json.Unmarshal(binary.StringToBytes(ctx.GetGroupFilesByFolder(ctx.Event.GroupID, flagapp.Folder).Get("folders|@pretty").Raw), &folders)
			}
			for _, v := range folders {
				text += v.FolderName + ": " + v.FolderID + "\n"
			}
			ctx.SendChain(message.Text(text))
		case "upload":
			app, err = adb.getApp(flagapp)
			if err != nil {
				ctx.SendChain(message.Text("ERROR:", err))
				return
			}
			var targets []string
			maxDepth := strings.Count(filepath.Join(app.Directory, app.Appname), string(os.PathSeparator))
			_ = filepath.WalkDir(filepath.Join(app.Directory, app.Appname), func(fPath string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if d.IsDir() && strings.Count(fPath, string(os.PathSeparator)) > maxDepth {
					return fs.SkipDir
				}
				if !d.IsDir() && filepath.Base(fPath) == flagapp.Upload {
					if app.GroupID > 0 {
						ctx.UploadGroupFile(int64(app.GroupID), fPath, filepath.Base(fPath), app.Folder)
					} else {
						ctx.UploadThisGroupFile(fPath, filepath.Base(fPath), app.Folder)
					}
				}
				if !d.IsDir() {
					targets = append(targets, filepath.Base(fPath))
				}
				return nil
			})
			ctx.SendChain(message.Text("文件列表:\n- ", strings.Join(targets, "\n- ")))
		case "ci":
			ctx.Send("少女祈祷中......")
			logtext := "运行日志: \n\n"
			logtext += fmt.Sprintf("命令参数: %+v \n\n", app)
			app, err = adb.getApp(flagapp)
			if err != nil {
				logtext += fmt.Sprintf("getApp 错误: %v\n\n", err)
				data, err := text.RenderToBase64(logtext, text.FontFile, 400, 20)
				if err != nil {
					ctx.SendChain(message.Text("ERROR:", err))
					return
				}
				if id := ctx.SendChain(message.Image("base64://" + binary.BytesToString(data))); id.ID() == 0 {
					ctx.SendChain(message.Text("ERROR:可能被风控了"))
				}
				return
			}
			logtext += fmt.Sprintf("数据库参数: %+v \n\n", app)
			index := strings.LastIndex(app.Gitrepo, "/")
			if index == -1 {
				logtext += "git的地址错误\n\n"
				data, err := text.RenderToBase64(logtext, text.FontFile, 400, 20)
				if err != nil {
					ctx.SendChain(message.Text("ERROR:", err))
					return
				}
				if id := ctx.SendChain(message.Image("base64://" + binary.BytesToString(data))); id.ID() == 0 {
					ctx.SendChain(message.Text("ERROR:可能被风控了"))
				}
				return
			}
			workdir := filepath.Join(cachePath, app.Appname, app.Gitrepo[index:])
			_ = os.MkdirAll(filepath.Dir(workdir), 0755)
			if flagapp.Gitbranch == "" {
				cmd = exec.Command("git", "clone", app.Gitrepo)
				app.Gitbranch = "latest"
			} else {
				cmd = exec.Command("git", "clone", "-b", flagapp.Gitbranch, app.Gitrepo)
				app.Gitbranch = flagapp.Gitbranch
			}
			cmd.Dir = filepath.Join(cachePath, app.Appname)
			err = cmd.Run()
			if err != nil {
				_ = os.RemoveAll(workdir)
				logtext += fmt.Sprintf("执行命令 %v 错误: %v,\n\n", cmd.Args, err)
				data, err := text.RenderToBase64(logtext, text.FontFile, 400, 20)
				if err != nil {
					ctx.SendChain(message.Text("ERROR:", err))
					return
				}
				if id := ctx.SendChain(message.Image("base64://" + binary.BytesToString(data))); id.ID() == 0 {
					ctx.SendChain(message.Text("ERROR:可能被风控了"))
				}
				return
			}
			logtext += fmt.Sprintf("执行命令 %v 成功\n\n", cmd.Args)
			makefileworkdir := filepath.Join(workdir, "Makefile")
			path := filepath.Join(file.BOTPATH, engine.DataFolder(), app.Appname, "Makefile")
			err = getConfigFile(path, makefileworkdir, engine.DataFolder()+"Makefile.tpl", app)
			if err != nil {
				_ = os.RemoveAll(workdir)
				logtext += fmt.Sprintf("加载 %v 错误: %v\n\n", makefileworkdir, err)
				data, err := text.RenderToBase64(logtext, text.FontFile, 400, 20)
				if err != nil {
					ctx.SendChain(message.Text("ERROR:", err))
					return
				}
				if id := ctx.SendChain(message.Image("base64://" + binary.BytesToString(data))); id.ID() == 0 {
					ctx.SendChain(message.Text("ERROR:可能被风控了"))
				}
				return
			}
			logtext += fmt.Sprintf("加载 %v 成功\n\n", makefileworkdir)
			cmd = exec.Command("make")
			cmd.Dir = workdir
			err = cmd.Run()
			if err != nil {
				_ = os.RemoveAll(workdir)
				logtext += fmt.Sprintf("执行命令 %v 错误: %v\n\n", cmd.Args, err)
				data, err := text.RenderToBase64(logtext, text.FontFile, 400, 20)
				if err != nil {
					ctx.SendChain(message.Text("ERROR:", err))
					return
				}
				if id := ctx.SendChain(message.Image("base64://" + binary.BytesToString(data))); id.ID() == 0 {
					ctx.SendChain(message.Text("ERROR:可能被风控了"))
				}
				return
			}
			logtext += fmt.Sprintf("执行命令 %v 成功\n\n", cmd.Args)
			tarPath := filepath.Join(file.BOTPATH, workdir, "_output", app.Appname+".tar.gz")
			newtarPath := filepath.Join(app.Directory, app.Appname, app.Appname+"@"+app.Gitbranch+".tar.gz")
			_ = os.MkdirAll(filepath.Dir(newtarPath), 0755)
			tf, _ := os.Open(tarPath)
			cf, _ := os.Create(newtarPath)
			_, _ = io.Copy(cf, tf)
			cf.Close()
			tf.Close()
			if flagapp.Upload != "" {
				if app.GroupID > 0 {
					ctx.UploadGroupFile(int64(app.GroupID), newtarPath, app.Appname+"@"+app.Gitbranch+".tar.gz", app.Folder)
				} else {
					ctx.UploadThisGroupFile(newtarPath, app.Appname+"@"+app.Gitbranch+".tar.gz", app.Folder)
				}
			}
			_ = os.RemoveAll(workdir)
			loadfileworkdir := filepath.Join(app.Directory, app.Appname, "load.sh")
			if file.IsExist(loadfileworkdir) {
				cmd = exec.Command("./load.sh", "stop")
				cmd.Dir = filepath.Join(app.Directory, app.Appname)
				_ = cmd.Run()
			}
			err = deCompress(newtarPath, app.Directory)
			if err != nil {
				_ = os.RemoveAll(workdir)
				logtext += fmt.Sprintf("解压 %v 到 %v 错误: %v\n\n", newtarPath, app.Directory, err)
				data, err := text.RenderToBase64(logtext, text.FontFile, 400, 20)
				if err != nil {
					ctx.SendChain(message.Text("ERROR:", err))
					return
				}
				if id := ctx.SendChain(message.Image("base64://" + binary.BytesToString(data))); id.ID() == 0 {
					ctx.SendChain(message.Text("ERROR:可能被风控了"))
				}
				return
			}
			logtext += fmt.Sprintf("解压 %v 到 %v 成功\n\n", newtarPath, app.Directory)
			path = filepath.Join(file.BOTPATH, engine.DataFolder(), app.Appname, "load.sh")
			err = getConfigFile(path, loadfileworkdir, engine.DataFolder()+"load.tpl", app)
			if err != nil {
				logtext += fmt.Sprintf("加载 %v 文件错误: %v\n\n", loadfileworkdir, err)
				data, err := text.RenderToBase64(logtext, text.FontFile, 400, 20)
				if err != nil {
					ctx.SendChain(message.Text("ERROR:", err))
					return
				}
				if id := ctx.SendChain(message.Image("base64://" + binary.BytesToString(data))); id.ID() == 0 {
					ctx.SendChain(message.Text("ERROR:可能被风控了"))
				}
				return
			}
			logtext += fmt.Sprintf("加载 %v 文件成功\n\n", loadfileworkdir)
			_ = os.Chmod(loadfileworkdir, 0777)
			cmd = exec.Command("./load.sh", "install")
			cmd.Dir = filepath.Join(app.Directory, app.Appname)
			err = cmd.Run()
			if err != nil {
				logtext += fmt.Sprintf("执行命令 %v 错误: %v \n\n", cmd.Args, err)
				data, err := text.RenderToBase64(logtext, text.FontFile, 400, 20)
				if err != nil {
					ctx.SendChain(message.Text("ERROR:", err))
					return
				}
				if id := ctx.SendChain(message.Image("base64://" + binary.BytesToString(data))); id.ID() == 0 {
					ctx.SendChain(message.Text("ERROR:可能被风控了"))
				}
				return
			}
			logtext += fmt.Sprintf("执行命令 %v 成功\n\n", cmd.Args)
			cmd = exec.Command("./load.sh", "start")
			cmd.Dir = filepath.Join(app.Directory, app.Appname)
			err = cmd.Run()
			if err != nil {
				logtext += fmt.Sprintf("执行命令 %v 错误: %v \n\n", cmd.Args, err)
				data, err := text.RenderToBase64(logtext, text.FontFile, 400, 20)
				if err != nil {
					ctx.SendChain(message.Text("ERROR:", err))
					return
				}
				if id := ctx.SendChain(message.Image("base64://" + binary.BytesToString(data))); id.ID() == 0 {
					ctx.SendChain(message.Text("ERROR:可能被风控了"))
				}
				return
			}
			logtext += fmt.Sprintf("执行命令 %v 成功\n\n", cmd.Args)
			data, err := text.RenderToBase64(logtext, text.FontFile, 400, 20)
			if err != nil {
				ctx.SendChain(message.Text("ERROR:", err))
				return
			}
			if id := ctx.SendChain(message.Image("base64://" + binary.BytesToString(data))); id.ID() == 0 {
				ctx.SendChain(message.Text("ERROR:可能被风控了"))
			}
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
		_ = os.Chmod(filename, 0777)
	}
	return nil
}
