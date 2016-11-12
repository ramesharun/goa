package dsl

import (
	"github.com/goadesign/goa/design"
	"github.com/goadesign/goa/eval"
)

// Metadata is a set of key/value pairs that can be assigned to an object. Each value consists of a
// slice of strings so that multiple invocation of the Metadata function on the same target using
// the same key builds up the slice. Metadata may be set on fields, endpoints, responses, endpoint
// groups, and service expressions.
//
// While keys can have any value the following names are handled explicitly by goagen when set on
// fields.
//
// `struct:field:name`: overrides the Go struct field name generated by exprault by goagen.
// Applicable to fields only.
//
//        Metadata("struct:field:name", "MyName")
//
// `struct:tag:xxx`: sets the struct field tag xxx on generated Go structs.  Overrides tags that
// goagen would otherwise set.  If the metadata value is a slice then the strings are joined with
// the space character as separator.
// Applicable to fields only.
//
//        Metadata("struct:tag:json", "myName,omitempty")
//        Metadata("struct:tag:xml", "myName,attr")
//
// `swagger:tag:xxx`: sets the Swagger object field tag xxx.
// Applicable to resources and actions.
//
//        Metadata("swagger:tag:Backend")
//        Metadata("swagger:tag:Backend:desc", "Quick description of what 'Backend' is")
//        Metadata("swagger:tag:Backend:url", "http://example.com")
//        Metadata("swagger:tag:Backend:url:desc", "See more docs here")
//
// `swagger:summary`: sets the Swagger operation summary field.
// Applicable to actions.
//
//        Metadata("swagger:summary", "Short summary of what action does")
//
// `swagger:extension:xxx`: defines a swagger extension value.
// Applicable to all constructs that support Metadata.
//
//        Metadata("swagger:extension:x-apis-json", `{"URL": "http://goa.design"}`)
//
// The special key names listed above may be used as follows:
//
//        var Account = Type("Account", func() {
//                Attribute("service", String, "Name of service", func() {
//                        // Override exprault name to avoid clash with built-in 'Service' field.
//                        Metadata("struct:field:name", "ServiceName")
//                })
//        })
//
func Metadata(name string, value ...string) {
	appendMetadata := func(metadata design.MetadataExpr, name string, value ...string) design.MetadataExpr {
		if metadata == nil {
			metadata = make(map[string][]string)
		}
		metadata[name] = append(metadata[name], value...)
		return metadata
	}
	switch expr := eval.Current().(type) {
	case design.CompositeExpr:
		att := expr.Attribute()
		att.Metadata = appendMetadata(att.Metadata, name, value...)
	case *design.AttributeExpr:
		expr.Metadata = appendMetadata(expr.Metadata, name, value...)
	case *design.APIExpr:
		expr.Metadata = appendMetadata(expr.Metadata, name, value...)
	case *design.EndpointExpr:
		expr.Metadata = appendMetadata(expr.Metadata, name, value...)
	default:
		eval.IncompatibleDSL()
	}
}