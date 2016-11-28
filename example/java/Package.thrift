namespace * Package
namespace java com.daigou.sg.rpc.tpackage
namespace csharp Zen.DataAccess.Package
namespace javascript TRPC
namespace swift TR

include "Payment.thrift"
include "Order.thrift"

struct TPackageSummary {
	1:required i32 pendingReplyCount;		//等待回复包裹数
	2:required i32 prepareShipmentCount;	//等待发货包裹数
	3:required i32 shippedCount;			//起运中包裹数
	4:required i32 arrangeDeliveryCount;	//等待安排取货包裹数
	5:required i32 pendingPaymentCount;		//等待付款包裹数
	6:required i32 pendingDeliveryCount;	//等待派送包裹数
	7:required i32 acknowledgeCount;		//等待评论包裹数
	8:required i32 completedCount;			//完成包裹数
}

struct TPackage {
	1:required i32 id;						//包裹id
	2:required string purchaseType;			//采购方式
	3:required string packageNumber;		//包裹号
	4:required string packageEtaDate;		//包裹到达预期日期
	5:required string chargeWeight;			//计费重
	6:required string shipmentTypeCode;		//运输方式编号
	7:required string altShipmentTypeName;	//运输方式名
	8:required string warehouseCode;		//仓库
	9:optional string arrivedDate;			//到达日期
	10:optional string shippedDate;			//起运日期
	11:optional string collectionDate;		//揽收日期
	12:required double totalFee;			//总费用
	13:required double packageWeight;		//包裹重量
	14:required string createDate;			//包裹生成日期
}

struct TPackageDetail {
	1:required TPackage tPackage;						//包裹信息
	2:required Payment.TPaymentBill paymentBill;		//账单信息
	3:required list<Order.TOrder> orders;				//订单列表
	4:optional list<Order.TOrderRemark> orderRemarks;	//备注
}

struct TPendingDeliveryPackage {
	1:required i32 shipmentId;					//派送方式id
	2:required string localDeliveryMethod;		//派送方式
	3:required string stationName;				//地铁站名
	4:required string periodName;				//派送时间段
	5:required string localDeliveryDate;		//派送日期
	6:required string shipToAddress;			//派送地址
	7:required bool canBeEdited;				//是否可修改
	8:optional i32 mrtStationItemId;			//地铁站id
	9:optional i32 neighbourhoodStationItemId;	//邻里点id
	10:optional i32 selfStationId;				//上门自取点id
	11:required double totalFee;				//派送总费用
	12:required list<TPackage> packages;			//包裹列表
}

struct TCompletedPackage {
	1:required i32 customerCommentId;	//评论id	
	2:required string satisfaction;		//满意度
	3:required string caption;			//标题
	4:required string comment;			//内容
	5:required list<TPackage> packages;	//包裹列表
}

struct TProductComment {
	1:required i32 id;				//商品评论id
	2:required i32 agentProductId;	//代购商品id
	3:required string originCode;	//采购国家
	4:required string attachments;	//图片
	5:required string comment;		//内容
	6:required i32 rating;			//满意度
	7:required string productUrl;	//商品url
	8:required string productImage;	//商品图片
	9:required string productSku;	//sku名称
	10:required string productName;	//商品名
	11:required i32 orderId;		//订单id
	12:required double subTotal;	//总价	
}

struct TInvoice {
	1:required TPackage package;					//包裹信息
	2:required Payment.TPaymentBill paymentBill;	//账单信息
}

struct TArrangeDeliveryBill {
	1:required double totalFee;						//总费用
	2:required list<TPackageDetail> packageBills;	//包裹账单列表
	3:required bool couponUsed;						//折扣券是否可使用
	4:required string couponErrorMessage;			//折扣券错误信息
}

struct TArrangeDeliveryPackage {
	1:required i32 arrivedCount;		//可安排包裹数
	2:required i32 notArrivedCount;		//未能安排包裹数
	3:required list<TPackage> packages;	//包裹列表
}

service Package {

	/// <summary>
	/// 获取包裹统计
	/// </summary>
	/// <returns>包裹统计</returns>
	TPackageSummary GetPackageSummary(),

	/// <summary>
	/// 获取包裹列表
	/// </summary>
	/// <param name="status">状态</param>
	/// <returns>包裹列表</returns>
	list<TPackage> GetPackageListByStatus(1:string status),

	/// <summary>
	/// 获取包裹明细
	/// </summary>
	/// <param name="packageId">包裹id</param>
	/// <returns>包裹明细</returns>
	TPackageDetail GetPackageDetail(1:i32 packageId),

	/// <summary>
	/// 获取等待派送包裹列表
	/// </summary>
	/// <returns>等待派送包裹列表</returns>
	list<TPendingDeliveryPackage> GetPendingDeliveryPackages(),

	/// <summary>
	/// 获取等待评论包裹列表
	/// </summary>
	/// <returns>等待评论包裹列</returns>
	list<TPackage> GetAcknowledgePackages(),

	/// <summary>
	/// 获取完成的包裹列表
	/// </summary>
	/// <returns>完成的包裹列表</returns>
	list<TCompletedPackage> GetCompletedPackages(),

	/// <summary>
	/// 获取待评论商品列表
	/// </summary>
	/// <param name="packageId">包裹id</param>
	/// <returns>待评论商品列表</returns>
	list<TProductComment> GetProductsByPackageId(1:i32 packageId),

	/// <summary>
	/// 获取发票
	/// </summary>
	/// <param name="packageId">包裹id</param>
	/// <returns>发票</returns>
	TInvoice GetInvoiceByPackageId(1:i32 packageId),

	/// <summary>
	/// 评论商品
	/// </summary>
	/// <param name="productComment">商品评论信息</param>
	/// <returns>评论id</returns>
	i32 CommentAgentProducts(1:TProductComment productComment);

	/// <summary>
	/// <summary>
	/// 获取可安排取货包裹数
	/// </summary>
	/// <returns>可安排取货包裹数</returns>
	TArrangeDeliveryPackage UserFindForArrangeDelivery(),

	/// <summary>
	/// 评论包裹
	/// </summary>
	/// <param name="packageIds">包裹id列表</param>
	/// <param name="subject">主题</param>
	/// <param name="content">内容</param>
	/// <param name="level">满意度</param>
	void SaveAcknowledge(1:string packageIds, 2:string subject, 3:string content, 4:string level),


	/// <summary>
	/// 回复包裹备注
	/// </summary>
	/// <param name="packageId">包裹id</param>
	/// <param name="orderRemarkId">orderremarkid</param>
	/// <param name="remark">内容</param>
	/// <param name="pictures">图片</param>
	void ReplyPackageRemark(1:i32 packageId, 2:i32 orderRemarkId, 3:string remark, 4:string pictures)
}
