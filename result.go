package result

// Result is a generic type that can be used to represent a successful or failed operation.
type Result[S any] struct {
	success S
	failure error
}

// Success creates a new Result with a successful value.
func Success[S any](value S) Result[S] {
	return Result[S]{success: value}
}

// Failure creates a new Result with a failure value.
func Failure[S any](err error) Result[S] {
	return Result[S]{failure: err}
}

// IsSuccess returns true if the Result is a success.
func (r Result[S]) IsSuccess() bool {
	return r.failure == nil
}

// IsFailure returns true if the Result is a failure.
func (r Result[S]) IsFailure() bool {
	return r.failure != nil
}

// GetSuccess returns the successful value of the Result.
func (r Result[S]) GetSuccess() S {
	return r.success
}

// GetFailure returns the failure value of the Result.
func (r Result[S]) GetFailure() error {
	return r.failure
}

// GetSuccessOrElse returns the success value or a default value if it's a failure
func (r Result[S]) GetSuccessOrElse(defaultValue S) S {
	if r.IsSuccess() {
		return r.success
	}
	return defaultValue
}

// Try converts a typical Go function return pattern (value, error) into a Result
// This is useful when you have a function that can return an error and you want to convert it into a Result.
// Example:
//
//	func returnsError() (int, error) {
//		return 42, nil
//	}
//
//	func main() {
//		value, err := returnsError()
//		result := result.Try(value, err)
//		if result.IsSuccess() {
//			fmt.Println("Success:", result.GetSuccess())
//		} else {
//			fmt.Println("Error:", result.GetFailure())
//		}
//	}
func Try[S any](value S, err error) Result[S] {
	if err != nil {
		return Failure[S](err)
	}
	return Success(value)
}

// ThenTry is a method that converts a typical Go function return pattern (value, error) into a Result on the Result type.
// This is useful when you have a function that can return an error and you want to use it with the Result type.
// Example:
//
//	func returnsError() (value, error) {
//		...
//	}
//
//	func main() {
//		value, err := returnsError()
//		result := result.Success(Person{Name: "John", Age: 30}).Try(returnsError())
//		if result.IsSuccess() {
//			fmt.Println("Success:", result.GetSuccess())
//		} else {
//			fmt.Println("Error:", result.GetFailure())
//		}
//	}
func (r Result[S]) ThenTry(value S, err error) Result[S] {
	if r.IsFailure() {
		return r
	}

	if err != nil {
		return Failure[S](err)
	}
	
	return Success(value)
}

// Transform applies a function to the successful value of the Result if it is a success
// and returns a new Result with a new success type.
//
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
//		r := result.Success(Person{Name: "John", Age: 30})
//		plusAgeResult := result.Transform(r, incrementAge)
//		fmt.Println(plusAgeResult.GetSuccess().Age) // 31
//	}
func Transform[S, NS any](r Result[S], f func(S) NS) Result[NS] {
	if r.IsFailure() {
		return Result[NS]{failure: r.failure}
	}
	return Success[NS](f(r.GetSuccess()))
}

// Then applies a function to the successful value of the Result if it is a success
// and returns a new Result with the same success type.
//
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
//		incrementedAge := r.Then(incrementAge).Then(incrementAge)
//		fmt.Println(incrementedAge.GetSuccess().Age) // 32, because the age was incremented twice
//	}
func (r Result[S]) Then(f func(S) S) Result[S] {
	if r.IsFailure() {
		return Result[S]{failure: r.failure}
	}
	return Success[S](f(r.GetSuccess()))
}

// TransformWith applies a function to the successful value of the Result if it is a success
// and returns a new Result with a new success type.
//
// Should be used when the mapping function returns a Result.
// This is useful when you have a function that can return an error.
//
// Example:
//
//	func canReturnError(p Person) result.Result[Person] {
//		p, err := returnsError(p)
//		if err != nil {
//			return result.Failure(err)
//		}
//		return result.Success(p)
//	}
//
//	func main() {
//		r := result.Success(Person{Name: "John", Age: 30})
//		potentialErrorResult := result.TransformWith(r, canReturnError)
//	}
func TransformWith[S, NS any](r Result[S], f func(S) Result[NS]) Result[NS] {
	if r.IsFailure() {
		return Result[NS]{failure: r.failure}
	}
	return f(r.GetSuccess())
}

// ThenWith applies a function to the successful value of the Result if it is a success
// and returns the Result returned by the function.
//
// This is equivalent to TransformWith but maintains the same type.
// It's useful when you have a function that returns a Result of the same type.
//
// Example:
//
//	func mayFail(i int) result.Result[int] {
//		if i < 0 {
//			return result.Failure[int](errors.New("negative value"))
//		}
//		return result.Success(i * 2)
//	}
//
//	func main() {
//		r := result.Success(42)
//		result := r.ThenWith(mayFail)
//		// result is Success(84)
//	}
func (r Result[S]) ThenWith(f func(S) Result[S]) Result[S] {
	if r.IsFailure() {
		return r
	}
	return f(r.GetSuccess())
}
