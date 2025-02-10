package docker

import (
	"bytes"
	"fmt"
	"net"
	"sync"
	log "github.com/sirupsen/logrus"
)

var (
	GlobalIPAllocator *IPAllocator
	// 默认的网段配置
	DefaultCIDR = "10.92.114.30/24"
)

// InitIPAllocator 初始化全局IP分配器
func InitIPAllocator(cidr string, usedIPs []string) error {
	if cidr == "" {
		cidr = DefaultCIDR
	}

	allocator, err := NewIPAllocator(cidr)
	if err != nil {
		return fmt.Errorf("初始化IP分配器失败: %v", err)
	}

	for _, ip := range usedIPs {
		allocator.RegisterIP(ip)
	}

	for ip := range allocator.usedIPs {
		fmt.Printf("已使用的ip: %v\n", ip)
	}

	GlobalIPAllocator = allocator
	return nil
}

// IPAllocator IP地址管理器
type IPAllocator struct {
	mu          sync.RWMutex
	usedIPs     map[string]bool
	network     *net.IPNet      // 网段信息
	startIP     net.IP          // 起始IP
	endIP       net.IP          // 结束IP
	totalIPs    int             // 可分配IP总数
	excludedIPs map[string]bool // 不允许分配的IP列表
}

// NewIPAllocator 创建新的IP分配器
// cidr 格式为 "IP/掩码"，例如 "10.92.114.30/24"
// IP 部分将作为起始分配点，1 到这个 IP 之前的所有 IP 都会被自动排除
func NewIPAllocator(cidr string) (*IPAllocator, error) {
	// 解析CIDR
	startIP, network, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, fmt.Errorf("无效的网段: %v", err)
	}

	// 确保起始IP在网段内
	if !network.Contains(startIP) {
		return nil, fmt.Errorf("起始IP %s 不在网段 %s 内", startIP, network)
	}

	// 计算可用IP数量
	ones, bits := network.Mask.Size()
	totalIPs := 1 << (bits - ones)
	// 排除网络地址和广播地址
	totalIPs -= 2

	// 获取网段的结束IP
	endIP := decrementIP(calculateBroadcastIP(network))

	// 如果起始IP在广播地址之后，返回错误
	if bytes.Compare(startIP, endIP) > 0 {
		return nil, fmt.Errorf("起始IP %s 超出网段 %s 的可用范围", startIP, network)
	}

	// 创建分配器
	allocator := &IPAllocator{
		usedIPs:     make(map[string]bool),
		network:     network,
		startIP:     startIP,
		endIP:       endIP,
		totalIPs:    totalIPs,
		excludedIPs: make(map[string]bool),
	}
	sIP := startIP.String()
	_ = sIP
	// 排除从网段第一个可用IP到起始IP之前的所有IP
	firstIP := incrementIP(network.IP)
	beforeStartIP := decrementIP(copyIP(startIP))
	allocator.AddExcludedIPRange(firstIP.String(), beforeStartIP.String())

	return allocator, nil
}

// AllocateIP 分配一个可用的IP地址
func (a *IPAllocator) AllocateIP() (string, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	currentIP := a.startIP
	for !currentIP.Equal(a.endIP) {
		ipStr := currentIP.String()

		// 跳过已使用的IP和被排除的IP
		if !a.usedIPs[ipStr] && !a.excludedIPs[ipStr] {
			if a.isValidIP(currentIP) {
				a.usedIPs[ipStr] = true
				return ipStr, nil
			}
		}

		currentIP = incrementIP(currentIP)
	}

	return "", fmt.Errorf("网段 %v 中没有可用的IP地址", a.network)
}

// RegisterIP 注册一个已使用的IP地址
// 如果IP已经被注册或不在网段内，将返回错误
func (a *IPAllocator) RegisterIP(ip string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	// 解析IP地址
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return fmt.Errorf("无效的IP地址格式: %s", ip)
	}

	// 检查IP是否在网段内
	if !a.network.Contains(parsedIP) {
		return fmt.Errorf("IP %s 不在网段 %s 内", ip, a.network)
	}

	// 检查IP是否已被使用或排除
	if a.usedIPs[ip] {
		return fmt.Errorf("IP %s 已被使用", ip)
	}
	if a.excludedIPs[ip] {
		// return fmt.Errorf("IP %s 在排除列表中", ip)
		log.Warnf("IP %s 在排除列表中，但被注册", ip)
	}

	// 注册IP
	a.usedIPs[ip] = true
	return nil
}

// ReleaseIP 释放一个IP地址
func (a *IPAllocator) ReleaseIP(ip string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if !a.usedIPs[ip] {
		return fmt.Errorf("IP %s 未被使用", ip)
	}

	delete(a.usedIPs, ip)
	return nil
}

// AddExcludedIP 添加不允许分配的IP
func (a *IPAllocator) AddExcludedIP(ip string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return fmt.Errorf("无效的IP地址: %s", ip)
	}

	if !a.network.Contains(parsedIP) {
		return fmt.Errorf("IP %s 不在网段 %v 内", ip, a.network)
	}

	a.excludedIPs[ip] = true
	a.totalIPs--
	return nil
}

// AddExcludedIPRange 添加一个范围的不允许分配的IP
// 例如：AddExcludedIPRange("10.92.114.1", "10.92.114.29") 将排除 1-29 的IP
func (a *IPAllocator) AddExcludedIPRange(startIP, endIP string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	start := net.ParseIP(startIP)
	if start == nil {
		return fmt.Errorf("无效的起始IP: %s", startIP)
	}

	end := net.ParseIP(endIP)
	if end == nil {
		return fmt.Errorf("无效的结束IP: %s", endIP)
	}

	// 确保IP在网段内
	if !a.network.Contains(start) || !a.network.Contains(end) {
		return fmt.Errorf("IP范围 %s-%s 不在网段 %v 内", startIP, endIP, a.network)
	}

	// 确保起始IP小于等于结束IP
	if bytes.Compare(start, end) > 0 {
		return fmt.Errorf("起始IP %s 大于结束IP %s", startIP, endIP)
	}

	current := start
	for !current.Equal(incrementIP(end)) {
		ipStr := current.String()
		if !a.excludedIPs[ipStr] {
			a.excludedIPs[ipStr] = true
			a.totalIPs--
		}
		current = incrementIP(current)
	}

	return nil
}

// RemoveExcludedIP 移除不允许分配的IP
func (a *IPAllocator) RemoveExcludedIP(ip string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if !a.excludedIPs[ip] {
		return fmt.Errorf("IP %s 不在排除列表中", ip)
	}

	delete(a.excludedIPs, ip)
	a.totalIPs++
	return nil
}

// GetExcludedIPs 获取所有被排除的IP
func (a *IPAllocator) GetExcludedIPs() []string {
	a.mu.RLock()
	defer a.mu.RUnlock()

	excludedList := make([]string, 0, len(a.excludedIPs))
	for ip := range a.excludedIPs {
		excludedList = append(excludedList, ip)
	}
	return excludedList
}

// isValidIP 检查IP是否在网段内且不是网络地址或广播地址
func (a *IPAllocator) isValidIP(ip net.IP) bool {
	return a.network.Contains(ip) &&
		!ip.Equal(a.network.IP) &&
		!ip.Equal(calculateBroadcastIP(a.network))
}

// incrementIP 获取下一个IP地址
func incrementIP(ip net.IP) net.IP {
	result := make(net.IP, len(ip))

	copy(result, ip)
	for i := len(result) - 1; i >= 0; i-- {
		result[i]++
		if result[i] > 0 {
			break
		}
	}
	return result
}

// decrementIP 获取前一个IP地址
func decrementIP(ip net.IP) net.IP {
	result := make(net.IP, len(ip))
	copy(result, ip)
	for i := len(result) - 1; i >= 0; i-- {
		result[i]--
		if result[i] < 255 {
			break
		}
	}
	return result
}

// calculateBroadcastIP 计算网段的广播地址
func calculateBroadcastIP(network *net.IPNet) net.IP {
	broadcastIP := make(net.IP, len(network.IP))
	copy(broadcastIP, network.IP)

	// 将IP与反掩码做或运算，得到广播地址
	for i := range broadcastIP {
		broadcastIP[i] |= ^network.Mask[i]
	}

	return broadcastIP
}

// copyIP 创建IP的副本
func copyIP(ip net.IP) net.IP {
	dup := make(net.IP, len(ip))
	copy(dup, ip)
	return dup
}
