package apierror

func ValidationError(err error) Error {
	return Error{err: err, status: 400, message: "リクエストに何らかの間違いがあります"}
}

func AuthorizationError() Error {
	return Error{status: 401, message: "ユーザの認証に失敗しました"}
}

func ClientDisconnectedError() Error {
	return Error{status: 499, message: "クライアントから接続が切断されました"}
}

func UnknownError(err error) Error {
	return Error{err: err, status: 500, message: "サーバ側で何らかのエラーが発生しました。時間を置いてから再度お試しください"}
}

func NotImplementedError() Error {
	return Error{status: 501, message: "この機能はまだ実装されていません"}
}

func DeadlineExceededError() Error {
	return Error{status: 504, message: "リクエストは規定時間内に処理されませんでした"}
}
