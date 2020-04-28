package crypto

import (
	"encoding/json"
	"fmt"
	"github.com/tendermint/go-amino"
	"math/big"
)

// This file implements large prime-order subgroups of cyclic, multiplicative groups of integers
// modulo p for some prime p. The subgroups are named GStarModPrime (G_p) and their respective
// super groups ZStarModPrime (Z*_p).

import (
	"crypto/rand"
)

const (
	// Election parameters (security parameter and groups).

	// Security parameter used in the double discrete log proof system.
	SecurityParam = 80

	// Group orders and modulus for the two required groups G_p and G_q of the UEP protocol.
	// 160 bit, order q of G_q
	Q string = "1081119563825030427708677600856959359670713108783"
	// P = a * Q + 1, 1024 bit, order p of G_p and modulus of G_q
	P string = "132981118064499312972124229719551507064282251442693318094413647002876359530119444044769383265695686373097209253015503887096288112369989708235068428214124661556800389180762828009952422599372290980806417384771730325122099441368051976156139223257233269955912341167062173607119895128870594055324929155200165347329"
	// O = 981 * P + 1, 1034 bit, modulus of G_p
	O string = "130321495703209326712681745125160476922996606413839451732525374062818832339517055163873995600381772645635265067955193809354362350122589914070367059649842168325664381397147571449753374147384845161190289037076295718619657452540690936633016438792088604556794094343720930134977497226293182174218430572096162040382421"

	// 1023 bit, order q of G_q
	Q1 string = "62419754450729612647565739452383276575857601899739936725159851334944150841968063259516646199602063446032019699733384807429406029957259175802747488347623169473252390835604139741777023566843030652585465424928737851640689453666654947197163915037792214904944077385094485372296355878304667660119111336076574627993"
	// P1 = 2 * Q1 + 1, 1024 bit, order p of G_p and modulus of G_q
	P1 string = "124839508901459225295131478904766553151715203799479873450319702669888301683936126519033292399204126892064039399466769614858812059914518351605494976695246338946504781671208279483554047133686061305170930849857475703281378907333309894394327830075584429809888154770188970744592711756609335320238222672153149255987"
	// O1 = 325 * P1 + 1, 1024 bit, modulus of G_p
	O1 string = "40448000884072788995622599165144363221155726031031478997903583665043809745595304992166786737342137113028748765427233355214255107412303945920180372449259813818667549261471482552671511271314283862875381595353822127863166765975992405783762216944489355258403762145541226521248038609141424643757184145777620358939789"

	// 2047 bit, order q of G_q
	Q2 string = "16158503035655503650357438344334975980222051334857742016065172713762327569433945446598600705761456731844358980460949009747059779575245460547544076193224141560315438683650498045875098875194826053398028819192033784138396109321309878080919047169238085235290822926018152521443787945770532904303776199561965192760957166694834171210342487393282284747428088017663161029038902829665513096354230157075129296432088558362971801859230928678799175576150822952201848806616643615613562842355410104862578550863465661734839271290328348967522998634176499319107762583194718667771801067716614802322659239302476074096777926805529797144183"
	// P2 = 2 * Q2 + 1, 2048 bit, order p of G_p and modulus of G_q
	P2 string = "32317006071311007300714876688669951960444102669715484032130345427524655138867890893197201411522913463688717960921898019494119559150490921095088152386448283120630877367300996091750197750389652106796057638384067568276792218642619756161838094338476170470581645852036305042887575891541065808607552399123930385521914333389668342420684974786564569494856176035326322058077805659331026192708460314150258592864177116725943603718461857357598351152301645904403697613233287231227125684710820209725157101726931323469678542580656697935045997268352998638215525166389437335543602135433229604645318478604952148193555853611059594288367"
	// O2 = 3157 * P2 + 1, 2048+ bit, modulus of G_p
	O2 string = "101992471161057539041056150829442368387161588025622067605403370169267811618267063658930367654766314891401593884669510149523441328678949346976098208931630781528711048971201943665563624100229742049048357906740117245481556242036107950446761025732230794005155674309026578715353189513703603691965435371635124296707161636177793288679681780426397781325766091567489872415293554660848718664187900751458216119079342980387078013335465621820580396236663994474298069667364254501752808660947348581892595813050195256870305480384552538683005167378922063702208197425125064230975608339427272632260625118477228979698862273996504079574086253"
)

//----------------------------------------------------------------------------------------------
// GStarModPrime

// GStarModPrime represents the concept of a prime-order subgroup of the multiplicative
// group of integers modulo a prime. The decisional Diffie-Hellman assumption is believed to hold
// in such a group.
type GStarModPrime struct {
	Modulus *big.Int `json:"mod"` // The modulus of the group. A prime number.
	Order   *big.Int `json:"ord"` // The order of the group. A prime number.
}

// NewGStarModPrime creates a new instance of the group with the given order and modulus. Both
// arguments must be prime numbers.
func NewGStarModPrime(modulus *big.Int, order *big.Int) GStarModPrime {
	return GStarModPrime{
		Modulus: modulus,
		Order:   order,
	}
}

// ZModOrder returns the ring of integers modulo this group's order.
func (g *GStarModPrime) ZModOrder() ZModPrime {
	return NewZModPrime(g.Order)
}

// ZStartModModulus returns the multiplicative group of integers modulo this group's modulus, i.e.
// the multiplicative group (Z*_p) of which this group (G*_q) is a subgroup.
func (g *GStarModPrime) ZStarModModulus() ZStarModPrime {
	return NewZStarModPrime(g.Modulus)
}

// ZStarModOrder returns multiplicative group of integers modulo this group's order.
func (g *GStarModPrime) ZStarModOrder() ZStarModPrime {
	return NewZStarModPrime(g.Order)
}

// RandomGenerator gets a randomly selected generator for this group. Implemented according to
// Appendix A.2.1 of FIPS 186-4
func (g *GStarModPrime) RandomGenerator() *big.Int {
	cofactor := g.Cofactor()
	one := big.NewInt(1)
	generator := big.NewInt(1)
	for generator.Cmp(one) == 0 {
		// Get a random element of the multiplicative group of integers with the same modulus.
		rndElem := g.ZStarModModulus().RandomElement()
		generator.Exp(rndElem, cofactor, g.Modulus)
	}
	return generator
}

// DefaultGenerator gets the default generator for this group. Implemented according to
// http://en.wikipedia.org/wiki/Schnorr_group
func (g *GStarModPrime) DefaultGenerator() *big.Int {
	cofactor := g.Cofactor()
	one := big.NewInt(1)
	generator := big.NewInt(1)
	h := big.NewInt(1)
	for generator.Cmp(one) == 0 {
		generator.Exp(h, cofactor, g.Modulus)
		h.Add(h, one)
	}
	return generator
}

// IdentityElement returns the identity element of this group which is always 1.
func (g *GStarModPrime) IdentityElement() *big.Int {
	return big.NewInt(1)
}

// Cofactor calculates and returns the cofactor of this group. The cofactor r is calculated as
// r = (p-1)/q. Cf. the construction of Schnorr groups.
func (g *GStarModPrime) Cofactor() *big.Int {
	return new(big.Int).Div(g.ZStarModModulus().Order(), g.Order)
}

// Contains checks if the given value is an element of this group. The value v is an element of this
// group if 1 < v < p and v^q = 1 (mod p).
func (g *GStarModPrime) Contains(v *big.Int) bool {
	return v.Sign() > 0 &&
		v.Cmp(g.Modulus) < 0 &&
		g.Exp(v, g.Order).Cmp(big.NewInt(1)) == 0
}

// RandomElement gets a random element of this group.
func (g *GStarModPrime) RandomElement() *big.Int {
	rndElem := g.ZStarModModulus().RandomElement()
	return new(big.Int).Exp(rndElem, g.Cofactor(), g.Modulus)
}

// Mul computes x * y mod this group's modulus.
func (g *GStarModPrime) Mul(x, y *big.Int) *big.Int {
	r := new(big.Int)
	r.Mul(x, y)
	return r.Mod(r, g.Modulus)
}

// Exp computes base^exp mod this group's modulus.
func (g *GStarModPrime) Exp(base, exp *big.Int) *big.Int {
	return new(big.Int).Exp(base, exp, g.Modulus)
}

// Invert computes and returns the multiplicative inverse of the given group element.
func (g *GStarModPrime) Invert(i *big.Int) *big.Int {
	return new(big.Int).ModInverse(i, g.Modulus)
}

// RandomBits fetches uniform random random bits with a maximum of the given bit length. If 'exact'
// is true, exactly bitlen number of bits are set. The returned byte array is in big-endian order,
// ready for usage with big.Int.
// Adapted from https://github.com/dedis/kyber/blob/master/util/random/rand.go
func RandomBits(bitlen uint, exact bool) []byte {
	b := make([]byte, (bitlen+7)/8)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	highbits := bitlen & 7
	if highbits != 0 {
		b[0] &= ^(0xff << highbits)
	}
	if exact {
		if highbits != 0 {
			b[0] |= 1 << (highbits - 1)
		} else {
			b[0] |= 0x80
		}
	}
	return b
}

// RandInt chooses a uniform random big.Int less than a given number.
// Taken and adapted from https://github.com/dedis/kyber/blob/master/util/random/rand.go
func RandInt(exclMax *big.Int) *big.Int {
	bitlen := uint(exclMax.BitLen())
	i := new(big.Int)
	for {
		i.SetBytes(RandomBits(bitlen, false))
		if i.Sign() > 0 && i.Cmp(exclMax) < 0 {
			return i
		}
	}
}

// RandInt chooses a uniform random big.Int less than a given number.
// Taken and adapted from https://github.com/dedis/kyber/blob/master/util/random/rand.go
func RandIntInRange(min, exclMax *big.Int) *big.Int {
	bitlen := uint(exclMax.BitLen())
	i := new(big.Int)
	for {
		i.SetBytes(RandomBits(bitlen, false))
		if i.Cmp(min) >= 0 && i.Cmp(exclMax) < 0 {
			return i
		}
	}
}

// String returns a string representation of this group.
func (g *GStarModPrime) String() string {
	return fmt.Sprintf("GStarModPrime: {\n\tmodulus=%s,\n\t order=%s\n}", g.Modulus.String(),
		g.Order.String())
}

type gStarModPrimeDTO struct {
	Modulus Int `json:"mod"`
	Order   Int `json:"ord"`
}

func (g GStarModPrime) MarshalAmino() (string, error) {
	dto := gStarModPrimeDTO{NewInt(g.Modulus), NewInt(g.Order)}
	return string(amino.MustMarshalBinaryBare(dto)), nil
}

func (g *GStarModPrime) UnmarshalAmino(bytes []byte) error {
	var dto gStarModPrimeDTO
	amino.MustUnmarshalBinaryBare(bytes, &dto)
	g.Modulus = dto.Modulus.BigInt()
	g.Order = dto.Order.BigInt()
	return nil
}

func (g GStarModPrime) MarshalJSON() ([]byte, error) {
	dto := gStarModPrimeDTO{NewInt(g.Modulus), NewInt(g.Order)}
	return json.Marshal(dto)
}

func (g *GStarModPrime) UnmarshalJSON(bytes []byte) error {
	var dto gStarModPrimeDTO
	err := amino.UnmarshalJSON(bytes, &dto)
	if err != nil {
		return err
	}
	g.Modulus = dto.Modulus.BigInt()
	g.Order = dto.Order.BigInt()
	return nil
}

//----------------------------------------------------------------------------------------------
// ZModPrime

// ZModPrime represents the ring of integers modulo a prime p (abbreviated as Z_p). Addition,
// subtraction and multiplication are performed modulo p.
type ZModPrime struct {
	Modulus *big.Int // the rings's modulus p. A prime number.
}

// NewZModPrime constructs a new instance of a ring using p as its modulus.
func NewZModPrime(p *big.Int) ZModPrime {
	return ZModPrime{Modulus: p}
}

// Contains checks if the given integer is an element of this ring.
func (g ZModPrime) Contains(n *big.Int) bool {
	return n.Sign() >= 0 && n.Cmp(g.Modulus) < 0
}

// Order gets the order of this ring, i.e. p - 1 because p is prime.
func (g ZModPrime) Order() *big.Int {
	return new(big.Int).Sub(g.Modulus, big.NewInt(1))
}

// RandomElement gets a random element of this ring.
func (g ZModPrime) RandomElement() *big.Int {
	return RandInt(g.Modulus)
}

// AdditiveIdentity returns the additive identity element of this ring, which is always 0.
func (g ZModPrime) AdditiveIdentity() *big.Int {
	return big.NewInt(0)
}

// MultiplicativeIdentity returns the multiplicative identity element of this ring, which is
// always 1.
func (g ZModPrime) MultiplicativeIdentity() *big.Int {
	return big.NewInt(1)
}

// Add computes x + y mod this ring's modulus.
func (g ZModPrime) Add(x, y *big.Int) *big.Int {
	r := new(big.Int)
	r.Add(x, y)
	r.Mod(r, g.Modulus)
	return r
}

// Mul computes x * y mod this ring's modulus.
func (g ZModPrime) Mul(x, y *big.Int) *big.Int {
	r := new(big.Int)
	r.Mul(x, y)
	r.Mod(r, g.Modulus)
	return r
}

// AdditiveInvert computes the additive inverse of the given element.
func (g ZModPrime) AdditiveInvert(i *big.Int) *big.Int {
	tmp := new(big.Int).Sub(g.Modulus, i)
	return tmp.Mod(tmp, g.Modulus)
}

// Exp computes base^exp mod this ring's modulus.
func (g ZModPrime) Exp(base, exp *big.Int) *big.Int {
	return new(big.Int).Exp(base, exp, g.Modulus)
}

func (g ZModPrime) MarshalAmino() (string, error) {
	return string(amino.MustMarshalBinaryBare(NewInt(g.Modulus))), nil
}

func (g *ZModPrime) UnmarshalAmino(bytes []byte) error {
	var i Int
	amino.MustUnmarshalBinaryBare(bytes, &i)
	g.Modulus = i.BigInt()
	return nil
}

func (g ZModPrime) MarshalJSON() ([]byte, error) {
	mod := NewInt(g.Modulus)
	return json.Marshal(mod)
}

func (g *ZModPrime) UnmarshalJSON(bytes []byte) error {
	var mod Int
	err := amino.UnmarshalJSON(bytes, &mod)
	if err != nil {
		return err
	}
	g.Modulus = mod.BigInt()
	return nil
}

//----------------------------------------------------------------------------------------------
// ZStarModPrime

// ZStarModPrime represents the multiplicative group of integers modulo prime p, abbreviated as
// Z*_p. The group's order is equal to p-1 because p is prime.
type ZStarModPrime struct {
	Modulus *big.Int // the group's modulus p. A prime number.
}

// NewZStarModPrime creates a new instance of this group with the given modulus p which must be
// a prime number.
func NewZStarModPrime(p *big.Int) ZStarModPrime {
	return ZStarModPrime{Modulus: p}
}

// Order gets the order of this group, i.e. p - 1 because p is prime.
func (g ZStarModPrime) Order() *big.Int {
	return new(big.Int).Sub(g.Modulus, big.NewInt(1))
}

// RandomElement gets a random element of this group.
func (g ZStarModPrime) RandomElement() *big.Int {
	return RandIntInRange(big.NewInt(1), g.Modulus)
}

// IdentityElement returns the identity element of this group which is always 1.
func (g ZStarModPrime) IdentityElement() *big.Int {
	return big.NewInt(1)
}

// Mul computes x * y mod this group's modulus.
func (g ZStarModPrime) Mul(x, y *big.Int) *big.Int {
	r := new(big.Int)
	r.Mul(x, y)
	r.Mod(r, g.Modulus)
	return r
}

// Exp computes base^exp mod this group's modulus.
func (g ZStarModPrime) Exp(base, exp *big.Int) *big.Int {
	return new(big.Int).Exp(base, exp, g.Modulus)
}

func (g ZStarModPrime) MarshalAmino() (string, error) {
	return string(amino.MustMarshalBinaryBare(NewInt(g.Modulus))), nil
}

func (g *ZStarModPrime) UnmarshalAmino(bytes []byte) error {
	var i Int
	amino.MustUnmarshalBinaryBare(bytes, &i)
	g.Modulus = i.BigInt()
	return nil
}

func (g ZStarModPrime) MarshalJSON() ([]byte, error) {
	mod := NewInt(g.Modulus)
	return json.Marshal(mod)
}

func (g *ZStarModPrime) UnmarshalJSON(bytes []byte) error {
	var mod Int
	err := amino.UnmarshalJSON(bytes, &mod)
	if err != nil {
		return err
	}
	g.Modulus = mod.BigInt()
	return nil
}
