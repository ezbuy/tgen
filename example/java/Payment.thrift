
/// <summary>
/// implements by c#
/// </summary>

namespace * Payment
namespace java com.daigou.sg.rpc.payment
namespace csharp Zen.DataAccess.Payment
namespace javascript TRPC


include "Package.thrift"
include "Order.thrift"

struct TPaymentSummary {
	1:required double amountAvailable;			//可用余额
	2:required double rebateAmountAvailable;	//可用返利余额
	3:optional double amountPendingVerification;//等待确认金额
	4:optional double pendingWithdrawAmount;	//等待提现金额
	5:optional double unpaidAmount;				//未付款金额
}

struct TPackageDetail {
	1:required Package.TPackage tPackage;						//包裹信息
	2:required TPaymentBill paymentBill;		//账单信息
	3:required list<Order.TOrder> orders;				//订单列表
	4:optional list<Order.TOrderRemark> orderRemarks;	//备注
}

struct TPaymentBillSummary {
	1:required i32 id;					//付款id
	2:required string paymentNumber;	//付款号
	3:required string createDate;		//创建日期
	4:required double total;			//金额
	5:required string type;				//付款类型
	6:required string paymentType;		//付款显示类型
	7:optional i32 shipmentId;			//派送信息id
}

struct TPaymentBill {
	1:required i32 id;						//付款id
	2:required string paymentNumber;		//付款号
	3:required string total;				//金额
	4:required string chargeWeight;			//计费重
	5:required string packageWeight;		//重量
	6:required list<TPaymentBillCategory> paymentBillDetails	//付款明细项
	7:required string paymentStatus;		//付款状态
	8:required string createDate;			//创建时间
	9:required string paymentType;			//付款类型
}

struct TPaymentBillDetail {
	1:required TPaymentBill paymentBill;		//付款明细
	2:optional list<Package.TPackageDetail> packageInfo;	//包裹信息
	3:optional list<Order.TOrderBill> orderInfo;		//订单信息
}

struct TPaymentBillCategory {
	1:required string billCategoryCode;		//付款明细编号
	2:required string billCategoryName;		//英文名
	3:required string altBillCategoryName;	//中文名
	4:required string total;				//金额
}

struct TCreditCardInfo {
	1:required bool enableiPhoneCreditCard;		//是否开启iPhone充值
	2:required bool enableAndroidCreditCard;	//是否开启android充值
	3:required bool enableMobileWebCreditCard;	//是否开启mobile web充值
	4:required string disableUOBMsg;			//关闭提示信息
	5:required double creditCardFee;			//手续费
	6:required string creditCardDesc;           //文字说明
}

struct TPayParcelPaymentResult {
	1:required bool Result;		
	2:required string Message;	
	3:required string PaymentType;	
	4:required bool NeedTopUp;
}

struct TPrimePaymentSummary {
	1:required list<TPrimeType> primeTypes;		
	2:required string tip;	
	3:required string prepay;	
	4:required bool isPrimeMonthlyBought;
	5:required bool isPrimeYearlyBought;
}

struct TPrimeType {
	1:required string primeTypeName;		
	2:required string price;	
}

struct TPrimePaymentResult {
	1:required bool result;		
	2:required string message;	
	3:required string paymentType;	
	4:required bool needTopUp;	
	5:required list<string> paymentNumber; 
	6:required bool hasOtherUnpaid;	
	7:required list<i32> paymentId;
}

service Payment {
	
	/// <summary>
	/// 获取用户付款统计
	/// </summary>
	/// <returns>用户付款统计</returns>
	TPaymentSummary GetPaymentSummary(),

	/// <summary>
	/// 根据状态获取付款列表
	/// </summary>
	/// <param name="status">状态</param>
	/// <param name="offset">数据请求位置</param>
	/// <param name="limit">请求个数</param>
	/// <returns>付款列表</returns>
	list<TPaymentBillSummary> GetPaymentListByStatus(1:string status, 2:i32 offset, 3:i32 limit),

	/// <summary>
	/// 获取付款明细
	/// </summary>
	/// <param name="paymentId">付款id</param>
	/// <returns>付款明细</returns>
	TPaymentBillDetail GetPaymentDetail(1:i32 paymentId),

	/// <summary>
	/// 确认付款
	/// </summary>
	/// <param name="paymentIds">付款id</param>
	/// <returns>用户预付款余额</returns>
	double ConfirmPayments(1:list<i32> paymentIds),

	/// <summary>
	/// 获取用户预付款余额
	/// </summary>
	/// <returns>用户预付款余额</returns>
	double GetPrepayBalance(),

	/// <summary>
	/// 获取提现银行
	/// </summary>
	/// <returns>提现银行</returns>
	list<string> GetWithdrawBanks(),

	/// <summary>
	/// 提现
	/// </summary>
	/// <param name="bankName">提现银行</param>
	/// <param name="account">提现账号</param>
	/// <param name="amount">提现金额</param>
	/// <param name="reason">原因</param>
	/// <returns>提现结果</returns>
	string AddWithdrawReqeust(1:string bankName, 2:string account, 3:double amount, 4:string reason),

	/// <summary>
	/// 获取信用卡信息
	/// </summary>
	/// <returns>信用卡信息</returns>
	TCreditCardInfo GetCreditCardFee(),

	/// <summary>
	/// 充值
	/// </summary>
	/// <param name="transactionNumber">交易号</param>
	/// <param name="bankName">银行</param>
	/// <param name="telephone">电话</param>
	/// <param name="amount">金额</param>
	/// <param name="paymentMethod">付款方式</param>
	/// <param name="paymentIds">付款id</param>
	/// <param name="payDate">充值日期</param>
	/// <returns>充值结果</returns>
	string TopUp(1:string transactionNumber, 2:string bankName, 3:string telephone, 4:double amount, 5:string paymentMethod,
	 6:list<i32> paymentIds, 7:string payDate),

	/// <summary>
	/// 获取充值信息
	/// </summary>
	/// <returns>充值信息</returns>
	list<string> GetTopUpDescription(),

	/// <summary>
	/// 获取信用卡充值请求链接 implements by c#
	/// </summary>
	/// <param name="total">金额</param>
	/// <param name="creditCardFee">手续费</param>
	/// <param name="paymentIds">付款id</param>
	/// <param name="telephone">电话</param>
	/// <returns>信用卡充值请求链接</returns>
	string UserDoCreditCardTopUp(1:double total, 2:double creditCardFee, 3:list<string> paymentIds, 4:string telephone),
	
	/// <summary>
	/// 根据pending payment的包裹Id付款
	/// </summary>
	/// <param name="PaymentBillId">包裹Id</param>
	/// <returns>是否成功</returns>
	TPayParcelPaymentResult UserPayParcelPayment(1:list<i32> paymentBillIds),
	
	/// <summary>
	/// 获取Prime的支付内容
	/// </summary>
	/// <returns>Prime的支付内容</returns>
	TPrimePaymentSummary UserGetPrimePaymentSummary(),
	
	/// <summary>
	/// 根据Prime类型付款
	/// </summary>
	/// <param name="primeType">Prime类型</param>
	/// <returns>是否成功</returns>
	TPrimePaymentResult UserPayPrimePayment(1:string primeType),

	/// <summary>
	/// Prime会员续费
	/// </summary>
	/// <returns>是否成功</returns>
	TPrimePaymentResult UserRenewPrime(),
}
