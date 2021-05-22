# Interview-test

# go-id
分布式id生成方法大致有：UUID，数据库自增，snowflake等

时间有序，使用标准库准备采用Snowflake算法
原生Snowflake算法会有时间回拨问题，会生成重复的ID

每个节点为一个Worker
```
type Worker struct {
	mu           sync.Mutex
	LastStamp    int64 // 记录时间戳
	WorkerID     int64
	DataCenterID int64
	Sequence     int64
}
```
NewWorker函数返回Worker实例，
getMilliSeconds方法用于获取当前时间，
NextID获取下一个ID，
- 先用getMilliSeconds获取当前时间戳
- 先把当前时间戳和上次ID 的时间戳进行比较，防止产生时间回溯的ID，
- 时间戳相同判断Sequence是否溢出，溢出则等待下一毫秒
- 新的时间戳的话sequence置0
- 记录下当前的时间戳
- 或运算生成ID并返回

# go-queue

# gscheduler

定时任务调度器，每个任务包含一个时间戳，任务需要在该时间点开始执行，精度为分钟级

定义一个Cron struct包含任务实例，运行状态，添加任务和停止指令

初始化后用go协程运行

在Cron上添加了Add，Stop和run方法
定义一个byTime切片，把Entry按时间排序
用select匹配执行任务时间点