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
- RTU
- ASCII
- RTU_OVER_TCP

# 使用插件

- go.bug.st/serial v1.5.0

# 测试使用工具

- Modbus Poll
- Modbus Slave
- HHD Virtual Serial Port Tools

# 使用说明
- TCP
```go
st := NewTcpTransporter("127.0.0.1:502")
st.TimeOut=1*time.Second
pk := NewTcpPackager(1)
defer func() { _ = st.Close() }()
pk := NewTcpPackager(1)
c := NewClient(pk, st)
request, results, err := c.ReadHoldingRegisters(1, 10)
```
- RTU
```go
st := NewSerialTransporter("COM3")
st.Mode=serial.Mode{BaudRate: defaultBaudRate}
defer func() { _ = st.Close() }()
pk := NewRtuPackager(1)
c := NewClient(pk, st)
request, results, err := c.ReadHoldingRegisters(1, 10)
```
- ASCII
```go
st := NewSerialTransporter("COM3")
st.Mode=serial.Mode{BaudRate: defaultBaudRate}
defer func() { _ = st.Close() }()
pk := NewAsciiPackager(1)
c := NewClient(pk, st)
request, results, err := c.ReadHoldingRegisters(1, 10)
```
- RTU_OVER_TCP
```go
st := NewTcpTransporter("127.0.0.1:502")
st.TimeOut=1*time.Second
defer func() { _ = st.Close() }()
pk := NewRtuPackager(1)
c := NewClient(pk, st)
request, results, err := c.ReadHoldingRegisters(0, 10)
```

