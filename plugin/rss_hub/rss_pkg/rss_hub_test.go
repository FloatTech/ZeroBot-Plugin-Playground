package rss_pkg

import (
	"context"
	"encoding/json"
	"testing"
)

func TestNewRssDomain(t *testing.T) {
	dm, err := newRssDomain("")
	if err != nil {
		t.Fatal(err)
		return
	}
	if dm == nil {
		t.Fatal("domain is nil")
	}
}

//var testRssHubChannelUrl = "https://rsshub.rssforever.com/bangumi/tv/calendar/today"

var dm *rssDomain

func TestSub(t *testing.T) {
	testCases := []struct {
		name     string
		feedLink string
		gid      int64
	}{
		{
			name:     "test1",
			feedLink: "/bangumi/tv/calendar/today",
			gid:      99,
		},
		{
			name:     "test2",
			feedLink: "/go-weekly",
			gid:      99,
		},
		{
			name:     "test3",
			feedLink: "/go-weekly",
			gid:      123,
		},
	}
	dm, _ = newRssDomain("")
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			channel, ifExisted, ifSub, err := dm.Subscribe(ctx, tc.gid, tc.feedLink)
			if err != nil {
				t.Fatal(err)
				return
			}
			t.Logf("[TEST] add sub res: %+v,%+v,%+v\n", channel, ifExisted, ifSub)
			res, ext, err := dm.storage.GetIfExistedSubscribe(ctx, tc.gid, tc.feedLink)
			if err != nil {
				t.Fatal(err)
				return
			}
			t.Logf("[TEST] if exist: %+v,%+v", res, ext)
			channels, err := dm.GetSubscribedChannelsByGroupId(ctx, 2)
			if err != nil {
				t.Fatal(err)
				return
			}
			t.Logf("[TEST] 2 channels: %+v", channels)
		})
	}
}

func TestSub_2(t *testing.T) {
	dm, _ = newRssDomain("")
	ctx := context.Background()
	channel, ifExisted, ifSub, err := dm.Subscribe(ctx, 99, "/bangumi/tv/calendar/today")
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Logf("add: %+v,%+v,%+v\n", channel, ifExisted, ifSub)
	res, ext, err := dm.storage.GetIfExistedSubscribe(ctx, 99, "/bangumi/tv/calendar/today")
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Logf("if exist: %+v,%+v", res, ext)
	channels, err := dm.GetSubscribedChannelsByGroupId(ctx, 2)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Logf("2 channels: %+v", channels)

	err = dm.Unsubscribe(ctx, 2, "/bangumi/tv/calendar/today")
	if err != nil {
		t.Fatal(err)
		return
	}
	res, ext, err = dm.storage.GetIfExistedSubscribe(ctx, 2, "/bangumi/tv/calendar/today")
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Logf("after del: %+v,%+v", res, ext)
}

func Test_SyncFeed(t *testing.T) {
	//t.Run("test1", TestSub)
	dm, _ = newRssDomain("")
	feed, err := dm.SyncJobTrigger(context.Background())
	if err != nil {
		return
	}
	rs, _ := json.Marshal(feed)
	t.Logf("[Test] feed: %+v", string(rs))
}
