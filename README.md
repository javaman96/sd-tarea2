# sd-tarea2
Sistema de almacenamiento Distribuido de libros para Biblioteca

## Integrantes

Sebastián Sanchez Lagos 201504022-2
Diego Córdova Opazo 201403009-6

## Instrucciones ejecución

Definir el funcionamiento general del sistema distribuido final

IMPORTANTE: Consideraciones importantes en la tarea


## Estructura del proyecto

```bash
.
├── chat
├── chunk_split.go
├── clientes
├── dataNodes
├── nameNode
├── README.md
└── recombine_chunks.go
```

OBS: Makefile hará uno de los siguientes comandos, pero
     dependiendo de la máquina en que nos encontremos.
     "go run programa máquina 1"
     "go run programa máquina 2"
     "go run programa máquina 3"
     "go run programa máquina 4"

## Funcionamiento a grandes rasgos

Clientes: Funcionamiento Clientes

NameNode: Funcionamiento NameNode

DataNodes: Funcionamiento DataNodes

## Coneccion con Red DI y Máquinas virtuales

IMPORTANTE: Consideraciones acerca del orden de ejecucion de las máquinas

* Para ejecutar cada maquina, hacer "cd sd_tarea2" para ir al git
* luego solo correr el comando "make"

+ Máquina 1: Que programa se ejecuta en esta máquina(?)
	+ ip:         10.10.28.121
	+ contraseña: DSmvWkQsyIkaJzU


+ Máquina 2: Que programa se ejecuta en esta máquina(?)
	+ ip:         10.10.28.123
	+ contraseña: kndFkwYEQRdcTTu


+ Máquina 3: Que programa se ejecuta en esta máquina(?)
	+ ip:         10.10.28.122
	+ contraseña: eDtGthpFSmaypHj


+ Máquina 4: Que programa se ejecuta en esta máquina(?)
	+ ip:         10.10.28.124
	+ contraseña: XvfDKuTBWbAXEgj


# Usar proto
## Revisar esto a medida que se vaya avanzando!!

se creo un go mod en "github.com/dcordova/sd_tarea_2/data_service"

si tocamos el codigo de chat/chat.go debemos usar el comando "go mod tidy"

si cambiamos las definiciones de las funciones de chat.go en chat.proto3
debemos usar "protoc --go_out=plugins=grpc:chat chat.proto"
