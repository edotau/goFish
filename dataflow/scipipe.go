package dataflow

/*
import (
	sp "github.com/scipipe/scipipe"
)

func scipipeDemo() {

}

func MkWorkFlow(name string, threads int) {
	workflow := sp.NewWorkflow("gofishing", 6)

	wget := workflow.NewProc("wget", "wget  http://trackhub.genome.duke.edu/lowelab/edotau/atacseq_nextflow_results_21-05-03/bwa/mergedLibrary/macs/broadPeak/CL12w16-3_atac_R1.mLb.clN_peaks.annotatePeaks.txt")
	move := workflow.NewProc("mv", "mv *annotatePeaks.txt testdata/")
	cat := workflow.NewProc("cat ", "cat testdata/*annotatePeaks.txt | cut -f 2,3,4| less")

	move.In("in").From(wget.Out("out"))
	cat.In("in").From(move.Out("out"))

	// Run workflow
	workflow.Run()
}
*/
