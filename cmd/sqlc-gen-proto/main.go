package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/ettle/strcase"
	"github.com/sqlc-dev/sqlc/internal/codegen/sdk"
	"github.com/sqlc-dev/sqlc/internal/plugin"
	"google.golang.org/protobuf/proto"

	"github.com/bufbuild/protocompile/parser"
	"github.com/bufbuild/protocompile/reporter"
	"github.com/jhump/protoreflect/v2/protobuilder"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	// "github.com/ryboe/q"
	// "google.golang.org/genproto/googleapis/type/expr"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/genproto/googleapis/type/decimal"
	"google.golang.org/genproto/googleapis/type/money"

	"github.com/jhump/protoreflect/v2/protoprint"

	"google.golang.org/protobuf/types/descriptorpb"
)

var (
	DEFAULT_OUTDIR           = "./sqlcgen"
	DEFAULT_USER_DEFINED_DIR = "./user_defined"
	DEFAULT_DEFAULT_PACKAGE  = "sqlcgen"
	DEFAULT_ONE_OF_ID        = "identifier"
	SYNTAX_PROTO3            = "proto3"
	DO_NOT_GENERATE          = "DO_NOT_GENERATE"

	METHOD_NAMES = []string{"Create", "Get", "Update", "Delete", "List"}
)

type options struct {
	OutDir         string `json:"out_dir,omitempty"          yaml:"out_dir"`
	UserDefinedDir string `json:"user_defined_dir,omitempty" yaml:"user_defined_dir"`
	OneOfID        string `json:"one_of_id,omitempty"        yaml:"one_of_id"`
	DefaultPackage string `json:"default_package,omitempty"  yaml:"default_package"`
}

func getGenRequest() (*plugin.GenerateRequest, error) {
	var req plugin.GenerateRequest
	reqBlob, err := io.ReadAll(os.Stdin)
	if err != nil {
		return nil, err
	}
	if err := proto.Unmarshal(reqBlob, &req); err != nil {
		return nil, err
	}

	return &req, nil
}

func parseOptions(req *plugin.GenerateRequest) (*options, error) {
	var options *options
	if len(req.PluginOptions) == 0 {
		return options, nil
	}
	if err := json.Unmarshal(req.PluginOptions, &options); err != nil {
		return nil, err
	}
	if options.OutDir == "" {
		options.OutDir = DEFAULT_OUTDIR
	}
	if options.UserDefinedDir == "" {
		options.UserDefinedDir = DEFAULT_USER_DEFINED_DIR
	}
	if options.OneOfID == "" {
		options.OneOfID = DEFAULT_ONE_OF_ID
	}
	if options.DefaultPackage == "" {
		options.DefaultPackage = DEFAULT_DEFAULT_PACKAGE
	}

	DEFAULT_OUTDIR = options.OutDir
	DEFAULT_DEFAULT_PACKAGE = options.DefaultPackage
	DEFAULT_ONE_OF_ID = options.OneOfID
	DEFAULT_USER_DEFINED_DIR = options.UserDefinedDir

	return options, nil
}

func main() {
	fdm := make(map[fpath]*protobuilder.FileBuilder)
	req, err := getGenRequest()
	if err != nil {
		log.Fatal(err)
	}
	opts, err := parseOptions(req)
	if err != nil {
		log.Fatal(err)
	}
	p := &Protos{
		files:   fdm,
		tables:  make([]*table, 0),
		enums:   make([]*enum, 0),
		queries: make([]*query, 0),
		options: opts,
	}

	if err := p.run(req); err != nil {
		log.Fatal(err)
	}
}

// used to inform what kind of string to use in maps
type fpath string
type messagename string
type methodname string

type Protofiles []*protobuilder.FileBuilder

// Sort order map
var sortOrder = map[string]int{
	"enum.proto":             1,
	"message.proto":          2,
	"request_response.proto": 3,
	"service.proto":          4,
}

func (pf Protofiles) Len() int {
	return len(pf)
}

func (pf Protofiles) Swap(i, j int) {
	pf[i], pf[j] = pf[j], pf[i]

}

func (pf Protofiles) Less(i, j int) bool {
	filenameI := filepath.Base(pf[i].Path())
	filenameJ := filepath.Base(pf[j].Path())

	orderI, okI := sortOrder[filenameI]
	orderJ, okJ := sortOrder[filenameJ]
	if okI && okJ {
		return orderI < orderJ
	}

	if okI {
		return true
	}

	if okJ {
		return false
	}

	return filenameI < filenameJ
}

// map[ $outdir/$package/$filename]
// map["./sqlcgen/foo/bar/baz/v1/message.proto"]
type Protos struct {
	queries []*query
	tables  []*table
	enums   []*enum
	files   map[fpath]*protobuilder.FileBuilder
	options *options
}

func handleSkip(s string, skips []string) bool {
	for _, skip := range skips {
		if s == skip {
			return true
		}
	}
	return false
}

func copyAnnotations(t *table) error {
	if t.a.ReqResp != nil {
		t.a.ReqResp.a = &Annotations{
			Generate: t.a.Generate,
			Package:  t.a.Package,
			OutDir:   t.a.OutDir,
		}
		if err := setProps(t.a.ReqResp); err != nil {
			return err
		}
	}
	if t.a.Service != nil {
		t.a.Service.a = &Annotations{
			Generate: t.a.Generate,
			Package:  t.a.Package,
			OutDir:   t.a.OutDir,
		}
		if err := setProps(t.a.Service); err != nil {
			return err
		}
	}
	return nil
}

func toRequestName(method, msgName string) string {
	return fmt.Sprintf("%s%s%s", method, msgName, "Request")

}
func toResponseName(method, msgName string) string {
	return fmt.Sprintf("%s%s%s", method, msgName, "Response")
}

func toHttpRule(method string, t *table) *annotations.HttpRule {
	var httpRule *annotations.HttpRule
	identifier := t.a.PrimaryKey

	// POST, LIST
	p := t.a.Service.Path.Path
	// GET, UPDATE, DELETE
	gp := fmt.Sprintf("%s/{%s}", p, identifier)

	switch method {
	case "Create":
		// Create are always POST with the path
		// Ex: /v1/users -X POST
		httpRule = &annotations.HttpRule{
			Pattern: &annotations.HttpRule_Post{
				Post: p,
			},
			Body: "*",
		}
	case "Get":
		// Get are always GET with path + primarykey,
		// OneOf must have oneof the primary key.
		// Connect cann't use a oneofs value.
		// Ex: /v1/users/{uuid}
		httpRule = &annotations.HttpRule{
			Pattern: &annotations.HttpRule_Get{
				Get: gp,
			},
		}
	case "Update":
		// Update are always PUT with path + primarykey,
		// OneOf must have oneof the primary key.
		// Connect cann't use a oneofs value.
		// Ex: /v1/users/{uuid}
		httpRule = &annotations.HttpRule{
			Pattern: &annotations.HttpRule_Put{
				Put: gp,
			},
			Body: "*",
		}
	case "Delete":
		// Delete are always DELETE with path + primarykey,
		// OneOf must have oneof the primary key.
		// Connect cann't use a oneofs value.
		// Ex: /v1/users/{uuid}
		httpRule = &annotations.HttpRule{
			Pattern: &annotations.HttpRule_Delete{
				Delete: gp,
			},
		}
	case "List":
		// List are always GET with path,
		// OneOf must have oneof the primary key.
		// Connect cann't use a oneofs value.
		// Ex: /v1/users/{uuid}
		httpRule = &annotations.HttpRule{
			Pattern: &annotations.HttpRule_Get{
				Get: p,
			},
		}
	}

	return httpRule
}

func (p Protos) createServices(
	mName string,
	rrMap map[methodname]*protobuilder.MessageBuilder,
	t *table,
) (err error) {
	svcfb := p.getFD(t.a.Service.a)
	sName := toPascal(t.a.Service.Name)
	n := protoreflect.Name(*sName)

	sb := svcfb.GetService(n)
	if sb == nil {
		sb = protobuilder.NewService(n)
		defer func() {
			if iErr := svcfb.TryAddService(sb); err != nil {
				err = iErr
			}
		}()
	}
	for _, method := range METHOD_NAMES {
		req := rrMap[methodname(toRequestName(method, mName))]
		resp := rrMap[methodname(toResponseName(method, mName))]

		reqRPC := protobuilder.RpcTypeMessage(req, false)
		respRPC := protobuilder.RpcTypeMessage(resp, false)
		methName := fmt.Sprintf("%s%s", method, mName)

		mb := protobuilder.NewMethod(
			protoreflect.Name(methName),
			reqRPC,
			respRPC,
		)

		httpRule := toHttpRule(method, t)

		methodOptions := &descriptorpb.MethodOptions{}
		proto.SetExtension(methodOptions, annotations.E_Http, httpRule)
		mb.SetOptions(methodOptions)

		// mtdOpts := (*descriptorpb.MethodOptions{}(nil).ProtoReflect().Descriptor())
		// mb.Options.GetFeatures
		// mb.SetOptions()

		if err := sb.TryAddMethod(mb); err != nil {
			return err
		}
	}

	return nil
}

func (p Protos) createRequestResponses(
	messageb *protobuilder.MessageBuilder,
	t *table,
) (map[methodname]*protobuilder.MessageBuilder, error) {
	mName := string(messageb.Name())
	// Copy over Annotations Into New Pointer
	// and Set Properties for ReqResp type.
	// if err := copyAnnotations(t); err != nil {
	// 	return nil, err
	// }

	// Get a new FileDescriptor
	rrfb := p.getFD(t.a.ReqResp.a)

	reqrespMap := make(map[methodname]*protobuilder.MessageBuilder)
	for _, method := range METHOD_NAMES {
		reqName := toRequestName(method, mName)
		respName := toResponseName(method, mName)

		reqb := protobuilder.NewMessage(protoreflect.Name(reqName))
		respb := protobuilder.NewMessage(protoreflect.Name(respName))

		// Add Annotedated Additional Fields
		for aType, aField := range t.a.ReqResp.ReqFields {
			at, err := p.convertType(aType)
			if err != nil {
				return nil, err
			}
			ab := protobuilder.NewField(
				protoreflect.Name(aField),
				at,
			)
			if err := reqb.TryAddField(ab); err != nil {
				return nil, err
			}

		}

		// Add OneOf
		var oneof *protobuilder.OneofBuilder
		if len(*t.a.ReqResp.OneOf) > 0 {
			oneof = protobuilder.NewOneof(protoreflect.Name(p.options.OneOfID))
		}
		for _, ooField := range *t.a.ReqResp.OneOf {
			ooName := protoreflect.Name(ooField)
			oob := messageb.GetField(ooName)
			if oob != nil {
				fCopy := protobuilder.NewField(oob.Name(), oob.Type())
				if err := oneof.TryAddChoice(fCopy); err != nil {
					return nil, err
				}
			}
		}

		if method == "Get" || method == "Update" || method == "Delete" {
			if err := reqb.TryAddOneOf(oneof); err != nil {
				return nil, err
			}
		}

		// only for Create, Update,
		if method == "Create" || method == "Update" {
			fb := protobuilder.NewField(
				protoreflect.Name(*toLowerSnake(mName)),
				protobuilder.FieldTypeMessage(messageb),
			)
			if err := reqb.TryAddField(fb); err != nil {
				return nil, err
			}
		}

		if method == "List" {
			psb := protobuilder.NewField(
				protoreflect.Name("page_size"),
				protobuilder.FieldTypeInt32(),
			)
			if err := reqb.TryAddField(psb); err != nil {
				return nil, err
			}
			ptb := protobuilder.NewField(
				protoreflect.Name("page_token"),
				protobuilder.FieldTypeString(),
			)
			if err := reqb.TryAddField(ptb); err != nil {
				return nil, err
			}
		}

		if err := rrfb.TryAddMessage(reqb); err != nil {
			return nil, err
		}

		skip := false
		for m, enabled := range t.a.ReqResp.RespEmpty {
			if strings.ToLower(method) == strings.ToLower(m) && enabled {
				skip = true
				continue
			}
		}
		if skip || method != "Delete" {
			fb := protobuilder.NewField(
				protoreflect.Name(*toLowerSnake(mName)),
				protobuilder.FieldTypeMessage(messageb),
			)
			if method == "List" {
				fb.SetRepeated()
				npt := protobuilder.NewField(
					protoreflect.Name("next_page_token"),
					protobuilder.FieldTypeString(),
				)
				if err := respb.TryAddField(npt); err != nil {
					return nil, err
				}
			}
			if err := respb.TryAddField(fb); err != nil {
				return nil, err
			}
		}

		if err := rrfb.TryAddMessage(respb); err != nil {
			return nil, err
		}

		reqrespMap[methodname(reqName)] = reqb
		reqrespMap[methodname(respName)] = respb
	}

	return reqrespMap, nil
}

func (p Protos) tableToMessage(
	fileb *protobuilder.FileBuilder,
	t *table,
) error {
	mName := *toPascal(t.i.Rel.Name)
	messageb := protobuilder.NewMessage(protoreflect.Name(mName))

	for _, c := range t.i.Columns {
		if handleSkip(c.Name, t.a.Skips) {
			continue
		}
		if c.PrimaryKey {
			t.a.PrimaryKey = c.Name
		}
		cName := protoreflect.Name(c.Name)
		t, err := p.convertType(c)
		if err != nil {
			return err
		}

		fieldb := protobuilder.NewField(cName, t)
		if c.IsArray {
			fieldb.SetRepeated()
		}

		if err := messageb.TryAddField(fieldb); err != nil {
			return err
		}
	}

	if err := fileb.TryAddMessage(messageb); err != nil {
		return err
	}

	// Copy over Annotations Into New Pointer
	// and Set Properties for ReqResp type.
	if err := copyAnnotations(t); err != nil {
		return err
	}
	// Handle request_response
	var rrMap map[methodname]*protobuilder.MessageBuilder
	if t.a.ReqResp != nil {
		rr, err := p.createRequestResponses(messageb, t)
		if err != nil {
			return err
		}
		rrMap = rr
	}

	// Handle service
	if t.a.Service != nil {
		if err := p.createServices(mName, rrMap, t); err != nil {
			return err
		}
	}

	return nil
}

func (p Protos) queryToMessage(
	messageb *protobuilder.MessageBuilder,
	q *query,
) error {
	for _, c := range q.i.Columns {
		if handleSkip(c.Name, q.a.Skips) {
			continue
		}
		// Append missing Columns from queries to Target
		cName := protoreflect.Name(c.Name)
		if messageb.GetField(cName) == nil {
			t, err := p.convertType(c)
			if err != nil {
				return err
			}
			fieldb := protobuilder.NewField(cName, t)
			if c.IsArray {
				fieldb.SetRepeated()
			}
			if err := messageb.TryAddField(fieldb); err != nil {
				return err
			}
		}
	}
	return nil
}

// lgetFD will reutrn a new empty FD if one does not exist
func (p Protos) getFD(a *Annotations) *protobuilder.FileBuilder {
	op := fpath(a.FullPath)
	if p.files[op] == nil {
		file := protobuilder.NewFile("")
		file.SetSyntax(protoreflect.Proto3)
		file.SetPath(a.FullPath)
		file.SetPackageName(protoreflect.FullName(a.Package))
		p.files[op] = file

	}
	return p.files[op]
}

// Responsible for Constructing message.proto
func (p Protos) Messages() error {
	for _, t := range p.tables {
		fb := p.getFD(t.a)

		if err := p.tableToMessage(fb, t); err != nil {
			return err
		}
	}

	return nil
}

func (p Protos) Queries() error {
	for _, q := range p.queries {
		mb, err := p.GetMessage(protoreflect.Name(*toPascal(q.a.Target)))
		if err != nil {
			return err
		}
		if err := p.queryToMessage(mb, q); err != nil {
			return err
		}
	}

	return nil
}

func (p Protos) appendFieldsToMessage(
	desc *descriptorpb.DescriptorProto,
	b *protobuilder.MessageBuilder,
) error {
	for _, uf := range desc.GetField() {
		var fType *protobuilder.FieldType
		if b.GetField(protoreflect.Name(*uf.Name)) != nil {
			continue
		}
		if uf.Type != nil {
			ft, err := p.convertType(uf.Type.String())
			if err != nil {
				return err
			}
			fType = ft
		} else {
			if uf.TypeName != nil {
				ft, err := p.convertType(userDefinedString(uf.TypeName))
				if err != nil {
					return err
				}
				fType = ft
			}
		}
		b.AddField(protobuilder.NewField(
			protoreflect.Name(*uf.Name),
			fType,
		))
	}

	return nil
}

func (p Protos) appendMethodsToService(
	desc *descriptorpb.ServiceDescriptorProto,
	b *protobuilder.ServiceBuilder,
) error {
	for _, mm := range desc.GetMethod() {
		if b.GetMethod(protoreflect.Name(*mm.Name)) != nil {
			continue
		}
		itb, err := p.GetMessage(protoreflect.Name(*mm.InputType))
		if err != nil {
			return err
		}
		otb, err := p.GetMessage(protoreflect.Name(*mm.OutputType))
		if err != nil {
			return err
		}
		reqb := protobuilder.RpcTypeMessage(itb, false)
		respb := protobuilder.RpcTypeMessage(otb, false)

		nmb := protobuilder.NewMethod(protoreflect.Name(*mm.Name), reqb, respb)
		nmb.SetOptions(mm.Options)

		if err := b.TryAddMethod(nmb); err != nil {
			return err
		}
	}

	return nil
}

func (p Protos) appendValuesToEnum(
	desc *descriptorpb.EnumDescriptorProto,
	b *protobuilder.EnumBuilder,
) error {
	for _, em := range desc.GetValue() {
		if b.GetValue(protoreflect.Name(*em.Name)) != nil {
			continue
		}
		evb := protobuilder.NewEnumValue(protoreflect.Name(*em.Name))
		if err := b.TryAddValue(evb); err != nil {
			return err
		}
	}

	return nil
}

// Always sorts enum.proto, message.proto, request_response.proto, service.proto
func (p Protos) GetFiles() []*protobuilder.FileBuilder {
	var filesSlice Protofiles

	for _, f := range p.files {
		filesSlice = append(filesSlice, f)
	}

	sort.Sort(filesSlice)

	return filesSlice
}

func (p Protos) UserDefined() error {
	for _, file := range p.GetFiles() {
		mp := fmt.Sprintf("%s/%s", p.options.UserDefinedDir, file.Path())

		if _, err := os.Stat(mp); err != nil {
			continue
		}
		mFile, err := os.Open(mp)
		if err != nil {
			return err
		}
		filename := filepath.Base(mp)
		node, err := parser.Parse(filename, mFile, reporter.NewHandler(nil))
		if err != nil {
			return err
		}
		result, err := parser.ResultFromAST(node, false, reporter.NewHandler(nil))
		if err != nil {
			return err
		}
		udFile := result.FileDescriptorProto()

		// Enums
		for _, ue := range udFile.GetEnumType() {
			e := file.GetEnum(protoreflect.Name(*ue.Name))
			if e == nil {
				neb := protobuilder.NewEnum(protoreflect.Name(*ue.Name))
				if err := p.appendValuesToEnum(ue, neb); err != nil {
					return err
				}
				if err := file.TryAddEnum(neb); err != nil {
					return err
				}
				continue
			}
			if err := p.appendValuesToEnum(ue, e); err != nil {
				return err
			}
		}

		// Messaages
		for _, um := range udFile.GetMessageType() {
			m := file.GetMessage(protoreflect.Name(*um.Name))
			if m == nil {
				nmb := protobuilder.NewMessage(protoreflect.Name(*um.Name))
				if err := p.appendFieldsToMessage(um, nmb); err != nil {
					return err
				}
				if err := file.TryAddMessage(nmb); err != nil {
					return err
				}
				continue
			}

			if err := p.appendFieldsToMessage(um, m); err != nil {
				return err
			}
		}

		// Services
		for _, sm := range udFile.GetService() {
			s := file.GetService(protoreflect.Name(*sm.Name))
			if s == nil {
				nsb := protobuilder.NewService(protoreflect.Name(*sm.Name))
				if err := p.appendMethodsToService(sm, nsb); err != nil {
					return err
				}
				if err := file.TryAddService(nsb); err != nil {
					return err
				}
				// make new service
				continue
			}
			if err := p.appendMethodsToService(sm, s); err != nil {
				return err
			}

		}
	}

	return nil
}

func (p Protos) GetMessage(name protoreflect.Name) (*protobuilder.MessageBuilder, error) {
	for _, file := range p.files {
		m := file.GetMessage(name)
		if m == nil {
			continue
		}
		return m, nil
	}

	return nil, fmt.Errorf("%q: Message Not Found.", name)
}

// Responsible for Creating enum.proto files
func (p Protos) Enums() error {
	for _, enum := range p.enums {
		pn := *toPascal(enum.i.Name)
		fb := p.getFD(enum.a)
		eName := protoreflect.Name(pn)

		eb := protobuilder.NewEnum(eName)

		vals := convertEnumValues(pn, enum.i.Vals)
		for _, val := range vals {
			vName := protoreflect.Name(val)
			if err := eb.TryAddValue(protobuilder.NewEnumValue(vName)); err != nil {
				return err
			}
		}
		if err := fb.TryAddEnum(eb); err != nil {
			return err
		}
	}
	return nil
}

func (p Protos) WriteFiles() error {
	fdSlice := []protoreflect.FileDescriptor{}
	printer := protoprint.Printer{}
	for _, file := range p.files {
		b, err := file.Build()
		if err != nil {
			return err
		}

		fdSlice = append(fdSlice, b)
	}

	if err := os.MkdirAll(p.options.OutDir, os.ModePerm); err != nil {
		return err
	}
	if err := printer.PrintProtosToFileSystem(fdSlice, p.options.OutDir); err != nil {
		return err
	}

	return nil
}

type Annotations struct {
	Generate bool   // All:   Generate Protos.
	Package  string // All: Package Name.
	// Replace  map[string]PType    // Tables -> Messages: Type Replacement
	Skips    []string // Tables -> Messages: Skip Field
	ReqResp  *ReqResp // Tables -> Messaes:  Information for generating Request and Responses
	Service  *Service // Tables -> Messaes:  Information for generating Services
	Target   string   // Applies only to querys
	FileName string   // Override output filename
	OutDir   string   // Override base output directory

	FullPath     string // Generated from Package + FilenName
	FullTypeName string // Generated from Package + FilenName
	OutputPath   string // Generated from OutDir + FullPath
	PrimaryKey   string
}

// enumWrapper attaches parsed comment Annotations for plugin.Enum
type enum struct {
	i *plugin.Enum
	a *Annotations
}

// queryWrapper attaches parsed comment Annotations to plugin.Query
type query struct {
	i *plugin.Query
	a *Annotations
}

// tableWrapper attaches parsed comment Annotations to plugin.Table
type table struct {
	i *plugin.Table
	a *Annotations
}

func wrapTable(i *plugin.Table) (*table, error) {
	a, err := parseAnnotations(i.RawComments)
	if err != nil {
		return nil, err
	}
	x := &table{
		i: i,
		a: a,
	}

	if err := setProps(x); err != nil {
		return nil, err
	}
	return x, nil
}

func wrapQuery(i *plugin.Query) (*query, error) {
	a, err := parseAnnotations(i.RawComments)
	if err != nil {
		return nil, err
	}
	x := &query{
		i: i,
		a: a,
	}
	if err := setProps(x); err != nil {
		return nil, err
	}
	return x, nil
}

func wrapEnum(i *plugin.Enum) (*enum, error) {
	a, err := parseAnnotations(i.RawComments)
	if err != nil {
		return nil, err
	}
	x := &enum{
		i: i,
		a: a,
	}
	if err := setProps(x); err != nil {
		return nil, err
	}
	return x, nil
}

func (p *Protos) run(req *plugin.GenerateRequest) error {

	schemas := req.GetCatalog().GetSchemas()
	queries := req.GetQueries()

	for _, schema := range schemas {
		if schema.Name == "pg_catalog" || schema.Name == "information_schema" {
			continue
		}
		// enums = append(enums, schema.GetEnums()...)
		for _, table := range schema.GetTables() {
			t, err := wrapTable(table)
			if err != nil {
				if err.Error() == DO_NOT_GENERATE {
					continue
				}
				return err
			}
			p.tables = append(p.tables, t)
		}
		for _, enum := range schema.GetEnums() {
			e, err := wrapEnum(enum)
			if err != nil {
				if err.Error() == DO_NOT_GENERATE {
					continue
				}
				return err
			}
			p.enums = append(p.enums, e)
		}
	}

	for _, query := range queries {
		q, err := wrapQuery(query)
		if err != nil {
			if err.Error() == DO_NOT_GENERATE {
				continue
			}
			return err
		}
		p.queries = append(p.queries, q)
	}

	if err := p.Enums(); err != nil {
		return err
	}
	if err := p.Messages(); err != nil {
		return err
	}
	if err := p.Queries(); err != nil {
		return err
	}

	if err := p.UserDefined(); err != nil {
		return err
	}

	if err := p.WriteFiles(); err != nil {
		return err
	}

	return nil
}

func parseAnnotations(comments []string) (*Annotations, error) {
	// replace := make(map[string]PType)
	a := &Annotations{
		// Replace: replace,
	}

	for _, line := range comments {
		var prefix string
		if strings.HasPrefix(line, "--") {
			prefix = "--"
		}
		if strings.HasPrefix(line, "/*") {
			prefix = "/*"
		}
		if strings.HasPrefix(line, "#") {
			prefix = "#"
		}
		if prefix == "" {
			continue
		}
		rest := line[len(prefix):]
		if !strings.Contains(rest, ":") {
			continue
		}
		for _, flagOpt := range []string{
			"generate",
			"service",
		} {
			if !strings.HasPrefix(strings.TrimSpace(rest), flagOpt) {
				continue
			}
			opt := fmt.Sprintf(" %s:", flagOpt)

			if !strings.HasPrefix(rest, opt) {
				return nil, fmt.Errorf("invalid metadata: %s", line)
			}
			switch flagOpt {
			case "generate":
				a.Generate = true
			case "service":
			}
		}

		for _, cmdOption := range []string{
			"package",
			"replace",
			"filename",
			"target",
			"skip",
			"request_response",
			"service",
		} {
			if !strings.HasPrefix(strings.TrimSpace(rest), cmdOption) {
				continue
			}
			opt := fmt.Sprintf(" %s: ", cmdOption)

			if !strings.HasPrefix(rest, opt) {
				return nil, fmt.Errorf("invalid metadata: %s", line)
			}

			part := strings.Split(strings.TrimSpace(line), " ")

			switch cmdOption {
			case "package":
				if len(part) != 3 {
					return nil, fmt.Errorf("-- package: <package>... takes exactly 1 argument")
				}
				packageName := part[2]
				a.Package = packageName
			case "target":
				if len(part) != 3 {
					return nil, fmt.Errorf(
						"-- target: <target>... takes exactly 1 argument",
					)
				}
				a.Target = part[2]
			case "skip":
				if len(part) != 3 {
					return nil, fmt.Errorf(
						"-- skip: <skip>... takes exactly 1 argument",
					)
				}
				skipField := part[2]
				a.Skips = append(a.Skips, skipField)
			case "request_response":
				if len(part) < 3 {
					return nil, fmt.Errorf(
						"-- request_response: takes at minimum 2 argument",
					)
				}
				if a.ReqResp == nil {
					es := []string{}
					a.ReqResp = &ReqResp{
						OneOf:     &es,
						ReqFields: make(map[string]string),
						RespEmpty: make(map[string]bool),
					}
				}
				switch part[2] {
				case "oneof":
					if len(part) >= 5 {
						*a.ReqResp.OneOf = append(*a.ReqResp.OneOf, part[3:]...)
					}
					if len(part) == 3 {
						emptySlice := []string{}
						a.ReqResp.OneOf = &emptySlice
					}
				case "req_field":
					if len(part) != 5 {
						return nil, fmt.Errorf(
							"-- request_response: req_field <type> <name> takes exactly 2 arguments.",
						)
					}
					a.ReqResp.ReqFields = make(map[string]string)
					a.ReqResp.ReqFields[part[3]] = part[4]
					// a.ReqResp = rr
				case "resp_empty":
					if len(part) != 4 {
						return nil, fmt.Errorf(
							"-- request_response: resp_empty <method> takes exactly 1 arguments.",
						)
					}
					a.ReqResp.RespEmpty = make(map[string]bool)
					a.ReqResp.RespEmpty[part[3]] = true
				}
			case "service":
				if len(part) != 4 {
					return nil, fmt.Errorf(
						"-- service: <name> <path> ... takes exactly 2 argument",
					)
				}
				name := part[2]
				path := part[3]
				p, err := url.Parse(path)
				if err != nil {
					return nil, err
				}
				a.Service = &Service{
					Path: p,
					Name: name,
				}
			}
		}
	}

	return a, nil
}

func toImportPath(pkg string, filename string) string {
	return fmt.Sprintf("%s/%s", pkgToPath(pkg), filename)
}

func pkgToPath(s string) string {
	// validate against protected type names in importsl
	return strings.ReplaceAll(s, ".", "/")
}

func validatePackageName(s string) error {
	tokens := strings.Split(s, ".")
	for _, reserved := range []string{"enum", "message", "import", "syntax", "repeated"} {
		if len(tokens) > 0 {
			if tokens[0] == reserved {
				return fmt.Errorf(
					"%q: is a reserved word and cannot be used as first part in package: %q",
					tokens[0],
					s,
				)
			}
		}
	}

	return nil
}

func toPascal(s string) *string {
	r := strings.ToLower(s)
	r = strcase.ToPascal(s)

	return &r
}

func toLowerSnake(s string) *string {
	s = strings.ToLower(s)
	s = strcase.ToSnake(s)
	return &s
}

func setCommonProps(a *Annotations) error {
	if a.Package == "" {
		a.Package = DEFAULT_DEFAULT_PACKAGE
	}

	if err := validatePackageName(a.Package); err != nil {
		return err
	}

	if a.OutDir == "" {
		a.OutDir = DEFAULT_OUTDIR
	}
	// Example Package: foo.bar.baz.v1
	// Example Filename: message.proto
	// Path relative from root: foo/bar/baz/v1/message.proto
	a.FullPath = fmt.Sprintf("%s/%s", pkgToPath(a.Package), a.FileName)

	// Example Package: foo.bar.baz.v1
	// Example Filename: message.proto
	// Path relative from root: foo/bar/baz/v1/message.proto
	// a.FullTypeName = fmt.Sprintf("%s.%s", a.Package, a.FileName)

	// Full Filesystem path to be written to
	a.OutputPath = fmt.Sprintf("%s/%s", a.OutDir, a.FullPath)

	return nil
}

func setProps(input interface{}) error {
	switch i := input.(type) {
	case *query:
		if !i.a.Generate {
			return fmt.Errorf(DO_NOT_GENERATE)
		}
		if i.a.Target == "" {
			return fmt.Errorf(
				"To append columns from Queries to protobufs  you must declare a target: <message-to-append-to>",
			)

		}
		if err := setCommonProps(i.a); err != nil {
			return err
		}
	case *enum:
		if !i.a.Generate {
			return fmt.Errorf(DO_NOT_GENERATE)
		}
		if i.a.FileName == "" {
			i.a.FileName = "enum.proto"
		}
		if err := setCommonProps(i.a); err != nil {
			return err
		}
	case *table:
		if !i.a.Generate {
			return fmt.Errorf(DO_NOT_GENERATE)
		}
		if i.a.FileName == "" {
			i.a.FileName = "message.proto"
		}
		if err := setCommonProps(i.a); err != nil {
			return err
		}
	case *ReqResp:
		if !i.a.Generate {
			return fmt.Errorf(DO_NOT_GENERATE)
		}
		if i.a.FileName == "" {
			i.a.FileName = "request_response.proto"
		}
		if err := setCommonProps(i.a); err != nil {
			return err
		}
	case *Service:
		if !i.a.Generate {
			return fmt.Errorf(DO_NOT_GENERATE)
		}
		if i.a.FileName == "" {
			i.a.FileName = "service.proto"
		}
		if err := setCommonProps(i.a); err != nil {
			return err
		}
	}

	return nil
}

type userDefinedString *string

func (p *Protos) convertType(input interface{}) (*protobuilder.FieldType, error) {
	var ct string
	var notNull bool

	switch i := input.(type) {
	case string:
		ct = i
		notNull = true
	case userDefinedString:
		ct = *i
	case *plugin.Column:
		s := sdk.DataType(i.Type)
		notNull = i.NotNull || i.IsArray
		ct = strings.ToLower(s)
	}

	tAny := (*anypb.Any)(nil).ProtoReflect().Descriptor()
	tI32 := (*wrapperspb.Int32Value)(nil).ProtoReflect().Descriptor()
	tI64 := (*wrapperspb.Int64Value)(nil).ProtoReflect().Descriptor()
	tFloat := (*wrapperspb.FloatValue)(nil).ProtoReflect().Descriptor()
	tBytes := (*wrapperspb.BytesValue)(nil).ProtoReflect().Descriptor()
	tBool := (*wrapperspb.BoolValue)(nil).ProtoReflect().Descriptor()
	tString := (*wrapperspb.StringValue)(nil).ProtoReflect().Descriptor()
	tDecimal := (*decimal.Decimal)(nil).ProtoReflect().Descriptor()
	tMoney := (*money.Money)(nil).ProtoReflect().Descriptor()
	tStruct := (*structpb.Struct)(nil).ProtoReflect().Descriptor()
	tTimestamp := (*timestamppb.Timestamp)(nil).ProtoReflect().Descriptor()

	switch ct {
	// Int32
	case "integer",
		"int",
		"int4",
		"pg_catalog.int4",
		"serial",
		"serial4",
		"pg_catalog.serial4",
		"smallserial",
		"smallint", "int2", "pg_catalog.int2", "serial2",
		"pg_catalog.serial2",
		WellKnownInt32Value:
		if notNull {
			return protobuilder.FieldTypeInt32(), nil
		}
		return protobuilder.FieldTypeImportedMessage(
			tI32,
		), nil

	// Int64
	case "interval",
		"pg_catalog.interval",
		"bigint",
		"int8",
		"pg_catalog.int8",
		"bigserial",
		"serial8",
		"pg_catalog.serial8",
		"TYPE_INT64",
		WellKnownInt64Value:
		if notNull {
			return protobuilder.FieldTypeInt64(), nil
		}
		return protobuilder.FieldTypeImportedMessage(
			tI64,
		), nil

	// Float
	case "real",
		"float4",
		"pg_catalog.float4",
		"float",
		"double precision",
		"float8",
		"pg_catalog.float8",
		"TYPE_DOUBLE",
		"TYPE_FLOAT",
		WellKnownFloatValue:
		if notNull {
			return protobuilder.FieldTypeFloat(), nil
		}
		return protobuilder.FieldTypeImportedMessage(
			tFloat,
		), nil

	case "numeric", "pg_catalog.numeric", WellKnownDecimal:
		return protobuilder.FieldTypeImportedMessage(
			tDecimal,
		), nil

	case "money", WellKnownMoney:
		return protobuilder.FieldTypeImportedMessage(
			tMoney,
		), nil

	case "boolean", "bool", "pg_catalog.bool", WellKnownBoolValue:
		if notNull {
			return protobuilder.FieldTypeBool(), nil
		}
		return protobuilder.FieldTypeImportedMessage(
			tBool,
		), nil

	case "json", WellKnownStruct:
		return protobuilder.FieldTypeImportedMessage(
			tStruct,
		), nil

	case "uuid", "jsonb", "bytea", "blob", "pg_catalog.bytea", WellKnownBytesValue:
		if notNull {
			return protobuilder.FieldTypeBytes(), nil
		}
		return protobuilder.FieldTypeImportedMessage(
			tBytes,
		), nil

	case "pg_catalog.timestamptz",
		"date",
		"timestamptz",
		"pg_catalog.timestamp",
		"pg_catalog.timetz",
		"pg_catalog.time",
		WellKnownTimestamp:
		return protobuilder.FieldTypeImportedMessage(
			tTimestamp,
		), nil

	case "citext",
		"lquery",
		"ltree",
		"ltxtquery",
		"name",
		"inet",
		"cidr",
		"macaddr",
		"macaddr8",
		"pg_catalog.bpchar",
		"pg_catalog.varchar",
		"string",
		"text",
		WellKnownStringValue:
		if notNull {
			return protobuilder.FieldTypeString(), nil
		}
		return protobuilder.FieldTypeImportedMessage(
			tString,
		), nil

	// All these PG Range Types Required FieldOptions
	// Handle this Later
	case "daterange":
		// switch driver {
		// case opts.SQLDriverPGXV4:
		// 	return "pgtype.Daterange"
		// case opts.SQLDriverPGXV5:
		// 	return "pgtype.Range[pgtype.Date]"
		// default:
		// 	return "interface{}"
		// }

	case "datemultirange":
		// switch driver {
		// case opts.SQLDriverPGXV5:
		// 	return "pgtype.Multirange[pgtype.Range[pgtype.Date]]"
		// default:
		// 	return "interface{}"
		// }

	case "tsrange":
		// switch driver {
		// case opts.SQLDriverPGXV4:
		// 	return "pgtype.Tsrange"
		// case opts.SQLDriverPGXV5:
		// 	return "pgtype.Range[pgtype.Timestamp]"
		// default:
		// 	return "interface{}"
		// }

	case "tsmultirange":
		// switch driver {
		// case opts.SQLDriverPGXV5:
		// 	return "pgtype.Multirange[pgtype.Range[pgtype.Timestamp]]"
		// default:
		// 	return "interface{}"
		// }

	case "tstzrange":
		// switch driver {
		// case opts.SQLDriverPGXV4:
		// 	return "pgtype.Tstzrange"
		// case opts.SQLDriverPGXV5:
		// 	return "pgtype.Range[pgtype.Timestamptz]"
		// default:
		// 	return "interface{}"
		// }

	case "tstzmultirange":
		// switch driver {
		// case opts.SQLDriverPGXV5:
		// 	return "pgtype.Multirange[pgtype.Range[pgtype.Timestamptz]]"
		// default:
		// 	return "interface{}"
		// }

	case "numrange":
		// switch driver {
		// case opts.SQLDriverPGXV4:
		// 	return "pgtype.Numrange"
		// case opts.SQLDriverPGXV5:
		// 	return "pgtype.Range[pgtype.Numeric]"
		// default:
		// 	return "interface{}"
		// }

	case "nummultirange":
		// switch driver {
		// case opts.SQLDriverPGXV5:
		// 	return "pgtype.Multirange[pgtype.Range[pgtype.Numeric]]"
		// default:
		// 	return "interface{}"
		// }

	case "int4range":
		// switch driver {
		// case opts.SQLDriverPGXV4:
		// 	return "pgtype.Int4range"
		// case opts.SQLDriverPGXV5:
		// 	return "pgtype.Range[pgtype.Int4]"
		// default:
		// 	return "interface{}"
		// }

	case "int4multirange":
		// switch driver {
		// case opts.SQLDriverPGXV5:
		// 	return "pgtype.Multirange[pgtype.Range[pgtype.Int4]]"
		// default:
		// 	return "interface{}"
		// }

	case "int8range":
		// switch driver {
		// case opts.SQLDriverPGXV4:
		// 	return "pgtype.Int8range"
		// case opts.SQLDriverPGXV5:
		// 	return "pgtype.Range[pgtype.Int8]"
		// default:
		// 	return "interface{}"
		// }

	case "int8multirange":
		// switch driver {
		// case opts.SQLDriverPGXV5:
		// 	return "pgtype.Multirange[pgtype.Range[pgtype.Int8]]"
		// default:
		// 	return "interface{}"
		// }

	case "hstore":
		// if driver.IsPGX() {
		// 	return "pgtype.Hstore"
		// }
		return protobuilder.FieldTypeImportedMessage(
			tAny,
		), nil

	case "bit", "varbit", "pg_catalog.bit", "pg_catalog.varbit":
		// if driver == opts.SQLDriverPGXV5 {
		// 	return "pgtype.Bits"
		// }
		// if driver == opts.SQLDriverPGXV4 {
		// 	return "pgtype.Varbit"
		// }

	case "cid":
		// if driver == opts.SQLDriverPGXV5 {
		// 	return "pgtype.Uint32"
		// }
		// if driver == opts.SQLDriverPGXV4 {
		// 	return "pgtype.CID"
		// }

	case "oid":
		// if driver == opts.SQLDriverPGXV5 {
		// 	return "pgtype.Uint32"
		// }
		// if driver == opts.SQLDriverPGXV4 {
		// 	return "pgtype.OID"
		// }

	case "tid":
		// if driver.IsPGX() {
		// 	return "pgtype.TID"
		// }

	case "xid":
		// if driver == opts.SQLDriverPGXV5 {
		// 	return "pgtype.Uint32"
		// }
		// if driver == opts.SQLDriverPGXV4 {
		// 	return "pgtype.XID"
		// }

	case "box":
		// if driver.IsPGX() {
		// 	return "pgtype.Box"
		// }

	case "circle":
		// if driver.IsPGX() {
		// 	return "pgtype.Circle"
		// }

	case "line":
		// if driver.IsPGX() {
		// 	return "pgtype.Line"
		// }

	case "lseg":
		// if driver.IsPGX() {
		// 	return "pgtype.Lseg"
		// }

	case "path":
		// if driver.IsPGX() {
		// 	return "pgtype.Path"
		// }

	case "point":
		// if driver.IsPGX() {
		// 	return "pgtype.Point"
		// }

	case "polygon":
		// if driver.IsPGX() {
		// 	return "pgtype.Polygon"
		// }

	case "vector":
		// if driver == opts.SQLDriverPGXV5 {
		// 	if emitPointersForNull {
		// 		return "*pgvector.Vector"
		// 	} else {
		// 		return "pgvector.Vector"
		// 	}
		// }

	case "void":
		// A void value can only be scanned into an empty interface.
		return protobuilder.FieldTypeImportedMessage(
			tAny,
		), nil

	case "any", WellKnownAny:
		return protobuilder.FieldTypeImportedMessage(
			tAny,
		), nil

	// If we are here check for enum
	default:
		tKind, err := stringToKind(ct)
		if err == nil {
			return protobuilder.FieldTypeScalar(tKind), nil
		}

		for _, f := range p.files {
			// We have to handle for when UserDefined comes in with EnumType
			// Ex:  foo.bar.baz.v1.$EnumType
			ct = filepath.Base(pkgToPath(ct))
			tName := protoreflect.Name(*toPascal(ct))
			// Check for Enums
			eb := f.GetEnum(tName)
			if eb != nil {
				return protobuilder.FieldTypeEnum(eb), nil
			}
			// Check For Messages
			mb := f.GetMessage(tName)
			if mb != nil {
				return protobuilder.FieldTypeMessage(mb), nil
			}
		}
	}

	return nil, fmt.Errorf("%s: Type Conversion Not Implemented.  Use --replace: or --skip:", ct)
}

func stringToKind(typeString string) (protoreflect.Kind, error) {
	switch strings.ToUpper(typeString) {
	case "TYPE_DOUBLE":
		return protoreflect.DoubleKind, nil
	case "TYPE_FLOAT":
		return protoreflect.FloatKind, nil
	case "TYPE_INT64":
		return protoreflect.Int64Kind, nil
	case "TYPE_UINT64":
		return protoreflect.Uint64Kind, nil
	case "TYPE_INT32":
		return protoreflect.Int32Kind, nil
	case "TYPE_FIXED64":
		return protoreflect.Fixed64Kind, nil
	case "TYPE_FIXED32":
		return protoreflect.Fixed32Kind, nil
	case "TYPE_BOOL":
		return protoreflect.BoolKind, nil
	case "TYPE_STRING":
		return protoreflect.StringKind, nil
	case "TYPE_GROUP":
		return protoreflect.GroupKind, nil
	case "TYPE_MESSAGE":
		return protoreflect.MessageKind, nil
	case "TYPE_BYTES":
		return protoreflect.BytesKind, nil
	case "TYPE_UINT32":
		return protoreflect.Uint32Kind, nil
	case "TYPE_ENUM":
		return protoreflect.EnumKind, nil
	case "TYPE_SFIXED32":
		return protoreflect.Sfixed32Kind, nil
	case "TYPE_SFIXED64":
		return protoreflect.Sfixed64Kind, nil
	case "TYPE_SINT32":
		return protoreflect.Sint32Kind, nil
	case "TYPE_SINT64":
		return protoreflect.Sint64Kind, nil
	// Add more cases for different types as needed
	default:
		return protoreflect.Kind(0), fmt.Errorf("unknown type: %s", typeString)
	}
}

// Takes a EnumName and Values and PREFIXS them in the proto style
// Example: ResourceType_UNSPECIFIED
func convertEnumValues(n string, s []string) []string {
	u := strcase.ToSNAKE(n)
	x := []string{u + "_UNSPECIFIED"}
	for _, z := range s {
		y := strcase.ToSNAKE(z)
		x = append(x, (u + "_" + y))
	}
	return x
}

type ReqResp struct {
	OneOf     *[]string
	ReqFields map[string]string
	RespEmpty map[string]bool
	a         *Annotations
}

type Service struct {
	Path *url.URL
	Name string
	a    *Annotations
}

type httpOptions struct {
	Method string
	Body   string
	Path   *url.URL
}

func parseDynamicPath(path string) []string {
	// Split both the template and path into segments
	pathParts := strings.Split(path, "/")

	// Create a map to store parameter values
	var params []string

	// Regex to identify path parameters (enclosed in { })
	re := regexp.MustCompile(`^{([^}]+)}$`)

	for _, pathPart := range pathParts {
		if matches := re.FindStringSubmatch(pathPart); len(matches) > 0 {
			// If the path segment is a parameter (e.g., {org}), extract the name
			paramName := matches[1]
			params = append(params, paramName)
		}
	}

	return params
}

const (
	wkprefix = "google.protobuf."

	// proto represents builtin types
	protoDouble   = "double"
	protoFloat    = "float"
	protoInt32    = "int32"
	protoInt64    = "int64"
	protoUint32   = "uint32"
	protoUint64   = "uint64"
	protoSint32   = "sint32"
	protoSint64   = "sint64"
	protoFixed32  = "fixed32"
	protoFixed64  = "fixed64"
	protoSFixed32 = "sfixed32"
	protoSFixed64 = "sfixed64"
	protoBool     = "bool"
	protoString   = "string"
	protoBytes    = "bytes"

	// well known represents google.proto. well known types
	WellKnownAny              = wkprefix + "Any"
	WellKnownBoolValue        = wkprefix + "BoolValue"
	WellKnownBytesValue       = wkprefix + "BytesValue"
	WellKnownDecimal          = wkprefix + "Decimal"
	WellKnownDoubleValue      = wkprefix + "DoubleValue"
	WellKnownDuration         = wkprefix + "Duration"
	WellKnownEmpty            = wkprefix + "Empty"
	WellKnownEnum             = wkprefix + "Enum"
	WellKnownEnumValue        = wkprefix + "EnumValue"
	WellKnownField            = wkprefix + "Field"
	WellKnownFieldCardinality = wkprefix + "Field.Cardinality"
	WellKnownFieldKind        = wkprefix + "Field.Kind"
	WellKnownFieldMask        = wkprefix + "FieldMask"
	WellKnownFloatValue       = wkprefix + "FloatValue"
	WellKnownInt32Value       = wkprefix + "Int32Value"
	WellKnownInt64Value       = wkprefix + "Int64Value"
	WellKnownListValue        = wkprefix + "ListValue"
	WellKnownMethod           = wkprefix + "Method"
	WellKnownMixin            = wkprefix + "Mixin"
	WellKnownMoney            = wkprefix + "Money"
	WellKnownNullValue        = wkprefix + "NullValue"
	WellKnownOption           = wkprefix + "Option"
	WellKnownSourceContext    = wkprefix + "SourceContext"
	WellKnownStringValue      = wkprefix + "StringValue"
	WellKnownStruct           = wkprefix + "Struct"
	WellKnownSyntax           = wkprefix + "Syntax"
	WellKnownTimestamp        = wkprefix + "Timestamp"
	WellKnownType             = wkprefix + "Type"
	WellKnownUInt32Value      = wkprefix + "UInt32Value"
	WellKnownUInt64Value      = wkprefix + "UInt64Value"
)
