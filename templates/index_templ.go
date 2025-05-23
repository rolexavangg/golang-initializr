// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.865
package templates

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import (
	"fmt"
)

type Dependency struct {
	ID          string
	Name        string
	Description string
	Category    string
}

type ProjectForm struct {
	Name         string
	Dependencies []Dependency
}

func Index() templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Var2 := templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
			templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
			templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
			if !templ_7745c5c3_IsBuffer {
				defer func() {
					templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
					if templ_7745c5c3_Err == nil {
						templ_7745c5c3_Err = templ_7745c5c3_BufErr
					}
				}()
			}
			ctx = templ.InitializeContext(ctx)
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<div class=\"hero\"><h1>Golang Initializr</h1><p>Quickly generate Go project skeleton with the dependencies you need</p></div><div class=\"project-form\"><form id=\"project-form\" action=\"/generate\" method=\"post\"><div class=\"form-group\"><label for=\"project-name\">Project Name</label> <input type=\"text\" id=\"project-name\" name=\"name\" placeholder=\"github.com/username/project\" required></div><div class=\"dependencies-section\"><h2>Dependencies</h2><p class=\"note\">All projects include: Uber FX, Zap Logger, Clean Architecture</p><div class=\"dependency-categories\"><div class=\"category\"><h3>Databases</h3><div class=\"dependency-list\"><div class=\"dependency-item\"><input type=\"checkbox\" id=\"postgres\" name=\"dependencies\" value=\"postgres\"> <label for=\"postgres\">PostgreSQL</label></div><div class=\"dependency-item\"><input type=\"checkbox\" id=\"redis\" name=\"dependencies\" value=\"redis\"> <label for=\"redis\">Redis</label></div></div></div><div class=\"category\"><h3>Messaging</h3><div class=\"dependency-list\"><div class=\"dependency-item\"><input type=\"checkbox\" id=\"kafka\" name=\"dependencies\" value=\"kafka\"> <label for=\"kafka\">Kafka</label></div></div></div><div class=\"category\"><h3>API</h3><div class=\"dependency-list\"><div class=\"dependency-item\"><input type=\"checkbox\" id=\"http\" name=\"dependencies\" value=\"http\" checked> <label for=\"http\">HTTP (Echo)</label></div><div class=\"dependency-item\"><input type=\"checkbox\" id=\"grpc\" name=\"dependencies\" value=\"grpc\"> <label for=\"grpc\">gRPC</label></div></div></div><div class=\"category\"><h3>Tools</h3><div class=\"dependency-list\"><div class=\"dependency-item\"><input type=\"checkbox\" id=\"docker\" name=\"dependencies\" value=\"docker\" checked> <label for=\"docker\">Docker</label></div></div></div></div></div><div class=\"form-actions\"><button type=\"submit\" class=\"btn-primary\">Generate Project</button></div></form><!-- Form submits directly to generate endpoint for immediate download --></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			return nil
		})
		templ_7745c5c3_Err = Layout("Home").Render(templ.WithChildren(ctx, templ_7745c5c3_Var2), templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

func downloadLink(projectName string, dependencies []string) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var3 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var3 == nil {
			templ_7745c5c3_Var3 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 2, "<a href=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var4 templ.SafeURL = templ.SafeURL("/download?session=" + projectName + "-" + fmt.Sprint(len(dependencies)))
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(string(templ_7745c5c3_Var4)))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 3, "\" class=\"btn-download\">Download Project</a>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

// No longer needed as we directly return the ZIP file
var _ = templruntime.GeneratedTemplate
