package redis

import (
	"github.com/gomodule/redigo/redis"
	"testing"
	"time"
)

func TestStore_Set(t *testing.T) {
	type fields struct {
		pool    *redis.Pool
		timeout time.Duration
	}
	type args struct {
		key   string
		value string
	}

	tests := []struct {
		name    string
		fields  *Store
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{name: "set values string", fields: DialRedisPoolBlocked("localhost:6379", 1, 20, time.Minute*3), args: args{key: "key1-test-string", value: "key1-values"}, wantErr: false},
		{name: "set values int", fields: DialRedisPoolBlocked("localhost:6379", 1, 20, time.Minute*3), args: args{key: "key1-test-int", value: "key1-values"}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Store{
				pool:    tt.fields.pool,
				timeout: tt.fields.timeout,
			}
			if err := s.Set(tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStore_Get(t *testing.T) {
	type fields struct {
		pool    *redis.Pool
		timeout time.Duration
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  *Store
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{name: "get values string", fields: DialRedisPoolBlocked("localhost:6379", 1, 20, time.Minute*3), args: args{key: "key1-test-string"}, want: "key1-values", wantErr: false},
		{name: "get values int", fields: DialRedisPoolBlocked("localhost:6379", 1, 20, time.Minute*3), args: args{key: "key1-test-int"}, want: "key1-values", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Store{
				pool:    tt.fields.pool,
				timeout: tt.fields.timeout,
			}
			got, err := s.Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}
