package main

import (
	"fmt"
	"slices"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"encoding/json"
	"time"
)

// Struct used to import and export todo list as JSON
// TODO refactor whole thing to make single struct per entry and store as JSON on device instead of string slices
type todoDataImportExportStruct struct {
	TodoDataShortVal []string `json:"todoDataShort"`
	TodoDataLongVal []string `json:"todoDataLong"`
}

var list *widget.List // Main screen list widget
var todoDataShort []string // Entry titles
var todoDataLong []string // Entry descriptions
var a fyne.App // App instance
var titleLabel *widget.Label // Top bar "title" label widget (used for errors and messages)

// Returns default top bar message with amount of entries in list
func getUpdatedTitleLabel() string {
	return fmt.Sprintf("%v elements in your TODO list!", len(todoDataShort))
}

// Goroutine that starts a timer after custom error/message is shown on top bar
// TODO refactor to single function that shows specified error/message and starts timer
// TODO stop running timer if new error/message is shown before current timer runs out
func updateTitleLabelAfterMessage() {
	timer := time.NewTimer(5*time.Second)
	<-timer.C
	fyne.Do(func() {
		titleLabel.SetText(getUpdatedTitleLabel())
		titleLabel.Refresh()
	})
}

// Export list as JSON string
func exportJson() string {
	dataStruct := todoDataImportExportStruct{}
	for _, v := range todoDataShort {
		dataStruct.TodoDataShortVal = append(dataStruct.TodoDataShortVal, v)
	}
	for _, v := range todoDataLong {
		dataStruct.TodoDataLongVal = append(dataStruct.TodoDataLongVal, v)
	}
	retval, err := json.Marshal(dataStruct)
	if err != nil {
		titleLabel.SetText("JSON Parsing error!")
		titleLabel.Refresh()
		go updateTitleLabelAfterMessage()
		return ""
	}
	titleLabel.SetText("JSON copied to clipboard!")
	titleLabel.Refresh()
	go updateTitleLabelAfterMessage()
	return string(retval)
}

func main() {
	a = app.NewWithID("pl.mdukat.dukatodo")
	w := a.NewWindow("TODO list app!")
	todoDataShort = make([]string, 0)
	todoDataLong = make([]string, 0)

	// Load data into memory
	todoLoadedShort := a.Preferences().StringList("todoListShort")
	for _, v := range todoLoadedShort {
		todoDataShort = append(todoDataShort, v)
	}
	todoLoadedLong := a.Preferences().StringList("todoListLong")
	for _, v := range todoLoadedLong {
		todoDataLong = append(todoDataLong, v)
	}

	// Backwards compatibility
	// Version 1.0 supported only titles of entries
	backCompatTodoLoaded := a.Preferences().StringList("todoList")
	if len(backCompatTodoLoaded) > 0 {
		for _, v := range backCompatTodoLoaded {
			todoDataShort = append(todoDataShort, v)
			todoDataLong = append(todoDataLong, "")
		}
		a.Preferences().SetStringList("todoListShort", todoDataShort)
		a.Preferences().SetStringList("todoListLong", todoDataLong)
		a.Preferences().RemoveValue("todoList")
	}

	// Main window
	titleLabel = widget.NewLabel(getUpdatedTitleLabel())
	exportJsonButton := widget.NewButtonWithIcon("", theme.LogoutIcon(), func() {
		cboard := a.Clipboard()
		jsonContent := exportJson()
		if len(jsonContent) > 0 {
			cboard.SetContent(jsonContent)
		}})
	importJsonButton := widget.NewButtonWithIcon("", theme.LoginIcon(), func() {
		// Import JSON window
		// TODO move me to dedicated function
		importJsonWindow := a.NewWindow("Import JSON")
		importJsonData := widget.NewMultiLineEntry()
		importJsonForm := &widget.Form{
			Items: []*widget.FormItem{
				{Text: "JSON", Widget: importJsonData}},
			OnSubmit: func() {
				// Validate JSON
				if json.Valid([]byte(importJsonData.Text)) != true {
					titleLabel.SetText("Invalid JSON structure!")
					titleLabel.Refresh()
					go updateTitleLabelAfterMessage()
					importJsonWindow.Close()
					return
				}
				// Check if length of titles and descriptions is the same
				var importedJson todoDataImportExportStruct
				err := json.Unmarshal([]byte(importJsonData.Text), &importedJson)
				if err != nil {
					titleLabel.SetText(fmt.Sprintf("Failed to parse JSON: %s", err.Error()))
					titleLabel.Refresh()
					go updateTitleLabelAfterMessage()
					importJsonWindow.Close()
					return
				}
				if len(importedJson.TodoDataShortVal) != len(importedJson.TodoDataLongVal) {
					titleLabel.SetText("JSON error: wrong length of arrays")
					titleLabel.Refresh()
					go updateTitleLabelAfterMessage()
					importJsonWindow.Close()
					return
				}
				// Add imported values to global slices
				for _, v := range importedJson.TodoDataShortVal {
					todoDataShort = append(todoDataShort, v)
				}
				for _, v := range importedJson.TodoDataLongVal {
					todoDataLong = append(todoDataLong, v)
				}
				// Save data to device
				// TODO move me to dedicated function
				a.Preferences().SetStringList("todoListShort", todoDataShort)
				a.Preferences().SetStringList("todoListLong", todoDataLong)
				list.UnselectAll()
				list.Refresh()
				titleLabel.SetText(getUpdatedTitleLabel())
				titleLabel.Refresh()
				importJsonWindow.Close()
			},
		}

		importJsonWindow.SetContent(importJsonForm)
		importJsonWindow.ShowAndRun()
	})
	
	topBar := container.NewHBox(importJsonButton, exportJsonButton, titleLabel)
	
	list = &widget.List{
		Length: func() int {
			return len(todoDataShort)
		},
		CreateItem: func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		UpdateItem: func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(todoDataShort[i])
		},
		OnSelected: func(i widget.ListItemID) {
			// Entry edit window
			// TODO move me to dedicated function
			entryEditWindow := a.NewWindow("Edit entry")
			entryEditWindow.SetOnClosed(func() {
				list.UnselectAll()
			})
			entryEditShort := widget.NewEntry()
			entryEditShort.Append(todoDataShort[i])
			entryEditLong := widget.NewMultiLineEntry()
			entryEditLong.Append(todoDataLong[i])
			entryEditDeleteButton := widget.NewButton("Delete entry", func() {
				todoDataShort = slices.Delete(todoDataShort, i, i+1)
				todoDataLong = slices.Delete(todoDataLong, i, i+1)
				// Save data to device
				// TODO move me to dedicated function
				a.Preferences().SetStringList("todoListShort", todoDataShort)
				a.Preferences().SetStringList("todoListLong", todoDataLong)
				list.UnselectAll()
				list.Refresh()
				titleLabel.SetText(getUpdatedTitleLabel())
				titleLabel.Refresh()
				entryEditWindow.Close()
			})

			entryEditForm := &widget.Form{
				Items: []*widget.FormItem{
					{Text: "Entry", Widget: entryEditShort},
					{Text: "Desc.", Widget: entryEditLong},
					{Text: "", Widget: entryEditDeleteButton}},
				OnSubmit: func() {
					todoDataShort[i] = entryEditShort.Text
					todoDataLong[i] = entryEditLong.Text
					// Save data to device
					// TODO move me to dedicated function
					a.Preferences().SetStringList("todoListShort", todoDataShort)
					a.Preferences().SetStringList("todoListLong", todoDataLong)
					list.UnselectAll()
					list.RefreshItem(i)
					entryEditWindow.Close()
				},
				OnCancel: func() {
					list.UnselectAll()
					entryEditWindow.Close()
				},
			}

			entryEditWindow.SetContent(entryEditForm)
			entryEditWindow.ShowAndRun()

		},
	}
	
	newElementButton := widget.NewButton("Add new TODO element", func() {
		// Add new element window
		// TODO move me to dedicated function
		newElementWindow := a.NewWindow("New entry")
		newElementEntry := widget.NewEntry()
		newElementDescription := widget.NewMultiLineEntry()
		newElementForm := &widget.Form{
			Items: []*widget.FormItem{
				{Text: "Entry", Widget: newElementEntry},
				{Text: "Desc.", Widget: newElementDescription}},
			OnSubmit: func() {
				todoDataShort = append(todoDataShort, newElementEntry.Text)
				todoDataLong = append(todoDataLong, newElementDescription.Text)
				// Save data to device
				// TODO move me to dedicated function
				a.Preferences().SetStringList("todoListShort", todoDataShort)
				a.Preferences().SetStringList("todoListLong", todoDataLong)
				list.Refresh()
				titleLabel.SetText(getUpdatedTitleLabel())
				titleLabel.Refresh()
				newElementWindow.Close()
			},
		}
		newElementWindow.SetContent(newElementForm)
		newElementWindow.ShowAndRun()
	})

	content := container.NewBorder(topBar, newElementButton, nil, nil, list)
	w.SetContent(content)
	w.ShowAndRun()
}
