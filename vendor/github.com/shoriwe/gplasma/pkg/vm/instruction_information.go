package vm

type ClassInformation struct {
	Name string
	Body []*Code
}

type FunctionInformation struct {
	Name              string
	Body              []*Code
	NumberOfArguments int
}

type ConditionInformation struct {
	Body     []*Code
	ElseBody []*Code
}

type LoopInformation struct {
	Body      []*Code
	Condition []*Code
	Receivers []string
}

type ExceptInformation struct {
	CaptureName string
	Targets     []*Code
	Body        []*Code
}

type TryInformation struct {
	Body    []*Code
	Excepts []ExceptInformation
	Finally []*Code
}

type CaseInformation struct {
	Targets []*Code
	Body    []*Code
}

type SwitchInformation struct {
	Cases   []CaseInformation
	Default []*Code
}
