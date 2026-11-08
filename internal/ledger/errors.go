package ledger

import "errors"

var ErrAccessDenied = errors.New("access denied: account does not belong to user")
var ErrNoBalanceForTransfer = errors.New("no money for transfer")
var ErrNoUserAccounts = errors.New("user haven`t accounts")
