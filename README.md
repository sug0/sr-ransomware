# sr-ransomware

Ransomware para Segurança de Redes

# Passos do algoritmo

* Quando a vítima inicia o instalador falso do Zoom, começa o instalador
  verdadeiro do Zoom no plano de fundo.
* Entretanto, aloja-se nos serviços do Windows, ou programas de startup,
  para correr sempre que a vítima liga o PC.
* Usa um mecanismo para determinar se a vítima já foi infetada, como:
    + Criar chaves no registry do Windows relativas ao ransomware.
        - Se existir a chave no registry, não infetar PC de novo.
    + Criar um ficheiro JSON ou outra forma de serialização de dados.
        - Se existir ficheiro, não infetar de novo.
    + Para além de ajudarem a perceber se um computador já foi infetado,
      registariam metadados como o tempo de infeção, mas útil à frente.
* Faz download do TOR (32-bits)
* Liga-se ao hidden service do atacante, num link especializado a gerar novos pares de
  chaves, e disponibiliza a chave pública à vítima.
    + Do lado do atacante, todas as chaves privadas das vítimas seriam cifradas
      com uma chave pública mestra
* A vítima usa a chave pública gerada pelo atacante para gerar as suas chaves locais.
    + Uma chave pública RSA (guardada em plaintext).
    + Uma chave simétrica AES-128 mais o nonce (cifrada com a chave pública gerada
      pelo atacante).
    + Uma chave privada RSA (cifrada com a chave AES gerada pela vítima).
* Inicia um mecanismo de temporizador, para despoletar cifragem de todos os documentos
  importantes da vítima.
* Quando chega à data prevista do ataque, cifra o sistema, e lança uma janela nova,
  persistente, a informar a vítima do ataque, junto com um endereço bitcoin para
  pagar o resgate dos dados.
