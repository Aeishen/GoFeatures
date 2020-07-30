package main

import "sync"

/*
轮询（Round Robin）：顺序循环将请求一次顺序循环地连接每个服务器。当其中某个服务器发生第二到第7 层的故障，BIG-IP 就把其从顺序循环队列中拿出，不参加下一次的轮询，直到其恢复正常。

以轮询的方式依次请求调度不同的服务器；实现时，一般为服务器带上权重；这样有两个好处：

针对服务器的性能差异可分配不同的负载；
当需要将某个结点剔除时，只需要将其权重设置为0即可；

优点：实现简单、高效；易水平扩展
缺点：请求到目的结点的不确定，造成其无法适用于有写的场景（缓存，数据库写）
应用场景：数据库或应用服务层中只有读的场景
*/




// 普通轮询算法: 假设有ABC 3台机器，那么请求过来将按照 ABCABC 这样的顺顺序将请求反向代理到后端服务器, 原理是记录当前的index值，每次请求+1 取模
type LoadBalance  struct {
	*sync.RWMutex
	Index int
	Servers []string
}

func InitLoadBalance(servers []string) *LoadBalance{
	return &LoadBalance{
		new(sync.RWMutex),
		0,
		servers,
	}
}

func (l *LoadBalance)GetHttpServerByRoundRobin() (curServer string) {
	l.Lock()
	defer l.Unlock()
	curServer = l.Servers[l.Index]
	l.Index = (l.Index + 1) % len(l.Servers)
	return
}

func main()  {
	
}
