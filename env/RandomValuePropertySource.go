package env

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-external-config/go/lang"
	"github.com/go-external-config/go/util/optional"
	"github.com/go-external-config/go/util/regex"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var keyPattern = regexp.MustCompile(`^random\.((?P<uuid>uuid)|(?P<string>string)\((?P<size>\d+)\)|(?P<value>value)(\((?P<bytes>\d+)\))?|(?P<int>int)(\((?P<max>\d+)\))?|(?P<int>int)(\((?P<min>-?\d+),(?P<max>\d+)\))?|(?P<int64>int64)(\((?P<max>\d+)\))?|(?P<int64>int64)(\((?P<min>-?\d+),(?P<max>\d+)\))?)$`)

// Custom property source as an additional logic for properties processing, like property=${random.uuid}
//
// ${random.value} - Random 64-bit hexadecimal value
//
// ${random.int} - Random 32-bit integer (full range)
//
// ${random.int(max)} - Random int in [0, max)
//
// ${random.int(min,max)} - Random int in [min, max)
//
// ${random.int64} - Random 64-bit integer (full range)
//
// ${random.int64(max)} - Random int64 in [0, max)
//
// ${random.int64(min,max)} - Random int64 in [min, max)
//
// ${random.uuid} - Random UUID
//
// ${random.string(length)} - Random alphanumeric string
type RandomValuePropertySource struct{}

func NewRandomValuePropertySource() *RandomValuePropertySource {
	return &RandomValuePropertySource{}
}

func (s *RandomValuePropertySource) Name() string {
	return "RandomValuePropertySource"
}

func (s *RandomValuePropertySource) HasProperty(key string) bool {
	return strings.HasPrefix(key, "random.")
}

func (s *RandomValuePropertySource) Property(key string) string {
	for _, m := range keyPattern.FindAllStringSubmatchIndex(key, -1) {
		match := regex.MatchOf(keyPattern, key, m)
		uuid := match.NamedGroup("uuid")
		if uuid.Present() {
			return s.RandomUuid()
		}
		str := match.NamedGroup("string")
		if str.Present() {
			size := optional.OfCommaErr(strconv.Atoi(match.NamedGroup("size").Value())).OrElsePanic("Cannot parse size %s", match.Expr())
			return s.RandomString(size)
		}
		value := match.NamedGroup("value")
		if value.Present() {
			bytes := optional.OfCommaErr(strconv.Atoi(match.NamedGroup("bytes").OrElse("8"))).OrElsePanic("Cannot parse bytes %s", match.Expr())
			return s.RandomValue(bytes)
		}
		intValue := match.NamedGroup("int")
		if intValue.Present() {
			max := match.NamedGroup("max")
			if max.Present() {
				min := optional.OfCommaErr(strconv.ParseInt(match.NamedGroup("min").OrElse("0"), 10, 32)).OrElsePanic("Cannot parse min %s", match.Expr())
				max := optional.OfCommaErr(strconv.ParseInt(max.Value(), 10, 32)).OrElsePanic("Cannot parse max %s", match.Expr())
				return fmt.Sprint(int32(s.RandomInt64Value(int64(min), int64(max))))
			} else {
				return fmt.Sprint(int32(s.RandomInt64Value(int64(math.MinInt), int64(math.MaxInt))))
			}
		}
		int64Value := match.NamedGroup("int64")
		if int64Value.Present() {
			max := match.NamedGroup("max")
			if max.Present() {
				min := optional.OfCommaErr(strconv.ParseInt(match.NamedGroup("min").OrElse("0"), 10, 64)).OrElsePanic("Cannot parse min %s", match.Expr())
				max := optional.OfCommaErr(strconv.ParseInt(max.Value(), 10, 64)).OrElsePanic("Cannot parse max %s", match.Expr())
				return fmt.Sprint(s.RandomInt64Value(min, max))
			} else {
				return fmt.Sprint(s.RandomInt64Value(math.MinInt64, math.MaxInt64))
			}
		}
	}
	panic("Cannot process " + key)
}

func (s *RandomValuePropertySource) RandomUuid() string {
	u := make([]byte, 16)
	optional.OfCommaErr(rand.Read(u)).OrElsePanic("Cannot generate random value")
	// Set version (4) and variant bits per RFC 4122
	u[6] = (u[6] & 0x0f) | 0x40 // Version 4 (0b0100xxxx)
	u[8] = (u[8] & 0x3f) | 0x80 // Variant 1 (0b10xxxxxx)
	// Format as canonical 8-4-4-4-12 hexadecimal UUID string
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:16])
}

func (s *RandomValuePropertySource) RandomValue(bytes int) string {
	buf := make([]byte, bytes)
	optional.OfCommaErr(rand.Read(buf)).OrElsePanic("Cannot generate random value")
	return fmt.Sprintf("%x", buf)
}

func (s *RandomValuePropertySource) RandomInt64Value(minInclusive, maxExclusive int64) int64 {
	lang.AssertState(maxExclusive > minInclusive, "Invalid range, [%d, %d)", minInclusive, maxExclusive)
	rng := new(big.Int).Sub(big.NewInt(maxExclusive), big.NewInt(minInclusive))
	n := optional.OfCommaErr(rand.Int(rand.Reader, rng)).OrElsePanic("Cannot generate random value")
	return n.Int64() + minInclusive
}

func (s *RandomValuePropertySource) RandomString(length int) string {
	bytes := make([]byte, length)
	optional.OfCommaErr(rand.Read(bytes)).OrElsePanic("Cannot generate random value")
	// Map random bytes to allowed characters
	for i, b := range bytes {
		bytes[i] = letters[int(b)%len(letters)]
	}
	return string(bytes)
}

func (s *RandomValuePropertySource) Properties() map[string]string {
	return nil
}
