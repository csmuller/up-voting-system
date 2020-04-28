package types

import (
	"fmt"
	"github.com/csmuller/uep-voting/crypto"
	"strings"
)

//--------------------------------------------------------------------------------------------------
// Voter Credentials

type QueryResVoterCredential struct {
	Credential  crypto.Int `json:"credential"`
	BlockHeight int64      `json:"block_height"`
}

type QueryResVoterCredentials []QueryResVoterCredential

func (credentials QueryResVoterCredentials) String() string {
	var str strings.Builder
	str.WriteString("QueryResVoterCredentials: {\n")
	for _, credential := range credentials {
		str.WriteString(fmt.Sprintf("%s,\n", credential.String()))
	}
	str.WriteString("}")
	return str.String()
}

func (credential QueryResVoterCredential) String() string {
	return fmt.Sprintf("%s posted at block height %d\n",
		credential.BlockHeight,
		credential.Credential.String())
}
