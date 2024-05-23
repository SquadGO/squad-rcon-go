# SquadRcon

This library is designed for the game Squad, it will give you the ability to easily connect to Rcon and parse/execute commands. I hope it will be useful to you!

## Install

```console
$ go get github.com/iamalone98/squad-rcon-go
```

## Quick start example

```golang
import (
  "fmt"
  "github.com/iamalone98/squad-rcon-go"
)

func main() {
  rcon, err := Dial("ip", "port", "password")
  if err != nil {
    return
  }

  defer rcon.Close()

  // Displays player messages, team kills, bans, kicks, squad creation and warns
  rcon.OnData(func(data string) {
    fmt.Println(data)
  })

  playersData := rcon.Execute("ListPlayers")
  fmt.Println(playersData)

  squadsData := rcon.Execute("ListSquads")
  fmt.Println(squadsData)

  // Use to prevent the program from ending
  for {

  }
}
```
