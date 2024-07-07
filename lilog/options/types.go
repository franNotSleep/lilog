package options

const (
  ServerNameType uint8 = iota + 1
)

type ServerName string

type LigLogOptions struct {
  ServerName *ServerName
}
