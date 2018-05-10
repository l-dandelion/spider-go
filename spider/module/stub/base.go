package stub

import (
	"github.com/l-dandelion/spider-go/spider/module"
)

/*
 * interface for internal module
 */
type ModuleInternal interface {
	module.Module
	IncrCalledCount()    // increase the called count by one
	IncrAcceptedCount()  // increase the accepted count by one
	IncrCompletedCount() // increase the completed count by one
	IncrHandlingNumber() // increase the handling number by one
	DecrHandlingNumber() // decrease the handling number by one
	Clear()              // clear all counts
}
