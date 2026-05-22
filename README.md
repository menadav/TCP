===================================================================================================================================================
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
FASE 1: Conectar server 
	Lograr que el Servidor acepte un cliente, le mande un saludo, y que el Cliente (CLI) lo reciba y lo pinte en pantalla usando tu paquete network.
FASE 2: Recibir llamadas cli -> server
	Lograr que el Cliente envie mensaje al server y le devuelva una respuesta simulada

===================================================================================================================================================

package main -> Obligatorio para poder compilar.
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

cmd/server/main.go:
    -:8080 ->                                    Leo en todas las interfazes de esta maquina   0.0.0.0:8080
    -Creamos servidor                          "listen, err := net.Listen("tcp", ":8080")"
    -Aceptamos los usuarios                    "user, err := listen.Accept()"
    -Enviamos mensaje al user                  "user.Write([]byte("Hola, bienvenido a The Answer Protocol\n"))"
    -Mantiene bucle infinito el servidor,
        siempre open.
===================================================================================================================================================

cmd/client/main.go:
    -Hacemos llamado server			"conn, err := net.Dial("tcp", ":8080")"
    -leemos respusta en el server		"n", err := conn.Read(buffer)"
===================================================================================================================================================
