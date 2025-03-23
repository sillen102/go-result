package result_test

import (
	"errors"
	"testing"

	"github.com/sillen102/go-result"
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
	thenWithR := result.TransformWith(r, func(i int) result.Result[string] {
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
	thenWithR = result.TransformWith(r, func(i int) result.Result[string] {
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
	thenWithR = result.TransformWith(r, func(i int) result.Result[string] {
		return result.Success("should not reach here")
	})

	if !thenWithR.IsFailure() {
		t.Error("Expected then-with result to be failure")
	}
}

func TestChaining(t *testing.T) {
	// Test chaining Then and ThenTry
	r := result.Success(42)

	chainedResult := r.
		Then(func(i int) int { return i + 1 }).
		ThenWith(func(i int) result.Result[int] {
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
		ThenWith(func(i int) result.Result[int] { return result.Success(i * 2) })

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
	value := r.GetSuccessOrElse(0)
	if value != 42 {
		t.Errorf("Expected 42, got %v", value)
	}

	// Test with failure
	r = result.Failure[int](errors.New("test error"))
	value = r.GetSuccessOrElse(0)
	if value != 0 {
		t.Errorf("Expected default value 0, got %v", value)
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

func TestResultTryWith(t *testing.T) {
	// Test Try on a success result with success value
	r1 := result.Success[int](10)
	result1 := r1.ThenTry(20, nil)

	if !result1.IsSuccess() {
		t.Error("Expected result to be success")
	}
	if result1.GetSuccess() != 20 {
		t.Errorf("Expected 20, got %v", result1.GetSuccess())
	}

	// Test Try on a success result with error
	testErr := errors.New("test error")
	r2 := result.Success[int](10)
	result2 := r2.ThenTry(20, testErr)

	if !result2.IsFailure() {
		t.Error("Expected result to be failure")
	}
	if !errors.Is(result2.GetFailure(), testErr) {
		t.Errorf("Expected error %v, got %v", testErr, result2.GetFailure())
	}

	// Test Try on a failure result - should ignore new values
	originalErr := errors.New("original error")
	r3 := result.Failure[int](originalErr)
	result3 := r3.ThenTry(30, nil)

	if !result3.IsFailure() {
		t.Error("Expected result to be failure")
	}
	if !errors.Is(result3.GetFailure(), originalErr) {
		t.Errorf("Expected original error to be preserved, got %v", result3.GetFailure())
	}

	// Test Try on a failure result with another error - should still ignore
	r4 := result.Failure[int](originalErr)
	newErr := errors.New("new error")
	result4 := r4.ThenTry(40, newErr)

	if !result4.IsFailure() {
		t.Error("Expected result to be failure")
	}
	if !errors.Is(result4.GetFailure(), originalErr) {
		t.Errorf("Expected original error to be preserved, got %v", result4.GetFailure())
	}
}
