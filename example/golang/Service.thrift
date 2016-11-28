namespace go github.com..ezbuy..tgen..thriftgotest..service

include "Types.thrift";
include "Enum.thrift"

service TestService {
    Enum.enumTest2 ReturnEnum()
    bool ThrowException(1:i32 arg) throws (1:Types.ExceptionsTest1 failure)
    bool ThrowException2(1:i32 arg) throws (1:Types.ExceptionsTest1 failure1, 2:Types.ExceptionsTest2 failure2)
    oneway void OneWayRequest()
    oneway void OneWayRequestWithArg(1:i32 arg)
}

service TestService2 {
    bool SomeMethod()
    oneway void SomeMethod2()

    # broadcast
    oneway void OnSomeMethod()
    oneway void OnSomeMethod2()

    bool DirectSomeMethod()
    oneway void DirectOnSomeMethod2()
}

service TestServiceValid {
    bool IsValid(1:Types.NestedArg arg1, 2:Types.NestedArg arg2)
}
