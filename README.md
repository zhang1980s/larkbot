# 企业飞书即时通信工具AWS工单系统接入方案
---
飞书AWS是一套基于飞书企业通信工具的方便用户和AWS售后工程师快捷文字沟通的工具。飞书用户可以通过简单的机器人关键字和飞书小卡片互动，向AWS售后工程师团队提交支持案例，更新案例内容，以及准实时接收来自后台工程师的更新。

## 架构图

![架构示意图](picture/larkbot_architecture_3.0.png)


## 操作方式
---



## 部署和配置
---
机器人通过SDK实现自动化部署及更新。部署和配置过程主要包含四个步骤，分别是：

1. CDK 部署机器人后端资源
2. 在飞书开放平台创建自定义机器人应用，设置消息卡片请求地址及事件订阅请求地址
3. 在DynamoDB中配置机器人的自定义参数
4. 创建SupportAPI角色

下面分别介绍每个步骤的详细操作方式。

#### CDK部署机器人后端
---

飞书机器人可以选择在一个AWS账号中部署，通过assume role的方式调用其他账号（包括本账号）的role进行support API操作。部署账号本身没有特殊要求。由于飞书服务器端在国内，并且发送请求时对回调地址（机器人服务器端）的请求有响应时间的要求，因此应该尽量选择距离国内较近的region。实测发现飞书服务器端到日本和新加坡Region的延迟相对较低，建议选在在这两个region部署。


0. 安装CDK工具

使用下面命令安装cdk工具。

```
npm install -g aws-cdk
```
参考下面官方文档安装cdk工具
https://docs.aws.amazon.com/cdk/v2/guide/getting_started.html


1. 从官方仓库下载源代码及lambda二进制文件（适用于没有go开发环境的部署环境）

```
git clone https://github.com/zhang1980s/larkbot.git
```


2. 初始化CDK部署环境（如果当前Region之前没有初始化过CDK环境）

```
$ cdk bootstrap aws://<accountID>/<region> --profile <profile>
```

用当前部署的账号ID和地区以及当前环境中的AWS profile （如有必要）替代上面的示例命令中的对应项目。

例如，下面命令在123456789012账号的ap-northeast-1地区，使用global profile初始化cdk环境。

```
$ cdk bootstrap aws//123456789012/ap-northeast-1 --profile global

  cdk bootstrap aws://123456789012/ap-northeast-1 --profile global
 ⏳  Bootstrapping environment aws://123456789012/ap-northeast-1...
Trusted accounts for deployment: (none)
Trusted accounts for lookup: (none)
Using default execution policy of 'arn:aws:iam::aws:policy/AdministratorAccess'. Pass '--cloudformation-execution-policies' to customize.
CDKToolkit: creating CloudFormation changeset...
 ✅  Environment aws://123456789012/ap-northeast-1 bootstrapped.
```

该命令会通过cloudformation创建用于cdk部署的相应iam policy，role以及用于存储状态数据的s3存储桶。

参考文档：

什么是CDK
https://docs.aws.amazon.com/cdk/v2/guide/home.html

Bootstrapping CDK
https://docs.aws.amazon.com/cdk/v2/guide/bootstrapping.html


3. 通过CDK部署飞书机器人后端环境

进入larkbot仓库的主目录，执行cdk-deploy-to.sh 脚本在指定账号的指定Region中部署飞书机器人后端环境。cdk命令会通过cloudformation的方式创建相关资源及对应的最小权限关系。


```
$ cd larkbot
$ ./cdk-deploy-to.sh <accountID> <region> --context stackName=<stackname> --profile <profile>
```

上面命令通过输入 --context stackName 参数自定义stackName，如果未输入此参数，飞书机器人会使用默认的"LarkbotAppStack"作为cloudformation的stack名称。

例如，下面命令在123456789012账号的ap-northeast-1地区，使用global profile创建飞书机器人后端。

 ```
 ./cdk-deploy-to.sh 123456789012 ap-northeast-1 --profile global --context stackName=larkbot

✨  Synthesis time: 7.73s

larkbot:  start: Building 8823c5122e6d34f5b8f013ff748df0c0e2f8d78e7d6fcb8e5dd9863f5f31cc95:123456789012-ap-northeast-1
larkbot:  success: Built 8823c5122e6d34f5b8f013ff748df0c0e2f8d78e7d6fcb8e5dd9863f5f31cc95:123456789012-ap-northeast-1
larkbot:  start: Building 63a30f564d7b72bdec248adf1074770947b5356568f272138db30aa8d7c781cc:123456789012-ap-northeast-1
larkbot:  success: Built 63a30f564d7b72bdec248adf1074770947b5356568f272138db30aa8d7c781cc:123456789012-ap-northeast-1
larkbot:  start: Publishing 8823c5122e6d34f5b8f013ff748df0c0e2f8d78e7d6fcb8e5dd9863f5f31cc95:123456789012-ap-northeast-1
larkbot:  start: Publishing 63a30f564d7b72bdec248adf1074770947b5356568f272138db30aa8d7c781cc:123456789012-ap-northeast-1
larkbot:  success: Published 63a30f564d7b72bdec248adf1074770947b5356568f272138db30aa8d7c781cc:123456789012-ap-northeast-1
larkbot:  success: Published 8823c5122e6d34f5b8f013ff748df0c0e2f8d78e7d6fcb8e5dd9863f5f31cc95:123456789012-ap-northeast-1
This deployment will make potentially sensitive changes according to your current security approval level (--require-approval broadening).
Please confirm you intend to make the following modifications:

...
(NOTE: There may be security-related changes not in this list. See https://github.com/aws/aws-cdk/issues/1299)

Do you wish to deploy these changes (y/n)? y
LarkbotAppStack (larkbot): deploying... [1/1]
larkbot: creating CloudFormation changeset...


...

 ✅  LarkbotAppStack (larkbot)

✨  Deployment time: 133.44s

Outputs:
LarkbotAppStack.msgEventRoleArn = arn:aws:iam::123456789012:role/larkbot-larkbotmsgeventServiceRoleC3080B6B-V1ESZLK7ODYY
LarkbotAppStack.msgEventapiEndpointAC31EC6D = https://t68l424zt0.execute-api.ap-northeast-1.amazonaws.com/prod/
Stack ARN:
arn:aws:cloudformation:ap-northeast-1:123456789012:stack/larkbot/b35ccb10-6c3a-11ee-bef1-02e3082fe481

✨  Total time: 141.17s

 ```

参考文档：
https://docs.aws.amazon.com/cdk/v2/guide/environments.html


在cdk部署完成后，程序会输出下面两个参数，保存这两个参数的输出，后面配置飞书自定义机器人应用及设置支持support API的role时会用到。

```
Outputs:
LarkbotAppStack.msgEventRoleArn = arn:aws:iam::123456789012:role/larkbot-larkbotmsgeventServiceRoleC3080B6B-V1ESZLK7ODYY
LarkbotAppStack.msgEventapiEndpointAC31EC6D = https://t68l424zt0.execute-api.ap-northeast-1.amazonaws.com/prod/
```

4. 删除飞书机器人后端（如果必要）

登陆飞书机器人部署的AWS账号，选择部署地区，进入cloudformation服务界面，可以看到对应的CloudFormation stack，删除此stack。如果之前已经使用了飞书机器人，还需要清理机器人产生的Cloudwatch log group以避免额外的费用。




#### 创建自定义机器人应用
---

1. 访问飞书开放平台https://open.feishu.cn, 确认已经使用飞书账号登陆飞书开放平台

2. 在页面的右上角点击开发者后台，然后在开发者后台主页中，点击创建企业自建应用按钮

![创建企业自建应用-1](picture/open-feishu-cn-1.png)

3. 在创建企业自建应用页面中输入应用的名称，应用描述，应用图标及背景色，选择完毕后点击创建

![创建企业自建应用-2](picture/create-custom-app-1.png)

4. 在添加应用能力页面中，点击机器人选项框左下角中的添加按钮
![添加应用能力](picture/create-custom-app-2.png)

5. 在机器人配置的主页中，点击机器人配置标题右侧的粉笔按钮编辑机器人配置

![机器人配置主页](picture/bot-config-1.jpeg)

6. 在消息卡片请求网址的选项框中，输入CDK程序部署完成后输出的msgEventapiEndpoint的URL，并且添加"/messages"路径

![消息卡片请求网址](picture/bot-config-msg-card-address.jpeg)

例如：
```
Outputs:
LarkbotAppStack.msgEventRoleArn = arn:aws:iam::123456789012:role/larkbot-larkbotmsgeventServiceRoleC3080B6B-V1ESZLK7ODYY
LarkbotAppStack.msgEventapiEndpointAC31EC6D = https://t68l424zt0.execute-api.ap-northeast-1.amazonaws.com/prod/

```

添加/messages路径：

```
https://t68l424zt0.execute-api.ap-northeast-1.amazonaws.com/prod/messages
```


点击验证没有任何输出表示飞书机器人后端响应正常。如果提示“请求URL验证未通过”，需要检查URL格式是否正常，或者选择距离国内更近的AWS Region部署机器人应用。

（飞书会向该URL发送一个challenge值并且要求1s回复challenge的值，如果无法及时返回则提示请求URL验证未通过）

7. 点击页面左侧中开发配置段落中的事件订阅功能，在事件订阅功能中配置请求地址，输入和消息卡片请求网址相同的URL


![事件订阅请求地址配置](picture/msg-subscription.jpeg)


例如：
```
Outputs:
LarkbotAppStack.msgEventRoleArn = arn:aws:iam::123456789012:role/larkbot-larkbotmsgeventServiceRoleC3080B6B-V1ESZLK7ODYY
LarkbotAppStack.msgEventapiEndpointAC31EC6D = https://t68l424zt0.execute-api.ap-northeast-1.amazonaws.com/prod/

```

添加/messages路径：

```
https://t68l424zt0.execute-api.ap-northeast-1.amazonaws.com/prod/messages
```

8. 权限管理

飞书管理员需要授权飞书机器人执行一些必要的单聊或群聊消息接收，发送，读取基本通讯录信息及附件处理功能的权限。在机器人配置的主页中，点击页面最左侧的开发配置段落中的权限管理功能。

**消息与群组权限**
在屏幕中央的权限配置的页面中，选择消息与群组，点击选择下面截屏中全部的权限。

![消息与群组权限1](picture/feishu-permission-1.jpeg)

进入权限列表的第二页,点击选择下面四个权限。

![消息与群组权限2](picture/feishu-permission-2.jpeg)

**通讯录权限**
选择权限配置中的通讯录权限，添加下列两个权限。

![通讯录权限1](picture/feishu-permission-2.5.jpeg)


**批量开通**

全部选择完毕后，点击页面右上角的批量开通按钮。在批量开通提示菜单中，可以选择确定回到机器人配置主页继续调整权限，或者点击确认并前往创建应用版本发布机器人。

![批量开通权限](picture/feishu-permission-3.jpeg)


**版本发布**

权限选择完毕后，需要发布机器人。进入机器人配置主页左侧的应用发布段落的版本管理与发布功能，点击右上角的创建新版本

![创建新版本](picture/release-management-1.jpeg)


输入相应的版本信息，然后点击发布。在机器人的可用范围选项中，选择添加所有可能用到机器人的飞书账号信息。但是如果飞书机器人被加入到群聊中，群聊中的用户即使没有被加入到机器人的可用范围列表中，仍然可以通过群聊添加机器人并且和机器人对话。

要限制机器人的使用，需要使用白名单功能。

![更新版本信息](picture/release-management-2.jpeg)



9.  事件订阅

在飞书机器人配置主页，在页面左侧开发配置段落中，选择事件订阅功能，然后选择添加事件按钮。

![事件订阅](picture/event-subscription-permission.jpeg)

在列表中选择消息与群组中的接受消息(v2.0),然后点击确认添加。

![事件订阅](picture/event-subscription-permission-1.png)

正确添加订阅消息权限后，页面应展示下面的内容。

![事件订阅列表](picture/event-subscription-permission-2.png)


10. 设置AppID和AppSecret




#### 在DynamoDB中配置机器人的自定义参数
---


#### 创建SupportAPI角色
---


## 成本预估
---


## TODO列表
---
[TODO List](TODO.md)