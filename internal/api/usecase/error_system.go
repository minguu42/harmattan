package usecase

func ValidationError(err error) Error {
	return Error{err: err, status: 400, message: "リクエストに何らかの間違いがあります"}
}

func AuthorizationError(err error) Error {
	return Error{err: err, status: 401, message: "ユーザの認証に失敗しました"}
}

func UnknownError(err error) Error {
	return Error{err: err, status: 500, message: "サーバ側で何らかのエラーが発生しました。時間を置いてから再度お試しください"}
}

func PanicError(err error) Error {
	return Error{
		err:     err,
		status:  500,
		message: "サーバ側で何らかのエラーが発生しました。時間を置いてから再度お試しください",
	}
}

func NotImplementedError() Error {
	return Error{status: 501, message: "この機能はまだ実装されていません"}
}

func DeadlineExceededError(err error) Error {
	return Error{err: err, status: 504, message: "リクエストは規定時間内に処理されませんでした"}
}
