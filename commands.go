package picol

import (
	"fmt"
	"strconv"
)

func ArityErr(i *Interp, name string, argv []string) error {
	return fmt.Errorf("Wrong number of args for %s %s", name, argv)
}

func CommandMath(i *Interp, argv []string, pd interface{}) (string, error) {
	if len(argv) != 3 {
		return "", ArityErr(i, argv[0], argv)
	}
	a, _ := strconv.Atoi(argv[1])
	b, _ := strconv.Atoi(argv[2])
	var c int
	switch {
	case argv[0] == "+":
		c = a + b
	case argv[0] == "-":
		c = a - b
	case argv[0] == "*":
		c = a * b
	case argv[0] == "/":
		c = a / b
	case argv[0] == ">":
		if a > b {
			c = 1
		}
	case argv[0] == ">=":
		if a >= b {
			c = 1
		}
	case argv[0] == "<":
		if a < b {
			c = 1
		}
	case argv[0] == "<=":
		if a <= b {
			c = 1
		}
	case argv[0] == "==":
		if a == b {
			c = 1
		}
	case argv[0] == "!=":
		if a != b {
			c = 1
		}
	default: // FIXME I hate warnings
		c = 0
	}
	return fmt.Sprintf("%d", c), nil
}

func CommandSet(i *Interp, argv []string, pd interface{}) (string, error) {
	if len(argv) != 3 {
		return "", ArityErr(i, argv[0], argv)
	}
	i.SetVar(argv[1], argv[2])
	return argv[2], nil
}

func CommandIf(i *Interp, argv []string, pd interface{}) (string, error) {
	if len(argv) != 3 && len(argv) != 5 {
		return "", ArityErr(i, argv[0], argv)
	}

	result, err := i.Eval(argv[1])
	if err != nil {
		return "", err
	}

	if r, _ := strconv.Atoi(result); r != 0 {
		return i.Eval(argv[2])
	} else if len(argv) == 5 {
		return i.Eval(argv[4])
	}

	return result, nil
}

func CommandWhile(i *Interp, argv []string, pd interface{}) (string, error) {
	if len(argv) != 3 {
		return "", ArityErr(i, argv[0], argv)
	}

	for {
		result, err := i.Eval(argv[1])
		if err != nil {
			return "", err
		}
		if r, _ := strconv.Atoi(result); r != 0 {
			result, err := i.Eval(argv[2])
			switch err {
			case PICOL_CONTINUE, nil:
				//pass
			case PICOL_BREAK:
				return result, nil
			default:
				return result, err
			}
		} else {
			return result, nil
		}
	}
}

func CommandRetCodes(i *Interp, argv []string, pd interface{}) (string, error) {
	if len(argv) != 1 {
		return "", ArityErr(i, argv[0], argv)
	}
	switch argv[0] {
	case "break":
		return "", PICOL_BREAK
	case "continue":
		return "", PICOL_CONTINUE
	}
	return "", nil
}

func DropCallFrame(i *Interp) {
	i.callframe = i.callframe.parent
}

func CommandCallProc(i *Interp, argv []string, pd interface{}) (string, error) {
	var x []string

	if pd, ok := pd.([]string); ok {
		x = pd
	} else {
		return "", nil
	}

	alist := x[0]
	body := x[1]
	p := alist[:]
	arity := 0

	done := false
	i.callframe = &CallFrame{vars: make(map[string]Var), parent: i.callframe}
	defer DropCallFrame(i) // remove the called proc callframe

	for {
		start := p
		for len(p) != 0 && p[0] != ' ' {
			p = p[1:]
		}
		if len(p) != 0 && p == start {
			p = p[1:]
			continue
		}

		if p == start {
			break
		}
		if len(p) == 0 {
			done = true
		} else {
			p = p[1:1]
		}
		arity++
		if arity > len(argv)-1 {
			return "", fmt.Errorf("Proc '%s' called with wrong arg num", argv[0])
		}
		i.SetVar(start, argv[arity])
		if len(p) != 0 {
			p = p[1:]
		}
		if done {
			break
		}
	}

	if arity != len(argv)-1 {
		return "", fmt.Errorf("Proc '%s' called with wrong arg num", argv[0])
	}

	result, err := i.Eval(body)
	if err == PICOL_RETURN {
		err = nil
	}
	return result, err
}

func CommandProc(i *Interp, argv []string, pd interface{}) (string, error) {
	if len(argv) != 4 {
		return "", ArityErr(i, argv[0], argv)
	}
	return "", i.RegisterCommand(argv[1], CommandCallProc, []string{argv[2], argv[3]})
}

func CommandReturn(i *Interp, argv []string, pd interface{}) (string, error) {
	if len(argv) != 1 && len(argv) != 2 {
		return "", ArityErr(i, argv[0], argv)
	}
	var r string
	if len(argv) == 2 {
		r = argv[1]
	}
	return r, PICOL_RETURN
}

func CommandError(i *Interp, argv []string, pd interface{}) (string, error) {
	if len(argv) != 1 && len(argv) != 2 {
		return "", ArityErr(i, argv[0], argv)
	}
	return "", fmt.Errorf(argv[1])
}

func (i *Interp) RegisterCoreCommands() {
	name := [...]string{"+", "-", "*", "/", ">", ">=", "<", "<=", "==", "!="}
	for _, n := range name {
		i.RegisterCommand(n, CommandMath, nil)
	}
	i.RegisterCommand("set", CommandSet, nil)
	i.RegisterCommand("if", CommandIf, nil)
	i.RegisterCommand("while", CommandWhile, nil)
	i.RegisterCommand("break", CommandRetCodes, nil)
	i.RegisterCommand("continue", CommandRetCodes, nil)
	i.RegisterCommand("proc", CommandProc, nil)
	i.RegisterCommand("return", CommandReturn, nil)
	i.RegisterCommand("error", CommandReturn, nil)
}