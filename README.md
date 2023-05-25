# go modbus

用go实现modbus协议。

# 支持功能码

Bit 操作:

- 功能码：1 读线圈 
- 功能码：2 读离散量输入
- 功能码：5 写单个线圈
- 功能码：15 写多个线圈

16-bit 操作:

- 功能码：4 读输入寄存器
- 功能码：3 读保持寄存器
- 功能码：6 写单个寄存器
- 功能码：16 写多个寄存器
- 功能码：23 读/写多个寄存器


# 支持格式

- TCP
- Serial (RTU, ASCII)

# 使用插件

- go.bug.st/serial v1.5.0

# 测试使用工具

- Modbus Poll
- Modbus Slave
- HHD Virtual Serial Port Tools

# 使用说明

Basic usage:

```go
// Modbus TCP
client := modbus.TCPClient("localhost:502")
// Read input register 9
results, err := client.ReadInputRegisters(8, 1)

// Modbus RTU/ASCII
// Default configuration is 9600, 8, 1, even
st := NewSerialTransporter("COM3")
pk := modbus.NewRtuPackager(1)
defer st.Close()
client := modbus.NewClient(pk,st)
results, err := client.ReadHoldingRegisters(1, 2)
```

Advanced usage:

```go
// Modbus TCP
handler := modbus.NewTCPClientHandler("localhost:502")
handler.Timeout = 10 * time.Second
handler.SlaveID = 0xFF
handler.Logger = log.New(os.Stdout, "test: ", log.LstdFlags)
// Connect manually so that multiple requests are handled in one connection session
err := handler.Connect()
defer handler.Close()

client := modbus.NewClient(handler)
results, err := client.ReadDiscreteInputs(15, 2)
results, err = client.WriteMultipleRegisters(1, 2, []byte{0, 3, 0, 4})
results, err = client.WriteMultipleCoils(5, 10, []byte{4, 3})
```

```go
// Modbus RTU/ASCII
st := NewSerialTransporterMode("COM3",serial.Mode{
BaudRate: 9600,
})
pk := modbus.NewRtuPackager(1)
defer st.Close()
client := modbus.NewClient(pk,st)
results, err := client.ReadHoldingRegisters(1, 2)
```

