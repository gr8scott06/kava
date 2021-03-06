package paychan

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func EndBlocker(ctx sdk.Context, k Keeper) sdk.Tags {
	var err sdk.Error
	var channelTags sdk.Tags
	tags := sdk.EmptyTags()

	// Iterate through submittedUpdatesQueue
	// TODO optimise so it doesn't pull every channel update from DB every block
	q := k.getSubmittedUpdatesQueue(ctx)
	var sUpdate SubmittedUpdate
	var found bool

	for _, id := range q {
		// close the channel if the update has reached its execution time.
		// Using >= in case some are somehow missed.
		sUpdate, found = k.getSubmittedUpdate(ctx, id)
		if !found {
			panic("can't find element in queue that should exist")
		}
		if ctx.BlockHeight() >= sUpdate.ExecutionTime {
			k.removeFromSubmittedUpdatesQueue(ctx, sUpdate.ChannelID)
			channelTags, err = k.closeChannel(ctx, sUpdate.Update)
			if err != nil {
				panic(err)
			}
			tags.AppendTags(channelTags)
		}
	}
	return tags
}
