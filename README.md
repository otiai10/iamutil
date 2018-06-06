# iamutil

Go utilities to control AWS IAM resources.

```go
package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/otiai10/iamutil"
)

func main() {

	// Specify configurations
	profile := &iamutil.InstanceProfile{
		Role: &iamutil.Role{
			Description: "Test Role by iamutil",
			PolicyArns: []string{
				"arn:aws:iam::aws:policy/AmazonS3FullAccess",
			},
		},
		Name: "otiai10-test",
	}

	// Setup API
	sess := session.New(&aws.Config{
		Region: aws.String("ap-northeast-1"),
	})

	// Execute
	err := profile.Create(sess)

	fmt.Printf("%+v\nERROR: %v\n", profile, err)

}
```