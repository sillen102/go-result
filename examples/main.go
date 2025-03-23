package main

import (
	"fmt"

	"github.com/sillen102/result"
)

type Person struct {
	Name string
	Age  int
}

func main() {
	r := result.Success(Person{Name: "John", Age: 30})
	incrementedAge := r.Then(incrementAge).Then(incrementAge)
	fmt.Println(incrementedAge.GetSuccess().Age) // 32, because the age was incremented twice

	r = result.Success(Person{Name: "John", Age: 30})
	potentialError := r.ThenTry(canReturnError).Then(incrementAge)
	fmt.Println(potentialError.IsFailure()) // true, because first FlatMap returns an error
}

func canReturnError(p Person) result.Result[Person] {
	p, err := returnsError(p)
	if err != nil {
		return result.Failure[Person](err)
	}
	return result.Success[Person](p)
}

func returnsError(_ Person) (Person, error) {
	return Person{}, fmt.Errorf("error")
}

func incrementAge(p Person) Person {
	p.Age++
	return p
}
