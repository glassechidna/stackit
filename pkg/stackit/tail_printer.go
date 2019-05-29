package stackit

import (
	"fmt"
	"github.com/fatih/color"
	"io"
	"strings"
)

type TailPrinter struct {
	timestampFormat string
	failureColor    *color.Color
	writer          io.Writer
}

func NewTailPrinter(writer io.Writer) TailPrinter {
	return TailPrinter{
		timestampFormat: "[03:04:05]",
		failureColor:    color.New(color.FgRed),
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
	fmt.Fprintln(tp.writer, line)
}

func fixedLengthString(length int, str string) string {
	verb := fmt.Sprintf("%%%d.%ds", length, length)
	return fmt.Sprintf(verb, str)
}

func isBadStatus(status string) bool {
	return strings.HasSuffix(status, "_FAILED")
}
