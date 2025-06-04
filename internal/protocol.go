package internal

type MsgKind int

const (
	ChallengeKind MsgKind = iota + 1
	ClientResponseKind
	QuoteResponseKind
	EchoKind
)

type Msg struct {
	Kind MsgKind

	Challenge     *Challenge
	Response      *Response
	QuoteResponse *QuoteResponse

	Echo *Echo
}

type Challenge struct {
	Difficulty int
	Challenge  []byte
}

type Response struct {
	Nonce []byte
}

type QuoteResponse struct {
	Text string
}

type Echo struct {
	Text string
}
