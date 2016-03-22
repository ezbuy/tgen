/// <summary>
/// implements by c#
/// </summary>

namespace * ShipForMe
namespace csharp Zen.DataAccess.ShipForMe
namespace javascript TRPC
namespace java com.daigou.sg.rpc.shipforme

include "Payment.thrift"

struct TShipForMeOrderHomeSummary {
	1:required i32 cnCount;			//中国订单数
	2:required i32 twCount;			//台湾订单数
	3:required i32 usCount;			//美国订单数
}

struct TShipForMeOrderSummary {
	1:required i32 notReceivedCount;			//未收到订单数
	2:required i32 problemOrdersCount;			//要回复的问题订单数
	3:required i32 notReplyProblemOrdersCount;	//不需要回复的问题订单数
	4:required i32 orderInParcelCount;			//生成包裹的订单数
	5:required i32 historyCount;				//历史订单数
	6:required i32 readyToShipSHCount;			//到达上海仓库订单数
	7:required i32 readyToShipGZCount;			//到达广州仓库订单数
	8:required i32 readyToShipTWCount;			//到达台湾仓库订单数
	9:required i32 readyToShipUSCount;			//到达美国仓库订单数
}

struct TShipForMeOrder {
	1:required i32 id;							//订单id
	2:required string orderNumber;				//订单号
	3:required string shipperName;				//快递公司
	4:required string wayBill;					//运单号
	5:required string alternative;				//备注
	6:required string warehouseCode;			//仓库
	7:required i32 arrivedDays;					//到达仓库的天数
	8:required double weight;					//重量
	9:required double volumeWeight;				//体积重
	10:required double unitPrice;				//申报价
	11:required bool hasPhotoService;			//是否有拍照服务
	12:required bool hasRepackService;			//是否有打包服务
	13:required bool hasOtherService;			//是否有其他服务
	14:required bool isProcessing;				//订单是否在处理
	15:optional list<TVendorName> vendorNames;	//可选运输方式
	16:optional string correctVendorName;		//原因
	17:optional double valueAddedCharge;		//增值服务费
	18:required string orderStatus;				//订单状态
	19:optional string attachments;				//附件
	20:optional string valueAddedService;		//增值服务
	21:optional string repackService;			//打包服务
	22:optional list<TOrderRemark> remarks;		//订单备注
}

struct TVendorName {
	1:required string shipmentTypeCode;			//运输方式编号	
	2:required string shipmentTypeName;			//运输方式名
}

struct TShipformeAddress {
	1:required string warehouse;			//仓库
	2:required string name;					//名字
	3:required string area;					//地区
	4:required string address;				//地址
	5:required string zipcode;				//邮编
	6:required string phoneNumber;			//联系电话
}

struct TShipformeOrderBill {
	1:required bool couponUsed;						//折扣券是否可使用
	2:required string couponErrorMessage;			//折扣券错误信息
	3:required Payment.TPaymentBill paymentBill;		//发货账单列表
	
}

service ShipfForMe {

	/// <summary>
	/// 获取各地区订单数
	/// </summary>
	/// <returns>各地区订单数</returns>
	TShipForMeOrderHomeSummary UserGetShipForMeHomeSummary(),

	/// <summary>
	/// 获取各状态订单数
	/// </summary>
	/// <param name="originCode">采购国家</param>
	/// <returns>各状态订单数</returns>
	TShipForMeOrderSummary UserGetShipForMeSummary(1:string originCode),

	/// <summary>
	/// 根据状态获取订单列表
	/// </summary>
	/// <param name="originCode">采购国家</param>
	/// <param name="warehouseCode">仓库</param>
	/// <param name="status">状态</param>
	/// <param name="offset">数据请求位置</param>
	/// <param name="limit">请求个数</param>
	/// <returns>订单列表</returns>
	list<TShipForMeOrder> UserGetShipForMeOrderListByStatus(1:string originCode, 2:string warehouseCode, 3:string status, 4:i32 offset,
	 5:i32 limit),

	/// <summary>
	/// 获取订单明细
	/// </summary>
	/// <param name="orderId">订单id</param>
	/// <returns>订单明细</returns>
	TShipForMeOrder UserGetShipForMeOrderDetailByOrderId(1:i32 orderId),

	/// <summary>
	/// 获取自助购地址
	/// </summary>
	/// <param name="originCode">采购国家</param>
	/// <returns>自助购地址</returns>
	list<TShipformeAddress> UserGetShipForMeAddressByRegion(1:string originCode),

	/// <summary>
	/// 获取自助购快递公司
	/// </summary>
	/// <param name="originCode">采购国家</param>
	/// <returns>自助购快递公司</returns>
	list<string> GetCourierCompaniesByRegion(1:string originCode),

	/// <summary>
	/// 新增自助购订单
	/// </summary>
	/// <param name="warehouseCode">仓库</param>
	/// <param name="shipperName">快递公司</param>
	/// <param name="wayBill">运单号</param>
	/// <param name="alternative">备注</param>
	/// <param name="takePhoto">是否拍照</param>
	/// <param name="originCode">采购国家</param>
	/// <param name="repack">是否打包</param>
	/// <returns>自助购订单id</returns>
	i32 UserAddNewShipForMeOrder(1:string warehouseCode, 2:string shipperName; 3:string wayBill, 4:string alternative,
		5:bool takePhoto, 6:string originCode, 7:string repack),

	/// <summary>
	/// 删除自助购订单
	/// </summary>
	/// <param name="orderId">订单id</param>
	void UserDeleteShipForMeOrder(1:i32 orderId),

	/// <summary>
	/// 获取包裹统计
	/// </summary>
	/// <param name="orderId">订单id</param>
	/// <param name="warehouseCode">仓库</param>
	/// <param name="shipperName">快递公司</param>
	/// <param name="wayBill">运单号</param>
	/// <param name="alternative">备注</param>
	/// <param name="takePhoto">是否拍照</param>
	/// <param name="repack">是否打包</param>
	void UserUpdateShipForMeOrder(1:i32 orderId, 2:string warehouseCode, 3:string shipperName; 4:string wayBill, 5:string alternative,
		6:bool takePhoto, 7:string repack),

	/// <summary>
	/// 发送短信验证码
	/// </summary>
	/// <param name="phoneNumber">手机号</param>
	/// <returns>是否发送成功</returns>
	bool UserSendToTelephone(1:string phoneNumber),

	/// <summary>
	/// 验证验证码
	/// </summary>
	/// <param name="phoneNumber">手机号</param>
	/// <param name="validationCode">验证码</param>
	/// <returns>是否验证成功</returns>
	bool UserValidationPhone(1:string phoneNumber, 2:string validationCode),

	/// <summary>
	/// 自助购订单是否打包
	/// </summary>
	/// <param name="confirm">是否打包</param>
	void UserConfirmShipForMeRepack(1:bool confirm, 2:i32 orderId),

	/// <summary>
	/// 设置订单申报价
	/// </summary>
	/// <param name="orderIds">订单id</param>
	/// <param name="price">申报价</param>
	void UserSaveShipForMeOrderPrice(1:string orderIds, 2:double price),
	
	/// <summary>
	/// 获取订单发货账单
	/// </summary>
	/// <param name="orderIds">订单id</param>
	/// <param name="insured">是否保险</param>
	/// <param name="deliveryMethod">派送方式</param>
	/// <param name="shipmentTypeCode">运输方式</param>
	/// <param name="customerAddressId">送货上门地址</param>
	/// <param name="originCode">采购国家</param>
	/// <param name="warehouseCode">仓库</param>
	/// <param name="couponCode">折扣码</param>
	/// <returns>订单发货账单</returns>
	TShipformeOrderBill UserGetShipForMeOrderFeeByOrderIds(1:list<string> orderIds, 2:bool insured, 3:string deliveryMethod, 4:string shipmentTypeCode,
		5:i32 customerAddressId, 6:string originCode, 7:string warehouseCode, 8:string couponCode)

	/// <summary>
	/// 订单发货生成包裹
	/// </summary>
	/// <param name="orderIds">订单id</param>
	/// <param name="insured">是否保险</param>
	/// <param name="deliveryMethod">派送方式</param>
	/// <param name="shipmentTypeCode">运输方式</param>
	/// <param name="customerAddressId">送货上门地址</param>
	/// <param name="originCode">采购国家</param>
	/// <param name="warehouseCode">仓库</param>
	/// <param name="couponCode">折扣码</param>
	/// <returns>包裹id</returns>
	i32 UserPackShipForMeOrder(1:list<string> orderIds, 2:bool insured, 3:string deliveryMethod, 4:string shipmentTypeCode,
		5:i32 customerAddressId, 6:string originCode, 7:string warehouseCode, 8:string couponCode)

}
