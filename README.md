# game-fun-be

## 数据结构

https://xsglh58murtx.sg.larksuite.com/wiki/DitdwzMCfiQrlXklQFVlr1Nsgdd

## 接口

https://xsglh58murtx.sg.larksuite.com/wiki/Nq3kw3YYDinDjnkW27nlHPNOgCe

## 服务器

https://hxny4q0lcre.feishu.cn/docx/E5s2dv7Ibo3a4Zx99XMcd1E2nZc

位置：日本
提供商：阿里云
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

---

## ssh config

```sh

# ~/.ssh/config

Host sg-jump
  HostName 8.129.133.197
  User ecs-user
  IdentityFile ~/.ssh/sg.pem
  StrictHostKeyChecking no

Host sg-be-api-node1
  HostName 172.16.21.135
  User ecs-user
  ProxyJump sg-jump
  StrictHostKeyChecking no
  IdentityFile ~/.ssh/sg.pem
  IdentitiesOnly yes
  # PermitLocalCommand yes
  # RequestTTY force
  # ForwardAgent yes

Host sg-be-consumer-node3
    HostName 172.16.21.1
    User ecs-user
    ProxyJump sg-jump
    StrictHostKeyChecking no
    IdentityFile ~/.ssh/sg.pem
    IdentitiesOnly yes

Host sg-be-block-node2
    HostName 172.16.21.134
    User ecs-user
    ProxyJump sg-jump
    StrictHostKeyChecking no
    IdentityFile ~/.ssh/sg.pem
    IdentitiesOnly yes
```
