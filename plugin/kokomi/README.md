# kokomi

kokomi是一个ZeroBot-Plugin的升级插件，提供包括角色查询等升级功能。
相比与喵喵菜单,本插件不依赖浏览器渲染,可以再树莓派等机器上运行,占用内存较低

具体功能可在安装插件后 通过 /用法kokomi 进行查看。
# 安装与更新

推荐使用git进行安装，以方便后续升级。在ZeroBot-Plugin根目录夹打开终端，运行

// 使用github

git clone https://github.com/lianhong2758/kokomi-plugin.git ./plugin/kokomi/

进行下载插件。 

然后在main.go中导入包	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/kokomi"  // kokomi原神面板

//在群文件下载

如果是手动下载的zip压缩包，请将解压后的kokomi文件夹放置在ZeroBot-Plugin目录下的plugin文件夹内。

然后在main.go中导入包	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/kokomi"  // kokomi原神面板

# 功能说明(移植喵喵菜单)
#绑定uid
#更新面板
#全部角色面板
#雷神面板

#更新面板 依赖于面板查询API，面板服务由 http://enka.network/ 提供

如果可以的话，也请在Patreon上支持Enka，或提供闲置的原神账户，具体可在Enka官网 Discord联系

国内网络如Enka服务访问不稳定，可尝试更换 @MiniGrayGay 大佬提供的中转服务
方法:将kokomi/kokomi.go中的url二级网址替换即可

    【链接1】：https://enka.microgg.cn/
    【链接2】：https://enka.minigg.cn/
# 未来可期 (以后将适配的功能)
#雷神伤害

喵喵本地计算

#雷神圣遗物

圣遗物评分为喵喵版评分规则

免责声明

    功能仅限交流技术使用，请勿将kokomi用于以盈利为目的的场景
    图片与其他素材均来借用喵喵菜单，在此感谢大佬，如有侵权请联系，会立即删除

其他

    喵喵插件 Miao-Plugin : Gitee / Github
    Enka: 感谢Enka提供的面板服务
    Snap.Genshin : 感谢 DGP Studio 开发的 胡桃API
    QQ群
       ZeroBot-Plugin官方二群609640932
