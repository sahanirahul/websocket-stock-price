package usecase

import (
	"context"
	"errors"
	"fmt"
	"sensibull/stocks-api/business/entities/core"
	"sensibull/stocks-api/business/utility"
	"sensibull/stocks-api/utils/logging"
	"sync"
)

// should run every 15min
func (is *instrumentservice) UpdateEquityStockDetails(ctx context.Context) (bool, error) {
	// fetch the equity instrument list
	shouldRetry := true
	underlyingsEQ, err := is.httpir.GetUnderLyings(ctx)
	if err != nil {
		logging.Logger.WriteLogs(ctx, "error_fetching_underlying_http", logging.ErrorLevel, logging.Fields{"error": err})
		return shouldRetry, err
	}
	if len(underlyingsEQ) == 0 {
		// could  not fetch derivative details
		return shouldRetry, errors.New("empty_underlying_response")
	}
	prevTokens, err := is.db.GetTokensAgainstToken(ctx, utility.TOKENFORALLUNDERLYING, utility.EQUITY)
	if err != nil {
		logging.Logger.WriteLogs(ctx, "error_getting_mapped_tokens_equity", logging.ErrorLevel, logging.Fields{"error": err})
		return shouldRetry, err
	}
	err = is.updateTokenSet(ctx, utility.TOKENFORALLUNDERLYING, utility.EQUITY, underlyingsEQ, &prevTokens)
	if err != nil {
		return shouldRetry, err
	}
	for _, val := range underlyingsEQ {
		if prevTokens.Set == nil || prevTokens.Set.Size() == 0 || !prevTokens.Set.Contains(val.Token) {
			// subscribe to websocket for this instrument
			is.websocket.AddSubscriptionRequest(core.WebsocketSubscription{MessageCommand: utility.MsgCommandSubscribe, DataType: utility.DataTypeQuote, Tokens: []int64{val.Token}})
			// create entry
			err = is.db.InsertInstrument(ctx, val)
			if err != nil {
				logging.Logger.WriteLogs(ctx, "error_insert_equity_instruments", logging.ErrorLevel, logging.Fields{"error": err})
				continue
			}
		}
	}
	return !shouldRetry, nil
}

// should run every 1min
func (is *instrumentservice) UpdateDerivativeStockDetails(ctx context.Context) error {
	curEqTokens, err := is.db.GetTokensAgainstToken(ctx, utility.TOKENFORALLUNDERLYING, utility.EQUITY)
	if err != nil {
		logging.Logger.WriteLogs(ctx, "error_getting_mapped_tokens_equity_derivative", logging.ErrorLevel, logging.Fields{"error": err})
		return err
	}
	retryCountMap := sync.Map{}
	var wg sync.WaitGroup
	for _, val := range curEqTokens.Set.Values() {
		token := int64(val.(float64))
		wg.Add(1)
		is.pool.AddJob(core.NewJob(func() {
			defer wg.Done()
			retryCountMap.Store(token, 1)
			for {
				shouldretry, err := is.updateDerivativeStockDetail(ctx, token)
				if err != nil {
					logging.Logger.WriteLogs(ctx, "error_updating_derivatives", logging.ErrorLevel, logging.Fields{"error": err, "token": token})
				}
				if shouldretry {
					val, ok := retryCountMap.Load(token)
					count := val.(int64)
					if ok && count > 3 {
						return
					}
					if !ok {
						count = 1
					}
					retryCountMap.Store(token, count+1)
				} else {
					return
				}
			}

		}))
	}
	wg.Wait()
	return nil
}

func (is *instrumentservice) updateDerivativeStockDetail(ctx context.Context, underlyingToken int64) (bool, error) {
	shouldRetry := true
	// fetch the derivatives instrument list
	underlyingsDvts, err := is.httpir.GetUnderLyingDerivatives(ctx, underlyingToken)
	if err != nil {
		logging.Logger.WriteLogs(ctx, "error_fetching_underlying_derivative_http", logging.ErrorLevel, logging.Fields{"error": err, "underlying_token": underlyingToken})
		return shouldRetry, err
	}
	if len(underlyingsDvts) == 0 {
		// could  not fetch derivative details
		return shouldRetry, errors.New("empty derivative response")
	}
	prevDvtsTokens, err := is.db.GetTokensAgainstToken(ctx, fmt.Sprint(underlyingToken), utility.DERIVATIVES)
	if err != nil {
		logging.Logger.WriteLogs(ctx, "error_fetching_mapped_tokens_for_derivatives", logging.ErrorLevel, logging.Fields{"error": err, "underlying_token": underlyingToken})
		return shouldRetry, err
	}
	err = is.updateTokenSet(ctx, fmt.Sprint(underlyingToken), utility.DERIVATIVES, underlyingsDvts, &prevDvtsTokens)
	if err != nil {
		return shouldRetry, err
	}
	for _, val := range underlyingsDvts {
		if prevDvtsTokens.Set == nil || prevDvtsTokens.Set.Size() == 0 || !prevDvtsTokens.Set.Contains(val.Token) {
			// subscribe to websocket for this instrument
			is.websocket.AddSubscriptionRequest(core.WebsocketSubscription{MessageCommand: utility.MsgCommandSubscribe, DataType: utility.DataTypeQuote, Tokens: []int64{val.Token}})
			// create entry in db
			err = is.db.InsertInstrument(ctx, val)
			if err != nil {
				logging.Logger.WriteLogs(ctx, "insert-instrument-failed", logging.ErrorLevel, logging.Fields{"error": err, "instrument": val})
				continue
			}
		}
	}
	return !shouldRetry, nil
}

func (is *instrumentservice) updateTokenSet(ctx context.Context, itoken, itype string, instrumnets []core.Instrument, prevTokens *core.Tokens) error {
	currentTokens := core.NewTokenSet()
	for _, val := range instrumnets {
		currentTokens.Set.Add(val.Token)
	}
	logging.Logger.WriteLogs(ctx, "saving-tokens-against-token", logging.InfoLevel, logging.Fields{"currentTokens": currentTokens})
	err := is.db.SaveTokensAgainstToken(ctx, itoken, itype, currentTokens)
	if err != nil {
		logging.Logger.WriteLogs(ctx, "error-saving-tokens-against-token", logging.ErrorLevel, logging.Fields{"error": err, "currentTokens": currentTokens})
		return err
	}
	listOfTokensToUnsubscribe := []int64{}
	if prevTokens.Set == nil || prevTokens.Set.Size() == 0 {
		return nil
	}
	for _, token := range prevTokens.Set.Values() {
		if !currentTokens.Set.Contains(token) {
			t := int64(token.(float64))
			err := is.db.DeleteInstrument(ctx, core.Instrument{Token: t})
			if err != nil {
				logging.Logger.WriteLogs(ctx, "instrument-delete-failed", logging.ErrorLevel, logging.Fields{"error": err, "token": t})
			}
			listOfTokensToUnsubscribe = append(listOfTokensToUnsubscribe, t)
		}
	}
	is.websocket.AddSubscriptionRequest(core.WebsocketSubscription{MessageCommand: utility.MsgCommandUnSubscribe, DataType: utility.DataTypeQuote, Tokens: listOfTokensToUnsubscribe})
	return nil
}
