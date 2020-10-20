package dstore

import "path/filepath"

//TemplateSubPath 获取模板路径
func TemplateSubPath(root, kind, name string) string {
	return filepath.Join(root, "template", kind, name)
}
