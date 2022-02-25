/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/bigkevmcd/go-configparser"
	"github.com/spf13/cobra"
)

type initStruct struct {
	path   string
	gitdir string
	conf   *configparser.ConfigParser
}

func New(path string, init initStruct) *initStruct {
	init.path = path
	init.gitdir = filepath.Join(path, ".git")

	//.git/config パスをリターン　cf = .git/config
	cf := repo_file(init, "config", false)
	_, Is_cf := os.Stat(cf)
	//configがあれば読み取る
	if Is_cf == nil {
		init.conf, _ = configparser.NewConfigParserFromFile(cf)
	} else {
		//コンフィグファイルが見つからない場合エラー停止する
		// log.Fatal(Is_cf)
		fmt.Println("configを読み取れませんでした")
	}

	return &init
}
func repo_file(init initStruct, path string, mkdir bool) string {
	//create dirname(path) if absent.
	dir_bool, _ := repo_dir(init, path, mkdir)
	if dir_bool {
		return repo_path(init, path)
	} else {
		//path先がディレクトリではない場合エラーで停止する
		log.Fatal(path)
		return path
	}

}

func repo_dir(init initStruct, path string, mkdir bool) (bool, string) {
	//mkdir path if absent if mkdir
	path = repo_path(init, path)

	path_stat, err := os.Stat(path)

	//path先が存在しない場合はディレクトリを作って返す
	if err == nil {
		if path_stat.IsDir() {
			return true, path
		} else {
			return false, "Not a directory" + path
		}
	}
	//mkdir=falseならばディレクトリを作らずパスのみリターン
	if mkdir {
		_ = os.Mkdir(path, 0777)
	}
	return true, path

}

func repo_path(init initStruct, path string) string {
	//.gitディレクトリ配下に渡されたpath名でディレクトリを作成
	return filepath.Join(init.gitdir, path)
}

func repo_create(path string) {
	/*引数に構造体を受け、その構造体を利用して他の構造体を初期化する関数になってしまっているため
	一つの構造体を初期化するのに2つも構造体を作らないといけなくなっている。要修正*/
	var repo initStruct
	repoed := New(path, repo)

	repo_info, err := os.Stat(repoed.path)
	//Argsに渡されたPATHが撮っている時
	if err == nil {
		//PATHが見つかったがディレクトリではない時
		if !repo_info.IsDir() {
			fmt.Printf("%s is not a directory", repoed.path)
			log.Fatal(path)
		}
		_, err := filepath.Glob(repo.path + "/*")
		//ディレクトリ内にファイルが見つかった時
		if err != nil {
			fmt.Printf("%s is not Empty!", repoed.path)
			log.Fatal(path)
		}
		//PATHが通っていない時
	} else {
		//PATHにディレクトリ作成
		os.Mkdir(repoed.path, 0777)
		fmt.Println(repoed.path)
		fmt.Println("コマンドライン引数に渡されたパスにディレクトリを作成しました")
	}

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
		repo_create(os.Args[2])
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
