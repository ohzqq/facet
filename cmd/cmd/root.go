package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/ohzqq/facet"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile  string
	dataFile string
	idx      = &facet.Index{}
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "facet",
	Short: "calculate facets for search",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		if cfgFile != "" {
			d, err := os.ReadFile(cfgFile)
			if err != nil {
				log.Fatal(err)
			}
			err = json.Unmarshal(d, idx)
			if err != nil {
				log.Fatal(err)
			}
		}
		if dataFile == "" {
			log.Fatalf("no index data provided")
		}
		d, err := os.ReadFile(dataFile)
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(d, &idx.Data)
		if err != nil {
			log.Fatal(err)
		}
		idx.Facets()

		filters := make(url.Values)
		filters.Add("authors", "Alice Winters")
		//filters.Add("tags", "abo")

		ids := idx.Filter(filters)
		fmt.Printf("%#V\n", len(ids))
		//d, err = json.Marshal(facets)
		//if err != nil {
		//log.Fatal(err)
		//}
		//println(string(d))
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

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "index config")
	rootCmd.PersistentFlags().StringVarP(&idx.Name, "name", "n", "", "index name")
	rootCmd.PersistentFlags().StringVarP(&dataFile, "data", "d", "", "data to index")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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
