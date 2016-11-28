namespace go github.com..ezbuy..tgen..thriftgotest..includeEnum

include "Enum.thrift";

struct IncludeEnumTest {
    1:required Enum.enumTest1 field1
    2:optional Enum.enumTest2 field2
    3:required Enum.EnumFieldTest field3
    4:optional Enum.EnumFieldTest field4
}

service IncludeEnumService {
    Enum.EnumFieldTest Test(1:Enum.EnumFieldTest arg)
}
