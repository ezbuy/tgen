namespace java com.ezbuy.basic
namespace swift Basic

struct TBasic {}

service BasicService {
  TBasic getStruct(1: i32 key, 2: i64 id),
  string getServiceName(1: i64 id),
  string getName()
}
