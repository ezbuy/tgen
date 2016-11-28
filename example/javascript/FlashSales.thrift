namespace go flashsales
namespace java com.daigou.sg.webapi.flashsales
namespace swift TR

struct TFlashSalesSetting {
	1:required string settingId;	//id
	2:required string salesName;	//闪购活动名称
	3:required string salesStartTime;	//闪购活动开始时间 格式 yyyy-MM-dd hh:mm:ss
	4:required string salesEndTime;		//闪购活动结束时间 格式 yyyy-MM-dd hh:mm:ss
	5:required bool	isActive;			//是否启用
	6:required i32	limitation;		//本次闪购活动单个商品的购买数量限制
	7:required bool	usableCoupon;	//是否可以使用Coupon
}

struct TFlashSalesProduct {
	1:required string productUrl;	//商品标准URL 格式 https://item.taobao.com/item.htm?id=xxxx
	2:required double productPrice;	//商品的闪购价格
	3:required i32	  productStock;	//设置商品的库存数量
	4:required double productOriginPrice; //商品原始价格
	5:required string productImage;	//商品图片
	6:required string catalogCode;	//闪购商品支持的国家
	7:required string discount;		//闪购折扣比例
	8:required string productLocalPrice; //商品的闪购本地价格
	9:required string productOriginLocalPrice; //闪购商品的原始价格
	10:required string rebateProductUrl;	//闪购商品的返利链接
	11:required string productName;			//闪购商品名称
	12:required string 	categoryName; // 闪购商品分类名
}

struct TFlashSalesCategoryProducts{
	1:required string categoryName;
	2:required list<TFlashSalesProduct> products; // 某一分类下的商品列表
	3:required i64 beginTimeSpan;					//当前时间距离闪购开始间隔毫秒数
	4:required i64 endTimeSpan;						//当前时间距离闪购结束间隔毫秒数
	5:required string settingId;					//闪购活动id
	6:required string salesName;					//闪购活动名称
	7:required string mobileSalesName;				//mobile salesName 闪购活动名称	

}

struct TFlashSalesSummary{
	1:required list<TFlashSalesProduct> products;	//闪购商品
	2:required i64 beginTimeSpan;					//当前时间距离闪购开始间隔毫秒数
	3:required i64 endTimeSpan;						//当前时间距离闪购结束间隔毫秒数
	4:required string settingId;					//闪购活动id
	5:required string salesName;					//闪购活动名称
	6:required string mobileSalesName;				//mobile salesName 闪购活动名称
}

service FlashSales {

	///<summary>
	/// 新增闪购活动
	///</summary>
	void UserAddFlashSalesSetting(1:TFlashSalesSetting setting),

	///<summary>
	/// 新增闪购商品
	///</summary>
	void UserAddFlashSalesProduct(1:string settingId, 2:TFlashSalesProduct product),

	///<summary>
	/// 删除闪购商品
	///</summary>
	void UserDeleteFlashSalesProduct(1:string settingId, 2:string productUrl),

	///<summary>
	/// 更新闪购活动
	///</summary>
	void UserUpdateFlashSalesSetting(1:TFlashSalesSetting setting),

	/// <summary>
	/// 获取闪购列表
	/// </summary>
	TFlashSalesSummary GetFlashSalesList(1:i32 offset, 2:i32 limit ,3:string area),

	/// <summary>
	/// 获取闪购活动列表
	/// </summary
	list<TFlashSalesSetting> UserGetFlashSalesSettingsListAdmin(1:i32 offset, 2:i32 limit),

	/// <summary>
	/// 获取闪购商品列表
	/// </summary>
	list<TFlashSalesProduct> UserGetFlashSalesListAdmin(1:string settingId, 2:i32 offset, 3:i32 limit),

	/// <summary>
	/// 启用闪购活动，并更新elastic search中的商品数据
	/// </summary>
	void SyncFlashSales(1:string settingId),

	/// <summary>
	/// 删除闪购活动
	/// </summary>
	void UserDeleteFlashSalesSetting(1:string settingId),

	///<summary>
	/// 更新闪购商品
	///</summary>
	void UserUpdateFlashSalesProduct(1:string settingId, 2:TFlashSalesProduct product),

	/// <summary>
	/// 获取一天的闪购活动
	/// date: 具体日期（2016-06-15）
	/// top: 前几个商品（0表示不取商品）
	/// </summary>
	list<TFlashSalesSummary> GetFlashSalesListByDate(1:string date, 2:i32 top, 3:string area),

	/// <summary>
	/// 根据id获取闪购活动
	/// </summary>
	TFlashSalesSummary GetFlashSalesListById(1:string settingId, 2:i32 offset, 3:i32 limit, 4:string area),
	
	// 获取指定分类顺序的闪购商品列表
	list<TFlashSalesCategoryProducts> GetFlashSalesListByCategories(1:string settingId,2:string area,3:list<string> categoryList),
	
}