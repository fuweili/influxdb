package bolt_test

import (
	"context"
	"testing"

	"github.com/influxdata/influxdb/kv"
)

type KVStoreFields struct {
	Bucket []byte
	Pairs  []kv.Pair
}

func initKVStore(f KVStoreFields, t *testing.T) (kv.Store, func()) {
	s, closeFn, err := NewTestKVStore(t)
	if err != nil {
		t.Fatalf("failed to create new kv store: %v", err)
	}

	mustCreateBucket(t, s, f.Bucket)

	err = s.Update(context.Background(), func(tx kv.Tx) error {
		b, err := tx.Bucket(f.Bucket)
		if err != nil {
			return err
		}

		for _, p := range f.Pairs {
			if err := b.Put(p.Key, p.Value); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		t.Fatalf("failed to put keys: %v", err)
	}
	return s, func() {
		closeFn()
	}
}

func TestKVStore(t *testing.T) {
	kv.KVStore(initKVStore, t)
}

func mustCreateBucket(t testing.TB, store kv.SchemaStore, bucket []byte) {
	t.Helper()

	// migrationName := fmt.Sprintf("create bucket %q", string(bucket))

	// if err := migration.CreateBuckets(migrationName, bucket).Up(context.Background(), store); err != nil {
	// 	t.Fatal(err)
	// }
}
