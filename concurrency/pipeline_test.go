package concurrency

import (
	"io"
	"strings"
	"testing"
)

func TestSettingUpPipeline(t *testing.T) {
	flow := New()
	// tasks
	clean := flow.MustRegister(taskClean())
	install := flow.MustRegister(taskInstall())
	build := flow.MustRegister(taskBuild())
	fmt := flow.MustRegister(taskFmt())
	lint := flow.MustRegister(taskLint())
	test := flow.MustRegister(taskTest())
	modTidy := flow.MustRegister(taskModTidy())
	diff := flow.MustRegister(taskDiff())

	// pipeline
	all := flow.MustRegister(Task{
		Name:        "all",
		Description: "build pipeline",
		Dependencies: Deps{
			clean,
			install,
			build,
			fmt,
			lint,
			test,
			modTidy,
			diff,
		},
	})
	flow.DefaultTask = all
	flow.Params.SetBool("ci", false)
	//flow.Main()
}

const toolsDir = "tools"

func taskClean() Task {
	return Task{
		Name:        "clean",
		Description: "remove git ignored files",
		Command:     Exec("git", "clean", "-fX"),
	}
}

func taskInstall() Task {
	return Task{
		Name:        "install",
		Description: "install build tools",
		Command: func(tf *TF) {
			installFmt := tf.Cmd("go", "install", "mvdan.cc/gofumpt/gofumports")
			installFmt.Dir = toolsDir
			if err := installFmt.Run(); err != nil {
				tf.Errorf("go install gofumports: %v", err)
			}

			installLint := tf.Cmd("go", "install", "github.com/golangci/golangci-lint/cmd/golangci-lint")
			installLint.Dir = toolsDir
			if err := installLint.Run(); err != nil {
				tf.Errorf("go install golangci-lint: %v", err)
			}
		},
	}
}

func taskBuild() Task {
	return Task{
		Name:        "build",
		Description: "go build",
		Command:     Exec("go", "build", "./..."),
	}
}

func taskFmt() Task {
	return Task{
		Name:        "fmt",
		Description: "gofumports",
		Command: func(tf *TF) {
			tf.Cmd("gofumports", strings.Split("-l -w -local github.com/pellared/taskflow .", " ")...).Run() //nolint // it is OK if it returns error
		},
	}
}

func taskLint() Task {
	return Task{
		Name:        "lint",
		Description: "golangci-lint",
		Command:     Exec("golangci-lint", "run"),
	}
}

func taskTest() Task {
	return Task{
		Name:        "test",
		Description: "go test with race detector and code covarage",
		Command:     Exec("go", "test", "-race", "-covermode=atomic", "-coverprofile=coverage.out", "./..."),
	}
}

func taskModTidy() Task {
	return Task{
		Name:        "mod-tidy",
		Description: "go mod tidy",
		Command: func(tf *TF) {
			if err := tf.Cmd("go", "mod", "tidy").Run(); err != nil {
				tf.Errorf("go mod tidy: %v", err)
			}

			toolsModTidy := tf.Cmd("go", "mod", "tidy")
			toolsModTidy.Dir = toolsDir
			if err := toolsModTidy.Run(); err != nil {
				tf.Errorf("go mod tidy: %v", err)
			}
		},
	}
}

func taskDiff() Task {
	return Task{
		Name:        "diff",
		Description: "git diff",
		Command: func(tf *TF) {
			if !tf.Params().Bool("ci") {
				tf.Skip("ci param is not set, skipping")
			}

			if err := tf.Cmd("git", "diff", "--exit-code").Run(); err != nil {
				tf.Errorf("git diff: %v", err)
			}

			cmd := tf.Cmd("git", "status", "--porcelain")
			sb := &strings.Builder{}
			cmd.Stdout = io.MultiWriter(tf.Output(), sb)
			if err := cmd.Run(); err != nil {
				tf.Errorf("git status --porcelain: %v", err)
			}
			if sb.Len() > 0 {
				tf.Error("git status --porcelain returned output")
			}
		},
	}
}

/*
import (
	"fmt"
	"github.com/goFish/simpleio"
	"github.com/goFish/vcf"
	"github.com/vertgenlab/gonomics/fileio"
	"testing"
	"time"
)

func Test_pipeline_Pipe(t *testing.T) {
	type fields struct {
		dataC     chan interface{}
		errC      chan error
		executors []Executor
	}
	type args struct {
		executor Executor
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		check  func(f fields) bool
	}{
		{
			name: "should success adding executor to executors",
			fields: fields{
				executors: []Executor{},
			},
			args: args{
				executor: func(in interface{}) interface{} {
					return 1
				},
			},
			check: func(f fields) bool {
				return len(f.executors) == 1
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &pipeline{
				dataC: tt.fields.dataC,
				//errC:      tt.fields.errC,
				executors: tt.fields.executors,
			}
			p.Pipe(tt.args.executor)
			if !tt.check(fields{p.dataC, nil, p.executors}) {
				t.Errorf("pipeline.Pipe() not run as expected")
			}
		})
	}
}

func BenchmarkWithPipelineModule(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	ReadPipe(b)

}

type VcfSlice []vcf.Vcf

func ReadPipe(b *testing.B) {
	outC := New(func(inC chan interface{}) {
		defer close(inC)
		for i := 0; i < b.N; i++ {
			inC <- i
		}
	}).
		Pipe(func(in interface{}) interface{} {
			var vcfs VcfSlice = vcf.ReadVcfs("testdata/small.vcf.gz")
			return vcfs
		}).
		Pipe(func(in interface{}) interface{} {
			WriteTest("testdata/rewrite.vcf", in.(VcfSlice))
			return nil
		}).
		Merge()

	for range outC {
		// Do nothing, just for  drain out channel
	}
}

func Read(filename string) []string {
	var ans []string
	reader := simpleio.NewReader(filename)
	for i, err := simpleio.ReadLine(reader); !err; i, err = simpleio.ReadLine(reader) {
		ans = append(ans, i.String())
	}
	return ans
}

func WriteTest(filename string, vcfs VcfSlice) {
	output := fileio.EasyCreate(filename)
	var err error

	for _, i := range vcfs {
		_, err = fmt.Fprintf(output, "%s\n", vcf.ToString(&i))
		simpleio.ErrorHandle(err)
	}

}

func BenchmarkWithoutPipelineModule(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		simpleio.ReadWriteVcfs("testdata/small.vcf.gz", "testdata/benchmark.vcf")
		//fmt.Printf("Number of records: %d\n", len(v))
	}
}

func FileioRead(filename string) []string {
	reader := fileio.EasyOpen(filename)
	var ans []string
	for i, err := fileio.EasyNextLine(reader); !err; i, err = fileio.EasyNextLine(reader) {
		ans = append(ans, i)
	}
	return ans
}

func multiplyTwo(v int) int {
	time.Sleep(100 * time.Millisecond)
	return v * 2
}

func square(v int) int {
	time.Sleep(200 * time.Millisecond)
	return v * v
}

func addQuoute(v int) string {
	time.Sleep(100 * time.Millisecond)
	return fmt.Sprintf("'%d'", v)
}

func addFoo(v string) string {
	time.Sleep(200 * time.Millisecond)
	return fmt.Sprintf("%s - Foo", v)
}*/
