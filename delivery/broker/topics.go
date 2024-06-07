package broker

var (
	RawBlock        Topic = newTopic("raw_block")
	RawBlockResults Topic = newTopic("raw_block_results")
	RawGenesis      Topic = newTopic("raw_genesis")
	RawTransaction  Topic = newTopic("raw_transaction")

	rawTopics = Topics{RawBlock, RawTransaction, RawBlockResults, RawGenesis}

	// allTopics is the list of all topics.
	allTopics = func(tcs []Topics) []string {
		stringTopics := make([]string, 0)
		for _, t := range tcs {
			stringTopics = append(stringTopics, t.ToStringSlice()...)
		}
		return removeDuplicates(stringTopics)
	}([]Topics{rawTopics})
)

type (
	Topic  *string
	Topics []Topic
)

func newTopic(t string) *string { return &t }

func (ts Topics) ToStringSlice() []string {
	res := make([]string, len(ts))

	for i, t := range ts {
		if t == nil {
			panic("topic is nil")
		}

		res[i] = *t
	}

	return res
}
