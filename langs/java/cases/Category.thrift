namespace * Category
namespace csharp Zen.DataAccess.Category
namespace java com.daigou.sg.rpc.category
namespace objc TRCategory
namespace javascript TRPC

struct TCategory {
	1: required i32 id;				//分类id
	2: required i64 pid;
	3: required string name;		// category名称
	4: required bool testbool;
	5: required string altName;		// 英文名称
}

struct TFloorCategory{
	1: required i32 id;				//分类id
	2: required string name;		// category名称
	3: required list<string> names;
	4: required list<TCategory> subCategories;
}

service Category {
	/// <summary>
	/// 获取热卖商品分类列表
	/// </summary>
	/// <param name="originCode">采购国家</param>
	/// <returns>热卖商品分类列表</returns>
	list<TCategory> GetTopLevelCategories(1:string originCode),

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
	list<TCategory> GetPrimeSubCategories(1:i32 categoryId)

}
