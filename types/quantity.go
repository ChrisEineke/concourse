package types

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
)

// Quantity is a arbitrary-precision decimal used to represent numbers in
// Concourse. The unit suffix is multiplied with the number, but carried along
// so marshalling can be done without changing the string representation.
type Quantity struct {
	// Decimal is the expanded value .
	decimal.Decimal
	// Number is the numeric value extracted from the source string.
	Number string
	// UnitSuffix is the unit suffix extracted from the source string.
	UnitSuffix string
}

var (
	One = Quantity{decimal.NewFromInt(1), "1", ""}
	Two = Quantity{decimal.NewFromInt(2), "2", ""}
	Ten = Quantity{decimal.NewFromInt(10), "10", ""}

	Pico  = Quantity{Ten.Pow(decimal.NewFromInt(-12)), Ten.Pow(decimal.NewFromInt(-12)).String(), "p"}
	Nano  = Quantity{Ten.Pow(decimal.NewFromInt(-9)), Ten.Pow(decimal.NewFromInt(-9)).String(), "n"}
	Micro = Quantity{Ten.Pow(decimal.NewFromInt(-6)), Ten.Pow(decimal.NewFromInt(-6)).String(), "u"}
	Milli = Quantity{Ten.Pow(decimal.NewFromInt(-3)), Ten.Pow(decimal.NewFromInt(-3)).String(), "m"}
	Centi = Quantity{Ten.Pow(decimal.NewFromInt(-2)), Ten.Pow(decimal.NewFromInt(-2)).String(), "c"}
	Deci  = Quantity{Ten.Pow(decimal.NewFromInt(-1)), Ten.Pow(decimal.NewFromInt(-1)).String(), "d"}

	Kilo = Quantity{Ten.Pow(decimal.NewFromInt(3)), Ten.Pow(decimal.NewFromInt(3)).String(), "k"}
	Mega = Quantity{Ten.Pow(decimal.NewFromInt(6)), Ten.Pow(decimal.NewFromInt(6)).String(), "M"}
	Giga = Quantity{Ten.Pow(decimal.NewFromInt(9)), Ten.Pow(decimal.NewFromInt(9)).String(), "G"}
	Tera = Quantity{Ten.Pow(decimal.NewFromInt(12)), Ten.Pow(decimal.NewFromInt(12)).String(), "T"}
	Peta = Quantity{Ten.Pow(decimal.NewFromInt(15)), Ten.Pow(decimal.NewFromInt(15)).String(), "P"}

	Kibi = Quantity{Two.Pow(decimal.NewFromInt(10)), Two.Pow(decimal.NewFromInt(10)).String(), "Ki"}
	Mebi = Quantity{Two.Pow(decimal.NewFromInt(20)), Two.Pow(decimal.NewFromInt(20)).String(), "Mi"}
	Gibi = Quantity{Two.Pow(decimal.NewFromInt(30)), Two.Pow(decimal.NewFromInt(30)).String(), "Gi"}
	Tebi = Quantity{Two.Pow(decimal.NewFromInt(40)), Two.Pow(decimal.NewFromInt(40)).String(), "Ti"}
	Pebi = Quantity{Two.Pow(decimal.NewFromInt(50)), Two.Pow(decimal.NewFromInt(50)).String(), "Pi"}

	// Suffixes is the list of supported unit suffixes sorted by descending
	// string length, then by multiplier magnitude.
	Suffixes = []Quantity{
		Pebi,
		Tebi,
		Gibi,
		Mebi,
		Kibi,
		Peta,
		Tera,
		Giga,
		Mega,
		Kilo,
		Deci,
		Centi,
		Milli,
		Micro,
		Nano,
		Pico,
	}
	LabelsToSuffix = map[string]Quantity{
		Pebi.UnitSuffix:  Pebi,
		Gibi.UnitSuffix:  Gibi,
		Mebi.UnitSuffix:  Mebi,
		Kibi.UnitSuffix:  Kibi,
		Peta.UnitSuffix:  Peta,
		Tera.UnitSuffix:  Tera,
		Giga.UnitSuffix:  Giga,
		Mega.UnitSuffix:  Mega,
		Kilo.UnitSuffix:  Kilo,
		Deci.UnitSuffix:  Deci,
		Centi.UnitSuffix: Centi,
		Milli.UnitSuffix: Milli,
		Micro.UnitSuffix: Micro,
		Nano.UnitSuffix:  Nano,
		Pico.UnitSuffix:  Pico,
	}
)

func (r *Quantity) ParseString(val string) error {
	if val == "" {
		return fmt.Errorf("val is empty")
	}
	suffix := One
	for _, s := range Suffixes {
		if !strings.HasSuffix(val, s.UnitSuffix) {
			continue
		}
		suffix = s
	}
	num := strings.TrimSuffix(val, suffix.UnitSuffix)
	decNum, err := decimal.NewFromString(num)
	if err != nil {
		return err
	}
	r.Decimal = decNum.Mul(suffix.Decimal).Normalize()
	r.Number = num
	r.UnitSuffix = suffix.UnitSuffix
	return nil
}

func (r *Quantity) UnmarshalRJSON(data []byte) error {
	var valueStr string
	if err := json.Unmarshal(data, &valueStr); err != nil {
		return err
	}
	return r.ParseString(valueStr)
}

func (r *Quantity) MarshalRJSON() ([]byte, error) {
	return json.Marshal(fmt.Sprintf("%s%s", r.Number, r.UnitSuffix))
}
