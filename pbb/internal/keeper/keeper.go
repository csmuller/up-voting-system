package keeper

import (
	"encoding/binary"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params/subspace"
	"github.com/csmuller/up-voting-system/crypto"
	"github.com/csmuller/up-voting-system/pbb/internal/types"
	"math/big"
)

// This key is used for the credential polynomial in the same store where the credentials are
// stored.
var polynomialKey = make([]byte, 8)

// BulletinBoardKeeper maintains the link to storage and exposes getter/setter methods for the various parts of
// the state machine
type BulletinBoardKeeper struct {
	credentialStoreKey sdk.StoreKey
	ballotStoreKey     sdk.StoreKey
	polynomialStoreKey sdk.StoreKey
	cdc                *codec.Codec // The wire codec for binary encoding/decoding.
	paramStore         subspace.Subspace
}

// NewBulletinBoardKeeper creates new instances of the pbb BulletinBoardKeeper
func NewBulletinBoardKeeper(credentialStoreKey sdk.StoreKey, ballotStoreKey sdk.StoreKey,
	polyStoreKey sdk.StoreKey, cdc *codec.Codec, paramStore subspace.Subspace) BulletinBoardKeeper {

	return BulletinBoardKeeper{
		credentialStoreKey: credentialStoreKey,
		ballotStoreKey:     ballotStoreKey,
		polynomialStoreKey: polyStoreKey,
		cdc:                cdc,
		paramStore:         paramStore.WithKeyTable(types.ParamKeyTable()),
	}
}

func (k BulletinBoardKeeper) GetBallot(ctx sdk.Context, electionCredential big.Int) *types.Ballot {
	store := ctx.KVStore(k.ballotStoreKey)
	if !store.Has(electionCredential.Bytes()) {
		return nil
	}
	var ballot *types.Ballot
	k.cdc.MustUnmarshalBinaryBare(store.Get(electionCredential.Bytes()), ballot)
	return ballot
}

func (k BulletinBoardKeeper) GetBallotsIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.ballotStoreKey)
	return sdk.KVStorePrefixIterator(store, nil)
}

func (k BulletinBoardKeeper) HasElectionCredential(ctx sdk.Context, uHat *big.Int) bool {
	store := ctx.KVStore(k.ballotStoreKey)
	return store.Has(uHat.Bytes())
}

func (k BulletinBoardKeeper) StoreBallot(ctx sdk.Context, b types.Ballot) error {
	store := ctx.KVStore(k.ballotStoreKey)
	if store.Has(b.UHat.BigInt().Bytes()) {
		return fmt.Errorf("a ballot has already been stored for voter with election credential %s",
			b.UHat.String())
	}
	//TODO: Using the the credential bytes directly as key, but might be better to use something
	// shorter. e.g. a hash of it.
	store.Set(b.UHat.BigInt().Bytes(), k.cdc.MustMarshalBinaryBare(b))
	return nil
}

func (k BulletinBoardKeeper) GetVoterCredentialsIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.credentialStoreKey)
	return sdk.KVStorePrefixIterator(store, nil)
}

func (k BulletinBoardKeeper) HasVoterCredential(ctx sdk.Context, credential big.Int) bool {
	store := ctx.KVStore(k.credentialStoreKey)
	credentialBytes := k.cdc.MustMarshalBinaryBare(credential)
	return store.Has(credentialBytes)
}

// StoreVoterCredential stores the given voter credential in the credentials KV store and updates
// the credentials polynomial, i.e includes the credential in the polynomial. Throws an error if the
// credential is already in the store.
func (k BulletinBoardKeeper) StoreVoterCredential(ctx sdk.Context, credential crypto.Int) error {
	store := ctx.KVStore(k.credentialStoreKey)
	credentialBytes := k.cdc.MustMarshalBinaryBare(credential)
	if store.Has(credentialBytes) {
		return fmt.Errorf("the credential %s is already set", credential.String())
	}
	// Save current block height with the credential.
	b := make([]byte, 8)
	binary.PutVarint(b, ctx.BlockHeight())
	store.Set(credentialBytes, b)

	// Update credential polynomial
	k.includeCredentialInPolynomial(ctx, credential)
	return nil
}

func (k BulletinBoardKeeper) includeCredentialInPolynomial(ctx sdk.Context, credential crypto.Int) {
	store := ctx.KVStore(k.polynomialStoreKey)
	poly := k.GetCredentialPolynomial(ctx)
	newPoly := poly.IncludeCredential(credential.BigInt())
	store.Set(polynomialKey, k.cdc.MustMarshalBinaryBare(newPoly))
}

func (k BulletinBoardKeeper) GetCredentialPolynomial(ctx sdk.Context) crypto.Polynomial {
	store := ctx.KVStore(k.polynomialStoreKey)
	if store.Has(polynomialKey) {
		polyBytes := store.Get(polynomialKey)
		var poly crypto.Polynomial
		k.cdc.MustUnmarshalBinaryBare(polyBytes, &poly)
		return poly
	} else {
		// Return the initial polynomial as long as no one has registered yet.
		params := k.GetParams(ctx)
		gP := params.CommP.G
		coeffs := []*big.Int{big.NewInt(1)}
		return crypto.NewPolynomial(coeffs, gP.ZModOrder())
	}
}

// SetParams sets the auth module's parameters.
func (k BulletinBoardKeeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramStore.SetParamSet(ctx, &params)
}

// GetParams gets the bulletin board module's parameters.
func (k BulletinBoardKeeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramStore.GetParamSet(ctx, &params)
	return
}
