
# minibus-go

Minibus-go is a Golang client library for communicating via the [ttyzero/minibus](https://github.com/ttyzero/minibus) 
message bus. You can use this to both listen for and send messages to minibus channels.


## Creating a client

Creating a client is very simple, use Client and pass in the fully qualified path
to the minibus working dir, or `minibus.Default` to connect to the default location
provided by `os.UserCacheDir() + /minibus`.  This is a non-blocking operation. 

```go
mb := minibus.Client(minibus.Default)
```

## Sending Messages

Sending messages is very straight forward, call Send with a channel and a message (string).

```go

err := mb.Send("channel", "This is my message")
if err != nil {
  fmt.Println("Failed to send", err)
}
```



## Listening for messages

Listen by opening a channel, this creates a background goroutine that will connect 
to the minibus service and begin outputting messages on the returned `chan string`

```go

exampleChan := mb.OpenChannel("example-channel")

for {
  select {
    msg, open := <- exampleChan:
    if !open {
      break 
    }
    fmt.Printf("(example-chan): %s", msg)
  }
}
```

## closing a channel

To close a channel, simply close() the `chan string` returned by OpenChannel, this
will cause the background goroutine to terminate. 

```go
mb := minibus.Client("default")
exampleChan := mb.OpenChannel("example-channel")
close(exampleChan)
```

<hr/>
<table border='none' width='100%'>
<tr><td>
<img src='https://raw.githubusercontent.com/ttyzero/logo/master/assets/ttyzero_animated.png' alt='ttyZero Logo' title='ttyZero Logo'/>
</td>
<td>
<h3>Minibus-go is part of the <a href='http://github.com/ttyzero'>ttyZero Project</a></h3>
<hr/>
<b>Minibus-go</b> is <i>(c) 2019 ttyZero authors</i> <br/>
 and is available under the <b>MIT license</b>. 
</td></tr>
</table>

