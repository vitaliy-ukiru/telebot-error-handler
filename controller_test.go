package telerr

import (
	"testing"

	"github.com/stretchr/testify/assert"
	tele "gopkg.in/telebot.v3"
)

var _ defaultCaseCatcher = (*Default)(nil)

func TestErrorController_setupCases(t *testing.T) {

	type testCase struct {
		name            string
		args            []Catcher
		want            bool
		wantDefaultCase assert.ValueAssertionFunc
	}
	tests := []testCase{
		{
			name:            "without default case",
			args:            []Catcher{},
			wantDefaultCase: assert.Nil,
		},
		{
			name:            "with default case",
			args:            []Catcher{Default(Ignore)},
			want:            true,
			wantDefaultCase: assert.NotNil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ec := &ErrorController{}

			assert.Equalf(t, tt.want, ec.setupCases(tt.args), "setupCases(%v)", tt.args)
			if !tt.wantDefaultCase(t, ec.defaultCase, "default case") {
				return
			}

		})
	}
}

func TestErrorController_process(t *testing.T) {
	//type fields struct {
	//	cases       []Catcher
	//	defaultCase Handler
	//}
	type args struct {
		err         error
		ctx         tele.Context
		catcherCall bool
	}
	//tests := []struct {
	//	name   string
	//	fields fields
	//	args   args
	//	want   bool
	//}{
	//	// TODO: Add test cases.
	//}
	//for _, tt := range tests {
	//	t.Run(tt.name, func(t *testing.T) {
	//		ec := &ErrorController{
	//			cases:       tt.fields.cases,
	//			defaultCase: tt.fields.defaultCase,
	//		}
	//		assert.Equalf(t, tt.want, ec.process(tt.args.err, tt.args.ctx, tt.args.catcherCall), "process(%v, %v, %v)", tt.args.err, tt.args.ctx, tt.args.catcherCall)
	//	})
	//}

	type testCase struct {
		name              string
		args              args
		caseMatches       []bool
		wantStack         []int
		shouldDefaultCase bool
		want              bool
	}

	tests := []testCase{
		{
			name:              "cases matches",
			caseMatches:       []bool{false, false, true, false},
			wantStack:         []int{0, 1, 2},
			shouldDefaultCase: false,
			want:              true,
		},
		{
			name:              "match default",
			caseMatches:       []bool{false, false},
			wantStack:         []int{0, 1},
			shouldDefaultCase: true,
			want:              true,
		},
		{
			name:              "catcher call",
			caseMatches:       []bool{false},
			wantStack:         []int{0},
			args:              args{catcherCall: true},
			shouldDefaultCase: false,
			want:              false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stack := new([]int)
			*stack = make([]int, 0, cap(tt.wantStack))

			cases := make([]Catcher, 0, len(tt.caseMatches))
			for i, match := range tt.caseMatches {
				i := i
				match := match

				cases = append(cases, &testCatcher{id: i, match: match, stack: stack})
			}

			ec := &ErrorController{
				cases: cases,
				defaultCase: func(_ error, _ tele.Context) {
					assert.True(t, tt.shouldDefaultCase, "default case execute")
				},
			}

			assert.Equalf(t, tt.want, ec.process(tt.args.err, tt.args.ctx, tt.args.catcherCall), "process(%v, %v, %v)", tt.args.err, tt.args.ctx, tt.args.catcherCall)
			assert.Equalf(t, tt.wantStack, *stack, "call stack")
		})
	}

}

type testCatcher struct {
	stack *[]int
	id    int
	match bool
}

func (t testCatcher) Catch(_ error, _ tele.Context) bool {
	*t.stack = append(*t.stack, t.id)
	return t.match

}
