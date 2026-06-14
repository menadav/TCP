===================================================================================================================================================

-Explicacion Breve.
Go es un lenguaje que compila a diferencia de python. Tambien tiene un Garbage Collector que no hace falta liberar memoria de variables.
- Concurrencia:
Go Destaca por su concurrencia es trabajo con Threads como con packete "net" crea hilos internamete eficientes ocupado casi nada de memoria.

C utiliza hilos del sistema operativo
Go utiliza Goroutines
    Gortuines: son hilos gestionados por el propio entorno de ejecución (runtime)

         C  |  Go
Memory  1MB    2KB  (Pero puede ir creciendo)
Create  slow   fast

===================================================================================================================================================

Flujo  Programa:
start -> go run cmd/server/main.go
- [server]Escucho puerto.                               ->	listen, err := net.Listen("tcp", ":8080")
- [server]Si cierra el programa cierra la net		    -> 	defer listen.Close()
- [server]Construimos un hub que contendra los canales  -> 	hub := constructor.NewHub()
- [server]La hacemos trabajar en segundo plano hilo     ->  go hub.run()
- [server]Entramos while infinito (for {})
- [server]Se freezea para esperar la llamada de un user ->	conn, err := listen.Accept()
go run cmd/client-cli/main.go
- [client]Llamamos al puerto del server		 	        ->	conn, err := net.Dial("tcp", ":8080")
- [client]Si cierra el programa cierra la net	 	    ->	defer conn.Close()
go run cmd/server/main.go
- [server]Registra al user al Hub de canales		    ->	hub.Register <- conn
- [server]Creamos Hilo que ejecutar lectura user        -> 	go network.ClientAtender(conn, hub)
- [server] -> [Cliente] comunicados
go run cmd/client-cli/main.go
- [client]Creare un segundo plano que leea server	    ->	go readServer
- [client]Podra escribir por terminal , que enviara [server] con su conexion de net.dial
Conseguimos una comunicacion Bidimensional [server]---[client]


===================================================================================================================================================


FASE 1: Conectar server 
	Lograr que el Servidor acepte un cliente, le mande un saludo, y que el Cliente (CLI) lo reciba y lo pinte en pantalla usando tu paquete network.

FASE 2: Recibir llamadas cli -> server y Concurrirlas con rutinas
	Lograr que el Cliente envie mensaje al server y le devuelva una respuesta simultanea con hilos

Fase 3: Recibir llamadas todos los clientes.
	Lograr que el cliente envie mensaje al server y devuelva a TODOS Broadcast .

Fase 4: Autentificador
    Lograr un bucle si no escribe correctamente los comandos para connectar

Fase 5: Crear Formato Chat.
    Depende lo que escriba el usuario hacer una accion o otra. Mirar, Escribir todos.

Fase 6: Crear Formato Grupo.
    Que jugador interacione con invitar, uniser dejar grupo.

Fase 7: Crear Parseo lineas jugador.
    Que jugador depende lo que escriba se parse.

Fase 8: Crear Struct World y constructor.

Fase 9: Unificar Data con Hub.

Fase 10: Crear Comandos Mundo.

Fase 11: Actualizar funciones para interacturar con el mundo
        Inventory, Look, Status, Who   Importante Checkear misiones Quests

Fase 12: Crear comandos de desplazamiento.

Fase 13: Crear comandos Take/Drop

===================================================================================================================================================

package main -> Obligatorio para poder compilar.
package: Si tienes mas de un script en la misma raiz , todos contendran el mismo name 
package <name>
Los files en la misma raiz se reconozen , no hace falta llamarlas de que package vienen.
- Diferente raiz:  network.function()
- Misma Raiz:	   function()


[go run <file>]
[go build <file>]

"fmt" -> Include prints (Parecido <stdio.h>)
"net" -> Include para trabajar con Net con Hilos
"bufio" -> Include para controlar \n  .Scanner

:=  -> DECLARAR variable, si queremos darle otro valor es con "="
defer -> Is finally like a python
<variable> <error> := <class>.<methode>("tcp", ":8080") -> segunda variable es para enviar el error como except pero hay que checkearlo.
<function> |make| -> Se usa para 3 tipos:
 - Slices:  ej buffer := make([]byte, 0, 1024) | numeros := make([]int, 5) 
 - Maps/Dict: ej edades := make(map[string]int) | jugadores := make(map[string]net.Conn, 100) 
 - Channels: ej mensaje := make(chan string) | colaTareas := make(chan int, 10)
if err != nil {return} -> No hay try,except. Se comprueba los errores similar a C.

===================================================================================================================================================

go.mod
    Es un requeriments.txt 

===================================================================================================================================================

go.sum 

===================================================================================================================================================

2: enp4s0f0: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc mq state UP group default qlen 1000
    link/ether 18:7e:b9:08:7a:a4 brd ff:ff:ff:ff:ff:ff
    inet 10.11.14.6/16 metric 100 brd 10.11.255.255 scope global enp4s0f0
       valid_lft forever preferred_lft forever
    inet6 fe80::1a7e:b9ff:fe08:7aa4/64 scope link 
       valid_lft forever preferred_lft forever