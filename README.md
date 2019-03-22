# picol.go

Original http://oldblog.antirez.com/post/picol.html

Sample use:
```golang
func CommandPuts(i *picol.Interp, argv []string, pd interface{}) (string, error) {
	if len(argv) != 2 {
		return "", fmt.Errorf("Wrong number of args for %s %s", argv[0], argv)
	}
	fmt.Println(argv[1])
	return "", nil
}
...
	interp := picol.InitInterp()
	// add core functions
	interp.RegisterCoreCommands()
	// add user function
	interp.RegisterCommand("puts", CommandPuts, nil)
	// eval
	result, err := interp.Eval(string(buf))
	if err != nil {
		fmt.Println("ERROR", err, result)
	} else {
		fmt.Println(result)
	}
```

## UPDATE

I forked this forked this project because there was a bug in the parser.go file that prevented it from compiling.
Also, the project name also janked the golang tools.

In the sample code above ```CommandPuts``` has a final param ```pd interface{}```. The value of this param is whatever value the user passed in when registering. I do not particularly see what the benefit is unless there is some sort of fork.

### PrivData

In the callframe there is privData that is set when calling a proc()
