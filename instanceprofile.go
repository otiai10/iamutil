package iamutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

// InstanceProfile ...
type InstanceProfile struct {
	// Path string
	Name    string
	Role    *Role
	Created *iam.InstanceProfile
}

// Create ...
func (instanceprofile *InstanceProfile) Create(sess *session.Session) error {

	client := iam.New(sess)

	out, err := client.CreateInstanceProfile(&iam.CreateInstanceProfileInput{
		// Path:             aws.String(instanceprofile.Path),
		InstanceProfileName: aws.String(instanceprofile.Name),
	})
	if err != nil {
		return err
	}
	instanceprofile.Created = out.InstanceProfile

	if instanceprofile.Role == nil {
		return nil
	}

	if instanceprofile.Role.Created == nil {
		instanceprofile.Role.Name = instanceprofile.Name
		if err := instanceprofile.Role.Create(sess); err != nil {
			if errOnDelete := instanceprofile.Delete(sess); errOnDelete != nil {
				return fmt.Errorf(
					"failed to create involved role: %v: and failed to delete instance profile: %v",
					err, errOnDelete,
				)
			}
			return err
		}
	}

	if _, err := client.AddRoleToInstanceProfile(&iam.AddRoleToInstanceProfileInput{
		InstanceProfileName: instanceprofile.Created.InstanceProfileName,
		RoleName:            instanceprofile.Role.Created.RoleName,
	}); err != nil {
		if errOnDelete := instanceprofile.Delete(sess); errOnDelete != nil {
			return fmt.Errorf(
				"failed to add role: %v: and failed to delete instance profile: %v",
				err, errOnDelete,
			)
		}
		return err
	}

	return nil
}

// Delete ...
func (instanceprofile *InstanceProfile) Delete(sess *session.Session) error {
	client := iam.New(sess)
	_, err := client.DeleteInstanceProfile(&iam.DeleteInstanceProfileInput{
		InstanceProfileName: instanceprofile.Created.InstanceProfileName,
	})
	return err
}

// TODO: Use multiple roles for an instance profile
// AddRole ...
// func (instanceprofile *InstanceProfile) AddRole(sess *session.Session, role *Role) error {
// 	client := iam.New(sess)
// 	_, err := client.AddRoleToInstanceProfile(&iam.AddRoleToInstanceProfileInput{
// 		InstanceProfileName: instanceprofile.Created.InstanceProfileName,
// 		RoleName:            role.Created.RoleName,
// 	})
//  return err
// }
