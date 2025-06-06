// Code generated by http://github.com/gojuno/minimock (v3.4.5). DO NOT EDIT.

package cart

//go:generate minimock -i route256/cart/internal/usecases/cart.productClient -o product_client_mock_test.go -n ProductClientMock -p cart

import (
	"context"
	"route256/cart/internal/models"
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock/v3"
)

// ProductClientMock implements productClient
type ProductClientMock struct {
	t          minimock.Tester
	finishOnce sync.Once

	funcGetItem          func(ctx context.Context, skuID int64) (p1 models.Product, err error)
	funcGetItemOrigin    string
	inspectFuncGetItem   func(ctx context.Context, skuID int64)
	afterGetItemCounter  uint64
	beforeGetItemCounter uint64
	GetItemMock          mProductClientMockGetItem
}

// NewProductClientMock returns a mock for productClient
func NewProductClientMock(t minimock.Tester) *ProductClientMock {
	m := &ProductClientMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetItemMock = mProductClientMockGetItem{mock: m}
	m.GetItemMock.callArgs = []*ProductClientMockGetItemParams{}

	t.Cleanup(m.MinimockFinish)

	return m
}

type mProductClientMockGetItem struct {
	optional           bool
	mock               *ProductClientMock
	defaultExpectation *ProductClientMockGetItemExpectation
	expectations       []*ProductClientMockGetItemExpectation

	callArgs []*ProductClientMockGetItemParams
	mutex    sync.RWMutex

	expectedInvocations       uint64
	expectedInvocationsOrigin string
}

// ProductClientMockGetItemExpectation specifies expectation struct of the productClient.GetItem
type ProductClientMockGetItemExpectation struct {
	mock               *ProductClientMock
	params             *ProductClientMockGetItemParams
	paramPtrs          *ProductClientMockGetItemParamPtrs
	expectationOrigins ProductClientMockGetItemExpectationOrigins
	results            *ProductClientMockGetItemResults
	returnOrigin       string
	Counter            uint64
}

// ProductClientMockGetItemParams contains parameters of the productClient.GetItem
type ProductClientMockGetItemParams struct {
	ctx   context.Context
	skuID int64
}

// ProductClientMockGetItemParamPtrs contains pointers to parameters of the productClient.GetItem
type ProductClientMockGetItemParamPtrs struct {
	ctx   *context.Context
	skuID *int64
}

// ProductClientMockGetItemResults contains results of the productClient.GetItem
type ProductClientMockGetItemResults struct {
	p1  models.Product
	err error
}

// ProductClientMockGetItemOrigins contains origins of expectations of the productClient.GetItem
type ProductClientMockGetItemExpectationOrigins struct {
	origin      string
	originCtx   string
	originSkuID string
}

// Marks this method to be optional. The default behavior of any method with Return() is '1 or more', meaning
// the test will fail minimock's automatic final call check if the mocked method was not called at least once.
// Optional() makes method check to work in '0 or more' mode.
// It is NOT RECOMMENDED to use this option unless you really need it, as default behaviour helps to
// catch the problems when the expected method call is totally skipped during test run.
func (mmGetItem *mProductClientMockGetItem) Optional() *mProductClientMockGetItem {
	mmGetItem.optional = true
	return mmGetItem
}

// Expect sets up expected params for productClient.GetItem
func (mmGetItem *mProductClientMockGetItem) Expect(ctx context.Context, skuID int64) *mProductClientMockGetItem {
	if mmGetItem.mock.funcGetItem != nil {
		mmGetItem.mock.t.Fatalf("ProductClientMock.GetItem mock is already set by Set")
	}

	if mmGetItem.defaultExpectation == nil {
		mmGetItem.defaultExpectation = &ProductClientMockGetItemExpectation{}
	}

	if mmGetItem.defaultExpectation.paramPtrs != nil {
		mmGetItem.mock.t.Fatalf("ProductClientMock.GetItem mock is already set by ExpectParams functions")
	}

	mmGetItem.defaultExpectation.params = &ProductClientMockGetItemParams{ctx, skuID}
	mmGetItem.defaultExpectation.expectationOrigins.origin = minimock.CallerInfo(1)
	for _, e := range mmGetItem.expectations {
		if minimock.Equal(e.params, mmGetItem.defaultExpectation.params) {
			mmGetItem.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmGetItem.defaultExpectation.params)
		}
	}

	return mmGetItem
}

// ExpectCtxParam1 sets up expected param ctx for productClient.GetItem
func (mmGetItem *mProductClientMockGetItem) ExpectCtxParam1(ctx context.Context) *mProductClientMockGetItem {
	if mmGetItem.mock.funcGetItem != nil {
		mmGetItem.mock.t.Fatalf("ProductClientMock.GetItem mock is already set by Set")
	}

	if mmGetItem.defaultExpectation == nil {
		mmGetItem.defaultExpectation = &ProductClientMockGetItemExpectation{}
	}

	if mmGetItem.defaultExpectation.params != nil {
		mmGetItem.mock.t.Fatalf("ProductClientMock.GetItem mock is already set by Expect")
	}

	if mmGetItem.defaultExpectation.paramPtrs == nil {
		mmGetItem.defaultExpectation.paramPtrs = &ProductClientMockGetItemParamPtrs{}
	}
	mmGetItem.defaultExpectation.paramPtrs.ctx = &ctx
	mmGetItem.defaultExpectation.expectationOrigins.originCtx = minimock.CallerInfo(1)

	return mmGetItem
}

// ExpectSkuIDParam2 sets up expected param skuID for productClient.GetItem
func (mmGetItem *mProductClientMockGetItem) ExpectSkuIDParam2(skuID int64) *mProductClientMockGetItem {
	if mmGetItem.mock.funcGetItem != nil {
		mmGetItem.mock.t.Fatalf("ProductClientMock.GetItem mock is already set by Set")
	}

	if mmGetItem.defaultExpectation == nil {
		mmGetItem.defaultExpectation = &ProductClientMockGetItemExpectation{}
	}

	if mmGetItem.defaultExpectation.params != nil {
		mmGetItem.mock.t.Fatalf("ProductClientMock.GetItem mock is already set by Expect")
	}

	if mmGetItem.defaultExpectation.paramPtrs == nil {
		mmGetItem.defaultExpectation.paramPtrs = &ProductClientMockGetItemParamPtrs{}
	}
	mmGetItem.defaultExpectation.paramPtrs.skuID = &skuID
	mmGetItem.defaultExpectation.expectationOrigins.originSkuID = minimock.CallerInfo(1)

	return mmGetItem
}

// Inspect accepts an inspector function that has same arguments as the productClient.GetItem
func (mmGetItem *mProductClientMockGetItem) Inspect(f func(ctx context.Context, skuID int64)) *mProductClientMockGetItem {
	if mmGetItem.mock.inspectFuncGetItem != nil {
		mmGetItem.mock.t.Fatalf("Inspect function is already set for ProductClientMock.GetItem")
	}

	mmGetItem.mock.inspectFuncGetItem = f

	return mmGetItem
}

// Return sets up results that will be returned by productClient.GetItem
func (mmGetItem *mProductClientMockGetItem) Return(p1 models.Product, err error) *ProductClientMock {
	if mmGetItem.mock.funcGetItem != nil {
		mmGetItem.mock.t.Fatalf("ProductClientMock.GetItem mock is already set by Set")
	}

	if mmGetItem.defaultExpectation == nil {
		mmGetItem.defaultExpectation = &ProductClientMockGetItemExpectation{mock: mmGetItem.mock}
	}
	mmGetItem.defaultExpectation.results = &ProductClientMockGetItemResults{p1, err}
	mmGetItem.defaultExpectation.returnOrigin = minimock.CallerInfo(1)
	return mmGetItem.mock
}

// Set uses given function f to mock the productClient.GetItem method
func (mmGetItem *mProductClientMockGetItem) Set(f func(ctx context.Context, skuID int64) (p1 models.Product, err error)) *ProductClientMock {
	if mmGetItem.defaultExpectation != nil {
		mmGetItem.mock.t.Fatalf("Default expectation is already set for the productClient.GetItem method")
	}

	if len(mmGetItem.expectations) > 0 {
		mmGetItem.mock.t.Fatalf("Some expectations are already set for the productClient.GetItem method")
	}

	mmGetItem.mock.funcGetItem = f
	mmGetItem.mock.funcGetItemOrigin = minimock.CallerInfo(1)
	return mmGetItem.mock
}

// When sets expectation for the productClient.GetItem which will trigger the result defined by the following
// Then helper
func (mmGetItem *mProductClientMockGetItem) When(ctx context.Context, skuID int64) *ProductClientMockGetItemExpectation {
	if mmGetItem.mock.funcGetItem != nil {
		mmGetItem.mock.t.Fatalf("ProductClientMock.GetItem mock is already set by Set")
	}

	expectation := &ProductClientMockGetItemExpectation{
		mock:               mmGetItem.mock,
		params:             &ProductClientMockGetItemParams{ctx, skuID},
		expectationOrigins: ProductClientMockGetItemExpectationOrigins{origin: minimock.CallerInfo(1)},
	}
	mmGetItem.expectations = append(mmGetItem.expectations, expectation)
	return expectation
}

// Then sets up productClient.GetItem return parameters for the expectation previously defined by the When method
func (e *ProductClientMockGetItemExpectation) Then(p1 models.Product, err error) *ProductClientMock {
	e.results = &ProductClientMockGetItemResults{p1, err}
	return e.mock
}

// Times sets number of times productClient.GetItem should be invoked
func (mmGetItem *mProductClientMockGetItem) Times(n uint64) *mProductClientMockGetItem {
	if n == 0 {
		mmGetItem.mock.t.Fatalf("Times of ProductClientMock.GetItem mock can not be zero")
	}
	mm_atomic.StoreUint64(&mmGetItem.expectedInvocations, n)
	mmGetItem.expectedInvocationsOrigin = minimock.CallerInfo(1)
	return mmGetItem
}

func (mmGetItem *mProductClientMockGetItem) invocationsDone() bool {
	if len(mmGetItem.expectations) == 0 && mmGetItem.defaultExpectation == nil && mmGetItem.mock.funcGetItem == nil {
		return true
	}

	totalInvocations := mm_atomic.LoadUint64(&mmGetItem.mock.afterGetItemCounter)
	expectedInvocations := mm_atomic.LoadUint64(&mmGetItem.expectedInvocations)

	return totalInvocations > 0 && (expectedInvocations == 0 || expectedInvocations == totalInvocations)
}

// GetItem implements productClient
func (mmGetItem *ProductClientMock) GetItem(ctx context.Context, skuID int64) (p1 models.Product, err error) {
	mm_atomic.AddUint64(&mmGetItem.beforeGetItemCounter, 1)
	defer mm_atomic.AddUint64(&mmGetItem.afterGetItemCounter, 1)

	mmGetItem.t.Helper()

	if mmGetItem.inspectFuncGetItem != nil {
		mmGetItem.inspectFuncGetItem(ctx, skuID)
	}

	mm_params := ProductClientMockGetItemParams{ctx, skuID}

	// Record call args
	mmGetItem.GetItemMock.mutex.Lock()
	mmGetItem.GetItemMock.callArgs = append(mmGetItem.GetItemMock.callArgs, &mm_params)
	mmGetItem.GetItemMock.mutex.Unlock()

	for _, e := range mmGetItem.GetItemMock.expectations {
		if minimock.Equal(*e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.p1, e.results.err
		}
	}

	if mmGetItem.GetItemMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmGetItem.GetItemMock.defaultExpectation.Counter, 1)
		mm_want := mmGetItem.GetItemMock.defaultExpectation.params
		mm_want_ptrs := mmGetItem.GetItemMock.defaultExpectation.paramPtrs

		mm_got := ProductClientMockGetItemParams{ctx, skuID}

		if mm_want_ptrs != nil {

			if mm_want_ptrs.ctx != nil && !minimock.Equal(*mm_want_ptrs.ctx, mm_got.ctx) {
				mmGetItem.t.Errorf("ProductClientMock.GetItem got unexpected parameter ctx, expected at\n%s:\nwant: %#v\n got: %#v%s\n",
					mmGetItem.GetItemMock.defaultExpectation.expectationOrigins.originCtx, *mm_want_ptrs.ctx, mm_got.ctx, minimock.Diff(*mm_want_ptrs.ctx, mm_got.ctx))
			}

			if mm_want_ptrs.skuID != nil && !minimock.Equal(*mm_want_ptrs.skuID, mm_got.skuID) {
				mmGetItem.t.Errorf("ProductClientMock.GetItem got unexpected parameter skuID, expected at\n%s:\nwant: %#v\n got: %#v%s\n",
					mmGetItem.GetItemMock.defaultExpectation.expectationOrigins.originSkuID, *mm_want_ptrs.skuID, mm_got.skuID, minimock.Diff(*mm_want_ptrs.skuID, mm_got.skuID))
			}

		} else if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmGetItem.t.Errorf("ProductClientMock.GetItem got unexpected parameters, expected at\n%s:\nwant: %#v\n got: %#v%s\n",
				mmGetItem.GetItemMock.defaultExpectation.expectationOrigins.origin, *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmGetItem.GetItemMock.defaultExpectation.results
		if mm_results == nil {
			mmGetItem.t.Fatal("No results are set for the ProductClientMock.GetItem")
		}
		return (*mm_results).p1, (*mm_results).err
	}
	if mmGetItem.funcGetItem != nil {
		return mmGetItem.funcGetItem(ctx, skuID)
	}
	mmGetItem.t.Fatalf("Unexpected call to ProductClientMock.GetItem. %v %v", ctx, skuID)
	return
}

// GetItemAfterCounter returns a count of finished ProductClientMock.GetItem invocations
func (mmGetItem *ProductClientMock) GetItemAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmGetItem.afterGetItemCounter)
}

// GetItemBeforeCounter returns a count of ProductClientMock.GetItem invocations
func (mmGetItem *ProductClientMock) GetItemBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmGetItem.beforeGetItemCounter)
}

// Calls returns a list of arguments used in each call to ProductClientMock.GetItem.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmGetItem *mProductClientMockGetItem) Calls() []*ProductClientMockGetItemParams {
	mmGetItem.mutex.RLock()

	argCopy := make([]*ProductClientMockGetItemParams, len(mmGetItem.callArgs))
	copy(argCopy, mmGetItem.callArgs)

	mmGetItem.mutex.RUnlock()

	return argCopy
}

// MinimockGetItemDone returns true if the count of the GetItem invocations corresponds
// the number of defined expectations
func (m *ProductClientMock) MinimockGetItemDone() bool {
	if m.GetItemMock.optional {
		// Optional methods provide '0 or more' call count restriction.
		return true
	}

	for _, e := range m.GetItemMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	return m.GetItemMock.invocationsDone()
}

// MinimockGetItemInspect logs each unmet expectation
func (m *ProductClientMock) MinimockGetItemInspect() {
	for _, e := range m.GetItemMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to ProductClientMock.GetItem at\n%s with params: %#v", e.expectationOrigins.origin, *e.params)
		}
	}

	afterGetItemCounter := mm_atomic.LoadUint64(&m.afterGetItemCounter)
	// if default expectation was set then invocations count should be greater than zero
	if m.GetItemMock.defaultExpectation != nil && afterGetItemCounter < 1 {
		if m.GetItemMock.defaultExpectation.params == nil {
			m.t.Errorf("Expected call to ProductClientMock.GetItem at\n%s", m.GetItemMock.defaultExpectation.returnOrigin)
		} else {
			m.t.Errorf("Expected call to ProductClientMock.GetItem at\n%s with params: %#v", m.GetItemMock.defaultExpectation.expectationOrigins.origin, *m.GetItemMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGetItem != nil && afterGetItemCounter < 1 {
		m.t.Errorf("Expected call to ProductClientMock.GetItem at\n%s", m.funcGetItemOrigin)
	}

	if !m.GetItemMock.invocationsDone() && afterGetItemCounter > 0 {
		m.t.Errorf("Expected %d calls to ProductClientMock.GetItem at\n%s but found %d calls",
			mm_atomic.LoadUint64(&m.GetItemMock.expectedInvocations), m.GetItemMock.expectedInvocationsOrigin, afterGetItemCounter)
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *ProductClientMock) MinimockFinish() {
	m.finishOnce.Do(func() {
		if !m.minimockDone() {
			m.MinimockGetItemInspect()
		}
	})
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *ProductClientMock) MinimockWait(timeout mm_time.Duration) {
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

func (m *ProductClientMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockGetItemDone()
}
