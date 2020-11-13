package snow

import (
	"errors"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

//获取节点Ip
func GetNodeIP() string {

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}

	ipCount := len(addrs)
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ipNum := strings.Split(ipnet.IP.String(), ".")
				if strings.Compare("10", ipNum[0]) == 0 ||
					strings.Compare("172", ipNum[0]) == 0 ||
					(strings.Compare("192", ipNum[0]) == 0 && strings.Compare("168", ipNum[1]) == 0) {

					if ipCount > 1 && strings.Compare(ipNum[3], "1") != 0 {
						return ipnet.IP.String()
					}
				}
			}
		}
	}

	return ""
}

// 转成int64
func ConvertToIntIP(ip string) (int64, error) {
	ips := strings.Split(ip, ".")
	E := errors.New("Not A IP.")
	if len(ips) != 4 {
		return 0, E
	}
	var intIP int
	for k, v := range ips {
		i, err := strconv.Atoi(v)
		if err != nil || i > 255 {
			return 0, E
		}
		intIP = intIP | i<<uint(8*(3-k))
	}
	return int64(intIP), nil
}

// 因为snowFlake目的是解决分布式下生成唯一id 所以ID中是包含集群和节点编号在内的
const (
	numberBits uint8 = 12 // 表示每个集群下的每个节点，1毫秒内可生成的id序号的二进制位 对应上图中的最后一段
	workerBits uint8 = 10 // 每台机器(节点)的ID位数 10位最大可以有2^10=1024个节点数 即每毫秒可生成 2^12-1=4096个唯一ID 对应上图中的倒数第二段
	// 这里求最大值使用了位运算，-1 的二进制表示为 1 的补码，感兴趣的同学可以自己算算试试 -1 ^ (-1 << nodeBits) 这里是不是等于 1023
	//workerMax   int64 = -1 ^ (-1 << workerBits) // 节点ID的最大值，用于防止溢出
	numberMax   int64 = -1 ^ (-1 << numberBits) // 同上，用来表示生成id序号的最大值
	timeShift         = workerBits + numberBits // 时间戳向左的偏移量
	workerShift       = numberBits              // 节点ID向左的偏移量
	// 41位字节作为时间戳数值的话，大约68年就会用完
	// 假如你2010年1月1日开始开发系统 如果不减去2010年1月1日的时间戳 那么白白浪费40年的时间戳啊！
	// 这个一旦定义且开始生成ID后千万不要改了 不然可能会生成相同的ID
	epoch int64 = 1525705533000
)

var Snowflake *Worker

func SnowflakeInit() {
	Snowflake, _ = newWorker()
}

// 实例化一个工作节点
// workerId 为当前节点的id,通过IP算出来
func newWorker() (*Worker, error) {
	ip := GetNodeIP()
	workerId, _ := ConvertToIntIP(ip)
	// 要先检测workerId是否在上面定义的范围内
	if workerId < 0 {
		return nil, errors.New("Worker id excess of quantity")
	}
	// 生成一个新节点
	return &Worker{
		timestamp: 0,
		workerId:  workerId,
		number:    0,
	}, nil
}

// 定义一个worker工作节点所需要的基本参数
type Worker struct {
	mu        sync.Mutex // 添加互斥锁 确保并发安全
	timestamp int64      // 记录上一次生成id的时间戳
	workerId  int64      // 该节点的ID
	number    int64      // 当前毫秒已经生成的id序列号(从0开始累加) 1毫秒内最多生成4096个ID
}

// 生成方法一定要挂载在某个worker下，这样逻辑会比较清晰 指定某个节点生成id
func (w *Worker) GetStringId() string {
	intA := w.GetId()
	return strconv.FormatInt(intA, 10)
}
func (w *Worker) GetId() int64 {
	w.mu.Lock()
	defer w.mu.Unlock()

	// 获取生成时的时间戳
	now := time.Now().UnixNano() / 1e6 // 纳秒转毫秒
	if w.timestamp == now {
		w.number++

		// 这里要判断，当前工作节点是否在1毫秒内已经生成numberMax个ID
		if w.number > numberMax {
			// 如果当前工作节点在1毫秒内生成的ID已经超过上限 需要等待1毫秒再继续生成
			for now <= w.timestamp {
				now = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		// 如果当前时间与工作节点上一次生成ID的时间不一致 则需要重置工作节点生成ID的序号
		w.number = 0
		// 下面这段代码看到很多前辈都写在if外面，无论节点上次生成id的时间戳与当前时间是否相同 都重新赋值  这样会增加一丢丢的额外开销 所以我这里是选择放在else里面
		w.timestamp = now // 将机器上一次生成ID的时间更新为当前时间
	}

	ID := (now-epoch)<<timeShift | (w.workerId << workerShift) | (w.number)
	return ID
}
