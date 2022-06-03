package emails

import "html/template"

var email = template.Must(template.New("vcm_template").Parse(`html file here`))
