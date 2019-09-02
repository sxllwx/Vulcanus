# vulcanus

To speed up the entire development process, Vulcanus provides some code generators or other commonly used code packages.


## code-gen

### Install Scaffold

```bash

export GO111MODULE=on
go install github.com/sxllwx/vulcanus/cmd/vulcanus

```
#### restful server

Vulcanus can help us generate restful style API-Server, the main use of the restful code package is "github.com/emicklei/go-restful"



#### gen restful container

vulcanus rest container -p {PKG_NAME} -k {RESOURCE_KIND}

PKG_NAME is the package name of the code package where the generated .go file is located
RESOURCE_KIND is the resource managed by this WebServer

#### gen restful ws

vulcanus rest ws -p {PKG_NAME} -k {RESOURCE_KIND}

PKG_NAME is the package name of the code package where the generated .go file is located
RESOURCE_KIND is the resource managed by this WebServer



