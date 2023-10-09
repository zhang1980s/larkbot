# larkbot-app

## Architecture

## Deployment 

## Bot Configuration 

## Code Structure

### lambda 

目录结构：

config: 基于lambda环境变量CFG_KEY的定义，获取bot_config表中的机器人配置信息
dao: 数据访问对象。定义访问Dynamodb和AWS support API的相关接口
model: 实体模型。飞书事件及飞书卡片数据结构定义
service: 核心逻辑
utils: 公共函数集


1. 事件订阅逻辑

接收消息：
https://open.feishu.cn/document/server-docs/im-v1/message/events/receive



2. todo

### DynamoDB