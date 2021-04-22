package concurrency

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
