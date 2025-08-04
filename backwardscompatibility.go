package main

import (
	"fyne.io/fyne/v2"
)

func RunBackwardsCompatibilityMigration() {
	// Backwards compatibility

	// 1.0 -> 1.2
	// Version 1.0 supported only titles of entries
	backCompatTodoLoaded := fyne.CurrentApp().Preferences().StringList("todoList")
	if len(backCompatTodoLoaded) > 0 {
		for _, v := range backCompatTodoLoaded {
			todoDataShort = append(todoDataShort, v)
			todoDataLong = append(todoDataLong, "")
		}
		fyne.CurrentApp().Preferences().SetStringList("todoListShort", todoDataShort)
		fyne.CurrentApp().Preferences().SetStringList("todoListLong", todoDataLong)
		fyne.CurrentApp().Preferences().RemoveValue("todoList")
	}
}
