package util

import (
	"github.com/jedib0t/go-pretty/progress"
	"github.com/jedib0t/go-pretty/table"
)

// NewTableWriter returns a table writer
func NewTableWriter() table.Writer {
	tw := table.NewWriter()
	tw.Style().Options.DrawBorder = false
	tw.Style().Options.SeparateColumns = false
	tw.Style().Options.SeparateHeader = false
	tw.Style().Options.SeparateRows = false
	return tw
}

// NewProgressWriter return a progress writer
func NewProgressWriter() progress.Writer {
	pw := progress.NewWriter()
	pw.SetAutoStop(false)
	pw.SetTrackerLength(25)
	pw.ShowOverallTracker(true)
	pw.ShowTime(true)
	pw.ShowTracker(true)
	pw.ShowValue(true)
	pw.SetMessageWidth(24)
	pw.SetSortBy(progress.SortByPercentDsc)
	pw.SetStyle(progress.StyleDefault)
	pw.SetTrackerPosition(progress.PositionRight)
	pw.Style().Colors = progress.StyleColorsExample
	pw.Style().Options.PercentFormat = "%4.1f%%"
	return pw
}
