package errors

type Error struct {
	Level   string
	Message string
}

type ErrorsBag struct {
	errors []Error
}

func (eb *ErrorsBag) Add(level string, message string) {
	eb.errors = append(eb.errors, Error{
		Level:   level,
		Message: message,
	})
}

func (eb *ErrorsBag) Any() bool {
	return len(eb.errors) > 0
}

func (eb *ErrorsBag) All() []Error {
	return eb.errors
}

func (eb *ErrorsBag) Get(index int) Error {
	return eb.errors[index]
}
