# sd-tarea2
Sistema de almacenamiento distribuido de libros para Biblioteca

## Integrantes

Sebastián Sanchez Lagos 201504022-2
Diego Córdova Opazo 201403009-6

## Instrucciones ejecución

Definir el funcionamiento general del sistema distribuido final
IMPORTANTE: Consideraciones importantes en la tarea


## Estructura del proyecto

```bash
.
├── data_node
├────── data_node.go
├── data_service.proto
├── data_service
├────── data_service.go
├────── data_service.pb.go
├────── chunks
├────────── (chunks_almacenados)
├── downloader_client
├────── client.go
├── name_node
├────── name_node.go
├────── logs
├────────── logs.txt
├── name_service.proto
├── name_service
├────── name_service.go
├────── name_service.pb.go
├── uploader_client
├────── client.go
├── README.md
├── MAKEFILE
```

OBS: Makefile correra cada nodo dependiendo de en que maquina nos encontremos.



## Funcionamiento a grandes rasgos


Para implementar las funcionalidades de los dos tipos de servidores, name_node y data_node, se crearon dos servicios grpc, name_service y data_service. 

Para cargar libros, estos se separan en chunks, y son enviados siempre al datanode1 (data node de la primera maquina), quien hace la función de maestro dentro de los data nodes (Se nos había dicho que era una opción viable). Este datanode genera una propuesta, y dependiendo de si el sistema es distribuido o centralizado se envía al namenode o a los datanodes para ser validada. La propuesta contiene las IP’s de los data nodes a las que deben ir estos nodos, y el datanode1 los reparte.

El descargador de libros tiene las funcionalidades de pedir la lista de libros, y de pedir descargar un libro en específico. 

El namenode se encarga de manejar el archivo logs.txt. Esta puede sacar el listado de libros, sacar los chunks con sus ips de un libro en específico, y tambien escribir los chunks de un nuevo libro que se sube al sistema.

Los datanodes se encargan de escribir en memoria los chunks, de recibir los chunks a ser subidos de parte de los clientes, generan propuestas (listas de chunks+ip) y también se encargan de leer chunks solicitados y devolverlos al cliente de descargas.

## Detalles y Consideraciones

* Sos nombres de los pdf no puede tener espacios


## Coneccion con Red DI y Ejecución en Máquinas virtuales

IMPORTANTE: Consideraciones acerca del orden de ejecucion de las máquinas

* Para ejecutar cada maquina, hacer "cd sd_tarea2" para ir al git luego solo correr el comando "make" (Esto correra un datanode o el namenode dependiendo de la máquina)
* Para correr un cliente downloader o uploader se debe hacer desde la máquina 1, debido a las ip's de las conexiones
* Para conectarse con un uploader o downloader desde otra maquina habría que cambiar las ip's en el código
de los archivos client.go correspondientes




+ Máquina 1: máquina con el data_node 1
	+ ip:         10.10.28.121
	+ contraseña: DSmvWkQsyIkaJzU


+ Máquina 2: máquina con el data_node_2
	+ ip:         10.10.28.122
	+ contraseña: eDtGthpFSmaypHj	


+ Máquina 3: máquina con el data_node_3
	+ ip:         10.10.28.123
	+ contraseña: kndFkwYEQRdcTTu


+ Máquina 4: máquina con el name_node
	+ ip:         10.10.28.124
	+ contraseña: XvfDKuTBWbAXEgj


# Usar proto
## Revisar esto a medida que se vaya avanzando!!

se creo un go mod en "github.com/dcordova/sd_tarea_2/data_service" y en "github.com/dcordova/sd_tarea_2/name_service"

si tocamos el codigo de chat/chat.go debemos usar el comando "go mod tidy"

si cambiamos las definiciones de las funciones de chat.go en chat.proto3
debemos usar "protoc --go_out=plugins=grpc:name_service name_service.proto"
