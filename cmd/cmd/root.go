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
	Use:   "facet [<cfg>] [<data>] [<filters>]",
	Short: "calculate facets for search",
	Long: `facet aggregates data on specified fields, with option filters. 

The command accepts stdin, flags, and positional arguments.

If a config file has a "data" field no other argument or flag is required. 

Without the "data" field, data must be specified through a flag or positional
argument.

By default, results are printed to stdout as json.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		log.SetFlags(log.Lshortfile)

		var (
			err       error
			filters   string
			hasFilter bool
		)

		if cmd.Flags().Changed("query") {
			filters, err = cmd.Flags().GetString("query")
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

		if len(args) > 0 {
			cfgFile = args[0]

			if len(args) > 1 {
				dataFiles = append(dataFiles, args[1])
			}

			if len(args) > 2 {
				filters = args[2]
			}
		}

		if cfgFile != "" {
			idx, err = facet.NewIndexFromFiles(cfgFile)
			if err != nil {
				log.Fatal(err)
			}

			err = idx.SetData(lo.ToAnySlice(dataFiles)...)
			if err != nil {
				log.Fatalf("error with data files: %v\n", err)
			}

			if hasFilter {
				idx = idx.Filter(filters)
			}

		} else {
			in := cmd.InOrStdin()
			err = idx.Decode(in)
			if err != nil {
				log.Fatal(err)
			}
		}
		if p, err := cmd.Flags().GetBool("pretty"); err == nil && p {
			idx.PrettyPrint()
		} else {
			idx.Print()
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "json formatted config file")
	rootCmd.PersistentFlags().StringSliceVarP(&dataFiles, "files", "f", []string{}, "list of data files to index")
	rootCmd.PersistentFlags().Bool("pretty", false, "pretty print json output")
	rootCmd.PersistentFlags().StringP("input", "i", "", "json formatted input")
	rootCmd.PersistentFlags().StringP("dir", "d", "", "directory of data files")
	rootCmd.PersistentFlags().StringP("query", "q", "", "encoded query/filter string (eg. color=red&color=pink&category=post")

	rootCmd.PersistentFlags().IntP("workers", "w", 1, "number of workers for computing facets")
	viper.BindPFlag("workers", rootCmd.Flags().Lookup("workers"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
