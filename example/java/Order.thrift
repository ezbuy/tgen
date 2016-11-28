namespace * Order
namespace java com.daigou.sg.rpc.order
namespace csharp Zen.DataAccess.Order
namespace javascript TRPC

include "Payment.thrift"

struct TOrderSummary {
	1:required i32 pendingReplyCount;		//等待回复订单数
	2:required i32 pendingPaymentCount;		//等待付款订单数
	3:required i32 pendingToPurchaseCount;	//等待购买订单数
	4:required i32 purchasingCount;			//正在购买订单数
	5:required i32 purchasedCount;			//已经购买订单数
	6:required i32 arrivedInShanghaiCount;	//达到上海仓库订单数
	7:required i32 arrivedInGuangzhouCount;	//到达广州仓库订单数
	8:required i32 orderInParcelCount;		//生成包裹订单数
	9:required i32 cancelledCount;			//取消订单数
	10:required i32 completedCount;			//完成订单数
	11:required i32 arrivedInTaiwanCount;	//到达台湾仓库订单数
	12:required i32 arrivedInUSACount;		//到达美国仓库订单数
}

struct TOrder{
	1:required i32 id;					//订单id
	2:required string orderNumber;		//订单号
	3:required string productImage;		//商品图片
	4:required string productName;		//商品名称
	5:optional bool gstFee;				//GST费
	6:required bool insured;			//保险
	7:optional string sellerDiscount;	//合作卖家国际运费折扣
	8:required i32 packageId;
}

struct TOrderDetail {
	1:required i32 id;							//订单id
	2:required string orderNumber;				//订单号
	3:required string orderStatus;				//订单状态
	4:required string productName;				//商品名称
	5:required string productImage;				//商品图片
	6:required string unitPrice;				//单价
	7:required double localUnitPrice;			//新币单价
	8:required i32 qty;							//商品数量
	9:required bool insured;					//保险
	10:required i32 shipmentTypeId;				//运输方式id	
	11:required string altShipmentTypeName;		//运输方式英文名
	12:required string shipmentTypeCode;		//运输方式编号
	13:optional string productRemark;			//商品备注
	14:required bool canEdit;					//订单是否可编辑
	15:required string warehouseCode;			//仓库
	16:optional string sellerDiscount;			//合作卖家国际运费折扣
	17:required list<Payment.TPaymentBillCategory> orderBillDetails;	//订单费用
	18:optional list<TOrderRemark> orderRemarks;			//订单备注
	19:optional list<TOrderItem> orderItems;				//订单sku
	20:required double localInternalShipmentFee;//国内运费
	21:required double totalPrice;				//总金额
	22:required string originCode;				//采购国家
	23:required string productUrl;				//商品url
}

struct TOrderBill {
	1:required i32 id;							//订单id
	2:required string orderNumber;				//订单号
	3:required string productName;				//商品名称
	4:required string productImage;				//商品图片
	5:required string unitPrice;				//单价
	6:required string localUnitPrice;			//新币单价
	7:required i32 qty;							//商品数量
	8:required bool insured;					//保险
	9:required list<Payment.TPaymentBillCategory> orderBillDetails;	//订单费用
	10:required string total;					//总金额
}

struct TOrderRemark {
	1:required i32 id;				//订单备注id
	2:required string remark;		//内容
	3:required bool needReply;		//是否需要回复
	4:optional string attachments;	//附件
	5:required string createDate;	//创建日期
	6:required string creator;		//创建者			
	7:optional i32 offsetId;		//差额补齐id
}

struct TOrderItem {
	1:required string sku;			//sku名称
	2:required string skuPrice;		//sku价格
	3:required string localSkuPrice;	//sku local价格
	4:required i32 qty;				//sku数量
}

struct TArrivedOrderSummary {
	1:required string shipmentTypeCode;		//运输方式编号
	2:required string altShipmentTypeName;	//运输方式英文名
	3:required i32 arrivedCount;			//到达仓库订单数
	4:required i32 notArrivedCount;			//未到达仓库订单数
}

struct TReadyToShipSummary {
	1:required i32 arrivedCount;		//到达仓库订单数
	2:required i32 notArrivedCount;		//未到达仓库订单数
	3:required list<TOrder> orders;		//订单列表	
	4:required Payment.TPaymentBill paymentBill;	//账单信息
}

struct TEzShipping {
	1:required bool on;						//是否开启
	2:required list<string> description;	//描述	
}

service Order {

	/// <summary>
	/// 购物车商品生成订单
	/// </summary>
	/// <param name="authorizeForBalance">授权差额补齐上线金额</param>
	/// <param name="originCode">采购国家</param>
	/// <returns>付款号</returns>
	i32 PayForCheckOut(1:bool authorizeForBalance, 2:string originCode),

	/// <summary>
	/// 更新订单
	/// </summary>
	/// <param name="orderId">订单id</param>
	/// <param name="shipmentTypeId">运输方式id</param>
	/// <param name="warehouseCode">仓库</param>
	/// <param name="internalShipmentFee">国内运费</param>
	/// <param name="insured">保险</param>
	/// <param name="productRemark">商品备注</param>
	void UpdateOrder(1:i32 orderId, 2:i32 shipmentTypeId, 3:string warehouseCode, 4:double internalShipmentFee, 5:bool insured, 6:string productRemark),

	/// <summary>
	/// 取消订单
	/// </summary>
	/// <param name="orderId">订单id</param>
	void CancelOrder(1:i32 orderId),

	/// <summary>
	/// 获取订单统计
	/// </summary>
	/// <param name="originCode">采购国家</param>
	/// <returns>订单统计</returns>
	TOrderSummary GetOrderSummary(1:string originCode),

	/// <summary>
	/// 根据状态获取订单列表
	/// </summary>
	/// <param name="originCode">采购国家</param>
	/// <param name="orderStatus">订单状态</param>
	/// <returns>订单列表</returns>
	list<TOrder> GetOrderListByStatus(1:string originCode, 2:string orderStatus, 3:string warehouseCode),

	/// <summary>
	/// 获取订单明细
	/// </summary>
	/// <param name="orderId">订单id</param>
	/// <returns>订单明细</returns>
	TOrderDetail GetOrderDetail(1:i32 orderId),

	/// <summary>
	/// 是否开启订单全到仓库后自动发货
	/// </summary>
	/// <param name="ezShipping">是否开启</param>
	/// <returns>是否开启</returns>
	TEzShipping UserChangeEZShipping(1:bool ezShipping),

	/// <summary>
	/// 获取订到全到仓库后自动发货的状态
	/// </summary>
	/// <returns>EzShipping状态信息</returns>
	TEzShipping UserGetEZShippingStatus(),

	/// <summary>
	/// 获取到达仓库订单的统计
	/// </summary>
	/// <param name="originCode">采购国家</param>
	/// <param name="warehouseCode">仓库</param>
	/// <returns>到达仓库订单的统计</returns>
	list<TArrivedOrderSummary> GetArrivedOrderSummary(1:string originCode, 2:string warehouseCode),

	/// <summary>
	/// 获取要发货的订单信息
	/// </summary>
	/// <param name="originCode">采购国家</param>
	/// <param name="warehouseCode">仓库</param>
	/// <param name="shipmentTypeCode">运输方式</param>
	/// <returns>要发货的订单信息</returns>
	TReadyToShipSummary GetArrivedOrders(1:string originCode, 2:string warehouseCode, 3:string shipmentTypeCode),

	/// <summary>
	/// 订单发货生成包裹
	/// </summary>
	/// <param name="originCode">采购国家</param>
	/// <param name="warehouseCode">仓库</param>
	/// <param name="shipmentTypeCode">运输方式</param>
	/// <returns>包裹id</returns>
	i32 WrapOrderToPackage(1:string originCode, 2:string warehouseCode, 3:string shipmentTypeCode),

	/// <summary>
	/// 回复订单备注
	/// </summary>
	/// <param name="orderId">订单id</param>
	/// <param name="orderRemarkParentId">要恢复的remarkid</param>
	/// <param name="remark">内容</param>
	/// <param name="pictures">图片</param>
	/// <returns>订单备注id</returns>
	void ReplyOrderRemark(1:i32 orderId, 2:i32 orderRemarkParentId, 3:string remark, 4:string pictures),

	/// <summary>
	/// 根据订单Id加入购物车 implements by c#
	/// </summary>
	void UserAddToCartByOrderId(1:i32 orderId)
}
