namespace go github.com..ezbuy..tgen..thriftgotest..unusedInclude

include "Const.thrift"
include "SimpleArguments.thrift"
include "IncludeEnum.thrift"
include "Types.thrift"
include "Service.thrift"

service UnusedInclude {
    void Ping()
    i32 Count()
}

service UnusedInclude2 {
    void Ping2()
    i32 Count2()
}
