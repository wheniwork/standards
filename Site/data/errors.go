package data
import (
	"github.com/kataras/iris"
)

type DError struct {
	Message string
	ServerError bool
	ClientError bool
	PermissionError bool
	NotFoundError bool
	ActualError error
}

type DErrorAPIResponse struct {
	Message string `json:"message,omitempty"`
	Success bool `json:"success"`
}

func NewClientError(Message string, err error) *DError {
	return &DError{
		Message:Message,
		ClientError:true,
		ActualError:err,
	}
}

func NewNotFoundError(Message string) *DError {
	return &DError{
		Message:Message,
		NotFoundError:true,
	}
}

func NewServerError(Message string, err error) *DError {
	return &DError{
		Message:Message,
		ServerError:true,
		ActualError:err,
	}
}

func NewPermissionError(Message string) *DError {
	return &DError{
		Message:Message,
		PermissionError:true,
	}
}

func ErrorResponse(ctx iris.Context, err *DError) {
	if err.PermissionError {
		ctx.StatusCode(403)
	} else if err.ServerError {
		ctx.StatusCode(500)
	} else if err.ClientError {
		ctx.StatusCode(400)
	} else if err.NotFoundError {
		ctx.StatusCode(404)
	}
	ctx.JSON(DErrorAPIResponse{
		Message:err.Message,
		Success:false,
	})
}

