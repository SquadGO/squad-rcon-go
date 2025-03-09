# SquadRcon

This library is designed for the game Squad, it will give you the ability to easily connect to Rcon and parse/execute commands. I hope it will be useful to you!

## Install

```text
go get github.com/SquadGO/squad-rcon-go
```

## Quick start example

```golang
import (
  "fmt"
  rcon "github.com/SquadGO/squad-rcon-go"
)

func main() {
	r, err := rcon.NewRcon(rcon.RconConfig{host: "127.0.0.1", password: "123456", port: "27165", autoReconnect: true, autoReconnectDelay: 5})
	if err != nil {
		fmt.Println(err)
		return
	}

	defer r.Close()

	fmt.Println("[RCON] Connection successful")

  /* Listeners works after first initialization */

	r.emitter.On("connected", func(_ interface{}) {
		fmt.Println("[RCON] Connection successful")
	})

  r.emitter.On("close", func(_ interface{}) {
		fmt.Println("[RCON] Connection closed")
	})

	r.emitter.On("error", func(err interface{}) {
		fmt.Println(err)
	})

	r.emitter.On("data", func(data interface{}) {
		fmt.Println("Data: ", data)
	})

  r.emitter.On("CHAT_MESSAGE", func(data interface{}) {
		if v, ok := data.(rcon.Message); ok {
			fmt.Println("Message: ", v.Message)
		}
	})

	r.emitter.On("ListPlayers", func(data interface{}) {
		if v, ok := data.(rcon.Players); ok {
			fmt.Println("Players: ", v)
		}
	})

	r.Execute("ListPlayers")

  // Use to prevent the program from ending
  select {}
}
```

## Rcon Events

| Function                     | Callback param type |
| ---------------------------- | ------------------- |
| **connected**                | **nil**             |
| **close**                    | **nil**             |
| **error**                    | **Error**           |
| **data**                     | **String**          |
| **PLAYER_WARNED**            | **Warn**            |
| **PLAYER_KICKED**            | **Kick**            |
| **CHAT_MESSAGE**             | **Message**         |
| **POSSESSED_ADMIN_CAMERA**   | **PosAdminCam**     |
| **UNPOSSESSED_ADMIN_CAMERA** | **UnposAdminCam**   |
| **SQUAD_CREATED**            | **SquadCreated**    |
| **ListPlayers**              | **Players**         |
| **ListSquads**               | **Squads**          |
