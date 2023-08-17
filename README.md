## kube-event-collector 集群内事件收集器
![](https://github.com/Kubernetes-Learning-Playground/kube-event-collector/blob/main/image/%E6%B5%81%E7%A8%8B%E5%9B%BE%20(1).jpg?raw=true)
### 项目思路与设计
设计背景：
在集群中k8s的event事件是通知资源对象的，用于记录系统的某些状态变化。使用自定义控制器的方式，监听
集群内的event事件，并进行后续操作，ex:(prometheus metrics, loki log, email message)

思路：自定义informer 控制器，并下发给对应的collector，执行各自流程。
### 项目功能

#### event-worker 
- 监听event资源
- 过滤特定event类型事件 (TODO)
- 记录event metrics，给后端prometheus
- 记录log结构化日志，给后端loki使用
- 下发qq邮箱
- 下发给elasticSearch

#### event-generator
- 模拟创建event事件

#### 配置文件
```yaml
kubeConfig: /Users/zhenyu.jiang/.kube/config   # k8s config目录，如果使用容器化部署，需要挂载kube config1
filterEventLevel:                               # 过滤的事件等级
elasticSearchEndpoint: http://127.0.0.1:9200    # es服务器endpoint
logFilePath: /Users/zhenyu.jiang/go/src/golanglearning/new_project/kube-event-collector   # 需要指定日志存放位置
# 通知模式：目前支持结构化日志、对接prometheus、发送email消息
mode:
  log: true
  prometheus: true
  message: true
  elasticSearch: false
# 邮箱配置信息 
sender:
  remote: smtp.qq.com
  port:  25
  email: <email>
  password: <password>
  targets: <targets>
```

### 项目测试
### event worker 
1. 进入项目根目录(--config <配置文件目录>)
```bash
➜  kube-event-collector git:(main) ✗ go run main.go kube-event-worker --config ./config.yaml
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
- using env:   export GIN_MODE=release
- using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /metrics                  --> github.com/practice/kube-event/pkg/server.PrometheusHandler.func1 (1 handlers)
[GIN-debug] [WARNING] You trusted all proxies, this is NOT safe. We recommend you to set a value.
Please check https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies for details.
[GIN-debug] Listening and serving HTTP on :8080

```

### event generator
需要在启动时指定必要参数，如果指定不正确，会报错
- kubeconfig kube config配置文件
- kind 触发事件的资源对象，ex: pods deployments
- name 事件名称
- namespace 
- type 事件程度：ex: Normal Warning等
- action 动作
- reason 事件根因 
- message 事件消息

1. 进入项目根目录
```bash
➜  kube-event-collector git:(main) ✗ go run main.go kube-event-generator --kubeconfig ~/.kube/config --kind pods --name kube-controller-manager-minikube --namespace kube-system
I0723 13:48:05.948372   32767 event_generator.go:115] Event generated successfully: 
&Event{ObjectMeta:{kube-controller-manager-minikube.17746911f0928000  kube-system  273659a2-efa3-484d-9404-0524cd869dee 2078909 0 2023-07-23 13:48:05 +0800 CST <nil> <nil> map[] map[] [] [] [{main Update v1 2023-07-23 13:48:05 +0800 CST FieldsV1 {"f:action":{},"f:eventTime":{},"f:firstTimestamp":{},"f:involvedObject":{},"f:lastTimestamp":{},"f:message":{},"f:reason":{},"f:reportingComponent":{},"f:reportingInstance":{},"f:type":{}} }]},InvolvedObject:ObjectReference{Kind:Pod,Namespace:kube-system,Name:kube-controller-manager-minikube,UID:5fdee821-fbb9-4933-8e25-5f34a0a42357,APIVersion:v1,ResourceVersion:2059197,FieldPath:,},Reason:Testing-Reason,Message:Testing-Message,Source:EventSource{Component:,Host:,},FirstTimestamp:2023-07-23 13:48:05 +0800 CST,LastTimestamp:2023-07-23 13:48:05 +0800 CST,Count:0,Type:Warning,EventTime:2023-07-23 13:48:05.942272 +0800 CST,Series:nil,Action:ttt,Related:nil,ReportingController:k8s-event-generator,ReportingInstance:k8s-event-generator,}
# 获取事件资源
➜  kube-event-collector git:(main) ✗ kubectl get event -A
NAMESPACE     LAST SEEN   TYPE      REASON               OBJECT                        MESSAGE
kube-system   6m9s        Warning   Testing-Reason       pod/kube-scheduler-minikube   Testing-Message
kube-system   5m7s        Warning   Testing-Reasonaaaa   pod/kube-scheduler-minikube   Testing-Message
```
