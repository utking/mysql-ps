package ui

import (
	"fmt"
	"testing"

	"github.com/rivo/tview"
)

func TestCheckTextType(t *testing.T) {
	tv := tview.NewTextView()
	val := tv.GetText(true)
	fmt.Printf("Type of GetText(): %T\n", val)
}
