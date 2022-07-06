package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/zeebo/errs"

	"project_template/pkg/fileutils"
)

// Error is a default error type for gen cli.
var Error = errs.Class("gen cli error")

// commands.
var (
	// root command
	rootCmd = &cobra.Command{
		Use:   "gen",
		Short: "cli for code generation tool",
		RunE:  cmdRoot,
	}

	// create entity
	createEntityCmd = &cobra.Command{
		Use:         "ent",
		Short:       "creates entity",
		RunE:        cmdCreateEntity,
		Annotations: map[string]string{"type": "run"},
	}
	// create entity db
	createEntityDBCmd = &cobra.Command{
		Use:         "ent-db",
		Short:       "creates entity db",
		RunE:        cmdCreateEntityDB,
		Annotations: map[string]string{"type": "run"},
	}
	// create controller
	createControllerCmd = &cobra.Command{
		Use:         "controller",
		Short:       "creates controller",
		RunE:        cmdCreateController,
		Annotations: map[string]string{"type": "run"},
	}
)

func init() {
	rootCmd.AddCommand(createEntityCmd)
	rootCmd.AddCommand(createEntityDBCmd)
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func cmdRoot(cmd *cobra.Command, args []string) error {
	const (
		ctrl  = "Controller"
		ent   = "Entity"
		entDB = "Entity DB"
	)

	prompt := promptui.Select{
		Label: "Generate",
		Items: []string{ctrl, ent, entDB},
	}

	_, res, err := prompt.Run()

	if err != nil {
		return Error.New("prompt failed")
	}

	switch res {
	case ctrl:
		return cmdCreateController(cmd, args)
	case ent:
		return cmdCreateEntity(cmd, args)
	case entDB:
		return cmdCreateEntityDB(cmd, args)
	}

	return nil
}

func cmdCreateEntity(cmd *cobra.Command, args []string) error {
	entityPC, err := promptEntityName()
	if err != nil {
		return Error.New("can not get entity name")
	}

	entitySnake := toSnakeCase(entityPC)

	path, rootPkg, err := getWorkingDirectoryName()
	if err != nil {
		return Error.New("can not get a current directory")
	}

	if res, err := fileutils.IsFileExist(path, entitySnake); err != nil {
		return Error.New("can not check target dir existence")
	} else if res == true {
		return Error.New("target dir already exists")
	}

	path = path + "/" + entitySnake
	err = os.MkdirAll(path, 0744)
	if err != nil {
		return Error.New(fmt.Sprintf("can not create target dir \"%s\"", entitySnake))
	}

	tData := newTemplateData(rootPkg, entityPC)
	tasks := map[string]string{
		entitySnake + ".go":      primeTemplate,
		"service.go":             serviceTemplate,
		entitySnake + "_test.go": testsTemplate,
	}

	for fName, fTemplate := range tasks {
		buf := &bytes.Buffer{}
		temp := template.Must(template.New("mapper").Parse(fTemplate))

		if err := temp.Execute(buf, tData); err != nil {
			fmt.Println(err)
			return Error.New(fmt.Sprintf("can not create prepare \"%s\" file", fName))
		}

		if err := writeFile(path+"/"+fName, buf); err != nil {
			return err
		}
	}

	buf := &bytes.Buffer{}
	temp := template.Must(template.New("mapper").Parse(dbInterfaceTemplate))
	if err := temp.Execute(buf, tData); err != nil {
		fmt.Println(err)
		return Error.New("can not prepare code snippet")
	}

	fmt.Println("\nCopy the code snippet below to DB interface declaration:")
	fmt.Println(buf)

	return nil
}

func cmdCreateEntityDB(cmd *cobra.Command, args []string) error {
	entityPC, err := promptEntityName()
	if err != nil {
		return Error.New("can not get entity name")
	}

	path, rootPkg, err := getWorkingDirectoryName()
	if err != nil {
		return Error.New("can not get a current directory")
	}

	path = path + "/database"
	entitySnake := toSnakeCase(entityPC)
	tData := newTemplateData(rootPkg, entityPC)
	tasks := map[string]string{
		entitySnake + ".go": dbRepoTemplate,
	}

	for fName, fTemplate := range tasks {
		buf := &bytes.Buffer{}
		temp := template.Must(template.New("mapper").Parse(fTemplate))

		if err := temp.Execute(buf, tData); err != nil {
			fmt.Println(err)
			return Error.New(fmt.Sprintf("can not create prepare \"%s\" file", fName))
		}

		if err := writeFile(path+"/"+fName, buf); err != nil {
			return err
		}
	}

	return nil
}

func cmdCreateController(cmd *cobra.Command, args []string) error {
	entityPC, err := promptEntityName()
	if err != nil {
		return Error.New("can not get entity name")
	}

	path, rootPkg, err := getWorkingDirectoryName()
	if err != nil {
		return Error.New("can not get a current directory")
	}

	path = path + "/console/consoleserver/controllers"
	entitySnake := toSnakeCase(entityPC)
	tData := newTemplateData(rootPkg, entityPC)
	tasks := map[string]string{
		entitySnake + ".go": controllerTemplate,
	}

	for fName, fTemplate := range tasks {
		buf := &bytes.Buffer{}
		temp := template.Must(template.New("mapper").Parse(fTemplate))

		if err := temp.Execute(buf, tData); err != nil {
			fmt.Println(err)
			return Error.New(fmt.Sprintf("can not create prepare \"%s\" file", fName))
		}

		if err := writeFile(path+"/"+fName, buf); err != nil {
			return err
		}
	}

	buf := &bytes.Buffer{}
	temp := template.Must(template.New("mapper").Parse(routesTemplate))
	if err := temp.Execute(buf, tData); err != nil {
		fmt.Println(err)
		return Error.New("can not prepare code snippet")
	}

	fmt.Println("\nCopy the code snippet with sample routes:")
	fmt.Println(buf)

	return nil
}

// promptEntityName implements cli interaction in order to get a name of entity'
func promptEntityName() (string, error) {
	validate := func(input string) error {

		if !pascalCaseRule.MatchString(input) {
			return Error.New("entity name must be in PascalCase")
		}

		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Entity",
		Validate: validate,
	}

	res, err := prompt.Run()

	if err != nil {
		return "", Error.New("prompt failed")
	}

	return res, nil
}

// writeFile writes generated content into the file
func writeFile(location string, buf *bytes.Buffer) error {
	if err := ioutil.WriteFile(location, buf.Bytes(), 0644); err != nil {
		return Error.New(fmt.Sprintf("can not write file \"%s\"", location))
	} else {
		fmt.Printf("new: %s\n", location)
	}

	return nil
}

// getWorkingDirectoryName returns name or working directory (not a full path, but only directory)
func getWorkingDirectoryName() (string, string, error) {
	path, err := os.Getwd()
	if err != nil {
		return "", "", Error.New("can not get wroking directory")
	}

	res := strings.Split(path, "/")
	if len(res) == 0 {
		return "", "", Error.New("can not get wroking directory")
	}

	fmt.Println(path)
	fmt.Println(res[len(res)-1])

	return path, res[len(res)-1], nil
}
