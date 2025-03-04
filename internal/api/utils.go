package api

import (
	"errors"
	"net/http"

	"game-fun-be/internal/model"
	"game-fun-be/internal/response"

	"github.com/gin-gonic/gin"
)

// GetUserIDFromContext 从上下文中获取用户 ID
func GetUserIDFromContext(c *gin.Context) (string, *response.Response) {
	userID, exists := c.Get("user_id")
	if !exists {
		errResp := response.Err(http.StatusUnauthorized, "Address not found in context", nil)
		return "", &errResp
	}

	userIDStr, ok := userID.(string)
	if !ok {
		errResp := response.Err(http.StatusUnauthorized, "Invalid address type in context", nil)
		return "", &errResp
	}

	return userIDStr, nil
}

// GetAddressFromContext 从上下文中获取地址
func GetAddressFromContext(c *gin.Context) (string, *response.Response) {
	address, exists := c.Get("address")
	if !exists {
		errResp := response.Err(http.StatusUnauthorized, "Address not found in context", nil)
		return "", &errResp
	}

	addressStr, ok := address.(string)
	if !ok {
		errResp := response.Err(http.StatusUnauthorized, "Invalid address type in context", nil)
		return "", &errResp
	}

	return addressStr, nil
}

// GetPageAndLimit 从上下文中获取并验证 page 和 limit 参数
func GetPageAndLimit(c *gin.Context) (page, limit string, errResp *response.Response) {
	page = c.Query("page")
	if page == "" {
		errResp = &response.Response{
			Code: http.StatusBadRequest,
			Msg:  "page parameter is required",
		}
		return
	}

	limit = c.Query("limit")
	if limit == "" {
		errResp = &response.Response{
			Code: http.StatusBadRequest,
			Msg:  "limit parameter is required",
		}
		return
	}

	return page, limit, nil
}

// GetPageAndLimit 从上下文中获取并验证 page 和 limit 参数
func GetLimit(c *gin.Context) (limit string, errResp *response.Response) {

	limit = c.Query("limit")
	if limit == "" {
		errResp = &response.Response{
			Code: http.StatusBadRequest,
			Msg:  "limit parameter is required",
		}
		return
	}

	return limit, nil
}

// ParseChainTypeWithResponse 解析并验证 chain_type 参数，并返回 HTTP 响应
func ParseChainTypeWithResponse(c *gin.Context) (model.ChainType, *response.Response) {
	chainType := c.Param("chain_type")
	if chainType == "" {
		errResp := response.Err(http.StatusBadRequest, "chain_type cannot be empty", errors.New("chain_type is required"))
		return model.ChainTypeUnknown, &errResp
	}
	return model.ChainTypeFromString(chainType), nil
}
