package main

import (
	"errors"

	ovirt "github.com/ovirt/go-ovirt"
)

func getTemplate(conn *ovirt.Connection, templateName string, templateVersion int64) (ovirt.Template, error) {
	template := ovirt.Template{}

	templatesService := conn.SystemService().TemplatesService()
	tempResp, err := templatesService.List().Send()
	if err != nil {
		return template, err
	}

	tempsSlice, ok := tempResp.Templates()
	if !ok {
		return template, errors.New("couldn't lookup list of templates")
	}

	for _, _template := range tempsSlice.Slice() {
		// name
		_name, ok := _template.Name()
		if !ok {
			return template, errors.New("couldn't get name of template")
		}

		// get version
		_tempVersion, ok := _template.Version()
		if !ok {
			return template, errors.New("couldn't get template version object")
		}

		// get version number
		_tempVersionNum, ok := _tempVersion.VersionNumber()
		if !ok {
			return template, errors.New("couldn't get template version number")
		}

		// name and version match; get actual template object
		if _name == templateName && _tempVersionNum == templateVersion {
			template = *_template
		}
	}

	return template, nil
}
