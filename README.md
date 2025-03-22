# SquadRcon

This library is designed for the game Squad, it will give you the ability to easily connect to Rcon and parse/execute commands. I hope it will be useful to you!

## Install

```text
go get -u github.com/SquadGO/squad-rcon-go/v2
```

## Quick start example

```golang
import (
  "fmt"
  rcon "github.com/SquadGO/squad-rcon-go/v2"
)

func main() {
  r, err := rcon.NewRcon(rcon.RconConfig{Host: "127.0.0.1", Password: "123456", Port: "27165", AutoReconnect: true, AutoReconnectDelay: 5})
  if err != nil {
    fmt.Println(err)
    return
  }

  defer r.Close()

  fmt.Println("[RCON] Connection successful")

  /* Listeners works after first initialization */

  r.Emitter.On("connected", func(_ interface{}) {
    fmt.Println("[RCON] Connection successful")
  })

  r.Emitter.On("close", func(_ interface{}) {
    fmt.Println("[RCON] Connection closed")
  })

  r.Emitter.On("error", func(err interface{}) {
    fmt.Println(err)
  })

  r.Emitter.On("data", func(data interface{}) {
    fmt.Println("Data: ", data)
  })

  r.Emitter.On("CHAT_MESSAGE", func(data interface{}) {
    if v, ok := data.(rcon.Message); ok {
      fmt.Println("Message: ", v.Message)
    }
  })

  r.Emitter.On("ListPlayers", func(data interface{}) {
    if v, ok := data.(rcon.Players); ok {
      fmt.Println("Players: ", v)
    }
  })

  r.Execute("ListPlayers")

  // Use to prevent the program from ending
  select {}
}
```

## Listeners

| Listener                     | Returns           |
| ---------------------------- | ----------------- |
| **connected**                | **nil**           |
| **close**                    | **nil**           |
| **error**                    | **Error**         |
| **data**                     | **String**        |
| **PLAYER_WARNED**            | **Warn**          |
| **PLAYER_KICKED**            | **Kick**          |
| **CHAT_MESSAGE**             | **Message**       |
| **POSSESSED_ADMIN_CAMERA**   | **PosAdminCam**   |
| **UNPOSSESSED_ADMIN_CAMERA** | **UnposAdminCam** |
| **SQUAD_CREATED**            | **SquadCreated**  |
| **ListPlayers**              | **Players**       |
| **ListSquads**               | **Squads**        |
