namespace go github.com..ezbuy..tgen..thriftgotest..simpleArguments
namespace webapi api.SimpleArguments.sub1.sub2.sub3

struct StructArg {

}

service SimpleArguments {
    void BoolArg(1: bool arg)
    void ByteArg(1: byte arg)
    void I16Arg(1: i16 arg)
    void I32Arg(1: i32 arg)
    void I64Arg(1: i64 arg)
    void DoubleArg(1: double arg)
    void BinaryArg(1: binary arg)
    void StringArg(1: string arg)
    void ListArg(1: list<string> arg)
    void MapArg(1: map<string, string> arg)
    void StructArg(1: StructArg arg)
}
