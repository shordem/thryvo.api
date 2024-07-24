package constants

const (
	// General Success Codes
	SuccessOperationCompleted  = 2000
	SuccessItemAddedToCart     = 2001
	SuccessItemRemovedFromCart = 2002
	SuccessCheckoutCompleted   = 2003
	SuccessPaymentProcessed    = 2004

	// General Client Errors
	ClientErrorBadRequest         = 4000
	ClientErrorUnauthorizedAccess = 4001
	ClientErrorPaymentRequired    = 4002
	ClientErrorForbidden          = 4003
	ClientErrorResourceNotFound   = 4004
	ClientRequestValidationError  = 4005
	ClientUnProcessableEntity     = 4006

	// General Server Errors
	ServerErrorInternal           = 5000
	ServerErrorServiceUnavailable = 5001
	ServerErrorGatewayTimeout     = 5002
	ServerErrorDatabase           = 5003
	ServerErrorExternalService    = 5004

	// User Account and Authentication
	UserRegistrationSuccessful  = 4100
	UserLoginSuccessful         = 4101
	UserLogoutSuccessful        = 4102
	PasswordResetSuccessful     = 4103
	UserNotFound                = 4104
	AccountLocked               = 4105
	AccountVerificationRequired = 4106
	EmailAlreadyInUse           = 4107
	InvalidEmailFormat          = 4108
	WeakPassword                = 4109
	InvalidCredentials          = 4110

	// Shopping Cart and Orders
	CartUpdatedSuccessfully         = 4200
	CartIsEmpty                     = 4201
	ItemOutOfStock                  = 4202
	ItemQuantityUpdatedSuccessfully = 4203
	InvalidItemID                   = 4204
	OrderPlacedSuccessfully         = 4205
	OrderCancelledSuccessfully      = 4206
	OrderAlreadyShipped             = 4207
	InvalidOrderID                  = 4208
	OrderNotFound                   = 4209
	MinimumOrderAmountNotMet        = 4210

	// Payment Processing
	PaymentAuthorized           = 4300
	PaymentDeclined             = 4301
	PaymentRefundedSuccessfully = 4302
	InvalidPaymentMethod        = 4303
	PaymentGatewayError         = 4304
	InsufficientFunds           = 4305
	CurrencyNotSupported        = 4306
	PaymentPending              = 4307

	// Inventory Management
	InventoryUpdatedSuccessfully     = 4400
	InventoryLevelLow                = 4401
	InventoryLevelHigh               = 4402
	InvalidInventoryData             = 4403
	InventoryItemNotFound            = 4404
	InventoryItemAddedSuccessfully   = 4405
	InventoryItemRemovedSuccessfully = 4406

	// Shipping and Delivery
	ShippingAddressAddedSuccessfully   = 4500
	ShippingAddressUpdatedSuccessfully = 4501
	InvalidShippingAddress             = 4502
	ShippingMethodNotAvailable         = 4503
	DeliveryScheduled                  = 4504
	DeliveryInProgress                 = 4505
	DeliveryCompleted                  = 4506
	DeliveryDelayed                    = 4507
	TrackingNumberNotFound             = 4508
)
