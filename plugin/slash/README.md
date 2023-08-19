## Slash Bot

A Bot Running On Telegram, However it was modified for Zerobot-plugin.

### Usage

* This plugin will cause conflict with other plugins, so please make sure that the prio is the loweset.

* 此功能需要管理员  /启用slash  才可使用

Regex Patterns,没有特定需求，**会和其他以 / 为初始头的Bot冲突**:

以下是示例:

> /rua 

=> {Username}rua了他自己~

Ep: MoeMagicMango💫 rua了他自己~

> @Lucy | HafuKo💫 /rua

* 因为正则影响相关，如果@Lucy的话可能会被Matcher的阻断器卡住不会回应

> => {username} rua了 {targetName}

Ep:

=> MoeMagicMango💫 rua了 Lucy | HafuKo💫

> /rua @Lucy | Hafuko 💫

Ep:

=> MoeMagicMango💫 rua了 Lucy | HafuKo💫

> @Lucy | HafuKo💫 /rua 捏捏

Ep:

=> MoeMagicMango💫 rua了 Lucy | HafuKo💫捏捏

>/rua 捏捏 @Lucy | HafuKo💫

Ep:

=> MoeMagicMango💫 rua了 Lucy | HafuKo💫捏捏

#### 回复环境下 (在回复某个人的对话，记得去掉qq自带的@)

>  /rua

=> {username} rua了 {targetName}

Ep:

=> MoeMagicMango💫 rua了 Lucy | HafuKo💫