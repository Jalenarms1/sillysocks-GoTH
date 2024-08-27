// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.771
package shared

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

func ProductSlides() templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
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
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"carousel carousel-center bg-zinc-100 rounded-md space-x-4 p-4\"><div class=\"carousel-item\"><img src=\"https://img.daisyui.com/images/stock/photo-1559703248-dcaaec9fab78.webp\" class=\"rounded-box\"></div><div class=\"carousel-item\"><img src=\"https://img.daisyui.com/images/stock/photo-1565098772267-60af42b81ef2.webp\" class=\"rounded-box\"></div><div class=\"carousel-item\"><img src=\"https://img.daisyui.com/images/stock/photo-1572635148818-ef6fd45eb394.webp\" class=\"rounded-box\"></div><div class=\"carousel-item\"><img src=\"https://img.daisyui.com/images/stock/photo-1494253109108-2e30c049369b.webp\" class=\"rounded-box\"></div><div class=\"carousel-item\"><img src=\"https://img.daisyui.com/images/stock/photo-1550258987-190a2d41a8ba.webp\" class=\"rounded-box\"></div><div class=\"carousel-item\"><img src=\"https://img.daisyui.com/images/stock/photo-1559181567-c3190ca9959b.webp\" class=\"rounded-box\"></div><div class=\"carousel-item\"><img src=\"https://img.daisyui.com/images/stock/photo-1601004890684-d8cbf643f5f2.webp\" class=\"rounded-box\"></div></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

var _ = templruntime.GeneratedTemplate
