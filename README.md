# JustaFlu ™

Ransomware escrito para Segurança de Redes 2019/20, que visa infetar
máquinas Windows com CPUs Intel, tanto em ambientes 32 como 64 bits.

# TODO

* Bundle crypto service in fake zoom installer exec.
* Check if victim has paid the ransom
    + https://ethplorer.io/address/0x2c84c8b62355c56c65f7ced493585311d23c649b
    + https://github.com/EverexIO/Ethplorer/wiki/Ethplorer-API#get-last-block
    + https://api.ethplorer.io/getAddressInfo/0xEF959C66dA176123244D3A737d8403251E3ea504?apiKey=freekey
* Maybe use SCHTASKS.EXE command to schedule encryption of system
* Maybe limit number of key generate requests on the oracle,
  so the server isn't overloaded with requests

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
* Descompacta Tor (32-bits) do próprio executável do instalador falso
  do Zoom.
* Liga-se ao hidden service do atacante, num link especializado a gerar novos pares de
  chaves, e disponibiliza chaves à vítima (como oráculo).
    + Atacante possui um par de chaves mestra RSA.
    + Quando chega um novo pedido, gera uma chave AES, e um par de chaves RSA.
    + Cifra a chave privada gerada com a chave AES gerada, e cifra a chave AES
      com a sua chave privada mestra.
    + Armazena a chave AES e a chave pública localmente.
        - Em alternativa a armazenar a chave pública na íntegra,
          pode simplesmente armazenar um hash da chave.
    + Envia a chave pública gerada mais a chave privada gerada, cifrada com AES,
      à vítima.
* Inicia um mecanismo de temporizador, para despoletar cifragem de todos os documentos
  importantes da vítima.
* Quando chega à data prevista do ataque, cifra o sistema, e lança uma janela nova,
  persistente, a informar a vítima do ataque, junto com um endereço bitcoin para
  pagar o resgate dos dados.
* A cifragem dos ficheiros decorre da seguinte forma:
    + A vítima gera uma nova chave AES por cada ficheiro a cifrar, e cifra o ficheiro
      com essa nova chave
    + A chave AES gerada é cifrada com a chave pública cedida à vítima pelo atacante.
* A decifragem dos ficheiros decorre da seguinte forma:
    + Depois da vítima pagar o resgate, o atacante decifra a chave AES que armazenou
      localmente relativa à chave pública que a vítima deu como argumento.
    + Envia a chave AES à vítima.
    + A vítima usa a chave AES para decifrar a chave privada que tem armazenada localmente.
    + Para cada ficheiro cifrado da vítima:
        - É decifrada a chave AES relativa ao ficheiro com a chave privada, agora
          decifrada.
        - É decifrado o conteúdo do ficheiro com a chave AES relativa ao ficheiro.
