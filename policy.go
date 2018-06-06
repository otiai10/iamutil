package iamutil

import (
	"bytes"
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

// Policy ...
type Policy struct {
	// Input fields
	Name        string
	Description string
	// Path     string
	Document *PolicyDocument

	// Result fields
	Created *iam.Policy
}

// Create ...
func (policy *Policy) Create(sess *session.Session) error {

	client := iam.New(sess)

	doc, err := policy.document()
	if err != nil {
		return err
	}

	out, err := client.CreatePolicy(&iam.CreatePolicyInput{
		PolicyName:  aws.String(policy.Name),
		Description: aws.String(policy.Description),
		// Path:     aws.String(policy.Path),
		PolicyDocument: aws.String(doc),
	})
	if err != nil {
		return err
	}

	policy.Created = out.Policy

	return nil
}

// document ...
func (policy *Policy) document() (string, error) {
	buf := bytes.NewBuffer(nil)
	if err := json.NewEncoder(buf).Encode(policy.Document); err != nil {
		return "", nil
	}
	return buf.String(), nil
}
