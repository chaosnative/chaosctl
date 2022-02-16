package hubgen

import (
	"errors"
	"fmt"
	"github.com/chaosnative/chaosctl/pkg/utils"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var importCmd = &cobra.Command{
	Use:   "generate",
	Short: "Imports the local ChaosHub with it's path",
	Long:  "Imports the local ChaosHub with it's path",
	Run: func(cmd *cobra.Command, args []string) {
		//User input to generate hub with current changes or add more charts
		var confirm string
		generatedChart := make(map[string][]string)

		//CLI input for the import-path of local chaos hub
		hubPath, err := cmd.Flags().GetString("import-path")
		utils.PrintError(err)
		if hubPath[len(hubPath)-1:] == "/" {
			hubPath = hubPath + "charts/"
		} else {
			hubPath = hubPath + "/charts/"
		}

		//CLI input for the export-path of generated hub
		exportPath, err := cmd.Flags().GetString("export-path")
		utils.PrintError(err)

		//Read the directories present in the local ChaosHub
		names, err := readDir(hubPath)
		if err != nil {
			panic(err)
		}

	LOOP:
		generatedChart = ListChartsAndExperiments(names, hubPath, generatedChart)
		utils.White_B.Print("\n📦 Experiments added successfully, do you want to add more charts and experiments? [Y/N] \n " +
			"Y to add new chart and experiments / N to generate the current configured ChaosHub: ")
		fmt.Scanln(&confirm)
		if strings.ToLower(confirm) == "y" {
			goto LOOP
		} else if strings.ToLower(confirm) == "n" {
			err = generateChaosHub(generatedChart, hubPath, exportPath)
			if err != nil {
				utils.PrintError(err)
			} else {
				utils.White_B.Print("\n🎉 ChaosHub generated successfully at " + exportPath)
			}
		} else {
			utils.PrintError(errors.New("invalid command, try again"))
		}
	},
}

func ListChartsAndExperiments(names []string, hubPath string, generatedCharts map[string][]string) map[string][]string {
	var (
		chartName   string
		err         error
		expIndex    string
		expIndexNum []int
		experiment  []string
	)

	//Select the available charts
	utils.White_B.Print("\n📝 Select the charts: ")
	chartPrompt := promptui.Select{
		Label: "Available charts",
		Items: names,
		Size:  len(names),
	}
	_, chartName, err = chartPrompt.Run()
	if err != nil {
		utils.Red.Println(errors.New("Prompt err:" + err.Error()))
		os.Exit(1)
	}

	//Display the available experiments from the selected charts
	utils.White_B.Print("\n☢ ️ Experiments in " + chartName + " category are: \n")
	expPath := hubPath + chartName
	expNames, err := readDir(expPath)
	if err != nil {
		utils.PrintError(err)
	}
	if len(expNames) == 0 {
		utils.Red.Println(errors.New("No experiments available in " + chartName + "\n"))
		return generatedCharts
	}
	for index, expName := range expNames {
		fmt.Printf("%d.  %s\n", index+1, expName)
	}

	//Select the list of experiments with comma separation
	utils.White_B.Print("\n✅ Select experiments from " + chartName + " charts: Range[1- " + fmt.Sprintf("%d", len(expNames)) + "] ex: 1,4,7: ")
	fmt.Scanln(&expIndex)

	//Split the user input string with commas and push the same in expIndex array
	res := strings.Split(expIndex, ",")
	for _, result := range res {
		num, err := strconv.Atoi(strings.Trim(result, " "))
		if err != nil {
			utils.PrintError(err)
		}
		if num <= len(expNames) {
			expIndexNum = append(expIndexNum, num-1)
		} else {
			utils.PrintError(errors.New("invalid input"))
		}
	}

	//Updated experiments array as per the experiment index
	for _, val := range expIndexNum {
		experiment = append(experiment, expNames[val])
	}

	//Map the charts according to the selected experiments
	generatedCharts[chartName] = experiment
	return generatedCharts
}

func generateChaosHub(generatedCharts map[string][]string, importPath string, exportPath string) error {
	updatedExportPath := exportPath + "/chaos-charts/charts"
	sourceInfo, err := os.Stat(importPath)
	if err != nil {
		return err
	}

	err = os.MkdirAll(updatedExportPath, sourceInfo.Mode())
	if err != nil {
		return err
	}

	for chart := range generatedCharts {
		err = os.MkdirAll(updatedExportPath+"/"+chart, sourceInfo.Mode())
		if err != nil {
			fmt.Errorf("Error while creating " + chart + " directory")
		}
		iconCmd := exec.Command("cp", "--recursive", importPath+"/"+chart+"/icons", updatedExportPath+"/"+chart+"/icons")
		iconCmd.Run()

		for _, experiment := range generatedCharts[chart] {
			cmd := exec.Command("cp", "--recursive", importPath+"/"+chart+"/"+experiment, updatedExportPath+"/"+chart+"/"+experiment)
			cmd.Run()
		}
	}
	return err
}

func readDir(path string) ([]string, error) {
	var names []string
	files, err := ioutil.ReadDir(path)
	if err != nil {
		utils.PrintError(err)
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() && file.Name() != "icons" {
			names = append(names, file.Name())
		}
	}
	return names, nil
}

func init() {
	HubgenCmd.AddCommand(importCmd)
	importCmd.Flags().String("import-path", "", "Hub Path of local ChaosHub")
	importCmd.Flags().String("export-path", "", "Path to save the generated Chaos Hub")
}
