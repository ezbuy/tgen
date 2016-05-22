namespace go github.com..ezbuy..tgen..thriftgotest..types

exception ExceptionsTest1 {
    1:required i32 code
    2:optional string message
}

exception ExceptionsTest2 {
    1:required bool success
    2:optional string errorMsg
}

struct InnerArg1 {
    1:required InnerArg2 in11
    2:optional InnerArg2 in12
}

struct InnerArg2 {
    1:required InnerArg3 in21 
    2:optional InnerArg3 in22
}

struct InnerArg3 {
    1:required string s1
    2:optional string s2
}

struct NestedArg {
    1:required InnerArg1 n1
    2:optional InnerArg1 n2
}
