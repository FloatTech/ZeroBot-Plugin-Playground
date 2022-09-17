# ZeroBot-Plugin-Playground
快来这里上传你的奇思妙想吧！好的想法将会被合并入主仓库哦~
## 简介
本仓库是[ZeroBot-Plugin](https://github.com/FloatTech/ZeroBot-Plugin)的插件试验场，无需任何pr即可上传并调试自己的插件。你可以在这里验证你的想法，与其他人进行交流，不会遇到任何阻碍。
## 我该怎么做？
1. 你需要成为[FloatTech](https://github.com/FloatTech)组织的成员。这可以通过在此仓库提出issue完成。
2. 同时，我们会邀请你加入组织，并加入[maintainer](https://github.com/orgs/FloatTech/teams/maintainer)组。你需要查看邮件确认邀请。
3. 现在你已经获得了对本仓库的写权限，可以推送你的插件了！
## 约法三章
1. 每个实现特定功能的插件应该放在根目录的一个文件夹下，不要越界，否则会被移除。
2. 不要在不经允许的情况下随意改动其他人的插件，否则会被踢出[maintainer](https://github.com/orgs/FloatTech/teams/maintainer)。
3. 礼貌交谈。
## 完善教程
目前开发教程还不是很全面，欢迎大家编辑此README完善它，让更多人轻松加入zbp大家庭！

<br><br><br>

## 插件基础 (参考/复制黏贴example文件夹下JiangRed的message.go)
### **Init**
---
#### 插件默认启用状态
*不影响插件注册*
```
DisableOnDefault: true | false
```
---
#### 插件帮助 
*/用法*
```
Help: STRING
```
示例:
``` 
Help: "This is what the plugin is about", +
"And this is some extra info", +
"- ThisIsTrigger And this is description", +
"- Trigger2 | Trigger3 | Trigger4", +
"- TriggerA | TriggerB | TriggerC", +
"- TriggerABC123 description",
```
---
#### 插件数据路径 (文件夹)  
公有数据:  
*名称首字母大写*
```
PublicDataFolder: "Example"
```
私有数据:  
*名称首字母小写*
```
PrivateDataFolder: "example"
```
---
#### 插件启用/禁用触发
启用:  
*任意Handle功能* //描述肯定不准确 大概意思是任何可以被塞进handle里面的都可以正常触发
```
OnEnable: func(ctx *zero.Ctx) {
    ctx.send("MESSAGE")
    ctx.sendChain(message.At(2401223))
}
```
禁用:  
*同上*
```
OnDisable: func(ctx *zero.Ctx) {
    ctx.DeleteMessage(message.NewMessageIDFromInteger(114514))
}
```
<br>

### **Engine** //这也不太对
---
常用功能: 
```
engine.On("notice/notify/poke") //没找到相关内容, 只在chat插件里面有这个.
任意戳一戳触发
```

***未完成***