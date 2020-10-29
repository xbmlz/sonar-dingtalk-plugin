# sonar-dingtalk-plugin

sonarqube 钉钉机器人插件

### 使用方法

#### 1. 部署插件

下载对应操作系统的插件，运行即可

#### 2. 添加钉钉群机器人

在钉钉群设置->智能群助手->添加自定义机器人

![图片描述](https://i.niupic.com/images/2020/10/29/8VEc.jpg)

复制`webhook`地址中的access_token=后面的内容，后面会用到

安全设置选择`自定义关键词`，添加`BUG`和`漏洞`

![图片描述](https://i.niupic.com/images/2020/10/29/8VFN.jpg)

#### 3. 设置SonarQube

点击`项目配置`下的`网络调用`

![图片描述](https://i.niupic.com/images/2020/10/29/8VFP.jpg)

点击右上角`创建`按钮

![图片描述](https://i.niupic.com/images/2020/10/29/8VFR.jpg)

名称随便填，URL填 `http://插件部署电脑的IP:9001/dingtalk?access_token=这里填刚才复制的机器人的token`