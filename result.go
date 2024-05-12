package result

// Result is a generic type that can be used to represent a successful or failed operation.
type Result[T any] struct {
	success T
	failure error
}

// Success creates a new Result with a successful value.
func Success[T any](value T) Result[T] {
	return Result[T]{success: value}
}

// Failure creates a new Result with a failure value.
func Failure[T any](err error) Result[T] {
	return Result[T]{failure: err}
}

// IsSuccess returns true if the Result is a success.
func (r Result[T]) IsSuccess() bool {
	return r.failure == nil
}

// IsFailure returns true if the Result is a failure.
func (r Result[T]) IsFailure() bool {
	return r.failure != nil
}

// GetSuccess returns the successful value of the Result.
func (r Result[T]) GetSuccess() T {
	return r.success
}

// GetFailure returns the failure value of the Result.
func (r Result[T]) GetFailure() error {
	return r.failure
}

// Map applies a function to the successful value of the Result if it is a success and returns a new Result.
// Should be used when the mapping function returns a value.
// This is useful when you have a function that can't return an error.
//
// Example:
//
//	func incrementAge(p Person) Person {
//		p.Age++
//		return p
//	}
//
//	func main() {
//		r := result.Success[Person](Person{Name: "John", Age: 30})
//		plusAgeResult := result.Map(r, incrementAge)
//		fmt.Println(plusAgeResult.GetSuccess().Age) // 31
//	}
func Map[T any, S any](r Result[T], f func(T) S) Result[S] {
	if r.IsSuccess() {
		return Success(f(r.GetSuccess()))
	}
	return Result[S]{failure: r.failure}
}

// Map applies a function to the successful value of the Result if it is a success and returns a new Result.
// Should be used when the mapping function returns a value of the same type, then you can chain the calls.
// This is useful when you have a function that can't return an error.
//
// Example:
//
//	func incrementAge(p Person) Person {
//		p.Age++
//		return p
//	}
//
//	func main() {
//		r := result.Success[Person](Person{Name: "John", Age: 30})
//		incrementedAge := r.Map(incrementAge).Map(incrementAge)
//		fmt.Println(incrementedAge.GetSuccess().Age) // 32, because the age was incremented twice
//	}
func (r Result[T]) Map(f func(T) T) Result[T] {
	if r.IsSuccess() {
		return Success(f(r.GetSuccess()))
	}
	return Result[T]{failure: r.failure}
}

// FlatMap applies a function to the successful value of the Result if it is a success and returns a new Result.
// Should be used when the mapping function returns a Result.
// This is useful when you have a function that can return an error.
//
// Example:
//
//	func canReturnError(p Person) result.Result[Person] {
//		p, err := returnsError(p)
//		if err != nil {
//			return result.Failure[Person](err)
//		}
//		return result.Success[Person](p)
//	}
//
//	func main() {
//		r := result.Success[Person](Person{Name: "John", Age: 30})
//		potentialErrorResult := result.FlatMap(r, canReturnError)
//	}
func FlatMap[T any, S any](r Result[T], f func(T) Result[S]) Result[S] {
	if r.IsSuccess() {
		return f(r.GetSuccess())
	}
	return Result[S]{failure: r.failure}
}

// FlatMap applies a function to the successful value of the Result if it is a success and returns a new Result.
// Should be used when the mapping function returns a Result of the same type, then you can chain the calls.
// This is useful when you have a function that can return an error.
//
// Example:
//
//	func canReturnError(p Person) result.Result[Person] {
//		p, err := returnsError(p)
//		if err != nil {
//			return result.Failure[Person](err)
//		}
//		return result.Success[Person](p)
//	}
//
//	func main() {
//		potentialError := r.FlatMap(canReturnError).Map(incrementAge)
//		fmt.Println(potentialError.IsFailure()) // true, because first FlatMap returns an error
//	}
func (r Result[T]) FlatMap(f func(T) Result[T]) Result[T] {
	if r.IsSuccess() {
		return f(r.GetSuccess())
	}
	return Result[T]{failure: r.failure}
}
