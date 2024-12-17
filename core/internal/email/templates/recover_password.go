package templates

import "fmt"

func RecoverPasswordTemplate(username string, url string, token string) string {
	return fmt.Sprintf(`
    <!DOCTYPE html>
    <html>

    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
    </head>

    <body style="font-family: Arial, sans-serif; margin: 0; padding: 0;">
        <div style="max-width: 600px; margin: 0 auto; border-radius: 8px; overflow: hidden;">
            <table style="width: 100%%; border-spacing: 0; border-collapse: collapse;">
                <tr>
                    <td>
                        <div style="padding: 20px;">
                            <h1 style="margin: 0 0 16px;">Olá, %s</h1>
                            <p style="margin: 0 0 16px;">
                                Recebemos uma solicitação para redefinir a senha da sua conta.
                                Se foi você que solicitou, clique no botão abaixo para redefinir sua senha:
                            </p>
                            <a href="%s" style="display: inline-block; padding: 12px 20px;">
                                Redefinir Senha
                            </a>

                            <span>
                                Aqui está o token %s
                            </span>
                            <p style="margin: 16px 0 0;">
                                Se você não fez essa solicitação, ignore este email.
                            </p>
                        </div>
                    </td>
                </tr>
                <tr>
                    <td>
                        <footer style="text-align: center; padding: 10px; font-size: 12px;">
                            <p style="margin: 0;">Este é um email automático. Por favor, não responda.</p>
                        </footer>
                    </td>
                </tr>
            </table>
        </div>
    </body>

    </html>
`, username, url, token)
}
