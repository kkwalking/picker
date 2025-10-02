package models

import (
	"encoding/json"
	"os"
)

const PersistFile = "students.json"

// StudentList 管理学生名单的数据结构
type StudentList struct {
	Students []string
}

// NewStudentList 创建新的学生名单
func NewStudentList() *StudentList {
	return &StudentList{
		Students: []string{},
	}
}

// LoadFromFile 从文件加载学生名单
func LoadFromFile() (*StudentList, error) {
	data, err := os.ReadFile(PersistFile)
	if err != nil {
		return NewStudentList(), err
	}

	var names []string
	err = json.Unmarshal(data, &names)
	if err != nil {
		return NewStudentList(), err
	}

	return &StudentList{Students: names}, nil
}

// SaveToFile 保存学生名单到文件
func (sl *StudentList) SaveToFile() error {
	data, err := json.MarshalIndent(sl.Students, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(PersistFile, data, 0644)
}

// AddStudents 添加学生名单
func (sl *StudentList) AddStudents(names []string) {
	sl.Students = append(sl.Students, names...)
}

// Clear 清空学生名单
func (sl *StudentList) Clear() {
	sl.Students = []string{}
}

// IsEmpty 检查名单是否为空
func (sl *StudentList) IsEmpty() bool {
	return len(sl.Students) == 0
}

// Count 返回学生数量
func (sl *StudentList) Count() int {
	return len(sl.Students)
}