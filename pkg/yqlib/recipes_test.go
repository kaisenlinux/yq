package yqlib

import (
	"testing"
)

var recipes = []expressionScenario{
	{
		description:    "Find items in an array",
		subdescription: "We have an array and we want to find the elements with a particular name.",
		explanation: []string{
			"`.[]` splats the array, and puts all the items in the context.",
			"These items are then piped (`|`) into `select(.name == \"Foo\")` which will select all the nodes that have a name property set to 'Foo'.",
			"See the [select](https://mikefarah.gitbook.io/yq/operators/select) operator for more information.",
		},
		document:   `[{name: Foo, numBuckets: 0}, {name: Bar, numBuckets: 0}]`,
		expression: `.[] | select(.name == "Foo")`,
		expected: []string{
			"D0, P[0], (!!map)::{name: Foo, numBuckets: 0}\n",
		},
	},
	{
		description:    "Find and update items in an array",
		subdescription: "We have an array and we want to _update_ the elements with a particular name.",
		document:       `[{name: Foo, numBuckets: 0}, {name: Bar, numBuckets: 0}]`,
		expression:     `(.[] | select(.name == "Foo") | .numBuckets) |= . + 1`,
		explanation: []string{
			"Following from the example above`.[]` splats the array, selects filters the items.",
			"We then pipe (`|`) that into `.numBuckets`, which will select that field from all the matching items",
			"Splat, select and the field are all in brackets, that whole expression is passed to the `|=` operator as the left hand side expression, with `. + 1` as the right hand side expression.",
			"`|=` is the operator that updates fields relative to their own value, which is referenced as dot (`.`).",
			"The expression `. + 1` increments the numBuckets counter.",
			"See the [assign](https://mikefarah.gitbook.io/yq/operators/assign-update) and [add](https://mikefarah.gitbook.io/yq/operators/add) operators for more information.",
		},
		expected: []string{
			"D0, P[], (doc)::[{name: Foo, numBuckets: 1}, {name: Bar, numBuckets: 0}]\n",
		},
	},
	{
		description:    "Multiple or complex updates to items in an array",
		subdescription: "We have an array and we want to _update_ the elements with a particular name in reference to its type.",
		document:       `myArray: [{name: Foo, type: cat}, {name: Bar, type: dog}]`,
		expression:     `with(.myArray[]; .name = .name + " - " + .type)`,
		explanation: []string{
			"The with operator will effectively loop through each given item in the first given expression, and run the second expression against it.",
			"`.myArray[]` splats the array in `myArray`. So `with` will run against each item in that array",
			"`.name = .name + \" - \" + .type` this expression is run against every item, updating the name to be a concatenation of the original name as well as the type.",
			"See the [with](https://mikefarah.gitbook.io/yq/operators/with) operator for more information and examples.",
		},
		expected: []string{
			"D0, P[], (doc)::myArray: [{name: Foo - cat, type: cat}, {name: Bar - dog, type: dog}]\n",
		},
	},
	{
		description: "Sort an array by a field",
		document:    `myArray: [{name: Foo, numBuckets: 1}, {name: Bar, numBuckets: 0}]`,
		expression:  `.myArray |= sort_by(.numBuckets)`,
		explanation: []string{
			"We want to resort `.myArray`.",
			"`sort_by` works by piping an array into it, and it pipes out a sorted array.",
			"So, we use `|=` to update `.myArray`. This is the same as doing `.myArray = (.myArray | sort_by(.numBuckets))`",
		},
		expected: []string{
			"D0, P[], (doc)::myArray: [{name: Bar, numBuckets: 0}, {name: Foo, numBuckets: 1}]\n",
		},
	},
	{
		description:    "Filter, flatten, sort and unique",
		subdescription: "Lets",
		document:       `[{type: foo, names: [Fred, Catherine]}, {type: bar, names: [Zelda]}, {type: foo, names: Fred}, {type: foo, names: Ava}]`,
		expression:     `[.[] | select(.type == "foo") | .names] | flatten | sort | unique`,
		explanation: []string{
			"`.[] | select(.type == \"foo\") | .names` will select the array elements of type \"foo\"",
			"Splat `.[]` will unwrap the array and match all the items. We need to do this so we can work on the child items, for instance, filter items out using the `select` operator.",
			"But we still want the final results back into an array. So after we're doing working on the children, we wrap everything back into an array using square brackets around the expression. `[.[] | select(.type == \"foo\") | .names]`",
			"Now have have an array of all the 'names' values. Which includes arrays of strings as well as strings on their own.",
			"Pipe `|` this array through `flatten`. This will flatten nested arrays. So now we have a flat list of all the name value strings",
			"Next we pipe `|` that through `sort` and then `unique` to get a sorted, unique list of the names!",
			"See the [flatten](https://mikefarah.gitbook.io/yq/operators/flatten), [sort](https://mikefarah.gitbook.io/yq/operators/sort) and [unique](https://mikefarah.gitbook.io/yq/operators/unique) for more information and examples.",
		},
		expected: []string{
			"D0, P[], (!!seq)::- Ava\n- Catherine\n- Fred\n",
		},
	},
}

func TestRecipes(t *testing.T) {
	for _, tt := range recipes {
		testScenario(t, &tt)
	}
	genericScenarios := make([]interface{}, len(recipes))
	for i, s := range recipes {
		genericScenarios[i] = s
	}
	documentScenarios(t, "usage", "recipes", genericScenarios, documentOperatorScenario)
}
