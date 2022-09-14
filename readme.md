# photoview 帮助文档

<p align="center" class="flex justify-center">
    <a href="https://www.serverless-devs.com" class="ml-1">
    <img src="http://editor.devsapp.cn/icon?package=photoview&type=packageType">
  </a>
  <a href="http://www.devsapp.cn/details.html?name=photoview" class="ml-1">
    <img src="http://editor.devsapp.cn/icon?package=photoview&type=packageVersion">
  </a>
  <a href="http://www.devsapp.cn/details.html?name=photoview" class="ml-1">
    <img src="http://editor.devsapp.cn/icon?package=photoview&type=packageDownload">
  </a>
</p>

<description>

> ***相册***

</description>

<table>

## 前期准备
使用该项目，推荐您拥有以下的产品权限 / 策略：

| 服务/业务 | 函数计算 |     
| --- |  --- |   
| 权限/策略 | AliyunFCFullAccess |     


</table>

<codepre id="codepre">



</codepre>

<deploy>

## 部署 & 体验

<appcenter>

- :fire: 通过 [Serverless 应用中心](https://fcnext.console.aliyun.com/applications/create?template=photoview) ，
[![Deploy with Severless Devs](https://img.alicdn.com/imgextra/i1/O1CN01w5RFbX1v45s8TIXPz_!!6000000006118-55-tps-95-28.svg)](https://fcnext.console.aliyun.com/applications/create?template=photoview)  该应用。 

</appcenter>

- 通过 [Serverless Devs Cli](https://www.serverless-devs.com/serverless-devs/install) 进行部署：
    - [安装 Serverless Devs Cli 开发者工具](https://www.serverless-devs.com/serverless-devs/install) ，并进行[授权信息配置](https://www.serverless-devs.com/fc/config) ；
    - 初始化项目：`s init photoview -d photoview`   
    - 进入项目，并进行项目部署：`cd photoview && s deploy -y`

</deploy>

<appdetail id="flushContent">

# 应用详情

本项目是一个基于函数计算和Serverless Mysql的相册体验样例,相较于传统架构用户可以体验到:
* FaaS+BaaS全链路Serverless服务，无需采购和管理服务器等基础设施，只需专注业务逻辑的开发
* 采用DataAPI的方式连接操作数据库，用户不再关心数据库的连接
* 采用Serverless Mysql数据库，根据业务负载变化自动弹性分配资源
* 实例自动启停，客户中断数据库访问时，进入静默状态，节省资源


# 执行过程
### a. 填写应用参数
| 参数                 | 类型   | 默认值                          | 名称       | 备注                                                                                                                                                   |
| ------------------- | ----- | ------------------------------ |----------|------------------------------------------------------------------------------------------------------------------------------------------------------|
| region              | string | cn-hangzhou                    | 地域       | 创建应用所在的地区，资源也会创在该地域，目前只支持 cn-hangzhou                                                                                                                                  |
| roleArn             | string | 无默认，必填                    | RAM角色ARN | 应用所属的函数计算服务配置的 role, 请提前创建好对应的 role, 授信函数计算服务, 并配置好 AliyunOSSFullAccess, AliyunFCDefaultRolePolicy, AliyunRDSFullAccess policy 和 AliyunECSFullAccess |
| ossBucket           | string | 无默认，非必填                           | OSS存储桶名  | OSS存储桶名(注意和函数同地域)，不填写即不使用 OSS 存储资源部署状态                                                                                                                                  |
| ossObjectName           | string | 无默认，非必填    | OSS 对象名  | OSS 对象名，不填写即不使用 OSS 存储资源部署状态                                                                                                                                            |
                                                                                                                       


### b. 部署应用
1. 应用中心会自动开始执行应用部署的的4个步骤：(前置环境，资源同步，资源检查，执行部署)
2. 其中在执行部署时会先后部署两个函数：资源创建函数，资源消费函数（PhotoView函数）
3. 在部署 PhotoView 函数时会先调用资源创建函数，创建出 Serverless MySQL 资源。\
`注意： 资源创建函数执行时会真实创建出 Serverless MySQL 数据库实例，这需要大概等待5到10分钟` \
具体表现为应用部署会一直在`Invoke resource creator function`：
![](https://img.alicdn.com/imgextra/i2/O1CN01d38Cwn1KwT9IOBZne_!!6000000001228-2-tps-1688-166.png)  
4. 资源部署完成后，点击访问域名，即可体验 PhotoView 应用
![](https://img.alicdn.com/imgextra/i3/O1CN01g83px31O2pUnuwWMS_!!6000000001648-2-tps-3948-1310.png)
### c. 资源创建汇总
将会创建：
| 资源                   | 备注                                                                                                                                                   |
| ------------------- |------------------------------------------------------------------------------------------------------------------------------------------------------|
| Serverless MySQL 实例| 包含 RDS 实例, 数据库，账户，dataAPI secret|
| VPC, vswitch，安全组| 支持函数计算和 Serverless MySQL 通信|
| 两个函数: 资源创建函数，PhotoView 函数| FC 提供资源创建和资源消费能力|


# 应用原理

## 一键部署 Serverless Mysql 资源及 PhotoView 函数
![](https://img.alicdn.com/imgextra/i2/O1CN01zJtHRa1HKxrrXeDCq_!!6000000000740-2-tps-2724-1004.png)  

函数 1：资源创建函数：
- 资源创建函数原理是将 Terraform 集成在 FC Custom Container 函数中，从而通过管理 FC 函数来管理 Terraform job。  
- 我们内置了 Serverless MySQL 的 TF 资源文件，用户只需要在创建应用是填入少量 Mysql 数据库配置即可创建出Serverless MySQL 资源。  
- 当填写有效 ossBucket 和 ossObjectName 时，用户将使用自己的 OSS 存储资源部署状态。

函数2：PhotoView函数
-   创建资源：
    1. Serverless Devs 工具具备 pre-action 能力，即在部署函数前完成某项工作。

    2. 利用 pre-action 能力在部署 PhotoView 函数前调用资源创建函数，创建出资源，并将资源配置传入到 PhotoView 函数的环境变量里。
- PhotoView 函数消费 Serverless MySQL 资源
  1. PhotoView 函数通过从环境变量中获取 Serverless MySQL 的资源配置，之后通过访问 Serverless MySQL 数据库从而完成业务需求。





















</appdetail>

<devgroup>

## 开发者社区

您如果有关于错误的反馈或者未来的期待，您可以在 [Serverless Devs repo Issues](https://github.com/serverless-devs/serverless-devs/issues) 中进行反馈和交流。如果您想要加入我们的讨论组或者了解 FC 组件的最新动态，您可以通过以下渠道进行：

<p align="center">

| <img src="https://serverless-article-picture.oss-cn-hangzhou.aliyuncs.com/1635407298906_20211028074819117230.png" width="130px" > | <img src="https://serverless-article-picture.oss-cn-hangzhou.aliyuncs.com/1635407044136_20211028074404326599.png" width="130px" > | <img src="https://serverless-article-picture.oss-cn-hangzhou.aliyuncs.com/1635407252200_20211028074732517533.png" width="130px" > |
|--- | --- | --- |
| <center>微信公众号：`serverless`</center> | <center>微信小助手：`xiaojiangwh`</center> | <center>钉钉交流群：`33947367`</center> | 

</p>

</devgroup>