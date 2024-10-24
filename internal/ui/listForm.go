package ui

import (
	"fmt"
	"strconv"

	"github.com/rivo/tview"
)

type ListForm struct {
	list        []string
	listView    *tview.TextView
	newItem     *tview.InputField
	updateIndex *tview.InputField
	updateItem  *tview.InputField
	removeIndex *tview.InputField
}

func newListForm(list []string) *ListForm {
	listView := tview.NewTextView().
		SetDynamicColors(true).
		SetText("Current List: " + fmt.Sprint(list))

	return &ListForm{
		list:        list,
		listView:    listView,
		newItem:     tview.NewInputField().SetLabel("Add Item: ").SetFieldWidth(20),
		updateIndex: tview.NewInputField().SetLabel("Update Index: ").SetFieldWidth(5),
		updateItem:  tview.NewInputField().SetLabel("New Value: ").SetFieldWidth(20),
		removeIndex: tview.NewInputField().SetLabel("Remove Index: ").SetFieldWidth(5),
	}
}

func (manager *ListForm) attachToForm(form *tview.Form) {
	form.AddFormItem(manager.newItem).
		AddButton("Add destination", manager.addItem).
		AddFormItem(manager.updateIndex).
		AddFormItem(manager.updateItem).
		AddButton("Update destination", manager.updateItemFunc).
		AddFormItem(manager.removeIndex).
		AddButton("Remove destination", manager.removeItem)
}

func (listForm *ListForm) getListView() *tview.TextView {
	return listForm.listView
}

func (listForm *ListForm) addItem() {
	newItem := listForm.newItem.GetText()
	if newItem != "" {
		listForm.list = append(listForm.list, newItem)
		listForm.refreshList()
		listForm.newItem.SetText("") 
	}
}

func (listForm *ListForm) updateItemFunc() {
	index, err := strconv.Atoi(listForm.updateIndex.GetText())
	if err == nil && index >= 0 && index < len(listForm.list) {
		listForm.list[index] = listForm.updateItem.GetText()
		listForm.refreshList()
		listForm.updateIndex.SetText("")
		listForm.updateItem.SetText("")
	}
}

func (listForm *ListForm) removeItem() {
	index, err := strconv.Atoi(listForm.removeIndex.GetText())
	if err == nil && index >= 0 && index < len(listForm.list) {
		listForm.list = append(listForm.list[:index], listForm.list[index+1:]...)
		listForm.refreshList()
		listForm.removeIndex.SetText("")
	}
}

func (listForm *ListForm) refreshList() {
	listForm.listView.SetText("Current List: " + fmt.Sprint(listForm.list))
}
