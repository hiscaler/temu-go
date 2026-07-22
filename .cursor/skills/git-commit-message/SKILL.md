---
name: git-commit-message
description: >-
  Drafts Git commit messages for this repo using Sign: Message (New/Chg/Enh/Bug/Doc)
  from staged diffs, with user confirmation before commit. Use when the user asks to
  撰写提交消息、git message、commit message
---

根据以下规则来撰写 Git 提交消息

## 消息格式

提交格式为 `Sign: Message`，提交内容过长的情况下采用结构化书写（单行摘要+空行+详细描述）

Sign 有以下几种方式：

- New: 添加了新功能
- Chg: 修改了某个已经存在的功能
- Enh: 优化或者重构了已有代码
- Bug: 修复了原有代码的问题
- Doc: 修改了文档、属性、方法注释

示例：New: 支持 Excel 读取

## 注意事项

Git 处理只针对 Staged 中的文件，如果 Staged 中没有暂存文件，直接跳过，并提示用户

- 仅提交 Staged 中的文件
- 如果发现文件用户修改过了，以用户的修改结果为准
- 添加了新功能，应该读取方法注释，说明方法含义
- 修改了调用方法的参数，应该读取参数注释，说明参数含义
- 提交前应该**格式化修改过的代码**，当前未修改的代码不需要格式化
- 读取并分析 staged 暂存文件修改内容，按照 Git Message 提交格式撰写提交内容，撰写完毕后提示用户确认，待用户确认完毕后再提交
- 不要撰写无意义的消息，比如修改了 10 个文件

