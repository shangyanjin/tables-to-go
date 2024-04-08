package main

import (
	"flag"
	"fmt"
	"os"

	"tables-to-go/internal/cli"
	"tables-to-go/pkg/database"
	"tables-to-go/pkg/output"
	"tables-to-go/pkg/settings"

	"github.com/spf13/viper"
)

// CmdArgs represents the supported command line args
type CmdArgs struct {
	Help bool
	*settings.Settings
}

// NewCmdArgs creates and prepares the command line arguments with default values
func NewCmdArgs() (args *CmdArgs) {

	args = &CmdArgs{
		Settings: settings.New(),
	}
	//alias for the command line arguments
	flag.StringVar(&args.User, "user", args.User, "user to connect to the database")
	flag.StringVar(&args.Pswd, "password", args.Pswd, "password of user")
	flag.StringVar(&args.DbName, "db", args.DbName, "database name")
	flag.StringVar(&args.Host, "host", args.Host, "host of database")
	flag.Var(&args.DbType, "type", fmt.Sprintf("type of database to use, currently supported: %v", settings.SprintfSupportedDbTypes()))
	flag.StringVar(&args.OutputFilePath, "o", args.OutputFilePath, "output file path, default is current working directory")
	flag.StringVar(&args.OutputFilePath, "output", args.OutputFilePath, "output file path, default is current working directory")

	//command line arguments
	flag.BoolVar(&args.Help, "?", false, "shows help and usage")
	flag.BoolVar(&args.Help, "help", false, "shows help and usage")
	flag.BoolVar(&args.Verbose, "v", args.Verbose, "verbose output")
	flag.BoolVar(&args.VVerbose, "vv", args.VVerbose, "more verbose output")
	flag.BoolVar(&args.Force, "f", args.Force, "force; skip tables that encounter errors")

	flag.Var(&args.DbType, "t", fmt.Sprintf("type of database to use, currently supported: %v", settings.SprintfSupportedDbTypes()))
	flag.StringVar(&args.User, "u", args.User, "user to connect to the database")
	flag.StringVar(&args.Pswd, "p", args.Pswd, "password of user")

	flag.StringVar(&args.DbName, "d", args.DbName, "database name")
	flag.StringVar(&args.Schema, "s", args.Schema, "schema name")
	flag.StringVar(&args.Host, "h", args.Host, "host of database")

	flag.StringVar(&args.Port, "port", args.Port, "port of database host, if not specified, it will be the default ports for the supported databases")
	flag.StringVar(&args.SSLMode, "sslmode", args.SSLMode, "Connect to database using secure connection. (default \"disable\")\nThe value will be passed as is to the underlying driver.\nRefer to this site for supported values: https://www.postgresql.org/docs/current/libpq-ssl.html")
	flag.StringVar(&args.Socket, "socket", args.Socket, "The socket file to use for connection. If specified, takes precedence over host:port.")

	flag.StringVar(&args.OutputFilePath, "of", args.OutputFilePath, "output file path, default is current working directory")

	flag.Var(&args.OutputFormat, "format", "format of struct fields (columns): camelCase (c) or original (o)")

	flag.Var(&args.FileNameFormat, "fn-format", "format of the filename: camelCase (c, default) or snake_case (s)")
	flag.StringVar(&args.Prefix, "pre", args.Prefix, "prefix for file- and struct names")
	flag.StringVar(&args.Suffix, "suf", args.Suffix, "suffix for file- and struct names")
	flag.StringVar(&args.PackageName, "pn", args.PackageName, "package name")
	flag.Var(&args.Null, "null", "representation of NULL columns: sql.Null* (sql) or primitive pointers (native|primitive)")

	flag.BoolVar(&args.NoInitialism, "no-initialism", args.NoInitialism, "disable the conversion to upper-case words in column names")

	flag.BoolVar(&args.TagsNoDb, "tags-no-db", args.TagsNoDb, "do not create db-tags")

	flag.BoolVar(&args.TagsMastermindStructable, "tags-structable", args.TagsMastermindStructable, "generate struct with tags for use in Masterminds/structable (https://github.com/Masterminds/structable)")
	flag.BoolVar(&args.TagsMastermindStructableOnly, "tags-structable-only", args.TagsMastermindStructableOnly, "generate struct with tags ONLY for use in Masterminds/structable (https://github.com/Masterminds/structable)")
	flag.BoolVar(&args.IsMastermindStructableRecorder, "structable-recorder", args.IsMastermindStructableRecorder, "generate a structable.Recorder field")

	// disable the print of usage when an error occurs
	flag.CommandLine.Usage = func() {}

	flag.Parse()

	return args
}

// main function to run the transformations
func main() {

	cmdArgs := NewCmdArgs()

	if cmdArgs.Help {
		flag.Usage()
		os.Exit(0)
	}

	// if exist config file, use it to update the args settings
	viper.SetConfigName("config.toml")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err == nil {
		cmdArgs.DbType = settings.DBType(viper.GetString("database.type"))
		cmdArgs.User = viper.GetString("database.user")
		cmdArgs.Pswd = viper.GetString("database.password")
		cmdArgs.DbName = viper.GetString("database.db")
		cmdArgs.Host = viper.GetString("database.host")
		cmdArgs.Port = viper.GetString("database.port")
	}

	// 使用 viper 更新 args 设置

	if err := cmdArgs.Verify(); err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	db := database.New(cmdArgs.Settings)

	if err := db.Connect(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	writer := output.NewFileWriter(cmdArgs.OutputFilePath)

	if err := cli.Run(cmdArgs.Settings, db, writer); err != nil {
		fmt.Printf("run error: %v\n", err)
		os.Exit(1)
	}
}
