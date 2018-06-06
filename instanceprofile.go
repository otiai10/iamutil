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

	out, err := client.GetInstanceProfile(&iam.GetInstanceProfileInput{
		InstanceProfileName: aws.String(instanceprofile.Name),
	})
	if err != nil {
		return err
	}
	instanceprofile.Created = out.InstanceProfile

	for _, awsrole := range instanceprofile.Created.Roles {
		role := &Role{Name: *awsrole.RoleName, Created: awsrole}
		_, err := client.RemoveRoleFromInstanceProfile(&iam.RemoveRoleFromInstanceProfileInput{
			InstanceProfileName: instanceprofile.Created.InstanceProfileName,
			RoleName:            role.Created.RoleName,
		})
		if err != nil {
			return err
		}
		if err := role.Delete(sess); err != nil {
			return err
		}
	}

	_, err = client.DeleteInstanceProfile(&iam.DeleteInstanceProfileInput{
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

// FindInstanceProfile ...
func FindInstanceProfile(sess *session.Session, name string) (*InstanceProfile, error) {
	client := iam.New(sess)
	out, err := client.GetInstanceProfile(&iam.GetInstanceProfileInput{
		InstanceProfileName: aws.String(name),
	})
	if err != nil {
		return nil, err
	}
	return &InstanceProfile{
		Name:    *out.InstanceProfile.InstanceProfileName,
		Role:    nil, // TODO: populate roles
		Created: out.InstanceProfile,
	}, nil
}
