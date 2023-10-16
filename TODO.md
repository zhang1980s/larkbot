## Deployment
- [x] CDK部署
apigw/lambda/ddb/ssm secrets/eventbridge rule 资源
支持默认region，默认account参数。

- [x] modify go runtime to custom runtime
provided.al2 on arm64

- [x] 不使用fixed的ddb表名

- [x] lambda version/alias

- [ ] stackset自动创建role
带support API role的org集中管理模式 （集成在larkbot app？ 分开的app？）

- [ ] Makefile
build lambda(s) code in single enter?

- [ ] Codepipeline
应用升级管理

- [ ] CDK代码优化
封装，抽象

## functions
- [x] go runtime update
1.17 -> 1.20

- [x] aws sdk update
aws-sdk-go-v2 (current version)

- [x] 替换不稳定老版本飞书API
https://github.com/larksuite/oapi-sdk-go
https://open.feishu.cn/document/home/index

send msg, send card， get contact已更新为sdk

未使用飞书sdk的接口：数据结构和接口url未变更
	downloadUrl      = "https://open.feishu.cn/open-apis/im/v1/messages/%s/resources/%s?type=%s"
	tokenUrl         = "https://open.feishu.cn/open-apis/auth/v3/tenant_access_token/internal/"
	createChatTabUrl = "https://open.feishu.cn/open-apis/im/v1/chats/%s/chat_tabs"

（官方没有API生命周期说明，未来会不会再变？）

- [x] assume role mode
通过获取目标role（account ID，role name）获取权限；
把目标账号的role写在ddb config表里面

- [ ] 拆分lambda（restAPI）为3个： 
  1. 请求事件地址API+核心lambda（核心逻辑）；（已完成）
https://open.feishu.cn/document/server-docs/event-subscription-guide/overview
  2. 消息卡片请求地址API+消息卡片处理lambda（卡片处理）；
https://open.feishu.cn/document/server-docs/im-v1/message-card/overview
  3. 周期性刷新case逻辑lambda;
数据结构未变，裁剪

- [ ] at bot in non-case group, (非工单群需要at，工单群不需要at -- 附件上传功能需要确认)
场景待确认

- [ ] Get user information from Contact (通讯录)
user_id 和用户名是否是一一对应关系？
https://open.feishu.cn/document/server-docs/contact-v3/user/get?appId=cli_a45b7fdd6cf8100b
https://open.feishu.cn/document/server-docs/contact-v3/user/batch_get_id?appId=cli_a45b7fdd6cf8100b


- [x] 增加可以使用机器人的飞书用户的白名单表 
https://open.feishu.cn/document/faq/trouble-shooting/how-to-obtain-user-id
当前用user_id, 是否有可读性比较好的维护方式？

- [x] store app key/app secret to secret manager
cdk部署时加入输入appkey和secret的参数,保存在ssm

- [x] 扫描历史工单功能

## ddb 数据生命周期管理
- [ ] 快速过期audit表的数据
CDK实现

## Instrument 
- [ ] xray-go
https://github.com/aws/aws-xray-sdk-go

## WAF
还没想好～