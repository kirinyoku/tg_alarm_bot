package sources

type Fetcher interface {
	Fetch() ([]Message, error)
}

type Processor interface {
	Process(message Message) error
}

type Message struct {
	ID   string
	Text string
}
