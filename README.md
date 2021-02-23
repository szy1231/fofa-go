# FOFA Pro - Go 使用说明文档

## FOFA Pro API   

<a href="https://fofa.so/api"><font face="menlo">`FOFA Pro API`</font></a> 是资产搜索引擎 <a href="https://fofa.so/">`FOFA Pro`</a> 为开发者提供的 `RESTful API` 接口, 允许开发者在自己的项目中集成 `FOFA Pro` 的功能。    


## FOFA Go SDK

基于 `FOFA Pro API` 编写的 `golang` 版 `SDK`, 方便 `golang` 开发者快速将 `FOFA Pro` 集成到自己的项目中。


### 环境

### 开发环境

``` zsh
$ go version
go version go1.15.3 windows/amd64
```

### 测试环境

``` zsh
$ go version
go version go1.15.3 windows/amd64
```

### 使用

```go
import "github.com/xiaoyu-0814/fofa-go/fofa"
```

### 获取

``` zsh
go get github.com/xiaoyu-0814/fofa-go/fofa
```

### 依赖

### Email & API Key   

| 字段  | 描述                                                         |
| ----- | ------------------------------------------------------------ |
| Email | 用户登陆 `FOFA Pro` 使用的 `Email`                           |
| Key   | 前往 [**`个人中心`**](https://fofa.so/my/users/info) 查看 `API Key` |

如果开发者经常使用固定的账号，建议将`email`与`key`添加到环境变量中。

`SDK` 提供的示例代码就是使用的这种形式。


### Example   

``` go
func FofaExample() {
	email := os.Getenv("FOFA_EMAIL")
	key := os.Getenv("FOFA_KEY")

	clt := fofa.NewFofaClient([]byte(email), []byte(key))
	if clt == nil {
		fmt.Printf("create fofa client\n")
		return
	}
    
	//QueryAsJSON
	ret, err := clt.QueryAsJSON(1, []byte(`body="小米"`))
	if err != nil {
		fmt.Printf("%v\n", err.Error())
		return
	}
	fmt.Printf("%s\n", ret)
    
	//QueryAsObject
	data, err := clt.QueryAsObject(1, []byte(`domain="163.com"`), []byte("ip,host,title"))
	if err != nil {
		fmt.Printf("%v\n", err.Error())
		return
	}
	fmt.Printf("count: %d\n", len(data.Results))
	fmt.Printf("\n%s\n", data.String())
}
```

## FOFA Go Cli

基于 `FOFA Pro API`与`FOFA Go SDK` 编写的 `golang` 版 `命令行工具`,  方便技术人员更便捷地搜索、筛选、导出 `FOFA` 的数据。 

### 使用

#### 0x01 下载

直接下载即可使用,链接：<a href="https://github.com/xiaoyu-0814/fofa-go/releases/tag/v1.0.0"><font face="menlo">`FOFA Cli`</font></a>

#### 0x02 使用方法

在成功下载之后，可直接在终端下使用`fofa`命令，如下：

```fofa_cli
$ fofa_cli

    Fofa is a tool for discovering assets.

    Usage:

            fofa init|info|search option argument ...

    The options are:

            init:
                    email           the email which you login to fofa.so

                    key             the md5 string which you can find on userinfo page

            search:
                    fields          fields which you want to select
                                    Use ip,port,protocol as default.

                    format          output format
                                    Default is /t splice, you can choose other.

                    query           query statement which is similar to the statement used in the fofa.so

                    page            page number you want to query, 100 records per page
                                    If page is not set or page is less than 1, page will be set to 1.

                    out             output file path
                                    Print to the terminal as default.

                    count           only count the total number of matches,true or false
                                    False as default.


```

##### 1.初始化

 `邮箱(email)`和 `API KEY (key)`请在[FOFA官网](https://fofa.so/)--->个人中心--->个人资料查看。


```fofa_cli
$ fofa_cli init -email example@fofa.so -key 32charsMD5String
[+] Successfully initialized

Email：example@fofa.so
UserName：fofa
Fcoin：0
Vip：true
VipLevel：1
```

##### 2.个人信息

```fofa_cli
$ fofa_cli info
Email：example@fofa.so
UserName：fofa
Fcoin：0
Vip：true
VipLevel：1
```

##### 3.查询

注意：query中&&对应+，||对应-，query参数中多个条件时不能有空格和引号，其他与FOFA 网页查询语法相同，具体请到[FOFA官网](https://fofa.so/)查看。

###### 基本查询

```fofa_cli
$ fofa_cli search -query domain=163.com+port=443
103.254.188.71  443
59.111.18.135   443
59.111.137.212  443
......

total: 181
```



###### 指定返回字段与页数

字段默认值：

fields：ip,port,protocol

page：1

```fofa_cli
$ fofa_cli search -query domain=163.com-domain=126.com -fields ip,port,protocol,title -page 2
101.71.154.230  80      nil     301 Moved Permanently
42.186.69.125   80      nil     nil
123.126.96.212  80      nil     nil
2408:8719:5200::24      80      nil     301 Moved Permanently
59.111.0.134    80      nil     301 Moved Permanently
123.126.97.207  80      nil     系统提示
......

total: 1434
```



###### 结果保存到文本

```fofa_cli
$ fofa_cli search -query domain=163.com -out ./fofa
[+] Successfully
```



###### 自定义格式化输出

```fofa_cli
$ fofa_cli search -query domain=163.com -format ------
59.111.181.60------80
123.128.14.183------80
101.71.154.225------80
123.126.97.202------80
123.134.184.218------80
1.71.150.8------80
103.254.188.71------443
......

total: 1111
```



###### 查询匹配总数

```fofa_cli
$ fofa_cli search -query domain=163.com -count true

total: 1111
```



###### 查看cli版本

```
$ fofa_cli version

Version：1.0.0
```



## 协议

`FOFA SDK` 遵循 `MIT` 协议 <a href="https://opensource.org/licenses/mit"><font face="menlo">https://opensource.org/licenses/mit</font></a>
