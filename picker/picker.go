package picker

import (
	"fmt"
	"math/rand"
	"time"
)

// Picker 点名器核心逻辑
type Picker struct {
	students []string
}

// NewPicker 创建新的点名器
func NewPicker(students []string) *Picker {
	return &Picker{
		students: students,
	}
}

// UpdateStudents 更新学生名单
func (p *Picker) UpdateStudents(students []string) {
	p.students = students
}

// HasStudents 检查是否有学生名单
func (p *Picker) HasStudents() bool {
	return len(p.students) > 0
}

// GetStudentCount 获取学生数量
func (p *Picker) GetStudentCount() int {
	return len(p.students)
}

// Pick 随机选择一个学生
func (p *Picker) Pick() string {
	if !p.HasStudents() {
		return ""
	}
	return p.students[rand.Intn(len(p.students))]
}

// StartAnimation 开始点名动画
func (p *Picker) StartAnimation(callback func(string), updateCallback func(string)) {
	if !p.HasStudents() {
		return
	}

	go func() {
		total := 7 * time.Second
		start := time.Now()
		interval := 50 * time.Millisecond

		var luckyOne string

		for time.Since(start) < total {
			luckyOne = p.Pick()
			if updateCallback != nil {
				updateCallback(luckyOne)
			}

			time.Sleep(interval)
			duration := time.Since(start)
			if duration > 3*time.Second && duration <= 4*time.Second {
				interval += 30 * time.Millisecond
			} else if duration > 4*time.Second && duration <= 5*time.Second {
				interval += 50 * time.Millisecond
			} else if duration > 5*time.Second {
				interval += 70 * time.Millisecond
			}
		}

		final := fmt.Sprintf("点到：%s", luckyOne)
		if callback != nil {
			callback(final)
		}
	}()
}