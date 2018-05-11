package scheduler

import "errors"

var (
	ErrEmptyAcceptedDomainList  = errors.New("Empty accepted domain list.")
	ErrSchedulerBeingInitilated = errors.New("The scheduler is being initializd.")
	ErrSchedulerBeingStarted    = errors.New("The scheduler is being started.")
	ErrSchedulerBeingStopped    = errors.New("The scheduler is being stopped.")
	ErrSchedulerBeingPaused     = errors.New("The scheduler is being paused.")
	ErrSchedulerInitialized     = errors.New("The scheduler has been initialized.")
	ErrSchedulerNotInitialized  = errors.New("The scheduler has not been initialized.")
	ErrSchedulerStarted         = errors.New("The scheduler has been started.")
	ErrSchedulerNotStarted      = errors.New("The scheduler has not been started.")
	ErrSchedulerPaused          = errors.New("The scheduler has been paused.")
	ErrSchedulerNotPaused       = errors.New("The scheduler has not been paused.")
	ErrSchedulerStopped         = errors.New("The scheduler has been stopped.")
	ErrSchedulerNotStopped      = errors.New("The scheduler has not been stopped.")
	ErrStatusUnsupported        = errors.New("Unsupported wanted status for check。")
)
