package task

// CancelCondition is a condition for canceling the task group.
type CancelCondition func(err error) (cancel bool)

// CancelOnError indicates that task group should return on non-nil error. This
// is the same behavior as the golang.org/x/sync/errgroup package.
func CancelOnError() CancelCondition {
	return cancelOnError
}

func cancelOnError(err error) (cancel bool) {
	return err != nil
}

// CancelOnReturn indicates that task group should be canceled if any
// subtask returns.
func CancelOnReturn() CancelCondition {
	return cancelOnReturn
}

func cancelOnReturn(err error) (cancel bool) {
	return true
}

// NeverCancel indicates that task group should not be canceled until
// all subtasks complete.
func NeverCancel() CancelCondition {
	return neverCancel
}

func neverCancel(err error) (cancel bool) {
	return false
}
