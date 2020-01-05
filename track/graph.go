package track

type ExchangeGraph struct {
	parties []Party
}

func (e ExchangeGraph) GetParty(id string) *Party {
	// might want to rewrite this since it could be O(1) if
	// you were sober and re-did some system architecture with
	// the party struct you alcoholic piece of shit(seriously use
	// a hash map)
	for _, party := range e.parties {
		if party.id == id {
			return &party
		}
	}
}
