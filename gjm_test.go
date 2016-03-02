package gjm

import (
	"fmt"
	"reflect"
	"testing"
)

type MapTest struct {
	in        map[string]interface{}
	path      string
	separator string
	value     interface{}
	out       interface{}
	err       error
}

func TestUpdateProperty(t *testing.T) {
	cases := []MapTest{
		{
			in:        setupDocument(),
			path:      "one",
			value:     "updated value",
			separator: ".",
			out: map[string]interface{}{
				"one": "updated value",
			},
			err: nil,
		},
		{
			in:        setupDocument(),
			path:      "one.three",
			value:     "updated value",
			separator: ".",
			out: map[string]interface{}{
				"one": map[string]interface{}{
					"three": "updated value",
					"two": map[string]interface{}{
						"three": []int{
							1, 2, 3,
						},
					},
					"four": map[string]interface{}{
						"five": []int{
							11, 22, 33,
						},
					},
				},
			},
			err: nil,
		},
		{
			in:        setupDocument(),
			path:      "one.two.three[3]",
			value:     "updated value",
			separator: ".",
			out: map[string]interface{}{
				"one": map[string]interface{}{
					"two": map[string]interface{}{
						"three": []interface{}{
							1, 2, 3, "updated value",
						},
					},
					"four": map[string]interface{}{
						"five": []int{
							11, 22, 33,
						},
					},
				},
			},
			err: nil,
		},
		{
			in:        setupDocument(),
			path:      "one.two.three[2]",
			value:     "updated value",
			separator: ".",
			out: map[string]interface{}{
				"one": map[string]interface{}{
					"two": map[string]interface{}{
						"three": []interface{}{
							1, 2, "updated value",
						},
					},
					"four": map[string]interface{}{
						"five": []int{
							11, 22, 33,
						},
					},
				},
			},
			err: nil,
		},
		{
			in:        setupDocument(),
			path:      "one.two.three[1]",
			value:     "updated value",
			separator: ".",
			out: map[string]interface{}{
				"one": map[string]interface{}{
					"two": map[string]interface{}{
						"three": []interface{}{
							1, "updated value", 3,
						},
					},
					"four": map[string]interface{}{
						"five": []int{
							11, 22, 33,
						},
					},
				},
			},
			err: nil,
		},
		{
			in:        setupDocument(),
			path:      "one.two.three",
			value:     "updated value",
			separator: ".",
			out: map[string]interface{}{
				"one": map[string]interface{}{
					"two": map[string]interface{}{
						"three": "updated value",
					},
					"four": map[string]interface{}{
						"five": []int{
							11, 22, 33,
						},
					},
				},
			},
			err: nil,
		},
		{
			in:        setupDocument_II(),
			path:      "one[3].three[0].four.nine",
			value:     "updated value",
			separator: ".",
			out: map[string]interface{}{
				"one": []map[string]interface{}{
					{
						"two": []map[string]interface{}{
							{"three": "got three"},
							{"four": "got four"},
						},
					},
					{
						"two": []map[string]interface{}{
							{"five": "got five"},
							{"six": "got six"},
						},
					},
					{
						"two": []map[string]interface{}{
							{"seven": "got seven"},
							{"eight": "got eight"},
						},
					},
					{
						"three": []map[string]interface{}{
							{"four": map[string]interface{}{
								"five": "six",
								"nine": "updated value",
							}},
							{"seven": map[string]interface{}{
								"eight": "ten",
							}},
						},
					},
				},
			},
			err: nil,
		},
	}

	num_cases := len(cases)
	for i, c := range cases {
		case_index := i + 1

		err_case := UpdateProperty(c.in, c.path, c.separator, c.value)
		out := c.in
		if !reflect.DeepEqual(c.err, err_case) {
			t.Errorf("\n[%d of %d: Errors should equal] \n\t%v \n \n\t%v", case_index, num_cases, err_case, c.err)
		}
		if !reflect.DeepEqual(out, c.out) {
			t.Errorf("\n[%d of %d: Results should equal] \n\t%v \n \n\t%v", case_index, num_cases, out, c.out)
		}
	}
}

func TestAddProperty(t *testing.T) {
	cases := []MapTest{
		{
			in:        setupDocument(),
			path:      "added",
			value:     "added value",
			separator: ".",
			out: map[string]interface{}{
				"one": map[string]interface{}{
					"two": map[string]interface{}{
						"three": []int{
							1, 2, 3,
						},
					},
					"four": map[string]interface{}{
						"five": []int{
							11, 22, 33,
						},
					},
				},
				"added": "added value",
			},
			err: nil,
		},
		{
			in:        setupDocument(),
			path:      "one.three",
			value:     "added value",
			separator: ".",
			out: map[string]interface{}{
				"one": map[string]interface{}{
					"three": "added value",
					"two": map[string]interface{}{
						"three": []int{
							1, 2, 3,
						},
					},
					"four": map[string]interface{}{
						"five": []int{
							11, 22, 33,
						},
					},
				},
			},
			err: nil,
		},
		{
			in:        setupDocument(),
			path:      "one.two.three[3]",
			value:     "added value",
			separator: ".",
			out: map[string]interface{}{
				"one": map[string]interface{}{
					"two": map[string]interface{}{
						"three": []interface{}{
							1, 2, 3, "added value",
						},
					},
					"four": map[string]interface{}{
						"five": []int{
							11, 22, 33,
						},
					},
				},
			},
			err: nil,
		},
		{
			in:        setupDocument_II(),
			path:      "one[3].three[0].four.nine",
			value:     "added value",
			separator: ".",
			out: map[string]interface{}{
				"one": []map[string]interface{}{
					{
						"two": []map[string]interface{}{
							{"three": "got three"},
							{"four": "got four"},
						},
					},
					{
						"two": []map[string]interface{}{
							{"five": "got five"},
							{"six": "got six"},
						},
					},
					{
						"two": []map[string]interface{}{
							{"seven": "got seven"},
							{"eight": "got eight"},
						},
					},
					{
						"three": []map[string]interface{}{
							{"four": map[string]interface{}{
								"five": "six",
								"nine": "added value",
							}},
							{"seven": map[string]interface{}{
								"eight": "ten",
							}},
						},
					},
				},
			},
			err: nil,
		},
		{
			in:        setupDocument(),
			path:      "one",
			value:     "added value",
			separator: ".",
			out:       setupDocument(),
			err:       fmt.Errorf("Property one already exists"),
		},
		{
			in:        setupDocument(),
			path:      "one.two.three",
			value:     "added value",
			separator: ".",
			out:       setupDocument(),
			err:       fmt.Errorf("Property one.two.three already exists"),
		},
		{
			in:        setupDocument(),
			path:      "one.two.three[1]",
			value:     "added value",
			separator: ".",
			out:       setupDocument(),
			err:       fmt.Errorf("Property one.two.three[1] already exists"),
		},
	}

	num_cases := len(cases)
	for i, c := range cases {
		case_index := i + 1

		err_case := AddProperty(c.in, c.path, c.separator, c.value)
		out := c.in
		if !reflect.DeepEqual(c.err, err_case) {
			t.Errorf("\n[%d of %d: Errors should equal] \n\t%v \n \n\t%v", case_index, num_cases, err_case, c.err)
		}
		if !reflect.DeepEqual(out, c.out) {
			t.Errorf("\n[%d of %d: Results should equal] \n\t%v \n \n\t%v", case_index, num_cases, out, c.out)
		}
	}
}

func TestDeleteProperty(t *testing.T) {
	cases := []MapTest{
		{
			in:        setupDocument(),
			path:      ".",
			separator: ".",
			out:       map[string]interface{}{},
			err:       nil,
		},
		{
			in:        setupDocument(),
			path:      "one",
			separator: ".",
			out:       map[string]interface{}{},
			err:       nil,
		},
		{
			in:        setupDocument(),
			path:      "one.two",
			separator: ".",
			out: map[string]interface{}{
				"one": map[string]interface{}{
					"four": map[string]interface{}{
						"five": []int{
							11, 22, 33,
						},
					},
				},
			},
			err: nil,
		},
		{
			in:        setupDocument(),
			path:      "one.two.three",
			separator: ".",
			out: map[string]interface{}{
				"one": map[string]interface{}{
					"two": map[string]interface{}{},
					"four": map[string]interface{}{
						"five": []int{
							11, 22, 33,
						},
					},
				},
			},
			err: nil,
		},
		{
			in:        setupDocument_I(),
			path:      "one[0]",
			separator: ".",
			out: map[string]interface{}{
				"one": []interface{}{
					map[string]interface{}{"map_b": []int{4, 5, 6}},
					map[string]interface{}{"map_c": []int{7, 8, 9}},
				},
			},
			err: nil,
		},
		{
			in:        setupDocument_I(),
			path:      "one[1]",
			separator: ".",
			out: map[string]interface{}{
				"one": []interface{}{
					map[string]interface{}{"map_a": []int{1, 2, 3}},
					map[string]interface{}{"map_c": []int{7, 8, 9}},
				},
			},
			err: nil,
		},
		{
			in:        setupDocument_II(),
			path:      "one[2].two[0]",
			separator: ".",
			out: map[string]interface{}{
				"one": []map[string]interface{}{
					{
						"two": []map[string]interface{}{
							{"three": "got three"},
							{"four": "got four"},
						},
					},
					{
						"two": []map[string]interface{}{
							{"five": "got five"},
							{"six": "got six"},
						},
					},
					{
						"two": []interface{}{
							map[string]interface{}{"eight": "got eight"},
						},
					},
					{
						"three": []map[string]interface{}{
							{"four": map[string]interface{}{
								"five": "six",
							}},
							{"seven": map[string]interface{}{
								"eight": "ten",
							}},
						},
					},
				},
			},
			err: nil,
		},
		{
			in:        setupDocument_II(),
			path:      "one[2].two[1]",
			separator: ".",
			out: map[string]interface{}{
				"one": []map[string]interface{}{
					{
						"two": []map[string]interface{}{
							{"three": "got three"},
							{"four": "got four"},
						},
					},
					{
						"two": []map[string]interface{}{
							{"five": "got five"},
							{"six": "got six"},
						},
					},
					{
						"two": []interface{}{
							map[string]interface{}{"seven": "got seven"},
						},
					},
					{
						"three": []map[string]interface{}{
							{"four": map[string]interface{}{
								"five": "six",
							}},
							{"seven": map[string]interface{}{
								"eight": "ten",
							}},
						},
					},
				},
			},
			err: nil,
		},
		{
			in:        setupDocument_II(),
			path:      "one[2].two[1].eight",
			separator: ".",
			out: map[string]interface{}{
				"one": []map[string]interface{}{
					{
						"two": []map[string]interface{}{
							{"three": "got three"},
							{"four": "got four"},
						},
					},
					{
						"two": []map[string]interface{}{
							{"five": "got five"},
							{"six": "got six"},
						},
					},
					{
						"two": []interface{}{
							map[string]interface{}{"seven": "got seven"},
						},
					},
					{
						"three": []map[string]interface{}{
							{"four": map[string]interface{}{
								"five": "six",
							}},
							{"seven": map[string]interface{}{
								"eight": "ten",
							}},
						},
					},
				},
			},
			err: nil,
		},
		{
			in:        setupDocument_II(),
			path:      "one[3].three[1].seven.eight",
			separator: ".",
			out: map[string]interface{}{
				"one": []map[string]interface{}{
					{
						"two": []map[string]interface{}{
							{"three": "got three"},
							{"four": "got four"},
						},
					},
					{
						"two": []map[string]interface{}{
							{"five": "got five"},
							{"six": "got six"},
						},
					},
					{
						"two": []map[string]interface{}{
							{"seven": "got seven"},
							{"eight": "got eight"},
						},
					},
					{
						"three": []map[string]interface{}{
							{"four": map[string]interface{}{
								"five": "six",
							}},
							{"seven": map[string]interface{}{}},
						},
					},
				},
			},
			err: nil,
		},
	}

	num_cases := len(cases)
	for i, c := range cases {
		case_index := i + 1

		err_case := DeleteProperty(c.in, c.path, c.separator)
		out := c.in
		if !reflect.DeepEqual(c.err, err_case) {
			t.Errorf("\n[%d of %d: Errors should equal] \n\t%v \n \n\t%v", case_index, num_cases, err_case, c.err)
		}
		if !reflect.DeepEqual(out, c.out) {
			t.Errorf("\n[%d of %d: Results should equal] \n\t%v \n \n\t%v", case_index, num_cases, out, c.out)
		}
	}
}

func TestGetProperty(t *testing.T) {
	cases := []MapTest{
		{
			in:        setupDocument(),
			path:      ".",
			separator: ".",
			out:       setupDocument(),
			err:       nil,
		},
		{
			in:        setupDocument(),
			path:      "one",
			separator: ".",
			out:       setupDocument()["one"].(map[string]interface{}),
			err:       nil,
		},
		{
			in:        setupDocument(),
			path:      "one.two",
			separator: ".",
			out:       setupDocument()["one"].(map[string]interface{})["two"],
			err:       nil,
		},
		{
			in:        setupDocument(),
			path:      "one.two.three",
			separator: ".",
			out:       setupDocument()["one"].(map[string]interface{})["two"].(map[string]interface{})["three"],
			err:       nil,
		},
		{
			in:        setupDocument(),
			path:      "one.two.three[0]",
			separator: ".",
			out:       setupDocument()["one"].(map[string]interface{})["two"].(map[string]interface{})["three"].([]int)[0],
			err:       nil,
		},
		{
			in:        setupDocument(),
			path:      "one.two.three[1]",
			separator: ".",
			out:       setupDocument()["one"].(map[string]interface{})["two"].(map[string]interface{})["three"].([]int)[1],
			err:       nil,
		},
		{
			in:        setupDocument(),
			path:      "one.two.three[2]",
			separator: ".",
			out:       setupDocument()["one"].(map[string]interface{})["two"].(map[string]interface{})["three"].([]int)[2],
			err:       nil,
		},
		{
			in:        setupDocument(),
			path:      "one.two.four",
			separator: ".",
			out:       nil,
			err:       fmt.Errorf("Property %s does not exist", "four"),
		},
		{
			in:        setupDocument(),
			path:      "one.two.four[0]",
			separator: ".",
			out:       nil,
			err:       fmt.Errorf("Property %s does not exist", "four"),
		},
		{
			in:        setupDocument_I(),
			path:      "one[0]",
			separator: ".",
			out:       setupDocument_I()["one"].([]map[string]interface{})[0],
			err:       nil,
		},
		{
			in:        setupDocument_I(),
			path:      "one[1]",
			separator: ".",
			out:       setupDocument_I()["one"].([]map[string]interface{})[1],
			err:       nil,
		},
		{
			in:        setupDocument_I(),
			path:      "one[2]",
			separator: ".",
			out:       setupDocument_I()["one"].([]map[string]interface{})[2],
			err:       nil,
		},
		{
			in:        setupDocument_I(),
			path:      "one[2].map_c",
			separator: ".",
			out:       setupDocument_I()["one"].([]map[string]interface{})[2]["map_c"],
			err:       nil,
		},
		{
			in:        setupDocument_II(),
			path:      "one[1].two[1]",
			separator: ".",
			out:       setupDocument_II()["one"].([]map[string]interface{})[1]["two"].([]map[string]interface{})[1],
			err:       nil,
		},
		{
			in:        setupDocument_II(),
			path:      "one[2].two[1].eight",
			separator: ".",
			out:       setupDocument_II()["one"].([]map[string]interface{})[2]["two"].([]map[string]interface{})[1]["eight"],
			err:       nil,
		},
		{
			in:        setupDocument_II(),
			path:      "one[1].two[1].eight",
			separator: ".",
			out:       nil,
			err:       fmt.Errorf("Property %s does not exist", "eight"),
		},
		{
			in:        setupDocument_II(),
			path:      "one[3].three[0].seven.eight",
			separator: ".",
			out:       nil,
			err:       fmt.Errorf("Property %s does not exist", "seven"),
		},
	}

	num_cases := len(cases)
	for i, c := range cases {
		case_index := i + 1

		out, err_case := GetProperty(c.in, c.path, c.separator)
		if !reflect.DeepEqual(c.err, err_case) {
			t.Errorf("\n[%d of %d: Errors should equal] \n\t%v \n \n\t%v", case_index, num_cases, err_case, c.err)
		}
		if !reflect.DeepEqual(out, c.out) {
			t.Errorf("\n[%d of %d: Results should equal] \n\t%v \n \n\t%v", case_index, num_cases, out, c.out)
		}
	}
}

func setupDocument() (document map[string]interface{}) {
	document = map[string]interface{}{
		"one": map[string]interface{}{
			"two": map[string]interface{}{
				"three": []int{
					1, 2, 3,
				},
			},
			"four": map[string]interface{}{
				"five": []int{
					11, 22, 33,
				},
			},
		},
	}

	return
}

func setupDocument_I() (document_I map[string]interface{}) {
	document_I = map[string]interface{}{
		"one": []map[string]interface{}{
			{"map_a": []int{1, 2, 3}},
			{"map_b": []int{4, 5, 6}},
			{"map_c": []int{7, 8, 9}},
		},
	}
	return
}

func setupDocument_II() (document_II map[string]interface{}) {
	document_II = map[string]interface{}{
		"one": []map[string]interface{}{
			{
				"two": []map[string]interface{}{
					{"three": "got three"},
					{"four": "got four"},
				},
			},
			{
				"two": []map[string]interface{}{
					{"five": "got five"},
					{"six": "got six"},
				},
			},
			{
				"two": []map[string]interface{}{
					{"seven": "got seven"},
					{"eight": "got eight"},
				},
			},
			{
				"three": []map[string]interface{}{
					{"four": map[string]interface{}{
						"five": "six",
					}},
					{"seven": map[string]interface{}{
						"eight": "ten",
					}},
				},
			},
		},
	}
	return
}

func setup() (document, document_I, document_II map[string]interface{}) {
	document = setupDocument()
	document_I = setupDocument_I()
	document_II = setupDocument_II()
	return
}
