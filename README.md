# funcmap

Go template functions.

    import "github.com/gomatic/funcmap"

...

    template.New(name).
        Funcs(funcmap.Map).
        Parse(templateSource).
        Execute(&result, templateVariables)
