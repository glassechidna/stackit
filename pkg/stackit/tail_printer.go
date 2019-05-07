package stackit

import (
	"fmt"
	"github.com/fatih/color"
	"io"
	"os"
	"strings"
)

type TailPrinter struct {
	timestampFormat string
	successColor    *color.Color
	failureColor    *color.Color
	writer          io.Writer
}

func NewTailPrinter() TailPrinter {
	return NewTailPrinterWithOptions(true, true, os.Stderr)
}

func NewTailPrinterWithOptions(showTimestamp, showColors bool, writer io.Writer) TailPrinter {
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
		successColor:    successColor,
		failureColor:    failureColor,
		writer:          writer,
	}
}

func (tp *TailPrinter) FormatTailEvent(tailEvent TailStackEvent) string {
	resourceNameLength := 20 // TODO: determine this from template/API

	timestampPrefix := tailEvent.Timestamp.Format(tp.timestampFormat)

	reasonPart := ""
	if tailEvent.ResourceStatusReason != nil {
		reasonPart = fmt.Sprintf("- %s", *tailEvent.ResourceStatusReason)
	}

	line := fmt.Sprintf("%s %s - %s %s", timestampPrefix, fixedLengthString(resourceNameLength, *tailEvent.LogicalResourceId), *tailEvent.ResourceStatus, reasonPart)

	if isBadStatus(*tailEvent.ResourceStatus) && tp.failureColor != nil {
		return tp.failureColor.Sprint(line)
	} else {
		return line
	}
}

func (tp *TailPrinter) PrintTailEvent(tailEvent TailStackEvent) {
	line := tp.FormatTailEvent(tailEvent)
	tp.writer.Write([]byte(line))
	tp.writer.Write([]byte("\n"))
}

func fixedLengthString(length int, str string) string {
	verb := fmt.Sprintf("%%%d.%ds", length, length)
	return fmt.Sprintf(verb, str)
}

func isBadStatus(status string) bool {
	return strings.HasSuffix(status, "_FAILED")
}
