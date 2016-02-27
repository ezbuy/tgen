namespace go webhomepage.abc
namespace java com.daigou.sg.rpc.collection
namespace objc TRCollection
namespace javascript TRPC

struct TWebHomepageBanner {
    1:required i32 bannerType;
    2:required string picture;
    3:optional string linkAddress;
}

struct TWebHomepageCampaign {
    1:required string picture;
    2:optional string linkAddress;
}

struct TWebHomepageCategory {
    1:required i32 id;
    2:required string name;
    3:required string image;
    4:optional list<TWebHomepageCategory> subCategories;
}

struct TWebHomepageProduct{
    1:required i32 id;
    2:required string name;
    3:required string url;
    4:required double price;
    5:required string image;
    6:required string originCode;
}

struct TWebHomepageCategoryCollection{
    1:required i32 id;
    2:required TWebHomepageCampaign campaign;
    3:required TWebHomepageCategory firstLevelCategory;
    4:required list<TWebHomepageCategory> secondLevelCategories;
    5:required list<TWebHomepageCategory> thirdLevelCategories;
}

struct TWebHomepagePromotion{
    1:required i32 topCategoryId;
    2:required TWebHomepageCampaign promotionBanner;
    3:required list<TWebHomepageProduct> products;
}

service WebHomepage{

    /// <summary>
    /// 获取合作卖家商品
    /// </summary>
    /// <param name="originCode">采购国家</param>
    /// <returns>合作卖家商品</returns>
    list<TWebHomepageProduct> GetSellerProducts(1:string originCode),

    
    /// <summary>
    /// 获取首页所需的所有分类集合
    /// </summary>
    /// <returns>合作卖家商品</returns>
    list<TWebHomepageCategoryCollection> GetCategoryCollections(),    


    /// <summary>
    /// 获取首页的Banner列表
    /// </summary>
    /// <returns>Banner列表</returns>
    list<TWebHomepageBanner> GetBannerList(),



    /// <summary>
    /// 获取促销信息和改分类下的商品信息
    /// </summary>
    /// <param name="topCategoryId">顶级分类ID</param>
    /// <returns>优惠和产品信息</returns>
    TWebHomepagePromotion GetMenuItemByTopCategoryId(1:i32 topCategoryId),

}
