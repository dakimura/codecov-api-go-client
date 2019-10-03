# codecov-api-go-client
CodeCov API client in Go. 

See https://docs.codecov.io/reference


# Usage
```go
import (
	"fmt"
	"github.com/dakimura/codecov-api-go-client/codecovapi"
)

func main() {
	token := "YOUR_API_TOKEN_GOES_HERE"
	cli := codecovapi.NewClient(token, nil)

	resp, err := cli.Get(codecovapi.GitHub, "owner name of the repository", "repository name")
	// resp, err := cli.Get(codecovapi.GitHub, "dakimura", "readthrough") // e.g
	// resp, err := cli.GetBranch(codecovapi.GitHub, "dakimura", "readthrough", "master")

	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}
```