namespace * Category
namespace csharp Zen.DataAccess.Category
namespace java com.daigou.sg.rpc.category
namespace objc TRCategory
namespace javascript TRPC

// Todo: duplicated with Collection.thrift
struct TProductSimple {
	2: required string url;			// 商品url
	3: required string name;		// 商品名称
	4: required string price;		// 商品价格
	5: required string picture;		// 商品图片
	6: required i32 favouriteCount;	// 收藏个数
	7: required string priceWithSymbol;		// 商品价格
	8: required string originCode;	//采购国家
	9: required string vendorName;	//卖家名字
	10: required string altProductName;		//商品英文名称
	11: required string customerLocalPrice;	//顾客本地货币商品价格，带货币符号
	12: required bool isEzBuy;	//是否是EzBuy商品
}

struct TCategory {
	1: required i32 id;				//分类id
	2: required i32 parentId;
	3: required string name;		// category名称
	4: required string picture;		// 图片
	5: required string altName;		// 英文名称
}

struct TFloorCategory{
	1: required i32 id;				//分类id
	2: required i32 parentId;
	3: required string name;		// category名称
	4: required string picture;		// 图片
	5: required string altName;		// 英文名称
	6: required list<TCategory> subCategories;
}

struct TRecentPrimeCustomer{
	1: required string customerName;
	2: required list<string> productImage;
	3: required string avatar;
	4: required i32 itemCount;
}

service Category {
	/// <summary>
	/// 获取首页热卖分类列表
	/// </summary>
	/// <param name="offset">请求位置</param>
	/// <param name="limit">请求个数</param>
	/// <param name="originCode">采购国家</param>
	/// <returns>首页热卖分类列表</returns>
	list<TCategory> GetHomePageCategories(1:i32 offset, 2:i32 limit, 3:string originCode),

	/// <summary>
	/// 获取热卖商品分类列表
	/// </summary>
	/// <param name="originCode">采购国家</param>
	/// <returns>热卖商品分类列表</returns>
	list<TCategory> GetTopLevelCategories(1:string originCode),

	/// <summary>
	/// 获取热卖商品子分类列表
	/// </summary>
	/// <param name="categoryId">父分类id</param>
	/// <returns>热卖商品子分类列表</returns>
	list<TCategory> GetSubCategories(1:i32 categoryId),

	/// <summary>
	/// 根据分类获取热卖商品列表
	/// </summary>
	/// <param name="id">分类id</param>
	/// <param name="offset">数据的起始位置</param>
	/// <param name="limit">一次请求要获取的个数</param>
	/// <param name="originCode">采购国家</param>
	/// <returns>热卖商品列表</returns>
	list<TProductSimple> GetProducts(1:i32 id, 2:i32 offset, 3:i32 limit, 4:string originCode),

	/// <summary>
	/// 根据分类获取热卖商品数量
	/// </summary>
	/// <param name="id">分类id</param>
	/// <param name="originCode">采购国家</param>
	/// <returns>热卖商品数量</returns>
	i32 GetAllProductCount(1:i32 id, 2:string originCode),

	// <summary>
	/// 搜索热卖商品列表
	/// </summary>
	/// <param name="keyword">关键字</param>
	/// <param name="offset">数据的起始位置</param>
	/// <param name="limit">一次请求要获取的个数</param>
	/// <param name="categoryId">分类id</param>
	/// <param name="originCode">采购国家</param>
	/// <returns>热卖商品列表</returns>
	list<TProductSimple> SearchCategoryProducts(1:string keyword, 2:i32 offset, 3:i32 limit, 4:i32 categoryId, 5:string originCode),

	/// <summary>
	/// 根据分类获取热卖prime商品列表
	/// </summary>
	/// <param name="id">分类id</param>
	/// <param name="offset">数据的起始位置</param>
	/// <param name="limit">一次请求要获取的个数</param>
	/// <returns>prime热卖商品列表</returns>
	list<TProductSimple> GetPrimeProducts(1:i32 id, 2:i32 offset, 3:i32 limit),

	// <summary>
	/// 获取晒单Summary列表
	/// </summary>
	/// <returns>晒单Summary列表</returns>
	list<TRecentPrimeCustomer> UserGetRecentPrimePurchaseSummaryList(1:i32 offset, 2:i32 limit),

	// <summary>
	/// 获取晒单用户列表
	/// </summary>
	/// <returns>晒单Summary列表</returns>
	list<TRecentPrimePurchase> UserGetRecentPrimePurchaseList(1:i32 offset, 2:i32 limit),

	// <summary>
	/// 获取晒单用户详情
	/// </summary>
	/// <returns>晒单用户详情</returns>
	TRecentPrimePurchase UserGetRecentPrimePurchaseDetail(1:i32 paymentBillId, 2:i32 offset, 3:i32 limit),


	/// <summary>
	/// 获取Pirme热卖商品分类列表
	/// </summary>
	/// <returns>热卖商品分类列表</returns>
	list<TCategory> GetPrimeTopLevelCategories(),

	/// <summary>
	/// 获取Prime热卖商品子分类列表
	/// </summary>
	/// <param name="categoryId">父分类id</param>
	/// <returns>热卖商品子分类列表</returns>
	list<TCategory> GetPrimeSubCategories(1:i32 categoryId),

	/// <summary>
	/// 获取Prime楼层分类列表
	/// </summary>
	/// <returns>楼层分类列表</returns>
	list<TFloorCategory> GetPrimeFloorCategories(),

}
