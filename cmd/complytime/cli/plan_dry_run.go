// SPDX-License-Identifier: Apache-2.0
package cli

import (
	"fmt"
	"os"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"gopkg.in/yaml.v3"

	"github.com/complytime/complytime/internal/complytime"
)

// PlanData sets up the yaml mapping type for writing to config file.
// Formats testdata as go struct.
type PlanData struct {
	FrameworkID string                  `yaml:"assessment_plan"`
	Components  []string                `yaml:"components"`
	Controls    [][]oscalTypes.Property `yaml:"controls"`
}

// planDryRun leverages the PlanData structure to populate tailoring config.
// The config is written to stdout.
func planDryRun(appDir complytime.ApplicationDirectory, frameworkId string, cds []oscalTypes.ComponentDefinition, ap *oscalTypes.AssessmentPlan) {
	basePlanData := PlanData{
		FrameworkID: frameworkId,
		Components:  make([]string, 0),
		Controls:    [][]oscalTypes.Property{},
	}
	if cds == nil {
		fmt.Fprintln(os.Stderr, "no component definitions found")
		return
	}
	for _, componentDef := range cds {
		if componentDef.Components == nil {
			continue
		}
		for _, component := range *componentDef.Components {
			if component.ControlImplementations == nil {
				continue
			}
			for _, ci := range *component.ControlImplementations {
				if ci.Props == nil {
					continue
				}
				// Rule_Id
				// Check_Id
				// UUID
				// [
				//   ["Rule1", "Checkid1", "123123-4234-1233245-12343245"]
				//   ["Rule1", "Checkid1", "123123-4234-1233245-12343245"]
				//   ["Rule1", "Checkid1", "123123-4234-1233245-12343245"]
				// ]
				controlDescription := []oscalTypes.Property{}
				if len(controlDescription) != 0 {
					basePlanData.Controls = append(basePlanData.Controls, *ci.Props)
				}
			}
		}
	}

	//fmt.Printf("%v\n", cds)
	//fmt.Printf("%v\n", ap)

	out, err := yaml.Marshal(&basePlanData)
	if err != nil {
		fmt.Println("error marshalling yaml content: ", err)
	}
	fmt.Println(string(out))
}
