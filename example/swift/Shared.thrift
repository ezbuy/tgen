namespace java com.ezbuy.basic
namespace swift TRS

struct TBasic {}

service Basic {
  TBasic getStruct(1: i32 key, 2: i64 id),
  list<TBasic> getObjects(1: i32 key, 2: i64 id),
  list<i16> getAges(1: i64 id),
  list<string> getNames()
}
