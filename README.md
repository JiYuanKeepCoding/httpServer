# httpServer
## 功能
/healthz 返回200

/env 获取当前环境


## 本地启动
```
go build
./httpServer
```

## 部署
```
make build-docker
make deploy
```