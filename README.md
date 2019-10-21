# vulcanus

为了加速整个开发流程（为了上班时间更方便的偷懒)，Vulcanus提供了一些我自己使用起来非常便利的工具

- 代码生成器
- 常用的代码包


## 代码生成器

当前已经支持生成

- restful server ( open-api&&swagger )

后续需要增加

- restful client
- orm (mysql, redis)



### 安装

```bash

export GO111MODULE=on
go install github.com/sxllwx/vulcanus/cmd/vulcanus

```
#### 生成Server侧代码

该代码生成工具主要是生成符合 [go-restful](https://github.com/emicklei/go-restful.git) 的代码

#### 生成 restful container

```bash
vulcanus rest container -p {PKG_NAME} -k {RESOURCE_KIND}
```

PKG_NAME 为生成的container.go 所在的代码包的包名(一般为main)
RESOURCE_KIND 为该 REST-style Server 管理的资源的类型(比如，Books啦, Users啦之类的)

#### 生成 restful webservice

```bash
vulcanus rest ws -p {PKG_NAME} -k {RESOURCE_KIND}
```

PKG_NAME 为生成的container.go 所在的代码包的包名(一般为api)
RESOURCE_KIND 为该 REST-style Server 管理的资源的类型(比如，Books啦, Users啦之类的)


#### 启动 http server

```go

func main(){

	c := NewContainer()

	// add web service
	m := NewbooksManager()
	c.Add(m.WebService())

	// regiser open api spec
	RegisterOpenAPI(c)

	if err := http.ListenAndServe(":8080", c); err != nil {
		panic(err)
	}
}

```

ok，```go build```

#### 为了让我们的REST-style server 更帅气，给他安排一下Swagger

docker run -it -p 8080:8080 -e API_URL=http://{你的IP}:8080/apidocs.json swaggerapi/swagger-ui

打开浏览器，帅气的REST-style的Server已经启动
