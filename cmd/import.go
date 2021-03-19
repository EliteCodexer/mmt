package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/erdaltsksn/cui"
	"github.com/konradit/mmt/pkg/android"
	"github.com/konradit/mmt/pkg/dji"
	"github.com/konradit/mmt/pkg/gopro"
	"github.com/konradit/mmt/pkg/insta360"
	"github.com/konradit/mmt/pkg/utils"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import media",
	Run: func(cmd *cobra.Command, args []string) {
		input, err := cmd.Flags().GetString("input")
		if err != nil {
			cui.Error("Something went wrong", err)
		}
		output, err := cmd.Flags().GetString("output")
		if err != nil {
			cui.Error("Something went wrong", err)
		}
		camera, err := cmd.Flags().GetString("camera")
		if err != nil {
			cui.Error("Something went wrong", err)
		}
		projectName, err := cmd.Flags().GetString("name")
		if err != nil {
			cui.Error("Something went wrong", err)
		}

		if projectName != "" {
			os.Mkdir(filepath.Join(output, projectName), 0755)
		}

		dateFormat, err := cmd.Flags().GetString("date")
		if err != nil {
			cui.Error("Something went wrong", err)
		}
		bufferSize, err := cmd.Flags().GetInt("buffer")
		if err != nil {
			cui.Error("Something went wrong", err)
		}
		prefix, err := cmd.Flags().GetString("prefix")
		if err != nil {
			cui.Error("Something went wrong", err)
		}

		dateRange, err := cmd.Flags().GetStringArray("range")
		if err != nil {
			cui.Error("Something went wrong", err)
		}

		if camera != "" {
			c, err := utils.CameraGet(camera)
			if err != nil {
				cui.Error("Something went wrong", err)
			}
			r, err := importFromCamera(c, input, filepath.Join(output, projectName), dateFormat, bufferSize, prefix, dateRange)
			if err != nil {
				cui.Error("Something went wrong", err)
			}
			data := [][]string{
				{strconv.Itoa(r.FilesImported), strconv.Itoa(len(r.FilesNotImported)), strconv.Itoa(len(r.Errors))},
			}
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Files Imported", "Files Skipped", "Errors"})

			for _, v := range data {
				table.Append(v)
			}
			table.Render() // Send output

			if len(r.Errors) != 0 {
				fmt.Println("Errors: ")
				for _, error := range r.Errors {
					fmt.Println(">> " + error.Error())
				}
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(importCmd)

	importCmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose")
	importCmd.Flags().StringP("input", "i", "", "Input directory for root, eg: E:\\")
	importCmd.Flags().StringP("output", "o", "", "Output directory for sorted media")
	importCmd.Flags().StringP("name", "n", "", "Project name")
	importCmd.Flags().StringP("camera", "c", "", "Camera type")
	importCmd.Flags().StringP("date", "d", "dd-mm-yyyy", "Date format, dd-mm-yyyy by default")
	importCmd.Flags().IntP("buffer", "b", 1000, "Buffer size for copying, default is 1000 bytes")
	importCmd.Flags().StringP("prefix", "p", "", "Prefix for each file, pass `cameraname` to prepend the camera name (eg: Hero9 Black)")
	importCmd.Flags().StringArray("range", []string{}, "A date range, eg: 01-05-2020,05-05-2020 -- also accepted: `today`, `yesterday`, `week`")

	for _, item := range []string{
		"input", "output",
	} {
		importCmd.MarkFlagRequired(item)
	}

}

func importFromCamera(c utils.Camera, input string, output string, dateFormat string, bufferSize int, prefix string, dateRange []string) (*utils.Result, error) {
	switch c {
	case utils.GoPro:
		return gopro.Import(input, output, dateFormat, bufferSize, prefix, dateRange)
	case utils.DJI:
		return dji.Import(input, output, dateFormat, bufferSize, prefix, dateRange)
	case utils.Insta360:
		return insta360.Import(input, output, dateFormat, bufferSize, prefix, dateRange)
	case utils.Android:
		return android.Import(input, output, dateFormat, bufferSize, prefix, dateRange)
	default:
		return nil, errors.New("Unsupported camera!")
	}
}
