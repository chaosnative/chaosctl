package hubgen

import (
	"archive/zip"
	"errors"
	"fmt"
	"github.com/chaosnative/chaosctl/pkg/types"
	"github.com/chaosnative/chaosctl/pkg/utils"
	"github.com/ghodss/yaml"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Imports the local ChaosHub with it's path and generates the configured ChaosHub",
	Long:  "Imports the local ChaosHub with it's path and generates the configured ChaosHub",
	Run: func(cmd *cobra.Command, args []string) {
		//User input to generate hub with current changes or add more charts
		var confirm string
		generatedChart := make(map[string][]string)

		//CLI input for the generated hub name
		hubName, err := cmd.Flags().GetString("hubname")
		utils.PrintError(err)

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
		if exportPath[len(exportPath)-1:] != "/" {
			exportPath = exportPath + "/" + hubName
		}

		//Read the directories present in the local ChaosHub
		names, err := readDir(hubPath)
		if err != nil {
			panic(err)
		}

	LOOP:
		generatedChart = listChartsAndExperiments(names, hubPath, generatedChart)
		fmt.Printf("\n")

		confirmPromptLabel := "📦 Experiments added successfully, select Y to generate the ChaosHub or N to add more charts"
		promptConfirm := promptui.Prompt{
			Label:     confirmPromptLabel,
			IsConfirm: true,
		}
		confirm, _ = promptConfirm.Run()
		if strings.ToLower(confirm) == "n" {
			goto LOOP
		} else if strings.ToLower(confirm) == "y" {
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

//copyDir is used to copy directory and its content using source and destination path
func copyDir(src string, dest string) error {
	f, err := os.Open(src)
	if err != nil {
		return err
	}

	file, err := f.Stat()
	if err != nil {
		return err
	}
	if !file.IsDir() {
		return fmt.Errorf("Source " + file.Name() + " is not a directory!")
	}

	err = os.Mkdir(dest, file.Mode())
	if err != nil {
		return err
	}

	files, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	for _, f := range files {
		if f.IsDir() {
			err = copyDir(src+"/"+f.Name(), dest+"/"+f.Name())
			if err != nil {
				return err
			}
		}

		if !f.IsDir() {
			content, err := ioutil.ReadFile(src + "/" + f.Name())
			if err != nil {
				return err

			}
			err = ioutil.WriteFile(dest+"/"+f.Name(), content, f.Mode())
			if err != nil {
				return err

			}

		}

	}
	return nil
}

//listChartsAndExperiments is used to list the charts and its related experiments
func listChartsAndExperiments(names []string, hubPath string, generatedCharts map[string][]string) map[string][]string {
	var (
		chartName   string
		err         error
		expIndex    string
		expIndexNum []int
		experiment  []string
	)

	//Select the available charts
	utils.White_B.Print("\n📝 Select the charts: ")
	prompt := promptui.Select{
		Label: "Available charts",
		Items: names,
		Size:  len(names),
	}
	_, chartName, err = prompt.Run()
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
	fmt.Printf("\n")
	//Select the list of experiments with comma separation
	expPromptLabel := "✅ Select experiments from " + chartName + " charts: Range[1- " + fmt.Sprintf("%d", len(expNames)) + "] ex: 1,4,7 "
	promptExpIndex := promptui.Prompt{
		Label: expPromptLabel,
	}
	expIndex, _ = promptExpIndex.Run()

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

//generateChaosHub is used to generate the configured ChaosHub
func generateChaosHub(generatedCharts map[string][]string, importPath string, exportPath string) error {
	updatedExportPath := exportPath + "/charts"
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
			fmt.Println(fmt.Errorf("Error while creating " + chart + " directory"))
		}
		err = generateCSV(updatedExportPath, chart, generatedCharts)
		if err != nil {
			fmt.Println(fmt.Errorf("error while creating CSV file"))
		}
		err = copyDir(importPath+"/"+chart+"/icons", updatedExportPath+"/"+chart+"/icons")
		if err != nil {
			fmt.Println(err.Error())
		}
		for _, experiment := range generatedCharts[chart] {
			err = copyDir(importPath+"/"+chart+"/"+experiment, updatedExportPath+"/"+chart+"/"+experiment)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}
	if err := zipSource(exportPath); err != nil {
		log.Fatal(err)
	}
	return err
}

//generateCSV is used to generate the CSV for a specific chart
func generateCSV(updatedExportPath string, chart string, generatedCharts map[string][]string) error {
	var experiments []string
	for _, experiment := range generatedCharts[chart] {
		experiments = append(experiments, experiment)
	}
	CSV := &types.CSV{
		ApiVersion: "litmuchaos.io/v1alpha1",
		Kind:       "ChartServiceVersion",
		Metadata:   types.Metadata{Name: chart},
		Spec:       types.Spec{DisplayName: chart + "chaos", Experiments: experiments},
	}
	CSVYaml, err := yaml.Marshal(CSV)
	if err != nil {
		fmt.Errorf(err.Error())
	}
	f, err := os.Create(updatedExportPath + "/" + chart + "/" + chart + ".chartserviceversion.yaml")
	if err != nil {
		fmt.Println(fmt.Errorf("error while creating CSV file"))
	}
	defer f.Close()

	_, err = f.Write(CSVYaml)
	if err != nil {
		return err
	}
	return nil
}

//readDir is used to read the directory from a specific path
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

//zipSource is used to zip a directory from a specific path
func zipSource(zipPath string) error {
	//Create a zip of the specific dir
	f, err := os.Create(zipPath + ".zip")
	if err != nil {
		return err
	}
	defer f.Close()

	writer := zip.NewWriter(f)
	defer writer.Close()

	//walk through the filepath and call each folder and file in the filepath
	return filepath.Walk(zipPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Method = zip.Deflate
		header.Name, err = filepath.Rel(filepath.Dir(zipPath), path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			header.Name += "/"
		}
		headerWriter, err := writer.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(headerWriter, f)
		return err
	})
}

func init() {
	HubgenCmd.AddCommand(generateCmd)
	generateCmd.Flags().String("hubname", "", "Name of the generated ChaosHub")
	generateCmd.Flags().String("import-path", "", "Hub Path of local ChaosHub")
	generateCmd.Flags().String("export-path", "", "Path to save the generated Chaos Hub")
}
