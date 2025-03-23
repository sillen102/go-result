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

func TestMap(t *testing.T) {
	// Test Map with success
	r := result.Success(42)
	mappedR := result.Map(r, func(i int) string {
		return "value: " + string(rune(i))
	})

	if !mappedR.IsSuccess() {
		t.Error("Expected mapped result to be success")
	}

	if mappedR.GetSuccess() != "value: *" {
		t.Errorf("Expected mapped value to be 'value: *', got %v", mappedR.GetSuccess())
	}

	// Test Map with failure
	testErr := errors.New("test error")
	r = result.Failure[int](testErr)
	mappedR = result.Map(r, func(i int) string {
		return "value: " + string(rune(i))
	})

	if !mappedR.IsFailure() {
		t.Error("Expected mapped result to be failure")
	}

	if !errors.Is(testErr, mappedR.GetFailure()) {
		t.Errorf("Expected failure to be preserved, got %v", mappedR.GetFailure())
	}
}

func TestResultMap(t *testing.T) {
	// Test Result.Map with success
	r := result.Success(42)
	mappedR := r.Map(func(i int) int {
		return i * 2
	})

	if !mappedR.IsSuccess() {
		t.Error("Expected mapped result to be success")
	}

	if mappedR.GetSuccess() != 84 {
		t.Errorf("Expected mapped value to be 84, got %v", mappedR.GetSuccess())
	}

	// Test Result.Map with failure
	testErr := errors.New("test error")
	r = result.Failure[int](testErr)
	mappedR = r.Map(func(i int) int {
		return i * 2
	})

	if !mappedR.IsFailure() {
		t.Error("Expected mapped result to be failure")
	}

	if !errors.Is(testErr, mappedR.GetFailure()) {
		t.Errorf("Expected failure to be preserved, got %v", mappedR.GetFailure())
	}
}

func TestFlatMap(t *testing.T) {
	// Test FlatMap with success -> success
	r := result.Success(42)
	flatMappedR := result.FlatMap(r, func(i int) result.Result[string] {
		return result.Success("value: " + string(rune(i)))
	})

	if !flatMappedR.IsSuccess() {
		t.Error("Expected flat-mapped result to be success")
	}

	if flatMappedR.GetSuccess() != "value: *" {
		t.Errorf("Expected flat-mapped value to be 'value: *', got %v", flatMappedR.GetSuccess())
	}

	// Test FlatMap with success -> failure
	r = result.Success(42)
	testErr := errors.New("function error")
	flatMappedR = result.FlatMap(r, func(i int) result.Result[string] {
		return result.Failure[string](testErr)
	})

	if !flatMappedR.IsFailure() {
		t.Error("Expected flat-mapped result to be failure")
	}

	if !errors.Is(testErr, flatMappedR.GetFailure()) {
		t.Errorf("Expected failure from function, got %v", flatMappedR.GetFailure())
	}

	// Test FlatMap with failure
	r = result.Failure[int](errors.New("initial error"))
	flatMappedR = result.FlatMap(r, func(i int) result.Result[string] {
		return result.Success("should not reach here")
	})

	if !flatMappedR.IsFailure() {
		t.Error("Expected flat-mapped result to be failure")
	}
}

func TestResultFlatMap(t *testing.T) {
	// Test Result.FlatMap with success -> success
	r := result.Success(42)
	flatMappedR := r.FlatMap(func(i int) result.Result[int] {
		return result.Success(i * 2)
	})

	if !flatMappedR.IsSuccess() {
		t.Error("Expected flat-mapped result to be success")
	}

	if flatMappedR.GetSuccess() != 84 {
		t.Errorf("Expected flat-mapped value to be 84, got %v", flatMappedR.GetSuccess())
	}

	// Test Result.FlatMap with success -> failure
	r = result.Success(42)
	testErr := errors.New("function error")
	flatMappedR = r.FlatMap(func(i int) result.Result[int] {
		return result.Failure[int](testErr)
	})

	if !flatMappedR.IsFailure() {
		t.Error("Expected flat-mapped result to be failure")
	}

	if !errors.Is(testErr, flatMappedR.GetFailure()) {
		t.Errorf("Expected failure from function, got %v", flatMappedR.GetFailure())
	}

	// Test Result.FlatMap with failure
	r = result.Failure[int](errors.New("initial error"))
	flatMappedR = r.FlatMap(func(i int) result.Result[int] {
		return result.Success(i * 2)
	})

	if !flatMappedR.IsFailure() {
		t.Error("Expected flat-mapped result to be failure")
	}
}

func TestChaining(t *testing.T) {
	// Test chaining Map and FlatMap
	r := result.Success(42)

	chainedResult := r.
		Map(func(i int) int { return i + 1 }).
		FlatMap(func(i int) result.Result[int] {
			if i > 40 {
				return result.Success(i * 2)
			}
			return result.Failure[int](errors.New("value too small"))
		}).
		Map(func(i int) int { return i - 10 })

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
		Map(func(i int) int { return i + 1 }).
		FlatMap(func(i int) result.Result[int] { return result.Success(i * 2) })

	if !chainedResult.IsFailure() {
		t.Error("Expected chained result to be failure")
	}

	if !errors.Is(testErr, chainedResult.GetFailure()) {
		t.Errorf("Expected failure to be preserved through chain, got %v", chainedResult.GetFailure())
	}
}
