# Rate Limiter Middleware

## Descrição
Este projeto implementa um middleware em Go que controla o número de requisições que podem ser feitas a cada segundo a partir de um IP ou token de acesso. Ele é útil para proteger APIs contra abusos e garantir que os recursos sejam utilizados de maneira justa.

## Documentação
Como Funciona
O middleware RateLimiterMiddleware controla o número de requisições que podem ser feitas a cada segundo a partir de um IP ou token de acesso. Se o limite de requisições for excedido, ele retorna um código 429 (Too Many Requests), com a mensagem designada.

## Para Testar
Use ferramentas como Postman ou cURL para simular requisições e observar o comportamento quando os limites são atingidos.

## Exemplo de Uso com cURL
curl -X GET http://localhost:8080/sua-endpoint -H "Authorization: Bearer seu_token"
Repita a requisição várias vezes para ver como o sistema responde quando o limite é atingido.
