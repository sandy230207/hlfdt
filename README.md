# hlfdt

## 如何建置可執行檔
```shell=
go build src/main.go
mv main hlfdt
```

## 安裝程式庫
```shell=
go get -v ./...
```

## 指令
![](https://github.com/sandy230207/fabric-tool/blob/master/doc/commands.png)
### 檢查並產生所需配置檔與部署腳本
生成包含憑證認證機構及節點的docker-compose文件、組織與策略配置檔(`configtx.yaml`)與若干部署所需之腳本
```shell
generate --file <DIR>
```
> `<DIR>` 為配置檔路徑名稱
> 若無指定之配置檔路徑，則預設之路徑為 `./config.yaml`

以可執行檔運行
```shell
./hlfdt generate --file config.yaml
```
或以GO編譯運行(建議開發時使用)
```shell
go run src/main.go generate --file config.yaml
```

### 創建與啟動 HLF 網路
```shell
up
```

以可執行檔運行
```shell
./hlfdt up
```
或以GO編譯運行(建議開發時使用)
```shell
go run src/main.go up
```

### 停止與清除 HLF 網路
```shell
down
```

以可執行檔運行
```shell
./hlfdt down
```
或以GO編譯運行(建議開發時使用)
```shell
go run src/main.go down
```

### 啟動可監控 HLF 網路的 UI
查看 HLF 網路運行之狀況及區塊數量
```shell
ui --port <PORT_NUMBER>
```
>  `<PORT_NUMBER>` 為欲運行 UI 的 Port Number

以可執行檔運行
```shell
./hlfdt ui --port 8000
```
或以GO編譯運行(建議開發時使用)
```shell
go run src/main.go ui --port 8000
```

## hlfdt 配置檔
使用者需根據文件要求之格式撰寫如 [config.yaml](https://github.com/sandy230207/fabric-tool/blob/master/config.yaml) 之YAML配置文件。
配置文件共分成四大區段，每區段皆為陣列型態
- [Channels（通道）](#Channels)
- [Chaincodes（鏈碼）](#Chaincodes)
- [Organizations（組織）](#Organizations)
- [CertificateAuthorities（CA）](#CertificateAuthorities)

### Channels
該區段主為宣告通道名稱(Name)及定義該通道各項政策(Policies)，若組織要加入該通道，則須先在此區段定義通道名稱，並於[Organizaiotns](#Organizations)區段之Channels中添加該通道之名稱。
### Chaincodes
主為宣告鏈碼名稱(Name)、設定鏈碼開發語言(Language)、設定鏈碼路徑(Path)及選定欲使用該鏈碼的通道名稱之區段(Channels)。
### Organizations
為宣告組織名稱及定義組織型態之區段，共分為兩種組織型態：
- #### PeerOrg
    節點組織，其下可宣告[背書節點(Endorsing Peer)](#EndorsingPeers)及[提交節點(Committing Peer)](CommittingPeers)，並設定該組織的讀政策、寫政策、管理者政策及背書政策，以及設定所連接之通道名稱(Channels)。
- #### OrderOrg
    排序組織，其下可宣告[排序節點(Ordering Peer)](#Peers)，並設定該組織的讀政策、寫政策及管理者政策。
    
組織中節點又依節點型態不同，分為三類型的鍵：
- #### EndorsingPeers
    宣告組織內的背書節點(Endorsing Peer)名稱及設定節點位址(Address)、節點埠號(Port)及節點資料庫埠號(DBPort)之區段。僅PeerOrg可使用此鍵。
- #### CommittingPeers
    宣告組織內的提交節點(Committing Peer)名稱及設定節點位址(Address)、節點埠號(Port)及節點資料庫埠號(DBPort)之區段。僅PeerOrg可使用此鍵。
- #### Peers
    宣告組織內的排序節點(Ordering Peer)名稱及設定節點位址(Address)及節點埠號(Port)之區段。僅OrderOrg可使用此鍵。
### CertificateAuthorities
宣告CA的名稱(Name)及設定CA位址(Address)、CA埠號(Port)，以及設定所連接之組織名稱(Organizations)。


## 自動生成之配置檔與腳本介紹
| 檔案路徑與名稱 | 功用說明 | 
| -------- | -------- |
| network.sh | 1. 入口點<br>2.	設定參數預設值 |
| scrpits/configUpdate.sh | 1.	取得通道資訊<br>2.	更新通道資訊 | 
| scripts/createChannel-{通道名稱}.sh | 1.	建立通道<br>2.	將節點加入通道 | 
| scripts/deployCC.sh | 1.	打包鏈碼<br>2.	安裝鏈碼<br>3.	以組織成員身份同意鏈碼的安裝<br>4.	提交鏈碼的更動至帳本<br>5.	 確認安裝鏈碼的更動是否已提交至帳本<br>6.	初始化鏈碼<br>7.	查詢鏈碼 | 
| scripts/envVar.sh | 1.	設定各項環境變數 | 
| scripts/setAnchorPeer.sh | 1.	建立錨節點的更新<br>2.	提交錨節點的更新至帳本 | 
| scripts/utils.sh | 1.	network.sh中各項flag的功能提示 |
| organizations/fabric-ca/registerEnroll.sh | 1.	向CA註冊身份 | 
| organizations/ccp-generate.sh | 1.	產生用於GRPC的連接配置文件 | 
| organizations/ccp-template.json | 1.	用於GRPC的連接配置文件JSON模板 | 
| organizations/ccp-template.yaml | 1.	用於GRPC的連接配置文件YAML模板 | 
| configtx/configtx-{通道名稱}.yaml | 1.	Hyperledger Fabric主要配置文件<br>2.	定義通道、各組織的政策 |
| docker/docker-compose-ca.yaml	| 1.	CA的docker container配置文件 | 
| docker/docker-compose-test-net.yaml | 1.	所有節點與cli的docker container配置文件 | 

## 檔案目錄
```
├── README.md
├── bin
│   ├── configtxgen
│   ├── configtxlator
│   ├── cryptogen
│   ├── discover
│   ├── fabric-ca-client
│   ├── fabric-ca-server
│   ├── idemixgen
│   ├── orderer
│   └── peer
├── config
│   ├── configtx.yaml
│   ├── core.yaml
│   └── orderer.yaml
├── config.yaml
├── configtx
│   └── configtx.yaml
├── configtx-mychannel
│   └── configtx.yaml
├── configtx-secondchannel
│   └── configtx.yaml
├── doc
│   └── commands.png
├── docker
│   ├── docker-compose-ca.yaml
│   ├── docker-compose-couch.yaml
│   └── docker-compose-test-net.yaml
├── go.mod
├── go.sum
├── hlfdt
├── mychannel.block
├── network.sh
├── organizations
│   ├── ccp-generate-original.sh
│   ├── ccp-generate-script.sh
│   ├── ccp-generate.sh
│   ├── ccp-template.json
│   ├── ccp-template.yaml
│   ├── cryptogen
│   │   ├── crypto-config-orderer.yaml
│   │   ├── crypto-config-org1.yaml
│   │   └── crypto-config-org2.yaml
│   ├── fabric-ca
│   │   ├── orderOrg
│   │   ├── org1
│   │   ├── org2
│   │   └── registerEnroll.sh
│   └── template
│       └── registerEnroll.sh
├── scripts
│   ├── configUpdate.sh
│   ├── createChannel-mychannel.sh
│   ├── createChannel-secondchannel.sh
│   ├── createChannel.sh
│   ├── deployCC.sh
│   ├── envVar.sh
│   ├── network-original.sh
│   ├── setAnchorPeer.sh
│   └── utils.sh
├── src
│   ├── config
│   │   ├── config.go
│   │   └── model.go
│   ├── configtx
│   │   ├── generater.go
│   │   └── model.go
│   ├── docker-ca
│   │   ├── generater.go
│   │   └── model.go
│   ├── docker-couch
│   │   ├── generater.go
│   │   └── model.go
│   ├── docker-net
│   │   ├── generater.go
│   │   └── model.go
│   ├── fabric-ca-server-config
│   │   ├── generater.go
│   │   └── model.go
│   ├── fabric-network
│   │   ├── ccp-generate.go
│   │   ├── configUpdate.go
│   │   ├── createChannel.go
│   │   ├── deployCC.go
│   │   ├── enrollRegister.go
│   │   ├── envVar.go
│   │   ├── network.go
│   │   └── setAnchorPeer.go
│   ├── main.go
│   ├── monitor
│   │   ├── index.html
│   │   ├── index.html.orig
│   │   ├── model.go
│   │   ├── model.go.orig
│   │   ├── monitor.go
│   │   └── monitor.go.orig
│   └── utils
│       └── utils.go
└── testfile
    ├── add-org4.yaml
    ├── asset-transfer-basic
    │   ├── application-go
    │   ├── application-java
    │   ├── application-javascript
    │   ├── application-typescript
    │   ├── chaincode-external
    │   ├── chaincode-go
    │   ├── chaincode-java
    │   ├── chaincode-javascript
    │   └── chaincode-typescript
    ├── chaincode.yaml
    └── demo-scope-original.yaml
```

