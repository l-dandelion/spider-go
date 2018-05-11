package scheduler

import (
	"github.com/l-dandelion/yi-ants-go/lib/constant"
)

/*
 * check whether status can be changed to wantedStatus from currentStatus
 */
func checkStatus(currentStatus, wantedStatus int8) (err error) {
	switch currentStatus {
	case constant.RUNNING_STATUS_PREPARING:
		err = ErrSchedulerBeingInitilated
	case constant.RUNNING_STATUS_STARTING:
		err = ErrSchedulerBeingStarted
	case constant.RUNNING_STATUS_STOPPING:
		err = ErrSchedulerBeingStopped
	case constant.RUNNING_STATUS_PAUSING:
		err = ErrSchedulerBeingPaused
	}
	if err != nil {
		return
	}

	if currentStatus == constant.RUNNING_STATUS_UNPREPARED &&
		wantedStatus != constant.RUNNING_STATUS_PREPARING {
		err = ErrSchedulerNotInitialized
		return
	}

	if currentStatus == constant.RUNNING_STATUS_STOPPED {
		err = ErrSchedulerStopped
		return
	}

	switch wantedStatus {
	case constant.RUNNING_STATUS_PREPARING:
		if currentStatus != constant.RUNNING_STATUS_UNPREPARED {
			err = ErrSchedulerInitialized
		}
	case constant.RUNNING_STATUS_STARTING:
		if currentStatus != constant.RUNNING_STATUS_PREPARED {
			err = ErrSchedulerStarted
		}
	case constant.RUNNING_STATUS_PAUSING:
		if currentStatus != constant.RUNNING_STATUS_STARTED {
			err = ErrSchedulerNotStarted
		}
	case constant.RUNNING_STATUS_STOPPING:
		if currentStatus != constant.RUNNING_STATUS_STARTED &&
			currentStatus != constant.RUNNING_STATUS_PAUSED {
			err = ErrSchedulerNotStarted
		}
	default:
		err = ErrStatusUnsupported
	}
	return
}

/*
 * get status description
 */
func GetStatusDescription(status int8) string {
	desc, ok := constant.GlobalArrRunningStatusDesc[status]
	if !ok {
		return "Unknow"
	}
	return desc
}
