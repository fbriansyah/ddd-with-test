package domain

import (
	"context"
	"ddd-with-test/mocks"
	"testing"

	"go.uber.org/mock/gomock"
)

func Test_CookiePurchase(t *testing.T) {
	t.Run(`Given a user tries to purchase a cookie and we have them in stock, 
		"when they tap their card, they get charged and then receive an email receipt a few moments later.`,
		func(t *testing.T) {
			var (
				ctrl = gomock.NewController(t)
				e    = mocks.NewMockEmailSender(ctrl)
				c    = mocks.NewMockCardCharger(ctrl)
				s    = mocks.NewMockCookieStockChecker(ctrl)
				ctx  = context.Background()
			)
			cookiesToBuy := 5
			totalExpectedCost := 250
			cs, err := NewCookieService(e, c, s)
			if err != nil {
				t.Fatalf("expected no error but got %v", err)
			}
			gomock.InOrder(
				s.EXPECT().AmountInStock(ctx).Times(1).Return(cookiesToBuy),
				c.EXPECT().ChargeCard(ctx, "some-token", totalExpectedCost).Times(1).Return(nil),
				e.EXPECT().SendEmailReceipt(ctx, "some-email").Times(1).Return(nil),
			)
			err = cs.PurchaseCookies(ctx, cookiesToBuy, "some-token", "some-email")
			if err != nil {
				t.Fatalf("expected no error but got %v", err)
			}
		},
	)
}
