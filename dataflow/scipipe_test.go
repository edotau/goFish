package dataflow

/*func TestSicipipeFeatures(t *testing.T) {
	// Init workflow and max concurrent tasks
	wf := sp.NewWorkflow("hello_world", 4)

	// Initialize processes, and file extensions
	hello := wf.NewProc("hello", "echo 'Hello ' > testdata/hello.txt")
	world := wf.NewProc("world", "echo $(cat World) >> testdata/hello.txt")

	// Define data flow
	world.In("in").From(hello.Out("out"))

	// Run workflow
	wf.Run()
	//MkWorkFlow("download atacseq data!", 2)
}*/
/*
func TestBasicRun(t *testing.T) {
	//initTestLogs()
	wf := sp.NewWorkflow("TestBasicRunWf", 16)

	p1 := wf.NewProc("p1", "echo foo > {o:foo}")
	assertIsType(t, p1.Out("foo"), sp.NewOutPort("foo"))
	p1.SetOutFunc("foo", func(t *sp.Task) string { return "foo.txt" })

	p2 := wf.NewProc("p2", "sed 's/foo/bar/g' {i:foo} > {o:bar}")

	assertIsType(t, p2.In("foo"), sp.NewInPort("foo"))
	assertIsType(t, p2.Out("bar"), sp.NewOutPort("bar"))
	p2.SetOut("bar", "{i:foo}.bar.txt")

	p2.In("foo").From(p1.Out("foo"))

	//assertIsType(t, p2.In("foo"), constNewInPort("foo"))
	assertIsType(t, p2.Out("bar"), sp.NewOutPort("bar"))

	wf.Run()

}
*/
