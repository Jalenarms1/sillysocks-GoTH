// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.747
package layouts

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

func Footer() templ.Component {
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<footer class=\"footer bg-zinc-950 text-neutral-content p-10\"><div><h6 class=\"footer-title\">Services</h6><a class=\"link link-hover\">Branding</a> <a class=\"link link-hover\">Design</a> <a class=\"link link-hover\">Marketing</a> <a class=\"link link-hover\">Advertisement</a></div><div><h6 class=\"footer-title\">Company</h6><a class=\"link link-hover\">About us</a> <a class=\"link link-hover\">Contact</a> <a class=\"link link-hover\">Jobs</a> <a class=\"link link-hover\">Press kit</a></div><div class=\"flex justify-between w-full items-end\"><div class=\"flex flex-col gap-2\"><h6 class=\"footer-title\">Legal</h6><a class=\"link link-hover\">Terms of use</a> <a class=\"link link-hover\">Privacy policy</a> <a class=\"link link-hover\">Cookie policy</a></div><img src=\"/public/sockslogo.png\" alt=\"\" class=\"w-16 h-16\"></div><p class=\"text-sm mx-auto text-center\">&copy; 2024 Silly Socks and More. All rights reserved.</p></footer>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}
