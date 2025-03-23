package result_test

import (
	"errors"
	"testing"

	"github.com/sillen102/result"
)

func TestSuccess(t *testing.T) {
	r := result.Success(42)

	if !r.IsSuccess() {
		t.Error("Expected IsSuccess to be true")
	}

	if r.IsFailure() {
		t.Error("Expected IsFailure to be false")
	}

	if r.GetSuccess() != 42 {
		t.Errorf("Expected success value to be 42, got %v", r.GetSuccess())
	}

	if r.GetFailure() != nil {
		t.Errorf("Expected failure value to be nil, got %v", r.GetFailure())
	}
}

func TestFailure(t *testing.T) {
	testErr := errors.New("test error")
	r := result.Failure[int](testErr)

	if r.IsSuccess() {
		t.Error("Expected IsSuccess to be false")
	}

	if !r.IsFailure() {
		t.Error("Expected IsFailure to be true")
	}

	if !errors.Is(testErr, r.GetFailure()) {
		t.Errorf("Expected failure value to be %v, got %v", testErr, r.GetFailure())
	}
}

func TestTransform(t *testing.T) {
	// Test Transform with success
	r := result.Success(42)
	transformedR := result.Transform(r, func(i int) string {
		return "value: " + string(rune(i))
	})

	if !transformedR.IsSuccess() {
		t.Error("Expected transformed result to be success")
	}

	if transformedR.GetSuccess() != "value: *" {
		t.Errorf("Expected transformed value to be 'value: *', got %v", transformedR.GetSuccess())
	}

	// Test Transform with failure
	testErr := errors.New("test error")
	r = result.Failure[int](testErr)
	transformedR = result.Transform(r, func(i int) string {
		return "value: " + string(rune(i))
	})

	if !transformedR.IsFailure() {
		t.Error("Expected transformed result to be failure")
	}

	if !errors.Is(testErr, transformedR.GetFailure()) {
		t.Errorf("Expected failure to be preserved, got %v", transformedR.GetFailure())
	}
}

func TestResultThen(t *testing.T) {
	// Test Result.Then with success
	r := result.Success(42)
	transformedR := r.Then(func(i int) int {
		return i * 2
	})

	if !transformedR.IsSuccess() {
		t.Error("Expected transformed result to be success")
	}

	if transformedR.GetSuccess() != 84 {
		t.Errorf("Expected transformed value to be 84, got %v", transformedR.GetSuccess())
	}

	// Test Result.Then with failure
	testErr := errors.New("test error")
	r = result.Failure[int](testErr)
	transformedR = r.Then(func(i int) int {
		return i * 2
	})

	if !transformedR.IsFailure() {
		t.Error("Expected transformed result to be failure")
	}

	if !errors.Is(testErr, transformedR.GetFailure()) {
		t.Errorf("Expected failure to be preserved, got %v", transformedR.GetFailure())
	}
}

func TestThenWith(t *testing.T) {
	// Test ThenWith with success -> success
	r := result.Success(42)
	thenWithR := result.ThenWith(r, func(i int) result.Result[string] {
		return result.Success("value: " + string(rune(i)))
	})

	if !thenWithR.IsSuccess() {
		t.Error("Expected then-with result to be success")
	}

	if thenWithR.GetSuccess() != "value: *" {
		t.Errorf("Expected then-with value to be 'value: *', got %v", thenWithR.GetSuccess())
	}

	// Test ThenWith with success -> failure
	r = result.Success(42)
	testErr := errors.New("function error")
	thenWithR = result.ThenWith(r, func(i int) result.Result[string] {
		return result.Failure[string](testErr)
	})

	if !thenWithR.IsFailure() {
		t.Error("Expected then-with result to be failure")
	}

	if !errors.Is(testErr, thenWithR.GetFailure()) {
		t.Errorf("Expected failure from function, got %v", thenWithR.GetFailure())
	}

	// Test ThenWith with failure
	r = result.Failure[int](errors.New("initial error"))
	thenWithR = result.ThenWith(r, func(i int) result.Result[string] {
		return result.Success("should not reach here")
	})

	if !thenWithR.IsFailure() {
		t.Error("Expected then-with result to be failure")
	}
}

func TestResultThenTry(t *testing.T) {
	// Test Result.ThenTry with success -> success
	r := result.Success(42)
	thenTryR := r.ThenTry(func(i int) result.Result[int] {
		return result.Success(i * 2)
	})

	if !thenTryR.IsSuccess() {
		t.Error("Expected then-try result to be success")
	}

	if thenTryR.GetSuccess() != 84 {
		t.Errorf("Expected then-try value to be 84, got %v", thenTryR.GetSuccess())
	}

	// Test Result.ThenTry with success -> failure
	r = result.Success(42)
	testErr := errors.New("function error")
	thenTryR = r.ThenTry(func(i int) result.Result[int] {
		return result.Failure[int](testErr)
	})

	if !thenTryR.IsFailure() {
		t.Error("Expected then-try result to be failure")
	}

	if !errors.Is(testErr, thenTryR.GetFailure()) {
		t.Errorf("Expected failure from function, got %v", thenTryR.GetFailure())
	}

	// Test Result.ThenTry with failure
	r = result.Failure[int](errors.New("initial error"))
	thenTryR = r.ThenTry(func(i int) result.Result[int] {
		return result.Success(i * 2)
	})

	if !thenTryR.IsFailure() {
		t.Error("Expected then-try result to be failure")
	}
}

func TestChaining(t *testing.T) {
	// Test chaining Then and ThenTry
	r := result.Success(42)

	chainedResult := r.
		Then(func(i int) int { return i + 1 }).
		ThenTry(func(i int) result.Result[int] {
			if i > 40 {
				return result.Success(i * 2)
			}
			return result.Failure[int](errors.New("value too small"))
		}).
		Then(func(i int) int { return i - 10 })

	if !chainedResult.IsSuccess() {
		t.Error("Expected chained result to be success")
	}

	if chainedResult.GetSuccess() != 76 { // (42+1)*2-10 = 76
		t.Errorf("Expected chained value to be 76, got %v", chainedResult.GetSuccess())
	}

	// Test chaining with failure
	testErr := errors.New("test error")
	r = result.Failure[int](testErr)

	chainedResult = r.
		Then(func(i int) int { return i + 1 }).
		ThenTry(func(i int) result.Result[int] { return result.Success(i * 2) })

	if !chainedResult.IsFailure() {
		t.Error("Expected chained result to be failure")
	}

	if !errors.Is(testErr, chainedResult.GetFailure()) {
		t.Errorf("Expected failure to be preserved through chain, got %v", chainedResult.GetFailure())
	}
}

func TestGetSuccessOr(t *testing.T) {
	// Test with success
	r := result.Success(42)
	value := r.GetSuccessOr(0)
	if value != 42 {
		t.Errorf("Expected 42, got %v", value)
	}

	// Test with failure
	r = result.Failure[int](errors.New("test error"))
	value = r.GetSuccessOr(0)
	if value != 0 {
		t.Errorf("Expected default value 0, got %v", value)
	}
}

func TestMatch(t *testing.T) {
	// Test with success
	r := result.Success(42)
	successCalled := false
	failureCalled := false

	r.Match(
		func(i int) {
			successCalled = true
			if i != 42 {
				t.Errorf("Expected 42, got %v", i)
			}
		},
		func(err error) {
			failureCalled = true
		},
	)

	if !successCalled {
		t.Error("Success function was not called")
	}
	if failureCalled {
		t.Error("Failure function was called unexpectedly")
	}

	// Test with failure
	testErr := errors.New("test error")
	r = result.Failure[int](testErr)
	successCalled = false
	failureCalled = false
	var capturedErr error

	r.Match(
		func(i int) {
			successCalled = true
		},
		func(err error) {
			failureCalled = true
			capturedErr = err
		},
	)

	if successCalled {
		t.Error("Success function was called unexpectedly")
	}
	if !failureCalled {
		t.Error("Failure function was not called")
	}
	if !errors.Is(capturedErr, testErr) {
		t.Errorf("Expected error %v, got %v", testErr, capturedErr)
	}
}

func TestTry(t *testing.T) {
	// Test with success case
	successFunc := func() (int, error) {
		return 42, nil
	}

	r := result.Try(successFunc())

	if !r.IsSuccess() {
		t.Error("Expected result to be success")
	}
	if r.GetSuccess() != 42 {
		t.Errorf("Expected 42, got %v", r.GetSuccess())
	}

	// Test with failure case
	testErr := errors.New("test error")
	failureFunc := func() (string, error) {
		return "", testErr
	}

	r2 := result.Try(failureFunc())

	if !r2.IsFailure() {
		t.Error("Expected result to be failure")
	}
	if !errors.Is(r2.GetFailure(), testErr) {
		t.Errorf("Expected error %v, got %v", testErr, r2.GetFailure())
	}
}
