package usecase

import (
	"context"
	"fmt"
	"sensibull/stocks-api/business/entities/core"
)

// should run every 15min
func (is *instrumentservice) UpdateEquityStockDetails(ctx context.Context) error {
	// fetch the equity instrument list
	underlyingsEQ, err := is.httpir.GetUnderLyings(ctx)
	if err != nil {
		// retry mechanism here
		return err
	}
	prevTokens, err := is.db.GetTokensAgainstToken(ctx, TOKENFORALLUNDERLYING, EQUITY)
	if err != nil {
		// retry mechanism here
		return err
	}
	go is.updateTokenSet(ctx, TOKENFORALLUNDERLYING, EQUITY, underlyingsEQ, &prevTokens)
	for _, val := range underlyingsEQ {
		if !prevTokens.Set.Contains(val.Token) {
			// subscribe to websocket for this instrument
			err = is.websocket.Subscribe(ctx, []int64{val.Token})
			if err != nil {
				//log
				// retry
				continue
			}
			// create entry
			err = is.db.UpsertInstrument(ctx, val)
			if err != nil {
				//log
				//retry
				continue
			}
		}
	}
	return nil
}

// should run every 1min
func (is *instrumentservice) UpdateDerivativeStockDetails(ctx context.Context) error {
	curEqTokens, err := is.db.GetTokensAgainstToken(ctx, TOKENFORALLUNDERLYING, EQUITY)
	if err != nil {
		// retry mechanism here
		return err
	}
	for _, val := range curEqTokens.Set.Values() {
		go is.updateDerivativeStockDetail(ctx, val.(int64))
	}
	return nil
}

func (is *instrumentservice) updateDerivativeStockDetail(ctx context.Context, underlyingToken int64) error {
	// fetch the derivatives instrument list
	underlyingsDvts, err := is.httpir.GetUnderLyingDerivatives(ctx, underlyingToken)
	if err != nil {
		// retry mechanism here
		return err
	}
	prevDvtsTokens, err := is.db.GetTokensAgainstToken(ctx, fmt.Sprint(underlyingToken), DERIVATIVES)
	if err != nil {
		// retry mechanism here
		return err
	}
	go is.updateTokenSet(ctx, fmt.Sprint(underlyingToken), DERIVATIVES, underlyingsDvts, &prevDvtsTokens)
	for _, val := range underlyingsDvts {
		if !prevDvtsTokens.Set.Contains(val.Token) {
			// subscribe to websocket for this instrument
			err = is.websocket.Subscribe(ctx, []int64{val.Token})
			if err != nil {
				//log
				// retry
				continue
			}
			// create entry
			err = is.db.UpsertInstrument(ctx, val)
			if err != nil {
				//log
				//retry
				continue
			}
		}
	}
	return nil
}

func (is *instrumentservice) UpdateInstrumentPrice(ctx context.Context) error {
	return nil
}

func (is *instrumentservice) updateTokenSet(ctx context.Context, itoken, itype string, instrumnets []core.Instrument, prevTokens *core.Tokens) {
	currentTokens := core.NewTokenSet()

	for _, val := range instrumnets {
		currentTokens.Set.Add(val.Token)
	}
	err := is.db.SaveTokensAgainstToken(ctx, itoken, itype, currentTokens)
	if err != nil {
		//retry mechanism here
		//log the error
	}
	listOfTokensToUnsubscribe := []int64{}
	for _, token := range prevTokens.Set.Values() {
		if !currentTokens.Set.Contains(token) {
			err := is.db.DeleteInstrument(ctx, core.Instrument{Token: token.(int64)})
			if err != nil {
				// log and continue
			}
			listOfTokensToUnsubscribe = append(listOfTokensToUnsubscribe, token.(int64))
		}
	}
	err = is.websocket.UnSubscribe(ctx, listOfTokensToUnsubscribe)
	if err != nil {
		//retry mechanism here
		//log the error
	}
}
