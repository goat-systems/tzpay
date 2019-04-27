package options

// Options is a struct to represent configuration options for payman
type Options struct {
	Delegate        string
	Secret          string
	Password        string
	Service         bool
	Cycle           int
	Node            string
	Port            string
	Fee             float32
	File            string
	NetworkFee      int
	NetworkGasLimit int
	Dry             bool
	RedditAgent     string
	RedditTitle     string
}
