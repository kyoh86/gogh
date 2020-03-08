package gen

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/dave/jennifer/jen"
	"github.com/kyoh86/gogh/appenv/internal"
	strcase "github.com/stoewer/go-strcase"
)

type Generator struct {
	PackageName        string
	DefaultEnvarPrefix string
	// TODO: make it DefaultServiceName ( with making ServiceName option can be accepted in LoadKeyring / SaveKeyring )
	ServiceName string

	name string

	storeFile    bool
	storeKeyring bool
	storeEnvar   bool
}

const (
	pkgYAML    = "gopkg.in/yaml.v3"
	pkgKeyring = "github.com/zalando/go-keyring"
	pkgProps   = "github.com/kyoh86/gogh/appenv/prop"
)

func (g *Generator) init() error {
	g.name = "env.generator"
	if _, file, _, ok := runtime.Caller(2); ok {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(cwd, file)
		if err != nil {
			return err
		}
		g.name = rel
	}

	g.storeFile = false
	g.storeKeyring = false
	g.storeEnvar = false
	g.DefaultEnvarPrefix = strcase.UpperSnakeCase(g.DefaultEnvarPrefix)

	return nil
}

func (g *Generator) createFile(packagePath string) *jen.File {
	var file *jen.File
	if g.PackageName != "" {
		file = jen.NewFilePathName(packagePath, g.PackageName)
	} else {
		file = jen.NewFilePath(packagePath)
	}
	file.HeaderComment(fmt.Sprintf("// Code generated by %s DO NOT EDIT.", g.name))
	return file
}

func (g *Generator) parseProps(properties []*internal.Property) {
	for _, p := range properties {
		g.storeFile = g.storeFile || p.StoreFile
		g.storeKeyring = g.storeKeyring || p.StoreKeyring
		g.storeEnvar = g.storeEnvar || p.StoreEnvar
	}
}

func (g *Generator) doMerge(file *jen.File, properties []*internal.Property) {
	/* TODO:
	- make the "Merge", LoadFile, SaveFile, LoadKeyring, SaveKeyring and GetEnvar function private
	- create new interface Config (like "Merged").
	- create new "Preference() Accessor" and "Get() Config" function.
	*/
	file.Func().Id("Merge").ParamsFunc(func(mergeParams *jen.Group) {
		if g.storeFile {
			mergeParams.Id("file").Id("File")
		}
		if g.storeKeyring {
			mergeParams.Id("keyring").Id("Keyring")
		}
		if g.storeEnvar {
			mergeParams.Id("envar").Id("Envar")
		}
	}).Params(jen.Id("merged").Id("Merged")).BlockFunc(func(mergeCodes *jen.Group) {
		file.Type().Id("Merged").StructFunc(func(mergedFields *jen.Group) {
			for _, p := range properties {
				mergedFields.Id(p.CamelName).Id(p.ValueType)

				file.Func().Params(jen.Id("m").Id("*Merged")).Id(p.Name).Params().Id(p.ValueType).Block(
					jen.Return(jen.Id("m").Dot(p.CamelName)),
				).Line()

				mergeCodes.Id("merged").Dot(p.CamelName).Op("=").New(jen.Qual(p.PkgPath, p.Name)).Dot("Default").Call().Assert(jen.Id(p.ValueType))
				if p.StoreFile {
					g.tryMerge(mergeCodes, "file", p)
				}
				if p.StoreKeyring {
					g.tryMerge(mergeCodes, "keyring", p)
				}
				if p.StoreEnvar {
					g.tryMerge(mergeCodes, "envar", p)
				}
				mergeCodes.Line()
			}
		})
		mergeCodes.Return()
	})
}

func (g *Generator) tryMerge(mergeCodes *jen.Group, srcName string, p *internal.Property) {
	mergeCodes.If(jen.Id(srcName).Dot(p.Name).Op("!=").Nil()).Block(
		jen.Id("merged").Dot(p.CamelName).Op("=").Id(srcName).Dot(p.Name).Dot("Value").Call().Assert(jen.Id(p.ValueType)),
	)
}

func (g *Generator) doAccess(file *jen.File, properties []*internal.Property) {
	file.Type().Id("Accessor").StructFunc(func(accessorFields *jen.Group) {
		if g.storeFile {
			accessorFields.Id("file").Id("File")
		}
		if g.storeKeyring {
			accessorFields.Id("keyring").Id("Keyring")
		}
	})

	file.Func().Params(jen.Id("a").Id("*Accessor")).Id("Names").Call().Params(jen.Index().String()).Block(
		jen.Return().Index().String().ValuesFunc(func(namesList *jen.Group) {
			file.Func().Params(jen.Id("a").Id("*Accessor")).Id("Property").Params(jen.Id("name").String()).Params(jen.Qual(pkgProps, "Accessor"), jen.Id("error")).Block(
				jen.Switch(jen.Id("name")).BlockFunc(func(propSwitch *jen.Group) {
					for _, p := range properties {
						// Add property name
						namesList.Lit(p.DottedName)

						// Add property case
						propSwitch.Case(jen.Lit(p.DottedName)).
							Block(jen.Return(jen.Id("&"+p.CamelName+"Accessor").Values(jen.Dict{
								jen.Id("parent"): jen.Id("a"),
							}), jen.Nil()))

						// Build Poperty Accessor
						file.Type().Id(p.CamelName + "Accessor").Struct(
							jen.Id("parent").Id("*Accessor"),
						)

						// Implement "Get" Func
						file.Func().Params(jen.Id("a").Id("*"+p.CamelName+"Accessor")).Id("Get").Params().Params(jen.String(), jen.Id("error")).BlockFunc(func(getCodes *jen.Group) {
							if p.StoreFile {
								g.tryGet(getCodes, "file", p)
							}
							if p.StoreKeyring {
								g.tryGet(getCodes, "keyring", p)
							}
							getCodes.Return(jen.Lit(""), jen.Nil())
						}).Line()

						// Implement "Set" Func
						file.Func().Params(jen.Id("a").Id("*" + p.CamelName + "Accessor")).Id("Set").Params(jen.Id("value").String()).Params(jen.Id("error")).BlockFunc(func(setCodes *jen.Group) {
							if p.StoreFile {
								g.trySet(setCodes, "file", p)
							}
							if p.StoreKeyring {
								g.trySet(setCodes, "keyring", p)
							}
							setCodes.Return(jen.Nil())
						}).Line()

						// Implement "Unset" Func
						file.Func().Params(jen.Id("a").Id("*" + p.CamelName + "Accessor")).Id("Unset").Params().BlockFunc(func(unsetCodes *jen.Group) {
							if p.StoreFile {
								g.tryUnset(unsetCodes, "file", p)
							}
							if p.StoreKeyring {
								g.tryUnset(unsetCodes, "keyring", p)
							}
						}).Line()
					}
				}),
				jen.Return(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("invalid propertye name %q"), jen.Id("name"))),
			)
		}),
	).Line()
}

func (g *Generator) tryGet(getCodes *jen.Group, srcName string, p *internal.Property) {
	getCodes.Block(
		jen.Id("p").Op(":=").Id("a").Dot("parent").Dot(srcName).Dot(p.Name),
		jen.If(jen.Id("p").Op("!=").Nil()).BlockFunc(func(ifBlock *jen.Group) {
			ifBlock.List(jen.Id("text"), jen.Err()).Op(":=").Id("p").Dot("MarshalText").Call()
			if p.Mask {
				ifBlock.Return(jen.Id("p").Dot("Mask").Call(jen.String().Call(jen.Id("text"))), jen.Err())
			} else {
				ifBlock.Return(jen.String().Call(jen.Id("text")), jen.Err())
			}
		}),
	)
}

func (g *Generator) trySet(setCodes *jen.Group, srcName string, p *internal.Property) {
	setCodes.Block(
		jen.Id("p").Op(":=").Id("a").Dot("parent").Dot(srcName).Dot(p.Name),
		jen.If(jen.Id("p").Op("==").Nil()).Block(
			jen.Id("p").Op("=").New(jen.Qual(p.PkgPath, p.Name)),
		),
		jen.If(
			jen.Err().Op(":=").Id("p").Dot("UnmarshalText").Call(jen.Id("[]byte").Call(jen.Id("value"))),
			jen.Err().Op("!=").Nil(),
		).Block(
			jen.Return(jen.Err()),
		),
		jen.Id("a").Dot("parent").Dot(srcName).Dot(p.Name).Op("=").Id("p"),
	)
}

func (g *Generator) tryUnset(unsetCodes *jen.Group, srcName string, p *internal.Property) {
	unsetCodes.Id("a").Dot("parent").Dot(srcName).Dot(p.Name).Op("=").Nil()
}

func (g *Generator) doFile(file *jen.File, properties []*internal.Property) {
	file.Type().Id("File").StructFunc(func(fileFields *jen.Group) {
		for _, p := range properties {
			if !p.StoreFile {
				continue
			}
			fileFields.Id(p.Name).
				Op("*").Qual(p.PkgPath, p.Name).
				Tag(map[string]string{"yaml": p.CamelName + ",omitempty"})
		}
	})
	file.Line()

	file.Func().Id("SaveFile").
		Params(
			jen.Id("w").Qual("io", "Writer"),
			jen.Id("file").Id("*File"),
		).
		Add(jen.Id("error")).
		Block(
			jen.Return(
				jen.Qual("gopkg.in/yaml.v3", "NewEncoder").Call(jen.Id("w")).
					Op(".").
					Id("Encode").Call(jen.Id("file")),
			),
		)
	file.Line()
	file.Func().Id("LoadFile").
		Params(
			jen.Id("r").Qual("io", "Reader"),
		).
		Params(
			jen.Id("file").Id("File"),
			jen.Err().Id("error"),
		).
		Block(
			jen.Err().Op("=").Qual("gopkg.in/yaml.v3", "NewDecoder").Call(jen.Id("r")).
				Op(".").
				Id("Decode").Call(jen.Op("&").Id("file")),
			jen.Return(),
		)
	file.Line()
}

func (g *Generator) doKeyring(file *jen.File, packagePath string, properties []*internal.Property) {
	file.Type().Id("Keyring").StructFunc(func(keyringFields *jen.Group) {
		file.Func().Id("LoadKeyring").Params().Params(jen.Id("key").Id("Keyring"), jen.Err().Id("error")).BlockFunc(func(loadKeyringCodes *jen.Group) {
			file.Func().Id("SaveKeyring").Params(jen.Id("key").Id("Keyring")).Params(jen.Err().Id("error")).BlockFunc(func(saveKeyringCodes *jen.Group) {
				for _, p := range properties {
					if !p.StoreKeyring {
						continue
					}
					keyringFields.Id(p.Name).
						Op("*").Qual(p.PkgPath, p.Name)
					loadKeyringCodes.Block(jen.List(jen.Id("v"), jen.Err()).Op(":=").Qual(pkgKeyring, "Get").
						Call(jen.Lit(packagePath), jen.Lit(p.KebabName)),
						jen.If(jen.Err().Op("==").Nil()).Block(
							jen.Var().Id("value").Qual(p.PkgPath, p.Name),
							jen.If(
								jen.Err().Op("=").Id("value").Dot("UnmarshalText").Call(jen.Index().Byte().Parens(jen.Id("v"))),
								jen.Err().Op("!=").Nil(),
							).Block(
								jen.Return(jen.Id("key"), jen.Err()),
							),
							jen.Id("key").Dot(p.Name).Op("=").Id("&value"),
						).Else().Block(
							jen.Qual("log", "Printf").Call(jen.Lit("info: there's no secret in "+p.KebabName+"@"+packagePath+" (%v)"), jen.Err()),
						),
					)
					saveKeyringCodes.Block(
						jen.List(jen.Id("buf"), jen.Err()).Op(":=").Id("key").Dot(p.Name).Dot("MarshalText").Call(),
						jen.If(jen.Err().Op("!=").Nil()).Block(
							jen.Return(jen.Err()),
						),
						jen.If(
							jen.Err().Op(":=").Qual(pkgKeyring, "Set").Call(jen.Lit(packagePath), jen.Lit(p.KebabName), jen.String().Call(jen.Id("buf"))),
							jen.Err().Op("!=").Nil(),
						).Block(
							jen.Return(jen.Err()),
						),
					)
				}
				loadKeyringCodes.Return()
				saveKeyringCodes.Return(jen.Nil())
			})
		}).Line()
	})
}

func (g *Generator) doEnvar(file *jen.File, properties []*internal.Property) {
	file.Type().Id("envarOption").Struct(
		jen.Id("envarPrefix").String(),
	).Line()
	file.Type().Id("GetEnvarOption").Func().Params(jen.Id("*envarOption")).Line()
	file.Func().Id("GetEnvarPrefix").Params(jen.Id("prefix").String()).Id("GetEnvarOption").Block(
		jen.Return(jen.Func().Params(jen.Id("o").Id("*envarOption")).Block(
			jen.Id("o").Dot("envarPrefix").Op("=").Id("prefix"),
		)),
	).Line()
	file.Type().Id("Envar").StructFunc(func(envarFields *jen.Group) {
		file.Func().Id("GetEnvar").Params(jen.Id("opt").Id("...GetEnvarOption")).Params(jen.Id("envar").Id("Envar"), jen.Err().Id("error")).BlockFunc(func(getEnvarCodes *jen.Group) {
			getEnvarCodes.Id("o").Op(":=").Id("envarOption").Values(jen.Dict{
				jen.Id("envarPrefix"): jen.Lit(g.DefaultEnvarPrefix),
			})
			for _, p := range properties {
				if !p.StoreEnvar {
					continue
				}
				envarFields.Id(p.Name).Op("*").Qual(p.PkgPath, p.Name)

				getEnvarCodes.Block(jen.List(jen.Id("v")).Op(":=").Qual("os", "Getenv").
					Call(jen.Id("o").Dot("envarPrefix").Op("+").Lit(p.SnakeName)),
					jen.If(jen.Id("v").Op("==").Lit("")).Block(
						jen.Qual("log", "Printf").Call(jen.Lit("info: there's no envar %s"+p.SnakeName+" (%v)"), jen.Id("o").Dot("envarPrefix"), jen.Err()),
					).Else().Block(
						jen.Var().Id("value").Qual(p.PkgPath, p.Name),
						jen.If(
							jen.Err().Op("=").Id("value").Dot("UnmarshalText").Call(jen.Index().Byte().Parens(jen.Id("v"))),
							jen.Err().Op("!=").Nil(),
						).Block(
							jen.Return(jen.Id("envar"), jen.Err()),
						),
						jen.Id("envar").Dot(p.Name).Op("=").Id("&value"),
					),
				)
			}
			getEnvarCodes.Return()
		}).Line()
	}).Line()
}

func (g *Generator) Do(packagePath, outDir string, properties ...*internal.Property) error {
	if err := g.init(); err != nil {
		return err
	}

	full, err := filepath.Abs(outDir)
	if err != nil {
		return err
	}

	g.parseProps(properties)

	mergeFile := g.createFile(packagePath)
	g.doMerge(mergeFile, properties)
	if err := mergeFile.Save(filepath.Join(full, "merge_gen.go")); err != nil {
		return err
	}

	accessFile := g.createFile(packagePath)
	g.doAccess(accessFile, properties)
	if err := accessFile.Save(filepath.Join(full, "access_gen.go")); err != nil {
		return err
	}

	if g.storeFile {
		fileFile := g.createFile(packagePath)
		fileFile.ImportAlias(pkgYAML, "yaml")
		g.doFile(fileFile, properties)
		if err := fileFile.Save(filepath.Join(full, "file_gen.go")); err != nil {
			return err
		}
	}

	if g.storeKeyring {
		keyringFile := g.createFile(packagePath)
		keyringFile.ImportAlias(pkgKeyring, "keyring")
		g.doKeyring(keyringFile, packagePath, properties)
		if err := keyringFile.Save(filepath.Join(full, "keyring_gen.go")); err != nil {
			return err
		}
	}

	if g.storeEnvar {
		envarFile := g.createFile(packagePath)
		g.doEnvar(envarFile, properties)
		if err := envarFile.Save(filepath.Join(full, "envar_gen.go")); err != nil {
			return err
		}
	}

	return nil
}
