package iamutil

// PolicyDocument ...
type PolicyDocument struct {
	Version   string
	Statement []Statement
}

// Statement ...
// TODO: implement more fields
type Statement struct {
	Effect   string
	Action   []string
	Resource []string
}
