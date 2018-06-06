package iamutil

import (
	"bytes"
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

// Role ...
type Role struct {
	// Input fields
	Name        string
	Description string
	Path        string
	PolicyArns  []string

	// Result fields
	Created            *iam.Role
	AttachedPolicyArns []string
}

// Create ...
func (role *Role) Create(sess *session.Session) error {

	client := iam.New(sess)

	apd, err := role.assumeRolePolicyDocument()
	if err != nil {
		return err
	}

	if role.Path == "" {
		role.Path = "/"
	}

	out, err := client.CreateRole(&iam.CreateRoleInput{
		Path:                     aws.String(role.Path),
		RoleName:                 aws.String(role.Name),
		Description:              aws.String(role.Description),
		AssumeRolePolicyDocument: aws.String(apd),
	})
	if err != nil {
		return err
	}

	role.Created = out.Role

	if len(role.PolicyArns) != 0 {
		if err := role.AttachPolicy(sess, role.PolicyArns...); err != nil {
			return err
		}
	}

	return nil
}

func (role *Role) assumeRolePolicyDocument() (string, error) {
	buf := bytes.NewBuffer(nil)
	err := json.NewEncoder(buf).Encode(map[string]interface{}{
		"Version": "2012-10-17",
		"Statement": []map[string]interface{}{
			{
				"Sid":       "",
				"Effect":    "Allow",
				"Action":    "sts:AssumeRole",
				"Principal": map[string]string{"Service": "ec2.amazonaws.com"},
			},
		},
	})
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// AttachPolicy ...
func (role *Role) AttachPolicy(sess *session.Session, arns ...string) error {

	client := iam.New(sess)

	for _, arn := range arns {
		_, err := client.AttachRolePolicy(&iam.AttachRolePolicyInput{
			PolicyArn: aws.String(arn),
			RoleName:  role.Created.RoleName,
		})
		if err != nil {
			return err
		}
		role.AttachedPolicyArns = append(role.AttachedPolicyArns, arn)
	}

	return nil
}
