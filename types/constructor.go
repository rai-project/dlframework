package types

import (
	"fmt"

	"github.com/chewxy/hm"
)

type Constructor struct {
	head hm.Type
	ts   []hm.Type
}

func NewConstructor(head hm.Type, ts ...hm.Type) *Constructor {
	return &Constructor{
		head: head,
		ts:   ts,
	}
}

func (t *Constructor) Name() string {
	return t.head.String()
}

func (t *Constructor) String() string {
	return fmt.Sprintf("%v(%#v)", t.head, t.ts)
}

func (t *Constructor) Apply(sub hm.Subs) hm.Substitutable {
	ts := make([]hm.Type, len(t.ts))
	for i, v := range t.ts {
		ts[i] = v.Apply(sub).(hm.Type)
	}
	head := t.head.Apply(sub).(hm.Type)
	return NewConstructor(head, ts...)
}

func (t *Constructor) FreeTypeVar() hm.TypeVarSet {
	tvs := t.head.FreeTypeVar()
	for _, v := range t.ts {
		tvs = v.FreeTypeVar().Union(tvs)
	}
	return tvs
}

func (t *Constructor) Normalize(k, v hm.TypeVarSet) (hm.Type, error) {
	ts := make([]hm.Type, len(t.ts))
	var err error
	for i, tt := range t.ts {
		if ts[i], err = tt.Normalize(k, v); err != nil {
			return nil, err
		}
	}
	head, err := t.head.Normalize(k, v)
	if err != nil {
		return nil, err
	}
	return NewConstructor(head, ts...), nil
}

func (t *Constructor) Types() hm.Types {
	ts := hm.BorrowTypes(len(t.ts) + 1)
	ts[0] = t.head
	copy(ts[1:], t.ts)
	return ts
}

func (t *Constructor) Eq(other hm.Type) bool {
	if ot, ok := other.(*Constructor); ok {
		if !t.head.Eq(ot.head) {
			return false
		}
		if len(ot.ts) != len(t.ts) {
			return false
		}
		for i, v := range t.ts {
			if !v.Eq(ot.ts[i]) {
				return false
			}
		}
		return true
	}
	return false
}

func (t *Constructor) Format(state fmt.State, c rune) {
	fmt.Fprintf(state, "%v", t.head)
	state.Write([]byte("("))
	for i, v := range t.ts {
		if i < len(t.ts)-1 {
			fmt.Fprintf(state, "%v, ", v)
		} else {
			fmt.Fprintf(state, "%v)", v)
		}
	}
}

// Clone implements Cloner
func (t *Constructor) Clone() interface{} {
	retVal := new(Constructor)
	ts := hm.BorrowTypes(len(t.ts))
	for i, tt := range t.ts {
		if c, ok := tt.(hm.Cloner); ok {
			ts[i] = c.Clone().(hm.Type)
		} else {
			ts[i] = tt
		}
	}
	retVal.ts = ts
	if c, ok := t.head.(hm.Cloner); ok {
		retVal.head = c.Clone().(hm.Type)
	} else {
		retVal.head = t.head
	}

	return retVal
}
