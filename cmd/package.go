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
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/glassechidna/stackit/pkg/stackit"
	"github.com/glassechidna/stackit/pkg/stackit/packager"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func packageTemplate(ctx context.Context, sess *session.Session, prefix string, templateReader packager.TemplateReader, writer io.Writer) (*string, error) {
	s3api := s3.New(sess)
	pkger := packager.New(s3api, sts.New(sess), *s3api.Config.Region)
	packagedTemplate, err := pkger.Package(ctx, prefix, templateReader, writer)
	if err != nil {
		return nil, errors.Wrap(err, "packaging template")
	}

	return packagedTemplate, nil

}

func writePackagedTemplateFile(absPath string, packagedTemplate *string, writer io.Writer) error {
	base := fmt.Sprintf("%s.packaged.yml", strings.TrimSuffix(filepath.Base(absPath), filepath.Ext(absPath)))
	packagedPath := filepath.Join(filepath.Dir(absPath), base)
	err := ioutil.WriteFile(packagedPath, []byte(*packagedTemplate), 0644)
	if err != nil {
		return errors.Wrap(err, "writing packaged template")
	}

	_, err = fmt.Fprintf(writer, "Wrote rendered template to %s\n", packagedPath)
	return err
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
		Short: "Upload artifacts and render template",
		Long: `
package will:

* Upload any local paths referenced in the template (complete list[1]) to S3
* Create and save the transformed template to <template>.packaged.yml

[1]: https://docs.aws.amazon.com/cli/latest/reference/cloudformation/package.html
`,
		Run: func(cmd *cobra.Command, args []string) {
			region, _ := cmd.PersistentFlags().GetString("region")
			profile, _ := cmd.PersistentFlags().GetString("profile")
			templatePath, _ := cmd.PersistentFlags().GetString("template")
			prefix, _ := cmd.PersistentFlags().GetString("prefix")

			template, err := pathToTemplate(templatePath)
			if err != nil {
				fmt.Fprintf(cmd.OutOrStderr(), "%+v\n", err)
				return
			}

			sess := awsSession(profile, region)
			packagedTemplate, err := packageTemplate(context.Background(), sess, prefix, template, cmd.OutOrStderr())
			if err != nil {
				fmt.Fprintf(cmd.OutOrStderr(), "%+v\n", err)
				return
			}

			writePackagedTemplateFile(template.Name(), packagedTemplate, cmd.OutOrStderr())
		},
	}

	cmd.PersistentFlags().String("template", "", "")
	cmd.PersistentFlags().String("prefix", "", "")
	RootCmd.AddCommand(cmd)
}

type templateReader struct {
	body string
	path string
}

func pathToTemplate(path string) (*templateReader, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, errors.Wrapf(err, "determining absolute path of '%s'", path)
	}

	if _, err := os.Stat(abs); os.IsNotExist(err) {
		return nil, errors.Errorf("no file exists at %s", abs)
	}

	body, err := ioutil.ReadFile(abs)
	if err != nil {
		return nil, err
	}

	return &templateReader{body: string(body), path: abs}, nil
}

func (t *templateReader) String() string {
	return t.body
}

func (t *templateReader) Name() string {
	return t.path
}
