https://www.jianshu.com/p/73c9ed3a4877



# 1、什么是RPC

RPC 代指远程过程调用（Remote Procedure Call），它的调用包含了传输协议和编码（对象序列号）协议等等。允许运行于一台计算机的程序调用另一台计算机的子程序，而开发人员无需额外地为这个交互作用编程



# 2、Protobuf



Protocol Buffers 是一种与语言、平台无关，可扩展的序列化结构化数据的方法，常用于通信协议，数据存储等等。相较于 JSON、XML，它更小、更快、更简单

## 1、语法

```go
syntax ="proto3";
service SearchService{
    rpc Search(SearchRequest) returns (SearchResponse);
}
message SearchRequest{
string query =1;
  int32 page_number =2;
  int32 result_per_page =3;
}
message SearchResponse{
...
}
```

1. 文件的第一行指定您正在使用`proto3`语法：如果不这样做，则协议缓冲区编译器将假定您正在使用[proto2](https://developers.google.com/protocol-buffers/docs/proto)。这必须是文件的第一行非空，非注释行。

2. 分配的字段编号用于标识消息的二进制格式，分配了最好不要修改，如果修改要及时通知使用者
3. 范围1-15的字段编号需要一个字节来编码，包括字段编码和字段类型，16-2047的字段编号需要占用两个字节
4. 最小是1 最大2^29 -1 ，保留字段不能使用 19000-19999

### <font color=red size=5x>分配字段编号</font>

1. 单数，proto3的默认规则
2. repeated 可以重复任意次



### <font color=red size=5x>注释</font>

支持单行注释//

多行注释/**/



### <font color=red size=5x>保留字段reserved</font>

当定义好字段后, 在后续开发中发现某个字段根本没用.

例如 `string userName = 2;` 字段, 这个时候最好不要进行注释或删除.

有可能以后加载相同的旧版本, 这可能会导致数据损坏, 隐私错误等. 确保不会发生这种情况的一种方法是指定要删除的字段为保留字段.

```
message SubscribeReq {
  
  reserved 2;
  
  int32 subReqID = 1;
  string userName = 2;
  string productName = 3;
  string address = 4;
}

或者
message Foo {
  reserved 2, 15, 9 to 11;
  reserved "foo", "bar";
}
```

顾名思义, 就是此字段会被保留可能在以后会使用此字段. 使用关键字 `reserved` 表示要保留字段编号为 `2`.



### <font color=red size=5x>基本类型对照表</font>

| .proto Type | Notes | C++ Type | Java/Kotlin Type[1] | Python Type[3] | Go Type | Ruby Type                      | C# Type    | PHP Type          | Dart Type |
| :---------- | :---- | :------- | :------------------ | :------------- | ------- | :----------------------------- | :--------- | :---------------- | :-------- |
| double      |       | double   | double              | float          | float64 | Float                          | double     | float             | double    |
| float       |       | float    | float               | float          | float32 | Float                          | float      | float             | double    |
| int32       |       | int32    | int                 | int            | int32   | Fixnum or Bignum (as required) | int        | integer           | int       |
| int64       |       | int64    | long                | int/long[4]    | int64   | Bignum                         | long       | integer/string[6] | Int64     |
| uint32      |       | uint32   | int[2]              | int/long[4]    | uint32  | Fixnum or Bignum (as required) | uint       | integer           | int       |
| uint64      |       | uint64   | long[2]             | int/long[4]    | uint64  | Bignum                         | ulong      | integer/string[6] | Int64     |
| sint32      |       | int32    | int                 | int            | int32   | Fixnum or Bignum (as required) | int        | integer           | int       |
| sint64      |       | int64    | long                | int/long[4]    | int64   | Bignum                         | long       | integer/string[6] | Int64     |
| fixed32     |       |          | int[2]              | int/long[4]    | uint32  | Fixnum or Bignum (as required) | uint       | integer           | int       |
| fixed64     |       | uint64   | long[2]             | int/long[4]    | uint64  | Bignum                         | ulong      | integer/string[6] | Int64     |
| sfixed32    |       | int32    | int                 | int            | int32   | Fixnum or Bignum (as required) | int        | integer           | int       |
| sfixed64    |       | int64    | long                | int/long[4]    | int64   | Bignum                         | long       | integer/string[6] | Int64     |
| bool        |       | bool     | boolean             | bool           | bool    | TrueClass/FalseClass           | bool       | boolean           | bool      |
| string      |       |          | String              | str/unicode[5] | string  | String (UTF-8)                 | string     | string            | String    |
| bytes       |       | string   | ByteString          | str            | []byte  | String (ASCII-8BIT)            | ByteString | string            | List      |



<font color=red size=5x>自定义类型</font>

```go
message SearchResponse {
  repeated Result results = 1;
}

message Result {
  string url = 1;
  string title = 2;
  repeated string snippets = 3;
}
```

Result 为自定义的类型

<font color=red size=5x>嵌套类型</font>

```go
message SearchResponse {
  message Result {
    string url = 1;
    string title = 2;
    repeated string snippets = 3;
  }
  repeated Result results = 1;
}
```

在其他地方调用

```go
message SomeOtherMessage {
  SearchResponse.Result result = 1;
}
```

深度嵌套

```go
message Outer {                  // Level 0
  message MiddleAA {  // Level 1
    message Inner {   // Level 2
      int64 ival = 1;
      bool  booly = 2;
    }
  }
  message MiddleBB {  // Level 1
    message Inner {   // Level 2
      int32 ival = 1;
      bool  booly = 2;
    }
  }
}
```





<font color=red size=5x>Any类型</font>

该`Any`消息类型，可以使用邮件作为嵌入式类型，而不必自己.proto定义。一个`Any`含有任意的序列化消息`bytes`，以充当一个全局唯一标识符和解析为消息的类型的URL一起。要使用该`Any`类型，您需要[导入](https://developers.google.com/protocol-buffers/docs/proto3#other) `google/protobuf/any.proto`

```go
import "google/protobuf/any.proto";

message ErrorStatus {
  string message = 1;
  repeated google.protobuf.Any details = 2;
}
```

给定消息类型的默认类型URL是`type.googleapis.com/_packagename_._messagename_`。























































