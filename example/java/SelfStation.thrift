namespace java com.daigou.selfstation.rpc.selfstation
namespace csharp SelfStation.Website.Controllers.Interfaces

struct TLoginResult {
	1: bool isSuccessful;
	2: string token;
	3: list<string> StationNames;
}

struct TParcelSection {
	1: required string name;
	2: required string value;
}

struct TParcel {
	1:required string userName;
	2:required string phone;
	3:required string parcelNumber;
	4:required string status;
	2:required list<TParcelSection> sections;
}

service SelfStation {
	// 只有登陆成功isSuccessful为true时，token跟StationNames才可能会有值
	TLoginResult Login(1: string username, 2: string password),

	// 获取上次图片去七牛的token
	string UserGetUploadToken(),
	string UserSetParcelReceived(1: string parcelNumber, 2: int rating, 3: string signatureImageKey),

	// date必须是yyyy-MM-dd 的格式
	list<TParcel> UserFindParcel(1: string stationName, 2: string date, 3: string userName, 3: string phone, 4: string parcelNumber, 5: int offset, 6: int limit),

	// status: incoming / arrived / completed
	list<TParcel> UserListParcel(1: string stationName, 2: string status, 5: int offset, 6: int limit),

	TParcel UserGetParcel(1: string parcelNumber),

	// 获得station有哪些货架
	list<string> UserGetShelfNumbers(1: string stationName),

	// 如果成功的话，返回空字符串；如果返回的字符串不是空，则说明有错误，直接alert给用户看
	string UserPutParcelToShelf(1: string shelfNumber, 2: string parcelNumber),

	// 修改密码,如果成功返回空字符串，否则直接显示错误信息
	string UserModifyPassword(1:string currentPassword, 2:string newPassword),
}
