package main

import (
	"fmt"
	"math/rand"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/kkwalking/picker/importer"
	"github.com/kkwalking/picker/models"
	"github.com/kkwalking/picker/picker"
	"github.com/kkwalking/picker/ui"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	a := app.New()
	w := a.NewWindow("课堂点名器")
	w.Resize(fyne.NewSize(600, 420))

	// 初始化组件
	display := ui.NewDisplay()
	studentList := loadOrCreateStudentList()
	pickerInstance := picker.NewPicker(studentList.Students)
	importerInstance := importer.NewStudentImporter(w)

	// 更新显示状态
	updateDisplayState(studentList, display)

	// 导入按钮
	btnImport := widget.NewButton("导入名单", func() {
		importerInstance.ShowImportDialog(func(students []string) {
			studentList.AddStudents(students)
			err := studentList.SaveToFile()
			if err != nil {
				dialog.ShowError(fmt.Errorf("保存名单失败: %v", err), w)
				return
			}

			pickerInstance.UpdateStudents(studentList.Students)
			dialog.ShowInformation("成功", fmt.Sprintf("导入成功，共 %d 名学生", len(students)), w)
			updateDisplayState(studentList, display)
		})
	})

	// 开始点名按钮
	btnStart := ui.NewPrimaryButton("开始点名", func() {
		if !pickerInstance.HasStudents() {
			dialog.ShowInformation("提示", "请先导入名单", w)
			return
		}

		// 这里需要重新获取按钮容器来控制显示/隐藏
		// 由于UI组件的复杂性，我们简化处理，不隐藏按钮

		pickerInstance.StartAnimation(
			// 动画结束回调
			func(result string) {
				display.ShowResult(result)
				display.Pulse()
			},
			// 动画更新回调
			func(name string) {
				display.UpdateText(name, 32, false)
			},
		)
	})

	// 布局
	buttons := container.NewHBox(btnImport, btnStart)

	content := container.NewVBox(
		layout.NewSpacer(),
		display.Container(),
		layout.NewSpacer(),
		buttons,
	)

	w.SetContent(content)
	w.ShowAndRun()
}

// loadOrCreateStudentList 加载或创建学生名单
func loadOrCreateStudentList() *models.StudentList {
	studentList, err := models.LoadFromFile()
	if err == nil && !studentList.IsEmpty() {
		fmt.Println("已从持久化文件加载名单")
		return studentList
	}

	fmt.Println("未找到持久化文件，请导入 Excel")
	return models.NewStudentList()
}

// updateDisplayState 更新显示状态
func updateDisplayState(studentList *models.StudentList, display *ui.Display) {
	if studentList.IsEmpty() {
		display.ShowText("请先导入学生名单")
	} else {
		display.ShowText("名单已加载，可以开始点名")
	}
}