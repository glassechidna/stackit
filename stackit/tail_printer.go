package stackit

import (
	"github.com/fatih/color"
	"io"
	"fmt"
	"os"
)

type TailPrinter struct {
	timestampFormat string
	successColor *color.Color
	failureColor *color.Color
	writer io.Writer
}

func NewTailPrinter() TailPrinter {
	return NewTailPrinterWithOptions(true, true)
}

func NewTailPrinterWithOptions(showTimestamp, showColors bool) TailPrinter {
	format := ""
	if showTimestamp {
		format = "[03:04:05]"
	}

	successColor := color.New(color.FgGreen)
	failureColor := color.New(color.FgRed)

	if !showColors {
		successColor = nil
		failureColor = nil
	}

	return TailPrinter{
		timestampFormat: format,
		successColor: successColor,
		failureColor: failureColor,
		writer: os.Stderr,
	}
}

func (tp *TailPrinter) PrintTailEvent(tailEvent TailStackEvent) {
	resourceNameLength := 20 // TODO: determine this from template/API

	event := tailEvent.Event

	timestampPrefix := event.Timestamp.Format(tp.timestampFormat)

	reasonPart := ""
	if event.ResourceStatusReason != nil {
		reasonPart = fmt.Sprintf("- %s", *event.ResourceStatusReason)
	}

	line := fmt.Sprintf("%s %s - %s %s", timestampPrefix, fixedLengthString(resourceNameLength, *event.LogicalResourceId), *event.ResourceStatus, reasonPart)

	if isBadStatus(*event.ResourceStatus) && tp.failureColor != nil {
		tp.failureColor.Fprintln(tp.writer, line)
	} else {
		fmt.Fprintln(tp.writer, line)
	}
}