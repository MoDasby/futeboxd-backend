package templates

import "fmt"

func NewRecoverPasswordTemplate(username, url string) string {
	return fmt.Sprintf(`
    <!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
    <html dir="ltr" lang="pt-BR">
    <head>
        <meta content="text/html; charset=UTF-8" http-equiv="Content-Type" />
        <meta name="x-apple-disable-message-reformatting" />
    </head>
    <div
        style="display:none;overflow:hidden;line-height:1px;opacity:0;max-height:0;max-width:0"
    >
        Futeboxd - Redefinição de senha
    </div>
    <body style="background-color:#f6f9fc;padding:10px 0">
        <table
        align="center"
        width="100%%"
        border="0"
        cellpadding="0"
        cellspacing="0"
        role="presentation"
        style="max-width:37.5em;background-color:#ffffff;border:1px solid #f0f0f0;padding:45px"
        >
        <tbody>
            <tr style="width:100%%">
            <td>
                <table
                align="center"
                width="100%%"
                border="0"
                cellpadding="0"
                cellspacing="0"
                role="presentation"
                >
                <tbody>
                    <tr>
                    <td>
                        <p
                        style="font-size:16px;line-height:26px;margin:16px 0;font-family:&#x27;Open Sans&#x27;, &#x27;HelveticaNeue-Light&#x27;, &#x27;Helvetica Neue Light&#x27;, &#x27;Helvetica Neue&#x27;, Helvetica, Arial, &#x27;Lucida Grande&#x27;, sans-serif;font-weight:300;color:#404040"
                        >
                        Olá,
                        <!-- -->%s<!-- -->,
                        </p>
                        <p
                        style="font-size:16px;line-height:26px;margin:16px 0;font-family:&#x27;Open Sans&#x27;, &#x27;HelveticaNeue-Light&#x27;, &#x27;Helvetica Neue Light&#x27;, &#x27;Helvetica Neue&#x27;, Helvetica, Arial, &#x27;Lucida Grande&#x27;, sans-serif;font-weight:300;color:#404040"
                        >
                        Alguém solicitou recentemente uma redefinição de senha para sua conta no Futeboxd. Se foi você, pode definir uma nova senha clicando no botão abaixo:
                        </p>
                        <a
                        href="%s"
                        style="line-height:100%%;text-decoration:none;display:block;max-width:100%%;mso-padding-alt:0px;background-color:#26ae5d;border-radius:4px;color:#fff;font-family:&#x27;Open Sans&#x27;, &#x27;Helvetica Neue&#x27;, Arial;font-size:15px;text-align:center;width:210px;padding:14px 7px 14px 7px"
                        target="_blank"
                        ><span
                            ><!--[if mso
                            ]><i
                                style="mso-font-width:350%%;mso-text-raise:21"
                                hidden
                                >&#8202;</i
                            ><!
                            [endif]--></span
                        ><span
                            style="max-width:100%%;display:inline-block;line-height:120%%;mso-padding-alt:0px;mso-text-raise:10.5px"
                            >Redefinir senha</span
                        ><span
                            ><!--[if mso
                            ]><i style="mso-font-width:350%%" hidden
                                >&#8202;&#8203;</i
                            ><!
                            [endif]--></span
                        ></a
                        >
                        <p
                        style="font-size:16px;line-height:26px;margin:16px 0;font-family:&#x27;Open Sans&#x27;, &#x27;HelveticaNeue-Light&#x27;, &#x27;Helvetica Neue Light&#x27;, &#x27;Helvetica Neue&#x27;, Helvetica, Arial, &#x27;Lucida Grande&#x27;, sans-serif;font-weight:300;color:#404040"
                        >
                        Se você não solicitou essa alteração ou não deseja mudar sua senha, ignore e exclua esta mensagem.
                        </p>
                        <p
                        style="font-size:16px;line-height:26px;margin:16px 0;font-family:&#x27;Open Sans&#x27;, &#x27;HelveticaNeue-Light&#x27;, &#x27;Helvetica Neue Light&#x27;, &#x27;Helvetica Neue&#x27;, Helvetica, Arial, &#x27;Lucida Grande&#x27;, sans-serif;font-weight:300;color:#404040"
                        >
                        Para manter sua conta segura, não encaminhe este e-mail para ninguém.
                        </p>
                        <p
                        style="font-size:16px;line-height:26px;margin:16px 0;font-family:&#x27;Open Sans&#x27;, &#x27;HelveticaNeue-Light&#x27;, &#x27;Helvetica Neue Light&#x27;, &#x27;Helvetica Neue&#x27;, Helvetica, Arial, &#x27;Lucida Grande&#x27;, sans-serif;font-weight:300;color:#404040"
                        >
                        Boas partidas e aproveite o Futeboxd!
                        </p>
                    </td>
                    </tr>
                </tbody>
                </table>
            </td>
            </tr>
        </tbody>
        </table>
    </body>
    </html>
`, username, url)
}
