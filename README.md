# its-a-key-value-database
This is an in-memory key value database similar to Redis (but a _bit_ more basic).
The program works by creating a tcp server on localhost:3332 that can then be connected to by a client. 
Commands are entered in the format **verb key value**

The accepted commands:
* SET key value
* GET key
* GETALL
* UPDATE key value
* DELETE key
* LEN
* GETVALUES
* GETKEYS
* EXISTS key

## Example inputs and outputs:
#### SET
```
> SET user1 Aden
> +OK
> SET user1 Eilers
> -ERR key already exists
```
#### GET
```
> GET user1
> +Aden
> GET user2
> -ERR key does not exist
```
#### GETALL
```
> GETALL
> +user1
> Aden
> user2
> Mike
> GETALL
> -ERR no entries in database
```
#### UPDATE
```
> UPDATE user1
> +OK
> UPDATE user2
> -ERR key does not exist
```
#### DELETE
```
> DELETE user1
> +OK
> DELETE user2
> -ERR key does not exist
```
#### LEN
```
> LEN
> +5
> LEN
> +0
```
#### GETVALUES
```
> GETVALUES
> +Aden
> Eilers
> Mike
> GETVALUES
> -ERR no entries in database
```
#### GETKEYS
```
> GETKEYS
> +user1
> user2
> user3
> GETKEYS
> -ERR no entries in database
```
#### EXISTS
```
> EXISTS user1
> +1
> EXISTS user2
> +0
```