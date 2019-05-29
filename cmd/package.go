// Copyright © 2017 Aidan Steele <aidan.steele@glassechidna.com.au>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/glassechidna/stackit/pkg/stackit"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func packageTemplate(ctx context.Context, region, profile, stackName, templatePath string, tags, parameters map[string]string, writer io.Writer) error {
	absPath, err := filepath.Abs(templatePath)
	if err != nil {
		return errors.Wrapf(err, "determining absolute path of '%s'", templatePath)
	}
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return errors.Errorf("no file exists at %s", absPath)
	}

	sess := awsSession(profile, region)
	sit := stackit.NewStackit(cloudformation.New(sess), sts.New(sess))

	s3api := s3.New(sess)
	if region == "" && s3api.Config.Region != nil {
		region = *s3api.Config.Region
	}

	packager := stackit.NewPackager(s3api, sts.New(sess), region)

	events := make(chan stackit.TailStackEvent)
	printer := stackit.NewTailPrinterWithOptions(true, true, writer)

	printerCtx, printerCancel := context.WithCancel(ctx)
	defer printerCancel()

	go func() {
		for {
			select {
			case <-printerCtx.Done():
				return
			case tailEvent := <-events:
				printer.PrintTailEvent(tailEvent)
			}
		}
	}()

	upInput, err := packager.Package(stackName, absPath, tags, parameters)
	if err != nil {
		return errors.Wrap(err, "packaging template")
	}

	prepared, err := sit.Prepare(ctx, *upInput, events)
	if err != nil {
		panic(err)
	}

	if prepared == nil {
		return nil // no-op change set
	}

	io.WriteString(writer, userFriendlyChangesOutput(prepared))

	err = savePreparedOutput(prepared)
	return errors.Wrap(err, "saving packaged output")
}

func userFriendlyChangesOutput(output *stackit.PrepareOutput) string {
	sbuf := &strings.Builder{}
	tbl := tablewriter.NewWriter(sbuf)
	tbl.SetHeader([]string{"Action", "Resource", "Type"})

	for _, change := range output.Changes {
		tbl.Append([]string{
			*change.ResourceChange.Action,
			*change.ResourceChange.LogicalResourceId,
			*change.ResourceChange.ResourceType,
		})
	}

	tbl.Render()

	buf := &bytes.Buffer{}

	tmpl, err := template.New("").Parse(`
Stack ID: {{ .StackId }}
Change Set ID: {{ .ChangeSetId }}
Changes:

{{ .Changes }}`)
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(buf, map[string]interface{}{
		"StackId":     *output.Output.StackId,
		"ChangeSetId": *output.Output.Id,
		"Changes":     sbuf.String(),
	})

	return buf.String()
}

func init() {
	cmd := &cobra.Command{
		Use:   "package",
		Short: "Package template and create change set",
		Long: `
package will:

* Upload any local paths referenced in the template (complete list[1]) to S3
* Create a changeset from the transformed template
* Print out a list of the estimated changes that executing the template will cause
* Save to disk the change set ID for later execution of 'stackit execute'

[1]: https://docs.aws.amazon.com/cli/latest/reference/cloudformation/package.html
`,
		Run: func(cmd *cobra.Command, args []string) {
			region, _ := cmd.PersistentFlags().GetString("region")
			profile, _ := cmd.PersistentFlags().GetString("profile")
			stackName, _ := cmd.PersistentFlags().GetString("stack-name")
			templatePath, _ := cmd.PersistentFlags().GetString("template")
			tagKvps, _ := cmd.PersistentFlags().GetStringSlice("tags")
			tags := keyvalSliceToMap(tagKvps)
			params := keyvalSliceToMap(args)

			err := packageTemplate(context.Background(), region, profile, stackName, templatePath, tags, params, cmd.OutOrStderr())
			if err != nil {
				fmt.Fprintf(cmd.OutOrStderr(), "%+v\n", err)
			}
		},
	}

	cmd.PersistentFlags().String("stack-name", "", "")
	cmd.PersistentFlags().String("template", "", "")
	cmd.PersistentFlags().StringSlice("tags", []string{}, "")
	RootCmd.AddCommand(cmd)
}
