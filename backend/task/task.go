package task

import (
	"fmt"
	"sync"
)

// Task 是一个任务函数类型
type Task func() error

// GlobalAsyncManager 是全局异步任务管理器
type GlobalAsyncManager struct {
	tasks    chan Task
	wg       sync.WaitGroup
	stopChan chan struct{}
}

var (
	managerInstance *GlobalAsyncManager
	once            sync.Once
)

// GetAsyncTaskManager 获取全局异步任务管理器实例
func GetAsyncTaskManager() *GlobalAsyncManager {
	if managerInstance == nil {
		panic("GlobalAsyncManager not initialized")
	}
	return managerInstance
}

// InitializeManager 初始化全局异步任务管理器
func InitializeGlobalManager(bufferSize, workerCount int) *GlobalAsyncManager {
	once.Do(func() {
		managerInstance = &GlobalAsyncManager{
			tasks:    make(chan Task, bufferSize),
			stopChan: make(chan struct{}),
		}
		managerInstance.run(workerCount)
	})
	return managerInstance
}

// run 启动 worker
func (m *GlobalAsyncManager) run(workerCount int) {
	for i := 0; i < workerCount; i++ {
		go m.worker()
	}
}

// worker 是任务消费者
func (m *GlobalAsyncManager) worker() {
	for {
		select {
		case task := <-m.tasks:
			if task != nil {
				if err := task(); err != nil {
					fmt.Printf("Task error: %v\n", err)
				}
				m.wg.Done()
			}
		case <-m.stopChan:
			return
		}
	}
}

// AddTask 添加任务到队列
func (m *GlobalAsyncManager) AddTask(task Task) {
	m.wg.Add(1)
	m.tasks <- task
}

// Stop 停止任务管理器
func (m *GlobalAsyncManager) Stop() {
	close(m.stopChan)
	m.wg.Wait() // 等待所有任务完成
	close(m.tasks)
}
