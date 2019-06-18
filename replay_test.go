package log

import (
	"github.com/aphistic/sweet"
	"github.com/efritz/glock"
	. "github.com/onsi/gomega"
)

type ReplaySuite struct{}

func (s *ReplaySuite) TestReplay(t sweet.T) {
	var (
		shim    = &testShim{}
		clock   = glock.NewMockClock()
		adapter = newReplayShim(adaptShim(shim), clock, LevelDebug)
	)

	adapter.LogWithFields(LevelDebug, LogFields{"x": "x"}, "foo", 12)
	adapter.LogWithFields(LevelDebug, LogFields{"y": "y"}, "bar", 43)
	adapter.LogWithFields(LevelDebug, LogFields{"z": "z"}, "baz", 74)
	adapter.Replay(LevelWarning)

	messages := shim.copy()
	Expect(messages).To(HaveLen(6))

	for i := 0; i < 3; i++ {
		Expect(messages[i+0].level).To(Equal(LevelDebug))
		Expect(messages[i+3].level).To(Equal(LevelWarning))
	}

	for i, format := range []string{"foo", "bar", "baz"} {
		Expect(messages[i+0].format).To(Equal(format))
		Expect(messages[i+3].format).To(Equal(format))
	}

	for i, expected := range []int{12, 43, 74} {
		Expect(messages[i+0].args[0]).To(Equal(expected))
		Expect(messages[i+3].args[0]).To(Equal(expected))
	}

	for i, field := range []string{"x", "y", "z"} {
		Expect(messages[i+0].fields[field]).To(Equal(field))
		Expect(messages[i+3].fields[field]).To(Equal(field))
	}
}

func (s *ReplaySuite) TestReplayTwice(t sweet.T) {
	var (
		shim    = &testShim{}
		clock   = glock.NewMockClock()
		adapter = newReplayShim(adaptShim(shim), clock, LevelDebug)
	)

	adapter.LogWithFields(LevelDebug, nil, "foo")
	adapter.LogWithFields(LevelDebug, nil, "bar")
	adapter.LogWithFields(LevelDebug, nil, "baz")
	adapter.Replay(LevelWarning)
	adapter.Replay(LevelError)

	messages := shim.copy()
	Expect(messages).To(HaveLen(9))
	Expect(messages[0].level).To(Equal(LevelDebug))
	Expect(messages[1].level).To(Equal(LevelDebug))
	Expect(messages[2].level).To(Equal(LevelDebug))
	Expect(messages[3].level).To(Equal(LevelWarning))
	Expect(messages[4].level).To(Equal(LevelWarning))
	Expect(messages[5].level).To(Equal(LevelWarning))
	Expect(messages[6].level).To(Equal(LevelError))
	Expect(messages[7].level).To(Equal(LevelError))
	Expect(messages[8].level).To(Equal(LevelError))

	for i, format := range []string{"foo", "bar", "baz", "foo", "bar", "baz", "foo", "bar", "baz"} {
		Expect(messages[i].format).To(Equal(format))
	}
}

func (s *ReplaySuite) TestReplayAtHigherlevelNoops(t sweet.T) {
	var (
		shim    = &testShim{}
		clock   = glock.NewMockClock()
		adapter = newReplayShim(adaptShim(shim), clock, LevelDebug)
	)

	adapter.LogWithFields(LevelDebug, nil, "foo")
	adapter.LogWithFields(LevelDebug, nil, "bar")
	adapter.LogWithFields(LevelDebug, nil, "baz")
	adapter.Replay(LevelError)
	adapter.Replay(LevelWarning)

	messages := shim.copy()
	Expect(messages).To(HaveLen(6))
	Expect(messages[0].level).To(Equal(LevelDebug))
	Expect(messages[1].level).To(Equal(LevelDebug))
	Expect(messages[2].level).To(Equal(LevelDebug))
	Expect(messages[3].level).To(Equal(LevelError))
	Expect(messages[4].level).To(Equal(LevelError))
	Expect(messages[5].level).To(Equal(LevelError))

	for i, format := range []string{"foo", "bar", "baz", "foo", "bar", "baz"} {
		Expect(messages[i].format).To(Equal(format))
	}
}

func (s *ReplaySuite) TestLogAfterReplaySendsImmediately(t sweet.T) {
	var (
		shim    = &testShim{}
		clock   = glock.NewMockClock()
		adapter = newReplayShim(adaptShim(shim), clock, LevelDebug)
	)

	adapter.LogWithFields(LevelDebug, nil, "foo")
	adapter.LogWithFields(LevelDebug, nil, "bar")
	adapter.LogWithFields(LevelDebug, nil, "baz")
	adapter.Replay(LevelWarning)
	adapter.LogWithFields(LevelDebug, nil, "bnk")
	adapter.LogWithFields(LevelDebug, nil, "qux")

	messages := shim.copy()
	Expect(messages).To(HaveLen(10))
	Expect(messages[0].level).To(Equal(LevelDebug))
	Expect(messages[1].level).To(Equal(LevelDebug))
	Expect(messages[2].level).To(Equal(LevelDebug))
	Expect(messages[3].level).To(Equal(LevelWarning))
	Expect(messages[4].level).To(Equal(LevelWarning))
	Expect(messages[5].level).To(Equal(LevelWarning))
	Expect(messages[6].level).To(Equal(LevelDebug))
	Expect(messages[7].level).To(Equal(LevelWarning))
	Expect(messages[8].level).To(Equal(LevelDebug))
	Expect(messages[9].level).To(Equal(LevelWarning))

	for i, format := range []string{"foo", "bar", "baz", "foo", "bar", "baz", "bnk", "bnk", "qux", "qux"} {
		Expect(messages[i].format).To(Equal(format))
	}
}

func (s *ReplaySuite) TestLogAfterSecondReplaySendsAtNewLevel(t sweet.T) {
	var (
		shim    = &testShim{}
		clock   = glock.NewMockClock()
		adapter = newReplayShim(adaptShim(shim), clock, LevelDebug)
	)

	adapter.LogWithFields(LevelDebug, nil, "foo")
	adapter.LogWithFields(LevelDebug, nil, "bar")
	adapter.Replay(LevelWarning)
	adapter.Replay(LevelError)
	adapter.LogWithFields(LevelDebug, nil, "baz")
	adapter.LogWithFields(LevelDebug, nil, "bnk")

	messages := shim.copy()
	Expect(messages).To(HaveLen(10))
	Expect(messages[0].level).To(Equal(LevelDebug))
	Expect(messages[1].level).To(Equal(LevelDebug))
	Expect(messages[2].level).To(Equal(LevelWarning))
	Expect(messages[3].level).To(Equal(LevelWarning))
	Expect(messages[4].level).To(Equal(LevelError))
	Expect(messages[5].level).To(Equal(LevelError))
	Expect(messages[6].level).To(Equal(LevelDebug))
	Expect(messages[7].level).To(Equal(LevelError))
	Expect(messages[8].level).To(Equal(LevelDebug))
	Expect(messages[9].level).To(Equal(LevelError))

	for i, format := range []string{"foo", "bar", "foo", "bar", "foo", "bar", "baz", "baz", "bnk", "bnk"} {
		Expect(messages[i].format).To(Equal(format))
	}
}

func (s *ReplaySuite) TestCheckReplayAddsAttribute(t sweet.T) {
	var (
		shim    = &testShim{}
		clock   = glock.NewMockClock()
		adapter = newReplayShim(adaptShim(shim), clock, LevelDebug, LevelInfo)
	)

	adapter.LogWithFields(LevelDebug, nil, "foo")
	adapter.LogWithFields(LevelInfo, nil, "bar")
	adapter.LogWithFields(LevelDebug, nil, "baz")
	adapter.Replay(LevelError)
	adapter.LogWithFields(LevelDebug, nil, "bonk")

	messages := shim.copy()
	Expect(messages).To(HaveLen(8))
	Expect(messages[0].fields).NotTo(HaveKey(FieldReplay))
	Expect(messages[1].fields).NotTo(HaveKey(FieldReplay))
	Expect(messages[2].fields).NotTo(HaveKey(FieldReplay))
	Expect(messages[3].fields[FieldReplay]).To(Equal(LevelDebug))
	Expect(messages[4].fields[FieldReplay]).To(Equal(LevelInfo))
	Expect(messages[5].fields[FieldReplay]).To(Equal(LevelDebug))
	Expect(messages[6].fields).NotTo(HaveKey(FieldReplay))
	Expect(messages[7].fields[FieldReplay]).To(Equal(LevelDebug))
}

func (s *ReplaySuite) TestCheckSecondReplayAddsAttribute(t sweet.T) {
	var (
		shim    = &testShim{}
		clock   = glock.NewMockClock()
		adapter = newReplayShim(adaptShim(shim), clock, LevelDebug, LevelInfo)
	)

	adapter.LogWithFields(LevelDebug, nil, "foo")
	adapter.LogWithFields(LevelInfo, nil, "bar")
	adapter.Replay(LevelWarning)
	adapter.Replay(LevelError)
	adapter.LogWithFields(LevelDebug, nil, "bnk")

	messages := shim.copy()
	Expect(messages).To(HaveLen(8))
	Expect(messages[0].fields).NotTo(HaveKey(FieldReplay))
	Expect(messages[1].fields).NotTo(HaveKey(FieldReplay))
	Expect(messages[2].fields[FieldReplay]).To(Equal(LevelDebug))
	Expect(messages[3].fields[FieldReplay]).To(Equal(LevelInfo))
	Expect(messages[4].fields[FieldReplay]).To(Equal(LevelDebug))
	Expect(messages[5].fields[FieldReplay]).To(Equal(LevelInfo))
	Expect(messages[6].fields).NotTo(HaveKey(FieldReplay))
	Expect(messages[7].fields[FieldReplay]).To(Equal(LevelDebug))
}
