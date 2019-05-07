- 不使用现成爬虫库/框架
- 使用ElasticSearch作为数据存储
- 使用Go语言标准模板库实现http数据展示部分


单任务 --> 并发版 --> 分布式

单任务
- 获取并打印所有城市第一页用户的详细信息

并发版
- Scheduler 1 : 所有worker公用一个输入（出现循环等待的问题）
- Scheduler 2 : 并发分发Requets (无法控制worker)
- Scheduler 3 : Request队列和Worker队列
- 正则表达式


URl去重(已经访问)
- 哈希表
- 