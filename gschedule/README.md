# gscheduler

定时任务调度器，每个任务包含一个时间戳，任务需要在该时间点开始执行，精度为分钟级

定义一个Cron struct包含任务实例，运行状态，添加任务和停止指令

初始化后用go协程运行

在Cron上添加了Add，Stop和run方法
定义一个byTime切片，把Entry按时间排序
用select匹配执行任务时间点