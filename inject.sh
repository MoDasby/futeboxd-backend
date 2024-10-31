#!/bin/bash

# Lista de tokens de usuários
tokens=(
    "a148cf5d368667e8514ca6375c934b6683fe2e1f5e8a24235e6dcd4c5e4dcd18db1063b63942cf083e95e733d19c1b26"
    "90773770612cf6e2dc2f1c0fda069f6de0c80c60f84d9a6c1afef5007c16d606aaaa3d0da2f8404939ef1e50c78c28c0"
    "1c825048803486dc7d5a2d6012719e273d475140ac86a66f5cfc07df76debdfd0b8736cd53ec27dce0b7f444f405e016"
    "9620b6ad11dbbb1cfa9c27a122c9ae7af03e08ade770313060843004e4bba591fd9957cc2a08069d4eb4bf37a87c41fe"
    "c71e26e4927a9b60c3c67195f32e472ef6b1bab532387b8bb283a294f93053689dfa8745fe537c2d1694c2b09178cc68"
    "6d891313acf9d8619218ad4a1e5fa265fcc02feae83d16072642c731867f234b6598ba0c22695f271695dad2696b9f2e"
    "7834de975ca15b85c57d0bd983bda5c53a4d868e9f504bb6322165effbb4c4b6ee902e1824b512c443e6ec50456a5205"
    "af7b2b4d9849110866a54827306350190da8030f65e15a9ff8463a9aacad9d90986812e7c685427370790053cde4a6c2"
    "d17208bd28a2ab4f279a64dda9e6a9c20550ad58674ba19fc45ec2e4a094136a52067cbc669d414d6fc28bf78bccc840"
    "39ab0d5411ec72398c83164d7c632f7a8ee6a47dae9385b3cfe704d6fdebaef7e55968b3b5cd898556184b75274ce0a5"
    "c97c73125917dd94eebcbe58c251ee2afd0a81079edaa06dcfd4935900a48bcf9e5e24281bcfc6a0e9182a5dfdb6b594"
    "744f13e057ea91c4a09d239f626601dc39b6956deffd6f5ac008bbcca9049007e036b68632b7b94a19d2d1dbe6817bbc"
    "66268a4edd6090ecdb57934582547cdb3ea32942bededa3d262f6509285b469d6d382ffc1b059e56ed1cece8ba44d345"
    "03f51c6ffbed9ad21a5fb2bd896c69b3357e31081953d7578e0b6277fbe321a30fec8adf9c83f25f4f78525837af18dc"
    "6c8501e732a37370d634585b406c41ea7de5164e9d82351c25cfd4c3cae25bd4774ba5a1fa1dcc7b38593cdad2539fac"
    "bf74c93c77d14c30682e9d9db576f404bf5cf757524b1de4d9486d14b0439fc4c189f1bf39a048501ab4237c70ec3384"
    "8360c20feabf79682cf4e4ad950b50f71cfb6e6cd69c14f08bc4d5963ebffc360b8d40bb9e66b82b4d5c4aeeb980edc5"
    "a3ddb9f1bde8d8e2add1a3817882a892e6bb9c8ab646ed765a72c52728480d0ae6f11fd71b541d527741d6cc46de37c6"
    "22c7a637c05b58b8df89cbbfae7c4dd00c991b3ba6e809d37ba867c79dde0cbdeec4b02616f73ed5eaba98607a9e5ec9"
    "4aa13f180731ff021d0b351b5d06dd0292ccb4a3109d33c33bb1b3ab2cc27d9ee5c0a0eaed3c1f5a7009fb62af9d5836"
    "c6ae959b1c9097f7e89038fa82b7f4c55925c041acbcc054a0eebc52a13c4cfd9ef0f5309f3bb5cb866bc008bd7fbdde"
    "0dfa01b411b228962ad8142a1b5f54c164f84487f81a5c0b3c5802ce947498d613c56874dd1c005a6a40e9760e25afcb"
    "63ea5c5aca0ed72e2e5c16de6ccbe56ec9d05a1b582427f77b52ec0f6cfcca5fc1f0e3543d28c6a81aec378701df58e8"
    "f39a0b4caf671df3413539d1215a86d97bbafc4380ddb2dbbdaecf80ed9d10d9e2a4958e121e69eb35499160f7907052"
    "358d529a4f23c4abd060f5e75468a237b3dda8d434a2a06074e895e5234d94462dced6aaade0737398dd54f4910fdbb0"
    "89ae6ed65b4b9ac698e72c2387f7ea1bfe9cf1e63a4cd32c135a8f02b6164f3d8e85a33c73a140d212af4702144e9e7c"
)


# Lista de partidas
matches=(
  699645  # Fortaleza 1 x Atletico MG 1
  699642  # Corinthians 5 x 2 Athletico PR
  699633  # Corinthians 2 x 2 Internacional
  704740  # Alaves 0 x 3 Barcelona
  704739  # Osasuna 4 x 2 Barcelona
  699636  # Fluminense 1 x 0 Cruzeiro
  684624  # Brasil 4 x 0 Peru
  721537
  721590
  721590
  694415
  704161
  694414
  694422
  720548
  708324
)

# Lista de descrições
descriptions=(
  "Foi uma ótima partida!"
  "Achei o jogo muito disputado!"
  "Bom desempenho, mas poderia ser melhor."
  "Uma atuação incrível do time visitante."
  "Partida emocionante do começo ao fim!"
  "Jogo fraco, mas valeu pelos gols."
  "A defesa foi o destaque."
)

# Função para gerar um número aleatório entre 1 e 5 para o rate
random_rate() {
  echo $(( (RANDOM % 5) + 1 ))
}

# Loop para enviar requisições
for i in {1..100}; do
  # Seleciona aleatoriamente um token, partida e descrição
  token=${tokens[$RANDOM % ${#tokens[@]}]}
  match=${matches[$RANDOM % ${#matches[@]}]}
  description=${descriptions[$RANDOM % ${#descriptions[@]}]}
  rate=$(random_rate)

  # Envia a requisição
  curl -X POST http://localhost:8080/reviews \
    -H "Content-Type: application/json" \
    -H "Authorization: $token" \
    -d "{
          \"rate\": $rate,
          \"description\": \"$description\",
          \"match_id\": $match
        }"
  
  echo "Review enviada para a partida $match com o rate $rate"
done
