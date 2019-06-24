# qtask
定时将任务放入消息队列，并提供延时，处理中，完结等控制任务处理过程。


示例:

* １. 创建任务

```go
//业务逻辑

//创建超时任务
qtask.Create(c,"订单绑定任务",map[string]interface{}{
    "order_no":"8973097380045"
},3600,"GCR:ORDER:BIND")

```


* 2. 编写MQC服务，该服务处理 `GCR:ORDER:BIND`消息队列数据

```go

func OrderBind(ctx *context.Context) (r interface{}) {
    //检查输入参数...
    
    //将任务修改为正在处理中
    qtask.Processing(ctx,ctx.Request.GetInt64("task_id"))


    //处理业务逻辑...


    //成功处理，结束任务
    qtask.Finish(ctx,ctx.Request.GetInt64("task_id"))
}

```

* 3. 定时任务

a. 注册服务

```go
app.Cron("/task/scan",qtask.Scan) //定时扫描任务
app.Cron("/task/clear",qtask.Clear(7)) //定时清理任务，删除７天的任务数据

app.MQC("/order/bind",OrderBind) //消息处理任务

```

b. 配置定时执行
```sh
app.Conf.CRON.SetSubConf("task", `{
		"tasks": [
			{
				"cron": "@every 10s",
				"service": "/task/scan"
			},
			{
				"cron": "@every 10s",
				"service": "/task/clear"
			}
        ]
    }`)
```


