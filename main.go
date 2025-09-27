// main.go
package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"math/rand"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"fyne.io/fyne/v2/storage"
	"github.com/xuri/excelize/v2"
)

var students []string

// tappableRect 在 canvas.Rectangle 基础上实现 fyne.Tappable
type tappableRect struct {
	*canvas.Rectangle
	tapped func()
}

func (t *tappableRect) Tapped(*fyne.PointEvent) {
	if t.tapped != nil {
		t.tapped()
	}
}
func (t *tappableRect) TappedSecondary(*fyne.PointEvent) {}

func main() {
	rand.Seed(time.Now().UnixNano())

	a := app.New()
	w := a.NewWindow("课堂点名器")
	w.Resize(fyne.NewSize(600, 420))

	// 大字号显示区域（canvas.Text）
	display := canvas.NewText("请先导入学生名单", color.NRGBA{R: 0, G: 0, B: 0, A: 255})
	display.Alignment = fyne.TextAlignCenter
	display.TextSize = 32

	students, err := loadStudentsFromFile()
	if err == nil && len(students) > 0 {
		fmt.Println("已从持久化文件加载名单")
		display.Text = "名单已加载，可以开始点名"

	} else {
		fmt.Println("未找到持久化文件，请导入 Excel")
	}

	// 导入按钮（用默认按钮即可）
	btnImport := widget.NewButton("导入名单", func() {
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if reader == nil {
				return
			}
			defer reader.Close()

			f, err := excelize.OpenReader(reader)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}

			rows, err := f.GetRows("Sheet1")
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
			// save to local file
			if len(students) > 0 {
				saveStudentsToFile(students)
			}

			dialog.ShowInformation("成功", fmt.Sprintf("导入成功，共 %d 名学生", len(students)), w)
			display.Text = "名单已导入，可以开始点名"
			display.TextSize = 28
			display.TextStyle = fyne.TextStyle{Bold: false}
			canvas.Refresh(display)
		}, w)

		fd.SetFilter(storage.NewExtensionFileFilter([]string{".xlsx", ".xls"}))
		fd.Show()
	})

	// 自定义“现代绿色”按钮（带按下动画）
	btnStart := newModernButton("开始点名", func() {
		// 点击处理逻辑（在主线程以外运行点名动画）
		if len(students) == 0 {
			dialog.ShowInformation("提示", "请先导入名单", w)
			return
		}

		go func() {
			total := 7 * time.Second
			start := time.Now()
			interval := 50 * time.Millisecond

			var luckyOne string

			for time.Since(start) < total {
				idx := rand.Intn(len(students))
				luckyOne = students[idx]
				// 更新 display（需要刷新）
				display.Text = luckyOne
				display.TextSize = 32
				display.TextStyle = fyne.TextStyle{Bold: false}
				canvas.Refresh(display)

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

			final := luckyOne
			display.Text = fmt.Sprintf("点到：%s", final)
			display.TextSize = 36
			display.TextStyle = fyne.TextStyle{Bold: true}
			canvas.Refresh(display)

			// 最终的简短放大效果（可视化强调）
			orig := display.TextSize
			for i := 0; i < 2; i++ {
				display.TextSize = orig + 8
				canvas.Refresh(display)
				time.Sleep(120 * time.Millisecond)
				display.TextSize = orig
				canvas.Refresh(display)
				time.Sleep(80 * time.Millisecond)
			}
		}()
	})

	// 两个按钮放在同一行
	buttons := container.NewHBox(
		btnImport,
		btnStart,
	)

	// 用透明矩形作为“高度占位器”，确保名字区域固定高度（约等于3行按钮）
	filler := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0}) // 透明占位
	filler.SetMinSize(fyne.NewSize(0, 120))                // 高度约为 3 行按钮

	// 名字显示区域：填充占位 + 垂直居中显示名字
	nameArea := container.NewMax(
		filler,
		container.NewVBox(
			layout.NewSpacer(),
			container.NewCenter(display),
			layout.NewSpacer(),
		),
	)

	// 页面整体布局
	content := container.NewVBox(
		layout.NewSpacer(),
		nameArea,
		layout.NewSpacer(),
		buttons,
	)

	w.SetContent(content)
	w.ShowAndRun()
}

// newModernButton 创建一个绿色圆角按钮（带按下动画），返回 fyne.CanvasObject
func newModernButton(text string, onClick func()) fyne.CanvasObject {
	normalColor := color.NRGBA{R: 33, G: 148, B: 83, A: 255}  // #219453
	pressedColor := color.NRGBA{R: 27, G: 122, B: 69, A: 255} // 深一点

	// 背景矩形，决定按钮大小
	rect := canvas.NewRectangle(normalColor)
	rect.CornerRadius = 8.0
	rect.SetMinSize(fyne.NewSize(140, 44)) // ✅ 用矩形设置尺寸

	tRect := &tappableRect{Rectangle: rect}

	label := canvas.NewText(text, color.NRGBA{R: 255, G: 255, B: 255, A: 255})
	label.Alignment = fyne.TextAlignCenter
	label.TextSize = 16

	// 点击回调（带短动画）
	tRect.tapped = func() {
		go func() {
			// 按下效果
			rect.FillColor = pressedColor
			label.TextSize = 15
			canvas.Refresh(rect)
			canvas.Refresh(label)

			time.Sleep(140 * time.Millisecond)

			// 恢复
			rect.FillColor = normalColor
			label.TextSize = 16
			canvas.Refresh(rect)
			canvas.Refresh(label)

			if onClick != nil {
				onClick()
			}
		}()
	}

	// ✅ 背景矩形控制大小，容器里放 label 和 tappable 区域
	btnContainer := container.NewMax(rect, container.NewCenter(label), tRect)
	return btnContainer
}

var persistFile = "students.json"

// 保存名单到 JSON 文件
func saveStudentsToFile(names []string) error {
	data, err := json.MarshalIndent(names, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(persistFile, data, 0644)
}

// 从 JSON 文件读取名单
func loadStudentsFromFile() ([]string, error) {
	data, err := os.ReadFile(persistFile)
	if err != nil {
		return nil, err
	}
	var names []string
	err = json.Unmarshal(data, &names)
	return names, err
}
