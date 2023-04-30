package main

type Dummy struct {
	Name  string         `json:"name"`
	Limbs []string       `json:arms`
	Stats map[string]int `json:stats`
	Age   int            `json:"age"`
	Mood  float64        `json:"mood"`
}

func makeDummy() Dummy {
	return Dummy{
		Name:  "John",
		Limbs: []string{"left_arm", "left_arm", "left_leg", "right_leg"},
		Stats: map[string]int{"goodness": 1, "badness": 0},
		Age:   4,
		Mood:  0.65,
	}
}
