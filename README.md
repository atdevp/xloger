# 离线日志采集工具

## 流程图

![github](image/xloger.png)

## 节点功能
* host
> 应用程序产生原始日志
* agent-m
> 采集主进程

* agent-s
> 采集备进程，若主进程有失败节点，则备进程补录

* HDFS
> 数据归档存储

* Mysql
> 元数据存储

* Web Operation Platform
> 人工操作平台进行任务的CRUD操作

