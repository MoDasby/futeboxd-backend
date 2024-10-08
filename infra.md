# Serviço de Autenticação e Autorização
## Responsabilidades:
- Gerenciar criação de contas e login.
- Emitir e validar tokens JWT.
- Manter o controle de permissões e acessos (se necessário).
## Endpoints:
- POST /auth/register: Criar conta.
- POST /auth/login: Realizar login.
- POST /auth/refresh-token: Renovar o token.
- GET /auth/user: Retornar informações sobre o usuário autenticado.
## Banco de Dados:
- Tabelas relacionadas a usuários, como users e tokens.
# Serviço de Times e Partidas
## Responsabilidades:
- Gerenciar a pesquisa de times e partidas.
- Buscar dados de partidas de um time específico.
- Integração com APIs de terceiros (ex: API da ESPN) para buscar dados atualizados.
- Endpoints:
- GET /teams/search: Pesquisar times pelo nome.
- GET /teams/{teamId}/matches: Listar partidas de um time.
- GET /matches/{matchId}: Obter detalhes de uma partida específica.
## Banco de Dados:
- Tabelas para armazenar informações de cache ou dados locais sobre times e partidas (opcional, caso dados venham de uma API externa).
# Serviço de Reviews
## Responsabilidades:
- Gerenciar a criação e visualização de reviews de partidas.
- Relacionar reviews com partidas e usuários.
- Permitir que os usuários curtam reviews.
## Endpoints:
- POST /reviews: Criar review de uma partida.
- GET /reviews/{matchId}: Obter reviews de uma partida específica.
- POST /reviews/{reviewId}/like: Curtir uma review.
- DELETE /reviews/{reviewId}/unlike: Descurtir uma review.
## Banco de Dados:
- Tabelas para armazenar reviews, likes, match_reviews, etc.
# Serviço de Usuários e Relacionamentos
## Responsabilidades:
- Gerenciar o perfil do usuário, como seus dados pessoais, reviews criados, e os relacionamentos (seguir/desseguir usuários).
- Manter informações sobre quem está seguindo quem.
## Endpoints:
- GET /users/{userId}: Obter perfil de um usuário.
- POST /users/{userId}/follow: Seguir um usuário.
- DELETE /users/{userId}/unfollow: Deixar de seguir um usuário.
- GET /users/{userId}/followers: Listar seguidores de um usuário.
- GET /users/{userId}/following: Listar usuários que um determinado usuário está seguindo.
## Banco de Dados:
- Tabelas para armazenar os dados de perfil, relações de seguidores, histórico de atividades, etc.
# Serviço de Notificações (Opcional)
## Responsabilidades:
- Gerenciar o envio de notificações para os usuários quando eles recebem curtidas nas reviews, quando alguém os segue, etc.
- Integração com sistemas de notificação por e-mail, push ou no próprio app.
## Endpoints:
- POST /notifications/send: Enviar uma notificação para um usuário.
- GET /notifications/{userId}: Listar notificações de um usuário.
## Banco de Dados:
- Tabelas para armazenar notifications, notification_status (lida/não lida), etc.

Diagrama Geral de Microsserviços
Autenticação e Autorização
Focado no gerenciamento de usuários e controle de acesso.
Pode ser responsável por emitir e validar tokens.
Times e Partidas

Serviço que se conecta com APIs externas para buscar dados de times e partidas, cacheando informações quando necessário.
Reviews

Focado na criação e interação com as reviews. Gera interações como curtidas.
Usuários e Relacionamentos

Serviço que gerencia o relacionamento entre os usuários (seguir/desseguir), além de perfis de usuário.
Notificações (Opcional)

Serviço dedicado para notificações que podem ser enviadas via e-mail ou push notifications.