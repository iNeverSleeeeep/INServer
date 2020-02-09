# INServer

#### 介绍
Golang Game Server

### 初始化
protobuf环境初始化
```
go get github.com/gogo/protobuf/protoc-gen-gofast
go get github.com/gogo/protobuf/proto
go get github.com/gogo/protobuf/jsonpb
go get github.com/gogo/protobuf/protoc-gen-gogo
go get github.com/gogo/protobuf/gogoproto
go get github.com/mitchellh/protoc-gen-go-json
go get github.com/divan/expvarmon
```

protoc 3.11.2

[一些原始想法](https://github.com/iNeverSleeeeep/INServer/blob/master/MIND.md)

目前版本存在几个核心问题没有解决，这些问题直接影响了服务器的可用性和易用性。
1. 服务器如何支持乱序启动，如何支持服务器宕机重开后恢复各个服务器的状态，服务器的关机流程。
2. **服务器中有几种大的类型的消息，他们的流向应该是如何，消息如何从一个起点到达终点。这个问题不解决，目前实现玩家移动的时候，整体的逻辑是混乱的，就算做出来了，以后也会越来越乱。**
3. 每个服务器应该保存玩家/角色的哪些数据，数据发生变化时的同步如何进行。

根据以上的问题，接下来顺序完成的功能
1. 完成整个服务器群的生命周期整理，清晰的开服关服流程，服务器重启流程。
2. 增加月台服务器 承接玩家进入游戏但是没有分到逻辑服务器的状态。
3. 整理清楚玩家和角色的概念，登录流程走完之后，一切都是角色，不应该有玩家了。完成一个玩家的生命周期整理，进入，退出游戏流程与关服时保存流程。
4. 玩家数据分布整理，玩家的哪些数据需要存储在哪些服务器。
5. 消息类别整理，不同类别的消息的在服务器之间如何流动的流程整理。
6. 完成玩家的移动和切地图逻辑。
