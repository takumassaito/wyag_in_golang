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
		fmt.Println(init.conf)
	} else {
		//コンフィグファイルが見つからない場合エラー停止する
		// log.Fatal(Is_cf)
		fmt.Println("configを読み取れませんでした")
	}

	return &init
}
func repo_file(init initStruct, path string, mkdir bool) string {
	//mkdirがtrueかつ指定したPATHにディレクトリが存在しない場合は、
	//ディレクトリを作成してnew_pathにディレクトリPATHを返す（ディレクトリ作成時の動作）
	//mkdirがfalseならばnew_pathにはファイルPATHのみ作成して返す（ファイル作成時の動作）
	new_path := repo_dir(init, path, mkdir)

	if len(new_path) != 0 {
		return repo_path(init, path)
	}

	return ""
}

func repo_dir(init initStruct, path string, mkdir bool) string {
	//mkdir path if absent if mkdir
	path = repo_path(init, path)

	path_stat, err := os.Stat(path)

	if err == nil {
		if path_stat.IsDir() {
			return path
		} else {
			//mkdirがtrueかつ指定されたPATHがファイルとして存在する時に停止する
			if mkdir {
				fmt.Println("ファイルPATHを指定してディレクトリを作ろうとしています")
				log.Fatal(path)
			}

		}
	} else {
		//path先が存在しない場合はディレクトリを作って返す
		//ただし、mkdir=falseならばディレクトリを作らずパスのみリターン
		if mkdir {
			_ = os.MkdirAll(path, 0777)
		}
	}
	return path

}

func repo_path(init initStruct, path string) string {
	//.gitディレクトリ配下へのPATHを作成して返す
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
			fmt.Printf("%s  is not a directory\n", repoed.path)
			log.Fatal()
		}
		_, err := filepath.Glob(repo.path + "/*.")
		//ディレクトリ内にファイルが見つかった時 できているか要確認！！！！！！！！！！！！！！！
		if err != nil {
			fmt.Printf("%s is not Empty!", repoed.path)
			log.Fatal()
		}
		//PATHが通っていない時
	} else {
		//PATHにディレクトリ作成
		os.Mkdir(repoed.path, 0777)
		fmt.Println(repoed.path)
		fmt.Println("コマンドライン引数に渡されたパスにディレクトリを作成しました")
	}
	//ディレクトリを作成するためmkdirにtrueを入れる
	mkdir := true
	//repoedは構造体へのポインタのため参照外し
	//必要ディレクトリの作成
	repo_dir(*repoed, "branches", mkdir)
	repo_dir(*repoed, "objects", mkdir)
	repo_dir(*repoed, "refs", mkdir)
	repo_dir(*repoed, "refs/tags", mkdir)
	repo_dir(*repoed, "refs/heads", mkdir)

	//.git/descriptionの作成
	file_des, err := os.OpenFile(repo_file(*repoed, "description", false), os.O_RDWR|os.O_CREATE, 0666)
	file_write_check(err)
	defer file_des.Close()
	fmt.Fprint(file_des, "Unnamed repository; edit this file 'description' to name the repository.\n")

	//.git/HEADの作成
	file_head, err2 := os.OpenFile(repo_file(*repoed, "HEAD", false), os.O_RDWR|os.O_CREATE, 0666)
	file_write_check(err2)
	defer file_head.Close()
	fmt.Fprint(file_head, "ref: refs/heads/master\n")

	config := repo_default_config()
	config.SaveWithDelimiter(repo_file(*repoed, "config", false), "=")

}

func file_write_check(err error) {
	if err != nil {
		fmt.Println("ファイルの書き込みにエラーが発生しました")
		log.Fatal(err)
	}
}

func repo_default_config() *configparser.ConfigParser {
	config_parser := configparser.New()

	config_parser.AddSection("core")
	config_parser.Set("core", "repositoryformatversion", "0")
	config_parser.Set("core", "filemode", "false")
	config_parser.Set("core", "bare", "false")

	return config_parser
}

//ルートディレクトリまで再帰的に.gitディレクトリを探す
func repo_find(path string) string {

	git_dir_path := filepath.Join(path, ".git")

	path_info, err := os.Stat(git_dir_path)

	if err == nil {
		if path_info.IsDir() {
			//.gitディレクトリが見つかったら親ディレクトリを返す
			//initStructの初期化に使われるのが.gitの親ディレクトリのため
			return filepath.Join(git_dir_path, "..")
		}
	}

	parent := filepath.Join(path, "..")

	if parent == path {
		fmt.Println("No git directory")
		return ""
	}
	recursived_path := repo_find(parent)

	return recursived_path

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
		//実行ファイルのPATHを取得
		path := os.Args[2]
		//repo_findで.gitディレクトリが見つかれがそのPATHを渡し、見つからなければ与えられた引数で作成
		exit_git_path := repo_find(path)
		if len(exit_git_path) != 0 {
			fmt.Println(".gitディレクトリが見つかりました")
			repo_create(exit_git_path)
			return
		}
		repo_create(os.Args[2])
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
