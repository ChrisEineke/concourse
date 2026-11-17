package types

import (
	"slices"
	"testing"

	"github.com/shopspring/decimal"
)

func TestQuantity(t *testing.T) {
	q := Quantity{Decimal: decimal.NewFromInt(10)}

	if q.IntPart() != 10 {
		t.Errorf("expected IntPart to be 10, got %d", q.IntPart())
	}

	q2 := Quantity{Decimal: decimal.NewFromFloat(1.5)}
	if !q2.Equal(decimal.NewFromFloat(1.5)) {
		t.Errorf("expected value to be 1.5, got %s", q2.String())
	}
}

func TestQuantityRJSONEncoding(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected decimal.Decimal
		wantErr  bool
	}{
		{
			name:     "empty string",
			input:    []byte(`""`),
			expected: decimal.Decimal{},
			wantErr:  true,
		},
		{
			name:     "plain number",
			input:    []byte(`"100"`),
			expected: decimal.NewFromInt(100),
			wantErr:  false,
		},
		{
			name:     "kilo prefix",
			input:    []byte(`"1.5k"`),
			expected: decimal.NewFromInt(1_500),
			wantErr:  false,
		},
		{
			name:     "mega prefix",
			input:    []byte(`"2M"`),
			expected: decimal.NewFromInt(2_000_000),
			wantErr:  false,
		},
		{
			name:     "giga prefix",
			input:    []byte(`"4G"`),
			expected: decimal.NewFromInt(4_000_000_000),
			wantErr:  false,
		},
		{
			name:     "tera prefix",
			input:    []byte(`"12T"`),
			expected: decimal.NewFromInt(12_000_000_000_000),
			wantErr:  false,
		},
		{
			name:     "peta prefix",
			input:    []byte(`"25P"`),
			expected: decimal.NewFromInt(25_000_000_000_000_000),
			wantErr:  false,
		},
		{
			name:     "kibi prefix",
			input:    []byte(`"1Ki"`),
			expected: decimal.NewFromInt(1_024),
			wantErr:  false,
		},
		{
			name:     "mebi prefix",
			input:    []byte(`"2Mi"`),
			expected: decimal.NewFromInt(2 * 1_024 * 1_024),
			wantErr:  false,
		},
		{
			name:     "gibi prefix",
			input:    []byte(`"4Gi"`),
			expected: decimal.NewFromInt(4 * 1_024 * 1_024 * 1_024),
			wantErr:  false,
		},
		{
			name:     "tebi prefix",
			input:    []byte(`"8Ti"`),
			expected: decimal.NewFromInt(8 * 1_024 * 1_024 * 1_024 * 1_024),
			wantErr:  false,
		},
		{
			name:     "pebi prefix",
			input:    []byte(`"1Pi"`),
			expected: decimal.NewFromInt(1125899906842624),
			wantErr:  false,
		},
		{
			name:     "invalid json",
			input:    []byte(`invalid`),
			expected: decimal.Zero,
			wantErr:  true,
		},
		{
			name:     "invalid number",
			input:    []byte(`"abc"`),
			expected: decimal.Zero,
			wantErr:  true,
		},
		{
			name:     "kubernetes 100 millicores",
			input:    []byte(`"100m"`),
			expected: decimal.NewFromInt(10).Pow(decimal.NewFromInt(-1)),
			wantErr:  false,
		},
		{
			name:     "kubernetes 1/4 Mibibytes memory",
			input:    []byte(`"0.25Mi"`),
			expected: decimal.NewFromInt(256).Mul(decimal.NewFromInt(1_024)),
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var q Quantity
			err := q.UnmarshalRJSON(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalRJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !q.Equal(tt.expected) {
				t.Errorf("UnmarshalRJSON() = %s, want %s", q.String(), tt.expected.String())
				return
			}
			if err != nil {
				return
			}

			input, err := q.MarshalRJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalRJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !slices.Equal(input, tt.input) {
				t.Errorf("MarshalRJSON() = %s, want %s", q.String(), tt.expected.String())
			}
		})
	}
}
