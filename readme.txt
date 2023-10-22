Desafio Pós Go-Expert - Client-Server-API 

Client 
  - O Client está disponível na porta 8081/
  - Ao chamar http://localhost:8081 o retorno deve ser Dólar: 5.0312(cotacao do momento).
  - Ao receber Dólar: 5.0312, é gerado um arquivo cotacao.txt com este valor.

Server 
  - O Server está disponível na porta 8080/cotacao
  - Ao chamar http://localhost:8080/cotacao
  o retorno deve ser um JSON semelhante a este:
  {
    "code": "USD",
    "codein": "BRL",
    "name": "Dólar Americano/Real Brasileiro",
    "high": "5.0984",
    "low": "5.0255",
    "varBid": "-0.0319",
    "pctChange": "-0.63",
    "bid": "5.0312",
    "ask": "5.0342",
    "timestamp": "1697835596",
    "create_date": "2023-10-20 17:59:56"
  }
  - Ao receber o JSON este é persistido no banco de dados.

OBS: É possível que aja algum erro de interpretação meu na
  alocação dos servidores e portas. Para isto e qualquer 
  outra questão, estou à disposição.

