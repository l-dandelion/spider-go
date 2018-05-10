package stub

import (
	"github.com/l-dandelion/spider-go/spider/module"
	"sync/atomic"
)

/*
 * implementation of interface ModuleInternal
 */
type myModule struct {
	calledCount     uint64                // called count
	acceptedCount   uint64                // accepted count
	completedCount  uint64                // completed count
	handlingNumber  uint64                // handling number
}

/*
 * create an instance of ModuleInternal
 */
func NewModuleInternal() (mi ModuleInternal) {
	return &myModule{}
}

/*
 * get the called count of module (concurrent security)
 */
func (m *myModule) CalledCount() uint64 {
	return atomic.LoadUint64(&m.calledCount)
}

/*
 * get the accepted count of module (concurrent security)
 */
func (m *myModule) AcceptedCount() uint64 {
	return atomic.LoadUint64(&m.acceptedCount)
}

/*
 * get the completed count of module (concurrent security)
 */
func (m *myModule) CompletedCount() uint64 {
	return atomic.LoadUint64(&m.completedCount)
}

/*
 * get the handling number of module (concurrent security)
 */
func (m *myModule) HandlingNumber() uint64 {
	return atomic.LoadUint64(&m.handlingNumber)
}

/*
 * get the counts of module (concurrent security)
 */
func (m *myModule) Counts() module.Counts {
	return module.Counts{
		CalledCount:    m.CalledCount(),
		AcceptedCount:  m.AcceptedCount(),
		CompletedCount: m.CompletedCount(),
		HandlingNumber: m.HandlingNumber(),
	}
}

/*
 * get the summary of module (concurrent security)
 */
func (m *myModule) Summary() module.SummaryStruct {
	counts := m.Counts()
	return module.SummaryStruct{
		Called:    counts.CalledCount,
		Accepted:  counts.AcceptedCount,
		Completed: counts.CompletedCount,
		Handling:  counts.HandlingNumber,
		Extra:     nil,
	}
}

/*
 * increase the called count by one (concurrent security)
 */
func (m *myModule) IncrCalledCount() {
	atomic.AddUint64(&m.calledCount, 1)
}

/*
 * increase the accepted count by one (concurrent security)
 */
func (m *myModule) IncrAcceptedCount() {
	atomic.AddUint64(&m.acceptedCount, 1)
}

/*
 * increase the completed count by one (concurrent security)
 */
func (m *myModule) IncrCompletedCount() {
	atomic.AddUint64(&m.completedCount, 1)
}

/*
 * increase the handling number by one (concurrent security)
 */
func (m *myModule) IncrHandlingNumber() {
	atomic.AddUint64(&m.handlingNumber, 1)
}

/*
 * increase the handling number by one (concurrent security)
 */
func (m *myModule) DecrHandlingNumber() {
	atomic.AddUint64(&m.handlingNumber, ^uint64(0))
}

/*
 * clear all counts (concurrent security)
 */
func (m *myModule) Clear() {
	atomic.StoreUint64(&m.calledCount, 0)
	atomic.StoreUint64(&m.acceptedCount, 0)
	atomic.StoreUint64(&m.completedCount, 0)
	atomic.StoreUint64(&m.handlingNumber, 0)
}
