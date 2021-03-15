package main

type accum interface {
	matcher(text *string) *[]string
	add(strs *[]string)
	print(mincount int32)
}
