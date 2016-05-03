namespace go github.com..ezbuy..tgen..thriftgotest..types

exception ExceptionsTest1 {
    1:required i32 code
    2:optional string message
}

exception ExceptionsTest2 {
    1:required bool success
    2:optional string errorMsg
}
