package utils

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin/binding"
)

// jsonBinding implementa a interface binding.Binding do Gin.
// Ela permite que usemos nossa própria lógica de Unmarshal (Sonic) 
// em vez da padrão do Go.
type jsonBinding struct{}

func (jsonBinding) Name() string {
	return "json"
}

func (jsonBinding) Bind(req *http.Request, obj interface{}) error {
	if req == nil || req.Body == nil {
		return io.EOF
	}
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}

	// Reutiliza a função Unmarshal do seu arquivo utils/json.go
	return Unmarshal(body, obj)
}

// NewJsonBinding cria uma nova instância do nosso binder.
func NewJsonBinding() binding.Binding {
	return jsonBinding{}
}
