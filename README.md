# SonarQube 钉钉机器人插件

一款简单实用的SonarQube钉钉消息推送插件

## 💡 使用

### 安装

#### docker方式(推荐)

```dockerfile
docker run \
-d \
--name=sonar-dingtalk-plugin \
--restart=always \
-p 9010:9010 \
xbmlz/sonar-dingtalk-plugin

# 使用代理
docker run \
-d \
--name=sonar-dingtalk-plugin \
--restart=always \
-p 9010:9010 \
-e HTTPS_PROXY=http://username:password@ip:port \
xbmlz/sonar-dingtalk-plugin
```

#### 二进制安装

根据对应操作系统下载 [Release](https://github.com/xbmlz/sonar-dingtalk-plugin/releases)文件，解压启动即可。

## 设置

### 添加钉钉机器人

1. 打开钉钉[群设置]——[智能群助手]——[添加机器人]——[自定义]

2. 机器人名字填写SonarQube，头像可根据需求自行更换

3. 安装设置选择[自定关键词]，添加`Bugs`、`漏洞`

![image-20220618135039580](https://cdn.jsdelivr.net/gh/xbmlz/static@main/img/202206181350643.png)

4. 点击完成后将Webhook地址中`access_token=`后面的内容复制保存下来（类似https://oapi.dingtalk.com/robot/send?access_token=xxxxxxx）

### SonarQube设置

1. 进入项目，点击[项目配置]——[网络调用]——[创建]

   ![image-20220618135444058](https://cdn.jsdelivr.net/gh/xbmlz/static@main/img/202206181354086.png)

2. 名称填写任意值，如`DingTalk`，URL填写 

   ```bash
   http://[插件安装电脑ip]:9010/dingtalk?access_token=[配置钉钉机器人时保存的access_token]&sonar_token=[sonar的token]
   ```

   ![](https://cdn.jsdelivr.net/gh/xbmlz/static@main/img/202206181500350.png)

   注意：sonar没有开启权限验证时不需要填写

   附：sonarqube 认证token获取方式，点击[配置]——[权限]——[用户]——[令牌]

### CI/CD

此步骤可自行百度配置，支持Gitlab、Jenkins等

参考文档：[Overview | SonarQube Docs](https://docs.sonarqube.org/8.3/analysis/branch-pr-analysis-overview/)

### 消息推送

完成上述步骤，就可以将sonarqube扫描结果，推送到钉钉群了

![](https://cdn.jsdelivr.net/gh/xbmlz/static@main/img/202206181406084.png)
