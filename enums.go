package hpc015

import "strconv"

// TimeVerifyMode
//
// In manual, written as `Commond Type`, command type of requesting for data:
//  0x00 exclude the verification hours and business hours
//  0x01 include the time of verifying the system
//  0x02 include the time of verifying the business hours
//  0x03 include the time of verifying the system and business hours
type TimeVerifyMode byte

const (
	Exclude TimeVerifyMode = iota
	System
	Business
	Both
)

var timeVerifyModeString = []string{
	"Exclude",
	"System",
	"Business",
	"Both",
}

func (m TimeVerifyMode) String() string {
	if Exclude <= m && m <= Both {
		return timeVerifyModeString[m]
	}
	return "%!Weekday(" + strconv.FormatUint(uint64(m), 10) + ")"
}

// NetworkType
//
// In manual, written as `Model`
type NetworkType byte

const (
	Online NetworkType = iota
	StandAlone
)

var networkTypeString = []string{
	"Online",
	"StandAlone",
}

func (m NetworkType) String() string {
	if Online <= m && m <= StandAlone {
		return networkTypeString[m]
	}
	return "%!NetworkType(" + strconv.FormatUint(uint64(m), 10) + ")"
}

// RespondingType represent whether configuration changed or not.
//
// Usually not to need use this type directly,
// `SetConfiguration()` will set this.
type RespondingType byte

const (
	NewParameterValue RespondingType = iota + 4 // new parameter value
	Confirmation                                // parameter confirmation, after confirmation and responding, the parameter will be neglected.
)

var respondingTypeString = []string{
	"NewParameterValue",
	"Confirmation",
}

func (m RespondingType) String() string {
	if NewParameterValue <= m && m <= Confirmation {
		return respondingTypeString[m-4]
	}
	return "%!RespondingType(" + strconv.FormatUint(uint64(m), 10) + ")"
}

// Speed represent Equipment detects speed.
type Speed byte

const (
	Low Speed = iota
	High
)

var speedString = []string{
	"Low",
	"High",
}

func (m Speed) String() string {
	if Low <= m && m <= High {
		return speedString[m]
	}
	return "%!Speed(" + strconv.FormatUint(uint64(m), 10) + ")"
}

// Display type
//
// In manual, written as `Disable Type`
//  0x00 the counting is not displayed on the screen.
//  0x01 display total amount
//  0x02 display bilateral
type DisplayType byte

const (
	None DisplayType = iota
	Unidirectinal
	Bilateral
)

var displayTypeString = []string{
	"None",
	"Unidirectinal",
	"Bilateral",
}

func (m DisplayType) String() string {
	if None <= m && m <= Bilateral {
		return displayTypeString[m]
	}
	return "%!DisplayType(" + strconv.FormatUint(uint64(m), 10) + ")"
}

// AnswerType represent whether status of cache response.
type AnswerType byte

const (
	Failed AnswerType = iota
	OK
)

var answerType = []string{
	"Failed",
	"OK",
}

func (m AnswerType) String() string {
	if Failed <= m && m <= OK {
		return answerType[m]
	}
	return "%!AnswerType(" + strconv.FormatUint(uint64(m), 10) + ")"
}

type Focus byte

const (
	Focused Focus = iota
	FocusOut
)

var focusString = []string{
	"Focused",
	"FocusOut",
}

func (m Focus) String() string {
	if Focused <= m && m <= FocusOut {
		return focusString[m]
	}
	return "%!Focus(" + strconv.FormatUint(uint64(m), 10) + ")"
}

type Charge byte

const (
	NotCharged Charge = iota
	_
	BeingCharged
)

var cargeString = []string{
	"NotCharged",
	"",
	"BeingCharged",
}

func (m Charge) String() string {
	if NotCharged == m || m == BeingCharged {
		return cargeString[m]
	}
	return "%!Charge(" + strconv.FormatUint(uint64(m), 10) + ")"
}
