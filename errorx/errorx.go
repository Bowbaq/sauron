package errorx

// Error encodes a constant error string
type Error string

// Error implements the error interface
func (e Error) Error() string {
	return string(e)
}

const (
	// ErrRepositoryOwnerRequired is returned when the repository owner is missing
	ErrRepositoryOwnerRequired = Error("Repository owner cannot be empty")
	// ErrRepositoryNameRequired is returned when the repository name is missing
	ErrRepositoryNameRequired = Error("Repository name cannot be empty")

	// ErrBucketNameRequired is returned when the bucket name is missing
	ErrBucketNameRequired = Error("Bucket name cannot be empty")
	// ErrBucketKeyRequired is returned when the bucket key is missing
	ErrBucketKeyRequired = Error("Key cannot be empty")
	// ErrSNSTopicRequired is returned when the SNS topic is missing
	ErrSNSTopicRequired = Error("SNS topic cannot be empty")
)
