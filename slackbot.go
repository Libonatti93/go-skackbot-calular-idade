package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "time"
    
    "github.com/slack-go/slack"
)

func main() {
    // Obtém o token do bot do Slack a partir da variável de ambiente
    slackToken := os.Getenv("SLACK_BOT_TOKEN")
    if slackToken == "" {
        log.Fatal("SLACK_BOT_TOKEN não definido")
    }

    // Inicializa o cliente Slack
    api := slack.New(slackToken)

    // Define o endereço do servidor
    http.HandleFunc("/slack/events", func(w http.ResponseWriter, r *http.Request) {
        // Verifica o tipo de solicitação
        if r.Method != http.MethodPost {
            w.WriteHeader(http.StatusMethodNotAllowed)
            return
        }

        // Verifica o token de verificação do Slack (opcional, mas recomendado)
        slackVerificationToken := os.Getenv("SLACK_VERIFICATION_TOKEN")
        s, err := slack.SlashCommandParse(r)
        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            return
        }

        if !s.ValidateToken(slackVerificationToken) {
            w.WriteHeader(http.StatusUnauthorized)
            return
        }

        // Trata o comando do Slack
        if s.Command == "/idade" {
            response := calculateAge(s.Text)
            w.Header().Set("Content-Type", "application/json")
            w.Write([]byte(response))
        }
    })

    fmt.Println("[INFO] Server listening on :3000")
    log.Fatal(http.ListenAndServe(":3000", nil))
}

// Função para calcular a idade
func calculateAge(input string) string {
    layout := "2006-01-02" // Formato esperado da data: YYYY-MM-DD
    birthdate, err := time.Parse(layout, input)
    if err != nil {
        return fmt.Sprintf("Formato de data inválido. Por favor, use o formato YYYY-MM-DD.")
    }

    today := time.Now()
    age := today.Year() - birthdate.Year()
    
    if today.YearDay() < birthdate.YearDay() {
        age--
    }

    return fmt.Sprintf("Você tem %d anos.", age)
}
