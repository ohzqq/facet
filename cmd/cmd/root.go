package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ohzqq/facet"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile   string
	dataFiles []string
	idx       = &facet.Index{}
)

var rootCmd = &cobra.Command{
	Use:   "facet",
	Short: "calculate facets for search",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetFlags(log.Lshortfile)
		var err error

		var q string
		var hasFilter bool
		if cmd.Flags().Changed("query") {
			q, err = cmd.Flags().GetString("query")
			if err != nil {
				hasFilter = false
			}
			hasFilter = true
		}

		if cmd.Flags().Changed("dir") {
			dir, err := cmd.Flags().GetString("dir")
			if err != nil {
				log.Fatal(err)
			}
			m, err := filepath.Glob(filepath.Join(dir, "/*"))
			if err != nil {
				log.Fatal(err)
			}
			dataFiles = m
		}

		if cmd.Flags().Changed("config") {
			switch cfgFile {
			case "":
				log.Fatalf("no config provided")
			default:
				idx, err = facet.NewIndexFromFiles(cfgFile)
				if err != nil {
					log.Fatal(err)
				}
				err = idx.SetData(lo.ToAnySlice(dataFiles)...)
				if err != nil {
					log.Fatalf("error with data files: %v\n", err)
				}
			}
		} else {
			in := cmd.InOrStdin()
			err = idx.Decode(in)
			if err != nil {
				log.Fatal(err)
			}
		}

		if hasFilter {
			idx = idx.Filter(q)
		}

		println(len(idx.Data))
		idx.Print()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "index config file in json format")
	rootCmd.PersistentFlags().StringSliceVarP(&dataFiles, "input", "i", []string{}, "data to index")
	rootCmd.PersistentFlags().StringP("dir", "d", "", "data dir")
	rootCmd.PersistentFlags().StringP("query", "q", "", "encoded query string")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// TODO: check for local file later
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
