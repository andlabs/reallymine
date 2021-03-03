// 25 september 2016
package decryptloop

import (
	"crypto/cipher"
	"fmt"
	"strings"

	"github.com/undeadbanegithub/reallymine/byteops"
)

// This complicated structure allows us to define a fixed set of Step objects and disallow nil at the same time, reducing the number of things that need validation.

type stepiface interface {
	name() string
	desc() string
	do(c cipher.Block, b []byte)
}

type step struct {
	s stepiface
}

type Step step

// for usage information
var validSteps []Step

var stepsByName = make(map[string]Step)

// TODO complain on duplicates
func addstep(s stepiface) Step {
	ss := Step{s}
	validSteps = append(validSteps, ss)
	stepsByName[ss.s.name()] = ss
	return ss
}

type stepDecryptType struct{}

func (stepDecryptType) name() string {
	return "decrypt"
}

func (stepDecryptType) desc() string {
	return "Decrypts the block"
}

func (stepDecryptType) do(c cipher.Block, b []byte) {
	c.Decrypt(b, b)
}

var stepDecrypt = addstep(&stepDecryptType{})
var StepDecrypt Step = stepDecrypt

type stepReverseType struct{}

func (stepReverseType) name() string {
	return "reverse"
}

func (stepReverseType) desc() string {
	return "Reverses the block"
}

func (stepReverseType) do(c cipher.Block, b []byte) {
	byteops.Reverse(b)
}

var stepReverse = addstep(&stepReverseType{})
var StepReverse Step = stepReverse

type stepSwapWordsType struct{}

func (stepSwapWordsType) name() string {
	return "swapwords"
}

func (stepSwapWordsType) desc() string {
	return "Reverses each word (two-byte pair) of the block"
}

func (stepSwapWordsType) do(c cipher.Block, b []byte) {
	byteops.SwapWords(b)
}

var stepSwapWords = addstep(&stepSwapWordsType{})
var StepSwapWords Step = stepSwapWords

type stepSwapLongsType struct{}

func (stepSwapLongsType) name() string {
	return "swaplongs"
}

func (stepSwapLongsType) desc() string {
	return "Reverses each long (four-byte group) of the block"
}

func (stepSwapLongsType) do(c cipher.Block, b []byte) {
	byteops.SwapLongs(b)
}

var stepSwapLongs = addstep(&stepSwapLongsType{})
var StepSwapLongs Step = stepSwapLongs

type stepSwapHalvesType struct{}

func (stepSwapHalvesType) name() string {
	return "swaphalves"
}

func (stepSwapHalvesType) desc() string {
	return "Switches the two halves of the block"
}

func (stepSwapHalvesType) do(c cipher.Block, b []byte) {
	byteops.SwapHalves(b)
}

var stepSwapHalves = addstep(&stepSwapHalvesType{})
var StepSwapHalves Step = stepSwapHalves

type StepList []Step

type UnknownStepNameError string

func (e UnknownStepNameError) Error() string {
	return fmt.Sprintf("unknown decrypt loop step name %q", string(e))
}

var ErrStepListStringEmpty = fmt.Errorf("step list string is empty/specifies no steps")

func StepListFromString(s string) (StepList, error) {
	names := strings.Split(s, " ")
	if len(names) == 0 {
		return nil, ErrStepListStringEmpty
	}
	steps := make(StepList, len(names))
	for i, name := range names {
		step, ok := stepsByName[name]
		if !ok {
			return nil, UnknownStepNameError(name)
		}
		steps[i] = step
	}
	return steps, nil
}

func (s StepList) String() string {
	// TODO remove this when PLX is done
	if len(s) == 0 {
		return "(unknown)"
	}
	names := make([]string, len(s))
	for i, step := range s {
		names[i] = step.s.name()
	}
	return strings.Join(names, " ")
}

// for diskusage.go

func (s StepList) runBlock(c cipher.Block, b []byte) {
	for _, step := range s {
		step.s.do(c, b)
	}
}

// for reallymine to use directly

// TODO merge with package command - can't use command directly since that imports us
func StepUsage() string {
	s := ""
	for _, step := range validSteps {
		s += fmt.Sprintf("  %s - %s\n", step.s.name(), step.s.desc())
	}
	return s
}
