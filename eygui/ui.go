package eygui

import (
	"github.com/ying32/govcl/vcl"
	"github.com/ying32/govcl/vcl/types"
	"github.com/ying32/govcl/vcl/types/colors"
)

func InitURLLabel(label vcl.TLabel) {
	label.Font().SetColor(colors.ClBlue)
	label.SetCursor(types.CrHandPoint)
	label.SetOnMouseEnter(func(sender vcl.IObject) {
		label.Font().SetColor(colors.ClRed)
		style := label.Font().Style()
		label.Font().SetStyle(style.Include(types.FsUnderline))
	})

	label.SetOnMouseLeave(func(sender vcl.IObject) {
		label.Font().SetColor(colors.ClBlue)
		style := label.Font().Style()
		label.Font().SetStyle(style.Exclude(types.FsUnderline))
	})
}
