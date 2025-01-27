package nject_test

import (
	"fmt"
	"reflect"

	"github.com/muir/nject/nject"
)

// ExampleGeneratedFromInjectionChain demonstrates how a special
// provider can be generated that builds types that are missing
// from an injection chain.
func ExampleGenerateFromInjectionChain() {
	type S struct {
		I int
	}
	fmt.Println(nject.Run("example",
		func() int {
			return 3
		},
		nject.GenerateFromInjectionChain(
			func(before nject.Collection, after nject.Collection) (nject.Provider, error) {
				full := before.Append("after", after)
				inputs, outputs := full.DownFlows()
				var n []interface{}
				for _, missing := range nject.ProvideRequireGap(outputs, inputs) {
					if missing.Kind() == reflect.Struct ||
						(missing.Kind() == reflect.Ptr &&
							missing.Elem().Kind() == reflect.Struct) {
						vp := reflect.New(missing)
						fmt.Println("Building filler for", missing)
						builder, err := nject.MakeStructBuilder(vp.Elem().Interface())
						if err != nil {
							return nil, err
						}
						n = append(n, builder)
					}
				}
				return nject.Sequence("build missing models", n...), nil
			}),
		func(s S, sp *S) {
			fmt.Println(s.I, sp.I)
		},
	))
	// Output: Building filler for nject_test.S
	// Building filler for *nject_test.S
	// 3 3
	// <nil>
}

func ExampleCollection_DownFlows_provider() {
	sequence := nject.Sequence("one provider", func(_ int, _ string) float64 { return 0 })
	inputs, outputs := sequence.DownFlows()
	fmt.Println("inputs", inputs)
	fmt.Println("outputs", outputs)
	// Output: inputs [int string]
	// outputs [float64]
}

func ExampleCollection_DownFlows_collection() {
	sequence := nject.Sequence("two providers",
		func(_ int, _ int64) float32 { return 0 },
		func(_ int, _ string) float64 { return 0 },
	)
	inputs, outputs := sequence.DownFlows()
	fmt.Println("inputs", inputs)
	fmt.Println("outputs", outputs)
	// Output: inputs [int int64 string]
	// outputs [float32 float64]
}

func ExampleCollection_ForEachProvider() {
	seq := nject.Sequence("example",
		func() int {
			return 10
		},
		func(_ int, _ string) {},
	)
	seq.ForEachProvider(func(p nject.Provider) {
		fmt.Println(p.DownFlows())
	})
	// Output: [] [int]
	// [int string] []
}
