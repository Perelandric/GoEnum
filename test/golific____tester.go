/****************************************************************************
	This file was generated by Golific.

	Do not edit this file. If you do, your changes will be overwritten the next
	time 'generate' is invoked.
******************************************************************************/

package main

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"
)

/******************************************************************************
	ENUM SUMMARY

Foo (type FooEnum, uint8)
	Bar 1 "bar" "bar"
	Baz 2 "baz" "This is the description"
	Buz 4 "Buz" "Buz"

Animal (type AnimalEnum, uint8)
	Dog 1 "doggy" "Your best friend, and you know it."
	Cat 2 "kitty" "Your best friend, but doesn't always show it."
	Horse 3 "horsie" "Everyone loves horses."
******************************************************************************/

/*****************************

Tester struct

******************************/
func NewTester() *Tester {
	return &Tester{
		private: private_anp1xa86ymnb{
			Test2: "foo",
		},
		AnimalEnum: Animal.Dog,
		Test4:      "bar",
	}
}

type Tester struct {
	private private_anp1xa86ymnb
	AnimalEnum
	Test4 string
	Test5 string
}

type private_anp1xa86ymnb struct {
	Test1 string `json:"test1"`
	Test2 string
	Test3 string
}

type json_anp1xa86ymnb struct {
	*private_anp1xa86ymnb
	AnimalEnum
	Test4 string
	Test5 string
}

func (self *Tester) Test2() string {
	return self.private.Test2
}

func (self *Tester) SetTest3(v string) {
	self.private.Test3 = v
}

func (self *Tester) MarshalJSON() ([]byte, error) {
	return json.Marshal(json_anp1xa86ymnb{
		&self.private,
		self.AnimalEnum,
		self.Test4,
		self.Test5,
	})
}

func (self *Tester) UnmarshalJSON(j []byte) error {
	var temp json_anp1xa86ymnb
	if err := json.Unmarshal(j, &temp); err != nil {
		return err
	}
	self.private = *temp.private_anp1xa86ymnb
	self.AnimalEnum = temp.AnimalEnum
	self.Test4 = temp.Test4
	self.Test5 = temp.Test5
	return nil
}

/*****************************

FooEnum - bit flags

******************************/

type FooEnum struct{ value_o9084kckqnl uint8 }

var Foo = struct {
	Bar FooEnum
	Baz FooEnum
	Buz FooEnum

	// Used to iterate in range loops
	foobar [3]FooEnum
}{
	Bar: FooEnum{value_o9084kckqnl: 1},
	Baz: FooEnum{value_o9084kckqnl: 2},
	Buz: FooEnum{value_o9084kckqnl: 4},
}

func init() {
	Foo.foobar = [3]FooEnum{
		Foo.Bar, Foo.Baz, Foo.Buz,
	}
}

// Value returns the numeric value of the variant as a uint8.
func (Fe FooEnum) Value() uint8 {
	return Fe.value_o9084kckqnl
}

// IntValue is the same as 'Value()', except that the value is cast to an 'int'.
func (Fe FooEnum) IntValue() int {
	return int(Fe.value_o9084kckqnl)
}

// Name returns the name of the variant as a string.
func (Fe FooEnum) Name() string {
	switch Fe.value_o9084kckqnl {
	case 1:
		return "Bar"
	case 2:
		return "Baz"
	case 4:
		return "Buz"
	}

	return ""
}

// String returns the given string value of the variant. If none has been set,
// its return value is as though 'Name()' had been called.
// If multiple bit values are assigned, the string values will be joined into a
// single string using "," as a separator.
func (Fe FooEnum) String() string {
	switch Fe.value_o9084kckqnl {
	case 1:
		return "bar"
	case 2:
		return "baz"
	case 4:
		return "Buz"
	}

	if Fe.value_o9084kckqnl == 0 {
		return ""
	}

	var vals = make([]string, 0, 3/2)

	for _, item := range Foo.foobar {
		if Fe.value_o9084kckqnl&item.value_o9084kckqnl == item.value_o9084kckqnl {
			vals = append(vals, item.String())
		}
	}
	return strings.Join(vals, ",")
}

// Description returns the description of the variant. If none has been set, its
// return value is as though 'String()' had been called.
func (Fe FooEnum) Description() string {
	switch Fe.value_o9084kckqnl {
	case 1:
		return "bar"
	case 2:
		return "This is the description"
	case 4:
		return "Buz"
	}
	return ""
}

// JSON marshaling methods
func (Fe FooEnum) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(Fe.String())), nil
}

func (Fe *FooEnum) UnmarshalJSON(b []byte) error {
	var s, err = strconv.Unquote(string(b))
	if err != nil {
		return err
	}

	if len(s) == 0 {
		return nil
	}

	switch s {
	case "bar":
		Fe.value_o9084kckqnl = 1
		return nil
	case "baz":
		Fe.value_o9084kckqnl = 2
		return nil
	case "Buz":
		Fe.value_o9084kckqnl = 4
		return nil
	}

	var val = 0

	for _, part := range strings.Split(string(b), ",") {
		switch part {
		case "bar":
			val |= 1
		case "baz":
			val |= 2
		case "Buz":
			val |= 4
		default:
			log.Printf("Unexpected value: %q while unmarshaling FooEnum\n", part)
		}
	}

	Fe.value_o9084kckqnl = uint8(val)
	return nil
}

// Bitflag enum methods

// Add returns a copy of the variant with the value of 'v' added to it.
func (Fe FooEnum) Add(v FooEnum) FooEnum {
	Fe.value_o9084kckqnl |= v.value_o9084kckqnl
	return Fe
}

// AddAll returns a copy of the variant with all the values of 'v' added to it.
func (Fe FooEnum) AddAll(v ...FooEnum) FooEnum {
	for _, item := range v {
		Fe.value_o9084kckqnl |= item.value_o9084kckqnl
	}
	return Fe
}

// Remove returns a copy of the variant with the value of 'v' removed from it.
func (Fe FooEnum) Remove(v FooEnum) FooEnum {
	Fe.value_o9084kckqnl &^= v.value_o9084kckqnl
	return Fe
}

// RemoveAll returns a copy of the variant with all the values of 'v' removed
// from it.
func (Fe FooEnum) RemoveAll(v ...FooEnum) FooEnum {
	for _, item := range v {
		Fe.value_o9084kckqnl &^= item.value_o9084kckqnl
	}
	return Fe
}

// Has returns 'true' if the receiver contains the value of 'v', otherwise
// 'false'.
func (Fe FooEnum) Has(v FooEnum) bool {
	return Fe.value_o9084kckqnl&v.value_o9084kckqnl == v.value_o9084kckqnl
}

// HasAny returns 'true' if the receiver contains any of the values of 'v',
// otherwise 'false'.
func (Fe FooEnum) HasAny(v ...FooEnum) bool {
	for _, item := range v {
		if Fe.value_o9084kckqnl&item.value_o9084kckqnl == item.value_o9084kckqnl {
			return true
		}
	}
	return false
}

// HasAll returns 'true' if the receiver contains all the values of 'v',
// otherwise 'false'.
func (Fe FooEnum) HasAll(v ...FooEnum) bool {
	for _, item := range v {
		if Fe.value_o9084kckqnl&item.value_o9084kckqnl != item.value_o9084kckqnl {
			return false
		}
	}
	return true
}

/*****************************

OofEnum

******************************/

type OofEnum struct{ value_19zqulavsbiwc uint8 }

var Oof = struct {
	Bar OofEnum
	Baz OofEnum
	Buz OofEnum

	// Used to iterate in range loops
	Values [3]OofEnum
}{
	Bar: OofEnum{value_19zqulavsbiwc: 1},
	Baz: OofEnum{value_19zqulavsbiwc: 123},
	Buz: OofEnum{value_19zqulavsbiwc: 3},
}

func init() {
	Oof.Values = [3]OofEnum{
		Oof.Bar, Oof.Baz, Oof.Buz,
	}
}

// Value returns the numeric value of the variant as a uint8.
func (Oe OofEnum) Value() uint8 {
	return Oe.value_19zqulavsbiwc
}

// IntValue is the same as 'Value()', except that the value is cast to an 'int'.
func (Oe OofEnum) IntValue() int {
	return int(Oe.value_19zqulavsbiwc)
}

// Name returns the name of the variant as a string.
func (Oe OofEnum) Name() string {
	switch Oe.value_19zqulavsbiwc {
	case 1:
		return "Bar"
	case 123:
		return "Baz"
	case 3:
		return "Buz"
	}

	return ""
}

// String returns the given string value of the variant. If none has been set,
// its return value is as though 'Name()' had been called.

func (Oe OofEnum) String() string {
	switch Oe.value_19zqulavsbiwc {
	case 1:
		return "bar"
	case 123:
		return "Baz"
	case 3:
		return "Buz"
	}

	return ""
}

// Description returns the description of the variant. If none has been set, its
// return value is as though 'String()' had been called.
func (Oe OofEnum) Description() string {
	switch Oe.value_19zqulavsbiwc {
	case 1:
		return "bar"
	case 123:
		return "Baz"
	case 3:
		return "Some description"
	}
	return ""
}

// JSON marshaling methods
func (Oe OofEnum) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Itoa(Oe.IntValue())), nil
}

func (Oe *OofEnum) UnmarshalJSON(b []byte) error {
	var n, err = strconv.ParseUint(string(b), 10, 64)
	if err != nil {
		return err
	}
	Oe.value_19zqulavsbiwc = uint8(n)
	return nil
}

/*****************************

AnimalEnum

******************************/

type AnimalEnum struct{ value_ofl4mtn9el8x uint8 }

var Animal = struct {
	Dog   AnimalEnum
	Cat   AnimalEnum
	Horse AnimalEnum

	// Used to iterate in range loops
	Values [3]AnimalEnum
}{
	Dog:   AnimalEnum{value_ofl4mtn9el8x: 1},
	Cat:   AnimalEnum{value_ofl4mtn9el8x: 2},
	Horse: AnimalEnum{value_ofl4mtn9el8x: 3},
}

func init() {
	Animal.Values = [3]AnimalEnum{
		Animal.Dog, Animal.Cat, Animal.Horse,
	}
}

// Value returns the numeric value of the variant as a uint8.
func (Ae AnimalEnum) Value() uint8 {
	return Ae.value_ofl4mtn9el8x
}

// IntValue is the same as 'Value()', except that the value is cast to an 'int'.
func (Ae AnimalEnum) IntValue() int {
	return int(Ae.value_ofl4mtn9el8x)
}

// Name returns the name of the variant as a string.
func (Ae AnimalEnum) Name() string {
	switch Ae.value_ofl4mtn9el8x {
	case 1:
		return "Dog"
	case 2:
		return "Cat"
	case 3:
		return "Horse"
	}

	return ""
}

// String returns the given string value of the variant. If none has been set,
// its return value is as though 'Name()' had been called.

func (Ae AnimalEnum) String() string {
	switch Ae.value_ofl4mtn9el8x {
	case 1:
		return "doggy"
	case 2:
		return "kitty"
	case 3:
		return "horsie"
	}

	return ""
}

// Description returns the description of the variant. If none has been set, its
// return value is as though 'String()' had been called.
func (Ae AnimalEnum) Description() string {
	switch Ae.value_ofl4mtn9el8x {
	case 1:
		return "Your best friend, and you know it."
	case 2:
		return "Your best friend, but doesn't always show it."
	case 3:
		return "Everyone loves horses."
	}
	return ""
}

// JSON marshaling methods
func (Ae AnimalEnum) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(Ae.String())), nil
}

func (Ae *AnimalEnum) UnmarshalJSON(b []byte) error {
	var s, err = strconv.Unquote(string(b))
	if err != nil {
		return err
	}

	if len(s) == 0 {
		return nil
	}

	switch s {
	case "doggy":
		Ae.value_ofl4mtn9el8x = 1
		return nil
	case "kitty":
		Ae.value_ofl4mtn9el8x = 2
		return nil
	case "horsie":
		Ae.value_ofl4mtn9el8x = 3
		return nil
	default:
		log.Printf("Unexpected value: %q while unmarshaling AnimalEnum\n", s)
	}

	return nil
}
