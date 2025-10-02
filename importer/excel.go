package importer

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"

	"github.com/xuri/excelize/v2"
)

// StudentImporter 学生名单导入器
type StudentImporter struct {
	window fyne.Window
}

// NewStudentImporter 创建新的导入器
func NewStudentImporter(window fyne.Window) *StudentImporter {
	return &StudentImporter{
		window: window,
	}
}

// ShowImportDialog 显示导入对话框
func (si *StudentImporter) ShowImportDialog(onSuccess func([]string)) {
	fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if reader == nil {
			return
		}
		defer reader.Close()

		students, err := si.importFromReader(reader)
		if err != nil {
			dialog.ShowError(err, si.window)
			return
		}

		if onSuccess != nil {
			onSuccess(students)
		}
	}, si.window)

	fd.SetFilter(storage.NewExtensionFileFilter([]string{".xlsx", ".xls"}))
	fd.Show()
}

// importFromReader 从读取器导入学生名单
func (si *StudentImporter) importFromReader(reader fyne.URIReadCloser) ([]string, error) {
	f, err := excelize.OpenReader(reader)
	if err != nil {
		return nil, err
	}

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return nil, err
	}

	students := []string{}
	for _, row := range rows {
		if len(row) > 0 && row[0] != "" {
			students = append(students, row[0])
		}
	}

	if len(students) == 0 {
		return nil, fmt.Errorf("未在Excel文件中找到学生名单")
	}

	return students, nil
}