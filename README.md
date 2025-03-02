# game-fun-be

## 数据结构

https://xsglh58murtx.sg.larksuite.com/wiki/DitdwzMCfiQrlXklQFVlr1Nsgdd

## 接口

https://xsglh58murtx.sg.larksuite.com/wiki/Nq3kw3YYDinDjnkW27nlHPNOgCe

## 服务器

位置：日本
操作系统：ubuntu 24

### 服务器（同 VPC 部署）
- 服务器 1（8 核 16g）：运行 web 后台服务，node 前端服务
- 服务器 2（8 核 16g）：运行链端采集服务和链上数据查询服务
- 服务器 3（8 核 16g）：运行后端消费和异步处理定时任务服务

### 其他服务
- es = 7.10 版本（指定该版本），3 master 4 data 节点，16g 内存
- kafka 3.0 以上版本，3 节点集群
- mysql 8.0 以上版本，一主一从
- redis 6.0 以上版本，一主一从，4g 内存
- clickhouse 23.8 以上版本，8 核 32g 双副本
