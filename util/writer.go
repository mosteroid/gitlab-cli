package util

import (
	"fmt"

	"github.com/jedib0t/go-pretty/text"

	"github.com/jedib0t/go-pretty/progress"
	"github.com/jedib0t/go-pretty/table"
)

// UnitTime defines the value as a time unit.
var UnitTime = progress.Units{
	Notation:  "",
	Formatter: FormatTime,
}

// FormatTime formats the given value as a "Time".
func FormatTime(value int64) string {
	return fmt.Sprintf("%ds", value)
}

// StatusWriter gitlab status writer
type StatusWriter struct {
	RunningColor            text.Color
	RunningBackgroundColor  text.Color
	PendingColor            text.Color
	PendingBackgroundColor  text.Color
	SuccessColor            text.Color
	SuccessBackgroundColor  text.Color
	FailedColor             text.Color
	FailedBackgroundColor   text.Color
	CanceledColor           text.Color
	CanceledBackgroundColor text.Color
	SkippedColor            text.Color
	SkippedBackgroundColor  text.Color
}

// NewTableWriter returns a new table writer
func NewTableWriter() table.Writer {
	tw := table.NewWriter()
	tw.Style().Options.DrawBorder = false
	tw.Style().Options.SeparateColumns = false
	tw.Style().Options.SeparateHeader = false
	tw.Style().Options.SeparateRows = false
	return tw
}

// NewProgressWriter returns a new progress writer
func NewProgressWriter() progress.Writer {
	pw := progress.NewWriter()
	pw.SetAutoStop(false)
	pw.SetTrackerLength(25)
	pw.ShowOverallTracker(true)
	pw.ShowTime(false)
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

//NewStatusWriter returns a new status writer
func NewStatusWriter() *StatusWriter {
	sw := &StatusWriter{
		RunningColor:            text.FgBlue,
		RunningBackgroundColor:  text.BgBlack,
		PendingColor:            text.FgYellow,
		PendingBackgroundColor:  text.BgBlack,
		SuccessColor:            text.FgGreen,
		SuccessBackgroundColor:  text.BgBlack,
		FailedColor:             text.FgRed,
		FailedBackgroundColor:   text.BgBlack,
		CanceledColor:           text.FgBlack,
		CanceledBackgroundColor: text.BgWhite,
		SkippedColor:            text.BgHiWhite,
		SkippedBackgroundColor:  text.BgBlack,
	}
	return sw
}

// Sprintf formats and colorizes the given status
func (writer *StatusWriter) Sprintf(status string) string {

	coloredStatus := status
	switch status {
	case "running":
		coloredStatus = writer.RunningColor.Sprint(coloredStatus)
		coloredStatus = writer.RunningBackgroundColor.Sprint(coloredStatus)
	case "pending":
		coloredStatus = writer.PendingColor.Sprint(coloredStatus)
		coloredStatus = writer.PendingBackgroundColor.Sprint(coloredStatus)
	case "success":
		coloredStatus = writer.SuccessColor.Sprint(coloredStatus)
		coloredStatus = writer.SuccessBackgroundColor.Sprint(coloredStatus)
	case "failed":
		coloredStatus = writer.FailedColor.Sprint(coloredStatus)
		coloredStatus = writer.FailedBackgroundColor.Sprint(coloredStatus)
	case "canceled":
		coloredStatus = writer.CanceledColor.Sprint(coloredStatus)
		coloredStatus = writer.CanceledBackgroundColor.Sprint(coloredStatus)
	case "skipped":
		coloredStatus = writer.SkippedColor.Sprint(coloredStatus)
		coloredStatus = writer.SkippedBackgroundColor.Sprint(coloredStatus)
	}
	return coloredStatus
}
