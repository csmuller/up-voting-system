package types

type GenesisState struct {
	Params Params `json:"params"`
}

func NewGenesisState(params Params) GenesisState {
	return GenesisState{
		Params: params,
	}
}
