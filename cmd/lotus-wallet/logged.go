package main

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"github.com/ipfs/go-cid"
	"golang.org/x/xerrors"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/crypto"

	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/chain/types"
)

type LoggedWallet struct {
	under api.Wallet
}

func (c *LoggedWallet) WalletNew(ctx context.Context, typ types.KeyType) (address.Address, error) {
	clientIp := ctx.Value("client_ip")

	log.Infow("WalletNew", "client", clientIp, "type", typ)

	//return address.Address{}, xerrors.Errorf("unsupported method")
	return c.under.WalletNew(ctx, typ)
}

func (c *LoggedWallet) WalletHas(ctx context.Context, addr address.Address) (bool, error) {
	clientIp := ctx.Value("client_ip")

	log.Infow("WalletHas", "address", addr, "client", clientIp)

	return c.under.WalletHas(ctx, addr)
}

func (c *LoggedWallet) WalletList(ctx context.Context) ([]address.Address, error) {
	clientIp := ctx.Value("client_ip")

	log.Infow("WalletList", "client", clientIp)

	return c.under.WalletList(ctx)
}

func (c *LoggedWallet) WalletSign(ctx context.Context, k address.Address, msg []byte, meta api.MsgMeta) (*crypto.Signature, error) {
	clientIp := ctx.Value("client_ip")
	uuid := ctx.Value("uuid")

	log.Infow("WalletSign",
		"uuid", uuid,
		"client", clientIp,
		"address", k,
		"type", meta.Type,
		"status", "start",
	)

	switch meta.Type {
	case api.MTChainMsg:
		var cmsg types.Message
		if err := cmsg.UnmarshalCBOR(bytes.NewReader(meta.Extra)); err != nil {
			return nil, xerrors.Errorf("unmarshalling message: %w", err)
		}

		_, bc, err := cid.CidFromBytes(msg)
		if err != nil {
			return nil, xerrors.Errorf("getting cid from signing bytes: %w", err)
		}

		if !cmsg.Cid().Equals(bc) {
			return nil, xerrors.Errorf("cid(meta.Extra).bytes() != msg")
		}

		log.Infow("WalletSign",
			"uuid", uuid,
			"cid", bc,
			"from", cmsg.From,
			"to", cmsg.To,
			"value", types.FIL(cmsg.Value),
			"feecap", types.FIL(cmsg.RequiredFunds()),
			"method", cmsg.Method,
			"params", hex.EncodeToString(cmsg.Params),
		)

		if !checkMethod(cmsg.Method) {
			log.Errorw("WalletSign",
				"uuid", uuid,
				"error", "unsupported method",
				"method", cmsg.Method,
			)
			return nil, xerrors.Errorf("unsupported method")
		}

		if !checkAddress(cmsg.From.String(), cmsg.To.String()) {
			log.Errorw("WalletSign",
				"uuid", uuid,
				"error", "address does not match",
				"from", cmsg.From,
				"to", cmsg.To,
			)
			return nil, xerrors.Errorf("address does not match")
		}
	case api.MTBlock:
		_, err := types.DecodeBlock(msg)
		if err != nil {
			log.Errorw("WalletSign",
				"uuid", uuid,
				"error", fmt.Sprintf("decode block failed, %v", err),
				"msg", msg)
			return nil, xerrors.Errorf("parsing block header error: %w", err)
		}
	default:
		_, err := cid.Parse(msg)
		if err == nil {
			log.Errorw("WalletSign",
				"uuid", uuid,
				"error", "other message types should not sign cid",
				"msg", msg)
			return nil, xerrors.Errorf("other message types should not sign cid")
		}
	}

	log.Infow("WalletSign",
		"uuid", uuid,
		"status", "finish",
	)

	return c.under.WalletSign(ctx, k, msg, meta)
}

func (c *LoggedWallet) WalletExport(ctx context.Context, a address.Address) (*types.KeyInfo, error) {
	clientIp := ctx.Value("client_ip")

	log.Infow("WalletExport", "client", clientIp, "address", a)

	return nil, xerrors.Errorf("unsupported method")
	//return c.under.WalletExport(ctx, a)
}

func (c *LoggedWallet) WalletImport(ctx context.Context, ki *types.KeyInfo) (address.Address, error) {
	clientIp := ctx.Value("client_ip")

	log.Infow("WalletImport", "client", clientIp, "type", ki.Type)

	//return address.Address{}, xerrors.Errorf("unsupported method")
	return c.under.WalletImport(ctx, ki)
}

func (c *LoggedWallet) WalletDelete(ctx context.Context, addr address.Address) error {
	clientIp := ctx.Value("client_ip")

	log.Infow("WalletDelete", "client", clientIp, "address", addr)

	return xerrors.Errorf("unsupported method")
	//return c.under.WalletDelete(ctx, addr)
}
