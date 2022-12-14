/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/nullc4t/og/internal/types"
	"github.com/nullc4t/og/pkg/extract"
	"github.com/nullc4t/og/pkg/generator"
	"github.com/nullc4t/og/pkg/names"
	"github.com/nullc4t/og/pkg/templates"
	"github.com/nullc4t/og/pkg/transform"
	"github.com/nullc4t/og/pkg/utils"
	"github.com/nullc4t/og/pkg/writer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"path/filepath"
	"strings"
	"text/template"
)

// protoCmd represents the proto command
var protoCmd = &cobra.Command{
	Use:   "proto -i interfaces.go output_dir/",
	Args:  cobra.ExactArgs(1),
	Short: "generate .proto file",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("proto called")

		outputDir, err := filepath.Abs(args[0])
		if err != nil {
			logger.Fatal(err)
		}

		tmpl := template.Must(template.New("").Funcs(templates.FuncMap).Parse(templates.Proto2))

		ifile := viper.GetString("interfaces_file")

		ifaceFile, err := extract.GoFile(ifile)
		if err != nil {
			logger.Fatal(err)
		}

		services, _ := extract.TypesFromASTFile(ifaceFile)

		protoPkg := "proto"
		var protoFile = types.ProtoFile{
			GoPackage:     protoPkg,
			GoPackagePath: strings.Replace(outputDir, ifaceFile.ModulePath, ifaceFile.Module, 1),
			//GoPackagePath: fmt.Sprintf("%s/%s", ifaceFile.ImportPath(), protoPkg),
			Package: "service",
		}

		for _, iface := range services {
			transform.NameEmptyArgsInInterface(iface)
			protoFile.Services = append(protoFile.Services, transform.Interface2ProtoService(*iface))
		}
		for _, service := range protoFile.Services {
			for _, field := range service.Fields {
				protoFile.Messages = append(protoFile.Messages, field.Request, field.Response)
			}
		}

		ctx := extract.NewContext()
		ifaces, structs, err := extract.ParseFile(ctx, ifile, "", 2)
		if err != nil {
			logger.Fatal(err)
		}

		logger.Println(ctx)

		ifaceSliceUtil := utils.NewSlice[*types.Interface](func(a, b *types.Interface) bool {
			return a.Name == b.Name
		})
		for _, iface := range ifaces {
			if ifaceSliceUtil.Contains(services, iface) {
				continue
			}
			protoFile.Messages = append(protoFile.Messages, types.ProtoMessage{
				Name: iface.Name,
				Fields: []types.ProtoField{
					{
						Name:  names.Camel2Snake(iface.Name),
						OneOf: true,
					},
				},
			})
		}

		protoImports := make(map[string]struct{})

		pmsu := utils.NewSlice[types.ProtoMessage](func(a, b types.ProtoMessage) bool {
			return a.Name == b.Name
		})
		for _, str := range structs {
			if pmsu.Contains(protoFile.Messages, types.ProtoMessage{Name: str.Name}) {
				continue
			}
			msg := transform.Struct2ProtoMessage(ctx, *str)
			protoFile.Messages = append(protoFile.Messages, msg)
			for _, field := range msg.Fields {
				switch field.Type {
				case "google.protobuf.Timestamp":
					protoImports["google/protobuf/timestamp.proto"] = struct{}{}
				case "google.protobuf.Any":
					protoImports["google/protobuf/any.proto"] = struct{}{}
				}
			}
		}

		for imp, _ := range protoImports {
			protoFile.Imports = append(protoFile.Imports, types.ProtoImport{Path: imp})
		}

		unit := generator.NewUnit(
			ifaceFile, tmpl, protoFile, nil, nil,
			filepath.Join(
				filepath.Join(args[0], "proto"),
				fmt.Sprintf("%s.proto", names.Camel2Snake(protoFile.Services[0].Name)),
			), writer.File,
		)
		err = unit.Generate()
		if err != nil {
			logger.Fatal(err)
		}
	},
}

func stripSliceAndPointer(s string) string {
	return strings.Replace(strings.Replace(s, "[]", "", 1), "*", "", 1)
}

func init() {
	rootCmd.AddCommand(protoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// protoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// protoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	protoCmd.Flags().StringP("interfaces_file", "i", "", "go file with interface(s)")
	_ = protoCmd.MarkFlagRequired("interfaces_file")
	_ = viper.BindPFlag("interfaces_file", protoCmd.Flag("interfaces_file"))

	protoCmd.Flags().StringSliceP("exclude_types", "x", nil, "exclude types from parsing")
	_ = viper.BindPFlag("exclude_types", protoCmd.Flag("exclude_types"))

}
