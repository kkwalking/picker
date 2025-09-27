package main

import (
	"fmt"
	"math/rand"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"fyne.io/fyne/v2/storage"
	"github.com/xuri/excelize/v2"
)

var students []string

func main() {
	a := app.New()
	w := a.NewWindow("课堂点名器")
	w.Resize(fyne.NewSize(400, 300))

	label := widget.NewLabel("请先导入学生名单")
	label.Alignment = fyne.TextAlignCenter

	btnImport := widget.NewButton("导入名单", func() {
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if reader == nil {
				return
			}
			defer reader.Close()

			file, err := excelize.OpenReader(reader)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}

			rows, err := file.GetRows("Sheet1")
			if err != nil {
				dialog.ShowError(err, w)
				return
			}

			students = []string{}
			for _, row := range rows {
				if len(row) > 0 {
					students = append(students, row[0])
				}
			}
			label.SetText(fmt.Sprintf("已导入 %d 名学生", len(students)))
		}, w)
		fd.SetFilter(storage.NewExtensionFileFilter([]string{".xlsx", ".xls"}))
		fd.Show()
	})

	btnStart := widget.NewButton("开始点名", func() {
		if len(students) == 0 {
			dialog.ShowInformation("提示", "请先导入名单", w)
			return
		}

		go func() {
			totalTime := 5 * time.Second
			start := time.Now()
			interval := 50 * time.Millisecond

			for time.Since(start) < totalTime {
				idx := rand.Intn(len(students))
				name := students[idx]

				a.SendNotification(&fyne.Notification{
					Title:   "当前候选",
					Content: name,
				})

				label.SetText(name)
				time.Sleep(interval)

				// 模拟减速
				if time.Since(start) > 3*time.Second {
					interval += 30 * time.Millisecond
				}
			}

			// 最终固定一个名字
			final := students[rand.Intn(len(students))]
			label.SetText(fmt.Sprintf("最终点到：%s", final))
		}()
	})

	content := container.NewVBox(
		label,
		btnImport,
		btnStart,
	)

	w.SetContent(content)
	w.ShowAndRun()
}
