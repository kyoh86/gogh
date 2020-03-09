package gen

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/dave/jennifer/jen"
)

type Generator struct {
	PackageName string

	name string

	storeFile    bool
	storeKeyring bool
	storeEnvar   bool
}

const (
	pkgYAML    = "gopkg.in/yaml.v3"
	pkgKeyring = "github.com/zalando/go-keyring"
	pkgTypes   = "github.com/kyoh86/gogh/appenv/types"
	pkgStrcase = "github.com/stoewer/go-strcase"
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

func (g *Generator) parseProps(properties []*Property) {
	for _, p := range properties {
		g.storeFile = g.storeFile || p.storeFile
		g.storeKeyring = g.storeKeyring || p.storeKeyring
		g.storeEnvar = g.storeEnvar || p.storeEnvar
	}
}

func (g *Generator) doMerge(file *jen.File, properties []*Property) {
	/* TODO:
	- make the "Merge", LoadFile, SaveFile, LoadKeyring, SaveKeyring and GetEnvar function private
	- create new interface Config (like "Merged").
	- create new "Preference() Accessor" and "Get() Config" function.
	- create "Save" method on Accessor
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
				mergedFields.Id(p.camelName).Id(p.valueType)

				file.Func().Params(jen.Id("m").Id("*Merged")).Id(p.name).Params().Id(p.valueType).Block(
					jen.Return(jen.Id("m").Dot(p.camelName)),
				).Line()

				mergeCodes.Id("merged").Dot(p.camelName).Op("=").New(jen.Qual(p.pkgPath, p.name)).Dot("Default").Call().Assert(jen.Id(p.valueType))
				if p.storeFile {
					g.tryMerge(mergeCodes, "file", p)
				}
				if p.storeKeyring {
					g.tryMerge(mergeCodes, "keyring", p)
				}
				if p.storeEnvar {
					g.tryMerge(mergeCodes, "envar", p)
				}
				mergeCodes.Line()
			}
		})
		mergeCodes.Return()
	})
}

func (g *Generator) tryMerge(mergeCodes *jen.Group, srcName string, p *Property) {
	mergeCodes.If(jen.Id(srcName).Dot(p.name).Op("!=").Nil()).Block(
		jen.Id("merged").Dot(p.camelName).Op("=").Id(srcName).Dot(p.name).Dot("Value").Call().Assert(jen.Id(p.valueType)),
	)
}

func (g *Generator) doAccess(file *jen.File, properties []*Property) {
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
			file.Func().Params(jen.Id("a").Id("*Accessor")).Id("Property").Params(jen.Id("name").String()).Params(jen.Qual(pkgTypes, "Accessor"), jen.Id("error")).Block(
				jen.Switch(jen.Id("name")).BlockFunc(func(propSwitch *jen.Group) {
					for _, p := range properties {
						// Add property name
						namesList.Lit(p.dottedName)

						// Add property case
						propSwitch.Case(jen.Lit(p.dottedName)).
							Block(jen.Return(jen.Id("&"+p.camelName+"Accessor").Values(jen.Dict{
								jen.Id("parent"): jen.Id("a"),
							}), jen.Nil()))

						// Build Poperty Accessor
						file.Type().Id(p.camelName + "Accessor").Struct(
							jen.Id("parent").Id("*Accessor"),
						)

						// Implement "Get" Func
						file.Func().Params(jen.Id("a").Id("*"+p.camelName+"Accessor")).Id("Get").Params().Params(jen.String(), jen.Id("error")).BlockFunc(func(getCodes *jen.Group) {
							if p.storeFile {
								g.tryGet(getCodes, "file", p)
							}
							if p.storeKeyring {
								g.tryGet(getCodes, "keyring", p)
							}
							getCodes.Return(jen.Lit(""), jen.Nil())
						}).Line()

						// Implement "Set" Func
						file.Func().Params(jen.Id("a").Id("*" + p.camelName + "Accessor")).Id("Set").Params(jen.Id("value").String()).Params(jen.Id("error")).BlockFunc(func(setCodes *jen.Group) {
							if p.storeFile {
								g.trySet(setCodes, "file", p)
							}
							if p.storeKeyring {
								g.trySet(setCodes, "keyring", p)
							}
							setCodes.Return(jen.Nil())
						}).Line()

						// Implement "Unset" Func
						file.Func().Params(jen.Id("a").Id("*" + p.camelName + "Accessor")).Id("Unset").Params().BlockFunc(func(unsetCodes *jen.Group) {
							if p.storeFile {
								g.tryUnset(unsetCodes, "file", p)
							}
							if p.storeKeyring {
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

func (g *Generator) tryGet(getCodes *jen.Group, srcName string, p *Property) {
	getCodes.Block(
		jen.Id("p").Op(":=").Id("a").Dot("parent").Dot(srcName).Dot(p.name),
		jen.If(jen.Id("p").Op("!=").Nil()).BlockFunc(func(ifBlock *jen.Group) {
			ifBlock.List(jen.Id("text"), jen.Err()).Op(":=").Id("p").Dot("MarshalText").Call()
			if p.mask {
				ifBlock.Return(jen.Id("p").Dot("Mask").Call(jen.String().Call(jen.Id("text"))), jen.Err())
			} else {
				ifBlock.Return(jen.String().Call(jen.Id("text")), jen.Err())
			}
		}),
	)
}

func (g *Generator) trySet(setCodes *jen.Group, srcName string, p *Property) {
	setCodes.Block(
		jen.Id("p").Op(":=").Id("a").Dot("parent").Dot(srcName).Dot(p.name),
		jen.If(jen.Id("p").Op("==").Nil()).Block(
			jen.Id("p").Op("=").New(jen.Qual(p.pkgPath, p.name)),
		),
		jen.If(
			jen.Err().Op(":=").Id("p").Dot("UnmarshalText").Call(jen.Id("[]byte").Call(jen.Id("value"))),
			jen.Err().Op("!=").Nil(),
		).Block(
			jen.Return(jen.Err()),
		),
		jen.Id("a").Dot("parent").Dot(srcName).Dot(p.name).Op("=").Id("p"),
	)
}

func (g *Generator) tryUnset(unsetCodes *jen.Group, srcName string, p *Property) {
	unsetCodes.Id("a").Dot("parent").Dot(srcName).Dot(p.name).Op("=").Nil()
}

func (g *Generator) doFile(file *jen.File, properties []*Property) {
	file.Type().Id("File").StructFunc(func(fileFields *jen.Group) {
		for _, p := range properties {
			if !p.storeFile {
				continue
			}
			fileFields.Id(p.name).
				Op("*").Qual(p.pkgPath, p.name).
				Tag(map[string]string{"yaml": p.camelName + ",omitempty"})
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

func (g *Generator) doKeyring(file *jen.File, properties []*Property) {
	file.Type().Id("Keyring").StructFunc(func(keyringFields *jen.Group) {
		file.Func().Id("LoadKeyring").Params(jen.Id("serviceName").String()).Params(jen.Id("key").Id("Keyring"), jen.Err().Id("error")).BlockFunc(func(loadKeyringCodes *jen.Group) {
			file.Func().Id("SaveKeyring").Params(jen.Id("serviceName").String(), jen.Id("key").Id("Keyring")).Params(jen.Err().Id("error")).BlockFunc(func(saveKeyringCodes *jen.Group) {
				for _, p := range properties {
					if !p.storeKeyring {
						continue
					}
					keyringFields.Id(p.name).
						Op("*").Qual(p.pkgPath, p.name)
					loadKeyringCodes.Block(jen.List(jen.Id("v"), jen.Err()).Op(":=").Qual(pkgKeyring, "Get").
						Call(jen.Id("serviceName"), jen.Lit(p.kebabName)),
						jen.If(jen.Err().Op("==").Nil()).Block(
							jen.Var().Id("value").Qual(p.pkgPath, p.name),
							jen.If(
								jen.Err().Op("=").Id("value").Dot("UnmarshalText").Call(jen.Index().Byte().Parens(jen.Id("v"))),
								jen.Err().Op("!=").Nil(),
							).Block(
								jen.Return(jen.Id("key"), jen.Err()),
							),
							jen.Id("key").Dot(p.name).Op("=").Id("&value"),
						).Else().Block(
							jen.Qual("log", "Printf").Call(jen.Lit("info: there's no secret in "+p.kebabName+"@%s (%v)"), jen.Id("serviceName"), jen.Err()),
						),
					)
					saveKeyringCodes.Block(
						jen.List(jen.Id("buf"), jen.Err()).Op(":=").Id("key").Dot(p.name).Dot("MarshalText").Call(),
						jen.If(jen.Err().Op("!=").Nil()).Block(
							jen.Return(jen.Err()),
						),
						jen.If(
							jen.Err().Op(":=").Qual(pkgKeyring, "Set").Call(jen.Id("serviceName"), jen.Lit(p.kebabName), jen.String().Call(jen.Id("buf"))),
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

func (g *Generator) doEnvar(file *jen.File, properties []*Property) {
	file.Type().Id("Envar").StructFunc(func(envarFields *jen.Group) {
		file.Func().Id("GetEnvar").Params(jen.Id("prefix").String()).Params(jen.Id("envar").Id("Envar"), jen.Err().Id("error")).BlockFunc(func(loadEnvarCodes *jen.Group) {
			loadEnvarCodes.Id("prefix").Op("=").Qual(pkgStrcase, "UpperSnakeCase").Call(jen.Id("prefix"))
			for _, p := range properties {
				if !p.storeEnvar {
					continue
				}
				envarFields.Id(p.name).
					Op("*").Qual(p.pkgPath, p.name)

				loadEnvarCodes.Block(jen.List(jen.Id("v")).Op(":=").Qual("os", "Getenv").
					Call(jen.Id("prefix").Op("+").Lit(p.snakeName)),
					jen.If(jen.Id("v").Op("==").Lit("")).Block(
						jen.Qual("log", "Printf").Call(jen.Lit("info: there's no envar %s"+p.snakeName+" (%v)"), jen.Id("prefix"), jen.Err()),
					).Else().Block(
						jen.Var().Id("value").Qual(p.pkgPath, p.name),
						jen.If(
							jen.Err().Op("=").Id("value").Dot("UnmarshalText").Call(jen.Index().Byte().Parens(jen.Id("v"))),
							jen.Err().Op("!=").Nil(),
						).Block(
							jen.Return(jen.Id("envar"), jen.Err()),
						),
						jen.Id("envar").Dot(p.name).Op("=").Id("&value"),
					),
				)
			}
			loadEnvarCodes.Return()
		}).Line()
	}).Line()
}

func (g *Generator) Do(packagePath, outDir string, properties ...*Property) error {
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
		g.doKeyring(keyringFile, properties)
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
