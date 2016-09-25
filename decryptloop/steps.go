// 25 september 2016
package decryptloop

import (
	"fmt"
	"crypto/cipher"

	"github.com/andlabs/reallymine/byteops"
)

// This complicated structure allows us to define a fixed set of Step objects and disallow nil at the same time, reducing the number of things that need validation.

type stepiface interface {
	name()					string
	desc()					string
	do(c cipher.Block, b []byte)
}

type step struct {
	s	stepiface
}

type Step step

// for usage information
var validSteps []Step

var stepsByName = make(map[string]Step)

func addstep(s stepiface) Step {
	ss := Step{s}
	validSteps = append(validSteps, ss)
	stepsByName[ss.name()] = ss
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

// for diskusage.go

func (s Step) name() string {
	return s.s.name()
}

func (s Step) do(c cipher.Block, b []byte) {
	s.s.do(c, b)
}

// for reallymine to use directly

// TODO merge with package command
func StepUsage() string {
	s := ""
	for _, step := range validSteps {
		s += fmt.Sprintf("  %s - %s\n", step.name(), step.s.desc())
	}
	return s
}
