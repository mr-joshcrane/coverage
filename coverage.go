package coverage

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"golang.org/x/tools/cover"
)

func ExtendCoverage(m *testing.M, name string) int {
	lastArg := os.Args[len(os.Args)-1]
	if lastArg == "ALL" {
		return m.Run()
	}
	path := "/tmp/" + name + ".cover"
	f, err := os.Create(path)
	if err != nil {
		fmt.Println(err)
		return 1
	}
	coverProfile := fmt.Sprintf("-coverprofile=%s", path)
	tags := []string{"test", "-coverpkg=./...", "./...", coverProfile, "-args", "ALL"}
	cmd := exec.Command("go", tags...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(out))
		f.Close()
		return 1
	}
	if strings.Contains(string(out), "FAIL") {
		fmt.Println(string(out))
	}
	profiles, err := cover.ParseProfiles(path)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return 1
	}
	var globalTested, globalTestable int
	for _, profile := range profiles {
		var tested, testable int
		for _, block := range profile.Blocks {
			lineCount := block.EndLine - block.StartLine
			if block.NumStmt > 0 {
				testable += lineCount
			}
			if block.Count > 0 && block.NumStmt > 0 {
				tested += lineCount
			}
		}
		percentageTested := float64(tested) / float64(testable) * 100
		fmt.Printf("%.2f%% - %s\n", percentageTested, profile.FileName)
		globalTested += tested
		globalTestable += testable
	}

	percentageTested := float64(globalTested) / float64(globalTestable) * 100
	fmt.Printf("\nOverall Coverage: %.2f%%\n", percentageTested)

	if err != nil {
		fmt.Println(err)
		f.Close()
		return 1
	}
	f.Close()
	return 0
}
