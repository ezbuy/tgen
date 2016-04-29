namespace go github.com..ezbuy..tgen..thriftgotest..enum

enum enumTest1 {
    EnumA = 5,
    EnumB,
    EnumC,
    EnumD = 20
}

enum enumTest2 {
    EnumA = 3,
    EnumB = 5,
    EnumC = 7,
    EnumD = 9
}

struct EnumFieldTest {
    1:required enumTest1 field1
    2:optional enumTest2 field2
}

service EnumServiceTest {
    enumTest1 Test1(1:enumTest1 arg)
    enumTest2 Test2(1:enumTest2 arg)
}
