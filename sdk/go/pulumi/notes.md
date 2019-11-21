- Replace use of interface{} with Pulumi-specific types
- These types will be ~derived from resource.PropertyValue
- Must be possible to compose types
	- Implies some sort of marshalling behavior
- These types will all be apply-able?
- input type -> output type -> resolved type
	- Any -> AnyOutput -> interface{}
	- Slice -> SliceOutput -> []interface{}
	- Map -> MapOutput -> map[string]interface{}

	- String -> StringOutput -> string
	- Strings -> StringsOutput -> []string
	- StringMap -> StringMapOutput -> map[string]string

	- structs:
		- resolved type: all fields are "normal" Go types; struct fields are tagged with Pulumi property names
		- input type: interface implemented by both input value type and output
		- input value type: the resolved type, but all field types are the corresponding input type
		- output type: output with an apply that takes a callback that accepts the resolved type

	- decoding into Output types:
		- use the type of the output to decide on decoding

	- casting from Output -> specific output:
		- type check before apply?
		- type check at cast time?
			- this has the benefit of catching errors during preview

- need to make sure that RPC errors are properly recorded s.t. the program fails appropriately
