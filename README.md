# Spawn

## Overview

Spawn is a lightweight Go package designed to simplify handling of concurrent tasks by providing a simple and intuitive way to spawn goroutines and wait for their completion. It offers a JoinHandle type that represents a handle to a spawned task, allowing you to wait for its completion, retrieve the result, and check if it has finished.
Installation

To use spawn in your Go project, simply import it using go get:

```bash
go get github.com/tembleking/spawn
```

## Usage

### Spawning a Task

To spawn a task, use the `Func` function, passing in the function that you want to execute concurrently. This function returns a `JoinHandle` representing the spawned task.

```go
joinHandle := spawn.Func(func() (result T, err error) {
    // Your concurrent task logic here
    return result, err
})
```

### Waiting for Task Completion

You can wait for the spawned task to complete using the `Wait` or `WaitCtx` methods of the `JoinHandle`.

```go
result, err := joinHandle.Wait() // Waits indefinitely for task completion
```

or

```go
result, err := joinHandle.WaitCtx(context.Background()) // Wait with context
```

### Checking Task Completion

You can also check whether the task has finished executing using the `IsFinished` method.

```go
if handle.IsFinished() {
    // Task has finished
} else {
    // Task is still running
}
```

## Examples

### Basic Example

```go
package main

import (
	"fmt"
	"time"

	"github.com/tembleking/spawn"
)

func main() {
	// Spawn a task
	handle := spawn.Func(func() (result int, err error) {
		time.Sleep(10 * time.Millisecond)
		return 42, nil
	})

	// Wait for the task to complete
	result, err := handle.Wait()
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Result:", result)
	}
}
```

### Example with function that only returns error

```go
package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/tembleking/spawn"
)

func main() {
	funcThatReturnsError := func() error {
		time.Sleep(10 * time.Millisecond)
		return errors.New("some error")
	}

	joinHandle := spawn.Func(func() (struct{}, error) {
		return struct{}{}, funcThatReturnsError()
	})

	_, err := joinHandle.Wait()
	if err != nil {
		fmt.Println("Error:", err)
	}
}
```

### Example with function that returns nothing

```go
package main

import (
    "fmt"
    "time"

    "github.com/tembleking/spawn"
)

func main() {
	funcThatReturnsNothing := func() {
		time.Sleep(10 * time.Millisecond)
	}

	joinHandle := spawn.Func(func() (struct{}, error) {
		funcThatReturnsNothing()
		return struct{}{}, nil
	})

	_, _ = joinHandle.Wait()
	fmt.Println("Task completed")
}
```


## Contributing

Contributions are welcome! If you find any issues or have suggestions for improvements, please open an issue or submit a pull request on GitHub.

## License

This package is licensed under the LGPL-3.0 License. See the [LICENSE](LICENSE) file for details.