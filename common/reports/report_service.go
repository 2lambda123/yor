package reports

import (
	"bridgecrewio/yor/common"
	"bridgecrewio/yor/common/structure"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"os"
	"sort"
)

type ReportService struct {
	report Report
}

const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
)

type Report struct {
	ScannedResources int
	NewResources     []structure.IBlock
	UpdatedResources []structure.IBlock
}

var ReportServiceInst *ReportService

func init() {
	ReportServiceInst = &ReportService{}
}

func (r *ReportService) GetReport() *Report {
	return &r.report
}

func (r *ReportService) CreateReport() *Report {
	changesAccumulator := TagChangeAccumulatorInstance
	r.report.ScannedResources = len(changesAccumulator.ScannedBlocks)
	r.report.NewResources = changesAccumulator.NewBlockTraces
	r.report.UpdatedResources = changesAccumulator.UpdatedBlockTraces
	return &r.report
}

func (r *ReportService) PrintToStdout() {
	PrintBanner()
	fmt.Println(colorReset, "Yor Findings Summary")
	fmt.Println(colorReset, "Scanned Resources:\t\t", colorBlue, r.report.ScannedResources)
	fmt.Println(colorReset, "New Resources Traced: \t", colorYellow, len(r.report.NewResources))
	fmt.Println(colorReset, "Updated Resources:\t\t", colorGreen, len(r.report.UpdatedResources))
	fmt.Println()
	if len(r.report.NewResources) > 0 {
		r.printNewResources()
	}
	fmt.Println()
	if len(r.report.UpdatedResources) > 0 {
		r.printUpdatedResources()
	}
}

func PrintBanner() {
	fmt.Printf("%v%vv%v\n", common.YorLogo, colorPurple, common.Version)
}

func (r *ReportService) printUpdatedResources() {
	fmt.Print(colorGreen, fmt.Sprintf("Updated Resource Traces (%v):\n", len(r.report.UpdatedResources)))
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"File", "Resource", "Tag Key", "Old Value", "Updated Value", "Yor ID"})
	table.SetColumnColor(tablewriter.Colors{},
		tablewriter.Colors{},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgRedColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgGreenColor},
		tablewriter.Colors{})

	table.SetRowLine(true)
	table.SetRowSeparator("-")

	for _, block := range r.report.UpdatedResources {
		diff := block.CalculateTagsDiff()
		sort.SliceStable(diff.Added, func(i, j int) bool {
			return diff.Added[i].GetKey() < diff.Added[j].GetKey()
		})
		for _, val := range diff.Added {
			table.Append([]string{block.GetFilePath(), block.GetResourceId(), val.GetKey(), "", val.GetValue(), block.GetTraceId()})
		}
		sort.SliceStable(diff.Updated, func(i, j int) bool {
			return diff.Updated[i].Key < diff.Updated[j].Key
		})
		for _, val := range diff.Updated {
			table.Append([]string{block.GetFilePath(), block.GetResourceId(), val.Key, val.PrevValue, val.NewValue, block.GetTraceId()})
		}
	}
	table.SetAutoMergeCellsByColumnIndex([]int{0, 1, 5})
	table.Render()
}

func (r *ReportService) printNewResources() {
	fmt.Print(colorYellow, fmt.Sprintf("New Resources Traced (%v):\n", len(r.report.NewResources)))
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"File", "Resource", "Owner", "Yor ID"})
	table.SetRowLine(true)
	table.SetRowSeparator("-")
	for _, block := range r.report.NewResources {
		table.Append([]string{block.GetFilePath(), block.GetResourceId(), block.GetNewOwner(), block.GetTraceId()})
	}

	table.Render()
}