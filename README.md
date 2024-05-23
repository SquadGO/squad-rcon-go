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
  "github.com/iamalone98/squad-rcon-go"
)

func main() {
  r, err := Dial("ip", "port", "password")
  if err != nil {
    return
  }

  defer r.Close()

  // Displays player messages, team kills, bans, kicks, squad creation and warns
  r.OnData(func(data string) {
    fmt.Println(data)
  })

  playersData := r.Execute("ListPlayers")
  fmt.Println(playersData)

  squadsData := r.Execute("ListSquads")
  fmt.Println(squadsData)

  // Use to prevent the program from ending
  for {

  }
}
```
