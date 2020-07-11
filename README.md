# CalculatePool

一个基于RESTful实现的、简单的PoW算力池模型。  

## 文件结构  

- calculatepool：  
  - calculatepool.proto：定义用到的结构
  - server.go：算力池模型的具体实现
- center.go：算力池的控制中心
- computenode.go：算力池的计算节点
- client.go：请求发起方  

## server  

server.go中提供了如下的方法：  

- NewServer()：创建一个Server的实例，用于之后的操作
- Start()：启动RESTful服务
- RegisterRoutes()：注册RESTful提供的接口
- SetHard()：设置计算难度
- GetHard()：获取当前的计算难度
- Hello()：用于判断节点是否存活的接口
- Register()：计算节点向中心节点注册
- PoW()：Proof of Work的简单实现
- DoWork()：中心节点接收PoW请求的接口

## center  

算力池的控制中心，调用server.go中的方法，对外提供PoW计算接口，以及算力池内部计算节点的注册。

## computenode  

向控制中心注册本节点信息，并提供PoW计算服务。


本作品采用[署名-非商业性使用-相同方式共享 4.0 国际 (CC BY-NC-SA 4.0)](https://creativecommons.org/licenses/by-nc-sa/4.0/deed.zh)进行许可，使用时请注明出处。