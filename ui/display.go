package ui

import (
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

// Display 显示组件，用于显示点名结果
type Display struct {
	text       *canvas.Text
	container  fyne.CanvasObject
}

// NewDisplay 创建新的显示组件
func NewDisplay() *Display {
	text := canvas.NewText("请先导入学生名单", color.NRGBA{R: 0, G: 0, B: 0, A: 255})
	text.Alignment = fyne.TextAlignCenter
	text.TextSize = 32

	// 用透明矩形作为"高度占位器"，确保名字区域固定高度
	filler := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	filler.SetMinSize(fyne.NewSize(0, 120))

	// 名字显示区域：填充占位 + 垂直居中显示名字
	nameArea := container.NewMax(
		filler,
		container.NewVBox(
			layout.NewSpacer(),
			container.NewCenter(text),
			layout.NewSpacer(),
		),
	)

	return &Display{
		text:      text,
		container: nameArea,
	}
}

// ShowText 显示文本
func (d *Display) ShowText(text string) {
	d.text.Text = text
	d.text.TextSize = 32
	d.text.TextStyle = fyne.TextStyle{Bold: false}
	canvas.Refresh(d.text)
}

// ShowResult 显示最终结果
func (d *Display) ShowResult(studentName string) {
	d.text.Text = studentName
	d.text.TextSize = 36
	d.text.TextStyle = fyne.TextStyle{Bold: true}
	canvas.Refresh(d.text)
}

// UpdateText 更新显示文本（带样式）
func (d *Display) UpdateText(text string, textSize float32, bold bool) {
	d.text.Text = text
	d.text.TextSize = textSize
	if bold {
		d.text.TextStyle = fyne.TextStyle{Bold: true}
	} else {
		d.text.TextStyle = fyne.TextStyle{Bold: false}
	}
	canvas.Refresh(d.text)
}

// Pulse 创建脉冲动画效果
func (d *Display) Pulse() {
	go func() {
		orig := d.text.TextSize
		for i := 0; i < 2; i++ {
			d.text.TextSize = orig + 8
			canvas.Refresh(d.text)
			time.Sleep(120 * time.Millisecond)
			d.text.TextSize = orig
			canvas.Refresh(d.text)
			time.Sleep(80 * time.Millisecond)
		}
	}()
}

// Container 返回显示组件的容器
func (d *Display) Container() fyne.CanvasObject {
	return d.container
}