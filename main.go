package main

type redirectSpec struct {
	prefix, target string
}

var (
	listener    = "localhost:8080"
	redirectors = []redirectSpec{
		{
			"", "http://localhost:8000",
		},
		{
			"git", "http://localhost:3000",
		},
	}
)

func main() {

}
