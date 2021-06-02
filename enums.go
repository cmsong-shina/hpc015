package hpc015

// TimeVerifyMode
//
// In manual, written as `Commond Type`
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
	return timeVerifyModeString[m]
}

// Operation mode
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
	return networkTypeString[m]
}

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
	return respondingTypeString[m-4]
}

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
	return speedString[m]
}

// Display type
//
//  0x00 the counting is not displayed on the screen.
//  0x01 display total amount
//  0x02 display bilateral
type DisplayType byte

const (
	None DisplayType = iota
	All
	Bilateral
)

var displayTypeString = []string{
	"None",
	"All",
	"Bilateral",
}

func (m DisplayType) String() string {
	return displayTypeString[m]
}

type AnswerType byte

// AnswerType represent wethere upload is failed or not
const (
	Failed AnswerType = iota
	OK
)

var answerType = []string{
	"Failed",
	"OK",
}

func (m AnswerType) String() string {
	return answerType[m]
}
