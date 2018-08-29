package stackit

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"math/rand"
	"os"
	"strings"
	"time"
)

func generateToken() string {
	rand.Seed(time.Now().UTC().UnixNano())
	return fmt.Sprintf("stackit-%d", rand.Uint32())
}

func mapToTags(tagMap map[string]string) []*cloudformation.Tag {
	tags := []*cloudformation.Tag{}

	prefix := "CFN_TAG_"
	for _, envvar := range os.Environ() {
		if strings.HasPrefix(envvar, prefix) {
			sansPrefix := envvar[len(prefix):]
			keyval := strings.SplitN(sansPrefix, "=", 2)
			tags = append(tags, &cloudformation.Tag{Key: &keyval[0], Value: &keyval[1]})

		}
	}

	for key, val := range tagMap {
		tags = append(tags, &cloudformation.Tag{Key: aws.String(key), Value: aws.String(val)})
	}

	return tags
}

func stringInSlice(slice []string, s string) bool {
	for _, ss := range slice {
		if s == ss {
			return true
		}
	}
	return false
}
