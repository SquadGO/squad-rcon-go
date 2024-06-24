# SquadRcon

This library is designed for the game Squad, it will give you the ability to easily connect to Rcon and parse/execute commands. I hope it will be useful to you!

## Install

```text
go get github.com/iamalone98/squad-rcon-go
```

## Quick start example

```golang
import (
  "fmt"
  rcon "github.com/iamalone98/squad-rcon-go"
)

func main() {
  r, err := rcon.Dial("ip", "port", "password")
  if err != nil {
    return
  }

  defer r.Close()

  r.OnClose(func(err error) {
    fmt.Println(err)
  })

  // Displays player messages, team kills, bans, kicks, squad creation and warns
  r.OnData(func(data string) {
    fmt.Println(data)
  })

  rcon.OnWarn(func(data rcon.Warn) {
    fmt.Println("Warn: ", data)
  })

  data := r.Execute("ListPlayers")
  fmt.Println(data)

  rcon.OnListSquads(func(data rcon.Squads) {
    fmt.Println(data)
  })

  r.Execute("ListSquads")

  // Use to prevent the program from ending
  select {}
}
```

## Rcon Events

| Function            | Callback param type |
| ------------------- | ------------------- |
| **OnClose**         | **Error**           |
| **OnData**          | **String**          |
| **OnWarn**          | **Warn**            |
| **OnKick**          | **Kick**            |
| **OnMessage**       | **Message**         |
| **OnPosAdminCam**   | **PosAdminCam**     |
| **OnUnposAdminCam** | **UnposAdminCam**   |
| **OnSquadCreated**  | **SquadCreated**    |
| **OnListPlayers**   | **Players**         |
| **OnListSquads**    | **Squads**          |
