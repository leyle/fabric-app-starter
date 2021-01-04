# fabric app starter service

这个主要是给 application 层一个对接 chaincode 的通用抽象层。本 service 提供的服务包含三部份内容：

- chaincode 代码（包含 public 及 private）
- 用户管理及 ca 证书管理
- chaincode 的通用 CRUD 操作抽象

下面的 `{{host}}`值是 `http://127.0.0.1:9000`，即本 service 启动时的监听地址。

用户 login 后，获取到一个 token 值，在其他需要验证的接口中，在 header 中配置 `X-TOKEN`为 key， 此 token 值为value 的一个属性。

下面分别介绍对接使用方法。

---

## chaincode 部署

chaincode 包含 public 和 private 内容，用户根据不同的需求，可以混合使用 public 和 private 或者仅使用 private 代码。

因为 fabric 网络使用 couchdb 来存储 ledger 的 world state 值。所以，为了支持 couchdb 的 find 等搜索操作，需要用户提供一个预先定义的 index 定义 json 文件，在部署 chaincode 时，一同打包到对应的 package 中。

此索引文件存放在 `public/META-INF/statedb/couchdb/indexes/xxxx.json` 中，文件中的内容，需要根据实际业务，添加相关的索引数据。

具体的 chaincode 部署，见 fabric-network-starter 中关于 chaincode 安装部署的说明。

---

## 程序启动

```shell
# 程序先要编译
cd cmd/apiserver/
go build

# 启动程序时，需要依赖一个配置文件
# 配置文件的模板放在 examples/example_config.yaml 文件中
# 这个文件中关于 fabric 的部分，依赖于 connection.yaml 文件及 wallet 等数据的存放路径
# 需要提前配置好
# 假设把 config.yaml 已经修改并存放在当前 cmd/apiserver/config.yaml 位置
./apiserver -c config.yaml
# 即可启动 service api
```

---

## 用户管理

使用 chaincode 相关 api 时，需要用户有相关凭据，这个凭据包含两部分

1. jwt user
2. ca enrolled

具体的来说，初始化时，配置了一个org 的 admin 用户，此用户可以对本 org 的 ca 进行用户的 register/enroll 操作。

所以，第一步就是 login 此用户

这部分的代码在 https://gitlab.com/emaliio/fabric-user-manager

需要添加新的接口及功能时，在相关地方添加即可。

### 用户 login 操作

```shell
POST {{host}}/api/jwt/user/login
# 输入的 json body
{
    "username": "orgadmin",
    "password": "passwd"
}

# 如果账户密码正确，返回的数据例子如下所示
# 其中 data.token 是我们后续需要的内容。
{
    "code": 200,
    "msg": "OK",
    "data": {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiJvcmdhZG1pbiIsInVzZXJuYW1lIjoib3JnYWRtaW4iLCJyb2xlIjoiYWRtaW4iLCJleHAiOjE2MTA1ODgxOTMsImlhdCI6MTYwOTcyNDE5M30.QHxRtl1SgWxbvTQCNV1NzCEbvglYlFn5OdB3eR3HpMU",
        "user": {
            "id": "orgadmin",
            "_rev": "1-8df046f034238bf71b853003d9fc9764",
            "username": "orgadmin",
            "role": "admin",
            "valid": true,
            "created": {
                "second": 1609317595,
                "humanTime": "2020-12-30 16:39:55"
            },
            "updated": {
                "second": 1609317595,
                "humanTime": "2020-12-30 16:39:55"
            }
        }
    }
}
```

---

### admin 分配新账户

```shell
POST {{host}}/api/jwt/user/create
# 输入的数据
{
    "username": "devtest2",
    "password": "passwd",
    "role": "client"
}

# 其中 role 可选项是 client/admin/peer/orderer
# 但是一般来说，我们分配的是 application 操作者，所以都是 client 值 
# 执行成功时的输出
{
    "code": 200,
    "msg": "OK",
    "data": {
        "id": "5ff271b8de18f9412cc3fbe1",
        "username": "devtest2",
        "salt": "20210104093904",
        "role": "client",
        "valid": true,
        "created": {
            "second": 1609724344,
            "humanTime": "2021-01-04 09:39:04"
        },
        "updated": {
            "second": 1609724344,
            "humanTime": "2021-01-04 09:39:04"
        }
    }
}

# 可能的错误是，用户名已经存在了
# 比如
{
    "code": 400,
    "msg": "username[devtest2] exist",
    "data": ""
}
```

---

### 封禁用户

```shell
# admin 可以停用一个账户
# 但是需要注意的是，这里的停用仅为停用此账户的 user/passwd 账户，不包含对 ca 的 revoke 操作。
# todo
```

---

### 取消封禁

```shell
# todo
```



---



## chaincode 操作 api

这部分的代码在 https://gitlab.com/emaliio/fabric-app-starter

需要添加更多的功能时，可以在这个基础上添加。

### 创建一个既有 public data 又有 private data 的数据

因为是一个通用的 api，所以主要做的事情是把调用者提交的数据存储到对应的 channel 及 chaincode 上。用户的业务数据需要先 json 序列化为字符串提交。同时需要根据情况把，把 public 和 private data 分开在不通 field 提交上来。

public 及 private data 的提交，不是全部存在，可以只存在一个，但是不能都不存在。

```shell
POST {{host}}/api/chaincode/publicandprivate/create
# app - 指的是调用者给自己的 app 取的名字
# dataId - 这个是提交的此 data 数据的唯一 id，也是后续查询等的凭证，需要调用者保证唯一性
# chaincode 会校验同一个 app 下的 dataId 的唯一性
# public - 的结构包含三部份内容，分别是  channel chaincode 及具体的业务data，业务 data 全部存放在 dataJson 中
# 下面分别写三个例子，分别是仅有 public data，仅有private data，public 和 private data 都存在
# 需要注意的是，public 与 private 的 channel 可以不是同一个，他们两者之间的 channel 及 chaincode 没有任何联系。

# 1. 仅有 public data
{
    "app": "customerConsent",
    "dataId": "5dc69a06-8dca-4e96-891b-d6e59cc82639",
    "public": {
        "channel": "one",
        "chaincode": "public",
        "dataJson": "{\"dataId\":\"5dc69a06-8dca-4e96-891b-d6e59cc82677\",\"customerId\":\"cid1\",\"providerId\":\"pid1\",\"array\":[\"a\",\"b\",\"c\"],\"dict\":{\"key1\":\"val1\",\"key2\":\"val2\"}}"
    }
}


# 2. 仅有 private data
# 与 public 结构基本一致，只是 private 结构中多了一个 collectionName 数据，指定使用 private policy
# 这个 collectionName 值是从 collections.json 文件中来的，这个是用户在部署 chaincode 时确定的。
{
    "app": "customerConsent",
    "dataId": "5dc69a06-8dca-4e96-891b-d6e59cc82639",
    "private": {
        "channel": "one",
        "chaincode": "private",
        "collectionName": "org1private",
        "dataJson": "{\"name\":\"john\",\"age\":22,\"class\":\"mca\"}"
    }
}

# 3. public 及 private 同时存在
{
    "app": "customerConsent",
    "dataId": "5dc69a06-8dca-4e96-891b-d6e59cc82639",
    "public": {
        "channel": "one",
        "chaincode": "public",
        "dataJson": "{\"dataId\":\"5dc69a06-8dca-4e96-891b-d6e59cc82677\",\"customerId\":\"cid1\",\"providerId\":\"pid1\",\"array\":[\"a\",\"b\",\"c\"],\"dict\":{\"key1\":\"val1\",\"key2\":\"val2\"}}"
    },
    "private": {
        "channel": "another",
        "chaincode": "private",
        "collectionName": "org1private",
        "dataJson": "{\"name\":\"john\",\"age\":22,\"class\":\"mca\"}"
    }
}
```

---

### 根据 id 读取提交的数据

```shell
GET {{host}}/api/chaincode/publicandprivate/public/info
# 因为查询的复杂性，所以把整个查询条件放在了 request body 中
# 下面是一个完整的同时查询同一个 dataId 的 public 及 private 数据的查询条件
# 也可以根据实际情况，仅有 public 或 private 的查询条件
{
    "app": "customerConsent",
    "dataId": "5dc69a06-8dca-4e96-891b-d6e59cc82659",
    "public": {
        "channel": "one",
        "chaincode": "p2"
    },
    "private": {
        "collectionName": "org1private",
        "channel": "ppp",
        "chaincode": "ppp"
    }
}

# 根据查询的情况，返回数据结构大概如下所示
# 其中 data.success 字段有三个可选值 none/partial/all
# none - 指的是查询没有数据/失败
# partial - 在同时查询 public 和 private 时，指仅 public 或 private 返回成功的数据
# all - 指根据查询条件，返回成功的数据
{
    "code": 200,
    "msg": "OK",
    "data": {
        "success": "all",
        "errMsg": "",
        "public": {
            "id": "customerConsent|5dc69a06-8dca-4e96-891b-d6e59cc82659",
            "app": "customerConsent",
            "dataId": "5dc69a06-8dca-4e96-891b-d6e59cc82659",
            "data": "{\"array\":[\"a\",\"b\",\"c\"],\"customerId\":\"cid1\",\"dataId\":\"5dc69a06-8dca-4e96-891b-d6e59cc82677\",\"dict\":{\"key1\":\"val1\",\"key2\":\"val2\"},\"providerId\":\"pid1\"}",
            "createdAt": 1609695763
        },
        "private": {
            "id": "customerConsent|5dc69a06-8dca-4e96-891b-d6e59cc82659",
            "collectionName": "org1private",
            "app": "customerConsent",
            "dataId": "5dc69a06-8dca-4e96-891b-d6e59cc82659",
            "data": "\"{\\\"name\\\":\\\"john\\\",\\\"age\\\":22,\\\"class\\\":\\\"mca\\\"}\"",
            "createdAt": 1609695761
        },
        "app": "customerConsent",
        "dataId": "5dc69a06-8dca-4e96-891b-d6e59cc82659"
    }
}
```

---

### 按 id 删除数据

```shell
# TODO
```

---

### 搜索

```shell
# TODO
```

