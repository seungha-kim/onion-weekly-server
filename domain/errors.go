package domain

type DuplicateError struct {
	fieldName string
}

func (err DuplicateError) Error() string {
	return err.fieldName + "is already taken"
}
