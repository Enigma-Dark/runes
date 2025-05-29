package templates

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed builtin/*.tmpl
var builtinTemplates embed.FS

// Manager handles template registration and retrieval
type Manager struct {
	templates map[string]*template.Template
}

// NewManager creates a new template manager
func NewManager() *Manager {
	return &Manager{
		templates: make(map[string]*template.Template),
	}
}

// LoadBuiltinTemplates loads all builtin templates from embedded files
func (m *Manager) LoadBuiltinTemplates() error {
	return fs.WalkDir(builtinTemplates, "builtin", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || !strings.HasSuffix(path, ".tmpl") {
			return nil
		}

		content, err := builtinTemplates.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read template %s: %w", path, err)
		}

		// Extract template name from filename (remove .tmpl extension)
		name := strings.TrimSuffix(filepath.Base(path), ".tmpl")

		tmpl, err := template.New(name).Parse(string(content))
		if err != nil {
			return fmt.Errorf("failed to parse template %s: %w", name, err)
		}

		m.templates[name] = tmpl
		return nil
	})
}

// LoadExternalTemplate loads a template from an external file
func (m *Manager) LoadExternalTemplate(name, filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read external template %s: %w", filePath, err)
	}

	tmpl, err := template.New(name).Parse(string(content))
	if err != nil {
		return fmt.Errorf("failed to parse external template %s: %w", name, err)
	}

	m.templates[name] = tmpl
	return nil
}

// GetTemplate retrieves a template by name
func (m *Manager) GetTemplate(name string) (*template.Template, error) {
	tmpl, exists := m.templates[name]
	if !exists {
		return nil, fmt.Errorf("template '%s' not found. Available templates: %s",
			name, strings.Join(m.ListTemplates(), ", "))
	}
	return tmpl, nil
}

// ListTemplates returns a list of all available template names
func (m *Manager) ListTemplates() []string {
	var names []string
	for name := range m.templates {
		names = append(names, name)
	}
	return names
}

// HasTemplate checks if a template exists
func (m *Manager) HasTemplate(name string) bool {
	_, exists := m.templates[name]
	return exists
}

// RegisterTemplate manually registers a template
func (m *Manager) RegisterTemplate(name string, tmpl *template.Template) {
	m.templates[name] = tmpl
}
