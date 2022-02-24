/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"log"
	"github.com/bigkevmcd/go-configparser"
)

type initStruct struct {
	path string
	gitdir string
	conf *configparser.ConfigParser
}

func (init initStruct) New(path string) {
	init.path = path
	init.gitdir = path + ".git"
	_, err := os.Stat(init.gitdir)
	if err != nil{
		fmt.Printf("Not a Git repository %s",path)
		log.Fatal(err)
	}
	//Read configuration file in .git/config
	cf := repo_file(init, "config")
	_, Is_cf := os.Stat(cf)
	if Is_cf == nil{
		init.conf, _ = configparser.NewConfigParserFromFile(cf)
	}else{
		//コンフィグファイルが見つからない場合エラー停止する
		log.Fatal(Is_cf)
	}
}
func repo_file (init initStruct, path string) string {
	//create dirname(path) if absent.
	dir_result, path := repo_dir(init, path)
	if dir_result{
		return repo_path(init, path)
	}
	//path先がディレクトリではない場合エラーで停止する
	log.Fatal(path)
	return path
}

func repo_dir(init initStruct, path string) (bool,string){
	//mkdir path if absent if mkdir
	path = repo_path(init, path)

	path_stat, err := os.Stat(path)
	
	//path先が存在しない場合はディレクトリを作って返す
	if err == nil{
		if path_stat.IsDir(){
			return true,path
		}else{
			return false,"Not a directory" + path
		}
	}

	_ = os.MkdirAll(path,0777)
	return true,path
}

func repo_path(init initStruct, path string) string {
	return init.gitdir + path
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		//構造体の宣言
		var test initStruct
		test.New(os.Args[1])

		fmt.Println("init called")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
