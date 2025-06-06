// Code generated by http://github.com/gojuno/minimock (v3.4.5). DO NOT EDIT.

package product

//go:generate minimock -i route256/cart/internal/services/product.Limiter -o limiter_mock_test.go -n LimiterMock -p product

import (
	"context"
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock/v3"
)

// LimiterMock implements Limiter
type LimiterMock struct {
	t          minimock.Tester
	finishOnce sync.Once

	funcWait          func(ctx context.Context) (err error)
	funcWaitOrigin    string
	inspectFuncWait   func(ctx context.Context)
	afterWaitCounter  uint64
	beforeWaitCounter uint64
	WaitMock          mLimiterMockWait
}

// NewLimiterMock returns a mock for Limiter
func NewLimiterMock(t minimock.Tester) *LimiterMock {
	m := &LimiterMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.WaitMock = mLimiterMockWait{mock: m}
	m.WaitMock.callArgs = []*LimiterMockWaitParams{}

	t.Cleanup(m.MinimockFinish)

	return m
}

type mLimiterMockWait struct {
	optional           bool
	mock               *LimiterMock
	defaultExpectation *LimiterMockWaitExpectation
	expectations       []*LimiterMockWaitExpectation

	callArgs []*LimiterMockWaitParams
	mutex    sync.RWMutex

	expectedInvocations       uint64
	expectedInvocationsOrigin string
}

// LimiterMockWaitExpectation specifies expectation struct of the Limiter.Wait
type LimiterMockWaitExpectation struct {
	mock               *LimiterMock
	params             *LimiterMockWaitParams
	paramPtrs          *LimiterMockWaitParamPtrs
	expectationOrigins LimiterMockWaitExpectationOrigins
	results            *LimiterMockWaitResults
	returnOrigin       string
	Counter            uint64
}

// LimiterMockWaitParams contains parameters of the Limiter.Wait
type LimiterMockWaitParams struct {
	ctx context.Context
}

// LimiterMockWaitParamPtrs contains pointers to parameters of the Limiter.Wait
type LimiterMockWaitParamPtrs struct {
	ctx *context.Context
}

// LimiterMockWaitResults contains results of the Limiter.Wait
type LimiterMockWaitResults struct {
	err error
}

// LimiterMockWaitOrigins contains origins of expectations of the Limiter.Wait
type LimiterMockWaitExpectationOrigins struct {
	origin    string
	originCtx string
}

// Marks this method to be optional. The default behavior of any method with Return() is '1 or more', meaning
// the test will fail minimock's automatic final call check if the mocked method was not called at least once.
// Optional() makes method check to work in '0 or more' mode.
// It is NOT RECOMMENDED to use this option unless you really need it, as default behaviour helps to
// catch the problems when the expected method call is totally skipped during test run.
func (mmWait *mLimiterMockWait) Optional() *mLimiterMockWait {
	mmWait.optional = true
	return mmWait
}

// Expect sets up expected params for Limiter.Wait
func (mmWait *mLimiterMockWait) Expect(ctx context.Context) *mLimiterMockWait {
	if mmWait.mock.funcWait != nil {
		mmWait.mock.t.Fatalf("LimiterMock.Wait mock is already set by Set")
	}

	if mmWait.defaultExpectation == nil {
		mmWait.defaultExpectation = &LimiterMockWaitExpectation{}
	}

	if mmWait.defaultExpectation.paramPtrs != nil {
		mmWait.mock.t.Fatalf("LimiterMock.Wait mock is already set by ExpectParams functions")
	}

	mmWait.defaultExpectation.params = &LimiterMockWaitParams{ctx}
	mmWait.defaultExpectation.expectationOrigins.origin = minimock.CallerInfo(1)
	for _, e := range mmWait.expectations {
		if minimock.Equal(e.params, mmWait.defaultExpectation.params) {
			mmWait.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmWait.defaultExpectation.params)
		}
	}

	return mmWait
}

// ExpectCtxParam1 sets up expected param ctx for Limiter.Wait
func (mmWait *mLimiterMockWait) ExpectCtxParam1(ctx context.Context) *mLimiterMockWait {
	if mmWait.mock.funcWait != nil {
		mmWait.mock.t.Fatalf("LimiterMock.Wait mock is already set by Set")
	}

	if mmWait.defaultExpectation == nil {
		mmWait.defaultExpectation = &LimiterMockWaitExpectation{}
	}

	if mmWait.defaultExpectation.params != nil {
		mmWait.mock.t.Fatalf("LimiterMock.Wait mock is already set by Expect")
	}

	if mmWait.defaultExpectation.paramPtrs == nil {
		mmWait.defaultExpectation.paramPtrs = &LimiterMockWaitParamPtrs{}
	}
	mmWait.defaultExpectation.paramPtrs.ctx = &ctx
	mmWait.defaultExpectation.expectationOrigins.originCtx = minimock.CallerInfo(1)

	return mmWait
}

// Inspect accepts an inspector function that has same arguments as the Limiter.Wait
func (mmWait *mLimiterMockWait) Inspect(f func(ctx context.Context)) *mLimiterMockWait {
	if mmWait.mock.inspectFuncWait != nil {
		mmWait.mock.t.Fatalf("Inspect function is already set for LimiterMock.Wait")
	}

	mmWait.mock.inspectFuncWait = f

	return mmWait
}

// Return sets up results that will be returned by Limiter.Wait
func (mmWait *mLimiterMockWait) Return(err error) *LimiterMock {
	if mmWait.mock.funcWait != nil {
		mmWait.mock.t.Fatalf("LimiterMock.Wait mock is already set by Set")
	}

	if mmWait.defaultExpectation == nil {
		mmWait.defaultExpectation = &LimiterMockWaitExpectation{mock: mmWait.mock}
	}
	mmWait.defaultExpectation.results = &LimiterMockWaitResults{err}
	mmWait.defaultExpectation.returnOrigin = minimock.CallerInfo(1)
	return mmWait.mock
}

// Set uses given function f to mock the Limiter.Wait method
func (mmWait *mLimiterMockWait) Set(f func(ctx context.Context) (err error)) *LimiterMock {
	if mmWait.defaultExpectation != nil {
		mmWait.mock.t.Fatalf("Default expectation is already set for the Limiter.Wait method")
	}

	if len(mmWait.expectations) > 0 {
		mmWait.mock.t.Fatalf("Some expectations are already set for the Limiter.Wait method")
	}

	mmWait.mock.funcWait = f
	mmWait.mock.funcWaitOrigin = minimock.CallerInfo(1)
	return mmWait.mock
}

// When sets expectation for the Limiter.Wait which will trigger the result defined by the following
// Then helper
func (mmWait *mLimiterMockWait) When(ctx context.Context) *LimiterMockWaitExpectation {
	if mmWait.mock.funcWait != nil {
		mmWait.mock.t.Fatalf("LimiterMock.Wait mock is already set by Set")
	}

	expectation := &LimiterMockWaitExpectation{
		mock:               mmWait.mock,
		params:             &LimiterMockWaitParams{ctx},
		expectationOrigins: LimiterMockWaitExpectationOrigins{origin: minimock.CallerInfo(1)},
	}
	mmWait.expectations = append(mmWait.expectations, expectation)
	return expectation
}

// Then sets up Limiter.Wait return parameters for the expectation previously defined by the When method
func (e *LimiterMockWaitExpectation) Then(err error) *LimiterMock {
	e.results = &LimiterMockWaitResults{err}
	return e.mock
}

// Times sets number of times Limiter.Wait should be invoked
func (mmWait *mLimiterMockWait) Times(n uint64) *mLimiterMockWait {
	if n == 0 {
		mmWait.mock.t.Fatalf("Times of LimiterMock.Wait mock can not be zero")
	}
	mm_atomic.StoreUint64(&mmWait.expectedInvocations, n)
	mmWait.expectedInvocationsOrigin = minimock.CallerInfo(1)
	return mmWait
}

func (mmWait *mLimiterMockWait) invocationsDone() bool {
	if len(mmWait.expectations) == 0 && mmWait.defaultExpectation == nil && mmWait.mock.funcWait == nil {
		return true
	}

	totalInvocations := mm_atomic.LoadUint64(&mmWait.mock.afterWaitCounter)
	expectedInvocations := mm_atomic.LoadUint64(&mmWait.expectedInvocations)

	return totalInvocations > 0 && (expectedInvocations == 0 || expectedInvocations == totalInvocations)
}

// Wait implements Limiter
func (mmWait *LimiterMock) Wait(ctx context.Context) (err error) {
	mm_atomic.AddUint64(&mmWait.beforeWaitCounter, 1)
	defer mm_atomic.AddUint64(&mmWait.afterWaitCounter, 1)

	mmWait.t.Helper()

	if mmWait.inspectFuncWait != nil {
		mmWait.inspectFuncWait(ctx)
	}

	mm_params := LimiterMockWaitParams{ctx}

	// Record call args
	mmWait.WaitMock.mutex.Lock()
	mmWait.WaitMock.callArgs = append(mmWait.WaitMock.callArgs, &mm_params)
	mmWait.WaitMock.mutex.Unlock()

	for _, e := range mmWait.WaitMock.expectations {
		if minimock.Equal(*e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.err
		}
	}

	if mmWait.WaitMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmWait.WaitMock.defaultExpectation.Counter, 1)
		mm_want := mmWait.WaitMock.defaultExpectation.params
		mm_want_ptrs := mmWait.WaitMock.defaultExpectation.paramPtrs

		mm_got := LimiterMockWaitParams{ctx}

		if mm_want_ptrs != nil {

			if mm_want_ptrs.ctx != nil && !minimock.Equal(*mm_want_ptrs.ctx, mm_got.ctx) {
				mmWait.t.Errorf("LimiterMock.Wait got unexpected parameter ctx, expected at\n%s:\nwant: %#v\n got: %#v%s\n",
					mmWait.WaitMock.defaultExpectation.expectationOrigins.originCtx, *mm_want_ptrs.ctx, mm_got.ctx, minimock.Diff(*mm_want_ptrs.ctx, mm_got.ctx))
			}

		} else if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmWait.t.Errorf("LimiterMock.Wait got unexpected parameters, expected at\n%s:\nwant: %#v\n got: %#v%s\n",
				mmWait.WaitMock.defaultExpectation.expectationOrigins.origin, *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmWait.WaitMock.defaultExpectation.results
		if mm_results == nil {
			mmWait.t.Fatal("No results are set for the LimiterMock.Wait")
		}
		return (*mm_results).err
	}
	if mmWait.funcWait != nil {
		return mmWait.funcWait(ctx)
	}
	mmWait.t.Fatalf("Unexpected call to LimiterMock.Wait. %v", ctx)
	return
}

// WaitAfterCounter returns a count of finished LimiterMock.Wait invocations
func (mmWait *LimiterMock) WaitAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmWait.afterWaitCounter)
}

// WaitBeforeCounter returns a count of LimiterMock.Wait invocations
func (mmWait *LimiterMock) WaitBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmWait.beforeWaitCounter)
}

// Calls returns a list of arguments used in each call to LimiterMock.Wait.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmWait *mLimiterMockWait) Calls() []*LimiterMockWaitParams {
	mmWait.mutex.RLock()

	argCopy := make([]*LimiterMockWaitParams, len(mmWait.callArgs))
	copy(argCopy, mmWait.callArgs)

	mmWait.mutex.RUnlock()

	return argCopy
}

// MinimockWaitDone returns true if the count of the Wait invocations corresponds
// the number of defined expectations
func (m *LimiterMock) MinimockWaitDone() bool {
	if m.WaitMock.optional {
		// Optional methods provide '0 or more' call count restriction.
		return true
	}

	for _, e := range m.WaitMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	return m.WaitMock.invocationsDone()
}

// MinimockWaitInspect logs each unmet expectation
func (m *LimiterMock) MinimockWaitInspect() {
	for _, e := range m.WaitMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to LimiterMock.Wait at\n%s with params: %#v", e.expectationOrigins.origin, *e.params)
		}
	}

	afterWaitCounter := mm_atomic.LoadUint64(&m.afterWaitCounter)
	// if default expectation was set then invocations count should be greater than zero
	if m.WaitMock.defaultExpectation != nil && afterWaitCounter < 1 {
		if m.WaitMock.defaultExpectation.params == nil {
			m.t.Errorf("Expected call to LimiterMock.Wait at\n%s", m.WaitMock.defaultExpectation.returnOrigin)
		} else {
			m.t.Errorf("Expected call to LimiterMock.Wait at\n%s with params: %#v", m.WaitMock.defaultExpectation.expectationOrigins.origin, *m.WaitMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcWait != nil && afterWaitCounter < 1 {
		m.t.Errorf("Expected call to LimiterMock.Wait at\n%s", m.funcWaitOrigin)
	}

	if !m.WaitMock.invocationsDone() && afterWaitCounter > 0 {
		m.t.Errorf("Expected %d calls to LimiterMock.Wait at\n%s but found %d calls",
			mm_atomic.LoadUint64(&m.WaitMock.expectedInvocations), m.WaitMock.expectedInvocationsOrigin, afterWaitCounter)
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *LimiterMock) MinimockFinish() {
	m.finishOnce.Do(func() {
		if !m.minimockDone() {
			m.MinimockWaitInspect()
		}
	})
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *LimiterMock) MinimockWait(timeout mm_time.Duration) {
	timeoutCh := mm_time.After(timeout)
	for {
		if m.minimockDone() {
			return
		}
		select {
		case <-timeoutCh:
			m.MinimockFinish()
			return
		case <-mm_time.After(10 * mm_time.Millisecond):
		}
	}
}

func (m *LimiterMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockWaitDone()
}
